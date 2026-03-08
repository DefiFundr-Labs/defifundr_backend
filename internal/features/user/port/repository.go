package userport

import (
	"context"

	userdomain "github.com/demola234/defifundr/internal/features/user/domain"
	"github.com/google/uuid"
)

// UserRepository defines data access operations for users.
type UserRepository interface {
	CreateUser(ctx context.Context, user userdomain.User) (*userdomain.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*userdomain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*userdomain.User, error)
	GetUserCompanyInfo(ctx context.Context, id uuid.UUID) (*userdomain.CompanyInfo, error)
	UpdateUser(ctx context.Context, user userdomain.User) (*userdomain.User, error)
	UpdateUserPersonalDetails(ctx context.Context, user userdomain.User) (*userdomain.User, error)
	UpdateUserAddressDetails(ctx context.Context, user userdomain.User) (*userdomain.User, error)
	UpdateUserBusinessDetails(ctx context.Context, companyInfo userdomain.CompanyInfo) (*userdomain.CompanyInfo, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	DeactivateUser(ctx context.Context, id uuid.UUID) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SetMFASecret(ctx context.Context, userID uuid.UUID, secret string) error
	GetMFASecret(ctx context.Context, userID uuid.UUID) (string, error)
}

// UserService defines the business operations for user management.
type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*userdomain.User, error)
	UpdateUser(ctx context.Context, user userdomain.User) (*userdomain.User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	ResetUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	UpdateKYC(ctx context.Context, kyc userdomain.KYC) error
}
