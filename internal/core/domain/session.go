package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID                 uuid.UUID  `json:"id"`
	UserID             uuid.UUID  `json:"user_id"`
	RefreshToken       string     `json:"refresh_token"`
	UserAgent          string     `json:"user_agent"`
	ClientIP           string     `json:"client_ip"`
	LastUsedAt         time.Time  `json:"last_used_at"`
	WebOAuthClientID   *string    `json:"web_oauth_client_id,omitempty"`
	OAuthAccessToken   string     `json:"oauth_access_token"`
	OAuthIDToken       *string    `json:"oauth_id_token,omitempty"`
	UserLoginType      string     `json:"user_login_type"`
	MFAVerified        bool       `json:"mfa_verified"`
	IsBlocked          bool       `json:"is_blocked"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
}
