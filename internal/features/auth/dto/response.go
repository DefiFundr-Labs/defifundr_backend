package authdto

import (
	"time"

	"github.com/google/uuid"
)

// SuccessResponse is a generic success response.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is a generic error response.
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginUserResponse represents a user in a login/register response.
type LoginUserResponse struct {
	ID                  string    `json:"id"`
	Email               string    `json:"email"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	ProfilePicture      string    `json:"profile_picture,omitempty"`
	AccountType         string    `json:"account_type"`
	AuthProvider        string    `json:"auth_provider"`
	ProviderID          string    `json:"provider_id,omitempty"`
	Nationality         string    `json:"nationality,omitempty"`
	PersonalAccountType string    `json:"personal_account_type,omitempty"`
	MFAEnabled          bool      `json:"mfa_enabled"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// SessionResponse represents a session.
type SessionResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token,omitempty"`
	UserLoginType string    `json:"user_login_type"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// DeviceResponse represents a user device/session.
type DeviceResponse struct {
	SessionID       string    `json:"session_id"`
	Browser         string    `json:"browser"`
	OperatingSystem string    `json:"operating_system"`
	DeviceType      string    `json:"device_type"`
	IPAddress       string    `json:"ip_address"`
	LoginType       string    `json:"login_type"`
	LastUsed        time.Time `json:"last_used"`
	CreatedAt       time.Time `json:"created_at"`
}

// UserWalletResponse represents a user wallet.
type UserWalletResponse struct {
	ID        string `json:"id"`
	Address   string `json:"address"`
	Type      string `json:"type"`
	Chain     string `json:"chain"`
	IsDefault bool   `json:"is_default"`
}

// ProfileCompletionResponse represents profile completion status.
type ProfileCompletionResponse struct {
	CompletionPercentage int      `json:"completion_percentage"`
	MissingFields        []string `json:"missing_fields,omitempty"`
	RequiredActions      []string `json:"required_actions,omitempty"`
}

// SecurityEventResponse represents a security event.
type SecurityEventResponse struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"event_type"`
	IPAddress string                 `json:"ip_address"`
	UserAgent string                 `json:"user_agent"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// MFASetupResponse represents the MFA setup response.
type MFASetupResponse struct {
	TOTPURI           string   `json:"totp_uri"`
	SetupInstructions []string `json:"setup_instructions"`
}
