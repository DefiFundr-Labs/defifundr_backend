package response

import (
	"time"

	"github.com/google/uuid"
)

// SuccessResponse is a generic success response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is a generic error response
type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginUserResponse represents a user after login
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
	UserAddress         string    `json:"user_address,omitempty"`
	City                string    `json:"city,omitempty"`
	PostalCode          string    `json:"postal_code,omitempty"`
	Country             string    `json:"country,omitempty"`
	PhoneNumber         string    `json:"phone_number,omitempty"`
	CompanyName         string    `json:"company_name,omitempty"`
	CompanyAddress      string    `json:"company_address,omitempty"`
	CompanyCity         string    `json:"company_city,omitempty"`
	CompanyPostalCode   string    `json:"company_postal_code,omitempty"`
	CompanyCountry      string    `json:"company_country,omitempty"`
	CompanyWebsite      string    `json:"company_website,omitempty"`
	EmploymentType      string    `json:"employment_type,omitempty"`
	MFAEnabled          bool      `json:"mfa_enabled"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// SessionResponse represents a session
type SessionResponse struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	AccessToken   string    `json:"access_token"`
	RefreshToken  string    `json:"refresh_token,omitempty"`
	UserLoginType string    `json:"user_login_type"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}

// DeviceResponse represents a user's device/session
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

// UserWalletResponse represents a user's blockchain wallet
type UserWalletResponse struct {
	ID        string `json:"id"`
	Address   string `json:"address"`
	Type      string `json:"type"`
	Chain     string `json:"chain"`
	IsDefault bool   `json:"is_default"`
}

// ProfileCompletionResponse represents profile completion status
type ProfileCompletionResponse struct {
	CompletionPercentage int      `json:"completion_percentage"`
	MissingFields        []string `json:"missing_fields,omitempty"`
	RequiredActions      []string `json:"required_actions,omitempty"`
}

// SecurityEventResponse represents a security event
type SecurityEventResponse struct {
	ID        string                 `json:"id"`
	EventType string                 `json:"event_type"`
	IPAddress string                 `json:"ip_address"`
	UserAgent string                 `json:"user_agent"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// MFASetupResponse represents the response for MFA setup
type MFASetupResponse struct {
	TOTPURI           string   `json:"totp_uri"`
	SetupInstructions []string `json:"setup_instructions"`
}

// AuthTokenResponse represents the response with tokens
type AuthTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// OnboardingResponse represents onboarding data
type OnboardingResponse struct {
	IsNewUser       bool     `json:"is_new_user"`
	OnboardingSteps []string `json:"onboarding_steps"`
}

// UserSummaryResponse represents a summary of user data
type UserSummaryResponse struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture,omitempty"`
	AccountType    string `json:"account_type"`
}

// UserResponse represents the user data in responses
type UserResponse struct {
	ID         string     `json:"id"`
	Email      string     `json:"email"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Provider   string     `json:"provider"`
	ProviderID string     `json:"provider_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	User          LoginUserResponse `json:"user"`
	AccessToken   SessionResponse   `json:"access_token,omitempty"`
	RefreshToken  string            `json:"refresh_token,omitempty"`
	SessionID     uuid.UUID         `json:"session_id"`
	ExpiresAt     time.Time         `json:"expires_at"`
	TokenType     string            `json:"token_type,omitempty"`
	EmailVerified bool              `json:"email_verified,omitempty"`
}

// LoginResponse represents the login response
type RegistrationResponse struct {
	User          UserResponse    `json:"user"`
	AccessToken   SessionResponse `json:"access_token,omitempty"`
	RefreshToken  string          `json:"refresh_token,omitempty"`
	SessionID     uuid.UUID       `json:"session_id"`
	ExpiresAt     time.Time       `json:"expires_at"`
	TokenType     string          `json:"token_type,omitempty"`
	EmailVerified bool            `json:"email_verified,omitempty"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}

// OTPResponse represents an OTP verification response
type OTPResponse struct {
	UserID       uuid.UUID `json:"user_id"`
	Purpose      string    `json:"purpose"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsVerified   bool      `json:"is_verified"`
	AttemptsMade int       `json:"attempts_made"`
	MaxAttempts  int       `json:"max_attempts"`
}

// KYCResponse represents a KYC verification response
type KYCResponse struct {
	ID                 uuid.UUID `json:"id"`
	UserID             uuid.UUID `json:"user_id"`
	VerificationStatus string    `json:"verification_status"`
	SubmittedAt        time.Time `json:"submitted_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// PageResponse represents a paginated response
type PageResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	Items      interface{} `json:"items"`
}
