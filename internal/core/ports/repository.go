package ports

import (
	"context"

	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"
)

// UserRepository defines the data access operations for User entities
type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (*domain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
}

// SessionRepository defines the data access operations for Session entities
type SessionRepository interface {
	CreateSession(ctx context.Context, session domain.Session) (*domain.Session, error)
	GetSessionByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	GetSessionByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	BlockSession(ctx context.Context, id uuid.UUID) error
}

// OTPRepository defines the data access operations for OTP entities
type OTPRepository interface {
	CreateOTP(ctx context.Context, otp domain.OTPVerification) (*domain.OTPVerification, error)
	GetOTPByUserIDAndPurpose(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose) (*domain.OTPVerification, error)
	VerifyOTP(ctx context.Context, id uuid.UUID, code string) error
	IncrementAttempts(ctx context.Context, id uuid.UUID) error
}

// KYCRepository defines the data access operations for KYC entities
type KYCRepository interface {
	CreateKYC(ctx context.Context, kyc domain.KYC) (*domain.KYC, error)
	GetKYCByUserID(ctx context.Context, userID uuid.UUID) (*domain.KYC, error)
	UpdateKYC(ctx context.Context, kyc domain.KYC) (*domain.KYC, error)
}

// EmailService defines operations for sending emails
type EmailRepository interface {
	SendVerificationEmail(ctx context.Context, email, name, code string) error
	SendPasswordResetEmail(ctx context.Context, email, name, code string) error
}


// WaitlistRepository defines the storage operations for the waitlist
type WaitlistRepository interface {
	CreateWaitlistEntry(ctx context.Context, entry domain.WaitlistEntry) (*domain.WaitlistEntry, error)
	GetWaitlistEntryByEmail(ctx context.Context, email string) (*domain.WaitlistEntry, error)
	GetWaitlistEntryByID(ctx context.Context, id uuid.UUID) (*domain.WaitlistEntry, error)
	GetWaitlistEntryByReferralCode(ctx context.Context, code string) (*domain.WaitlistEntry, error)
	ListWaitlistEntries(ctx context.Context, limit, offset int, filters map[string]string) ([]domain.WaitlistEntry, int64, error)
	ExportWaitlistToCsv(ctx context.Context) ([]byte, error)
}