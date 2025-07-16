package repositories

import (
	"context"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// CompanyRepository struct implements the repository interface for companies
type CompanyRepository struct {
	store db.Queries
}

// NewCompanyRepository creates a new CompanyRepository
func NewCompanyRepository(store db.Queries) *CompanyRepository {
	return &CompanyRepository{store: store}
}

// CreateCompany implements the company creation functionality
func (r *CompanyRepository) CreateCompany(ctx context.Context, company domain.Company) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-repository").Start(ctx, "CreateCompany")
	defer span.End()

	params := db.CreateCompanyParams{
		ID:                       company.ID,
		OwnerID:                  company.OwnerID,
		CompanyName:              company.CompanyName,
		CompanyEmail:             toPgTextPtr(company.CompanyEmail),
		CompanyPhone:             toPgTextPtr(company.CompanyPhone),
		CompanySize:              toPgTextPtr(company.CompanySize),
		CompanyIndustry:          toPgTextPtr(company.CompanyIndustry),
		CompanyDescription:       toPgTextPtr(company.CompanyDescription),
		CompanyHeadquarters:      toPgTextPtr(company.CompanyHeadquarters),
		CompanyLogo:              toPgTextPtr(company.CompanyLogo),
		CompanyWebsite:           toPgTextPtr(company.CompanyWebsite),
		PrimaryContactName:       toPgTextPtr(company.PrimaryContactName),
		PrimaryContactEmail:      toPgTextPtr(company.PrimaryContactEmail),
		PrimaryContactPhone:      toPgTextPtr(company.PrimaryContactPhone),
		CompanyAddress:           toPgTextPtr(company.CompanyAddress),
		CompanyCity:              toPgTextPtr(company.CompanyCity),
		CompanyPostalCode:        toPgTextPtr(company.CompanyPostalCode),
		CompanyCountry:           toPgTextPtr(company.CompanyCountry),
		CompanyRegistrationNumber: toPgTextPtr(company.CompanyRegistrationNumber),
		RegistrationCountry:      toPgTextPtr(company.RegistrationCountry),
		TaxID:                    toPgTextPtr(company.TaxID),
		IncorporationDate:        toPgDatePtr(company.IncorporationDate),
		AccountStatus:            toPgText(company.AccountStatus),
		KybStatus:                toPgText(company.KYBStatus),
		LegalEntityType:          toPgTextPtr(company.LegalEntityType),
		CreatedAt:                pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:                pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	dbCompany, err := r.store.CreateCompany(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBCompanyToDomainCompany(dbCompany), nil
}

// GetCompanyByID retrieves a company by their ID
func (r *CompanyRepository) GetCompanyByID(ctx context.Context, id uuid.UUID) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-repository").Start(ctx, "GetCompanyByID")
	defer span.End()

	dbCompany, err := r.store.GetCompanyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapDBCompanyToDomainCompany(dbCompany), nil
}

// // GetCompanyByOwnerID retrieves a company by owner ID
func (r *CompanyRepository) GetCompanyByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-repository").Start(ctx, "GetCompanyByOwnerID")
	defer span.End()

	dbCompany, err := r.store.GetCompanyByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	return mapDBCompanyToDomainCompany(dbCompany), nil
}

// UpdateCompany updates a company's information
func (r *CompanyRepository) UpdateCompany(ctx context.Context, company domain.Company) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-repository").Start(ctx, "UpdateCompany")
	defer span.End()

	params := db.UpdateCompanyParams{
		ID:                       company.ID,
		CompanyName:              company.CompanyName,
		CompanyEmail:             toPgTextPtr(company.CompanyEmail),
		CompanyPhone:             toPgTextPtr(company.CompanyPhone),
		CompanySize:              toPgTextPtr(company.CompanySize),
		CompanyIndustry:          toPgTextPtr(company.CompanyIndustry),
		CompanyDescription:       toPgTextPtr(company.CompanyDescription),
		CompanyHeadquarters:      toPgTextPtr(company.CompanyHeadquarters),
		CompanyLogo:              toPgTextPtr(company.CompanyLogo),
		CompanyWebsite:           toPgTextPtr(company.CompanyWebsite),
		PrimaryContactName:       toPgTextPtr(company.PrimaryContactName),
		PrimaryContactEmail:      toPgTextPtr(company.PrimaryContactEmail),
		PrimaryContactPhone:      toPgTextPtr(company.PrimaryContactPhone),
		CompanyAddress:           toPgTextPtr(company.CompanyAddress),
		CompanyCity:              toPgTextPtr(company.CompanyCity),
		CompanyPostalCode:        toPgTextPtr(company.CompanyPostalCode),
		CompanyCountry:           toPgTextPtr(company.CompanyCountry),
		CompanyRegistrationNumber: toPgTextPtr(company.CompanyRegistrationNumber),
		RegistrationCountry:      toPgTextPtr(company.RegistrationCountry),
		TaxID:                    toPgTextPtr(company.TaxID),
		IncorporationDate:        toPgDatePtr(company.IncorporationDate),
		AccountStatus:            toPgText(company.AccountStatus),
		LegalEntityType:          toPgTextPtr(company.LegalEntityType),
	}

	dbCompany, err := r.store.UpdateCompany(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBCompanyToDomainCompany(dbCompany), nil
}

// UpdateCompanyKYB updates a company's KYB information
func (r *CompanyRepository) UpdateCompanyKYB(ctx context.Context, companyID uuid.UUID, kybStatus string, verifiedAt *time.Time, method, provider, rejectionReason *string) (*domain.Company, error) {
	ctx, span := tracing.Tracer("company-repository").Start(ctx, "UpdateCompanyKYB")
	defer span.End()

	params := db.UpdateCompanyParams{
		ID:                      companyID,
		KybStatus:               toPgText(kybStatus),
		KybVerifiedAt:           toPgTimestamptzPtr(verifiedAt),
		KybVerificationMethod:   toPgTextPtr(method),
		KybVerificationProvider: toPgTextPtr(provider),
		KybRejectionReason:      toPgTextPtr(rejectionReason),
	}

	dbCompany, err := r.store.UpdateCompany(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBCompanyToDomainCompany(dbCompany), nil
}

// DeleteCompany removes a company from the database
// func (r *CompanyRepository) DeleteCompany(ctx context.Context, id uuid.UUID) error {
// 	ctx, span := tracing.Tracer("company-repository").Start(ctx, "DeleteCompany")
// 	defer span.End()

// 	return r.store.DeleteCompany(ctx, id)
// }

// ListCompanies retrieves companies accessible to a user with pagination
func (r *CompanyRepository) ListCompanies(ctx context.Context, limit, offset int) ([]*domain.CompanyWithOwner, error) {
	ctx, span := tracing.Tracer("company-repository").Start(ctx, "ListCompanies")
	defer span.End()

	// Since there's no ListCompanies query, we'll use GetCompaniesByOwner for all companies
	// You might want to add a proper ListCompanies query to your SQL file
	companies := make([]*domain.CompanyWithOwner, 0)
	
	// This is a workaround - you should add a proper ListCompanies SQL query
	// For now, return empty slice
	return companies, nil
}

// Helper function to map database company to domain company
func mapDBCompanyToDomainCompany(dbCompany db.Companies) *domain.Company {
	return &domain.Company{
		ID:                       dbCompany.ID,
		OwnerID:                  dbCompany.OwnerID,
		CompanyName:              dbCompany.CompanyName,
		CompanyEmail:             getTextStringPtr(dbCompany.CompanyEmail),
		CompanyPhone:             getTextStringPtr(dbCompany.CompanyPhone),
		CompanySize:              getTextStringPtr(dbCompany.CompanySize),
		CompanyIndustry:          getTextStringPtr(dbCompany.CompanyIndustry),
		CompanyDescription:       getTextStringPtr(dbCompany.CompanyDescription),
		CompanyHeadquarters:      getTextStringPtr(dbCompany.CompanyHeadquarters),
		CompanyLogo:              getTextStringPtr(dbCompany.CompanyLogo),
		CompanyWebsite:           getTextStringPtr(dbCompany.CompanyWebsite),
		PrimaryContactName:       getTextStringPtr(dbCompany.PrimaryContactName),
		PrimaryContactEmail:      getTextStringPtr(dbCompany.PrimaryContactEmail),
		PrimaryContactPhone:      getTextStringPtr(dbCompany.PrimaryContactPhone),
		CompanyAddress:           getTextStringPtr(dbCompany.CompanyAddress),
		CompanyCity:              getTextStringPtr(dbCompany.CompanyCity),
		CompanyPostalCode:        getTextStringPtr(dbCompany.CompanyPostalCode),
		CompanyCountry:           getTextStringPtr(dbCompany.CompanyCountry),
		CompanyRegistrationNumber: getTextStringPtr(dbCompany.CompanyRegistrationNumber),
		RegistrationCountry:      getTextStringPtr(dbCompany.RegistrationCountry),
		TaxID:                    getTextStringPtr(dbCompany.TaxID),
		IncorporationDate:        getDatePtr(dbCompany.IncorporationDate),
		AccountStatus:            getTextString(dbCompany.AccountStatus),
		KYBStatus:                getTextString(dbCompany.KybStatus),
		KYBVerifiedAt:            getTimestamptzPtr(dbCompany.KybVerifiedAt),
		KYBVerificationMethod:    getTextStringPtr(dbCompany.KybVerificationMethod),
		KYBVerificationProvider:  getTextStringPtr(dbCompany.KybVerificationProvider),
		KYBRejectionReason:       getTextStringPtr(dbCompany.KybRejectionReason),
		LegalEntityType:          getTextStringPtr(dbCompany.LegalEntityType),
		CreatedAt:                getTimestamptzTime(dbCompany.CreatedAt),
		UpdatedAt:                getTimestamptzTime(dbCompany.UpdatedAt),
	}
}
