package repositories

import (
	"bytes"
	"context"
	"encoding/csv"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type WaitlistRepository struct {
	store db.Queries
}

func NewWaitlistRepository(store db.Queries) *WaitlistRepository {
	return &WaitlistRepository{store: store}
}

func (r *WaitlistRepository) CreateWaitlistEntry(ctx context.Context, entry domain.WaitlistEntry) (*domain.WaitlistEntry, error) {


	// Prepare invited and registered dates
	var invitedDate pgtype.Timestamptz
	if entry.InvitedDate != nil {
		invitedDate = pgtype.Timestamptz{
			Time:   *entry.InvitedDate,
			Valid:  true,
		}
	}

	var registeredDate pgtype.Timestamptz
	if entry.RegisteredDate != nil {
		registeredDate = pgtype.Timestamptz{
			Time:   *entry.RegisteredDate,
			Valid:  true,
		}
	}

	// Create params for query
	params := db.CreateWaitlistEntryParams{
		ID:              entry.ID,
		Email:           entry.Email,
		FullName:        pgtype.Text{String: entry.FullName, Valid: entry.FullName != ""},
		ReferralCode:    entry.ReferralCode,
		ReferralSource:  pgtype.Text{String: entry.ReferralSource, Valid: entry.ReferralSource != ""},
		Status:          entry.Status,
		SignupDate:      entry.SignupDate,
		InvitedDate:     invitedDate,
		RegisteredDate:  registeredDate,
		Metadata:        nil,
	}

	dbEntry, err := r.store.CreateWaitlistEntry(ctx, params)
	if err != nil {
		return nil, err
	}

	// Map DB result back to domain model
	result := &domain.WaitlistEntry{
		ID:             dbEntry.ID,
		Email:          dbEntry.Email,
		FullName:       dbEntry.FullName.String,
		ReferralCode:   dbEntry.ReferralCode,
		ReferralSource: dbEntry.ReferralSource.String,
		Status:         dbEntry.Status,
		SignupDate:     dbEntry.SignupDate,
		Metadata:       make(map[string]interface{}),
	}

	if dbEntry.InvitedDate.Valid {
		result.InvitedDate = &dbEntry.InvitedDate.Time
	}

	if dbEntry.RegisteredDate.Valid {
		result.RegisteredDate = &dbEntry.RegisteredDate.Time
	}


	return result, nil
}

func (r *WaitlistRepository) GetWaitlistEntryByEmail(ctx context.Context, email string) (*domain.WaitlistEntry, error) {
	dbEntry, err := r.store.GetWaitlistEntryByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	// Map DB result to domain model
	result := &domain.WaitlistEntry{
		ID:             dbEntry.ID,
		Email:          dbEntry.Email,
		FullName:       dbEntry.FullName.String,
		ReferralCode:   dbEntry.ReferralCode,
		ReferralSource: dbEntry.ReferralSource.String,
		Status:         dbEntry.Status,
		SignupDate:     dbEntry.SignupDate,
		Metadata:       make(map[string]interface{}),
	}

	if dbEntry.InvitedDate.Valid {
		result.InvitedDate = &dbEntry.InvitedDate.Time
	}

	if dbEntry.RegisteredDate.Valid {
		result.RegisteredDate = &dbEntry.RegisteredDate.Time
	}

	return result, nil
}

func (r *WaitlistRepository) GetWaitlistEntryByID(ctx context.Context, id uuid.UUID) (*domain.WaitlistEntry, error) {
	dbEntry, err := r.store.GetWaitlistEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Map DB result to domain model
	result := &domain.WaitlistEntry{
		ID:             dbEntry.ID,
		Email:          dbEntry.Email,
		FullName:       dbEntry.FullName.String,
		ReferralCode:   dbEntry.ReferralCode,
		ReferralSource: dbEntry.ReferralSource.String,
		Status:         dbEntry.Status,
		SignupDate:     dbEntry.SignupDate,
		Metadata:       make(map[string]interface{}),
	}

	if dbEntry.InvitedDate.Valid {
		result.InvitedDate = &dbEntry.InvitedDate.Time
	}

	if dbEntry.RegisteredDate.Valid {
		result.RegisteredDate = &dbEntry.RegisteredDate.Time
	}

	return result, nil
}

func (r *WaitlistRepository) GetWaitlistEntryByReferralCode(ctx context.Context, code string) (*domain.WaitlistEntry, error) {
	dbEntry, err := r.store.GetWaitlistEntryByReferralCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Map DB result to domain model
	result := &domain.WaitlistEntry{
		ID:             dbEntry.ID,
		Email:          dbEntry.Email,
		FullName:       dbEntry.FullName.String,
		ReferralCode:   dbEntry.ReferralCode,
		ReferralSource: dbEntry.ReferralSource.String,
		Status:         dbEntry.Status,
		SignupDate:     dbEntry.SignupDate,
		Metadata:       make(map[string]interface{}),
	}

	if dbEntry.InvitedDate.Valid {
		result.InvitedDate = &dbEntry.InvitedDate.Time
	}

	if dbEntry.RegisteredDate.Valid {
		result.RegisteredDate = &dbEntry.RegisteredDate.Time
	}

	return result, nil
}


func (r *WaitlistRepository) ListWaitlistEntries(ctx context.Context, limit, offset int, filters map[string]string) ([]domain.WaitlistEntry, int64, error) {
	// Build filter conditions based on the provided filters
	params := db.ListWaitlistEntriesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}


	// Get entries from DB
	dbEntries, err := r.store.ListWaitlistEntries(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	// Map DB results to domain models
	entries := make([]domain.WaitlistEntry, len(dbEntries))
	for i, dbEntry := range dbEntries {
		entries[i] = domain.WaitlistEntry{
			ID:             dbEntry.ID,
			Email:          dbEntry.Email,
			FullName:       dbEntry.FullName.String,
			ReferralCode:   dbEntry.ReferralCode,
			ReferralSource: dbEntry.ReferralSource.String,
			Status:         dbEntry.Status,
			SignupDate:     dbEntry.SignupDate,
			Metadata:       make(map[string]interface{}),
		}

		if dbEntry.InvitedDate.Valid {
			entries[i].InvitedDate = &dbEntry.InvitedDate.Time
		}

		if dbEntry.RegisteredDate.Valid {
			entries[i].RegisteredDate = &dbEntry.RegisteredDate.Time
		}

	}

	// Get total count for pagination
	countParams := db.CountWaitlistEntriesParams{}

	total, err := r.store.CountWaitlistEntries(ctx, countParams)
	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (r *WaitlistRepository) ExportWaitlistToCsv(ctx context.Context) ([]byte, error) {
	// Get all waitlist entries
	dbEntries, err := r.store.ExportWaitlistEntries(ctx)
	if err != nil {
		return nil, err
	}

	// Create a buffer to write the CSV data to
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write CSV header
	headers := []string{"ID", "Email", "Full Name", "Referral Code", "Referral Source", "Status", "Signup Date", "Invited Date", "Registered Date"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

	// Write each waitlist entry as a CSV row
	for _, entry := range dbEntries {
		invitedDate := ""
		if entry.InvitedDate.Valid {
			invitedDate = entry.InvitedDate.Time.Format(time.RFC3339)
		}

		registeredDate := ""
		if entry.RegisteredDate.Valid {
			registeredDate = entry.RegisteredDate.Time.Format(time.RFC3339)
		}

		row := []string{
			entry.ID.String(),
			entry.Email,
			entry.FullName.String,
			entry.ReferralCode,
			entry.ReferralSource.String,
			entry.Status,
			entry.SignupDate.Format(time.RFC3339),
			invitedDate,
			registeredDate,
		}

		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}