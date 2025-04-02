// infrastructure/mail/rabbitmq_email_sender.go
package mail

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	emailExchange     = "emails"
	emailQueue        = "email_delivery"
	emailRoutingKey   = "email.send"
	connectionTimeout = 5 * time.Second
)

// EmailMessage represents an email message in the queue
type EmailMessage struct {
	ID          string                 `json:"id"`
	Recipient   string                 `json:"recipient"`
	Subject     string                 `json:"subject"`
	TemplateName string                `json:"template_name"`
	Data        map[string]interface{} `json:"data"`
	Attachments []ports.EmailAttachment `json:"attachments,omitempty"`
	Priority    ports.EmailPriority    `json:"priority"`
	CreatedAt   time.Time              `json:"created_at"`
}

// RabbitMQEmailSender is an implementation of EmailSender that uses RabbitMQ
type RabbitMQEmailSender struct {
	config  config.Config
	logger  logging.Logger
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQEmailSender creates a new RabbitMQ-based email sender
func NewRabbitMQEmailSender(config config.Config, logger logging.Logger) (ports.EmailSender, error) {
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

	// Declare the exchange
	err = channel.ExchangeDeclare(
		emailExchange, // name
		"direct",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare an exchange: %w", err)
	}

	// Declare the queue
	_, err = channel.QueueDeclare(
		emailQueue, // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	// Bind the queue to the exchange
	err = channel.QueueBind(
		emailQueue,      // queue name
		emailRoutingKey, // routing key
		emailExchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind a queue: %w", err)
	}

	return &RabbitMQEmailSender{
		config:  config,
		logger:  logger,
		conn:    conn,
		channel: channel,
	}, nil
}

// SendEmail sends an email directly
func (s *RabbitMQEmailSender) SendEmail(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}) error {
    _, err := s.QueueEmailWithAttachments(ctx, recipient, subject, templateName, data, nil, ports.NormalPriority)
    return err
}

// SendEmailWithAttachment sends an email with attachments directly
func (s *RabbitMQEmailSender) SendEmailWithAttachment(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}, attachments []ports.EmailAttachment) error {
    _, err := s.QueueEmailWithAttachments(ctx, recipient, subject, templateName, data, attachments, ports.NormalPriority)
    return err
}

// QueueEmail queues an email for asynchronous delivery
func (s *RabbitMQEmailSender) QueueEmail(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}, priority ports.EmailPriority) (string, error) {
	return s.QueueEmailWithAttachments(ctx, recipient, subject, templateName, data, nil, priority)
}

// QueueEmailWithAttachments queues an email with attachments for asynchronous delivery
func (s *RabbitMQEmailSender) QueueEmailWithAttachments(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}, attachments []ports.EmailAttachment, priority ports.EmailPriority) (string, error) {
	// Create a unique ID for the email
	id := uuid.New().String()

	// Create the email message
	message := EmailMessage{
		ID:           id,
		Recipient:    recipient,
		Subject:      subject,
		TemplateName: templateName,
		Data:         data,
		Attachments:  attachments,
		Priority:     priority,
		CreatedAt:    time.Now(),
	}

	// Convert the message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("failed to marshal email message: %w", err)
	}

	// Set message properties
	headers := make(amqp.Table)
	headers["priority"] = int(priority)

	// Publish the message to RabbitMQ
	ctxTimeout, cancel := context.WithTimeout(ctx, connectionTimeout)
	defer cancel()

	err = s.channel.PublishWithContext(
		ctxTimeout,
		emailExchange,   // exchange
		emailRoutingKey, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Priority:     uint8(priority),
			Headers:      headers,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to publish email message: %w", err)
	}

	s.logger.Info("Email queued", map[string]interface{}{
		"id":        id,
		"recipient": recipient,
		"template":  templateName,
		"priority":  priority,
	})

	return id, nil
}

// Close closes the RabbitMQ connection
func (s *RabbitMQEmailSender) Close() error {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}