package waitlistusecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	waitlistdomain "github.com/demola234/defifundr/internal/features/waitlist/domain"
	waitlistport "github.com/demola234/defifundr/internal/features/waitlist/port"
	appErrors "github.com/demola234/defifundr/pkg/apperrors"
	"github.com/demola234/defifundr/pkg/random"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/google/uuid"
)

type waitlistUseCase struct {
	repo         waitlistport.Repository
	emailSender  waitlistport.EmailSender
}

// New creates a new waitlist use case.
func New(repo waitlistport.Repository, emailSender waitlistport.EmailSender) waitlistport.Service {
	return &waitlistUseCase{
		repo:        repo,
		emailSender: emailSender,
	}
}

// JoinWaitlist implements waitlistport.Service
func (uc *waitlistUseCase) JoinWaitlist(ctx context.Context, email, fullName, referralSource string) (*waitlistdomain.WaitlistEntry, error) {
	ctx, span := tracing.Tracer("waitlist-usecase").Start(ctx, "JoinWaitlist")
	defer span.End()

	existing, err := uc.repo.GetWaitlistEntryByEmail(ctx, email)
	if err == nil && existing != nil {
		conflictErr := appErrors.NewConflictError("Email already on waitlist")
		span.RecordError(conflictErr)
		return nil, conflictErr
	}

	referralCode := generateReferralCode(fullName)

	entry := waitlistdomain.WaitlistEntry{
		ID:             uuid.New(),
		Email:          email,
		FullName:       fullName,
		ReferralCode:   referralCode,
		ReferralSource: referralSource,
		Status:         "waiting",
		SignupDate:     time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	saved, err := uc.repo.CreateWaitlistEntry(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("failed to create waitlist entry: %w", err)
	}

	position, err := uc.GetWaitlistPosition(ctx, saved.ID)
	if err != nil {
		position = 0
	}

	if sendErr := uc.emailSender.SendWaitlistConfirmation(ctx, email, fullName, referralCode, position); sendErr != nil {
		fmt.Printf("Failed to send waitlist confirmation email: %v\n", sendErr)
	}

	return saved, nil
}

// GetWaitlistPosition implements waitlistport.Service
func (uc *waitlistUseCase) GetWaitlistPosition(ctx context.Context, id uuid.UUID) (int, error) {
	ctx, span := tracing.Tracer("waitlist-usecase").Start(ctx, "GetWaitlistPosition")
	defer span.End()

	_, err := uc.repo.GetWaitlistEntryByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}

	entries, _, err := uc.repo.ListWaitlistEntries(ctx, 1000000, 0, map[string]string{"status": "waiting"})
	if err != nil {
		span.RecordError(err)
		return 0, err
	}

	for i, e := range entries {
		if e.ID == id {
			return i + 1, nil
		}
	}

	notFoundErr := appErrors.NewNotFoundError("Entry not found in waitlist")
	span.RecordError(notFoundErr)
	return 0, notFoundErr
}

// GetWaitlistStats implements waitlistport.Service
func (uc *waitlistUseCase) GetWaitlistStats(ctx context.Context) (map[string]interface{}, error) {
	ctx, span := tracing.Tracer("waitlist-usecase").Start(ctx, "GetWaitlistStats")
	defer span.End()

	waiting, _, err := uc.repo.ListWaitlistEntries(ctx, 1000000, 0, map[string]string{"status": "waiting"})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	invited, _, err := uc.repo.ListWaitlistEntries(ctx, 1000000, 0, map[string]string{"status": "invited"})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	registered, _, err := uc.repo.ListWaitlistEntries(ctx, 1000000, 0, map[string]string{"status": "registered"})
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	sources := make(map[string]int)
	for _, e := range waiting {
		if e.ReferralSource != "" {
			sources[e.ReferralSource]++
		}
	}

	stats := map[string]interface{}{
		"total_signups":    len(waiting) + len(invited) + len(registered),
		"waiting_count":    len(waiting),
		"invited_count":    len(invited),
		"registered_count": len(registered),
		"conversion_rate":  calculateConversionRate(len(invited), len(registered)),
		"sources":          sources,
	}

	return stats, nil
}

// ListWaitlist implements waitlistport.Service
func (uc *waitlistUseCase) ListWaitlist(ctx context.Context, page, pageSize int, filters map[string]string) ([]waitlistdomain.WaitlistEntry, int64, error) {
	ctx, span := tracing.Tracer("waitlist-usecase").Start(ctx, "ListWaitlist")
	defer span.End()

	offset := (page - 1) * pageSize
	entries, total, err := uc.repo.ListWaitlistEntries(ctx, pageSize, offset, filters)
	if err != nil {
		span.RecordError(err)
		return nil, 0, err
	}
	return entries, total, nil
}

// ExportWaitlist implements waitlistport.Service
func (uc *waitlistUseCase) ExportWaitlist(ctx context.Context) ([]byte, error) {
	ctx, span := tracing.Tracer("waitlist-usecase").Start(ctx, "ExportWaitlist")
	defer span.End()

	data, err := uc.repo.ExportWaitlistToCsv(ctx)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return data, nil
}

// generateReferralCode creates a referral code based on name and a random string.
func generateReferralCode(name string) string {
	if name == "" {
		return random.RandomString(8)
	}

	prefix := ""
	parts := strings.Fields(name)
	if len(parts) > 0 {
		first := parts[0]
		if len(first) >= 3 {
			prefix = strings.ToUpper(first[0:3])
		} else if len(first) > 0 {
			prefix = strings.ToUpper(first)
		}
	}

	return fmt.Sprintf("%s%s", prefix, random.RandomString(5))
}

// calculateConversionRate returns invited→registered conversion as a percentage.
func calculateConversionRate(invited, registered int) float64 {
	if invited == 0 {
		return 0.0
	}
	return float64(registered) / float64(invited) * 100.0
}
