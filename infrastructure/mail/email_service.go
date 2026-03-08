package mail

import (
	"context"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
)

// EmailService implements email sending operations required by the auth feature.
type EmailService struct {
	config config.Config
	logger logging.Logger
	sender *AsyncQEmailSender
}

// NewEmailService creates a new EmailService.
func NewEmailService(cfg config.Config, logger logging.Logger, sender *AsyncQEmailSender) *EmailService {
	return &EmailService{config: cfg, logger: logger, sender: sender}
}

// SendPasswordResetEmail sends a password reset OTP to the user.
func (s *EmailService) SendPasswordResetEmail(ctx context.Context, email, name, otpCode string) error {
	return s.sender.SendEmail(ctx, email, "Password Reset", "password_reset", map[string]interface{}{
		"name":    name,
		"otp":     otpCode,
		"support": s.config.SenderEmail,
	})
}

// SendWaitlistConfirmation sends a waitlist confirmation email.
func (s *EmailService) SendWaitlistConfirmation(ctx context.Context, email, name, referralCode string, position int) error {
	return s.sender.SendEmail(ctx, email, "Welcome to the Waitlist", "waitlist_confirmation", map[string]interface{}{
		"name":         name,
		"referralCode": referralCode,
		"position":     position,
	})
}

// SendBatchUpdate sends a batch notification to multiple recipients.
func (s *EmailService) SendBatchUpdate(ctx context.Context, emails []string, subject, message string) error {
	for _, email := range emails {
		if err := s.sender.SendEmail(ctx, email, subject, "batch_update", map[string]interface{}{
			"message": message,
		}); err != nil {
			s.logger.Error("failed to send batch email", err, map[string]interface{}{"email": email})
		}
	}
	return nil
}
