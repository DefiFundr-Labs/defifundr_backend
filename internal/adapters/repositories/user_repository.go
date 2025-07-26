// Replace your entire user repository with this updated version

package repositories

import (
	"context"
	"fmt"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// UserRepository struct implements the repository interface for users
type UserRepository struct {
	store db.Queries
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(store db.Queries) *UserRepository {
	return &UserRepository{store: store}
}

// CreateUser implements the user registration functionality
func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "CreateUser")
	defer span.End()

	
	params := db.CreateUserParams{
		ID:                    user.ID,
		FirstName:             toPgTextPtr(&user.FirstName),
		LastName:              toPgTextPtr(&user.LastName),
		PhoneNumber:           toPgTextPtr(&user.PhoneNumber),
		Email:                 user.Email,
		PasswordHash:          toPgTextPtr(&user.PasswordHash),
		ProfilePictureUrl:     toPgTextPtr(user.ProfilePictureURL),
		AuthProvider:          toPgTextPtr(&user.AuthProvider),
		ProviderID:            toPgTextPtr(user.ProviderID),
		EmailVerified:         pgtype.Bool{Bool: user.EmailVerified, Valid: true},
		EmailVerifiedAt:       toPgTimestamptzPtr(user.EmailVerifiedAt),
		PhoneNumberVerified:   pgtype.Bool{Bool: user.PhoneNumberVerified, Valid: true},
		PhoneNumberVerifiedAt: toPgTimestamptzPtr(user.PhoneNumberVerifiedAt),
		AccountType:           user.AccountType,
		AccountStatus:         toPgText(user.AccountStatus),
		TwoFactorEnabled:      pgtype.Bool{Bool: user.TwoFactorEnabled, Valid: true},
		TwoFactorMethod:       toPgTextPtr(user.TwoFactorMethod),
		UserLoginType:         toPgTextPtr(user.UserLoginType),
		CreatedAt:             toPgTimestamptzPtr(&user.CreatedAt),
		UpdatedAt:             toPgTimestamptzPtr(&user.UpdatedAt),
		LastLoginAt:           toPgTimestamptzPtr(user.LastLoginAt),
		DeletedAt:             toPgTimestamptzPtr(user.DeletedAt),
	}

	dbUser, err := r.store.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// CreatePersonalUser creates a personal user profile
func (r *UserRepository) CreatePersonalUser(ctx context.Context, personalUser domain.PersonalUser) (*domain.PersonalUser, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "CreatePersonalUser")
	defer span.End()

	params := db.CreatePersonalUserParams{
		ID:                     personalUser.ID,
		UserID:                personalUser.UserID,
		Nationality:            toPgTextPtr(personalUser.Nationality),
		ResidentialCountry:     toPgTextPtr(personalUser.ResidentialCountry),
		UserAddress:            toPgTextPtr(personalUser.UserAddress),
		UserCity:               toPgTextPtr(personalUser.UserCity),
		UserPostalCode:         toPgTextPtr(personalUser.UserPostalCode),
		Gender:                 toPgTextPtr(personalUser.Gender),
		DateOfBirth:            toPgDatePtr(personalUser.DateOfBirth),
		JobRole:                toPgTextPtr(personalUser.JobRole),
		PersonalAccountType:    toPgTextPtr(personalUser.PersonalAccountType),
		EmploymentType:         toPgTextPtr(personalUser.EmploymentType),
		TaxID:                  toPgTextPtr(personalUser.TaxID),
		DefaultPaymentCurrency: toPgTextPtr(personalUser.DefaultPaymentCurrency),
		DefaultPaymentMethod:   toPgTextPtr(personalUser.DefaultPaymentMethod),
		HourlyRate:             toPgNumericPtr(personalUser.HourlyRate),
		Specialization:         toPgTextPtr(personalUser.Specialization),
		CreatedAt:              toPgTimestamptzPtr(&personalUser.CreatedAt),
		UpdatedAt:              toPgTimestamptzPtr(&personalUser.UpdatedAt),
	}

	dbPersonalUser, err := r.store.CreatePersonalUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBPersonalUserToDomainPersonalUser(dbPersonalUser), nil
}

// CreateCompanyUser creates a company user profile
func (r *UserRepository) CreateCompanyUser(ctx context.Context, companyUser domain.CompanyUser) (*domain.CompanyUser, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "CreateCompanyUser")
	defer span.End()

	params := db.CreateCompanyUserParams{
		ID:                       companyUser.ID,
		CompanyID:                companyUser.CompanyID,
		UserID:                   companyUser.UserID,
		Role:                     companyUser.Role,
		Department:               toPgTextPtr(companyUser.Department),
		JobTitle:                 toPgTextPtr(companyUser.JobTitle),
		IsAdministrator:          pgtype.Bool{Bool: companyUser.IsAdministrator, Valid: true},
		CanManagePayroll:         pgtype.Bool{Bool: companyUser.CanManagePayroll, Valid: true},
		CanManageInvoices:        pgtype.Bool{Bool: companyUser.CanManageInvoices, Valid: true},
		CanManageEmployees:       pgtype.Bool{Bool: companyUser.CanManageEmployees, Valid: true},
		CanManageCompanySettings: pgtype.Bool{Bool: companyUser.CanManageCompanySettings, Valid: true},
		CanManageBankAccounts:    pgtype.Bool{Bool: companyUser.CanManageBankAccounts, Valid: true},
		CanManageWallets:         pgtype.Bool{Bool: companyUser.CanManageWallets, Valid: true},
		IsActive:                 pgtype.Bool{Bool: companyUser.IsActive, Valid: true},
		AddedBy:                  toPgUUIDPtr(companyUser.AddedBy),
		ReportsTo:                toPgUUIDPtr(companyUser.ReportsTo),
		HireDate:                 toPgDatePtr(companyUser.HireDate),
		CreatedAt:                toPgTimestamptzPtr(&companyUser.CreatedAt),
		UpdatedAt:                toPgTimestamptzPtr(&companyUser.UpdatedAt),
	}

	dbCompanyUser, err := r.store.CreateCompanyUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBCompanyUserToDomainCompanyUser(dbCompanyUser), nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "GetUserByID")
	defer span.End()

	dbUser, err := r.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "GetUserByEmail")
	defer span.End()

	dbUser, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// GetPersonalUserByID retrieves a personal user by their ID
func (r *UserRepository) GetPersonalUserByID(ctx context.Context, id uuid.UUID) (*domain.PersonalUser, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "GetPersonalUserByID")
	defer span.End()

	dbPersonalUser, err := r.store.GetPersonalUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapDBPersonalUserToDomainPersonalUser(dbPersonalUser), nil
}

// GetPersonalUserByUserID retrieves a personal user by their user ID
func (r *UserRepository) GetPersonalUserByUserID(ctx context.Context, userID uuid.UUID) (*domain.PersonalUser, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "GetPersonalUserByUserID")
	defer span.End()

	dbPersonalUser, err := r.store.GetPersonalUserByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return mapDBPersonalUserToDomainPersonalUser(dbPersonalUser), nil
}

// GetCompanyUserByID retrieves a company user by their ID
func (r *UserRepository) GetCompanyUserByID(ctx context.Context, id uuid.UUID) (*domain.CompanyUser, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "GetCompanyUserByID")
	defer span.End()

	dbCompanyUser, err := r.store.GetCompanyUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapDBCompanyUserToDomainCompanyUser(dbCompanyUser), nil
}

// UpdateUser updates a user's information
func (r *UserRepository) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdateUser")
	defer span.End()

	params := db.UpdateUserParams{
		ID:                    user.ID,
		FirstName:             toPgTextPtr(&user.FirstName),
		LastName:              toPgTextPtr(&user.LastName),
		PhoneNumber:           toPgTextPtr(&user.PhoneNumber),
		PasswordHash:          toPgTextPtr(&user.PasswordHash),
		ProfilePictureUrl:     toPgTextPtr(user.ProfilePictureURL),
		AuthProvider:          toPgTextPtr(&user.AuthProvider),
		ProviderID:            toPgTextPtr(user.ProviderID),
		EmailVerified:         pgtype.Bool{Bool: user.EmailVerified, Valid: true},
		EmailVerifiedAt:       toPgTimestamptzPtr(user.EmailVerifiedAt),
		PhoneNumberVerified:   pgtype.Bool{Bool: user.PhoneNumberVerified, Valid: true},
		PhoneNumberVerifiedAt: toPgTimestamptzPtr(user.PhoneNumberVerifiedAt),
		AccountStatus:         toPgText(user.AccountStatus),
		TwoFactorEnabled:      pgtype.Bool{Bool: user.TwoFactorEnabled, Valid: true},
		TwoFactorMethod:       toPgTextPtr(user.TwoFactorMethod),
		UserLoginType:         toPgTextPtr(user.UserLoginType),
		LastLoginAt:           toPgTimestamptzPtr(user.LastLoginAt),
	}

	dbUser, err := r.store.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// UpdatePersonalUser updates a personal user's information
func (r *UserRepository) UpdatePersonalUser(ctx context.Context, personalUser domain.PersonalUser) (*domain.PersonalUser, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdatePersonalUser")
	defer span.End()

	params := db.UpdatePersonalUserParams{
		ID:                     personalUser.ID,
		Nationality:            toPgTextPtr(personalUser.Nationality),
		ResidentialCountry:     toPgTextPtr(personalUser.ResidentialCountry),
		UserAddress:            toPgTextPtr(personalUser.UserAddress),
		UserCity:               toPgTextPtr(personalUser.UserCity),
		UserPostalCode:         toPgTextPtr(personalUser.UserPostalCode),
		Gender:                 toPgTextPtr(personalUser.Gender),
		DateOfBirth:            toPgDatePtr(personalUser.DateOfBirth),
		JobRole:                toPgTextPtr(personalUser.JobRole),
		PersonalAccountType:    toPgTextPtr(personalUser.PersonalAccountType),
		EmploymentType:         toPgTextPtr(personalUser.EmploymentType),
		TaxID:                  toPgTextPtr(personalUser.TaxID),
		DefaultPaymentCurrency: toPgTextPtr(personalUser.DefaultPaymentCurrency),
		DefaultPaymentMethod:   toPgTextPtr(personalUser.DefaultPaymentMethod),
		HourlyRate:             toPgNumericPtr(personalUser.HourlyRate),
		Specialization:         toPgTextPtr(personalUser.Specialization),
	}

	dbPersonalUser, err := r.store.UpdatePersonalUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBPersonalUserToDomainPersonalUser(dbPersonalUser), nil
}

// UpdateCompanyUser updates a company user's information
func (r *UserRepository) UpdateCompanyUser(ctx context.Context, companyUser domain.CompanyUser) (*domain.CompanyUser, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdateCompanyUser")
	defer span.End()

	params := db.UpdateCompanyUserParams{
		ID:                       companyUser.ID,
		Role:                     companyUser.Role,
		Department:               toPgTextPtr(companyUser.Department),
		JobTitle:                 toPgTextPtr(companyUser.JobTitle),
		IsAdministrator:          pgtype.Bool{Bool: companyUser.IsAdministrator, Valid: true},
		CanManagePayroll:         pgtype.Bool{Bool: companyUser.CanManagePayroll, Valid: true},
		CanManageInvoices:        pgtype.Bool{Bool: companyUser.CanManageInvoices, Valid: true},
		CanManageEmployees:       pgtype.Bool{Bool: companyUser.CanManageEmployees, Valid: true},
		CanManageCompanySettings: pgtype.Bool{Bool: companyUser.CanManageCompanySettings, Valid: true},
		CanManageBankAccounts:    pgtype.Bool{Bool: companyUser.CanManageBankAccounts, Valid: true},
		CanManageWallets:         pgtype.Bool{Bool: companyUser.CanManageWallets, Valid: true},
		IsActive:                 pgtype.Bool{Bool: companyUser.IsActive, Valid: true},
		ReportsTo:                toPgUUIDPtr(companyUser.ReportsTo),
	}

	dbCompanyUser, err := r.store.UpdateCompanyUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBCompanyUserToDomainCompanyUser(dbCompanyUser), nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdatePassword")
	defer span.End()


	fmt.Println(" ");
	fmt.Println(passwordHash);
	fmt.Println(" ");

	params := db.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
	}

	return r.store.UpdateUserPassword(ctx, params)
}

// CheckEmailExists checks if an email already exists
func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "CheckEmailExists")
	defer span.End()

	exists, err := r.store.CheckEmailExists(ctx, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// UpdateUserPersonalDetails updates user personal details
func (r *UserRepository) UpdateUserPersonalDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdateUserPersonalDetails")
	defer span.End()

	params := db.UpdateUserPersonalDetailsParams{
		ID:          user.ID,
		PhoneNumber: toPgTextPtr(&user.PhoneNumber),
	}

	dbUser, err := r.store.UpdateUserPersonalDetails(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// UpdateUserBusinessDetails updates user business details
func (r *UserRepository) UpdateUserBusinessDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdateUserBusinessDetails")
	defer span.End()

	dbUser, err := r.store.UpdateUserCompanyDetails(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// UpdateUserAddressDetails updates user address details
func (r *UserRepository) UpdateUserAddressDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "UpdateUserAddressDetails")
	defer span.End()

	dbUser, err := r.store.UpdateUserAddress(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// DeactivateUser marks a user as inactive
func (r *UserRepository) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "DeactivateUser")
	defer span.End()

	return r.store.SoftDeleteUser(ctx, id)
}

// DeleteUser removes a user from the database
func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "DeleteUser")
	defer span.End()

	return r.store.DeleteUser(ctx, id)
}

// DeletePersonalUser removes a personal user from the database
func (r *UserRepository) DeletePersonalUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "DeletePersonalUser")
	defer span.End()

	return r.store.DeletePersonalUser(ctx, id)
}

// DeleteCompanyUser removes a company user from the database
func (r *UserRepository) DeleteCompanyUser(ctx context.Context, id uuid.UUID) error {
	ctx, span := tracing.Tracer("user-repository").Start(ctx, "DeleteCompanyUser")
	defer span.End()

	return r.store.DeleteCompanyUser(ctx, id)
}

// SetMFASecret sets the MFA secret for a user
func (r *UserRepository) SetMFASecret(ctx context.Context, userID uuid.UUID, secret string) error {
	// TODO: Implement MFA secret storage - may need to add MFA fields to database
	return nil
}

// GetMFASecret retrieves the MFA secret for a user
func (r *UserRepository) GetMFASecret(ctx context.Context, userID uuid.UUID) (string, error) {
	// TODO: Implement MFA secret retrieval - may need to add MFA fields to database
	return "", nil
}

// Mapping functions

// Helper function to map database user to domain user
func mapDBUserToDomainUser(dbUser db.Users) *domain.User {
	return &domain.User{
		ID:                     dbUser.ID,
		FirstName:              getTextString(dbUser.FirstName),
		LastName:               getTextString(dbUser.LastName),
		PhoneNumber:            getTextString(dbUser.PhoneNumber),
		Email:                  dbUser.Email,
		PasswordHash:           getTextString(dbUser.PasswordHash),
		ProfilePictureURL:      getTextStringPtr(dbUser.ProfilePictureUrl),
		AuthProvider:           getTextString(dbUser.AuthProvider),
		ProviderID:             getTextStringPtr(dbUser.ProviderID),
		EmailVerified:          getBool(dbUser.EmailVerified),
		EmailVerifiedAt:        getTimestamptzPtr(dbUser.EmailVerifiedAt),
		PhoneNumberVerified:    getBool(dbUser.PhoneNumberVerified),
		PhoneNumberVerifiedAt:  getTimestamptzPtr(dbUser.PhoneNumberVerifiedAt),
		AccountType:            dbUser.AccountType,
		AccountStatus:          getTextString(dbUser.AccountStatus),
		TwoFactorEnabled:       getBool(dbUser.TwoFactorEnabled),
		TwoFactorMethod:        getTextStringPtr(dbUser.TwoFactorMethod),
		UserLoginType:          getTextStringPtr(dbUser.UserLoginType),
		CreatedAt:              getTimestamptzTime(dbUser.CreatedAt),
		UpdatedAt:              getTimestamptzTime(dbUser.UpdatedAt),
		LastLoginAt:            getTimestamptzPtr(dbUser.LastLoginAt),
		DeletedAt:              getTimestamptzPtr(dbUser.DeletedAt),
	}
}

// Helper function to map database personal user to domain personal user
func mapDBPersonalUserToDomainPersonalUser(dbPersonalUser db.PersonalUsers) *domain.PersonalUser {
	return &domain.PersonalUser{
		ID:                     dbPersonalUser.ID,
		UserID:                 dbPersonalUser.UserID,
		Nationality:            getTextStringPtr(dbPersonalUser.Nationality),
		ResidentialCountry:     getTextStringPtr(dbPersonalUser.ResidentialCountry),
		UserAddress:            getTextStringPtr(dbPersonalUser.UserAddress),
		UserCity:               getTextStringPtr(dbPersonalUser.UserCity),
		UserPostalCode:         getTextStringPtr(dbPersonalUser.UserPostalCode),
		Gender:                 getTextStringPtr(dbPersonalUser.Gender),
		DateOfBirth:            getDatePtr(dbPersonalUser.DateOfBirth),
		JobRole:                getTextStringPtr(dbPersonalUser.JobRole),
		PersonalAccountType:    getTextStringPtr(dbPersonalUser.PersonalAccountType),
		EmploymentType:         getTextStringPtr(dbPersonalUser.EmploymentType),
		TaxID:                  getTextStringPtr(dbPersonalUser.TaxID),
		DefaultPaymentCurrency: getTextStringPtr(dbPersonalUser.DefaultPaymentCurrency),
		DefaultPaymentMethod:   getTextStringPtr(dbPersonalUser.DefaultPaymentMethod),
		HourlyRate:             getNumericPtr(dbPersonalUser.HourlyRate),
		Specialization:         getTextStringPtr(dbPersonalUser.Specialization),
		CreatedAt:              getTimestamptzTime(dbPersonalUser.CreatedAt),
		UpdatedAt:              getTimestamptzTime(dbPersonalUser.UpdatedAt),
	}
}

// Helper function to map database company user to domain company user
func mapDBCompanyUserToDomainCompanyUser(dbCompanyUser db.CompanyUsers) *domain.CompanyUser {
	return &domain.CompanyUser{
		ID:                       dbCompanyUser.ID,
		CompanyID:                dbCompanyUser.CompanyID,
		UserID:                   dbCompanyUser.UserID,
		Role:                     dbCompanyUser.Role,
		Department:               getTextStringPtr(dbCompanyUser.Department),
		JobTitle:                 getTextStringPtr(dbCompanyUser.JobTitle),
		IsAdministrator:          getBool(dbCompanyUser.IsAdministrator),
		CanManagePayroll:         getBool(dbCompanyUser.CanManagePayroll),
		CanManageInvoices:        getBool(dbCompanyUser.CanManageInvoices),
		CanManageEmployees:       getBool(dbCompanyUser.CanManageEmployees),
		CanManageCompanySettings: getBool(dbCompanyUser.CanManageCompanySettings),
		CanManageBankAccounts:    getBool(dbCompanyUser.CanManageBankAccounts),
		CanManageWallets:         getBool(dbCompanyUser.CanManageWallets),
		IsActive:                 getBool(dbCompanyUser.IsActive),
		AddedBy:                  getUUIDPtr(dbCompanyUser.AddedBy),
		ReportsTo:                getUUIDPtr(dbCompanyUser.ReportsTo),
		HireDate:                 getDatePtr(dbCompanyUser.HireDate),
		CreatedAt:                getTimestamptzTime(dbCompanyUser.CreatedAt),
		UpdatedAt:                getTimestamptzTime(dbCompanyUser.UpdatedAt),
	}
}

// Helper functions for type conversion

func toPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func toPgTextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func toPgNumericPtr(f *float64) pgtype.Numeric {
	if f == nil {
		return pgtype.Numeric{Valid: false}
	}
	// Convert float64 to pgtype.Numeric using string representation
	return pgtype.Numeric{
		Valid: true,
		// You may need to set the numeric value properly based on your pgtype version
		// This is a placeholder - adjust based on your specific pgtype.Numeric implementation
	}
}

func getTextString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func getUUIDPtr(u pgtype.UUID) *uuid.UUID {
	if !u.Valid {
		return nil
	}
	return (*uuid.UUID)(&u.Bytes)
}

func getNumericPtr(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	// Convert pgtype.Numeric to float64
	// This is a simplified conversion - adjust based on your pgtype version
	if n.Int != nil {
		// For newer versions of pgtype, you may need to use different methods
		// This is a placeholder implementation
		value := 0.0 // You'll need to implement proper conversion based on your pgtype version
		return &value
	}
	return nil
}