package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity in the domain model
type User struct {
	ID                     uuid.UUID  `json:"id"`
	FirstName              string    `json:"first_name,omitempty"`
	LastName               string    `json:"last_name,omitempty"`
	PhoneNumber            string    `json:"phone_number,omitempty"`
	Email                  string     `json:"email"`
	PasswordHash           string    `json:"-"`
	ProfilePictureURL      *string    `json:"profile_picture_url,omitempty"`
	AuthProvider           string    `json:"auth_provider,omitempty"`
	ProviderID             *string    `json:"provider_id,omitempty"`
	EmailVerified          bool       `json:"email_verified"`
	EmailVerifiedAt        *time.Time `json:"email_verified_at,omitempty"`
	PhoneNumberVerified    bool       `json:"phone_number_verified"`
	PhoneNumberVerifiedAt  *time.Time `json:"phone_number_verified_at,omitempty"`
	AccountType            string     `json:"account_type"`
	AccountStatus          string     `json:"account_status"`
	TwoFactorEnabled       bool       `json:"two_factor_enabled"`
	TwoFactorMethod        *string    `json:"two_factor_method,omitempty"`
	UserLoginType          *string    `json:"user_login_type,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	UpdatedAt              time.Time  `json:"updated_at"`
	LastLoginAt            *time.Time `json:"last_login_at,omitempty"`
	DeletedAt              *time.Time `json:"deleted_at,omitempty"`
}

// PersonalUser represents additional personal user information
type PersonalUser struct {
	ID                      uuid.UUID  `json:"id"`
	UserID                  uuid.UUID  `json:"user_id"`
	Nationality             *string    `json:"nationality,omitempty"`
	ResidentialCountry      *string    `json:"residential_country,omitempty"`
	UserAddress             *string    `json:"user_address,omitempty"`
	UserCity                *string    `json:"user_city,omitempty"`
	UserPostalCode          *string    `json:"user_postal_code,omitempty"`
	Gender                  *string    `json:"gender,omitempty"`
	DateOfBirth             *time.Time `json:"date_of_birth,omitempty"`
	JobRole                 *string    `json:"job_role,omitempty"`
	PersonalAccountType     *string    `json:"personal_account_type,omitempty"`
	EmploymentType          *string    `json:"employment_type,omitempty"`
	TaxID                   *string    `json:"tax_id,omitempty"`
	DefaultPaymentCurrency  *string    `json:"default_payment_currency,omitempty"`
	DefaultPaymentMethod    *string    `json:"default_payment_method,omitempty"`
	HourlyRate              *float64   `json:"hourly_rate,omitempty"`
	Specialization          *string    `json:"specialization,omitempty"`
	KYCStatus               string     `json:"kyc_status"`
	KYCVerifiedAt           *time.Time `json:"kyc_verified_at,omitempty"`
	CreatedAt               time.Time  `json:"created_at"`
	UpdatedAt               time.Time  `json:"updated_at"`
}


// CompanyUser represents the relationship between users and companies
type CompanyUser struct {
	ID                        uuid.UUID  `json:"id"`
	CompanyID                 uuid.UUID  `json:"company_id"`
	UserID                    uuid.UUID  `json:"user_id"`
	Role                      string     `json:"role"`
	Department                *string    `json:"department,omitempty"`
	JobTitle                  *string    `json:"job_title,omitempty"`
	IsAdministrator           bool       `json:"is_administrator"`
	CanManagePayroll          bool       `json:"can_manage_payroll"`
	CanManageInvoices         bool       `json:"can_manage_invoices"`
	CanManageEmployees        bool       `json:"can_manage_employees"`
	CanManageCompanySettings  bool       `json:"can_manage_company_settings"`
	CanManageBankAccounts     bool       `json:"can_manage_bank_accounts"`
	CanManageWallets          bool       `json:"can_manage_wallets"`
	Permissions               *string    `json:"permissions,omitempty"` // JSONB as string
	IsActive                  bool       `json:"is_active"`
	AddedBy                   *uuid.UUID `json:"added_by,omitempty"`
	ReportsTo                 *uuid.UUID `json:"reports_to,omitempty"`
	HireDate                  *time.Time `json:"hire_date,omitempty"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

// CompanyEmployee represents employee information within a company
type CompanyEmployee struct {
	ID               uuid.UUID  `json:"id"`
	CompanyID        uuid.UUID  `json:"company_id"`
	UserID           *uuid.UUID `json:"user_id,omitempty"`
	EmployeeID       *string    `json:"employee_id,omitempty"`
	Department       *string    `json:"department,omitempty"`
	Position         *string    `json:"position,omitempty"`
	EmploymentStatus string     `json:"employment_status"`
	EmploymentType   *string    `json:"employment_type,omitempty"`
	StartDate        *time.Time `json:"start_date,omitempty"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	ManagerID        *uuid.UUID `json:"manager_id,omitempty"`
	SalaryAmount     *float64   `json:"salary_amount,omitempty"`
	SalaryCurrency   *string    `json:"salary_currency,omitempty"`
	SalaryFrequency  *string    `json:"salary_frequency,omitempty"`
	HourlyRate       *float64   `json:"hourly_rate,omitempty"`
	PaymentMethod    *string    `json:"payment_method,omitempty"`
	PaymentSplit     *string    `json:"payment_split,omitempty"`     // JSONB as string
	TaxInformation   *string    `json:"tax_information,omitempty"`   // JSONB as string
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// UserWithPersonalInfo represents a combined view of user and personal information
type UserWithPersonalInfo struct {
	User         User         `json:"user"`
	PersonalUser PersonalUser `json:"personal_user"`
}

// CompanyWithOwner represents a company with owner information
type CompanyWithOwner struct {
	Company        Company `json:"company"`
	OwnerFirstName *string `json:"owner_first_name,omitempty"`
	OwnerLastName  *string `json:"owner_last_name,omitempty"`
	OwnerEmail     string  `json:"owner_email"`
}

// Constructor functions

// NewUser creates a new User instance with default values
func NewUser(email, accountType string) *User {
	return &User{
		ID:               uuid.New(),
		Email:            email,
		AccountType:      accountType,
		AccountStatus:    "pending",
		EmailVerified:    false,
		TwoFactorEnabled: false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// NewPersonalUser creates a new PersonalUser instance
func NewPersonalUser(userID uuid.UUID) *PersonalUser {
	return &PersonalUser{
		ID:        uuid.New(),
		UserID:    userID,
		KYCStatus: "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewCompanyUser creates a new CompanyUser instance
func NewCompanyUser(companyID, userID uuid.UUID, role string) *CompanyUser {
	return &CompanyUser{
		ID:                        uuid.New(),
		CompanyID:                 companyID,
		UserID:                    userID,
		Role:                      role,
		IsAdministrator:           false,
		CanManagePayroll:          false,
		CanManageInvoices:         false,
		CanManageEmployees:        false,
		CanManageCompanySettings:  false,
		CanManageBankAccounts:     false,
		CanManageWallets:          false,
		IsActive:                  true,
		CreatedAt:                 time.Now(),
		UpdatedAt:                 time.Now(),
	}
}

// NewCompanyEmployee creates a new CompanyEmployee instance
func NewCompanyEmployee(companyID uuid.UUID) *CompanyEmployee {
	return &CompanyEmployee{
		ID:               uuid.New(),
		CompanyID:        companyID,
		EmploymentStatus: "active",
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// Service interfaces

// UserService defines the operations that can be performed on the User entity
type UserService interface {
	// Core user operations
	CreateUser(ctx context.Context, user *User) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	UpdateUserLoginTime(ctx context.Context, userID uuid.UUID) error
	SoftDeleteUser(ctx context.Context, userID uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
	SearchUsersByEmail(ctx context.Context, searchTerm string, limit, offset int) ([]*User, error)
	GetUsersByAccountType(ctx context.Context, accountType string, limit, offset int) ([]*User, error)
	CountUsers(ctx context.Context) (int64, error)
	CountUsersByAccountType(ctx context.Context, accountType string) (int64, error)
	
	// Authentication operations
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPasswordHash string) error
	CheckEmailExists(ctx context.Context, email string) (bool, error)
	VerifyEmail(ctx context.Context, userID uuid.UUID) error
	VerifyPhoneNumber(ctx context.Context, userID uuid.UUID) error
}

// PersonalUserService defines operations for personal user information
type PersonalUserService interface {
	CreatePersonalUser(ctx context.Context, personalUser *PersonalUser) (*PersonalUser, error)
	GetPersonalUserByID(ctx context.Context, id uuid.UUID) (*PersonalUser, error)
	GetPersonalUserByUserID(ctx context.Context, userID uuid.UUID) (*PersonalUser, error)
	GetPersonalUserWithUserDetails(ctx context.Context, id uuid.UUID) (*UserWithPersonalInfo, error)
	UpdatePersonalUser(ctx context.Context, personalUser *PersonalUser) (*PersonalUser, error)
	DeletePersonalUser(ctx context.Context, id uuid.UUID) error
	ListPersonalUsers(ctx context.Context, limit, offset int) ([]*PersonalUser, error)
	GetPersonalUsersByKYCStatus(ctx context.Context, kycStatus string, limit, offset int) ([]*PersonalUser, error)
	UpdateKYCStatus(ctx context.Context, id uuid.UUID, status string, verifiedAt *time.Time) (*PersonalUser, error)
}

// CompanyUserService defines operations for company-user relationships
type CompanyUserService interface {
	CreateCompanyUser(ctx context.Context, companyUser *CompanyUser) (*CompanyUser, error)
	GetCompanyUserByID(ctx context.Context, id uuid.UUID) (*CompanyUser, error)
	GetCompanyUserByCompanyAndUser(ctx context.Context, companyID, userID uuid.UUID) (*CompanyUser, error)
	UpdateCompanyUser(ctx context.Context, companyUser *CompanyUser) (*CompanyUser, error)
	DeactivateCompanyUser(ctx context.Context, id uuid.UUID) (*CompanyUser, error)
	DeleteCompanyUser(ctx context.Context, id uuid.UUID) error
	ListCompanyUsers(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*CompanyUser, error)
	GetCompanyAdministrators(ctx context.Context, companyID uuid.UUID) ([]*CompanyUser, error)
	GetUserCompanies(ctx context.Context, userID uuid.UUID) ([]*CompanyUser, error)
	UpdatePermissions(ctx context.Context, id uuid.UUID, permissions string) (*CompanyUser, error)
}

// CompanyEmployeeService defines operations for employee management
type CompanyEmployeeService interface {
	CreateCompanyEmployee(ctx context.Context, employee *CompanyEmployee) (*CompanyEmployee, error)
	GetCompanyEmployeeByID(ctx context.Context, id uuid.UUID) (*CompanyEmployee, error)
	GetCompanyEmployeeByEmployeeID(ctx context.Context, companyID uuid.UUID, employeeID string) (*CompanyEmployee, error)
	UpdateCompanyEmployee(ctx context.Context, employee *CompanyEmployee) (*CompanyEmployee, error)
	UpdateEmployeeStatus(ctx context.Context, id uuid.UUID, status string, endDate *time.Time) (*CompanyEmployee, error)
	DeleteCompanyEmployee(ctx context.Context, id uuid.UUID) error
	ListCompanyEmployees(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*CompanyEmployee, error)
	GetActiveEmployees(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*CompanyEmployee, error)
	GetEmployeesByDepartment(ctx context.Context, companyID uuid.UUID, department string, limit, offset int) ([]*CompanyEmployee, error)
	GetEmployeesByManager(ctx context.Context, managerID uuid.UUID) ([]*CompanyEmployee, error)
	CountCompanyEmployees(ctx context.Context, companyID uuid.UUID) (int64, error)
	CountActiveEmployees(ctx context.Context, companyID uuid.UUID) (int64, error)
}

// Repository interfaces (for dependency injection)

// UserRepository defines the data access interface for users
type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	UpdateLoginTime(ctx context.Context, userID uuid.UUID) error
	SoftDelete(ctx context.Context, userID uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*User, error)
	SearchByEmail(ctx context.Context, searchTerm string, limit, offset int) ([]*User, error)
	GetByAccountType(ctx context.Context, accountType string, limit, offset int) ([]*User, error)
	Count(ctx context.Context) (int64, error)
	CountByAccountType(ctx context.Context, accountType string) (int64, error)
}

// PersonalUserRepository defines the data access interface for personal users
type PersonalUserRepository interface {
	Create(ctx context.Context, personalUser *PersonalUser) (*PersonalUser, error)
	GetByID(ctx context.Context, id uuid.UUID) (*PersonalUser, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*PersonalUser, error)
	Update(ctx context.Context, personalUser *PersonalUser) (*PersonalUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*PersonalUser, error)
	GetByKYCStatus(ctx context.Context, kycStatus string, limit, offset int) ([]*PersonalUser, error)
}


// CompanyUserRepository defines the data access interface for company users
type CompanyUserRepository interface {
	Create(ctx context.Context, companyUser *CompanyUser) (*CompanyUser, error)
	GetByID(ctx context.Context, id uuid.UUID) (*CompanyUser, error)
	GetByCompanyAndUser(ctx context.Context, companyID, userID uuid.UUID) (*CompanyUser, error)
	Update(ctx context.Context, companyUser *CompanyUser) (*CompanyUser, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*CompanyUser, error)
	GetAdministrators(ctx context.Context, companyID uuid.UUID) ([]*CompanyUser, error)
	GetUserCompanies(ctx context.Context, userID uuid.UUID) ([]*CompanyUser, error)
}

// CompanyEmployeeRepository defines the data access interface for company employees
type CompanyEmployeeRepository interface {
	Create(ctx context.Context, employee *CompanyEmployee) (*CompanyEmployee, error)
	GetByID(ctx context.Context, id uuid.UUID) (*CompanyEmployee, error)
	GetByEmployeeID(ctx context.Context, companyID uuid.UUID, employeeID string) (*CompanyEmployee, error)
	Update(ctx context.Context, employee *CompanyEmployee) (*CompanyEmployee, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByCompany(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*CompanyEmployee, error)
	GetActive(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*CompanyEmployee, error)
	GetByDepartment(ctx context.Context, companyID uuid.UUID, department string, limit, offset int) ([]*CompanyEmployee, error)
	GetByManager(ctx context.Context, managerID uuid.UUID) ([]*CompanyEmployee, error)
	CountByCompany(ctx context.Context, companyID uuid.UUID) (int64, error)
	CountActive(ctx context.Context, companyID uuid.UUID) (int64, error)
}