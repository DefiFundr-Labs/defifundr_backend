package mail

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	"gopkg.in/gomail.v2"
)

// EmailWorker processes emails from the queue and sends them
type EmailWorker struct {
	config       config.Config
	logger       logging.Logger
	dialer       *gomail.Dialer
	templates    map[string]*template.Template
	templatesDir string
	sender       *AsyncQEmailSender
}

// NewEmailWorker creates a new email worker
func NewEmailWorker(config config.Config, logger logging.Logger, sender *AsyncQEmailSender) (*EmailWorker, error) {
	// Set up the dialer for sending emails
	dialer := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
	)

	// Determine templates directory - try different paths
	templatesDir := ""
	if templatesDir == "" {
		templatesDir = "./templates/emails"
	}

	// Load email templates - don't fail if templates directory doesn't exist
	templates, err := loadTemplates(templatesDir)
	if err != nil {
		logger.Warn("Failed to load email templates, will use text-only emails", map[string]any{
			"error": err.Error(),
			"path":  templatesDir,
		})
		templates = make(map[string]*template.Template)
	}

	worker := &EmailWorker{
		config:       config,
		logger:       logger,
		dialer:       dialer,
		templates:    templates,
		templatesDir: templatesDir,
		sender:       sender,
	}

	// Set up the processor for the async queue
	sender.SetProcessor(func(item any) error {
		emailMsg, ok := item.(EmailMessage)
		if !ok {
			return fmt.Errorf("invalid message type, expected EmailMessage")
		}
		return worker.processEmail(emailMsg)
	})

	return worker, nil
}

// Start starts the email worker
func (w *EmailWorker) Start() {
	w.sender.Start()
	w.logger.Info("Email worker started")
}

// Stop stops the email worker
func (w *EmailWorker) Stop() {
	w.sender.Stop()
	w.logger.Info("Email worker stopped")
}

// processEmail processes an email message from the queue
func (w *EmailWorker) processEmail(emailMsg EmailMessage) error {
	w.logger.Info("Processing email", map[string]any{
		"id":        emailMsg.ID,
		"recipient": emailMsg.Recipient,
		"template":  emailMsg.TemplateName,
	})

	// Send the email
	err := w.sendEmail(emailMsg)
	if err != nil {
		w.logger.Error("Failed to send email", err, map[string]any{
			"id":        emailMsg.ID,
			"recipient": emailMsg.Recipient,
		})
		return fmt.Errorf("failed to send email: %w", err)
	}

	w.logger.Info("Email sent successfully", map[string]any{
		"id":        emailMsg.ID,
		"recipient": emailMsg.Recipient,
	})

	return nil
}

// sendEmail sends an email using GoMail
func (w *EmailWorker) sendEmail(emailMsg EmailMessage) error {
	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", w.config.SenderEmail)
	m.SetHeader("To", emailMsg.Recipient)
	m.SetHeader("Subject", emailMsg.Subject)

	// Check if this is a direct text email
	if textBody, ok := emailMsg.Data["TextBody"].(string); ok {
		m.SetBody("text/plain", textBody)
	} else {
		// Try to use a template if available
		templateName := emailMsg.TemplateName
		tmpl, ok := w.templates[templateName]

		// If template found, render it
		if ok && tmpl != nil {
			var htmlBuffer strings.Builder
			err := tmpl.Execute(&htmlBuffer, emailMsg.Data)
			if err != nil {
				w.logger.Warn("Failed to render template, using plain text", map[string]any{
					"template": templateName,
					"error":    err.Error(),
				})
			} else {
				htmlBody := htmlBuffer.String()
				m.AddAlternative("text/html", htmlBody)
				
				// Generate plain text version from HTML
				plainText := generatePlainTextFromHTML(htmlBody)
				m.SetBody("text/plain", plainText)
			}
		} else {
			// No template found, create a simple text version
			content := fmt.Sprintf("Hello,\n\nThis is a message from %s.\n\n", w.config.SenderName)
			
			// Add any data values as simple text
			for k, v := range emailMsg.Data {
				content += fmt.Sprintf("%s: %v\n", k, v)
			}
			
			content += fmt.Sprintf("\n\nBest regards,\nThe %s Team", w.config.SenderName)
			m.SetBody("text/plain", content)
		}
	}

	// Add attachments if any
	for _, attachment := range emailMsg.Attachments {
		m.Attach(attachment.Filename,
			gomail.SetCopyFunc(func(w io.Writer) error {
				_, err := w.Write(attachment.Content)
				return err
			}),
			gomail.SetHeader(map[string][]string{
				"Content-Type": {attachment.MimeType},
			}),
		)
	}

	// Send the email
	return w.dialer.DialAndSend(m)
}

// loadTemplates loads all email templates from the templates directory
func loadTemplates(templatesDir string) (map[string]*template.Template, error) {
	templates := make(map[string]*template.Template)

	// Check if the directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("templates directory does not exist: %s", templatesDir)
	}

	// Walk through the templates directory
	err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-HTML files
		if filepath.Ext(path) != ".html" && filepath.Ext(path) != ".tmpl" {
			return nil
		}

		// Read the template file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		// Parse the template
		tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template %s: %w", path, err)
		}

		// Get the template name (filename without extension)
		templateName := filepath.Base(path)
		templateName = templateName[:len(templateName)-len(filepath.Ext(templateName))]

		// Store the template
		templates[templateName] = tmpl

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return templates, nil
}

// generatePlainTextFromHTML generates a plain text version from HTML
func generatePlainTextFromHTML(html string) string {
	// Remove HTML tags
	plainText := html
	plainText = tagRegex.ReplaceAllString(plainText, "")
	
	// Replace common HTML entities
	plainText = strings.ReplaceAll(plainText, "&nbsp;", " ")
	plainText = strings.ReplaceAll(plainText, "&amp;", "&")
	plainText = strings.ReplaceAll(plainText, "&lt;", "<")
	plainText = strings.ReplaceAll(plainText, "&gt;", ">")
	
	// Normalize whitespace
	plainText = strings.TrimSpace(plainText)
	plainText = whitespaceRegex.ReplaceAllString(plainText, " ")
	
	return plainText
}

// Regular expressions for HTML to text conversion
var (
	tagRegex       = regexp.MustCompile(`<[^>]*>`)
	whitespaceRegex = regexp.MustCompile(`\s+`)
)