package repositories

import (
	"context"
	"fmt"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type SecurityRepository struct {
	store db.Queries
}

func NewSecurityRepository(store db.Queries) *SecurityRepository {
	return &SecurityRepository{
		store: store,
	}
}

// LogSecurityEvent logs a security event
func (r *SecurityRepository) LogSecurityEvent(ctx context.Context, event domain.SecurityEvent) error {
	params := db.CreateSecurityEventParams{
		ID:        event.ID,
		UserID:    event.UserID,
		EventType: event.EventType,
		IpAddress: event.IPAddress,
		UserAgent: toPgText(event.UserAgent),
		Timestamp: pgtype.Timestamp{Time: event.Timestamp, Valid: true},
	}

	_, err := r.store.CreateSecurityEvent(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to log security event: %w", err)
	}

	return nil
}

// GetRecentLoginsByUserID gets recent login events for a user
func (r *SecurityRepository) GetRecentLoginsByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.SecurityEvent, error) {
	params := db.GetRecentLoginEventsByUserIDParams{
		UserID: userID,
		Limit:  int32(limit),
	}

	events, err := r.store.GetRecentLoginEventsByUserID(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent logins: %w", err)
	}

	result := make([]domain.SecurityEvent, len(events))
	for i, event := range events {
		result[i] = domain.SecurityEvent{
			ID:        event.ID,
			UserID:    event.UserID,
			EventType: event.EventType,
			IPAddress: event.IpAddress,
			UserAgent: event.UserAgent.String,
			Timestamp: event.Timestamp.Time,
			// Parse metadata from JSON
			Metadata: make(map[string]interface{}),
		}
	}

	return result, nil
}

// GetSecurityEventsByUserID gets security events by type and time range
func (r *SecurityRepository) GetSecurityEventsByUserID(ctx context.Context, userID uuid.UUID, eventType string, startTime, endTime time.Time) ([]domain.SecurityEvent, error) {
	params := db.GetSecurityEventsByUserIDAndTypeParams{
		UserID:      userID,
		EventType:   eventType,
		Timestamp:   pgtype.Timestamp{Time: startTime, Valid: true},
		Timestamp_2: pgtype.Timestamp{Time: endTime, Valid: true},
	}

	events, err := r.store.GetSecurityEventsByUserIDAndType(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get security events: %w", err)
	}

	result := make([]domain.SecurityEvent, len(events))
	for i, event := range events {
		result[i] = domain.SecurityEvent{
			ID:        event.ID,
			UserID:    event.UserID,
			EventType: event.EventType,
			IPAddress: event.IpAddress,
			UserAgent: event.UserAgent.String,
			Timestamp: event.Timestamp.Time,
			// Parse metadata from JSON
			Metadata: make(map[string]interface{}),
		}
	}

	return result, nil
}
