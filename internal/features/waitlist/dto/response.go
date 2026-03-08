package waitlistdto

import (
	"time"

	"github.com/google/uuid"
)

// EntryResponse represents a single waitlist entry in API responses.
type EntryResponse struct {
	ID             uuid.UUID  `json:"id"`
	Email          string     `json:"email"`
	FullName       string     `json:"full_name,omitempty"`
	ReferralCode   string     `json:"referral_code"`
	ReferralSource string     `json:"referral_source,omitempty"`
	Status         string     `json:"status"`
	Position       int        `json:"position,omitempty"`
	SignupDate     time.Time  `json:"signup_date"`
	InvitedDate    *time.Time `json:"invited_date,omitempty"`
}

// StatsResponse represents waitlist statistics.
type StatsResponse struct {
	TotalSignups    int            `json:"total_signups"`
	WaitingCount    int            `json:"waiting_count"`
	InvitedCount    int            `json:"invited_count"`
	RegisteredCount int            `json:"registered_count"`
	ConversionRate  float64        `json:"conversion_rate"`
	Sources         map[string]int `json:"sources"`
}

// SuccessResponse is a generic success envelope.
type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse is a generic error envelope.
type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PageResponse is a generic paginated response envelope.
type PageResponse struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
	Items      interface{} `json:"items"`
}
