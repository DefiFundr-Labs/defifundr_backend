package ports

import (
	"context"
	"time"

	"github.com/demola234/defifundr/internal/core/domain"
	emailEnums "github.com/demola234/defifundr/pkg/utils"
	"github.com/google/uuid"
)

type AuthService interface {
	// Web3Auth authentication
	AuthenticateWithWeb3(ctx context.Context, webAuthToken string, userAgent string, clientIP string) (*domain.User, *domain.Session, error)

	// User authentication
	Login(ctx context.Context, email string, provider string, providerId string, webAuthToken string, password string) (*domain.User, error)
	RegisterUser(ctx context.Context, email string, firstName string, lastName string, authProvider string, webAuthToken string, password string) (*domain.User, error)

	// User profile completion
	RegisterPersonalDetails(ctx context.Context, userId uuid.UUID, nationality string, dateOfBirth time.Time, gender string, personalAccountType string, phoneNumber string) (*domain.User, error)
	RegisterAddressDetails(ctx context.Context, userId uuid.UUID, userAddress, userCity, userPostalCode, residentialCountry string) (*domain.User, error) 
	RegisterBusinessDetails(ctx context.Context, companyInfo domain.Company) (*domain.Company, error)
	GetProfileCompletionStatus(ctx context.Context, userID uuid.UUID) (*domain.ProfileCompletion, error)

	// Multi-factor authentication
	SetupMFA(ctx context.Context, userID uuid.UUID) (string, error)
	VerifyMFA(ctx context.Context, userID uuid.UUID, code string) (bool, error)

	// Wallet management
	LinkWallet(ctx context.Context, userID uuid.UUID, walletAddress string, walletType string, chain string) error
	GetUserWallets(ctx context.Context, userID uuid.UUID) ([]domain.UserWallet, error)

	// Session management
	CreateSession(ctx context.Context, userID uuid.UUID, userAgent, clientIP string, webOAuthClientID string, email string, loginType string) (*domain.Session, error)
	GetActiveDevices(ctx context.Context, userID uuid.UUID) ([]domain.DeviceInfo, error)
	RevokeSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error
	Logout(ctx context.Context, sessionID uuid.UUID) error
	RefreshToken(ctx context.Context, refreshToken, userAgent, clientIP string) (*domain.Session, string, error)

	// User operations
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	CheckEmailExists(ctx context.Context, email string) (bool, error)

	// Security
	LogSecurityEvent(ctx context.Context, eventType string, userID uuid.UUID, metadata map[string]interface{}) error

	// Forgot password
	InitiatePasswordReset(ctx context.Context, email string) error
	VerifyResetOTP(ctx context.Context, email string, otp string) error  // Just verify, don't invalidate
	ResetPassword(ctx context.Context, email string, otp string, newPassword string) error  // Verify OTP and reset password
}

// UserService defines the use cases for user operations
type UserService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error
	UpdateKYC(ctx context.Context, kyc domain.KYC) error
	ResetUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	GetPersonalUserByUserID(ctx context.Context, userID uuid.UUID) (*domain.PersonalUser, error)
}

// WaitlistService defines the use cases for the waitlist feature
type WaitlistService interface {
	JoinWaitlist(ctx context.Context, email, fullName, referralSource string) (*domain.WaitlistEntry, error)
	GetWaitlistPosition(ctx context.Context, id uuid.UUID) (int, error)
	GetWaitlistStats(ctx context.Context) (map[string]interface{}, error)
	ListWaitlist(ctx context.Context, page, pageSize int, filters map[string]string) ([]domain.WaitlistEntry, int64, error)
	ExportWaitlist(ctx context.Context) ([]byte, error)
}

// EmailService defines methods for sending application emails
type EmailSender interface {
	SendEmail(ctx context.Context, recipient string, subject string, templateName string, data map[string]interface{}) error
	SendEmailWithAttachment(ctx context.Context, recipient string, subject string, templateName string,
		data map[string]interface{}, attachments []emailEnums.EmailAttachment) error
	QueueEmail(ctx context.Context, recipient string, subject string, templateName string,
		data map[string]interface{}, priority emailEnums.EmailPriority) (string, error)
}

// EmailService defines methods for sending application emails
type EmailService interface {
	SendWaitlistConfirmation(ctx context.Context, email, name, referralCode string, position int) error
	SendPasswordResetEmail(ctx context.Context, email, name, otpCode string) error
	SendWaitlistInvitation(ctx context.Context, email, name string, inviteLink string) error
	SendBatchUpdate(ctx context.Context, emails []string, subject, message string) error
}

type CompanyService interface {
	CreateCompany(ctx context.Context, company domain.Company) (*domain.Company, error)
	GetCompanyByID(ctx context.Context, companyID uuid.UUID) (*domain.Company, error)
	GetCompanyByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Company, error)
	UpdateCompany(ctx context.Context, company domain.Company) (*domain.Company, error)
	UpdateCompanyKYB(ctx context.Context, companyID uuid.UUID, kybStatus string, verifiedAt *time.Time, method, provider, rejectionReason *string) (*domain.Company, error)
	// DeleteCompany(ctx context.Context, companyID uuid.UUID) error
	ListCompanies(ctx context.Context, limit, offset int) ([]*domain.CompanyWithOwner, error)
}

type PersonalUserService interface {
	CreatePersonalUser(ctx context.Context, personalUser domain.PersonalUser) (*domain.PersonalUser, error)
	GetPersonalUserByID(ctx context.Context, personalUserID uuid.UUID) (*domain.PersonalUser, error)
	GetPersonalUserByUserID(ctx context.Context, userID uuid.UUID) (*domain.PersonalUser, error)
	GetPersonalUserWithUserDetails(ctx context.Context, personalUserID uuid.UUID) (*domain.UserWithPersonalInfo, error)
	UpdatePersonalUser(ctx context.Context, personalUser domain.PersonalUser) (*domain.PersonalUser, error)
	UpdateKYCStatus(ctx context.Context, personalUserID uuid.UUID, status string, verifiedAt *time.Time) (*domain.PersonalUser, error)
	DeletePersonalUser(ctx context.Context, personalUserID uuid.UUID) error
	ListPersonalUsers(ctx context.Context, limit, offset int) ([]*domain.PersonalUser, error)
	GetPersonalUsersByKYCStatus(ctx context.Context, kycStatus string, limit, offset int) ([]*domain.PersonalUser, error)
}


type CompanyUserService interface {
	CreateCompanyUser(ctx context.Context, companyUser domain.CompanyUser) (*domain.CompanyUser, error)
	GetCompanyUserByID(ctx context.Context, companyUserID uuid.UUID) (*domain.CompanyUser, error)
	GetCompanyUserByCompanyAndUser(ctx context.Context, companyID, userID uuid.UUID) (*domain.CompanyUser, error)
	UpdateCompanyUser(ctx context.Context, companyUser domain.CompanyUser) (*domain.CompanyUser, error)
	DeactivateCompanyUser(ctx context.Context, companyUserID uuid.UUID) (*domain.CompanyUser, error)
	DeleteCompanyUser(ctx context.Context, companyUserID uuid.UUID) error
	ListCompanyUsers(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domain.CompanyUser, error)
	GetCompanyAdministrators(ctx context.Context, companyID uuid.UUID) ([]*domain.CompanyUser, error)
	GetUserCompanies(ctx context.Context, userID uuid.UUID) ([]*domain.CompanyUser, error)
	UpdatePermissions(ctx context.Context, companyUserID uuid.UUID, permissions string) (*domain.CompanyUser, error)
}


type CompanyEmployeeService interface {
	CreateCompanyEmployee(ctx context.Context, employee domain.CompanyEmployee) (*domain.CompanyEmployee, error)
	GetCompanyEmployeeByID(ctx context.Context, employeeID uuid.UUID) (*domain.CompanyEmployee, error)
	GetCompanyEmployeeByEmployeeID(ctx context.Context, companyID uuid.UUID, employeeID string) (*domain.CompanyEmployee, error)
	UpdateCompanyEmployee(ctx context.Context, employee domain.CompanyEmployee) (*domain.CompanyEmployee, error)
	UpdateEmployeeStatus(ctx context.Context, employeeID uuid.UUID, status string, endDate *time.Time) (*domain.CompanyEmployee, error)
	DeleteCompanyEmployee(ctx context.Context, employeeID uuid.UUID) error
	ListCompanyEmployees(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domain.CompanyEmployee, error)
	GetActiveEmployees(ctx context.Context, companyID uuid.UUID, limit, offset int) ([]*domain.CompanyEmployee, error)
	GetEmployeesByDepartment(ctx context.Context, companyID uuid.UUID, department string, limit, offset int) ([]*domain.CompanyEmployee, error)
	GetEmployeesByManager(ctx context.Context, managerID uuid.UUID) ([]*domain.CompanyEmployee, error)
	CountCompanyEmployees(ctx context.Context, companyID uuid.UUID) (int64, error)
	CountActiveEmployees(ctx context.Context, companyID uuid.UUID) (int64, error)
}