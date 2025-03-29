package request

import (
	"errors"
	"regexp"
	"strings"
)

// WaitlistJoinRequest represents the request to join the waitlist
type WaitlistJoinRequest struct {
	Email          string `json:"email" binding:"required,email"`
	FullName       string `json:"full_name"`
	ReferralSource string `json:"referral_source"`
	ReferralCode   string `json:"referral_code"`
}

// Validate validates the waitlist join request
func (r *WaitlistJoinRequest) Validate() error {
	// Validate email
	if !isValidEmail(r.Email) {
		return errors.New("invalid email format")
	}

	// Trim whitespace
	r.Email = strings.TrimSpace(r.Email)
	r.FullName = strings.TrimSpace(r.FullName)
	r.ReferralSource = strings.TrimSpace(r.ReferralSource)
	r.ReferralCode = strings.TrimSpace(r.ReferralCode)

	return nil
}

// WaitlistInviteRequest represents the request to invite users from the waitlist
type WaitlistInviteRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

// Validate validates the waitlist invite request
func (r *WaitlistInviteRequest) Validate() error {
	if len(r.IDs) == 0 {
		return errors.New("at least one ID must be provided")
	}

	// Validate that all IDs are UUIDs
	for _, id := range r.IDs {
		if !isValidUUID(id) {
			return errors.New("invalid UUID format")
		}
	}

	return nil
}

// WaitlistListRequest represents the request to list waitlist entries
type WaitlistListRequest struct {
	Page     int               `json:"page" form:"page"`
	PageSize int               `json:"page_size" form:"page_size"`
	Filters  map[string]string `json:"filters" form:"filters"`
}

// Validate validates the waitlist list request
func (r *WaitlistListRequest) Validate() error {
	// Set default values if not provided
	if r.Page <= 0 {
		r.Page = 1
	}

	if r.PageSize <= 0 {
		r.PageSize = 10
	} else if r.PageSize > 100 {
		r.PageSize = 100 // Limit max page size
	}

	return nil
}

// Helper function for UUID validation
func isValidUUID(u string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}$")
	return r.MatchString(u)
}