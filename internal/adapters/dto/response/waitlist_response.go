package response

import (
	"time"

	"github.com/google/uuid"
)

type WaitlistEntryResponse struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	FullName      string     `json:"full_name,omitempty"`
	ReferralCode  string     `json:"referral_code"`
	ReferralSource string    `json:"referral_source,omitempty"`
	Status        string     `json:"status"`
	Position      int        `json:"position,omitempty"`
	SignupDate    time.Time  `json:"signup_date"`
	InvitedDate   *time.Time `json:"invited_date,omitempty"`
}

type WaitlistStatsResponse struct {
	TotalSignups   int                 `json:"total_signups"`
	WaitingCount   int                 `json:"waiting_count"`
	InvitedCount   int                 `json:"invited_count"`
	RegisteredCount int                `json:"registered_count"`
	ConversionRate float64             `json:"conversion_rate"`
	Sources        map[string]int      `json:"sources"`
}