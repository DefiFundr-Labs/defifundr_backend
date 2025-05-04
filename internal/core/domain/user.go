package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents the user entity in the domain model
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

// NewUser creates a new User instance with default values
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

// UserService defines the operations that can be performed on the User entity
type UserService interface {
	RegisterUser(ctx context.Context, user User) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user User) (*User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
}
