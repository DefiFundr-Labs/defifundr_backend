package userrepo

import (
	"context"
	"errors"
	"fmt"

	db "github.com/demola234/defifundr/db/sqlc"
	userdomain "github.com/demola234/defifundr/internal/features/user/domain"
	userport "github.com/demola234/defifundr/internal/features/user/port"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UserRepository implements userport.UserRepository using SQLC.
type UserRepository struct {
	store db.Queries
}

// New creates a new UserRepository.
func New(store db.Queries) userport.UserRepository {
	return &UserRepository{store: store}
}

func (r *UserRepository) CreateUser(ctx context.Context, user userdomain.User) (*userdomain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "CreateUser")
	defer span.End()

	params := db.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: toPgText(user.PasswordHash),
		AccountType:  user.AccountType,
		AuthProvider: toPgText(user.AuthProvider),
		ProviderID:   toPgText(user.ProviderID),
	}

	dbUser, err := r.store.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapDBUserToDomain(dbUser), nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*userdomain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "GetUserByID")
	defer span.End()

	dbUser, err := r.store.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u := mapDBUserToDomain(dbUser)
	if personalUser, err := r.store.GetPersonalUserByID(ctx, id); err == nil {
		enrichWithPersonal(u, personalUser)
	}
	return u, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*userdomain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "GetUserByEmail")
	defer span.End()

	dbUser, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	u := mapDBUserToDomain(dbUser)
	if personalUser, err := r.store.GetPersonalUserByID(ctx, dbUser.ID); err == nil {
		enrichWithPersonal(u, personalUser)
	}
	return u, nil
}

func (r *UserRepository) GetUserCompanyInfo(ctx context.Context, id uuid.UUID) (*userdomain.CompanyInfo, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "GetUserCompanyInfo")
	defer span.End()
	companies, err := r.store.GetCompaniesByOwner(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get company info: %w", err)
	}
	if len(companies) == 0 {
		return &userdomain.CompanyInfo{UserID: id}, nil
	}
	return mapCompanyToDomain(companies[0], id), nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user userdomain.User) (*userdomain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdateUser")
	defer span.End()
	dbUser, err := r.store.UpdateUser(ctx, db.UpdateUserParams{
		PasswordHash: toPgText(user.PasswordHash),
		AuthProvider: toPgText(user.AuthProvider),
		ProviderID:   toPgText(user.ProviderID),
		ID:           user.ID,
	})
	if err != nil {
		return nil, err
	}
	return mapDBUserToDomain(dbUser), nil
}

func (r *UserRepository) UpdateUserPersonalDetails(ctx context.Context, user userdomain.User) (*userdomain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdateUserPersonalDetails")
	defer span.End()

	params := db.UpdatePersonalUserParams{
		ID:                 user.ID,
		FirstName:          toPgText(user.FirstName),
		LastName:           toPgText(user.LastName),
		Nationality:        toPgText(user.Nationality),
		ResidentialCountry: toPgTextPtr(user.ResidentialCountry),
		JobRole:            toPgTextPtr(user.JobRole),
		EmploymentType:     toPgTextPtr(user.EmploymentType),
		Gender:             toPgTextPtr(user.Gender),
		UserAddress:        toPgTextPtr(user.UserAddress),
		UserCity:           toPgTextPtr(user.UserCity),
		UserPostalCode:     toPgTextPtr(user.UserPostalCode),
		PhoneNumber:        toPgTextPtr(user.PhoneNumber),
	}

	_, err := r.store.UpdatePersonalUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(ctx, user.ID)
}

func (r *UserRepository) UpdateUserAddressDetails(ctx context.Context, user userdomain.User) (*userdomain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdateUserAddressDetails")
	defer span.End()

	params := db.UpdatePersonalUserParams{
		ID:             user.ID,
		UserAddress:    toPgTextPtr(user.UserAddress),
		UserCity:       toPgTextPtr(user.UserCity),
		UserPostalCode: toPgTextPtr(user.UserPostalCode),
	}

	_, err := r.store.UpdatePersonalUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.GetUserByID(ctx, user.ID)
}

func (r *UserRepository) UpdateUserBusinessDetails(ctx context.Context, companyInfo userdomain.CompanyInfo) (*userdomain.CompanyInfo, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdateUserBusinessDetails")
	defer span.End()
	name := ""
	if companyInfo.CompanyName != nil {
		name = *companyInfo.CompanyName
	}
	// Check if company exists for this owner
	existing, err := r.store.GetCompaniesByOwner(ctx, companyInfo.UserID)
	if err != nil || len(existing) == 0 {
		// Create new company
		params := db.CreateCompanyParams{
			OwnerID:     companyInfo.UserID,
			CompanyName: name,
		}
		if companyInfo.CompanySize != nil {
			params.CompanySize = pgtype.Text{String: *companyInfo.CompanySize, Valid: true}
		}
		if companyInfo.CompanyIndustry != nil {
			params.CompanyIndustry = pgtype.Text{String: *companyInfo.CompanyIndustry, Valid: true}
		}
		if companyInfo.CompanyDescription != nil {
			params.CompanyDescription = pgtype.Text{String: *companyInfo.CompanyDescription, Valid: true}
		}
		if companyInfo.CompanyHeadquarters != nil {
			params.CompanyHeadquarters = pgtype.Text{String: *companyInfo.CompanyHeadquarters, Valid: true}
		}
		c, err := r.store.CreateCompany(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to create company: %w", err)
		}
		return mapCompanyToDomain(c, companyInfo.UserID), nil
	}
	// Update existing company
	params := db.UpdateCompanyParams{
		ID:          existing[0].ID,
		CompanyName: name,
	}
	if companyInfo.CompanySize != nil {
		params.CompanySize = pgtype.Text{String: *companyInfo.CompanySize, Valid: true}
	}
	if companyInfo.CompanyIndustry != nil {
		params.CompanyIndustry = pgtype.Text{String: *companyInfo.CompanyIndustry, Valid: true}
	}
	if companyInfo.CompanyDescription != nil {
		params.CompanyDescription = pgtype.Text{String: *companyInfo.CompanyDescription, Valid: true}
	}
	if companyInfo.CompanyHeadquarters != nil {
		params.CompanyHeadquarters = pgtype.Text{String: *companyInfo.CompanyHeadquarters, Valid: true}
	}
	c, err := r.store.UpdateCompany(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update company: %w", err)
	}
	return mapCompanyToDomain(c, companyInfo.UserID), nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdatePassword")
	defer span.End()
	_, err := r.store.UpdateUser(ctx, db.UpdateUserParams{
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
		ID:           userID,
	})
	return err
}

func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	_, err := r.store.GetUserByEmail(ctx, email)
	return err == nil, nil
}

func (r *UserRepository) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "DeactivateUser")
	defer span.End()
	_, err := r.store.UpdateUser(ctx, db.UpdateUserParams{
		AccountStatus: pgtype.Text{String: "deactivated", Valid: true},
		ID:            id,
	})
	return err
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "DeleteUser")
	defer span.End()
	return r.store.SoftDeleteUser(ctx, id)
}

func (r *UserRepository) SetMFASecret(ctx context.Context, userID uuid.UUID, secret string) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "SetMFASecret")
	defer span.End()
	_, err := r.store.UpdateUser(ctx, db.UpdateUserParams{
		TwoFactorEnabled: pgtype.Bool{Bool: true, Valid: true},
		TwoFactorMethod:  pgtype.Text{String: secret, Valid: true},
		ID:               userID,
	})
	return err
}

func (r *UserRepository) GetMFASecret(ctx context.Context, userID uuid.UUID) (string, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "GetMFASecret")
	defer span.End()
	dbUser, err := r.store.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}
	if !dbUser.TwoFactorMethod.Valid {
		return "", errors.New("MFA not configured for user")
	}
	return dbUser.TwoFactorMethod.String, nil
}

// --- helpers ---

func toPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func toPgTextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func getStr(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func mapDBUserToDomain(dbUser db.Users) *userdomain.User {
	u := &userdomain.User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		PasswordHash: getStr(dbUser.PasswordHash),
		AccountType:  dbUser.AccountType,
		AuthProvider: getStr(dbUser.AuthProvider),
		ProviderID:   getStr(dbUser.ProviderID),
	}
	if dbUser.PasswordHash.Valid {
		pw := dbUser.PasswordHash.String
		u.Password = &pw
	}
	if dbUser.CreatedAt.Valid {
		u.CreatedAt = dbUser.CreatedAt.Time
	}
	if dbUser.UpdatedAt.Valid {
		u.UpdatedAt = dbUser.UpdatedAt.Time
	}
	return u
}

func mapCompanyToDomain(c db.Companies, userID uuid.UUID) *userdomain.CompanyInfo {
	info := &userdomain.CompanyInfo{UserID: userID}
	if c.CompanyName != "" {
		info.CompanyName = &c.CompanyName
	}
	if c.CompanySize.Valid {
		info.CompanySize = &c.CompanySize.String
	}
	if c.CompanyIndustry.Valid {
		info.CompanyIndustry = &c.CompanyIndustry.String
	}
	if c.CompanyDescription.Valid {
		info.CompanyDescription = &c.CompanyDescription.String
	}
	if c.CompanyHeadquarters.Valid {
		info.CompanyHeadquarters = &c.CompanyHeadquarters.String
	}
	if c.CompanyWebsite.Valid {
		info.CompanyWebsite = &c.CompanyWebsite.String
	}
	return info
}

func enrichWithPersonal(u *userdomain.User, p db.GetPersonalUserByIDRow) {
	u.FirstName = getStr(p.FirstName)
	u.LastName = getStr(p.LastName)
	u.Nationality = getStr(p.Nationality)
	u.PersonalAccountType = getStr(p.PersonalAccountType)
	if p.Gender.Valid {
		u.Gender = &p.Gender.String
	}
	if p.ProfilePicture.Valid {
		u.ProfilePicture = &p.ProfilePicture.String
	}
	if p.ResidentialCountry.Valid {
		u.ResidentialCountry = &p.ResidentialCountry.String
	}
	if p.JobRole.Valid {
		u.JobRole = &p.JobRole.String
	}
	if p.EmploymentType.Valid {
		u.EmploymentType = &p.EmploymentType.String
	}
	if p.UserAddress.Valid {
		u.UserAddress = &p.UserAddress.String
	}
	if p.UserCity.Valid {
		u.UserCity = &p.UserCity.String
	}
	if p.UserPostalCode.Valid {
		u.UserPostalCode = &p.UserPostalCode.String
	}
	if p.PhoneNumber.Valid {
		u.PhoneNumber = &p.PhoneNumber.String
	}
	if p.PhoneNumberVerified.Valid {
		u.PhoneNumberVerified = &p.PhoneNumberVerified.Bool
	}
}
