// repositories/session_repository.go (enhanced)
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

type SessionRepository struct {
	store db.Queries
}

func NewSessionRepository(store db.Queries) *SessionRepository {
	return &SessionRepository{
		store: store,
	}
}

// CreateSession creates a new session
func (r *SessionRepository) CreateSession(ctx context.Context, session domain.Session) (*domain.Session, error) {
	params := db.CreateSessionParams{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshToken:     session.RefreshToken,
		OauthAccessToken: pgtype.Text{String: session.OAuthAccessToken, Valid: true},
		UserAgent:        session.UserAgent,
		UserLoginType:    session.UserLoginType,
		MfaEnabled:       session.MFAEnabled,
		ClientIp:         session.ClientIP,
		IsBlocked:        session.IsBlocked,
		LastUsedAt:       pgtype.Timestamp{Time: time.Now(), Valid: true},
		ExpiresAt:        pgtype.Timestamp{Time: session.ExpiresAt, Valid: true},
	}

	// Add WebOAuthClientID if provided
	if session.WebOAuthClientID != nil {
		params.WebOauthClientID = pgtype.Text{String: *session.WebOAuthClientID, Valid: true}
	}

	// Add OAuthIDToken if provided
	if session.OAuthIDToken != nil {
		params.OauthIDToken = pgtype.Text{String: *session.OAuthIDToken, Valid: true}
	}

	dbSession, err := r.store.CreateSession(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return mapDbSessionToDomain(dbSession), nil
}

// GetSessionByID gets a session by ID
func (r *SessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	dbSession, err := r.store.GetSessionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by ID: %w", err)
	}

	return mapDbSessionToDomain(dbSession), nil
}

// GetSessionByRefreshToken gets a session by refresh token
func (r *SessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	dbSession, err := r.store.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	return mapDbSessionToDomain(dbSession), nil
}

// GetActiveSessionsByUserID gets all active sessions for a user
func (r *SessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	dbSessions, err := r.store.GetActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	sessions := make([]domain.Session, len(dbSessions))
	for i, dbSession := range dbSessions {
		sessions[i] = *mapDbSessionToDomain(dbSession)
	}

	return sessions, nil
}

// UpdateRefreshToken updates a session's refresh token
func (r *SessionRepository) UpdateRefreshToken(ctx context.Context, sessionID uuid.UUID, refreshToken string) (*domain.Session, error) {
	params := db.UpdateSessionRefreshTokenParams{
		ID:           sessionID,
		RefreshToken: refreshToken,
		LastUsedAt:   pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	dbSession, err := r.store.UpdateSessionRefreshToken(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	return mapDbSessionToDomain(dbSession), nil
}

// UpdateSession updates a session
func (r *SessionRepository) UpdateSession(ctx context.Context, session domain.Session) error {
	params := db.UpdateSessionParams{
		ID:         session.ID,
		IsBlocked:  session.IsBlocked,
		MfaEnabled: session.MFAEnabled,
		LastUsedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	_, err := r.store.UpdateSession(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// BlockSession blocks a session
func (r *SessionRepository) BlockSession(ctx context.Context, id uuid.UUID) error {
	err := r.store.BlockSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to block session: %w", err)
	}

	return nil
}

// BlockAllUserSessions blocks all sessions for a user
func (r *SessionRepository) BlockAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	err := r.store.BlockAllUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to block all user sessions: %w", err)
	}

	return nil
}

// DeleteSession deletes a session
func (r *SessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	err := r.store.DeleteSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

// Helper function to map database session model to domain model
func mapDbSessionToDomain(dbSession db.Sessions) *domain.Session {
	var webOAuthClientID *string
	if dbSession.WebOauthClientID.Valid {
		webOAuthClientID = &dbSession.WebOauthClientID.String
	}

	var oAuthIDToken *string
	if dbSession.OauthIDToken.Valid {
		oAuthIDToken = &dbSession.OauthIDToken.String
	}

	var oAuthAccessToken string
	if dbSession.OauthAccessToken.Valid {
		oAuthAccessToken = dbSession.OauthAccessToken.String
	}

	var expiresAt time.Time
	if dbSession.ExpiresAt.Valid {
		expiresAt = dbSession.ExpiresAt.Time
	}

	var lastUsedAt time.Time
	if dbSession.LastUsedAt.Valid {
		lastUsedAt = dbSession.LastUsedAt.Time
	}

	return &domain.Session{
		ID:               dbSession.ID,
		UserID:           dbSession.UserID,
		RefreshToken:     dbSession.RefreshToken,
		WebOAuthClientID: webOAuthClientID,
		OAuthIDToken:     oAuthIDToken,
		OAuthAccessToken: oAuthAccessToken,
		UserAgent:        dbSession.UserAgent,
		UserLoginType:    dbSession.UserLoginType,
		MFAEnabled:       dbSession.MfaEnabled,
		ClientIP:         dbSession.ClientIp,
		IsBlocked:        dbSession.IsBlocked,
		ExpiresAt:        expiresAt,
		LastUsedAt:       lastUsedAt,
		CreatedAt:        dbSession.CreatedAt.Time,
	}
}
