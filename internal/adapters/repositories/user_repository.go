package repositories

import (
	"context"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
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

// RegisterUser implements the user registration functionality
func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {

	params := db.CreateUserParams{
		ID:                  user.ID,
		Email:               user.Email,
		PasswordHash:        toPgText(user.PasswordHash),
		ProfilePicture:      toPgTextPtr(user.ProfilePicture),
		Gender:              toPgTextPtr(user.Gender),
		AccountType:         user.AccountType,
		PersonalAccountType: user.PersonalAccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Nationality:         user.Nationality,
		ResidentialCountry:  toPgTextPtr(user.ResidentialCountry),
		JobRole:             toPgTextPtr(user.JobRole),
		AuthProvider:        toPgText(user.AuthProvider),
		ProviderID:          user.ProviderID,
		EmployeeType:        toPgText(user.EmployeeType),
		EmploymentType:      toPgTextPtr(user.EmploymentType),
		CreatedAt:           pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt:           pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	// Call the database to create user
	dbUser, err := r.store.CreateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	// Map database user back to domain user
	return mapDBUserToDomainUser(dbUser), nil
}

// Create User Company Info implements the user registration functionality
func (r *UserRepository) UpdateCompanyInfo(ctx context.Context, companyInfo domain.CompanyInfo) (*domain.CompanyInfo, error) {
	params := db.UpdateUserCompanyDetailsParams{
		ID:                    companyInfo.UserID,
		CompanyName:           toPgTextPtr(companyInfo.CompanyName),
		CompanyHeadquarters:   toPgTextPtr(companyInfo.CompanyHeadquarters),
		CompanySize:           toPgTextPtr(companyInfo.CompanySize),
		CompanyIndustry:       toPgTextPtr(companyInfo.CompanyIndustry),
		CompanyDescription:    toPgTextPtr(companyInfo.CompanyDescription),
		CompanyHeadquarters_2: toPgTextPtr(companyInfo.CompanyHeadquarters),
		AccountType:           companyInfo.AccountType,
	}

	dbCompanyInfo, err := r.store.UpdateUserCompanyDetails(ctx, params)

	if err != nil {
		return nil, err
	}

	return mapDBCompanyInfoToDomainCompanyInfo(dbCompanyInfo), nil
}

func (r *UserRepository) GetUserCompanyInfo(ctx context.Context, id uuid.UUID) (*domain.CompanyInfo, error) {
	dbUser, err := r.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapDBCompanyInfoToDomainCompanyInfo(dbUser), nil
}

// GetUserByID retrieves a user by their ID
func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	dbUser, err := r.store.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// GetUserByEmail retrieves a user by their email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	dbUser, err := r.store.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// UpdateUser updates a user's information
func (r *UserRepository) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	params := db.UpdateUserParams{
		ID:                  user.ID,
		Email:               user.Email,
		ProfilePicture:      toPgTextPtr(user.ProfilePicture),
		AccountType:         user.AccountType,
		Gender:              toPgTextPtr(user.Gender),
		PersonalAccountType: user.PersonalAccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Nationality:         user.Nationality,
		ResidentialCountry:  toPgTextPtr(user.ResidentialCountry),
		JobRole:             toPgTextPtr(user.JobRole),
		EmploymentType:      toPgTextPtr(user.EmploymentType),
		AuthProvider:        toPgText(user.AuthProvider),
		ProviderID:          user.ProviderID,
		UserAddress:         toPgTextPtr(user.UserAddress),
		UserCity:            toPgTextPtr(user.UserCity),
		UserPostalCode:      toPgTextPtr(user.UserPostalCode),
	}

	dbUser, err := r.store.UpdateUser(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, passwordHash string) error {
	params := db.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: pgtype.Text{String: passwordHash, Valid: true},
	}

	return r.store.UpdateUserPassword(ctx, params)
}

// CheckEmailExists checks if an email already exists
func (r *UserRepository) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := r.store.CheckEmailExists(ctx, email)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// Helper function to get a string pointer
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// Helper function to get a string from pgtype.Text
func getTextString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

// Helper function to map database user to domain user
func mapDBUserToDomainUser(dbUser db.Users) *domain.User {
	var passwordPtr *string
	if dbUser.PasswordHash.Valid {
		passwordPtr = &dbUser.PasswordHash.String
	}

	return &domain.User{
		ID:                  dbUser.ID,
		Email:               dbUser.Email,
		Password:            passwordPtr,
		PasswordHash:        getTextString(dbUser.PasswordHash),
		ProfilePicture:      strPtr(getTextString(dbUser.ProfilePicture)),
		AccountType:         dbUser.AccountType,
		Gender:              strPtr(getTextString(dbUser.Gender)),
		PersonalAccountType: dbUser.PersonalAccountType,
		FirstName:           dbUser.FirstName,
		LastName:            dbUser.LastName,
		Nationality:         dbUser.Nationality,
		ResidentialCountry:  strPtr(getTextString(dbUser.ResidentialCountry)),
		JobRole:             strPtr(getTextString(dbUser.JobRole)),
		AuthProvider:        getTextString(dbUser.AuthProvider),
		ProviderID:          dbUser.ProviderID,
		EmployeeType:        getTextString(dbUser.EmployeeType),
		EmploymentType:      strPtr(getTextString(dbUser.EmploymentType)),
		// Fill in missing fields with empty values
		Address:      getTextString(dbUser.UserAddress),
		City:         getTextString(dbUser.UserCity),
		PostalCode:   getTextString(dbUser.UserPostalCode),
		WebAuthToken: "",
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}
}

func mapDBCompanyInfoToDomainCompanyInfo(dbCompanyInfo db.Users) *domain.CompanyInfo {
	return &domain.CompanyInfo{
		CompanyName:         strPtr(getTextString((dbCompanyInfo.CompanyName))),
		CompanyWebsite:      strPtr(getTextString((dbCompanyInfo.CompanyWebsite))),
		CompanySize:         strPtr(getTextString((dbCompanyInfo.CompanySize))),
		CompanyIndustry:     strPtr(getTextString((dbCompanyInfo.CompanyIndustry))),
		CompanyDescription:  strPtr(getTextString((dbCompanyInfo.CompanyDescription))),
		CompanyHeadquarters: strPtr(getTextString((dbCompanyInfo.CompanyHeadquarters))),
	}
}

func (r *UserRepository) UpdateUserPersonalDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	params := db.UpdateUserPersonalDetailsParams{
		ID:                  user.ID,
		PhoneNumber:         toPgTextPtr(user.PhoneNumber),
		Nationality:         user.Nationality,
		ResidentialCountry:  toPgTextPtr(user.ResidentialCountry),
		AccountType:         user.AccountType,
		PersonalAccountType: user.PersonalAccountType,
	}

	dbUser, err := r.store.UpdateUserPersonalDetails(ctx, params)
	if err != nil {
		return nil, err
	}

	return mapDBUserToDomainUser(dbUser), nil
}
func (r *UserRepository) UpdateUserBusinessDetails(ctx context.Context, companyInfo domain.CompanyInfo) (*domain.CompanyInfo, error) {
	params := db.UpdateUserCompanyDetailsParams{
		ID:                  companyInfo.UserID,
		CompanyName:         toPgText(*companyInfo.CompanyName),
		CompanySize:         toPgText(*companyInfo.CompanySize),
		CompanyIndustry:     toPgText(*companyInfo.CompanyIndustry),
		CompanyDescription:  toPgText(*companyInfo.CompanyDescription),
		CompanyHeadquarters: toPgText(*companyInfo.CompanyHeadquarters),
		AccountType:         companyInfo.AccountType,
	}

	dbUser, err := r.store.UpdateUserCompanyDetails(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapDBCompanyInfoToDomainCompanyInfo(dbUser), nil
}

func (r *UserRepository) UpdateUserAddressDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	params := db.UpdateUserAddressParams{
		ID:             user.ID,
		UserAddress:    toPgTextPtr(user.UserAddress),
		UserCity:       toPgTextPtr(user.UserCity),
		UserPostalCode: toPgTextPtr(user.UserPostalCode),
	}
	dbUser, err := r.store.UpdateUserAddress(ctx, params)
	if err != nil {
		return nil, err
	}
	return mapDBUserToDomainUser(dbUser), nil
}

func toPgText(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func toPgTextPtr(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

// DeactivateUser marks a user as inactive
func (r *UserRepository) DeactivateUser(ctx context.Context, id uuid.UUID) error {
	// // Use UpdateUserIsActiveStatus method from the store
	// params := db.UpdateUserIsActiveStatusParams{
	//     ID:      id,
	//     IsActive: false,
	// }

	// return r.store.UpdateUserIsActiveStatus(ctx, params)
	panic("implementation")
}

// DeleteUser removes a user from the database
func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.store.DeleteUser(ctx, id)
}

// SetMFASecret sets the MFA secret for a user
func (r *UserRepository) SetMFASecret(ctx context.Context, userID uuid.UUID, secret string) error {
	// params := db.UpdateUserMFASecretParams{
	//     ID:       userID,
	//     MfaSecret: pgtype.Text{String: secret, Valid: true},
	// }

	// return r.store.UpdateUserMFASecret(ctx, params)
	panic("implementation")
}

// GetMFASecret retrieves the MFA secret for a user
func (r *UserRepository) GetMFASecret(ctx context.Context, userID uuid.UUID) (string, error) {
	// // Get the user from the database
	// dbUser, err := r.store.GetUser(ctx, userID)
	// if err != nil {
	//     return "", err
	// }

	// // Check if MFA secret exists
	// if !dbUser.MfaSecret.Valid {
	//     return "", nil
	// }

	// return dbUser.MfaSecret.String, nil
	panic("implementation")
}

// Fix the UpdatePassword method to match the interface signature
// UpdatePassword updates a user's password after verifying the old password
// func (r *UserRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
// // First get the user to verify the old password
// dbUser, err := r.store.GetUser(ctx, userID)
// if err != nil {
//     return err
// }

// // Here you would typically verify that oldPassword matches the stored hash
// // This would require a password verification function that's not shown in the code
// // For example:
// // if !verifyPassword(oldPassword, dbUser.PasswordHash.String) {
// //     return errors.New("old password does not match")
// // }

// // Update with the new password hash
// params := db.UpdateUserPasswordParams{
//     ID:           userID,
//     PasswordHash: pgtype.Text{String: newPassword, Valid: true},
// }

// return r.store.UpdateUserPassword(ctx, params)

// 		panic("implementation")

// }
