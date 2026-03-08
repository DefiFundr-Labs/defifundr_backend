package mail

import (
	"context"
	"fmt"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	emailEnums "github.com/demola234/defifundr/pkg/utils"

	"github.com/google/uuid"
)

// AsyncQEmailSender is an implementation of EmailSender that uses AsyncQueue
type AsyncQEmailSender struct {
	config config.Config
	logger logging.Logger
	queue  *AsyncQueue
}

// NewAsyncQEmailSender creates a new AsyncQ-based email sender
func NewAsyncQEmailSender(config config.Config, logger logging.Logger) (*AsyncQEmailSender, error) {
	// Create a processor for the queue
	processor := func(item interface{}) error {
		//TODO
		return nil
	}

	// Create the async queue
	queue := NewAsyncQueue(1000, 5, logger, processor)

	sender := &AsyncQEmailSender{
		config: config,
		logger: logger,
		queue:  queue,
	}

	return sender, nil
}

// SendEmail sends an email directly
func (s *AsyncQEmailSender) SendEmail(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}) error {
	_, err := s.QueueEmail(ctx, recipient, subject, templateName, data, emailEnums.NormalPriority)
	return err
}

// SendEmailWithAttachment sends an email with attachments directly
func (s *AsyncQEmailSender) SendEmailWithAttachment(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}, attachments []emailEnums.EmailAttachment) error {
	_, err := s.QueueEmail(ctx, recipient, subject, templateName, data, emailEnums.NormalPriority)
	return err
}

// QueueEmail queues an email for asynchronous delivery
func (s *AsyncQEmailSender) QueueEmail(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}, priority emailEnums.EmailPriority) (string, error) {
	// Create a unique ID for the email
	id := uuid.New().String()

	// Create the email message
	message := EmailMessage{
		ID:           id,
		Recipient:    recipient,
		Subject:      subject,
		TemplateName: templateName,
		Data:         data,
		Attachments:  nil,
		Priority:     priority,
		CreatedAt:    time.Now(),
	}

	// Enqueue the message
	err := s.queue.EnqueueWithContext(ctx, message)
	if err != nil {
		return "", fmt.Errorf("failed to enqueue email message: %w", err)
	}

	s.logger.Info("Email queued", map[string]interface{}{
		"id":        id,
		"recipient": recipient,
		"template":  templateName,
		"priority":  priority,
	})

	return id, nil
}

// SetProcessor sets the processor function for the queue
func (s *AsyncQEmailSender) SetProcessor(processor func(interface{}) error) {
	s.queue.processor = processor
}

// Start starts the queue
func (s *AsyncQEmailSender) Start() {
	s.queue.Start()
}

// Stop stops the queue
func (s *AsyncQEmailSender) Stop() {
	s.queue.Stop()
}
