package waitlistport

import (
	"context"

	waitlistdomain "github.com/demola234/defifundr/internal/features/waitlist/domain"
	"github.com/google/uuid"
)

// Service defines the business operations for the waitlist feature.
type Service interface {
	JoinWaitlist(ctx context.Context, email, fullName, referralSource string) (*waitlistdomain.WaitlistEntry, error)
	GetWaitlistPosition(ctx context.Context, id uuid.UUID) (int, error)
	GetWaitlistStats(ctx context.Context) (map[string]any, error)
	ListWaitlist(ctx context.Context, page, pageSize int, filters map[string]string) ([]waitlistdomain.WaitlistEntry, int64, error)
	ExportWaitlist(ctx context.Context) ([]byte, error)
}

// EmailSender is the waitlist-scoped subset of the email service.
type EmailSender interface {
	SendWaitlistConfirmation(ctx context.Context, email, name, referralCode string, position int) error
}
