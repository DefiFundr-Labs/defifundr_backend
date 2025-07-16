// services/auth_service.go (improved)
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/demola234/defifundr/pkg/tracing"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	commons "github.com/demola234/defifundr/infrastructure/hash"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	random "github.com/demola234/defifundr/pkg/random"
	tokenMaker "github.com/demola234/defifundr/pkg/token_maker"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
)

type authService struct {
	userRepo     ports.UserRepository
	sessionRepo  ports.SessionRepository
	oauthRepo    ports.OAuthRepository
	walletRepo   ports.WalletRepository
	securityRepo ports.SecurityRepository
	emailService ports.EmailService
	companyRepo ports.CompanyRepository
	tokenMaker   tokenMaker.Maker
	config       config.Config
	logger       logging.Logger
	otpRepo      ports.OTPRepository
	userService  ports.UserService 
}

// SetupMFA sets up multi-factor authentication for a user
func (a *authService) SetupMFA(ctx context.Context, userID uuid.UUID) (string, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "SetupMFA")
	defer span.End()
	// Check if user exists
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Generate a new TOTP secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "DefiFundr",
		AccountName: user.Email,
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Store the secret in the database
	err = a.userRepo.SetMFASecret(ctx, userID, key.Secret())
	if err != nil {
		return "", fmt.Errorf("failed to store MFA secret: %w", err)
	}

	// Log the MFA setup
	a.LogSecurityEvent(ctx, "mfa_setup_initiated", userID, map[string]interface{}{
		"time": time.Now().Format(time.RFC3339),
	})

	// Return the TOTP URI for QR code generation
	return key.URL(), nil
}

// VerifyMFA verifies a TOTP code
func (a *authService) VerifyMFA(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "VerifyMFA")
	defer span.End()
	// Get the MFA secret for the user
	secret, err := a.userRepo.GetMFASecret(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get MFA secret: %w", err)
	}

	// Validate the TOTP code
	valid := totp.Validate(code, secret)

	// Log the verification attempt
	a.LogSecurityEvent(ctx, "mfa_verification", userID, map[string]interface{}{
		"success": valid,
		"time":    time.Now().Format(time.RFC3339),
	})

	return valid, nil
}

// NewAuthService creates a new instance of authService
func NewAuthService(
	userRepo ports.UserRepository,
	sessionRepo ports.SessionRepository,
	oauthRepo ports.OAuthRepository,
	walletRepo ports.WalletRepository,
	securityRepo ports.SecurityRepository,
	emailService ports.EmailService,
	companyRepo ports.CompanyRepository,
	tokenMaker tokenMaker.Maker,
	config config.Config,
	logger logging.Logger,
	otpRepo ports.OTPRepository,
	userService ports.UserService,
) ports.AuthService {
	return &authService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		oauthRepo:    oauthRepo,
		walletRepo:   walletRepo,
		securityRepo: securityRepo,
		emailService: emailService,
		companyRepo: companyRepo,
		tokenMaker:   tokenMaker,
		config:       config,
		logger:       logger,
		otpRepo:      otpRepo,
		userService:  userService,
	}
}

// Login implements ports.AuthService.
func (a *authService) Login(ctx context.Context, email string,  provider string, providerId string, webAuthToken string, password string) (*domain.User, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "Login")
	defer span.End()
	a.logger.Info("Starting user registration process", map[string]interface{}{
		"email":    email,
		"provider": provider,
	})

	

	if provider == "email" {
		// Email-based authentication requires a password
		if password == "" {
			a.logger.Error("Password required for email authentication", nil, map[string]interface{}{
				"email": email,
			})
			pwErr := errors.New("password is required for email authentication")
			span.RecordError(pwErr)
			return nil, pwErr
		}

		// Check if the user exists
		existingUser, err := a.userRepo.GetUserByEmail(ctx, email)
		if err != nil {
			a.logger.Error("Failed to get user by email", err, map[string]interface{}{
				"email": email,
			})
			span.RecordError(err)
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}
		if existingUser == nil {
			a.logger.Error("User not found", nil, map[string]interface{}{
				"email": email,
			})
			notFoundErr := errors.New("user not found")
			span.RecordError(notFoundErr)
			return nil, notFoundErr
		}

		// Verify the password
		checkedPassword, err := commons.CheckPassword(password, existingUser.PasswordHash)
		if err != nil {
			a.logger.Error("Validation Failed", err, map[string]interface{}{
				"email": email,
			})
			return nil, fmt.Errorf("failed to check password: %w", err)
		}

		if !checkedPassword {
			a.logger.Error("Invalid password", nil, map[string]interface{}{
				"email": email,
			})
			invPwErr := errors.New("invalid password")
			span.RecordError(invPwErr)
			return nil, invPwErr
		}
	} else if provider != "" && webAuthToken != "" {
		// For OAuth or Web3Auth, validate the token and fill user data
		claims, err := a.oauthRepo.ValidateWebAuthToken(ctx, webAuthToken)
		if err != nil {
			a.logger.Error("Failed to validate WebAuth token", err, map[string]interface{}{
				"provider": provider,
			})
			return nil, fmt.Errorf("invalid authentication token: %w", err)
		}

		// Extract user information from OAuth claims
		if claims.Email != "" {
			email = claims.Email
		}
	} else {
		a.logger.Error("Missing authentication credentials", nil, map[string]interface{}{
			"provider": provider,
		})
		credErr := errors.New("missing authentication credentials")
		span.RecordError(credErr)
		return nil, credErr
	}
	// Step 2: Check if user with same email already exists
	existingUser, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		a.logger.Error("Failed to get user by email", err, map[string]interface{}{
			"email": email,
		})
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if existingUser == nil {
		a.logger.Warn("Login attempt for non-existing email", map[string]interface{}{
			"email": email,
		})
		notRegisteredErr := errors.New("email not registered")
		span.RecordError(notRegisteredErr)
		return nil, notRegisteredErr
	}

	// Return the user
	return existingUser, nil
}

// RegisterUser implements the user registration process with Web3Auth integration
func (a *authService) RegisterUser(ctx context.Context, email, firstName, lastName, authProvider, webAuthToken, passwordStr string) (*domain.User, error) {
	user := domain.User{
		ID:           uuid.New(),
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		AuthProvider: authProvider,
		AccountType:  "personal",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "RegisterUser")
	defer span.End()
	a.logger.Info("Starting user registration process", map[string]interface{}{
		"email":    user.Email,
		"provider": user.AuthProvider,
	})

	existingUser, err := a.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		a.logger.Warn("Registration attempt for existing email", map[string]interface{}{
			"email": user.Email,
		})
		span.RecordError(err)
		return nil, errors.New("email already registered")
	}

	// Step 1: Handle authentication based on provider
	if user.AuthProvider == "email" {
		// Email-based authentication requires a password
		if passwordStr == "" {
			a.logger.Error("Password required for email authentication", nil, map[string]interface{}{
				"email": user.Email,
			})
			pwErr := errors.New("password is required for email authentication")
			span.RecordError(pwErr)
			return nil, pwErr
		}

		// Hash the password
		hashedPassword, err := commons.HashPassword(passwordStr)
		if err != nil {
			a.logger.Error("Failed to hash password", err, map[string]interface{}{
				"email": user.Email,
			})
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = hashedPassword
	} else if user.AuthProvider != "" && webAuthToken != "" {
		// For OAuth or Web3Auth, validate the token and fill user data
		claims, err := a.oauthRepo.ValidateWebAuthToken(ctx, webAuthToken)
		if err != nil {
			a.logger.Error("Failed to validate WebAuth token", err, map[string]interface{}{
				"provider": user.AuthProvider,
			})
			span.RecordError(err)
			return nil, fmt.Errorf("invalid authentication token: %w", err)
		}

		// Extract user information from OAuth claims
		// Now claims is a *Web3AuthClaims struct, not a map
		if claims.Email != "" {
			user.Email = claims.Email
		}

		if claims.Name != "" {
			nameParts := strings.Split(claims.Name, " ")
			user.FirstName = nameParts[0]
			if len(nameParts) > 1 {
				user.LastName = strings.Join(nameParts[1:], " ")
			}
		}

		if claims.ProfileImage != "" {
			profileImage := claims.ProfileImage
			user.ProfilePictureURL = &profileImage
		}

		// Set provider ID (usually the email for Google OAuth)
		if claims.VerifierID != "" {
			user.ProviderID = &claims.VerifierID
		}

		// Refine provider information based on verifier
		if claims.Verifier != "" {
			if strings.Contains(claims.Verifier, "google") {
				user.AuthProvider = "google"
			} else if strings.Contains(claims.Verifier, "facebook") {
				user.AuthProvider = "facebook"
			} else if strings.Contains(claims.Verifier, "twitter") {
				user.AuthProvider = "twitter"
			}
			// Add more provider mappings as needed
		}

		// For OAuth users, no password is needed
		user.PasswordHash = ""
	} else {
		// If not email auth and no token provided, it's an error
		a.logger.Error("Missing authentication credentials", nil, map[string]interface{}{
			"provider": user.AuthProvider,
		})
		return nil, errors.New("missing authentication credentials")
	}


	// Step 3: Set default values if not provided
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	// Step 4: Register the user in the database
	createdUser, err := a.userRepo.CreateUser(ctx, user)
	if err != nil {
		a.logger.Error("Failed to register user", err, map[string]interface{}{
			"email": user.Email,
		})
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return createdUser, nil
}

// RegisterBusiness implements ports.AuthService.
func (a *authService) RegisterBusiness(ctx context.Context, companyInfo domain.Company) (*domain.Company, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "RegisterBusiness")
	defer span.End()
	
	a.logger.Info("Starting business registration process", map[string]interface{}{
		"owner_id": companyInfo.OwnerID,
	})

	// Check if company already exists for this owner
	existingCompany, err := a.companyRepo.GetCompanyByOwnerID(ctx, companyInfo.OwnerID)
	if err != nil {
		// If company doesn't exist, create a new one
		if err.Error() == "no rows in result set" || err.Error() == "sql: no rows in result set" {
			// Verify that the owner exists
			owner, err := a.userRepo.GetUserByID(ctx, companyInfo.OwnerID)
			if err != nil {
				a.logger.Error("Failed to get owner by ID", err, map[string]interface{}{
					"owner_id": companyInfo.OwnerID,
				})
				span.RecordError(err)
				return nil, fmt.Errorf("failed to get owner by ID: %w", err)
			}

			// Set additional fields for new company
			companyInfo.ID = uuid.New()
			companyInfo.CreatedAt = time.Now()
			companyInfo.UpdatedAt = time.Now()
			companyInfo.KYBStatus = "pending" // Default KYB status

			// Create the company
			createdCompany, err := a.companyRepo.CreateCompany(ctx, companyInfo)
			if err != nil {
				a.logger.Error("Failed to create company", err, map[string]interface{}{
					"owner_id": companyInfo.OwnerID,
				})
				span.RecordError(err)
				return nil, fmt.Errorf("failed to create company: %w", err)
			}

			// Update the user's business details timestamp
			owner.UpdatedAt = time.Now()
			_, err = a.userRepo.UpdateUserBusinessDetails(ctx, *owner)
			if err != nil {
				a.logger.Error("Failed to update user business details timestamp", err, map[string]interface{}{
					"user_id": companyInfo.OwnerID,
				})
				// Don't fail the whole operation for this
			}

			a.logger.Info("Company created successfully", map[string]interface{}{
				"company_id": createdCompany.ID,
				"owner_id":   companyInfo.OwnerID,
			})

			return createdCompany, nil

		} else {
			a.logger.Error("Failed to get company by owner ID", err, map[string]interface{}{
				"owner_id": companyInfo.OwnerID,
			})
			span.RecordError(err)
			return nil, fmt.Errorf("failed to get company by owner ID: %w", err)
		}
	}

	// Company exists, update it
	updatedCompany := *existingCompany
	
	// Update company fields if provided
	if companyInfo.CompanyName != "" {
		updatedCompany.CompanyName = companyInfo.CompanyName
	}
	if companyInfo.CompanyDescription != nil {
		updatedCompany.CompanyDescription = companyInfo.CompanyDescription
	}
	if companyInfo.CompanySize != nil {
		updatedCompany.CompanySize = companyInfo.CompanySize
	}
	if companyInfo.CompanyHeadquarters != nil {
		updatedCompany.CompanyHeadquarters = companyInfo.CompanyHeadquarters
	}
	if companyInfo.CompanyIndustry != nil {
		updatedCompany.CompanyIndustry = companyInfo.CompanyIndustry
	}
	if companyInfo.CompanyWebsite != nil {
		updatedCompany.CompanyWebsite = companyInfo.CompanyWebsite
	}
	
	updatedCompany.UpdatedAt = time.Now()

	// Update the company in the database
	company, err := a.companyRepo.UpdateCompany(ctx, updatedCompany)
	if err != nil {
		a.logger.Error("Failed to update company", err, map[string]interface{}{
			"company_id": existingCompany.ID,
			"owner_id":   companyInfo.OwnerID,
		})
		span.RecordError(err)
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	// Update the user's business details timestamp
	owner, err := a.userRepo.GetUserByID(ctx, companyInfo.OwnerID)
	if err == nil {
		owner.UpdatedAt = time.Now()
		_, err = a.userRepo.UpdateUserBusinessDetails(ctx, *owner)
		if err != nil {
			a.logger.Error("Failed to update user business details timestamp", err, map[string]interface{}{
				"user_id": companyInfo.OwnerID,
			})
			// Don't fail the whole operation for this
		}
	}

	a.logger.Info("Company updated successfully", map[string]interface{}{
		"company_id": company.ID,
		"owner_id":   companyInfo.OwnerID,
	})

	return company, nil
}

// RegisterPersonalDetails implements ports.AuthService
func (a *authService) RegisterPersonalDetails(ctx context.Context, userId uuid.UUID, nationality string, dateOfBirth time.Time, gender string, personalAccountType string, phoneNumber string) (*domain.User, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "RegisterPersonalDetails")
	defer span.End()
	
	a.logger.Info("Starting user personal details update process", map[string]interface{}{
		"user_id": userId,
	})

	// Get the existing user by ID
	existingUser, err := a.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		a.logger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": userId,
		})
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Update the main user record with phone number if provided
	if phoneNumber != "" {
		updatedUser := *existingUser
		updatedUser.PhoneNumber = phoneNumber
		updatedUser.UpdatedAt = time.Now()

		result, err := a.userRepo.UpdateUserPersonalDetails(ctx, updatedUser)
		if err != nil {
			a.logger.Error("Failed to update user phone number", err, map[string]interface{}{
				"user_id": userId,
			})
			span.RecordError(err)
			return nil, fmt.Errorf("failed to update user phone number: %w", err)
		}
		existingUser = result
	}

	// Check if personal user already exists
	existingPersonalUser, err := a.userRepo.GetPersonalUserByUserID(ctx, userId)
	if err != nil {
		// If personal user doesn't exist, create a new one
		if err.Error() == "no rows in result set" || err.Error() == "sql: no rows in result set" {
			personalUser := domain.PersonalUser{
				ID:                  uuid.New(), // Generate new ID for personal user
				UserID:              userId,
				Nationality:         &nationality,
				Gender:              &gender,
				DateOfBirth:         &dateOfBirth,
				PersonalAccountType: &personalAccountType,
				CreatedAt:           time.Now(),
				UpdatedAt:           time.Now(),
			}

			_, err := a.userRepo.CreatePersonalUser(ctx, personalUser)
			if err != nil {
				a.logger.Error("Failed to create personal user", err, map[string]interface{}{
					"user_id": userId,
				})
				span.RecordError(err)
				return nil, fmt.Errorf("failed to create personal user: %w", err)
			}

			a.logger.Info("Personal user created successfully", map[string]interface{}{
				"user_id": userId,
			})

		} else {
			a.logger.Error("Failed to get personal user by user ID", err, map[string]interface{}{
				"user_id": userId,
			})
			span.RecordError(err)
			return nil, fmt.Errorf("failed to get personal user: %w", err)
		}
	} else {
		// Update existing personal user
		updatedPersonalUser := *existingPersonalUser
		
		// Update fields if provided
		if nationality != "" {
			updatedPersonalUser.Nationality = &nationality
		}
		if gender != "" {
			updatedPersonalUser.Gender = &gender
		}
		if !dateOfBirth.IsZero() {
			updatedPersonalUser.DateOfBirth = &dateOfBirth
		}
		if personalAccountType != "" {
			updatedPersonalUser.PersonalAccountType = &personalAccountType
		}
		
		updatedPersonalUser.UpdatedAt = time.Now()

		_, err := a.userRepo.UpdatePersonalUser(ctx, updatedPersonalUser)
		if err != nil {
			a.logger.Error("Failed to update personal user", err, map[string]interface{}{
				"user_id": userId,
			})
			span.RecordError(err)
			return nil, fmt.Errorf("failed to update personal user: %w", err)
		}

		a.logger.Info("Personal user updated successfully", map[string]interface{}{
			"user_id": userId,
		})
	}

	// Return the updated main user
	finalUser, err := a.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		a.logger.Error("Failed to get updated user", err, map[string]interface{}{
			"user_id": userId,
		})
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	a.logger.Info("User personal details processed successfully", map[string]interface{}{
		"user_id": userId,
	})

	return finalUser, nil
}

// RegisterBusinessDetails implements ports.AuthService
func (a *authService) RegisterBusinessDetails(ctx context.Context, companyInfo domain.Company) (*domain.Company, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "RegisterBusinessDetails")
	defer span.End()
	
	a.logger.Info("Starting business details update process", map[string]interface{}{
		"owner_id": companyInfo.OwnerID,
	})

	// Get the existing company by owner ID
	existingCompany, err := a.companyRepo.GetCompanyByOwnerID(ctx, companyInfo.OwnerID)
	if err != nil {
		a.logger.Error("Failed to get company by owner ID", err, map[string]interface{}{
			"owner_id": companyInfo.OwnerID,
		})
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get company by owner ID: %w", err)
	}

	// Update only the business details fields, keeping other fields as they are
	updatedCompany := *existingCompany

	// Update company name (assuming it's always required)
	if companyInfo.CompanyName != "" {
		updatedCompany.CompanyName = companyInfo.CompanyName
	}

	// Update optional fields - check if they're pointers or regular strings
	// If CompanyDescription is a pointer
	if companyInfo.CompanyDescription != nil && *companyInfo.CompanyDescription != "" {
		updatedCompany.CompanyDescription = companyInfo.CompanyDescription
	}

	// If CompanyHeadquarters is a pointer
	if companyInfo.CompanyHeadquarters != nil && *companyInfo.CompanyHeadquarters != "" {
		updatedCompany.CompanyHeadquarters = companyInfo.CompanyHeadquarters
	}

	// If CompanyIndustry is a pointer
	if companyInfo.CompanyIndustry != nil && *companyInfo.CompanyIndustry != "" {
		updatedCompany.CompanyIndustry = companyInfo.CompanyIndustry
	}

	// If CompanySize is a pointer
	if companyInfo.CompanySize != nil && *companyInfo.CompanySize != "" {
		updatedCompany.CompanySize = companyInfo.CompanySize
	}

	// If CompanyWebsite is provided
	if companyInfo.CompanyWebsite != nil && *companyInfo.CompanyWebsite != "" {
		updatedCompany.CompanyWebsite = companyInfo.CompanyWebsite
	}

	// Set updated timestamp
	updatedCompany.UpdatedAt = time.Now()

	// Update the company in the database using company repository
	businessResult, err := a.companyRepo.UpdateCompany(ctx, updatedCompany)
	if err != nil {
		a.logger.Error("Failed to update company details", err, map[string]interface{}{
			"owner_id":   companyInfo.OwnerID,
			"company_id": existingCompany.ID,
		})
		span.RecordError(err)
		return nil, fmt.Errorf("failed to update company details: %w", err)
	}

	// Update the user's business details timestamp (optional)
	owner, err := a.userRepo.GetUserByID(ctx, companyInfo.OwnerID)
	if err == nil {
		owner.UpdatedAt = time.Now()
		_, err = a.userRepo.UpdateUserBusinessDetails(ctx, *owner)
		if err != nil {
			a.logger.Error("Failed to update user business details timestamp", err, map[string]interface{}{
				"user_id": companyInfo.OwnerID,
			})
			// Don't fail the whole operation for this
		}
	}

	a.logger.Info("Business details updated successfully", map[string]interface{}{
		"owner_id":   companyInfo.OwnerID,
		"company_id": businessResult.ID,
	})

	return businessResult, nil
}

// RegisterAddressDetails implements ports.AuthService
func (a *authService) RegisterAddressDetails(ctx context.Context, userId uuid.UUID, userAddress, userCity, userPostalCode, residentialCountry string) (*domain.User, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "RegisterAddressDetails")
	defer span.End()
	
	a.logger.Info("Starting address details update process", map[string]interface{}{
		"user_id": userId,
	})

	// Get the existing user by ID
	existingUser, err := a.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		a.logger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": userId,
		})
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Check if personal user already exists
	existingPersonalUser, err := a.userRepo.GetPersonalUserByUserID(ctx, userId)
	if err != nil {
		// If personal user doesn't exist, create a new one with address details
		if err.Error() == "no rows in result set" || err.Error() == "sql: no rows in result set" {
			personalUser := domain.PersonalUser{
				ID:                 uuid.New(),
				UserID:             userId,
				UserAddress:        &userAddress,
				UserCity:           &userCity,
				UserPostalCode:     &userPostalCode,
				ResidentialCountry: &residentialCountry,
				KYCStatus:          "pending", // Default KYC status
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}

			// Only set non-empty values
			if userAddress == "" {
				personalUser.UserAddress = nil
			}
			if userCity == "" {
				personalUser.UserCity = nil
			}
			if userPostalCode == "" {
				personalUser.UserPostalCode = nil
			}
			if residentialCountry == "" {
				personalUser.ResidentialCountry = nil
			}

			_, err := a.userRepo.CreatePersonalUser(ctx, personalUser)
			if err != nil {
				a.logger.Error("Failed to create personal user with address details", err, map[string]interface{}{
					"user_id": userId,
				})
				span.RecordError(err)
				return nil, fmt.Errorf("failed to create personal user with address details: %w", err)
			}

			a.logger.Info("Personal user created with address details", map[string]interface{}{
				"user_id": userId,
			})

		} else {
			a.logger.Error("Failed to get personal user by user ID", err, map[string]interface{}{
				"user_id": userId,
			})
			span.RecordError(err)
			return nil, fmt.Errorf("failed to get personal user: %w", err)
		}
	} else {
		// Update existing personal user with address details
		updatedPersonalUser := *existingPersonalUser
		
		// Update address fields if provided
		if userAddress != "" {
			updatedPersonalUser.UserAddress = &userAddress
		}
		if userCity != "" {
			updatedPersonalUser.UserCity = &userCity
		}
		if userPostalCode != "" {
			updatedPersonalUser.UserPostalCode = &userPostalCode
		}
		if residentialCountry != "" {
			updatedPersonalUser.ResidentialCountry = &residentialCountry
		}
		
		updatedPersonalUser.UpdatedAt = time.Now()

		_, err := a.userRepo.UpdatePersonalUser(ctx, updatedPersonalUser)
		if err != nil {
			a.logger.Error("Failed to update personal user address details", err, map[string]interface{}{
				"user_id": userId,
			})
			span.RecordError(err)
			return nil, fmt.Errorf("failed to update personal user address details: %w", err)
		}

		a.logger.Info("Personal user address details updated successfully", map[string]interface{}{
			"user_id": userId,
		})
	}

	// Update the main user record's updated_at timestamp
	result, err := a.userRepo.UpdateUserAddressDetails(ctx, *existingUser)
	if err != nil {
		a.logger.Error("Failed to update user address timestamp", err, map[string]interface{}{
			"user_id": userId,
		})
		span.RecordError(err)
		return nil, fmt.Errorf("failed to update user address timestamp: %w", err)
	}

	a.logger.Info("Address details updated successfully", map[string]interface{}{
		"user_id": userId,
	})

	return result, nil
}

// GetUserByEmail implements ports.AuthService.
func (a *authService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "GetUserByEmail")
	defer span.End()
	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		span.RecordError(err)
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByID implements ports.AuthService.
func (a *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "GetUserByID")
	defer span.End()
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if user == nil {
		span.RecordError(err)
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (a *authService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "CheckEmailExists")
	defer span.End()
	exists, err := a.userRepo.CheckEmailExists(ctx, email)
	if err != nil {
		span.RecordError(err)
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}

	return exists, nil
}

// AuthenticateWithWeb3 implements unified Web3Auth authentication flow
func (a *authService) AuthenticateWithWeb3(ctx context.Context, webAuthToken string, userAgent string, clientIP string) (*domain.User, *domain.Session, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "AuthenticateWithWeb3")
	defer span.End()
	
	// Validate the Web3Auth token
	claims, err := a.oauthRepo.ValidateWebAuthToken(ctx, webAuthToken)
	if err != nil {
		a.logger.Error("Web3Auth token validation failed", err, nil)
		span.RecordError(err)
		return nil, nil, err
	}

	// Extract identity information
	email := claims.Email
	if email == "" {
		err := errors.New("email not provided in Web3Auth token")
		span.RecordError(err)
		return nil, nil, err
	}

	// Check if user exists
	existingUser, err := a.userRepo.GetUserByEmail(ctx, email)
	var user *domain.User
	isNewUser := false

	if err != nil || existingUser == nil {
		// This is a new user - create account
		a.logger.Info("Creating new user from Web3Auth", map[string]interface{}{
			"email":    email,
			"verifier": claims.Verifier,
		})

		// Extract profile data from claims
		firstName, lastName := extractNameFromClaims(claims)

		// Determine provider from verifier
		authProvider := mapVerifierToProvider(claims.Verifier)

		// Create new user
		newUser := domain.User{
			ID:                uuid.New(),
			Email:             email,
			FirstName:         firstName,
			LastName:          lastName,
			ProfilePictureURL: &claims.ProfileImage, // Use the correct field name
			ProviderID:        &claims.VerifierID,   // This should be a pointer
			AuthProvider:      string(authProvider),
			AccountType:       "personal", // Default value, can be updated later
			EmailVerified:     true,       // Web3Auth emails are typically verified
			AccountStatus:     "active",   // Set default account status
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		user, err = a.userRepo.CreateUser(ctx, newUser)
		if err != nil {
			span.RecordError(err)
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}

		isNewUser = true

		// Create personal user profile with additional details
		personalUser := domain.PersonalUser{
			ID:                  uuid.New(),
			UserID:              user.ID,
			PersonalAccountType: stringPtr("user"), // Default value, can be updated later
			Nationality:         stringPtr("unknown"), // Default value, can be updated later
			KYCStatus:           "pending",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		_, err = a.userRepo.CreatePersonalUser(ctx, personalUser)
		if err != nil {
			// Log error but don't fail the authentication
			a.logger.Error("Failed to create personal user profile", err, map[string]interface{}{
				"user_id": user.ID,
			})
		}

		// Track registration event
		a.LogSecurityEvent(ctx, "user_registered", user.ID, map[string]interface{}{
			"provider": user.AuthProvider,
			"email":    user.Email,
		})
	} else {
		// Existing user - return user data
		user = existingUser

		// Update any profile info that may have changed
		updateNeeded := false

		// Update profile picture if new one available
		if claims.ProfileImage != "" && (user.ProfilePictureURL == nil || *user.ProfilePictureURL != claims.ProfileImage) {
			user.ProfilePictureURL = &claims.ProfileImage
			updateNeeded = true
		}

		// Update name if it was empty before
		if user.FirstName == "" && user.LastName == "" {
			firstName, lastName := extractNameFromClaims(claims)
			user.FirstName = firstName
			user.LastName = lastName
			updateNeeded = true
		}

		if updateNeeded {
			user.UpdatedAt = time.Now()
			_, err := a.userRepo.UpdateUser(ctx, *user)
			if err != nil {
				a.logger.Error("Failed to update user profile", err, map[string]interface{}{
					"user_id": user.ID,
				})
			}
		}
	}

	// Process wallets from Web3Auth claims if available
	if len(claims.Wallets) > 0 {
		for _, wallet := range claims.Wallets {
			err := a.processWallet(ctx, user.ID, wallet)
			if err != nil {
				// Log but continue - wallet linking is non-critical
				a.logger.Warn("Failed to process wallet", map[string]interface{}{
					"user_id": user.ID,
					"wallet":  wallet.PublicKey,
					"error":   err.Error(),
				})
			}
		}
	}

	// Create session for the user
	session, err := a.CreateSession(ctx, user.ID, userAgent, clientIP, webAuthToken, user.Email, "web3auth")
	if err != nil {
		span.RecordError(err)
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Track login event
	a.LogSecurityEvent(ctx, "user_login", user.ID, map[string]interface{}{
		"provider":    user.AuthProvider,
		"ip":          clientIP,
		"is_new_user": isNewUser,
	})

	// Check for suspicious activity in background
	go a.detectSuspiciousActivity(context.Background(), user.ID, clientIP, userAgent)

	return user, session, nil
}

// Helper function to create string pointer
func stringPtr(s string) *string {
	return &s
}

// processWallet handles wallet data from Web3Auth
func (a *authService) processWallet(ctx context.Context, userID uuid.UUID, wallet domain.Wallet) error {
	// Skip empty wallets
	if wallet.PublicKey == "" {
		return nil
	}

	// Check if wallet already exists
	existingWallet, err := a.walletRepo.GetWalletByAddress(ctx, wallet.PublicKey)
	if err != nil {
		return fmt.Errorf("error checking wallet existence: %w", err)
	}

	// If wallet exists and belongs to user, nothing to do
	if existingWallet != nil && existingWallet.UserID == userID {
		return nil
	}

	// If wallet exists but belongs to another user, log security event and don't link
	if existingWallet != nil && existingWallet.UserID != userID {
		a.LogSecurityEvent(ctx, "wallet_conflict", userID, map[string]interface{}{
			"wallet_address": wallet.PublicKey,
			"existing_user":  existingWallet.UserID.String(),
		})
		return fmt.Errorf("wallet already linked to another account")
	}

	// Determine chain from wallet type
	chain := "ethereum"
	if wallet.Type != "" && wallet.Type != "hex" {
		chain = strings.ToLower(wallet.Type)
	}

	// Create new wallet
	return a.LinkWallet(ctx, userID, wallet.PublicKey, wallet.Type, chain)
}

// LinkWallet links a blockchain wallet to a user account (continued)
func (a *authService) LinkWallet(ctx context.Context, userID uuid.UUID, walletAddress string, walletType string, chain string) error {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "LinkWallet")
	defer span.End()
	// ... (previous code)

	// Log security event
	a.LogSecurityEvent(ctx, "wallet_linked", userID, map[string]interface{}{
		"wallet_address": walletAddress,
		"wallet_type":    walletType,
		"chain":          chain,
	})

	return nil
}

// GetUserWallets retrieves all wallets for a user
func (a *authService) GetUserWallets(ctx context.Context, userID uuid.UUID) ([]domain.UserWallet, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "GetUserWallets")
	defer span.End()
	return a.walletRepo.GetWalletsByUserID(ctx, userID)
}

// GetProfileCompletionStatus calculates profile completion percentage
func (a *authService) GetProfileCompletionStatus(ctx context.Context, userID uuid.UUID) (*domain.ProfileCompletion, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "GetProfileCompletionStatus")
	defer span.End()

	// Get user data
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Get personal user data (may not exist)
	personalUser, err := a.userRepo.GetPersonalUserByUserID(ctx, userID)
	if err != nil {
		// If personal user doesn't exist, we'll treat all personal fields as incomplete
		a.logger.Info("Personal user not found, treating personal fields as incomplete", map[string]interface{}{
			"user_id": userID,
			"error":   err.Error(),
		})
		personalUser = nil
	}

	// Define required fields
	type fieldCheck struct {
		name     string
		required bool
		value    interface{}
	}

	// Common required fields from main user table
	fields := []fieldCheck{
		{"First Name", true, user.FirstName != ""},
		{"Last Name", true, user.LastName != ""},
		{"Phone Number", true, user.PhoneNumber != ""},
		{"Email Verified", true, user.EmailVerified},
	}

	// Personal user fields (if personal user exists)
	if personalUser != nil {
		fields = append(fields, []fieldCheck{
			{"Nationality", true, personalUser.Nationality != nil && *personalUser.Nationality != "" && *personalUser.Nationality != "unknown"},
			{"Address", true, personalUser.UserAddress != nil && *personalUser.UserAddress != ""},
			{"City", true, personalUser.UserCity != nil && *personalUser.UserCity != ""},
			{"Postal Code", true, personalUser.UserPostalCode != nil && *personalUser.UserPostalCode != ""},
			{"Gender", false, personalUser.Gender != nil && *personalUser.Gender != ""}, // Optional field
			{"Date of Birth", false, personalUser.DateOfBirth != nil}, // Optional field
		}...)
	} else {
		// If no personal user, add these as incomplete required fields
		fields = append(fields, []fieldCheck{
			{"Nationality", true, false},
			{"Address", true, false},
			{"City", true, false},
			{"Postal Code", true, false},
			{"Gender", false, false},
			{"Date of Birth", false, false},
		}...)
	}

	// Account type specific additional fields
	switch user.AccountType {
	case "business", "company":
		// Add business-specific fields if needed
		// fields = append(fields, []fieldCheck{
		//     {"Business Name", true, /* check business fields */},
		// }...)
	case "personal":
		// Personal account might need additional fields
		if personalUser != nil {
			fields = append(fields, []fieldCheck{
				{"Employment Type", false, personalUser.EmploymentType != nil && *personalUser.EmploymentType != ""},
				{"Job Role", false, personalUser.JobRole != nil && *personalUser.JobRole != ""},
			}...)
		}
	}

	// Calculate completion percentage
	var completedFields, requiredFields int
	var missingFields []string

	for _, field := range fields {
		if field.required {
			requiredFields++

			// Check if the field has a value
			isCompleted := false
			switch v := field.value.(type) {
			case bool:
				isCompleted = v
			default:
				isCompleted = field.value != nil
			}

			if isCompleted {
				completedFields++
			} else {
				missingFields = append(missingFields, field.name)
			}
		}
	}

	// Calculate percentage
	percentage := 0
	if requiredFields > 0 {
		percentage = (completedFields * 100) / requiredFields
	}

	// Determine required actions
	var requiredActions []string

	if len(missingFields) > 0 {
		requiredActions = append(requiredActions, "complete_profile")
	}

	// Add specific actions based on missing fields
	if personalUser == nil {
		requiredActions = append(requiredActions, "create_personal_profile")
	}

	if !user.EmailVerified {
		requiredActions = append(requiredActions, "verify_email")
	}

	if user.PhoneNumber != "" && !user.PhoneNumberVerified {
		requiredActions = append(requiredActions, "verify_phone")
	}

	// Create profile completion response
	completion := &domain.ProfileCompletion{
		UserID:               userID,
		CompletionPercentage: percentage,
		MissingFields:        missingFields,
		RequiredActions:      requiredActions,
	}

	a.logger.Info("Profile completion calculated", map[string]interface{}{
		"user_id":              userID,
		"completion_percentage": percentage,
		"missing_fields":       len(missingFields),
		"required_actions":     len(requiredActions),
	})

	return completion, nil
}

// detectSuspiciousActivity monitors for suspicious login activity
func (a *authService) detectSuspiciousActivity(ctx context.Context, userID uuid.UUID, clientIP string, userAgent string) {
	// Get user's previous logins
	previousLogins, err := a.securityRepo.GetRecentLoginsByUserID(ctx, userID, 5)
	if err != nil {
		return // Don't block authentication on error
	}

	// If this is the first login, nothing to check
	if len(previousLogins) == 0 {
		return
	}

	// Check if this is a login from a new location/device
	isNewIP := true
	isNewDevice := true

	for _, login := range previousLogins {
		if login.IPAddress == clientIP {
			isNewIP = false
		}

		if login.UserAgent == userAgent {
			isNewDevice = false
		}
	}

	// If this is a new location or device, send notification
	if isNewIP || isNewDevice {
		// Get user for email notification
		user, err := a.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			a.logger.Error("Failed to get user for security notification", err, map[string]interface{}{
				"user_id": userID,
			})
			return
		}

		// Send security notification
		deviceInfo := parseUserAgent(userAgent)
		loginTime := time.Now().Format(time.RFC1123)

		// Send email alert
		if a.emailService != nil {
			emailData := map[string]interface{}{
				"name":       user.FirstName,
				"ip":         clientIP,
				"device":     deviceInfo,
				"time":       loginTime,
				"login_type": "Web3Auth",
			}

			fmt.Printf("Email Data: %s\n", emailData)

			// Use a new context to avoid cancellation
			go func() {
				bgCtx := context.Background()
				err := a.emailService.SendBatchUpdate(
					bgCtx,
					[]string{user.Email},
					"New Login Detected",
					fmt.Sprintf("We noticed a new login to your DefiFundr account from %s at %s. If this was you, you can ignore this message.", deviceInfo, loginTime),
				)
				if err != nil {
					a.logger.Error("Failed to send security notification email", err, nil)
				}
			}()
		}

		// Send Email Alert
		fmt.Printf("Security alert: New login detected from %s at %s\n", deviceInfo, loginTime)
		fmt.Printf("Login Time: %s\n", loginTime)
		fmt.Printf("Send Email Alert: %s\n", user.Email)

		securityTreat := domain.SecurityEvent{
			UserID:    userID,
			IPAddress: clientIP,
			UserAgent: userAgent,
			ID:        uuid.New(),
			EventType: "New IP/Device Detected",
			Metadata: map[string]interface{}{
				"device":     deviceInfo,
				"time":       loginTime,
				"login_type": "Web3Auth",
			},
			Timestamp: time.Now(),
		}

		// Log the suspicious activity
		a.securityRepo.LogSecurityEvent(ctx, securityTreat)
	}
}

// CreateSession creates a new session for the user
func (a *authService) CreateSession(ctx context.Context, userID uuid.UUID, userAgent, clientIP string, webOAuthClientID string, email string, loginType string) (*domain.Session, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "CreateSession")
	defer span.End()
	a.logger.Info("Creating new session", map[string]interface{}{
		"user_id": userID,
		"ip":      clientIP,
	})

	// Generate a new access token
	accessToken, payload, err := a.tokenMaker.CreateToken(
		email,
		userID,
		a.config.AccessTokenDuration,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate a refresh token
	refreshToken, payload, err := a.tokenMaker.CreateToken(
		email,
		userID,
		a.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	expiresAt := time.Now().Add(a.config.RefreshTokenDuration)

	// Create session
	session := domain.Session{
		ID:               uuid.New(),
		UserID:           userID,
		RefreshToken:     refreshToken,
		OAuthAccessToken: accessToken,
		UserAgent:        userAgent,
		ClientIP:         clientIP,
		IsBlocked:        false,
		UserLoginType:    loginType,
		ExpiresAt:        &expiresAt,
		CreatedAt:        time.Now(),
	}

	// Set Web3Auth token if provided
	if webOAuthClientID != "" {
		session.WebOAuthClientID = &webOAuthClientID
	}

	// Create session in database
	userSession, err := a.sessionRepo.CreateSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	a.logger.Info("Session created successfully", map[string]interface{}{
		"session_id": session.ID,
		"expires_at": payload.ExpiredAt,
	})

	return userSession, nil
}

// GetActiveDevices returns all active devices for a user
func (a *authService) GetActiveDevices(ctx context.Context, userID uuid.UUID) ([]domain.DeviceInfo, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "GetActiveDevices")
	defer span.End()
	// Get active sessions
	activeSessions, err := a.sessionRepo.GetActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	// Create device info list
	devices := make([]domain.DeviceInfo, 0, len(activeSessions))

	for _, session := range activeSessions {
		// Parse user agent
		deviceInfo := parseUserAgent(session.UserAgent)

		// Create device info
		device := domain.DeviceInfo{
			SessionID:       session.ID,
			Browser:         deviceInfo,
			OperatingSystem: extractOSFromUserAgent(session.UserAgent),
			DeviceType:      determineDeviceType(session.UserAgent),
			IPAddress:       session.ClientIP,
			LoginType:       session.UserLoginType,
			LastUsed:        time.Now(), // Update to use lastUsedAt when available
			CreatedAt:       session.CreatedAt,
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// RevokeSession revokes a specific session
func (a *authService) RevokeSession(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID) error {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "RevokeSession")
	defer span.End()
	// Get session to verify ownership
	session, err := a.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Verify session belongs to user
	if session.UserID != userID {
		span.RecordError(err)
		return errors.New("session does not belong to user")
	}

	// Block session
	err = a.sessionRepo.BlockSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to block session: %w", err)
	}

	// Log security event
	a.LogSecurityEvent(ctx, "session_revoked", userID, map[string]interface{}{
		"session_id": sessionID,
	})

	return nil
}

// Logout logs out a user by revoking their session
func (a *authService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "Logout")
	defer span.End()
	return a.sessionRepo.DeleteSession(ctx, sessionID)
}

// RefreshToken refreshes an access token
func (a *authService) RefreshToken(ctx context.Context, refreshToken, userAgent, clientIP string) (*domain.Session, string, error) {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "RefreshToken")
	defer span.End()
	// Get session by refresh token
	session, err := a.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	// Validate session
	if session == nil || session.IsBlocked || time.Now().After(*session.ExpiresAt) {
		span.RecordError(err)
		return nil, "", errors.New("invalid or expired refresh token")
	}

	// Get the user
	user, err := a.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	// Generate a new access token
	accessToken, _, err := a.tokenMaker.CreateToken(
		user.Email,
		user.ID,
		a.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate a refresh token
	newRefreshToken, _, err := a.tokenMaker.CreateToken(
		user.Email,
		user.ID,
		a.config.RefreshTokenDuration,
	)

	if err != nil {
		return nil, "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Update session with new refresh token
	updatedSession, err := a.sessionRepo.UpdateRefreshToken(ctx, session.ID, newRefreshToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to update refresh token: %w", err)
	}

	// Log security event
	a.LogSecurityEvent(ctx, "token_refreshed", user.ID, map[string]interface{}{
		"session_id": session.ID,
		"ip":         clientIP,
	})

	return updatedSession, accessToken, nil
}

// LogSecurityEvent logs a security event
func (a *authService) LogSecurityEvent(ctx context.Context, eventType string, userID uuid.UUID, metadata map[string]interface{}) error {
	ctx, span := tracing.Tracer("auth-service").Start(ctx, "LogSecurityEvent")
	defer span.End()
	// Get client IP from context if available
	clientIP := ""
	if ipValue := ctx.Value("client_ip"); ipValue != nil {
		if ip, ok := ipValue.(string); ok {
			clientIP = ip
		}
	}

	// Get user agent from context if available
	userAgent := ""
	if uaValue := ctx.Value("user_agent"); uaValue != nil {
		if ua, ok := uaValue.(string); ok {
			userAgent = ua
		}
	}

	// Create security event
	event := domain.SecurityEvent{
		ID:        uuid.New(),
		UserID:    userID,
		EventType: eventType,
		IPAddress: clientIP,
		UserAgent: userAgent,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}

	// Log event
	a.logger.Info("Security event", map[string]interface{}{
		"event_type": eventType,
		"user_id":    userID.String(),
		"ip":         clientIP,
		"metadata":   metadata,
	})

	// Store event in database
	return a.securityRepo.LogSecurityEvent(ctx, event)
}

// Helper functions

// extractNameFromClaims extracts first and last name from Web3Auth claims
func extractNameFromClaims(claims *domain.Web3AuthClaims) (string, string) {
	if claims.Name == "" {
		return "User", ""
	}

	nameParts := strings.Split(claims.Name, " ")
	firstName := nameParts[0]

	var lastName string
	if len(nameParts) > 1 {
		lastName = strings.Join(nameParts[1:], " ")
	}

	return firstName, lastName
}

// mapVerifierToProvider maps Web3Auth verifier to auth provider
func mapVerifierToProvider(verifier string) domain.AuthProvider {
	lowerVerifier := strings.ToLower(verifier)

	if strings.Contains(lowerVerifier, "google") {
		return domain.GoogleProvider
	} else if strings.Contains(lowerVerifier, "facebook") {
		return domain.FacebookProvider
	} else if strings.Contains(lowerVerifier, "apple") {
		return domain.AppleProvider
	} else if strings.Contains(lowerVerifier, "twitter") {
		return domain.TwitterProvider
	} else if strings.Contains(lowerVerifier, "discord") {
		return domain.DiscordProvider
	}

	return domain.Web3AuthProvider
}

// parseUserAgent extracts browser and device info from user agent
func parseUserAgent(userAgent string) string {
	lowerUA := strings.ToLower(userAgent)

	// Extract browser
	var browser string
	switch {
	case strings.Contains(lowerUA, "chrome"):
		browser = "Chrome"
	case strings.Contains(lowerUA, "firefox"):
		browser = "Firefox"
	case strings.Contains(lowerUA, "safari") && !strings.Contains(lowerUA, "chrome"):
		browser = "Safari"
	case strings.Contains(lowerUA, "edge"):
		browser = "Edge"
	default:
		browser = "Unknown Browser"
	}

	// Extract device type
	var device string
	switch {
	case strings.Contains(lowerUA, "iphone"):
		device = "iPhone"
	case strings.Contains(lowerUA, "ipad"):
		device = "iPad"
	case strings.Contains(lowerUA, "android"):
		device = "Android Device"
	case strings.Contains(lowerUA, "macintosh") || strings.Contains(lowerUA, "mac os"):
		device = "Mac"
	case strings.Contains(lowerUA, "windows"):
		device = "Windows PC"
	case strings.Contains(lowerUA, "linux"):
		device = "Linux PC"
	default:
		device = "Unknown Device"
	}

	return fmt.Sprintf("%s on %s", browser, device)
}

// extractOSFromUserAgent extracts OS from user agent
func extractOSFromUserAgent(userAgent string) string {
	lowerUA := strings.ToLower(userAgent)

	switch {
	case strings.Contains(lowerUA, "windows"):
		return "Windows"
	case strings.Contains(lowerUA, "macintosh") || strings.Contains(lowerUA, "mac os"):
		return "MacOS"
	case strings.Contains(lowerUA, "linux") && !strings.Contains(lowerUA, "android"):
		return "Linux"
	case strings.Contains(lowerUA, "android"):
		return "Android"
	case strings.Contains(lowerUA, "iphone") || strings.Contains(lowerUA, "ipad") || strings.Contains(lowerUA, "ios"):
		return "iOS"
	default:
		return "Unknown OS"
	}
}

// determineDeviceType determines the device type from user agent
func determineDeviceType(userAgent string) string {
	lowerUA := strings.ToLower(userAgent)

	switch {
	case strings.Contains(lowerUA, "iphone") || strings.Contains(lowerUA, "android") && strings.Contains(lowerUA, "mobile"):
		return "Mobile"
	case strings.Contains(lowerUA, "ipad") || strings.Contains(lowerUA, "android") && !strings.Contains(lowerUA, "mobile"):
		return "Tablet"
	default:
		return "Desktop"
	}
}

// isValidWalletAddress validates a wallet address format
func isValidWalletAddress(address string) bool {
	// Basic validation - can be expanded for different chains
	if len(address) < 10 {
		return false
	}

	// If Ethereum-style address (0x...)
	if strings.HasPrefix(address, "0x") {
		return len(address) == 42
	}

	return true
}
// InitiatePasswordReset starts the password reset process for email-based accounts
func (a *authService) InitiatePasswordReset(ctx context.Context, email string) error {
	// Check if email exists and is email-based account
	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		// Return generic message for security - don't reveal if email exists
		a.logger.Info("Password reset requested", map[string]interface{}{
			"email": email,
		})
		return nil
	}

	// Check if account was created with email/password
	if user.AuthProvider != "email" {
		a.logger.Info("Password reset attempted for OAuth account", map[string]interface{}{
			"email": email,
			"provider": user.AuthProvider,
		})
		// Return nil instead of error for security - don't reveal details
		return nil
	}

	// Generate OTP
	otpCode := random.RandomOtp()
	otp := domain.OTPVerification{
		ID:           uuid.New(),
		UserID:       user.ID,
		Purpose:      domain.OTPPurposePasswordReset,
		OTPCode:      otpCode,
		ExpiresAt:    time.Now().Add(15 * time.Minute),

	}

	// Store OTP
	_, err = a.otpRepo.CreateOTP(ctx, otp)
	if err != nil {
		a.logger.Error("Failed to create OTP", err, map[string]interface{}{
			"email": email,
		})
		return nil // Don't reveal internal errors
	}

	// Send password reset email
	err = a.emailService.SendPasswordResetEmail(ctx, email, user.FirstName, otp.OTPCode)
	if err != nil {
		a.logger.Error("Failed to send password reset email", err, map[string]interface{}{
			"email": email,
		})
		// Email failure shouldn't be exposed to the user
		return nil
	}

	// Log security event
	a.LogSecurityEvent(ctx, "password_reset_initiated", user.ID, map[string]interface{}{
		"email": email,
	})

	return nil
}

// VerifyResetOTP verifies the OTP but doesn't invalidate it
func (a *authService) VerifyResetOTP(ctx context.Context, email string, code string) error {
	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("invalid email or OTP")
	}

	// Get OTP
	otp, err := a.otpRepo.GetOTPByUserIDAndPurpose(ctx, user.ID, domain.OTPPurposePasswordReset)
	if err != nil {
		return errors.New("invalid or expired OTP")
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return errors.New("OTP has expired")
	}

	// Check attempts
	if otp.AttemptsMade >= otp.MaxAttempts {
		return errors.New("maximum attempts exceeded")
	}

	// Verify code - just check if it's correct without invalidating
	if otp.OTPCode != code {
		// Increment attempts on failure
		a.otpRepo.IncrementAttempts(ctx, otp.ID)
		return errors.New("invalid OTP")
	}

	// Log security event for verification success
	a.LogSecurityEvent(ctx, "password_reset_otp_verified", user.ID, map[string]interface{}{
		"email": email,
	})

	return nil
}

// ResetPassword verifies OTP and resets the user's password in one step
func (a *authService) ResetPassword(ctx context.Context, email string, code string, newPassword string) error {
	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return errors.New("invalid email")
	}

	// Get OTP
	otp, err := a.otpRepo.GetOTPByUserIDAndPurpose(ctx, user.ID, domain.OTPPurposePasswordReset)
	if err != nil {
		return errors.New("invalid or expired OTP")
	}

	// Check if OTP is expired
	if time.Now().After(otp.ExpiresAt) {
		return errors.New("OTP has expired")
	}

	// Check attempts
	if otp.AttemptsMade >= otp.MaxAttempts {
		return errors.New("maximum attempts exceeded")
	}

	// Verify code
	if otp.OTPCode != code {
		// Increment attempts on failure
		a.otpRepo.IncrementAttempts(ctx, otp.ID)
		return errors.New("invalid OTP")
	}

	// Now proceed with password reset
	err = a.userService.ResetUserPassword(ctx, user.ID, newPassword)
	if err != nil {
		return fmt.Errorf("failed to reset password: %w", err)
	}

	// Invalidate the OTP after successful password reset
	err = a.otpRepo.VerifyOTP(ctx, otp.ID, code)
	if err != nil {
		a.logger.Error("Failed to invalidate OTP after password reset", err, map[string]interface{}{
			"otp_id": otp.ID,
		})
	}

	// Block all user sessions
	err = a.sessionRepo.BlockAllUserSessions(ctx, user.ID)
	if err != nil {
		a.logger.Error("Failed to block user sessions after password reset", err, map[string]interface{}{
			"user_id": user.ID,
		})
	}

	// Log security event
	a.LogSecurityEvent(ctx, "password_reset_completed", user.ID, map[string]interface{}{
		"email": user.Email,
	})

	return nil
}