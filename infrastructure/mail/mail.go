package mail

import (
	"time"

	emailEnums "github.com/demola234/defifundr/pkg/utils"
)

// EmailMessage represents an email message in the queue
type EmailMessage struct {
	ID           string                       `json:"id"`
	Recipient    string                       `json:"recipient"`
	Subject      string                       `json:"subject"`
	TemplateName string                       `json:"template_name"`
	Data         map[string]any       `json:"data"`
	Attachments  []emailEnums.EmailAttachment `json:"attachments,omitempty"`
	Priority     emailEnums.EmailPriority     `json:"priority"`
	CreatedAt    time.Time                    `json:"created_at"`
}
