package request

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// Validate validates the email check request
func (r *CheckEmailRequest) Validate() error {
	if !isValidEmail(r.Email) {
		return errors.New("invalid email format")
	}
	return nil
}

// Validate validates the business details request
func (r *RegisterBusinessDetailsRequest) Validate() error {
	if r.CompanyName == "" {
		return errors.New("company name is required")
	}

	if r.CompanyIndustry == "" {
		return errors.New("company industry is required")
	}

	if r.AccountType == "" {
		return errors.New("account type is required")
	}

	if r.CompanyDescription == "" || len(r.CompanyDescription) < 50 {
		return errors.New("company description is required and should be more than 50 characters")
	}

	if r.CompanyCountry == "" {
		return errors.New("company country is required")
	}

	return nil
}

// Validate validates the login request
func (r *LoginRequest) Validate() error {
	if r.Provider != "email" && r.Provider != "google" && r.Provider != "apple" {
		return errors.New("invalid provider")
	}

	if r.Provider == "email" {
		if r.Email == "" {
			return errors.New("email is required for email provider")
		}
		if !isValidEmail(r.Email) {
			return errors.New("invalid email format")
		}
		if r.Password == "" {
			return errors.New("password is required for email provider")
		}
	}

	if r.Provider != "email" && r.ProviderID == "" {
		return errors.New("provider ID is required")
	}

	if (r.Provider == "apple" || r.Provider == "google") && r.WebAuthToken == "" {
		return errors.New("web auth token is required")
	}

	return nil
}

// Validate validates the profile update request
func (r *UpdateProfileRequest) Validate() error {
	if strings.TrimSpace(r.FirstName) == "" || strings.TrimSpace(r.LastName) == "" {
		return errors.New("first name and last name cannot be empty")
	}

	if strings.TrimSpace(r.Nationality) == "" {
		return errors.New("nationality cannot be empty")
	}

	return nil
}

// Validate validates the email verification request
func (r *VerifyEmailRequest) Validate() error {
	if strings.TrimSpace(r.UserID) == "" {
		return errors.New("user ID cannot be empty")
	}

	if strings.TrimSpace(r.OTPCode) == "" {
		return errors.New("OTP code cannot be empty")
	}

	return nil
}

// Validate validates the resend OTP request
func (r *ResendOTPRequest) Validate() error {
	if strings.TrimSpace(r.UserID) == "" {
		return errors.New("user ID cannot be empty")
	}

	if !isValidOTPPurpose(r.Purpose) {
		return errors.New("invalid OTP purpose")
	}

	if !isValidContactMethod(r.ContactMethod) {
		return errors.New("invalid contact method")
	}

	return nil
}

// Validate validates the KYC update request
func (r *UpdateKYCRequest) Validate() error {
	if !isValidIDType(r.IDType) {
		return errors.New("invalid ID type")
	}

	if strings.TrimSpace(r.IDNumber) == "" {
		return errors.New("ID number cannot be empty")
	}

	if strings.TrimSpace(r.IDIssuingCountry) == "" {
		return errors.New("ID issuing country cannot be empty")
	}

	if r.IDExpiryDate.Before(time.Now()) {
		return errors.New("ID expiry date must be in the future")
	}

	if strings.TrimSpace(r.IDFrontImage) == "" {
		return errors.New("ID front image is required")
	}

	if strings.TrimSpace(r.SelfieImage) == "" {
		return errors.New("selfie image is required")
	}

	return nil
}

// Helper functions for validation

// isValidEmail checks if the email format is valid
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// validatePassword checks if the password meets security requirements
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecialChar := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	if !hasUpperCase || !hasLowerCase || !hasNumber || !hasSpecialChar {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}

	return nil
}

// isValidAccountType checks if the account type is valid
func isValidAccountType(accountType string) bool {
	validTypes := []string{"personal", "business"}
	for _, validType := range validTypes {
		if accountType == validType {
			return true
		}
	}
	return false
}

// isValidOTPPurpose checks if the OTP purpose is valid
func isValidOTPPurpose(purpose string) bool {
	validPurposes := []string{"email_verification", "password_reset", "two_factor_auth"}
	for _, validPurpose := range validPurposes {
		if purpose == validPurpose {
			return true
		}
	}
	return false
}

// isValidContactMethod checks if the contact method is valid
func isValidContactMethod(method string) bool {
	validMethods := []string{"email", "phone"}
	for _, validMethod := range validMethods {
		if method == validMethod {
			return true
		}
	}
	return false
}

// isValidIDType checks if the ID type is valid
func isValidIDType(idType string) bool {
	validTypes := []string{"passport", "national_id", "drivers_license", "residence_permit"}
	for _, validType := range validTypes {
		if idType == validType {
			return true
		}
	}
	return false
}