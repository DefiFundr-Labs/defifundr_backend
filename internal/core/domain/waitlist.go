package domain

import (
	"time"

	"github.com/google/uuid"
)


type WaitlistEntry struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	FullName      string     `json:"full_name,omitempty"`
	ReferralCode  string     `json:"referral_code"`
	ReferralSource string    `json:"referral_source,omitempty"`
	Status        string     `json:"status"`
	SignupDate    time.Time  `json:"signup_date"`
	InvitedDate   *time.Time `json:"invited_date,omitempty"`
	RegisteredDate *time.Time `json:"registered_date,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}