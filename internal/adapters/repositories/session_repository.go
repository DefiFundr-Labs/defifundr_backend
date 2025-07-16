package repositories

import (
	"context"
	"fmt"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"
	"github.com/demola234/defifundr/pkg/tracing"
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

func (r *SessionRepository) CreateSession(ctx context.Context, session domain.Session) (*domain.Session, error) {
	ctx, span := tracing.Tracer("session-repository").Start(ctx, "CreateSession")
	defer span.End()

	params := db.CreateSessionParams{
		ID:               session.ID,
		UserID:           session.UserID,
		RefreshToken:     toPgTextPtr(&session.RefreshToken),
		UserAgent:        toPgTextPtr(&session.UserAgent),
		ClientIp:         toPgTextPtr(&session.ClientIP),
		LastUsedAt:       toPgTimestamptzPtr(&session.LastUsedAt),
		WebOauthClientID: toPgTextPtr(session.WebOAuthClientID),  // This is already *string
		OauthAccessToken: toPgTextPtr(&session.OAuthAccessToken),
		OauthIDToken:     toPgTextPtr(session.OAuthIDToken),      // This is already *string
		UserLoginType:    toPgTextPtr(&session.UserLoginType),
		MfaVerified:      session.MFAVerified,
		IsBlocked:        session.IsBlocked,
		ExpiresAt:        toPgTimestamptzPtr(session.ExpiresAt),         // This is already *time.Time
		CreatedAt:        toPgTimestamptzPtr(&session.CreatedAt),
	}

	dbSession, err := r.store.CreateSession(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return mapDBSessionToDomain(dbSession), nil
}

// Fixed mapDBSessionToDomain function
func mapDBSessionToDomain(dbSession db.Sessions) *domain.Session {
	return &domain.Session{
		ID:                dbSession.ID,
		UserID:            dbSession.UserID,
		RefreshToken:      getTextString(dbSession.RefreshToken),
		UserAgent:         getTextString(dbSession.UserAgent),
		ClientIP:          getTextString(dbSession.ClientIp),
		LastUsedAt:        getTimestamptzTime(dbSession.LastUsedAt),
		WebOAuthClientID:  getTextStringPtr(dbSession.WebOauthClientID),
		OAuthAccessToken:  getTextString(dbSession.OauthAccessToken),
		OAuthIDToken:      getTextStringPtr(dbSession.OauthIDToken),
		UserLoginType:     getTextString(dbSession.UserLoginType),
		MFAVerified:       getBool(dbSession.MfaVerified),        // Fixed: use getBool helper
		IsBlocked:         getBool(dbSession.IsBlocked),          // Fixed: use getBool helper
		ExpiresAt:         getTimestamptz(dbSession.ExpiresAt),   // Fixed: use getTimestamptz helper
		CreatedAt:         getTimestamptzTime(dbSession.CreatedAt),
	}
}

// GetSessionByID gets a session by ID
func (r *SessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error) {
	ctx, span := tracing.Tracer("session-repository").Start(ctx, "GetSessionByID")
	defer span.End()

	dbSession, err := r.store.GetSessionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by ID: %w", err)
	}

	return mapDBSessionToDomain(dbSession), nil
}

func (r *SessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error) {
	ctx, span := tracing.Tracer("session-repository").Start(ctx, "GetSessionByRefreshToken")
	defer span.End()

	dbSession, err := r.store.GetSessionByRefreshToken(ctx, pgtype.Text{String: refreshToken, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	return mapDBSessionToDomain(dbSession), nil
}

// GetActiveSessionsByUserID gets all active sessions for a user
func (r *SessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Session, error) {
	ctx, span := tracing.Tracer("session-repository").Start(ctx, "GetActiveSessionsByUserID")
	defer span.End()

	dbSessions, err := r.store.GetSessionsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	sessions := make([]domain.Session, len(dbSessions))
	for i, dbSession := range dbSessions {
		sessions[i] = *mapDBSessionToDomain(dbSession)
	}

	return sessions, nil
}

// UpdateRefreshToken updates a session's refresh token
func (r *SessionRepository) UpdateRefreshToken(ctx context.Context, sessionID uuid.UUID, refreshToken string) (*domain.Session, error) {
	ctx, span := tracing.Tracer("session-repository").Start(ctx, "UpdateRefreshToken")
	defer span.End()

	params := db.UpdateRefreshTokenParams{
		ID:           sessionID,
		RefreshToken: pgtype.Text{String: refreshToken, Valid: true},
	}

	dbSession, err := r.store.UpdateRefreshToken(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}

	return mapDBSessionToDomain(dbSession), nil
}

// UpdateSession updates a session
func (r *SessionRepository) UpdateSession(ctx context.Context, session domain.Session) error {
	ctx, span := tracing.Tracer("session-repository").Start(ctx, "UpdateSession")
	defer span.End()

	err := r.store.UpdateSessionLastUsed(ctx, session.ID)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// BlockSession blocks a session
func (r *SessionRepository) BlockSession(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("session-repository").Start(ctx, "BlockSession")
	defer span.End()

	err := r.store.RevokeSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to block session: %w", err)
	}

	return nil
}

// BlockAllUserSessions blocks all sessions for a user
func (r *SessionRepository) BlockAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	ctx, span := tracing.Tracer("session-repository").Start(ctx, "BlockAllUserSessions")
	defer span.End()

	err := r.store.RevokeAllUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to block all user sessions: %w", err)
	}

	return nil
}

// DeleteSession deletes a session
func (r *SessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("session-repository").Start(ctx, "DeleteSession")
	defer span.End()

	err := r.store.DeleteSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}


