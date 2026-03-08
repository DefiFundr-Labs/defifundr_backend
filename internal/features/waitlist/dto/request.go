package waitlistdto

import (
	"errors"
	"regexp"
	"strings"
)

// JoinRequest represents the request to join the waitlist.
type JoinRequest struct {
	Email          string `json:"email" binding:"required,email"`
	FullName       string `json:"full_name"`
	ReferralSource string `json:"referral_source"`
	ReferralCode   string `json:"referral_code"`
}

// Validate sanitises and validates the join request.
func (r *JoinRequest) Validate() error {
	if !isValidEmail(r.Email) {
		return errors.New("invalid email format")
	}

	r.Email = strings.TrimSpace(r.Email)
	r.FullName = strings.TrimSpace(r.FullName)
	r.ReferralSource = strings.TrimSpace(r.ReferralSource)
	r.ReferralCode = strings.TrimSpace(r.ReferralCode)

	return nil
}

// isValidEmail checks if the email format is valid.
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
