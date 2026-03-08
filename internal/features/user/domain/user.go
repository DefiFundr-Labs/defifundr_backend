// Package userdomain contains domain types for the user feature.
package userdomain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents the user entity in the domain model.
type User struct {
	ID                  uuid.UUID `json:"id"`
	Email               string    `json:"email"`
	Password            *string   `json:"-"`
	PasswordHash        string    `json:"-"`
	ProfilePicture      *string   `json:"profile_picture,omitempty"`
	AccountType         string    `json:"account_type"`
	Gender              *string   `json:"gender,omitempty"`
	PersonalAccountType string    `json:"personal_account_type"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	Nationality         string    `json:"nationality"`
	ResidentialCountry  *string   `json:"residential_country,omitempty"`
	JobRole             *string   `json:"job_role,omitempty"`
	EmploymentType      *string   `json:"employment_type,omitempty"`
	UserCity            *string   `json:"user_city,omitempty"`
	UserAddress         *string   `json:"user_address,omitempty"`
	UserPostalCode      *string   `json:"user_postal_code,omitempty"`
	PhoneNumber         *string   `json:"phone_number,omitempty"`
	PhoneNumberVerified *bool     `json:"phone_number_verified,omitempty"`
	Address             string    `json:"address"`
	City                string    `json:"city"`
	PostalCode          string    `json:"postal_code"`
	AuthProvider        string    `json:"auth_provider"`
	ProviderID          string    `json:"provider_id"`
	EmployeeType        string    `json:"employee_type"`
	WebAuthToken        string    `json:"webauth_token"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// NewUser creates a new User instance with default values.
func NewUser(
	email, firstName, lastName, nationality string,
	accountType, personalAccountType string,
) *User {
	return &User{
		ID:                  uuid.New(),
		Email:               email,
		FirstName:           firstName,
		LastName:            lastName,
		Nationality:         nationality,
		AccountType:         accountType,
		PersonalAccountType: personalAccountType,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
}

// UserService defines the operations that can be performed on the User entity.
type UserService interface {
	RegisterUser(ctx context.Context, user User) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user User) (*User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}

// CompanyInfo holds company-related information for a business user.
type CompanyInfo struct {
	UserID              uuid.UUID `json:"user_id"`
	CompanyName         *string   `json:"company_name,omitempty"`
	CompanySize         *string   `json:"company_size,omitempty"`
	CompanyIndustry     *string   `json:"company_industry,omitempty"`
	CompanyDescription  *string   `json:"company_description,omitempty"`
	CompanyHeadquarters *string   `json:"company_headquarters,omitempty"`
	CompanyWebsite      *string   `json:"company_website,omitempty"`
	AccountType         string    `json:"account_type,omitempty"`
}

// KYC holds KYC verification information for a user.
type KYC struct {
	ID                   uuid.UUID `json:"id"`
	UserID               uuid.UUID `json:"user_id"`
	FaceVerification     bool      `json:"face_verification"`
	IdentityVerification bool      `json:"identity_verification"`
	VerificationType     string    `json:"verification_type"`
	VerificationNumber   string    `json:"verification_number"`
	VerificationStatus   string    `json:"verification_status"`
	UpdatedAt            time.Time `json:"updated_at"`
	CreatedAt            time.Time `json:"created_at"`
}
