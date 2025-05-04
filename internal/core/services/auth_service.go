// services/auth_service.go (improved)
package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	commons "github.com/demola234/defifundr/infrastructure/hash"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
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
	tokenMaker   tokenMaker.Maker
	config       config.Config
	logger       logging.Logger
}

// SetupMFA sets up multi-factor authentication for a user
func (a *authService) SetupMFA(ctx context.Context, userID uuid.UUID) (string, error) {
	// Check if user exists
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
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
	tokenMaker tokenMaker.Maker,
	config config.Config,
	logger logging.Logger,
) ports.AuthService {
	return &authService{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		oauthRepo:    oauthRepo,
		walletRepo:   walletRepo,
		securityRepo: securityRepo,
		emailService: emailService,
		tokenMaker:   tokenMaker,
		config:       config,
		logger:       logger,
	}
}

// Login implements ports.AuthService.
func (a *authService) Login(ctx context.Context, email string, user domain.User, password string) (*domain.User, error) {
	a.logger.Info("Starting user registration process", map[string]interface{}{
		"email":    email,
		"provider": user.AuthProvider,
	})

	if user.AuthProvider == "email" {
		// Email-based authentication requires a password
		if password == "" {
			a.logger.Error("Password required for email authentication", nil, map[string]interface{}{
				"email": user.Email,
			})
			return nil, errors.New("password is required for email authentication")
		}

		// Check if the user exists
		existingUser, err := a.userRepo.GetUserByEmail(ctx, email)
		if err != nil {
			a.logger.Error("Failed to get user by email", err, map[string]interface{}{
				"email": email,
			})
			return nil, fmt.Errorf("failed to get user by email: %w", err)
		}
		if existingUser == nil {
			a.logger.Error("User not found", nil, map[string]interface{}{
				"email": email,
			})
			return nil, errors.New("user not found")
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
			return nil, errors.New("invalid password")
		}
	} else if user.AuthProvider != "" && user.WebAuthToken != "" {
		// For OAuth or Web3Auth, validate the token and fill user data
		claims, err := a.oauthRepo.ValidateWebAuthToken(ctx, user.WebAuthToken)
		if err != nil {
			a.logger.Error("Failed to validate WebAuth token", err, map[string]interface{}{
				"provider": user.AuthProvider,
			})
			return nil, fmt.Errorf("invalid authentication token: %w", err)
		}

		// Extract user information from OAuth claims
		if claims.Email != "" {
			user.Email = claims.Email
		}
	} else {
		a.logger.Error("Missing authentication credentials", nil, map[string]interface{}{
			"provider": user.AuthProvider,
		})
		return nil, errors.New("missing authentication credentials")
	}
	// Step 2: Check if user with same email already exists
	existingUser, err := a.userRepo.GetUserByEmail(ctx, user.Email)
	if err != nil {
		a.logger.Error("Failed to get user by email", err, map[string]interface{}{
			"email": user.Email,
		})
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if existingUser == nil {
		a.logger.Warn("Login attempt for non-existing email", map[string]interface{}{
			"email": user.Email,
		})
		return nil, errors.New("email not registered")
	}

	// Return the user
	return existingUser, nil
}

// RegisterUser implements the user registration process with Web3Auth integration
func (a *authService) RegisterUser(ctx context.Context, user domain.User, passwordStr string) (*domain.User, error) {
	a.logger.Info("Starting user registration process", map[string]interface{}{
		"email":    user.Email,
		"provider": user.AuthProvider,
	})

	// Step 1: Handle authentication based on provider
	if user.AuthProvider == "email" {
		// Email-based authentication requires a password
		if passwordStr == "" {
			a.logger.Error("Password required for email authentication", nil, map[string]interface{}{
				"email": user.Email,
			})
			return nil, errors.New("password is required for email authentication")
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
	} else if user.AuthProvider != "" && user.WebAuthToken != "" {
		// For OAuth or Web3Auth, validate the token and fill user data
		claims, err := a.oauthRepo.ValidateWebAuthToken(ctx, user.WebAuthToken)
		if err != nil {
			a.logger.Error("Failed to validate WebAuth token", err, map[string]interface{}{
				"provider": user.AuthProvider,
			})
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
			user.ProfilePicture = &profileImage
		}

		// Set provider ID (usually the email for Google OAuth)
		if claims.VerifierID != "" {
			user.ProviderID = claims.VerifierID
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

	// Step 2: Check if user with same email already exists
	existingUser, err := a.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		a.logger.Warn("Registration attempt for existing email", map[string]interface{}{
			"email": user.Email,
		})
		return nil, errors.New("email already registered")
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
func (a *authService) RegisterBusiness(ctx context.Context, user domain.User) (*domain.User, error) {
	// Add Users business details
	// Update the user with business details
	a.logger.Info("Starting user personal details update process", map[string]interface{}{
		"user_id": user.ID,
	})

	// Get the existing user by ID
	existingUser, err := a.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		a.logger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Update only the personal details fields, keeping other fields as they are
	updatedUser := *existingUser
	// Update only the company details fields, keeping other fields as they are
	updatedUser.CompanyName = user.CompanyName
	updatedUser.CompanyAddress = user.CompanyAddress
	updatedUser.CompanyCity = user.CompanyCity
	updatedUser.CompanyCountry = user.CompanyCountry
	updatedUser.CompanyPostalCode = user.CompanyPostalCode
	updatedUser.CompanyWebsite = user.CompanyWebsite

	// Update the user in the database
	users, err := a.userRepo.UpdateUserBusinessDetails(ctx, updatedUser)
	if err != nil {
		a.logger.Error("Failed to update user", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return users, nil
}

// RegisterPersonalDetails implements ports.AuthService
func (a *authService) RegisterPersonalDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	a.logger.Info("Starting user personal details update process", map[string]interface{}{
		"user_id": user.ID,
	})

	// Get the existing user by ID
	existingUser, err := a.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		a.logger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Update only the personal details fields, keeping other fields as they are
	updatedUser := *existingUser
	updatedUser.Nationality = user.Nationality

	if user.AccountType != "" {
		updatedUser.AccountType = user.AccountType
	}

	if user.PersonalAccountType != "" {
		updatedUser.PersonalAccountType = user.PersonalAccountType
	}

	if user.PhoneNumber != nil {
		updatedUser.PhoneNumber = user.PhoneNumber
	}

	// Update the user with personal details
	result, err := a.userRepo.UpdateUserPersonalDetails(ctx, updatedUser)
	if err != nil {
		a.logger.Error("Failed to update user personal details", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to update user personal details: %w", err)
	}

	a.logger.Info("User personal details updated successfully", map[string]interface{}{
		"user_id": user.ID,
	})

	return result, nil
}

// RegisterBusinessDetails implements ports.AuthService
func (a *authService) RegisterBusinessDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	a.logger.Info("Starting business details update process", map[string]interface{}{
		"user_id": user.ID,
	})

	// Get the existing user by ID
	existingUser, err := a.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		a.logger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Update only the business details fields, keeping other fields as they are
	updatedUser := *existingUser

	if user.CompanyName != "" {
		updatedUser.CompanyName = user.CompanyName
	}

	if user.CompanyAddress != "" {
		updatedUser.CompanyAddress = user.CompanyAddress
	}

	if user.CompanyCity != "" {
		updatedUser.CompanyCity = user.CompanyCity
	}

	if user.CompanyPostalCode != "" {
		updatedUser.CompanyPostalCode = user.CompanyPostalCode
	}

	if user.CompanyCountry != "" {
		updatedUser.CompanyCountry = user.CompanyCountry
	}

	if user.CompanyWebsite != nil {
		updatedUser.CompanyWebsite = user.CompanyWebsite
	}

	if user.EmploymentType != nil {
		updatedUser.EmploymentType = user.EmploymentType
	}

	// Update the user with business details
	result, err := a.userRepo.UpdateUserBusinessDetails(ctx, updatedUser)
	if err != nil {
		a.logger.Error("Failed to update business details", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to update business details: %w", err)
	}

	a.logger.Info("Business details updated successfully", map[string]interface{}{
		"user_id": user.ID,
	})

	return result, nil
}

// RegisterAddressDetails implements ports.AuthService
func (a *authService) RegisterAddressDetails(ctx context.Context, user domain.User) (*domain.User, error) {
	a.logger.Info("Starting address details update process", map[string]interface{}{
		"user_id": user.ID,
	})

	// Get the existing user by ID
	existingUser, err := a.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		a.logger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	// Update only the address details fields, keeping other fields as they are
	updatedUser := *existingUser

	if user.UserAddress != nil {
		updatedUser.UserAddress = user.UserAddress
	}

	if user.City != "" {
		updatedUser.City = user.City
	}

	if user.PostalCode != "" {
		updatedUser.PostalCode = user.PostalCode
	}

	// Update the user with address details
	result, err := a.userRepo.UpdateUserAddressDetails(ctx, updatedUser)
	if err != nil {
		a.logger.Error("Failed to update address details", err, map[string]interface{}{
			"user_id": user.ID,
		})
		return nil, fmt.Errorf("failed to update address details: %w", err)
	}

	a.logger.Info("Address details updated successfully", map[string]interface{}{
		"user_id": user.ID,
	})

	return result, nil
}

// GetUserByEmail implements ports.AuthService.
func (a *authService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// GetUserByID implements ports.AuthService.
func (a *authService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

func (a *authService) CheckEmailExists(ctx context.Context, email string) (bool, error) {
	exists, err := a.userRepo.CheckEmailExists(ctx, email)
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}

	return exists, nil
}

// AuthenticateWithWeb3 implements unified Web3Auth authentication flow
func (a *authService) AuthenticateWithWeb3(ctx context.Context, webAuthToken string, userAgent string, clientIP string) (*domain.User, *domain.Session, error) {
	// Validate the Web3Auth token
	claims, err := a.oauthRepo.ValidateWebAuthToken(ctx, webAuthToken)
	if err != nil {
		a.logger.Error("Web3Auth token validation failed", err, nil)
		return nil, nil, err
	}

	// Extract identity information
	email := claims.Email
	if email == "" {
		return nil, nil, errors.New("email not provided in Web3Auth token")
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
		profileImage := claims.ProfileImage

		// Determine provider from verifier
		authProvider := mapVerifierToProvider(claims.Verifier)

		// Create new user
		newUser := domain.User{
			ID:                  uuid.New(),
			Email:               email,
			FirstName:           firstName,
			LastName:            lastName,
			ProfilePicture:      &profileImage,
			ProviderID:          claims.VerifierID,
			AuthProvider:        string(authProvider),
			AccountType:         "personal", // Default value, can be updated later
			PersonalAccountType: "user",     // Default value, can be updated later
			Nationality:         "unknown",  // Default value, can be updated later
		}

		user, err = a.userRepo.CreateUser(ctx, newUser)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create user: %w", err)
		}

		isNewUser = true

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
		if claims.ProfileImage != "" && (user.ProfilePicture == nil || *user.ProfilePicture != claims.ProfileImage) {
			profileImage := claims.ProfileImage
			user.ProfilePicture = &profileImage
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
			a.userRepo.UpdateUser(ctx, *user)
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
	return a.walletRepo.GetWalletsByUserID(ctx, userID)
}

// GetProfileCompletionStatus calculates profile completion percentage
func (a *authService) GetProfileCompletionStatus(ctx context.Context, userID uuid.UUID) (*domain.ProfileCompletion, error) {
	// Get user data
	user, err := a.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Define required fields
	type fieldCheck struct {
		name     string
		required bool
		value    interface{}
	}

	// Common required fields
	fields := []fieldCheck{
		{"First Name", true, user.FirstName != ""},
		{"Last Name", true, user.LastName != ""},
		{"Nationality", true, user.Nationality != "" && user.Nationality != "unknown"},
	}

	// Account type specific fields
	if user.AccountType == "business" {
		fields = append(fields, []fieldCheck{
			{"Company Name", true, user.CompanyName != ""},
			{"Company Address", true, user.CompanyAddress != ""},
			{"Company City", true, user.CompanyCity != ""},
			{"Company Country", true, user.CompanyCountry != ""},
		}...)
	} else {
		fields = append(fields, []fieldCheck{
			{"Address", true, user.UserAddress != nil && *user.UserAddress != ""},
			{"City", true, user.City != ""},
			{"Postal Code", true, user.PostalCode != ""},
		}...)
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

	// Create profile completion response
	completion := &domain.ProfileCompletion{
		UserID:               userID,
		CompletionPercentage: percentage,
		MissingFields:        missingFields,
		RequiredActions:      requiredActions,
	}

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

	// Create session
	session := domain.Session{
		ID:               uuid.New(),
		UserID:           userID,
		RefreshToken:     refreshToken,
		OAuthAccessToken: accessToken,
		UserAgent:        userAgent,
		ClientIP:         clientIP,
		IsBlocked:        false,
		MFAEnabled:       false,
		UserLoginType:    loginType,
		ExpiresAt:        time.Now().Add(a.config.AccessTokenDuration),
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
	// Get session to verify ownership
	session, err := a.sessionRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get session: %w", err)
	}

	// Verify session belongs to user
	if session.UserID != userID {
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
	return a.sessionRepo.DeleteSession(ctx, sessionID)
}

// RefreshToken refreshes an access token
func (a *authService) RefreshToken(ctx context.Context, refreshToken, userAgent, clientIP string) (*domain.Session, string, error) {
	// Get session by refresh token
	session, err := a.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	// Validate session
	if session == nil || session.IsBlocked || time.Now().After(session.ExpiresAt) {
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
