package authport

import (
	"context"
	"time"

	userdomain "github.com/demola234/defifundr/internal/features/user/domain"
	authdomain "github.com/demola234/defifundr/internal/features/auth/domain"
	userport "github.com/demola234/defifundr/internal/features/user/port"
	"github.com/google/uuid"
)

// AuthService defines all authentication use cases.
type AuthService interface {
	AuthenticateWithWeb3(ctx context.Context, webAuthToken, userAgent, clientIP string) (*userdomain.User, *authdomain.Session, error)
	Login(ctx context.Context, email string, user userdomain.User, password string) (*userdomain.User, error)
	RegisterUser(ctx context.Context, user userdomain.User, password string) (*userdomain.User, error)
	RegisterPersonalDetails(ctx context.Context, user userdomain.User) (*userdomain.User, error)
	RegisterAddressDetails(ctx context.Context, user userdomain.User) (*userdomain.User, error)
	RegisterBusinessDetails(ctx context.Context, companyInfo userdomain.CompanyInfo) (*userdomain.CompanyInfo, error)
	GetProfileCompletionStatus(ctx context.Context, userID uuid.UUID) (*authdomain.ProfileCompletion, error)
	SetupMFA(ctx context.Context, userID uuid.UUID) (string, error)
	VerifyMFA(ctx context.Context, userID uuid.UUID, code string) (bool, error)
	LinkWallet(ctx context.Context, userID uuid.UUID, walletAddress, walletType, chain string) error
	GetUserWallets(ctx context.Context, userID uuid.UUID) ([]authdomain.UserWallet, error)
	CreateSession(ctx context.Context, userID uuid.UUID, userAgent, clientIP, webOAuthClientID, email, loginType string) (*authdomain.Session, error)
	GetActiveDevices(ctx context.Context, userID uuid.UUID) ([]authdomain.DeviceInfo, error)
	RevokeSession(ctx context.Context, userID, sessionID uuid.UUID) error
	Logout(ctx context.Context, sessionID uuid.UUID) error
	RefreshToken(ctx context.Context, refreshToken, userAgent, clientIP string) (*authdomain.Session, string, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*userdomain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*userdomain.User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	LogSecurityEvent(ctx context.Context, eventType string, userID uuid.UUID, metadata map[string]interface{}) error
	GetUserSecurityEvents(ctx context.Context, userID uuid.UUID) ([]authdomain.SecurityEvent, error)
	InitiatePasswordReset(ctx context.Context, email string) error
	VerifyResetOTP(ctx context.Context, email, otp string) error
	ResetPassword(ctx context.Context, email, otp, newPassword string) error
	GetUserRepository() userport.UserRepository
}

// SessionRepository defines data access for sessions.
type SessionRepository interface {
	CreateSession(ctx context.Context, session authdomain.Session) (*authdomain.Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (*authdomain.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*authdomain.Session, error)
	GetActiveSessionsByUserID(ctx context.Context, userID uuid.UUID) ([]authdomain.Session, error)
	UpdateRefreshToken(ctx context.Context, sessionID uuid.UUID, refreshToken string) (*authdomain.Session, error)
	UpdateSession(ctx context.Context, session authdomain.Session) error
	BlockSession(ctx context.Context, id uuid.UUID) error
	BlockAllUserSessions(ctx context.Context, userID uuid.UUID) error
	DeleteSession(ctx context.Context, id uuid.UUID) error
}

// OAuthRepository defines operations for OAuth/Web3Auth.
type OAuthRepository interface {
	ValidateWebAuthToken(ctx context.Context, tokenString string) (*authdomain.Web3AuthClaims, error)
	GetUserInfoFromProviderToken(ctx context.Context, provider, token string) (*userdomain.User, error)
}

// WalletRepository defines data access for wallets.
type WalletRepository interface {
	CreateWallet(ctx context.Context, wallet authdomain.UserWallet) error
	GetWalletByAddress(ctx context.Context, address string) (*authdomain.UserWallet, error)
	GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]authdomain.UserWallet, error)
	UpdateWallet(ctx context.Context, wallet authdomain.UserWallet) error
	DeleteWallet(ctx context.Context, walletID uuid.UUID) error
}

// SecurityRepository defines data access for security events.
type SecurityRepository interface {
	LogSecurityEvent(ctx context.Context, event authdomain.SecurityEvent) error
	GetRecentLoginsByUserID(ctx context.Context, userID uuid.UUID, limit int) ([]authdomain.SecurityEvent, error)
	GetSecurityEventsByUserID(ctx context.Context, userID uuid.UUID, eventType string, startTime, endTime time.Time) ([]authdomain.SecurityEvent, error)
}

// OTPRepository defines data access for OTP verifications.
type OTPRepository interface {
	CreateOTP(ctx context.Context, otp authdomain.OTPVerification) (*authdomain.OTPVerification, error)
	GetOTPByUserIDAndPurpose(ctx context.Context, userID uuid.UUID, purpose authdomain.OTPPurpose) (*authdomain.OTPVerification, error)
	VerifyOTP(ctx context.Context, id uuid.UUID, code string) error
	IncrementAttempts(ctx context.Context, id uuid.UUID) error
}

// EmailService defines email operations for auth.
type EmailService interface {
	SendPasswordResetEmail(ctx context.Context, email, name, otpCode string) error
	SendBatchUpdate(ctx context.Context, emails []string, subject, message string) error
}
