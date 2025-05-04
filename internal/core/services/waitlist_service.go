// internal/core/services/waitlist_service.go
package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/demola234/defifundr/pkg/app_errors"
	"github.com/demola234/defifundr/pkg/random"
	"github.com/google/uuid"
)

type waitlistService struct {
	waitlistRepo ports.WaitlistRepository
	emailService ports.EmailService
}

// NewWaitlistService creates a new waitlist service
func NewWaitlistService(waitlistRepo ports.WaitlistRepository, emailService ports.EmailService) ports.WaitlistService {
	return &waitlistService{
		waitlistRepo: waitlistRepo,
		emailService: emailService,
	}
}

// JoinWaitlist implements ports.WaitlistService
func (s *waitlistService) JoinWaitlist(ctx context.Context, email, fullName, referralSource string) (*domain.WaitlistEntry, error) {
	// Check if email is already on waitlist
	existingEntry, err := s.waitlistRepo.GetWaitlistEntryByEmail(ctx, email)
	if err == nil && existingEntry != nil {
		return nil, appErrors.NewConflictError("Email already on waitlist")
	}
	// Generate a unique referral code
	referralCode := generateReferralCode(fullName)

	// Create waitlist entry
	entry := domain.WaitlistEntry{
		ID:             uuid.New(),
		Email:          email,
		FullName:       fullName,
		ReferralCode:   referralCode,
		ReferralSource: referralSource,
		Status:         "waiting",
		SignupDate:     time.Now(),
		Metadata:       make(map[string]interface{}),
	}

	// Store in database
	savedEntry, err := s.waitlistRepo.CreateWaitlistEntry(ctx, entry)
	if err != nil {
		return nil, fmt.Errorf("failed to create waitlist entry: %w", err)
	}

	// Get waitlist position
	position, err := s.GetWaitlistPosition(ctx, savedEntry.ID)
	if err != nil {
		// Non-fatal error, continue with default position value
		position = 0
	}

	// Send confirmation email
	err = s.emailService.SendWaitlistConfirmation(ctx, email, fullName, referralCode, position)
	if err != nil {
		// Log error but continue - email sending is non-critical
		fmt.Printf("Failed to send waitlist confirmation email: %v\n", err)
	}

	return savedEntry, nil
}

// GetWaitlistPosition implements ports.WaitlistService
func (s *waitlistService) GetWaitlistPosition(ctx context.Context, id uuid.UUID) (int, error) {
	// Get the entry to check its signup date
	_, err := s.waitlistRepo.GetWaitlistEntryByID(ctx, id)
	if err != nil {
		return 0, err
	}

	// Get all waiting entries sorted by signup date
	entries, _, err := s.waitlistRepo.ListWaitlistEntries(ctx, 1000000, 0, map[string]string{"status": "waiting"})
	if err != nil {
		return 0, err
	}

	// Find position in the list
	for i, e := range entries {
		if e.ID == id {
			return i + 1, nil
		}
	}

	return 0, appErrors.NewNotFoundError("Entry not found in waitlist")
}


// GetWaitlistStats implements ports.WaitlistService
func (s *waitlistService) GetWaitlistStats(ctx context.Context) (map[string]interface{}, error) {
	// Get waitlist entries for different statuses
	waiting, _, err := s.waitlistRepo.ListWaitlistEntries(ctx, 1000000, 0, map[string]string{"status": "waiting"})
	if err != nil {
		return nil, err
	}

	invited, _, err := s.waitlistRepo.ListWaitlistEntries(ctx, 1000000, 0, map[string]string{"status": "invited"})
	if err != nil {
		return nil, err
	}

	registered, _, err := s.waitlistRepo.ListWaitlistEntries(ctx, 1000000, 0, map[string]string{"status": "registered"})
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	sources := make(map[string]int)
	for _, entry := range waiting {
		if entry.ReferralSource != "" {
			sources[entry.ReferralSource]++
		}
	}

	// Compile stats
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

// ListWaitlist implements ports.WaitlistService
func (s *waitlistService) ListWaitlist(ctx context.Context, page, pageSize int, filters map[string]string) ([]domain.WaitlistEntry, int64, error) {
	offset := (page - 1) * pageSize
	return s.waitlistRepo.ListWaitlistEntries(ctx, pageSize, offset, filters)
}

// ExportWaitlist implements ports.WaitlistService
func (s *waitlistService) ExportWaitlist(ctx context.Context) ([]byte, error) {
	return s.waitlistRepo.ExportWaitlistToCsv(ctx)
}

// GetWaitlistEntryByID gets a waitlist entry by ID
func (s *waitlistService) GetWaitlistEntryByID(ctx context.Context, id uuid.UUID) (*domain.WaitlistEntry, error) {
	entry, err := s.waitlistRepo.GetWaitlistEntryByID(ctx, id)
	if err != nil {
		return nil, appErrors.NewNotFoundError("Waitlist entry not found")
	}
	return entry, nil
}

// Helper functions

// generateReferralCode creates a referral code based on name and a random string
func generateReferralCode(name string) string {
	if name == "" {
		return random.RandomString(8)
	}

	// Take first 3 characters of name (if available)
	prefix := ""
	nameParts := strings.Fields(name)
	if len(nameParts) > 0 {
		firstPart := nameParts[0]
		if len(firstPart) >= 3 {
			prefix = strings.ToUpper(firstPart[0:3])
		} else if len(firstPart) > 0 {
			prefix = strings.ToUpper(firstPart)
		}
	}
	
	// Add random characters to make it unique
	suffix := random.RandomString(5)
	return fmt.Sprintf("%s%s", prefix, suffix)
}

// calculateConversionRate calculates the conversion rate as a percentage
func calculateConversionRate(invited, registered int) float64 {
	if invited == 0 {
		return 0.0
	}
	return float64(registered) / float64(invited) * 100.0
}