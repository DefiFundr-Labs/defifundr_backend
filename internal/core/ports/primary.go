package ports

import (
	"context"

	"github.com/defifundr/defifundr-backend/internal/core/domain"
	"github.com/google/uuid"
)

// Authentication Service Interface
type AuthService interface {
	RegisterUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error)
	LoginUser(ctx context.Context, params domain.LoginUserParams) (string, error) // Returns JWT token
	VerifyToken(ctx context.Context, token string) (domain.User, error)
}

// User Service Interface
type UserService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, params domain.UpdateUserParams) (domain.User, error)
	ListUsers(ctx context.Context, limit, offset int32) ([]domain.User, error)
}

// Wallet Service Interface
type WalletService interface {
	CreateWallet(ctx context.Context, params domain.CreateWalletParams) (domain.Wallet, error)
	GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Wallet, error)
	GetPrimaryWallet(ctx context.Context, userID uuid.UUID) (domain.Wallet, error)
	UpdateWalletPrimaryStatus(ctx context.Context, params domain.UpdateWalletPrimaryParams) (domain.Wallet, error)
}

// Organization Service Interface
type OrganizationService interface {
	CreateOrganization(ctx context.Context, params domain.CreateOrganizationParams) (domain.Organization, error)
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (domain.Organization, error)
	ListOrganizationsByEmployerID(ctx context.Context, employerID uuid.UUID) ([]domain.Organization, error)
	UpdateOrganization(ctx context.Context, id uuid.UUID, params domain.UpdateOrganizationParams) (domain.Organization, error)
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
	AddMemberToOrganization(ctx context.Context, params domain.AddMemberParams) (domain.Member, error)
	GetOrganizationMembers(ctx context.Context, organizationID uuid.UUID) ([]domain.Member, error)
	RemoveOrganizationMember(ctx context.Context, organizationID, employeeID uuid.UUID) error
	UpdateOrganizationMemberRole(ctx context.Context, params domain.UpdateMemberRoleParams) (domain.Member, error)
}

// Payroll Service Interface
type PayrollService interface {
	CreatePayroll(ctx context.Context, params domain.CreatePayrollParams) (domain.Payroll, error)
	GetPayrollByID(ctx context.Context, id uuid.UUID) (domain.Payroll, error)
	ListPayrollsByEmployerID(ctx context.Context, employerID uuid.UUID) ([]domain.Payroll, error)
	ListPayrollsByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]domain.Payroll, error)
	UpdatePayrollStatus(ctx context.Context, params domain.UpdatePayrollStatusParams) (domain.Payroll, error)
	DeployPayrollContract(ctx context.Context, payrollID uuid.UUID) (domain.Payroll, error)
	ProcessPayroll(ctx context.Context, payrollID uuid.UUID) ([]domain.Transaction, error)
	AddEmployeeToPayroll(ctx context.Context, params domain.AddEmployeeToPayrollParams) (domain.PayrollEmployee, error)
	GetEmployeesByPayrollID(ctx context.Context, payrollID uuid.UUID) ([]domain.PayrollEmployee, error)
	GetPayrollsByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]domain.Payroll, error)
	RemoveEmployeeFromPayroll(ctx context.Context, payrollID, employeeID uuid.UUID) error
}

// Invoice Service Interface
type InvoiceService interface {
	CreateInvoice(ctx context.Context, params domain.CreateInvoiceParams) (domain.Invoice, error)
	GetInvoiceByID(ctx context.Context, id uuid.UUID) (domain.Invoice, error)
	ListInvoicesByFreelancerID(ctx context.Context, freelancerID uuid.UUID) ([]domain.Invoice, error)
	ListInvoicesByEmployerID(ctx context.Context, employerID uuid.UUID) ([]domain.Invoice, error)
	UpdateInvoiceStatus(ctx context.Context, params domain.UpdateInvoiceStatusParams) (domain.Invoice, error)
	DeployInvoiceContract(ctx context.Context, invoiceID uuid.UUID) (domain.Invoice, error)
	PayInvoice(ctx context.Context, invoiceID uuid.UUID) (domain.Transaction, error)
}

// Transaction Service Interface
type TransactionService interface {
	GetTransactionByID(ctx context.Context, id uuid.UUID) (domain.Transaction, error)
	ListTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Transaction, error)
	VerifyTransaction(ctx context.Context, txHash string, blockchain domain.BlockchainType) (bool, error)
}

// Notification Service Interface
type NotificationService interface {
	CreateNotification(ctx context.Context, params domain.CreateNotificationParams) (domain.Notification, error)
	GetNotificationsByUserID(ctx context.Context, params domain.NotificationListParams) ([]domain.Notification, error)
	GetUnreadNotificationCount(ctx context.Context, userID uuid.UUID) (int64, error)
	MarkNotificationAsRead(ctx context.Context, id uuid.UUID) (domain.Notification, error)
	MarkAllNotificationsAsRead(ctx context.Context, userID uuid.UUID) error
}
