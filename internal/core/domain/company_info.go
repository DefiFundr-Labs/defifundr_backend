package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Company represents a company entity
type Company struct {
	ID                        uuid.UUID  `json:"id"`
	OwnerID                   uuid.UUID  `json:"owner_id"`
	CompanyName               string     `json:"company_name"`
	CompanyEmail              *string    `json:"company_email,omitempty"`
	CompanyPhone              *string    `json:"company_phone,omitempty"`
	CompanySize               *string    `json:"company_size,omitempty"`
	CompanyIndustry           *string    `json:"company_industry,omitempty"`
	CompanyDescription        *string    `json:"company_description,omitempty"`
	CompanyHeadquarters       *string    `json:"company_headquarters,omitempty"`
	CompanyLogo               *string    `json:"company_logo,omitempty"`
	CompanyWebsite            *string    `json:"company_website,omitempty"`
	PrimaryContactName        *string    `json:"primary_contact_name,omitempty"`
	PrimaryContactEmail       *string    `json:"primary_contact_email,omitempty"`
	PrimaryContactPhone       *string    `json:"primary_contact_phone,omitempty"`
	CompanyAddress            *string    `json:"company_address,omitempty"`
	CompanyCity               *string    `json:"company_city,omitempty"`
	CompanyPostalCode         *string    `json:"company_postal_code,omitempty"`
	CompanyCountry            *string    `json:"company_country,omitempty"`
	CompanyRegistrationNumber *string    `json:"company_registration_number,omitempty"`
	RegistrationCountry       *string    `json:"registration_country,omitempty"`
	TaxID                     *string    `json:"tax_id,omitempty"`
	IncorporationDate         *time.Time `json:"incorporation_date,omitempty"`
	AccountStatus             string     `json:"account_status"`
	KYBStatus                 string     `json:"kyb_status"`
	KYBVerifiedAt             *time.Time `json:"kyb_verified_at,omitempty"`
	KYBVerificationMethod     *string    `json:"kyb_verification_method,omitempty"`
	KYBVerificationProvider   *string    `json:"kyb_verification_provider,omitempty"`
	KYBRejectionReason        *string    `json:"kyb_rejection_reason,omitempty"`
	LegalEntityType           *string    `json:"legal_entity_type,omitempty"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

// CompanyService defines operations for company management
type CompanyService interface {
	CreateCompany(ctx context.Context, company *Company) (*Company, error)
	GetCompanyByID(ctx context.Context, id uuid.UUID) (*Company, error)
	GetCompanyByOwnerID(ctx context.Context, ownerID uuid.UUID) (*Company, error)
	GetCompanyWithOwnerDetails(ctx context.Context, id uuid.UUID) (*CompanyWithOwner, error)
	UpdateCompany(ctx context.Context, company *Company) (*Company, error)
	UpdateCompanyKYB(ctx context.Context, id uuid.UUID, status, method, provider, rejectionReason string, verifiedAt *time.Time) (*Company, error)
	DeleteCompany(ctx context.Context, id uuid.UUID) error
	ListCompanies(ctx context.Context, limit, offset int) ([]*CompanyWithOwner, error)
	GetCompaniesByKYBStatus(ctx context.Context, kybStatus string, limit, offset int) ([]*Company, error)
	SearchCompaniesByName(ctx context.Context, searchTerm string, limit, offset int) ([]*Company, error)
	GetCompaniesByIndustry(ctx context.Context, industry string, limit, offset int) ([]*Company, error)
}

// CompanyRepository defines the data access interface for companies
type CompanyRepository interface {
	Create(ctx context.Context, company *Company) (*Company, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Company, error)
	GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*Company, error)
	Update(ctx context.Context, company *Company) (*Company, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, limit, offset int) ([]*Company, error)
	GetByKYBStatus(ctx context.Context, kybStatus string, limit, offset int) ([]*Company, error)
	SearchByName(ctx context.Context, searchTerm string, limit, offset int) ([]*Company, error)
	GetByIndustry(ctx context.Context, industry string, limit, offset int) ([]*Company, error)
}


// NewCompany creates a new Company instance
func NewCompany(ownerID uuid.UUID, companyName string) *Company {
	return &Company{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		CompanyName:   companyName,
		AccountStatus: "pending",
		KYBStatus:     "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
