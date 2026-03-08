package authrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/MicahParks/keyfunc"
	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	userdomain "github.com/demola234/defifundr/internal/features/user/domain"
	authdomain "github.com/demola234/defifundr/internal/features/auth/domain"
	authport "github.com/demola234/defifundr/internal/features/auth/port"
	"github.com/demola234/defifundr/pkg/tracing"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// ────────────────────────────── SessionRepository ──────────────────────────────

// SessionRepository implements authport.SessionRepository.
type SessionRepository struct {
	store db.Queries
	pool  db.DBTX
}

// NewSessionRepository creates a new SessionRepository.
func NewSessionRepository(store db.Queries, pool db.DBTX) authport.SessionRepository {
	return &SessionRepository{store: store, pool: pool}
}

func (r *SessionRepository) CreateSession(ctx context.Context, session authdomain.Session) (*authdomain.Session, error) {
	ctx, span := tracing.Tracer("auth-session-repo").Start(ctx, "CreateSession")
	defer span.End()

	params := db.CreateSessionParams{
		ID:            session.ID,
		UserID:        session.UserID,
		RefreshToken:  pgtype.Text{String: session.RefreshToken, Valid: true},
		UserAgent:     pgtype.Text{String: session.UserAgent, Valid: true},
		ClientIp:      pgtype.Text{String: session.ClientIP, Valid: true},
		UserLoginType: pgtype.Text{String: session.UserLoginType, Valid: true},
		MfaVerified:   session.MFAEnabled,
		IsBlocked:     session.IsBlocked,
		ExpiresAt:     pgtype.Timestamptz{Time: session.ExpiresAt, Valid: true},
	}
	if session.WebOAuthClientID != nil {
		params.WebOauthClientID = pgtype.Text{String: *session.WebOAuthClientID, Valid: true}
	}
	if session.OAuthIDToken != nil {
		params.OauthIDToken = pgtype.Text{String: *session.OAuthIDToken, Valid: true}
	}
	params.OauthAccessToken = pgtype.Text{String: session.OAuthAccessToken, Valid: true}

	dbSession, err := r.store.CreateSession(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return mapDbSessionToDomain(dbSession), nil
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, id uuid.UUID) (*authdomain.Session, error) {
	ctx, span := tracing.Tracer("auth-session-repo").Start(ctx, "GetSessionByID")
	defer span.End()

	dbSession, err := r.store.GetSessionByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by ID: %w", err)
	}
	return mapDbSessionToDomain(dbSession), nil
}

func (r *SessionRepository) GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*authdomain.Session, error) {
	ctx, span := tracing.Tracer("auth-session-repo").Start(ctx, "GetSessionByRefreshToken")
	defer span.End()
	const q = `SELECT id, user_id, refresh_token, user_agent, client_ip, last_used_at, web_oauth_client_id, oauth_access_token, oauth_id_token, user_login_type, mfa_verified, is_blocked, expires_at, created_at FROM sessions WHERE refresh_token = $1 AND is_blocked = FALSE LIMIT 1`
	var s db.Sessions
	row := r.pool.QueryRow(ctx, q, refreshToken)
	err := row.Scan(
		&s.ID, &s.UserID, &s.RefreshToken, &s.UserAgent, &s.ClientIp,
		&s.LastUsedAt, &s.WebOauthClientID, &s.OauthAccessToken, &s.OauthIDToken,
		&s.UserLoginType, &s.MfaVerified, &s.IsBlocked, &s.ExpiresAt, &s.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by refresh token: %w", err)
	}
	return mapDbSessionToDomain(s), nil
}

func (r *SessionRepository) GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]authdomain.Session, error) {
	ctx, span := tracing.Tracer("auth-session-repo").Start(ctx, "GetActiveSessionsByUserID")
	defer span.End()

	dbSessions, err := r.store.GetSessionsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}
	sessions := make([]authdomain.Session, len(dbSessions))
	for i, s := range dbSessions {
		sessions[i] = *mapDbSessionToDomain(s)
	}
	return sessions, nil
}

func (r *SessionRepository) UpdateRefreshToken(ctx context.Context, id uuid.UUID, newToken string) (*authdomain.Session, error) {
	ctx, span := tracing.Tracer("auth-session-repo").Start(ctx, "UpdateRefreshToken")
	defer span.End()
	const q = `UPDATE sessions SET refresh_token = $1, last_used_at = NOW() WHERE id = $2 RETURNING id, user_id, refresh_token, user_agent, client_ip, last_used_at, web_oauth_client_id, oauth_access_token, oauth_id_token, user_login_type, mfa_verified, is_blocked, expires_at, created_at`
	var s db.Sessions
	row := r.pool.QueryRow(ctx, q, newToken, id)
	err := row.Scan(
		&s.ID, &s.UserID, &s.RefreshToken, &s.UserAgent, &s.ClientIp,
		&s.LastUsedAt, &s.WebOauthClientID, &s.OauthAccessToken, &s.OauthIDToken,
		&s.UserLoginType, &s.MfaVerified, &s.IsBlocked, &s.ExpiresAt, &s.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update refresh token: %w", err)
	}
	return mapDbSessionToDomain(s), nil
}

func (r *SessionRepository) UpdateSession(ctx context.Context, session authdomain.Session) error {
	ctx, span := tracing.Tracer("auth-session-repo").Start(ctx, "UpdateSession")
	defer span.End()
	return r.store.UpdateSessionLastUsed(ctx, session.ID)
}

func (r *SessionRepository) BlockSession(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("auth-session-repo").Start(ctx, "BlockSession")
	defer span.End()
	if err := r.store.RevokeSession(ctx, id); err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}
	return nil
}

func (r *SessionRepository) BlockAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	ctx, span := tracing.Tracer("auth-session-repo").Start(ctx, "BlockAllUserSessions")
	defer span.End()
	if err := r.store.RevokeAllUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}
	return nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("auth-session-repo").Start(ctx, "DeleteSession")
	defer span.End()
	return r.store.RevokeSession(ctx, id)
}

func mapDbSessionToDomain(s db.Sessions) *authdomain.Session {
	var webOAuthClientID *string
	if s.WebOauthClientID.Valid {
		webOAuthClientID = &s.WebOauthClientID.String
	}
	var oAuthIDToken *string
	if s.OauthIDToken.Valid {
		oAuthIDToken = &s.OauthIDToken.String
	}
	var oAuthAccessToken string
	if s.OauthAccessToken.Valid {
		oAuthAccessToken = s.OauthAccessToken.String
	}
	var expiresAt time.Time
	if s.ExpiresAt.Valid {
		expiresAt = s.ExpiresAt.Time
	}
	var lastUsedAt time.Time
	if s.LastUsedAt.Valid {
		lastUsedAt = s.LastUsedAt.Time
	}
	var createdAt time.Time
	if s.CreatedAt.Valid {
		createdAt = s.CreatedAt.Time
	}
	var mfaEnabled bool
	if s.MfaVerified.Valid {
		mfaEnabled = s.MfaVerified.Bool
	}
	var isBlocked bool
	if s.IsBlocked.Valid {
		isBlocked = s.IsBlocked.Bool
	}
	return &authdomain.Session{
		ID:               s.ID,
		UserID:           s.UserID,
		RefreshToken:     s.RefreshToken.String,
		WebOAuthClientID: webOAuthClientID,
		OAuthIDToken:     oAuthIDToken,
		OAuthAccessToken: oAuthAccessToken,
		UserAgent:        s.UserAgent.String,
		UserLoginType:    s.UserLoginType.String,
		MFAEnabled:       mfaEnabled,
		ClientIP:         s.ClientIp.String,
		IsBlocked:        isBlocked,
		ExpiresAt:        expiresAt,
		LastUsedAt:       lastUsedAt,
		CreatedAt:        createdAt,
	}
}

// ────────────────────────────── OAuthRepository ──────────────────────────────

// OAuthRepository implements authport.OAuthRepository.
type OAuthRepository struct {
	store       db.Queries
	logger      logging.Logger
	jwksCache   map[string]*keyfunc.JWKS
	cacheExpiry map[string]time.Time
	cacheMutex  sync.RWMutex
}

// NewOAuthRepository creates a new OAuthRepository.
func NewOAuthRepository(store db.Queries, logger logging.Logger) authport.OAuthRepository {
	return &OAuthRepository{
		store:       store,
		logger:      logger,
		jwksCache:   make(map[string]*keyfunc.JWKS),
		cacheExpiry: make(map[string]time.Time),
	}
}

func (r *OAuthRepository) getJWKS(jwksURL string) (*keyfunc.JWKS, error) {
	r.cacheMutex.RLock()
	jwks, found := r.jwksCache[jwksURL]
	expiry, _ := r.cacheExpiry[jwksURL]
	r.cacheMutex.RUnlock()

	if found && time.Now().Before(expiry) {
		return jwks, nil
	}

	r.cacheMutex.Lock()
	defer r.cacheMutex.Unlock()

	jwks, found = r.jwksCache[jwksURL]
	expiry, _ = r.cacheExpiry[jwksURL]
	if found && time.Now().Before(expiry) {
		return jwks, nil
	}

	options := keyfunc.Options{
		RefreshInterval: time.Hour,
		RefreshErrorHandler: func(err error) {
			r.logger.Error("Error refreshing JWKS", err, map[string]any{"jwks_url": jwksURL})
		},
	}
	newJWKS, err := keyfunc.Get(jwksURL, options)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %v", err)
	}
	r.jwksCache[jwksURL] = newJWKS
	r.cacheExpiry[jwksURL] = time.Now().Add(time.Hour)
	return newJWKS, nil
}

func (r *OAuthRepository) ValidateWebAuthToken(ctx context.Context, tokenString string) (*authdomain.Web3AuthClaims, error) {
	_, span := tracing.Tracer("auth-oauth-repo").Start(ctx, "ValidateWebAuthToken")
	defer span.End()

	jwks, err := r.getJWKS("https://api-auth.web3auth.io/jwks")
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %v", err)
	}

	claims := &authdomain.Web3AuthClaims{}
	parser := jwtv4.NewParser(jwtv4.WithValidMethods([]string{"ES256"}))
	tok, err := parser.ParseWithClaims(tokenString, claims, jwks.Keyfunc)
	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			return nil, errors.New("token has expired, please re-authenticate")
		}
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}
	if !tok.Valid {
		return nil, errors.New("invalid token")
	}
	if claims.Verifier == "" || claims.VerifierID == "" {
		return nil, errors.New("missing required Web3Auth claims")
	}
	if claims.Issuer != "https://api-auth.web3auth.io" {
		return nil, fmt.Errorf("invalid issuer: %v", claims.Issuer)
	}
	r.logger.Info("Successfully validated Web3Auth token", map[string]any{
		"email":        claims.Email,
		"verifier":     claims.Verifier,
		"wallet_count": len(claims.Wallets),
	})
	return claims, nil
}

func (r *OAuthRepository) GetUserInfoFromProviderToken(ctx context.Context, provider, token string) (*userdomain.User, error) {
	ctx, span := tracing.Tracer("auth-oauth-repo").Start(ctx, "GetUserInfoFromProviderToken")
	defer span.End()

	if provider == string(authdomain.Web3AuthProvider) {
		claims, err := r.ValidateWebAuthToken(ctx, token)
		if err != nil {
			return nil, err
		}
		firstName, lastName := oauthExtractName(claims)
		profileImage := claims.ProfileImage
		return &userdomain.User{
			Email:          claims.Email,
			FirstName:      firstName,
			LastName:       lastName,
			ProfilePicture: &profileImage,
			AuthProvider:   string(oauthMapVerifier(claims.Verifier)),
			ProviderID:     claims.VerifierID,
		}, nil
	}
	return nil, fmt.Errorf("unsupported provider: %s", provider)
}

// oauthExtractName extracts first/last name from Web3Auth claims.
func oauthExtractName(claims *authdomain.Web3AuthClaims) (string, string) {
	if claims.Name == "" {
		return "User", ""
	}
	parts := strings.Split(claims.Name, " ")
	first := parts[0]
	var last string
	if len(parts) > 1 {
		last = strings.Join(parts[1:], " ")
	}
	return first, last
}

// oauthMapVerifier maps a Web3Auth verifier string to a authdomain.AuthProvider.
func oauthMapVerifier(verifier string) authdomain.AuthProvider {
	lv := strings.ToLower(verifier)
	switch {
	case strings.Contains(lv, "google"):
		return authdomain.GoogleProvider
	case strings.Contains(lv, "facebook"):
		return authdomain.FacebookProvider
	case strings.Contains(lv, "apple"):
		return authdomain.AppleProvider
	case strings.Contains(lv, "twitter"):
		return authdomain.TwitterProvider
	case strings.Contains(lv, "discord"):
		return authdomain.DiscordProvider
	}
	return authdomain.Web3AuthProvider
}

// ────────────────────────────── WalletRepository ──────────────────────────────

// WalletRepository implements authport.WalletRepository.
type WalletRepository struct {
	store db.Queries
	pool  db.DBTX
}

// NewWalletRepository creates a new WalletRepository.
func NewWalletRepository(store db.Queries, pool db.DBTX) authport.WalletRepository {
	return &WalletRepository{store: store, pool: pool}
}

func (r *WalletRepository) CreateWallet(ctx context.Context, wallet authdomain.UserWallet) error {
	ctx, span := tracing.Tracer("auth-wallet-repo").Start(ctx, "CreateWallet")
	defer span.End()

	params := db.CreateUserWalletParams{
		ID:            wallet.ID,
		UserID:        wallet.UserID,
		WalletAddress: wallet.Address,
		WalletType:    wallet.Type,
		IsDefault:     wallet.IsDefault,
	}
	_, err := r.store.CreateUserWallet(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create user wallet: %w", err)
	}
	return nil
}

func (r *WalletRepository) GetWalletByAddress(ctx context.Context, address string) (*authdomain.UserWallet, error) {
	ctx, span := tracing.Tracer("auth-wallet-repo").Start(ctx, "GetWalletByAddress")
	defer span.End()
	const q = `SELECT uw.id, uw.user_id, uw.wallet_address, uw.wallet_type, uw.chain_id, uw.is_default, uw.is_verified, uw.verification_method, uw.verified_at, uw.nickname, uw.created_at, uw.updated_at, sn.name as network_name FROM user_wallets uw JOIN supported_networks sn ON uw.chain_id = sn.chain_id WHERE uw.wallet_address = $1 LIMIT 1`
	var w db.GetUserWalletsByUserRow
	row := r.pool.QueryRow(ctx, q, address)
	err := row.Scan(
		&w.ID, &w.UserID, &w.WalletAddress, &w.WalletType, &w.ChainID,
		&w.IsDefault, &w.IsVerified, &w.VerificationMethod, &w.VerifiedAt,
		&w.Nickname, &w.CreatedAt, &w.UpdatedAt, &w.NetworkName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet by address: %w", err)
	}
	result := mapWalletToDomain(w)
	return &result, nil
}

func (r *WalletRepository) GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]authdomain.UserWallet, error) {
	ctx, span := tracing.Tracer("auth-wallet-repo").Start(ctx, "GetWalletsByUserID")
	defer span.End()

	wallets, err := r.store.GetUserWalletsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}
	result := make([]authdomain.UserWallet, len(wallets))
	for i, w := range wallets {
		result[i] = mapWalletToDomain(w)
	}
	return result, nil
}

func (r *WalletRepository) UpdateWallet(ctx context.Context, wallet authdomain.UserWallet) error {
	ctx, span := tracing.Tracer("auth-wallet-repo").Start(ctx, "UpdateWallet")
	defer span.End()
	_, err := r.store.UpdateUserWallet(ctx, db.UpdateUserWalletParams{
		WalletType: wallet.Type,
		IsDefault:  pgtype.Bool{Bool: wallet.IsDefault, Valid: true},
		ID:         wallet.ID,
	})
	return err
}

func (r *WalletRepository) DeleteWallet(ctx context.Context, walletID uuid.UUID) error {
	if err := r.store.DeleteUserWallet(ctx, walletID); err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}
	return nil
}

func mapWalletToDomain(w db.GetUserWalletsByUserRow) authdomain.UserWallet {
	var isDefault bool
	if w.IsDefault.Valid {
		isDefault = w.IsDefault.Bool
	}
	result := authdomain.UserWallet{
		ID:        w.ID,
		UserID:    w.UserID,
		Address:   w.WalletAddress,
		Type:      w.WalletType,
		Chain:     strconv.Itoa(int(w.ChainID)),
		IsDefault: isDefault,
	}
	if w.CreatedAt.Valid {
		result.CreatedAt = w.CreatedAt.Time
	}
	if w.UpdatedAt.Valid {
		result.UpdatedAt = w.UpdatedAt.Time
	}
	return result
}

// ────────────────────────────── SecurityRepository ──────────────────────────────

// SecurityRepository implements authport.SecurityRepository.
type SecurityRepository struct {
	store db.Queries
}

// NewSecurityRepository creates a new SecurityRepository.
func NewSecurityRepository(store db.Queries) authport.SecurityRepository {
	return &SecurityRepository{store: store}
}

func (r *SecurityRepository) LogSecurityEvent(ctx context.Context, event authdomain.SecurityEvent) error {
	ctx, span := tracing.Tracer("auth-security-repo").Start(ctx, "LogSecurityEvent")
	defer span.End()

	metadataBytes, _ := json.Marshal(event.Metadata)
	params := db.CreateSecurityEventParams{
		UserID:    pgtype.UUID{Bytes: event.UserID, Valid: true},
		EventType: event.EventType,
		Severity:  "info",
		IpAddress: pgtype.Text{String: event.IPAddress, Valid: true},
		UserAgent: pgtype.Text{String: event.UserAgent, Valid: true},
		Metadata:  metadataBytes,
	}
	_, err := r.store.CreateSecurityEvent(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to log security event: %w", err)
	}
	return nil
}

func (r *SecurityRepository) GetRecentLoginsByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]authdomain.SecurityEvent, error) {
	ctx, span := tracing.Tracer("auth-security-repo").Start(ctx, "GetRecentLoginsByUserID")
	defer span.End()
	rows, err := r.store.GetSecurityEventsByUser(ctx, db.GetSecurityEventsByUserParams{
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		OffsetVal: 0,
		LimitVal:  int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get recent logins: %w", err)
	}
	return mapSecurityEvents(rows), nil
}

func (r *SecurityRepository) GetSecurityEventsByUserID(ctx context.Context, userID uuid.UUID, _ string, _, _ time.Time) ([]authdomain.SecurityEvent, error) {
	ctx, span := tracing.Tracer("auth-security-repo").Start(ctx, "GetSecurityEventsByUserID")
	defer span.End()
	rows, err := r.store.GetSecurityEventsByUser(ctx, db.GetSecurityEventsByUserParams{
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		OffsetVal: 0,
		LimitVal:  50,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get security events: %w", err)
	}
	return mapSecurityEvents(rows), nil
}

func mapSecurityEvents(rows []db.SecurityEvents) []authdomain.SecurityEvent {
	events := make([]authdomain.SecurityEvent, len(rows))
	for i, s := range rows {
		e := authdomain.SecurityEvent{
			ID:        s.ID,
			UserID:    s.UserID.Bytes,
			EventType: s.EventType,
			IPAddress: s.IpAddress.String,
			UserAgent: s.UserAgent.String,
			Metadata:  map[string]any{},
		}
		if s.CreatedAt.Valid {
			e.Timestamp = s.CreatedAt.Time
		}
		events[i] = e
	}
	return events
}

// ────────────────────────────── OTPRepository ──────────────────────────────

// OTPRepository implements authport.OTPRepository using an in-memory store.
type OTPRepository struct {
	store    db.Queries
	mu       sync.Mutex
	otpStore map[string]*authdomain.OTPVerification
}

// NewOTPRepository creates a new OTPRepository.
func NewOTPRepository(store db.Queries) authport.OTPRepository {
	return &OTPRepository{
		store:    store,
		otpStore: make(map[string]*authdomain.OTPVerification),
	}
}

func (r *OTPRepository) CreateOTP(_ context.Context, otp authdomain.OTPVerification) (*authdomain.OTPVerification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if otp.ID == uuid.Nil {
		otp.ID = uuid.New()
	}
	if otp.CreatedAt.IsZero() {
		otp.CreatedAt = time.Now()
	}
	if otp.ExpiresAt.IsZero() {
		otp.ExpiresAt = time.Now().Add(5 * time.Minute)
	}
	key := otp.UserID.String() + ":" + string(otp.Purpose)
	r.otpStore[key] = &otp
	return &otp, nil
}

func (r *OTPRepository) GetOTPByUserIDAndPurpose(_ context.Context, userID uuid.UUID, purpose authdomain.OTPPurpose) (*authdomain.OTPVerification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	key := userID.String() + ":" + string(purpose)
	record, ok := r.otpStore[key]
	if !ok {
		return nil, errors.New("OTP not found")
	}
	if time.Now().After(record.ExpiresAt) {
		return nil, errors.New("OTP has expired")
	}
	return record, nil
}

func (r *OTPRepository) VerifyOTP(_ context.Context, userID uuid.UUID, code string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, record := range r.otpStore {
		if record.UserID == userID {
			if record.OTPCode == code || record.HashedOTP == code {
				record.IsVerified = true
				return nil
			}
			return errors.New("invalid OTP code")
		}
	}
	return errors.New("OTP not found")
}

func (r *OTPRepository) IncrementAttempts(_ context.Context, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, record := range r.otpStore {
		if record.UserID == userID {
			record.AttemptsMade++
			return nil
		}
	}
	return errors.New("OTP not found")
}
