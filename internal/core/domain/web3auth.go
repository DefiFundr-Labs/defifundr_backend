// domain/web3auth.go
package domain

import (
	"time"

	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Web3AuthClaims represents JWT claims from Web3Auth
type Web3AuthClaims struct {
	jwtv4.StandardClaims
	Email             string   `json:"email"`
	Name              string   `json:"name"`
	ProfileImage      string   `json:"profileImage"`
	Verifier          string   `json:"verifier"`
	VerifierID        string   `json:"verifierId"`
	AggregateVerifier string   `json:"aggregateVerifier"`
	Wallets           []Wallet `json:"wallets"`
	Nonce             string   `json:"nonce"`
}

// Wallet represents a blockchain wallet in Web3Auth tokens
type Wallet struct {
	PublicKey string `json:"public_key"`
	Type      string `json:"type"`
	Curve     string `json:"curve"`
}

// UserWallet represents a user's blockchain wallet
type UserWallet struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Address   string    `json:"address"`
	Type      string    `json:"type"`
	Chain     string    `json:"chain"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// AuthProvider enumerates supported authentication providers
type AuthProvider string

const (
	EmailProvider    AuthProvider = "email"
	GoogleProvider   AuthProvider = "google"
	AppleProvider    AuthProvider = "apple"
	FacebookProvider AuthProvider = "facebook"
	DiscordProvider  AuthProvider = "discord"
	TwitterProvider  AuthProvider = "twitter"
	Web3AuthProvider AuthProvider = "web3auth"
)

// SecurityEvent represents a user security event
type SecurityEvent struct {
	ID        uuid.UUID              `json:"id"`
	UserID    uuid.UUID              `json:"user_id"`
	EventType string                 `json:"event_type"`
	IPAddress string                 `json:"ip_address"`
	UserAgent string                 `json:"user_agent"`
	Metadata  map[string]interface{} `json:"metadata"`
	Timestamp time.Time              `json:"timestamp"`
}

// ProfileCompletion tracks user profile completion status
type ProfileCompletion struct {
	UserID               uuid.UUID `json:"user_id"`
	CompletionPercentage int       `json:"completion_percentage"`
	MissingFields        []string  `json:"missing_fields,omitempty"`
	RequiredActions      []string  `json:"required_actions,omitempty"`
}

// DeviceInfo stores information about a user's device
type DeviceInfo struct {
	SessionID       uuid.UUID `json:"session_id"`
	DeviceID        uuid.UUID `json:"device_id,omitempty"`
	Browser         string    `json:"browser"`
	OperatingSystem string    `json:"operating_system"`
	DeviceType      string    `json:"device_type"`
	DeviceModel     string    `json:"device_model,omitempty"`
	IPAddress       string    `json:"ip_address"`
	LoginType       string    `json:"login_type"`
	PushEnabled     bool      `json:"push_enabled"`
	LastUsed        time.Time `json:"last_used"`
	CreatedAt       time.Time `json:"created_at"`
}

// MFASetup represents a multi-factor authentication setup
type MFASetup struct {
	UserID        uuid.UUID `json:"user_id"`
	Secret        string    `json:"secret"`
	Verified      bool      `json:"verified"`
	RecoveryCodes []string  `json:"recovery_codes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
