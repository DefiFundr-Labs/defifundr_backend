package authdto

import (
	"errors"
	"regexp"
	"time"
)

// Web3AuthLoginRequest is the request for Web3Auth login.
type Web3AuthLoginRequest struct {
	WebAuthToken string `json:"web_auth_token" binding:"required"`
}

// RegisterUserRequest is the request for email/password registration.
type RegisterUserRequest struct {
	Email        string `json:"email" binding:"omitempty"`
	Password     string `json:"password,omitempty" binding:"omitempty,min=8"`
	FirstName    string `json:"first_name" binding:"omitempty"`
	LastName     string `json:"last_name" binding:"omitempty"`
	Provider     string `json:"provider" binding:"omitempty"`
	ProviderID   string `json:"provider_id" binding:"omitempty"`
	WebAuthToken string `json:"web_auth_token" binding:"required"`
}

func (r *RegisterUserRequest) Validate() error {
	if !isValidEmail(r.Email) {
		return errors.New("invalid email format")
	}
	return validatePassword(r.Password)
}

// LoginRequest is the request for email/password login.
type LoginRequest struct {
	Email        string `json:"email" binding:"omitempty"`
	Password     string `json:"password,omitempty" binding:"omitempty,min=8"`
	Provider     string `json:"provider" binding:"omitempty"`
	ProviderID   string `json:"provider_id" binding:"omitempty"`
	WebAuthToken string `json:"web_auth_token" binding:"required"`
}

func (r *LoginRequest) Validate() error {
	if r.Provider != "email" && r.Provider != "google" && r.Provider != "apple" {
		return errors.New("invalid provider")
	}
	if r.Provider == "email" && r.Password == "" {
		return errors.New("password is required")
	}
	if r.WebAuthToken == "" {
		return errors.New("web auth token is required")
	}
	return nil
}

// RefreshTokenRequest is the request for token refresh.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RegisterPersonalDetailsRequest is the request for personal details.
type RegisterPersonalDetailsRequest struct {
	FirstName           string `json:"first_name" binding:"required"`
	LastName            string `json:"last_name" binding:"required"`
	Nationality         string `json:"nationality" binding:"required"`
	PersonalAccountType string `json:"personal_account_type"`
	PhoneNumber         string `json:"phone_number"`
}

// RegisterAddressDetailsRequest is the request for address details.
type RegisterAddressDetailsRequest struct {
	UserAddress string `json:"user_address" binding:"required"`
	City        string `json:"city" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Country     string `json:"country" binding:"required"`
}

// RegisterBusinessDetailsRequest is the request for business details.
type RegisterBusinessDetailsRequest struct {
	CompanyName        string `json:"company_name" binding:"required"`
	CompanyAddress     string `json:"company_address" binding:"required"`
	CompanyDescription string `json:"company_description" binding:"required"`
	CompanyIndustry    string `json:"company_industry"`
	CompanyCountry     string `json:"company_country" binding:"required"`
	CompanySize        string `json:"company_size"`
	AccountType        string `json:"account_type"`
}

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
	return nil
}

// LinkWalletRequest is the request for linking a wallet.
type LinkWalletRequest struct {
	Address string `json:"address" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Chain   string `json:"chain" binding:"required"`
}

// RevokeDeviceRequest is the request to revoke a device/session.
type RevokeDeviceRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

// LogoutRequest is the request for logout.
type LogoutRequest struct {
	SessionID string `json:"session_id"`
}

// VerifyMFARequest is the request to verify an MFA code.
type VerifyMFARequest struct {
	Code string `json:"code" binding:"required"`
}

// ForgotPasswordRequest is the request to initiate password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyResetOTPRequest is the request to verify a password reset OTP.
type VerifyResetOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

// CompletePasswordResetRequest is the request to complete a password reset.
type CompletePasswordResetRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// UpdateKYCRequest is the request for KYC update.
type UpdateKYCRequest struct {
	IDType           string    `json:"id_type" binding:"required"`
	IDNumber         string    `json:"id_number" binding:"required"`
	IDIssuingCountry string    `json:"id_issuing_country" binding:"required"`
	IDExpiryDate     time.Time `json:"id_expiry_date" binding:"required"`
	IDFrontImage     string    `json:"id_front_image" binding:"required"`
	IDBackImage      string    `json:"id_back_image"`
	SelfieImage      string    `json:"selfie_image" binding:"required"`
}

// --- validation helpers ---

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasUpperCase || !hasLowerCase || !hasNumber || !hasSpecial {
		return errors.New("password must contain at least one uppercase letter, one lowercase letter, one number, and one special character")
	}
	return nil
}
