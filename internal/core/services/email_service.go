package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/internal/core/ports"
)

type EmailService struct {
	config    config.Config
	logger    logging.Logger
	emailSender ports.EmailSender
}

// NewEmailService creates a new email service
func NewEmailService(config config.Config, logger logging.Logger, emailSender ports.EmailSender) ports.EmailService {
	return &EmailService{
		config:      config,
		logger:      logger,
		emailSender: emailSender,
	}
}

// SendWaitlistConfirmation sends a waitlist confirmation email as plain text
func (s *EmailService) SendWaitlistConfirmation(ctx context.Context, email, name, referralCode string, position int) error {
	if s.isTestMode() {
		s.logger.Info("Test mode: Would send waitlist confirmation")
		return nil
	}

	subject := "Welcome to DefiFundr Waitlist"
	
	// Create plain text email body
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
	
	// Prepare data for the email sender
	templateData := map[string]interface{}{
		"Name":         name,
		"Position":     position,
		"ReferralCode": referralCode,
		"AppName":      "DefiFundr",
		"TextBody":     body,
	}

	// Queue email with normal priority
	_, err := s.emailSender.QueueEmail(ctx, email, subject, "text_email", templateData, ports.NormalPriority)
	if err != nil {
		s.logger.Error("Failed to queue waitlist confirmation email", err, map[string]interface{}{
			"email": email,
		})
		return fmt.Errorf("failed to queue waitlist confirmation email: %w", err)
	}

	s.logger.Info("Queued waitlist confirmation email", map[string]interface{}{
		"email": email,
	})
	return nil
}

// SendWaitlistInvitation sends an invitation email to users from the waitlist
func (s *EmailService) SendWaitlistInvitation(ctx context.Context, email, name string, inviteLink string) error {
	if s.isTestMode() {
		s.logger.Info("Test mode: Would send waitlist invitation")
		return nil
	}

	subject := "You're In! Access DefiFundr Now"
	
	// Prepare template data
	templateData := map[string]interface{}{
		"Name":       name,
		"InviteLink": inviteLink,
		"AppName":    "DefiFundr",
	}

	// Queue email with high priority
	_, err := s.emailSender.QueueEmail(ctx, email, subject, "waitlist_invitation", templateData, ports.HighPriority)
	if err != nil {
		s.logger.Error("Failed to queue waitlist invitation email", err, map[string]interface{}{
			"email": email,
		})
		return fmt.Errorf("failed to queue waitlist invitation email: %w", err)
	}

	s.logger.Info("Queued waitlist invitation email", map[string]interface{}{
		"email": email,
	})
	return nil
}

// SendBatchUpdate sends a batch update email to multiple waitlist members
func (s *EmailService) SendBatchUpdate(ctx context.Context, emails []string, subject, message string) error {
	if s.isTestMode() {
		s.logger.Info("Test mode: Would send batch update")
		return nil
	}

	// Prepare template data
	templateData := map[string]interface{}{
		"Message": message,
		"AppName": "DefiFundr",
	}

	for _, email := range emails {
		// Queue email with low priority (for batch operations)
		_, err := s.emailSender.QueueEmail(ctx, email, subject, "general_update", templateData, ports.LowPriority)
		if err != nil {
			s.logger.Error("Failed to queue batch update email", err, map[string]interface{}{
				"email": email,
			})
			// Continue with other emails instead of returning early
			continue
		}
		s.logger.Debug("Queued batch update email for", map[string]interface{}{
			"email": email,
		})
	}

	s.logger.Info("Completed queuing batch update emails", map[string]interface{}{
		"total": len(emails),
	})
	return nil
}

// isTestMode checks if the service is running in test mode
func (s *EmailService) isTestMode() bool {
	return strings.ToLower(s.config.Environment) == "test" || 
		   strings.ToLower(s.config.Environment) == "development"
}