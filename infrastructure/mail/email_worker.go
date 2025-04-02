package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	amqp "github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
)

const (
	maxRetries    = 3
	retryInterval = 5 * time.Second
)

// EmailWorker processes emails from the queue and sends them
type EmailWorker struct {
	config         config.Config
	logger         logging.Logger
	conn           *amqp.Connection
	channel        *amqp.Channel
	dialer         *gomail.Dialer
	templates      map[string]*template.Template
	templatesDir   string
	workerCount    int
	shutdownCh     chan struct{}
	workerWg       sync.WaitGroup
}

// NewEmailWorker creates a new email worker
func NewEmailWorker(config config.Config, logger logging.Logger, workerCount int) (*EmailWorker, error) {
	// Connect to RabbitMQ
	conn, err := amqp.Dial(config.RabbitMqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Set up the dialer for sending emails
	dialer := gomail.NewDialer(
		config.SMTPHost,
		config.SMTPPort,
		config.SMTPUsername,
		config.SMTPPassword,
	)

	// Load email templates
	templatesDir := "./templates"
	templates, err := loadTemplates(templatesDir)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to load email templates: %w", err)
	}

	return &EmailWorker{
		config:        config,
		logger:        logger,
		conn:          conn,
		channel:       channel,
		dialer:        dialer,
		templates:     templates,
		templatesDir:  templatesDir,
		workerCount:   workerCount,
		shutdownCh:    make(chan struct{}),
	}, nil
}

// Start starts the email worker
func (w *EmailWorker) Start() error {
	// Ensure the queue exists
	_, err := w.channel.QueueDeclare(
		emailQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Prefetch count
	err = w.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start worker goroutines
	for i := 0; i < w.workerCount; i++ {
		workerID := i
		w.workerWg.Add(1)

		go func() {
			defer w.workerWg.Done()
			w.runWorker(workerID)
		}()
	}

	w.logger.Info("Email worker started", map[string]interface{}{
		"worker_count": w.workerCount,
	})

	return nil
}

// Stop stops the email worker
func (w *EmailWorker) Stop() error {
	// Signal all workers to stop
	close(w.shutdownCh)

	// Wait for all workers to finish
	w.workerWg.Wait()

	// Close the channel and connection
	if w.channel != nil {
		w.channel.Close()
	}
	if w.conn != nil {
		w.conn.Close()
	}

	w.logger.Info("Email worker stopped")
	return nil
}

// runWorker runs a worker to process emails from the queue
func (w *EmailWorker) runWorker(workerID int) {
	w.logger.Info("Starting email worker", map[string]interface{}{
		"worker_id": workerID,
	})

	// Create a consumer
	msgs, err := w.channel.Consume(
		emailQueue,        // queue
		fmt.Sprintf("email-worker-%d", workerID), // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	if err != nil {
		w.logger.Error("Failed to register a consumer", err, map[string]interface{}{
			"worker_id": workerID,
		})
		return
	}

	for {
		select {
		case <-w.shutdownCh:
			// Worker is being shut down
			w.logger.Info("Worker shutting down", map[string]interface{}{
				"worker_id": workerID,
			})
			return

		case msg, ok := <-msgs:
			if !ok {
				// Channel was closed
				w.logger.Error("Channel closed", nil, map[string]interface{}{
					"worker_id": workerID,
				})
				return
			}

			// Process the message
			err := w.processMessage(msg)

			if err != nil {
				// Log the error
				w.logger.Error("Failed to process message", err, map[string]interface{}{
					"worker_id": workerID,
				})

				// Get retry count from headers
				retryCount := 0
				if headers, ok := msg.Headers["x-retry-count"].(int); ok {
					retryCount = headers
				}

				// If we haven't exceeded max retries, requeue the message
				if retryCount < maxRetries {
					// Increment retry count
					if msg.Headers == nil {
						msg.Headers = make(amqp.Table)
					}
					msg.Headers["x-retry-count"] = retryCount + 1

					// Requeue the message after a delay
					time.Sleep(retryInterval)
					err = w.requeueMessage(msg)
					if err != nil {
						w.logger.Error("Failed to requeue message", err, map[string]interface{}{
							"worker_id": workerID,
						})
						// Nack the message without requeuing
						msg.Nack(false, false)
					}
				} else {
					// Max retries exceeded, move to dead letter queue
					err = w.moveToDeadLetterQueue(msg)
					if err != nil {
						w.logger.Error("Failed to move message to dead letter queue", err, map[string]interface{}{
							"worker_id": workerID,
						})
					}
					// Nack the message without requeuing
					msg.Nack(false, false)
				}
			} else {
				// Acknowledge the message
				msg.Ack(false)
			}
		}
	}
}

// processMessage processes an email message from the queue
func (w *EmailWorker) processMessage(msg amqp.Delivery) error {
	// Parse the message
	var emailMsg EmailMessage
	err := json.Unmarshal(msg.Body, &emailMsg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal email message: %w", err)
	}

	w.logger.Info("Processing email", map[string]interface{}{
		"id":        emailMsg.ID,
		"recipient": emailMsg.Recipient,
		"template":  emailMsg.TemplateName,
	})

	// Send the email
	err = w.sendEmail(emailMsg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	w.logger.Info("Email sent successfully", map[string]interface{}{
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

	// Render the email template
	templateName := emailMsg.TemplateName
	tmpl, ok := w.templates[templateName]
	if !ok {
		return fmt.Errorf("template not found: %s", templateName)
	}

	// Render the HTML template
	var htmlBody string
	if tmpl != nil {
		var htmlBuffer strings.Builder
		err := tmpl.Execute(&htmlBuffer, emailMsg.Data)
		if err != nil {
			return fmt.Errorf("failed to render email template: %w", err)
		}
		htmlBody = htmlBuffer.String()
		m.AddAlternative("text/html", htmlBody)
	}

	// Generate plain text version from HTML
	plainText := generatePlainTextFromHTML(htmlBody)
	m.SetBody("text/plain", plainText)

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

// requeueMessage requeues a message to the main queue
func (w *EmailWorker) requeueMessage(msg amqp.Delivery) error {
	return w.channel.PublishWithContext(
		context.Background(),
		"",        // exchange
		emailQueue, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			DeliveryMode: msg.DeliveryMode,
			Priority:     msg.Priority,
			Headers:      msg.Headers,
		},
	)
}

// moveToDeadLetterQueue moves a message to the dead letter queue
func (w *EmailWorker) moveToDeadLetterQueue(msg amqp.Delivery) error {
	// Ensure the dead letter queue exists
	_, err := w.channel.QueueDeclare(
		emailQueue+".dead", // name
		true,              // durable
		false,             // delete when unused
		false,             // exclusive
		false,             // no-wait
		nil,               // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare dead letter queue: %w", err)
	}

	// Publish to the dead letter queue
	return w.channel.PublishWithContext(
		context.Background(),
		"",                 // exchange
		emailQueue+".dead", // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			DeliveryMode: msg.DeliveryMode,
			Priority:     msg.Priority,
			Headers:      msg.Headers,
		},
	)
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
		content, err := ioutil.ReadFile(path)
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
	// This is a simple implementation that removes HTML tags
	// For a more robust solution, consider using a proper HTML to text converter
	
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