package services

import (
	"context"

	"github.com/defifundr/defifundr-backend/internal/core/domain"
	"github.com/defifundr/defifundr-backend/internal/core/ports"
	"github.com/google/uuid"
)

type payrollService struct {
	payrollRepo     ports.PayrollRepository
	userRepo        ports.UserRepository
	walletRepo      ports.WalletRepository
	blockchainSvc   ports.BlockchainService
	transactionRepo ports.TransactionRepository
	notificationSvc ports.NotificationService
}

func NewPayrollService(
	payrollRepo ports.PayrollRepository,
	userRepo ports.UserRepository,
	walletRepo ports.WalletRepository,
	blockchainSvc ports.BlockchainService,
	transactionRepo ports.TransactionRepository,
	notificationSvc ports.NotificationService,
) ports.PayrollService {
	return &payrollService{
		payrollRepo:     payrollRepo,
		userRepo:        userRepo,
		walletRepo:      walletRepo,
		blockchainSvc:   blockchainSvc,
		transactionRepo: transactionRepo,
		notificationSvc: notificationSvc,
	}
}

func (s *payrollService) CreatePayroll(ctx context.Context, params domain.CreatePayrollParams) (domain.Payroll, error) {
	// Verify employer exists
	_, err := s.userRepo.GetUserByID(ctx, params.EmployerID)
	if err != nil {
		return domain.Payroll{}, domain.ErrUserNotFound
	}

	// Create payroll
	payroll, err := s.payrollRepo.CreatePayroll(ctx, params)
	if err != nil {
		return domain.Payroll{}, err
	}

	// Add employees to the payroll if provided
	if len(params.EmployeeIDs) > 0 {
		for _, employeeID := range params.EmployeeIDs {
			// Verify employee exists
			_, err := s.userRepo.GetUserByID(ctx, employeeID)
			if err != nil {
				continue // Skip if employee doesn't exist
			}

			// Add employee to payroll
			_, err = s.payrollRepo.AddEmployeeToPayroll(ctx, domain.AddEmployeeToPayrollParams{
				PayrollID:  payroll.ID,
				EmployeeID: employeeID,
			})
			if err != nil {
				continue // Skip if error adding employee
			}

			// Create notification for employee
			s.notificationSvc.CreateNotification(ctx, domain.CreateNotificationParams{
				UserID:  employeeID,
				Message: "You have been added to a new payroll",
				Type:    domain.NotificationTypePayroll,
			})
		}
	}

	return payroll, nil
}

func (s *payrollService) GetPayrollByID(ctx context.Context, id uuid.UUID) (domain.Payroll, error) {
	payroll, err := s.payrollRepo.GetPayrollByID(ctx, id)
	if err != nil {
		return domain.Payroll{}, domain.ErrPayrollNotFound
	}

	// Get employees for this payroll
	employees, err := s.payrollRepo.GetEmployeesByPayrollID(ctx, id)
	if err == nil {
		payroll.Employees = employees
	}

	return payroll, nil
}

func (s *payrollService) ListPayrollsByEmployerID(ctx context.Context, employerID uuid.UUID) ([]domain.Payroll, error) {
	return s.payrollRepo.ListPayrollsByEmployerID(ctx, employerID)
}

func (s *payrollService) ListPayrollsByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]domain.Payroll, error) {
	return s.payrollRepo.ListPayrollsByOrganizationID(ctx, organizationID)
}

func (s *payrollService) UpdatePayrollStatus(ctx context.Context, params domain.UpdatePayrollStatusParams) (domain.Payroll, error) {
	// Verify payroll exists
	_, err := s.payrollRepo.GetPayrollByID(ctx, params.ID)
	if err != nil {
		return domain.Payroll{}, domain.ErrPayrollNotFound
	}

	return s.payrollRepo.UpdatePayrollStatus(ctx, params)
}

func (s *payrollService) DeployPayrollContract(ctx context.Context, payrollID uuid.UUID) (domain.Payroll, error) {
	// Get payroll
	payroll, err := s.payrollRepo.GetPayrollByID(ctx, payrollID)
	if err != nil {
		return domain.Payroll{}, domain.ErrPayrollNotFound
	}

	// Get employees for this payroll
	employees, err := s.payrollRepo.GetEmployeesByPayrollID(ctx, payrollID)
	if err != nil || len(employees) == 0 {
		return domain.Payroll{}, domain.ErrNoEmployeesInPayroll
	}

	// Get wallet addresses for each employee
	var employeeAddresses []string
	for _, employee := range employees {
		wallet, err := s.walletRepo.GetPrimaryWallet(ctx, employee.EmployeeID)
		if err != nil {
			continue
		}
		employeeAddresses = append(employeeAddresses, wallet.WalletAddress)
	}

	if len(employeeAddresses) == 0 {
		return domain.Payroll{}, domain.ErrInvalidWalletAddress
	}

	// Deploy smart contract
	contractAddress, err := s.blockchainSvc.DeployPayrollContract(
		ctx,
		payrollID,
		employeeAddresses,
		payroll.SalaryAmount.String(),
		string(payroll.Currency),
	)
	if err != nil {
		return domain.Payroll{}, domain.ErrContractDeployment
	}

	// Update payroll with contract address
	return s.payrollRepo.UpdatePayrollContractAddress(ctx, domain.UpdatePayrollContractParams{
		ID:              payrollID,
		ContractAddress: contractAddress,
	})
}

func (s *payrollService) ProcessPayroll(ctx context.Context, payrollID uuid.UUID) ([]domain.Transaction, error) {
	// Get payroll
	payroll, err := s.payrollRepo.GetPayrollByID(ctx, payrollID)
	if err != nil {
		return nil, domain.ErrPayrollNotFound
	}

	if payroll.ContractAddress == nil {
		return nil, domain.ErrContractDeployment
	}

	// Get employees for this payroll
	employees, err := s.payrollRepo.GetEmployeesByPayrollID(ctx, payrollID)
	if err != nil || len(employees) == 0 {
		return nil, domain.ErrNoEmployeesInPayroll
	}

	// Process payroll for each employee
	var transactions []domain.Transaction
	for _, employee := range employees {
		// Get employee's wallet
		wallet, err := s.walletRepo.GetPrimaryWallet(ctx, employee.EmployeeID)
		if err != nil {
			continue
		}

		// Execute payroll contract for this employee
		txHash, err := s.blockchainSvc.ExecutePayroll(ctx, *payroll.ContractAddress, wallet.WalletAddress)
		if err != nil {
			continue
		}

		// Create transaction record
		transaction, err := s.transactionRepo.CreateTransaction(ctx, domain.CreateTransactionParams{
			UserID:   employee.EmployeeID,
			TxHash:   txHash,
			Amount:   payroll.SalaryAmount,
			Currency: payroll.Currency,
			Type:     domain.TransactionTypePayroll,
			Status:   domain.TransactionStatusPending,
		})
		if err != nil {
			continue
		}

		transactions = append(transactions, transaction)

		// Create notification for employee
		s.notificationSvc.CreateNotification(ctx, domain.CreateNotificationParams{
			UserID:  employee.EmployeeID,
			Message: "You have received a payroll payment",
			Type:    domain.NotificationTypePayroll,
		})
	}

	// Update payroll status if all transactions were successful
	if len(transactions) > 0 {
		_, err = s.payrollRepo.UpdatePayrollStatus(ctx, domain.UpdatePayrollStatusParams{
			ID:     payrollID,
			Status: domain.PayrollStatusCompleted,
		})

		if err != nil {
			return nil, err
		}
	}

	return transactions, nil
}

func (s *payrollService) AddEmployeeToPayroll(ctx context.Context, params domain.AddEmployeeToPayrollParams) (domain.PayrollEmployee, error) {
	// Verify payroll exists
	_, err := s.payrollRepo.GetPayrollByID(ctx, params.PayrollID)
	if err != nil {
		return domain.PayrollEmployee{}, domain.ErrPayrollNotFound
	}

	// Verify employee exists
	_, err = s.userRepo.GetUserByID(ctx, params.EmployeeID)
	if err != nil {
		return domain.PayrollEmployee{}, domain.ErrUserNotFound
	}

	employee, err := s.payrollRepo.AddEmployeeToPayroll(ctx, params)
	if err != nil {
		return domain.PayrollEmployee{}, err
	}

	// Notify employee
	s.notificationSvc.CreateNotification(ctx, domain.CreateNotificationParams{
		UserID:  params.EmployeeID,
		Message: "You have been added to a payroll",
		Type:    domain.NotificationTypePayroll,
	})

	return employee, nil
}

func (s *payrollService) GetEmployeesByPayrollID(ctx context.Context, payrollID uuid.UUID) ([]domain.PayrollEmployee, error) {
	return s.payrollRepo.GetEmployeesByPayrollID(ctx, payrollID)
}

func (s *payrollService) GetPayrollsByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]domain.Payroll, error) {
	return s.payrollRepo.GetPayrollsByEmployeeID(ctx, employeeID)
}

func (s *payrollService) RemoveEmployeeFromPayroll(ctx context.Context, payrollID, employeeID uuid.UUID) error {
	// Verify payroll exists
	_, err := s.payrollRepo.GetPayrollByID(ctx, payrollID)
	if err != nil {
		return domain.ErrPayrollNotFound
	}

	// Verify employee exists
	_, err = s.userRepo.GetUserByID(ctx, employeeID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	err = s.payrollRepo.RemoveEmployeeFromPayroll(ctx, payrollID, employeeID)
	if err != nil {
		return err
	}

	// Notify employee
	s.notificationSvc.CreateNotification(ctx, domain.CreateNotificationParams{
		UserID:  employeeID,
		Message: "You have been removed from a payroll",
		Type:    domain.NotificationTypePayroll,
	})

	return nil
}
