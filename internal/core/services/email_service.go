package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/internal/core/ports"
	"gopkg.in/gomail.v2"
)

type EmailService struct {
	config    config.Config
	logger    logging.Logger
	fromEmail string
	dialer    *gomail.Dialer
}

// NewEmailService creates a new email service using gomail
func NewEmailService(config config.Config, logger logging.Logger) ports.EmailService {
	dialer := gomail.NewDialer(
		config.SMTPHost, 
		config.SMTPPort, 
		config.SMTPUsername, 
		config.SMTPPassword,
	)
	
	return &EmailService{
		config:    config,
		logger:    logger,
		fromEmail: config.SMTPUsername,
		dialer:    dialer,
	}
}

// SendWaitlistConfirmation sends a waitlist confirmation email
func (s *EmailService) SendWaitlistConfirmation(ctx context.Context, email, name, referralCode string, position int) error {
	if s.isTestMode() {
		s.logger.Info("Test mode: Would send waitlist confirmation")
		return nil
	}

	subject := "Welcome to DefiFundr Waitlist"
	body := fmt.Sprintf(
		"Hello %s,\n\n"+
		"Thank you for joining the DefiFundr waitlist. We're excited to have you on board!\n\n"+
		"Your current position on the waitlist: %d\n\n"+
		"Want to move up? Share your unique referral code with friends: %s\n\n"+
		"We'll notify you as soon as we're ready to welcome you to the platform.\n\n"+
		"Best regards,\n"+
		"The DefiFundr Team",
		name, position, referralCode,
	)

	return s.sendEmail(email, subject, body)
}

// SendWaitlistInvitation sends an invitation email to users from the waitlist
func (s *EmailService) SendWaitlistInvitation(ctx context.Context, email, name string, inviteLink string) error {
	if s.isTestMode() {
		s.logger.Info("Test mode: Would send waitlist invitation")
		return nil
	}

	subject := "You're In! Access DefiFundr Now"
	body := fmt.Sprintf(
		"Hello %s,\n\n"+
		"Great news! You've been selected from our waitlist to access the DefiFundr platform.\n\n"+
		"Click this link to get started: %s\n\n"+
		"This invitation link is valid for the next 7 days.\n\n"+
		"Best regards,\n"+
		"The DefiFundr Team",
		name, inviteLink,
	)

	return s.sendEmail(email, subject, body)
}

// SendBatchUpdate sends a batch update email to multiple waitlist members
func (s *EmailService) SendBatchUpdate(ctx context.Context, emails []string, subject, message string) error {
	if s.isTestMode() {
		s.logger.Info("Test mode: Would send batch update")
		return nil
	}

	for _, email := range emails {
		err := s.sendEmail(email, subject, message)
		if err != nil {
			s.logger.Error("Failed to send batch update email", err, map[string]interface{}{
				"email": email,
			})
			// Continue with other emails instead of returning early
			continue
		}
		s.logger.Debug("Sent batch update email")
	}

	s.logger.Info("Completed sending batch update emails")
	return nil
}

// sendEmail is a helper function to send an email
func (s *EmailService) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.fromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if err := s.dialer.DialAndSend(m); err != nil {
		s.logger.Error("Failed to send email", err, map[string]interface{}{
			"to": to,
			"subject": subject,
		})
		return fmt.Errorf("failed to send email: %w", err)
	}

	s.logger.Info("Sent email successfully", map[string]interface{}{
		"to": to,
		"subject": subject,
	})
	return nil
}

// isTestMode checks if the service is running in test mode
func (s *EmailService) isTestMode() bool {
	return strings.ToLower(s.config.Environment) == "test" || 
		   strings.ToLower(s.config.Environment) == "development"
}