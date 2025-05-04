package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	RefreshToken     string    `json:"-"`
	UserAgent        string    `json:"user_agent"`
	WebOAuthClientID *string   `json:"web_oauth_client_id,omitempty"`
	OAuthAccessToken string    `json:"oauth_access_token,omitempty"`
	OAuthIDToken     *string   `json:"oauth_id_token,omitempty"`
	UserLoginType    string    `json:"user_login_type"`
	MFAEnabled       bool      `json:"mfa_enabled"`
	ClientIP         string    `json:"client_ip"`
	IsBlocked        bool      `json:"is_blocked"`
	ExpiresAt        time.Time `json:"expires_at"`
	CreatedAt        time.Time `json:"created_at"`
	LastUsedAt       time.Time `json:"last_used_at"`
}
