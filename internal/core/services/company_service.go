package services

import (
	"context"
	"fmt"
	"time"

	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/google/uuid"
)

type companyService struct {
	companyRepo ports.CompanyRepository
}

// NewCompanyService creates a new instance of companyService
func NewCompanyService(companyRepo ports.CompanyRepository) ports.CompanyService {
	return &companyService{
		companyRepo: companyRepo,
	}
}

// CreateCompany implements ports.CompanyService
func (c *companyService) CreateCompany(ctx context.Context, company domain.Company) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-service").Start(ctx, "CreateCompany")
	defer span.End()

	// Set default values
	if company.ID == uuid.Nil {
		company.ID = uuid.New()
	}
	if company.AccountStatus == "" {
		company.AccountStatus = "pending"
	}
	if company.KYBStatus == "" {
		company.KYBStatus = "pending"
	}
	company.CreatedAt = time.Now()
	company.UpdatedAt = time.Now()

	createdCompany, err := c.companyRepo.CreateCompany(ctx, company)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to create company: %w", err)
	}

	return createdCompany, nil
}

// GetCompanyByID implements ports.CompanyService
func (c *companyService) GetCompanyByID(ctx context.Context, companyID uuid.UUID) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-service").Start(ctx, "GetCompanyByID")
	defer span.End()

	company, err := c.companyRepo.GetCompanyByID(ctx, companyID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get company with ID %s: %w", companyID, err)
	}

	return company, nil
}

// GetCompanyByOwnerID implements ports.CompanyService
func (c *companyService) GetCompanyByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-service").Start(ctx, "GetCompanyByOwnerID")
	defer span.End()

	company, err := c.companyRepo.GetCompanyByOwnerID(ctx, ownerID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get company with owner ID %s: %w", ownerID, err)
	}

	return company, nil
}

// UpdateCompany implements ports.CompanyService
func (c *companyService) UpdateCompany(ctx context.Context, company domain.Company) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-service").Start(ctx, "UpdateCompany")
	defer span.End()

	// Verify company exists
	existingCompany, err := c.companyRepo.GetCompanyByID(ctx, company.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company with ID %s: %w", company.ID, err)
	}

	// Preserve fields that shouldn't be updated through this method
	company.OwnerID = existingCompany.OwnerID
	company.CreatedAt = existingCompany.CreatedAt
	company.UpdatedAt = time.Now()

	// Update the company
	updatedCompany, err := c.companyRepo.UpdateCompany(ctx, company)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	return updatedCompany, nil
}

// UpdateCompanyKYB implements ports.CompanyService
func (c *companyService) UpdateCompanyKYB(ctx context.Context, companyID uuid.UUID, kybStatus string, verifiedAt *time.Time, method, provider, rejectionReason *string) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-service").Start(ctx, "UpdateCompanyKYB")
	defer span.End()

	// Verify company exists
	_, err := c.companyRepo.GetCompanyByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company with ID %s: %w", companyID, err)
	}

	// Update KYB status
	updatedCompany, err := c.companyRepo.UpdateCompanyKYB(ctx, companyID, kybStatus, verifiedAt, method, provider, rejectionReason)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to update company KYB: %w", err)
	}

	return updatedCompany, nil
}

// DeleteCompany implements ports.CompanyService
// func (c *companyService) DeleteCompany(ctx context.Context, companyID uuid.UUID) error {
// 	ctx, span := tracing.Tracer("company-service").Start(ctx, "DeleteCompany")
// 	defer span.End()

// 	// Verify company exists
// 	_, err := c.companyRepo.GetCompanyByID(ctx, companyID)
// 	if err != nil {
// 		return fmt.Errorf("failed to get company with ID %s: %w", companyID, err)
// 	}

// 	// Delete the company
// 	err = c.companyRepo.DeleteCompany(ctx, companyID)
// 	if err != nil {
// 		span.RecordError(err)
// 		return fmt.Errorf("failed to delete company: %w", err)
// 	}

// 	return nil
// }

// ListCompanies implements ports.CompanyService
func (c *companyService) ListCompanies(ctx context.Context, limit, offset int) ([]*domain.CompanyWithOwner, error) {
	ctx, span := tracing.Tracer("company-service").Start(ctx, "ListCompanies")
	defer span.End()

	companies, err := c.companyRepo.ListCompanies(ctx, limit, offset)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}

	return companies, nil
}