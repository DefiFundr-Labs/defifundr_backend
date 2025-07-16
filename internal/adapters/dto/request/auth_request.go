package request

import "time"

// Web3AuthLoginRequest represents the login request for Web3Auth
type Web3AuthLoginRequest struct {
	WebAuthToken string `json:"web_auth_token" binding:"required"`
}

// RefreshTokenRequest represents the request to refresh an access token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RegisterPersonalDetailsRequest represents user personal details
type RegisterPersonalDetailsRequest struct {
	Nationality         string 		`json:"nationality" binding:"required"`
	PersonalAccountType string 		`json:"personal_account_type"`
	PhoneNumber         string 		`json:"phone_number"`
	Gender				string 		`json:"gender"`
	DateOfBirth			time.Time 	`json:"date_of_birth"`
}

// RegisterAddressDetailsRequest represents user address details
type RegisterAddressDetailsRequest struct {
	UserAddress string `json:"user_address" binding:"required"`
	City        string `json:"city" binding:"required"`
	PostalCode  string `json:"postal_code" binding:"required"`
	Country     string `json:"country" binding:"required"`
}

// RegisterBusinessDetailsRequest represents business details
type RegisterBusinessDetailsRequest struct {
	CompanyName        string `json:"company_name" binding:"required"`
	CompanyAddress     string `json:"company_address" binding:"required"`
	CompanyDescription string `json:"company_description" binding:"required"`
	CompanyIndustry    string `json:"company_industry"`
	CompanyCountry     string `json:"company_country" binding:"required"`
	CompanySize        string `json:"company_size"`
	AccountType        string `json:"account_type"`
}

// LinkWalletRequest represents the request to link a blockchain wallet
type LinkWalletRequest struct {
	Address string `json:"address" binding:"required"`
	Type    string `json:"type" binding:"required"`
	Chain   string `json:"chain" binding:"required"`
}

// RevokeDeviceRequest represents the request to revoke a device
type RevokeDeviceRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

// LogoutRequest represents the request to logout
type LogoutRequest struct {
	SessionID string `json:"session_id"`
}

// SetupMFARequest represents the request to setup MFA
type SetupMFARequest struct {
	// No fields needed, authentication is done through middleware
}

// VerifyMFARequest represents the request to verify an MFA code
type VerifyMFARequest struct {
	Code string `json:"code" binding:"required"`
}

// ConfirmResetPasswordRequest represents the request to confirm a password reset
type ConfirmResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePasswordRequest represents the request to change a password
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// RegisterUserRequest represents the user registration request
type RegisterUserRequest struct {
	Email        string `json:"email" binding:"omitempty"`
	Password     string `json:"password,omitempty" binding:"omitempty,min=8"`
	FirstName    string `json:"first_name" binding:"omitempty"`
	LastName     string `json:"last_name" binding:"omitempty"`
	AccountType  string `json:"account_type" binding:"omitempty"`
	Provider     string `json:"provider" binding:"omitempty"`
	ProviderID   string `json:"provider_id" binding:"omitempty"`
	WebAuthToken string `json:"web_auth_token" binding:"required"`
}

// UpdateUserPasswordRequest represents the request to update user password
type UpdateUserPasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
	OldPassword     string `json:"old_password" binding:"required"`
}

// CheckEmailRequest represents the request to check email availability
type CheckEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// LoginRequest represents the user login request
type LoginRequest struct {
	Email        string `json:"email" binding:"omitempty"`
	Password     string `json:"password,omitempty" binding:"omitempty,min=8"`
	Provider     string `json:"provider" binding:"omitempty"`
	ProviderID   string `json:"provider_id" binding:"omitempty"`
	WebAuthToken string `json:"web_auth_token" binding:"required"`
}

// UpdateProfileRequest represents the profile update request
type UpdateProfileRequest struct {
	FirstName          string `json:"first_name" binding:"required"`
	LastName           string `json:"last_name" binding:"required"`
	Nationality        string `json:"nationality" binding:"required"`
	Gender             string `json:"gender"`
	ResidentialCountry string `json:"residential_country"`
	JobRole            string `json:"job_role"`
	CompanyWebsite     string `json:"company_website"`
	EmploymentType     string `json:"employment_type"`
}

// VerifyEmailRequest represents the email verification request
type VerifyEmailRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	OTPCode string `json:"otp_code" binding:"required"`
}

// ResendOTPRequest represents the resend OTP request
type ResendOTPRequest struct {
	UserID        string `json:"user_id" binding:"required"`
	Purpose       string `json:"purpose" binding:"required"`
	ContactMethod string `json:"contact_method" binding:"required"`
}

// UpdateKYCRequest represents the KYC update request
type UpdateKYCRequest struct {
	IDType            string    `json:"id_type" binding:"required"`
	IDNumber          string    `json:"id_number" binding:"required"`
	IDIssuingCountry  string    `json:"id_issuing_country" binding:"required"`
	IDExpiryDate      time.Time `json:"id_expiry_date" binding:"required"`
	IDFrontImage      string    `json:"id_front_image" binding:"required"`
	IDBackImage       string    `json:"id_back_image"`
	SelfieImage       string    `json:"selfie_image" binding:"required"`
	AddressProofImage string    `json:"address_proof_image"`
}

// ForgotPasswordRequest represents the forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// VerifyResetOTPRequest represents the OTP verification request
type VerifyResetOTPRequest struct {
	Email string `json:"email" binding:"required,email"`
	OTP   string `json:"otp" binding:"required"`
}

// CompletePasswordResetRequest represents the final password reset request
type CompletePasswordResetRequest struct {
	Email       string `json:"email" binding:"required,email"`
	OTP         string `json:"otp" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}