package userrepo

import (
	"context"
	"errors"

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
	return nil, errors.New("not implemented")
}

func (r *UserRepository) UpdateUser(ctx context.Context, user userdomain.User) (*userdomain.User, error) {
	return nil, errors.New("not implemented")
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
	return nil, errors.New("not implemented")
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	return errors.New("not implemented")
}

func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	_, err := r.store.GetUserByEmail(ctx, email)
	return err == nil, nil
}

func (r *UserRepository) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	return errors.New("not implemented")
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "DeleteUser")
	defer span.End()
	return r.store.SoftDeleteUser(ctx, id)
}

func (r *UserRepository) SetMFASecret(ctx context.Context, userID uuid.UUID, secret string) error {
	return errors.New("not implemented")
}

func (r *UserRepository) GetMFASecret(ctx context.Context, userID uuid.UUID) (string, error) {
	return "", errors.New("not implemented")
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
