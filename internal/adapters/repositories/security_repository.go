package repositories

import (
	"context"
	"fmt"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"
	"github.com/demola234/defifundr/pkg/tracing"
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
	ctx, span := tracing.Tracer("security-repository").Start(ctx, "LogSecurityEvent")
	defer span.End()

	params := db.CreateSecurityEventParams{
		ID:        event.ID,
		UserID:    toPgUUIDPtr(&event.UserID),
		CompanyID: toPgUUIDPtr(&event.CompanyID),
		EventType: event.EventType,
		Severity:  event.Severity,
		IpAddress: toPgTextPtr(&event.IPAddress),
		UserAgent: toPgTextPtr(&event.UserAgent),
		Metadata:  toPgJSONB(event.Metadata),
		CreatedAt: toPgTimestamptzPtr(&event.Timestamp),
	}

	_, err := r.store.CreateSecurityEvent(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to log security event: %w", err)
	}

	return nil
}

// GetRecentLoginsByUserID gets recent login events for a user
func (r *SecurityRepository) GetRecentLoginsByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]domain.SecurityEvent, error) {
	ctx, span := tracing.Tracer("security-repository").Start(ctx, "GetRecentLoginsByUserID")
	defer span.End()

	params := db.GetSecurityEventsByUserParams{
		UserID:   toPgUUID(userID),
		LimitVal: int32(limit),
	}

	events, err := r.store.GetSecurityEventsByUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent logins: %w", err)
	}

	result := make([]domain.SecurityEvent, len(events))
	for i, event := range events {
		result[i] = mapDBSecurityEventToDomain(event)
	}

	return result, nil
}

// GetSecurityEventsByUserID gets security events by type and time range
func (r *SecurityRepository) GetSecurityEventsByUserID(ctx context.Context, userID uuid.UUID, companyID uuid.UUID, eventType string, startTime, endTime time.Time) ([]domain.SecurityEvent, error) {
	ctx, span := tracing.Tracer("security-repository").Start(ctx, "GetSecurityEventsByUserID")
	defer span.End()

	params := db.GetSecurityEventsByTypeParams{
		EventType:  eventType,
		UserID:     userID,
		CompanyID:  companyID, // NULL for company_id
		LimitVal:   100,                       // Default limit
		OffsetVal:  0,                         // Default offset
	}

	events, err := r.store.GetSecurityEventsByType(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get security events: %w", err)
	}

	// Filter by time range (since SQL query doesn't support time filtering)
	var filteredEvents []domain.SecurityEvent
	for _, event := range events {
		eventTime := getTimestamptzTime(event.CreatedAt)
		if eventTime.After(startTime) && eventTime.Before(endTime) {
			filteredEvents = append(filteredEvents, mapDBSecurityEventToDomain(event))
		}
	}

	return filteredEvents, nil
}

// Helper function to map database security event to domain security event
func mapDBSecurityEventToDomain(dbEvent db.SecurityEvents) domain.SecurityEvent {
	return domain.SecurityEvent{
		ID:        dbEvent.ID,
		UserID:    getUUIDFromPgUUID(dbEvent.UserID),
		CompanyID: getUUIDFromPgUUID(dbEvent.CompanyID),
		EventType: dbEvent.EventType,
		Severity:  dbEvent.Severity,
		IPAddress: getTextString(dbEvent.IpAddress),
		UserAgent: getTextString(dbEvent.UserAgent),
		Timestamp: getTimestamptzTime(dbEvent.CreatedAt),
	}
}
