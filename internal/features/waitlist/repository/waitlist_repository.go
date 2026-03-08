package waitlistrepo

import (
	"bytes"
	"context"
	"encoding/csv"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	waitlistdomain "github.com/demola234/defifundr/internal/features/waitlist/domain"
	waitlistport "github.com/demola234/defifundr/internal/features/waitlist/port"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type WaitlistRepository struct {
	store db.Queries
}

// New creates a new WaitlistRepository.
func New(store db.Queries) waitlistport.Repository {
	return &WaitlistRepository{store: store}
}

func (r *WaitlistRepository) CreateWaitlistEntry(ctx context.Context, entry waitlistdomain.WaitlistEntry) (*waitlistdomain.WaitlistEntry, error) {
	ctx, span := tracing.Tracer("waitlist-repository").Start(ctx, "CreateWaitlistEntry")
	defer span.End()

	var invitedDate pgtype.Timestamptz
	if entry.InvitedDate != nil {
		invitedDate = pgtype.Timestamptz{Time: *entry.InvitedDate, Valid: true}
	}

	var registeredDate pgtype.Timestamptz
	if entry.RegisteredDate != nil {
		registeredDate = pgtype.Timestamptz{Time: *entry.RegisteredDate, Valid: true}
	}

	params := db.CreateWaitlistEntryParams{
		ID:             entry.ID,
		Email:          entry.Email,
		FullName:       pgtype.Text{String: entry.FullName, Valid: entry.FullName != ""},
		ReferralCode:   entry.ReferralCode,
		ReferralSource: pgtype.Text{String: entry.ReferralSource, Valid: entry.ReferralSource != ""},
		Status:         entry.Status,
		SignupDate:     entry.SignupDate,
		InvitedDate:    invitedDate,
		RegisteredDate: registeredDate,
		Metadata:       nil,
	}

	dbEntry, err := r.store.CreateWaitlistEntry(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapToWaitlistEntry(dbEntry), nil
}

func (r *WaitlistRepository) GetWaitlistEntryByEmail(ctx context.Context, email string) (*waitlistdomain.WaitlistEntry, error) {
	ctx, span := tracing.Tracer("waitlist-repository").Start(ctx, "GetWaitlistEntryByEmail")
	defer span.End()

	dbEntry, err := r.store.GetWaitlistEntryByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return mapToWaitlistEntry(dbEntry), nil
}

func (r *WaitlistRepository) GetWaitlistEntryByID(ctx context.Context, id uuid.UUID) (*waitlistdomain.WaitlistEntry, error) {
	ctx, span := tracing.Tracer("waitlist-repository").Start(ctx, "GetWaitlistEntryByID")
	defer span.End()

	dbEntry, err := r.store.GetWaitlistEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapToWaitlistEntry(dbEntry), nil
}

func (r *WaitlistRepository) GetWaitlistEntryByReferralCode(ctx context.Context, code string) (*waitlistdomain.WaitlistEntry, error) {
	ctx, span := tracing.Tracer("waitlist-repository").Start(ctx, "GetWaitlistEntryByReferralCode")
	defer span.End()

	dbEntry, err := r.store.GetWaitlistEntryByReferralCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return mapToWaitlistEntry(dbEntry), nil
}

func (r *WaitlistRepository) ListWaitlistEntries(ctx context.Context, limit, offset int, filters map[string]string) ([]waitlistdomain.WaitlistEntry, int64, error) {
	ctx, span := tracing.Tracer("waitlist-repository").Start(ctx, "ListWaitlistEntries")
	defer span.End()

	params := db.ListWaitlistEntriesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	dbEntries, err := r.store.ListWaitlistEntries(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	entries := make([]waitlistdomain.WaitlistEntry, len(dbEntries))
	for i, dbEntry := range dbEntries {
		entries[i] = *mapToWaitlistEntry(dbEntry)
	}

	total, err := r.store.CountWaitlistEntries(ctx, db.CountWaitlistEntriesParams{})
	if err != nil {
		return nil, 0, err
	}

	return entries, total, nil
}

func (r *WaitlistRepository) ExportWaitlistToCsv(ctx context.Context) ([]byte, error) {
	ctx, span := tracing.Tracer("waitlist-repository").Start(ctx, "ExportWaitlistToCsv")
	defer span.End()

	dbEntries, err := r.store.ExportWaitlistEntries(ctx)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	headers := []string{"ID", "Email", "Full Name", "Referral Code", "Referral Source", "Status", "Signup Date", "Invited Date", "Registered Date"}
	if err := writer.Write(headers); err != nil {
		return nil, err
	}

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

// mapToWaitlistEntry maps a SQLC Waitlist row to a domain WaitlistEntry.
func mapToWaitlistEntry(dbEntry db.Waitlist) *waitlistdomain.WaitlistEntry {
	result := &waitlistdomain.WaitlistEntry{
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

	return result
}
