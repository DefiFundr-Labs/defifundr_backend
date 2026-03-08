package waitlistport

import (
	"context"

	waitlistdomain "github.com/demola234/defifundr/internal/features/waitlist/domain"
	"github.com/google/uuid"
)

// Repository defines the data access operations for the waitlist feature.
type Repository interface {
	CreateWaitlistEntry(ctx context.Context, entry waitlistdomain.WaitlistEntry) (*waitlistdomain.WaitlistEntry, error)
	GetWaitlistEntryByEmail(ctx context.Context, email string) (*waitlistdomain.WaitlistEntry, error)
	GetWaitlistEntryByID(ctx context.Context, id uuid.UUID) (*waitlistdomain.WaitlistEntry, error)
	GetWaitlistEntryByReferralCode(ctx context.Context, code string) (*waitlistdomain.WaitlistEntry, error)
	ListWaitlistEntries(ctx context.Context, limit, offset int, filters map[string]string) ([]waitlistdomain.WaitlistEntry, int64, error)
	ExportWaitlistToCsv(ctx context.Context) ([]byte, error)
}
