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
}

// NewSessionRepository creates a new SessionRepository.
func NewSessionRepository(store db.Queries) authport.SessionRepository {
	return &SessionRepository{store: store}
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

// GetSessionByRefreshToken is not yet implemented (missing SQLC query).
func (r *SessionRepository) GetSessionByRefreshToken(_ context.Context, _ string) (*authdomain.Session, error) {
	return nil, errors.New("not implemented: GetSessionByRefreshToken")
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

// UpdateRefreshToken is not yet implemented (missing SQLC query).
func (r *SessionRepository) UpdateRefreshToken(_ context.Context, _ uuid.UUID, _ string) (*authdomain.Session, error) {
	return nil, errors.New("not implemented: UpdateRefreshToken")
}

// UpdateSession is not yet implemented (missing SQLC query).
func (r *SessionRepository) UpdateSession(_ context.Context, _ authdomain.Session) error {
	return errors.New("not implemented: UpdateSession")
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

// DeleteSession is not yet implemented (missing SQLC query).
func (r *SessionRepository) DeleteSession(_ context.Context, _ uuid.UUID) error {
	return errors.New("not implemented: DeleteSession")
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
			r.logger.Error("Error refreshing JWKS", err, map[string]interface{}{"jwks_url": jwksURL})
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
	r.logger.Info("Successfully validated Web3Auth token", map[string]interface{}{
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
}

// NewWalletRepository creates a new WalletRepository.
func NewWalletRepository(store db.Queries) authport.WalletRepository {
	return &WalletRepository{store: store}
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

// GetWalletByAddress is not yet implemented (missing SQLC query).
func (r *WalletRepository) GetWalletByAddress(_ context.Context, _ string) (*authdomain.UserWallet, error) {
	return nil, errors.New("not implemented: GetWalletByAddress")
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

// UpdateWallet is not yet implemented (missing SQLC query).
func (r *WalletRepository) UpdateWallet(_ context.Context, _ authdomain.UserWallet) error {
	return errors.New("not implemented: UpdateWallet")
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

// GetRecentLoginsByUserID is not yet implemented (missing SQLC query).
func (r *SecurityRepository) GetRecentLoginsByUserID(_ context.Context, _ uuid.UUID, _ int) ([]authdomain.SecurityEvent, error) {
	return nil, errors.New("not implemented: GetRecentLoginsByUserID")
}

// GetSecurityEventsByUserID is not yet implemented (missing SQLC query).
func (r *SecurityRepository) GetSecurityEventsByUserID(_ context.Context, _ uuid.UUID, _ string, _, _ time.Time) ([]authdomain.SecurityEvent, error) {
	return nil, errors.New("not implemented: GetSecurityEventsByUserID")
}

// ────────────────────────────── OTPRepository ──────────────────────────────

// OTPRepository implements authport.OTPRepository.
type OTPRepository struct {
	store db.Queries
}

// NewOTPRepository creates a new OTPRepository.
func NewOTPRepository(store db.Queries) authport.OTPRepository {
	return &OTPRepository{store: store}
}

// CreateOTP is not yet implemented (missing SQLC queries for OTP).
func (r *OTPRepository) CreateOTP(_ context.Context, _ authdomain.OTPVerification) (*authdomain.OTPVerification, error) {
	return nil, errors.New("not implemented: CreateOTP")
}

// GetOTPByUserIDAndPurpose is not yet implemented (missing SQLC queries for OTP).
func (r *OTPRepository) GetOTPByUserIDAndPurpose(_ context.Context, _ uuid.UUID, _ authdomain.OTPPurpose) (*authdomain.OTPVerification, error) {
	return nil, errors.New("not implemented: GetOTPByUserIDAndPurpose")
}

// VerifyOTP is not yet implemented (missing SQLC queries for OTP).
func (r *OTPRepository) VerifyOTP(_ context.Context, _ uuid.UUID, _ string) error {
	return errors.New("not implemented: VerifyOTP")
}

// IncrementAttempts is not yet implemented (missing SQLC queries for OTP).
func (r *OTPRepository) IncrementAttempts(_ context.Context, _ uuid.UUID) error {
	return errors.New("not implemented: IncrementAttempts")
}
