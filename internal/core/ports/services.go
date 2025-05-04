package ports

import (
	"context"

	"github.com/demola234/defifundr/internal/core/domain"
	emailEnums "github.com/demola234/defifundr/pkg/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	// Web3Auth authentication
	AuthenticateWithWeb3(ctx context.Context, webAuthToken string, userAgent string, clientIP string) (*domain.User, *domain.Session, error)

	// User authentication
	Login(ctx context.Context, email string, user domain.User, password string) (*domain.User, error)
	RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error)

	// User profile completion
	RegisterPersonalDetails(ctx context.Context, user domain.User) (*domain.User, error)
	RegisterAddressDetails(ctx context.Context, user domain.User) (*domain.User, error)
	RegisterBusinessDetails(ctx context.Context, user domain.User) (*domain.User, error)
	GetProfileCompletionStatus(ctx context.Context, userID uuid.UUID) (*domain.ProfileCompletion, error)

	// Multi-factor authentication
	SetupMFA(ctx context.Context, userID uuid.UUID) (string, error)
	VerifyMFA(ctx context.Context, userID uuid.UUID, code string) (bool, error)

	// Wallet management
	LinkWallet(ctx context.Context, userID uuid.UUID, walletAddress string, walletType string, chain string) error
	GetUserWallets(ctx context.Context, userID uuid.UUID) ([]domain.UserWallet, error)

	// Session management
	CreateSession(ctx context.Context, userID uuid.UUID, userAgent, clientIP string, webOAuthClientID string, email string, loginType string) (*domain.Session, error)
	GetActiveDevices(ctx context.Context, userID uuid.UUID) ([]domain.DeviceInfo, error)
	RevokeSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	Logout(ctx context.Context, sessionID uuid.UUID) error
	RefreshToken(ctx context.Context, refreshToken, userAgent, clientIP string) (*domain.Session, string, error)

	// User operations
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)

	// Security
	LogSecurityEvent(ctx context.Context, eventType string, userID uuid.UUID, metadata map[string]interface{}) error
}

// UserService defines the use cases for user operations
type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	UpdateKYC(ctx context.Context, kyc domain.KYC) error
}

// WaitlistService defines the use cases for the waitlist feature
type WaitlistService interface {
	JoinWaitlist(ctx context.Context, email, fullName, referralSource string) (*domain.WaitlistEntry, error)
	GetWaitlistPosition(ctx context.Context, id uuid.UUID) (int, error)
	GetWaitlistStats(ctx context.Context) (map[string]interface{}, error)
	ListWaitlist(ctx context.Context, page, pageSize int, filters map[string]string) ([]domain.WaitlistEntry, int64, error)
	ExportWaitlist(ctx context.Context) ([]byte, error)
}

// EmailService defines methods for sending application emails
type EmailSender interface {
	SendEmail(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}) error
	SendEmailWithAttachment(ctx context.Context, recipient string, subject string, templateName string,
		data map[string]interface{}, attachments []emailEnums.EmailAttachment) error
	QueueEmail(ctx context.Context, recipient string, subject string, templateName string,
		data map[string]interface{}, priority emailEnums.EmailPriority) (string, error)
}

// EmailService defines methods for sending application emails
type EmailService interface {
	SendWaitlistConfirmation(ctx context.Context, email, name, referralCode string, position int) error
	SendWaitlistInvitation(ctx context.Context, email, name string, inviteLink string) error
	SendBatchUpdate(ctx context.Context, emails []string, subject, message string) error
}
