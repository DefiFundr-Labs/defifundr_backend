package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/internal/adapters/dto/request"
	"github.com/demola234/defifundr/internal/adapters/dto/response"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	appErrors "github.com/demola234/defifundr/pkg/app_errors"
	token "github.com/demola234/defifundr/pkg/token_maker"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	authService ports.AuthService
	logger      logging.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService ports.AuthService, logger logging.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) GetUserRepository() ports.UserRepository {
	// Cast to authService to access userRepo (add a method to the interface if preferred)
	if service, ok := h.authService.(interface{ GetUserRepository() ports.UserRepository }); ok {
		return service.GetUserRepository()
	}
	return nil
}

// Web3AuthLogin handles login/registration with Web3Auth
// @Summary Login or register with Web3Auth
// @Description Authenticate or create a new user with Web3Auth tokens
// @Tags authentication
// @Accept json
// @Produce json
// @Param loginRequest body request.Web3AuthLoginRequest true "Web3Auth token"
// @Success 200 {object} response.SuccessResponse "Authentication successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Authentication failed"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/web3auth/login [post]
func (h *AuthHandler) Web3AuthLogin(ctx *gin.Context) {
	// Extract request correlation ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	var req request.Web3AuthLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate request
	if req.WebAuthToken == "" {
		reqLogger.Error("Missing Web3Auth token", nil, nil)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Web3Auth token is required",
		})
		return
	}

	// Authenticate with Web3Auth
	user, session, err := h.authService.AuthenticateWithWeb3(
		ctx,
		req.WebAuthToken,
		ctx.Request.UserAgent(),
		ctx.ClientIP(),
	)

	if err != nil {
		reqLogger.Error("Failed to authenticate with Web3Auth", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Authentication failed: " + err.Error(),
		})
		return
	}

	// Check profile completion status
	profileCompletion, err := h.authService.GetProfileCompletionStatus(ctx, user.ID)
	var completionData *response.ProfileCompletionResponse

	if err == nil {
		completionData = &response.ProfileCompletionResponse{
			CompletionPercentage: profileCompletion.CompletionPercentage,
			MissingFields:        profileCompletion.MissingFields,
			RequiredActions:      profileCompletion.RequiredActions,
		}
	} else {
		reqLogger.Error("Failed to get profile completion status", err, map[string]interface{}{
			"user_id": user.ID,
		})
	}

	// Get user wallets
	wallets, err := h.authService.GetUserWallets(ctx, user.ID)
	var walletResponses []response.UserWalletResponse

	if err == nil {
		walletResponses = make([]response.UserWalletResponse, len(wallets))
		for i, wallet := range wallets {
			walletResponses[i] = response.UserWalletResponse{
				ID:        wallet.ID.String(),
				Address:   wallet.Address,
				Type:      wallet.Type,
				Chain:     wallet.Chain,
				IsDefault: wallet.IsDefault,
			}
		}
	} else {
		reqLogger.Error("Failed to get user wallets", err, map[string]interface{}{
			"user_id": user.ID,
		})
	}

	// Create the session response
	sessionResponse := response.SessionResponse{
		ID:            session.ID,
		UserID:        user.ID,
		AccessToken:   session.OAuthAccessToken,
		UserLoginType: session.UserLoginType,
		ExpiresAt:     session.ExpiresAt,
		CreatedAt:     session.CreatedAt,
	}

	// Create user response
	profilePicture := ""
	if user.ProfilePicture != nil {
		profilePicture = *user.ProfilePicture
	}

	userResponse := response.LoginUserResponse{
		ID:             user.ID.String(),
		Email:          user.Email,
		ProfilePicture: profilePicture,
		AccountType:    user.AccountType,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AuthProvider:   user.AuthProvider,
		ProviderID:     user.ProviderID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}

	// Determine if this is a new registration (user created in the same transaction)
	isNewUser := session.CreatedAt.Sub(user.CreatedAt) < time.Minute

	// Return success response
	responseData := map[string]interface{}{
		"user":               userResponse,
		"session":            sessionResponse,
		"profile_completion": completionData,
		"wallets":            walletResponses,
	}

	// Add onboarding data for new users
	if isNewUser {
		responseData["is_new_user"] = true
		responseData["onboarding_steps"] = []string{
			"complete_profile",
			"verify_email",
			"link_wallet",
		}
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Authentication successful",
		Data:    responseData,
	})

	reqLogger.Info("Web3Auth login successful", map[string]interface{}{
		"user_id":     user.ID,
		"is_new_user": isNewUser,
	})
}

// RegisterUser handles new user registration with email/password
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags authentication
// @Accept json
// @Produce json
// @Param registerRequest body request.RegisterUserRequest true "User registration details"
// @Success 201 {object} response.SuccessResponse "User registered successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 409 {object} response.ErrorResponse "Email already registered"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/register [post]
func (h *AuthHandler) RegisterUser(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	var req request.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid registration request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Create user domain object
	user := domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		AuthProvider: req.Provider,
		WebAuthToken: req.WebAuthToken,
		AccountType:  "personal",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Register user
	createdUser, err := h.authService.RegisterUser(ctx, user, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to register user"

		// Check for specific errors
		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		} else if err.Error() == "email already registered" {
			status = http.StatusConflict
			message = "Email already registered"
		}

		reqLogger.Error("User registration failed", err, map[string]interface{}{
			"email": req.Email,
		})

		ctx.JSON(status, response.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Create a session
	session, err := h.authService.CreateSession(
		ctx,
		createdUser.ID,
		ctx.Request.UserAgent(),
		ctx.ClientIP(),
		"",
		createdUser.Email,
		"email",
	)

	if err != nil {
		reqLogger.Error("Failed to create session", err, map[string]interface{}{
			"user_id": createdUser.ID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Registration successful but failed to create session",
		})
		return
	}

	// Create user response
	userResponse := response.LoginUserResponse{
		ID:           createdUser.ID.String(),
		Email:        createdUser.Email,
		FirstName:    createdUser.FirstName,
		LastName:     createdUser.LastName,
		AccountType:  createdUser.AccountType,
		AuthProvider: createdUser.AuthProvider,
		CreatedAt:    createdUser.CreatedAt,
		UpdatedAt:    createdUser.UpdatedAt,
	}

	// Create session response
	sessionResponse := response.SessionResponse{
		ID:            session.ID,
		UserID:        createdUser.ID,
		AccessToken:   session.OAuthAccessToken,
		UserLoginType: session.UserLoginType,
		ExpiresAt:     session.ExpiresAt,
		CreatedAt:     session.CreatedAt,
	}

	// Return success with onboarding steps
	ctx.JSON(http.StatusCreated, response.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data: map[string]interface{}{
			"user":    userResponse,
			"session": sessionResponse,
			"onboarding_steps": []string{
				"complete_profile",
				"verify_email",
			},
		},
	})

	reqLogger.Info("User registered successfully", map[string]interface{}{
		"user_id": createdUser.ID,
		"email":   createdUser.Email,
	})
}

// Login handles user login with email/password
// @Summary Login with email/password
// @Description Login with email and password credentials
// @Tags authentication
// @Accept json
// @Produce json
// @Param loginRequest body request.LoginRequest true "User login credentials"
// @Success 200 {object} response.SuccessResponse "Login successful"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Invalid email or password"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	var req request.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	if err := req.Validate(); err != nil {
		reqLogger.Error("Validation failed", err, nil)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	user := domain.User{
		Email:        req.Email,
		AuthProvider: req.Provider,
		ProviderID:   req.ProviderID,
		WebAuthToken: req.WebAuthToken,
	}

	loggedInUser, err := h.authService.Login(ctx, req.Email, user, req.Password)
	if err != nil {
		reqLogger.Error("Login failed", err, map[string]interface{}{
			"email": req.Email,
		})
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Invalid email or password",
		})
		return
	}

	session, err := h.authService.CreateSession(ctx, loggedInUser.ID, ctx.Request.UserAgent(), ctx.ClientIP(), req.WebAuthToken, loggedInUser.Email, "login")
	if err != nil {
		reqLogger.Error("Failed to create session", err, nil)
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Login successful but failed to create session",
		})
		return
	}

	profilePicture := ""
	if loggedInUser.ProfilePicture != nil {
		profilePicture = *loggedInUser.ProfilePicture
	}

	userResponse := response.LoginUserResponse{
		ID:             loggedInUser.ID.String(),
		Email:          loggedInUser.Email,
		ProfilePicture: profilePicture,
		AccountType:    loggedInUser.AccountType,
		FirstName:      loggedInUser.FirstName,
		LastName:       loggedInUser.LastName,
		AuthProvider:   loggedInUser.AuthProvider,
		ProviderID:     loggedInUser.ProviderID,
		CreatedAt:      loggedInUser.CreatedAt,
		UpdatedAt:      loggedInUser.UpdatedAt,
	}

	sessionResponse := response.SessionResponse{
		ID:            session.ID,
		UserID:        loggedInUser.ID,
		AccessToken:   session.OAuthAccessToken,
		UserLoginType: req.Provider,
		ExpiresAt:     session.ExpiresAt,
		CreatedAt:     session.CreatedAt,
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"user":    userResponse,
			"session": sessionResponse,
		},
	})
}

// RefreshToken refreshes an access token
// @Summary Refresh access token
// @Description Refresh an expired access token using a refresh token
// @Tags authentication
// @Accept json
// @Produce json
// @Param refreshRequest body request.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.SuccessResponse "Token refreshed successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Invalid or expired refresh token"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	var req request.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid refresh token request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Refresh token
	session, accessToken, err := h.authService.RefreshToken(
		ctx,
		req.RefreshToken,
		ctx.Request.UserAgent(),
		ctx.ClientIP(),
	)

	if err != nil {
		status := http.StatusUnauthorized
		message := "Invalid or expired refresh token"

		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		}

		reqLogger.Error("Failed to refresh token", err, nil)

		ctx.JSON(status, response.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Create session response
	sessionResponse := response.SessionResponse{
		ID:            session.ID,
		UserID:        session.UserID,
		AccessToken:   accessToken,
		UserLoginType: session.UserLoginType,
		ExpiresAt:     session.ExpiresAt,
		CreatedAt:     session.CreatedAt,
	}

	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: map[string]interface{}{
			"session": sessionResponse,
		},
	})

	reqLogger.Info("Token refreshed successfully", map[string]interface{}{
		"session_id": session.ID,
		"user_id":    session.UserID,
	})
}

// UpdatePersonalDetails updates user personal details
// @Summary Update personal details
// @Description Update personal details for a registered user
// @Tags profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param personalDetails body request.RegisterPersonalDetailsRequest true "Personal details"
// @Success 200 {object} response.SuccessResponse "Personal details updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/profile/personal-details [put]

func (h *AuthHandler) UpdatePersonalDetails(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var req request.RegisterPersonalDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid personal details request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Get current user data
	currentUser, err := h.authService.GetUserByID(ctx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": userUUID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve user data",
		})
		return
	}

	// Update user with new details
	currentUser.FirstName = req.FirstName
	currentUser.LastName = req.LastName
	currentUser.Nationality = req.Nationality
	currentUser.PersonalAccountType = req.PersonalAccountType

	if req.PhoneNumber != "" {
		phoneNumber := req.PhoneNumber
		currentUser.PhoneNumber = &phoneNumber
	}

	// Update user
	updatedUser, err := h.authService.RegisterPersonalDetails(ctx, *currentUser)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to update personal details"

		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		}

		reqLogger.Error("Failed to update personal details", err, map[string]interface{}{
			"user_id": userUUID,
		})

		ctx.JSON(status, response.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Create user response
	profilePicture := ""
	if updatedUser.ProfilePicture != nil {
		profilePicture = *updatedUser.ProfilePicture
	}

	userResponse := response.LoginUserResponse{
		ID:                  updatedUser.ID.String(),
		Email:               updatedUser.Email,
		ProfilePicture:      profilePicture,
		AccountType:         updatedUser.AccountType,
		FirstName:           updatedUser.FirstName,
		LastName:            updatedUser.LastName,
		Nationality:         updatedUser.Nationality,
		PersonalAccountType: updatedUser.PersonalAccountType,
		CreatedAt:           updatedUser.CreatedAt,
		UpdatedAt:           updatedUser.UpdatedAt,
	}

	// Get updated profile completion
	profileCompletion, err := h.authService.GetProfileCompletionStatus(ctx, updatedUser.ID)
	var completionData *response.ProfileCompletionResponse

	if err == nil {
		completionData = &response.ProfileCompletionResponse{
			CompletionPercentage: profileCompletion.CompletionPercentage,
			MissingFields:        profileCompletion.MissingFields,
			RequiredActions:      profileCompletion.RequiredActions,
		}
	}

	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Personal details updated successfully",
		Data: map[string]interface{}{
			"user":               userResponse,
			"profile_completion": completionData,
		},
	})

	reqLogger.Info("Personal details updated successfully", map[string]interface{}{
		"user_id": updatedUser.ID,
	})
}

// UpdateAddressDetails updates user address details
// @Summary Update address details
// @Description Update address details for a registered user
// @Tags profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param addressDetails body request.RegisterAddressDetailsRequest true "Address details"
// @Success 200 {object} response.SuccessResponse "Address details updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/profile/address [put]
func (h *AuthHandler) UpdateAddressDetails(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var req request.RegisterAddressDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid address details request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Get current user data
	currentUser, err := h.authService.GetUserByID(ctx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get user by ID", err, map[string]interface{}{
			"user_id": userUUID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve user data",
		})
		return
	}

	// Update user with new details
	userAddress := req.UserAddress
	currentUser.UserAddress = &userAddress
	currentUser.City = req.City
	currentUser.PostalCode = req.PostalCode
	currentUser.ResidentialCountry = &req.Country

	// Update user
	updatedUser, err := h.authService.RegisterAddressDetails(ctx, *currentUser)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to update address details"

		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		}

		reqLogger.Error("Failed to update address details", err, map[string]interface{}{
			"user_id": userUUID,
		})

		ctx.JSON(status, response.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Get updated profile completion
	profileCompletion, err := h.authService.GetProfileCompletionStatus(ctx, updatedUser.ID)
	var completionData *response.ProfileCompletionResponse

	if err == nil {
		completionData = &response.ProfileCompletionResponse{
			CompletionPercentage: profileCompletion.CompletionPercentage,
			MissingFields:        profileCompletion.MissingFields,
			RequiredActions:      profileCompletion.RequiredActions,
		}
	}

	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Address details updated successfully",
		Data: map[string]interface{}{
			"user":               updatedUser,
			"profile_completion": completionData,
		},
	})

	reqLogger.Info("Address details updated successfully", map[string]interface{}{
		"user_id": updatedUser.ID,
	})
}

// UpdateBusinessDetails handles updating a user's business details
// @Summary Update business details
// @Description Update business details for a registered user
// @Tags profile
// @Accept json
// @Produce json
// @Security Bearer
// @Param businessDetails body request.RegisterBusinessDetailsRequest true "Business details"
// @Success 200 {object} response.SuccessResponse "Business details updated successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/profile/business [put]
func (h *AuthHandler) UpdateBusinessDetails(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing update business details request")

	// Get userID from authorization payload in context
	authPayload, exists := ctx.Get("authorization_payload")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Authorization payload not found",
		})
		return
	}
	user := authPayload.(*token.Payload)

	// Parse request body
	var req request.RegisterBusinessDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Validate request
	if err := req.Validate(); err != nil {
		reqLogger.Error("Invalid business details", err, nil)
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Create user domain object
	companyWebsite := req.CompanyWebsite
	employmentType := req.EmploymentType
	userDetails := domain.User{
		ID:                user.UserID,
		CompanyName:       req.CompanyName,
		CompanyAddress:    req.CompanyAddress,
		CompanyCity:       req.CompanyCity,
		CompanyPostalCode: req.CompanyPostalCode,
		CompanyCountry:    req.CompanyCountry,
	}

	if companyWebsite != "" {
		userDetails.CompanyWebsite = &companyWebsite
	}
	if employmentType != "" {
		userDetails.EmploymentType = &employmentType
	}

	// Update user
	updatedUser, err := h.authService.RegisterBusinessDetails(ctx, userDetails)
	if err != nil {
		reqLogger.Error("Failed to update business details", err, map[string]interface{}{
			"user_id": user.UserID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to update business details",
		})
		return
	}

	// Build response
	resp := response.UserResponse{
		ID:         updatedUser.ID.String(),
		Email:      updatedUser.Email,
		FirstName:  updatedUser.FirstName,
		LastName:   updatedUser.LastName,
		Provider:   updatedUser.AuthProvider,
		ProviderID: updatedUser.ProviderID,
		CreatedAt:  updatedUser.CreatedAt,
		UpdatedAt:  updatedUser.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Business details updated successfully",
		Data:    resp,
	})

	reqLogger.Info("Business details updated successfully", map[string]interface{}{
		"user_id": user.UserID,
	})
}

// GetProfileCompletion returns the user's profile completion status
// @Summary Get profile completion status
// @Description Retrieve the profile completion status for the authenticated user
// @Tags profile
// @Produce json
// @Security Bearer
// @Success 200 {object} response.SuccessResponse{data=response.ProfileCompletionResponse} "Profile completion status retrieved"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/profile/completion [get]
func (h *AuthHandler) GetProfileCompletion(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Get profile completion status
	profileCompletion, err := h.authService.GetProfileCompletionStatus(ctx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get profile completion status", err, map[string]interface{}{
			"user_id": userUUID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to get profile completion status",
		})
		return
	}

	// Create response
	completionData := response.ProfileCompletionResponse{
		CompletionPercentage: profileCompletion.CompletionPercentage,
		MissingFields:        profileCompletion.MissingFields,
		RequiredActions:      profileCompletion.RequiredActions,
	}

	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Profile completion status retrieved",
		Data:    completionData,
	})

	reqLogger.Debug("Profile completion status retrieved", map[string]interface{}{
		"user_id":    userUUID,
		"completion": profileCompletion.CompletionPercentage,
		"missing":    len(profileCompletion.MissingFields),
		"actions":    len(profileCompletion.RequiredActions),
	})
}

// LinkWallet links a blockchain wallet to a user
// @Summary Link blockchain wallet
// @Description Link a blockchain wallet to the authenticated user
// @Tags wallet
// @Accept json
// @Produce json
// @Security Bearer
// @Param walletDetails body request.LinkWalletRequest true "Wallet details"
// @Success 200 {object} response.SuccessResponse "Wallet linked successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 409 {object} response.ErrorResponse "Wallet already linked to another account"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/wallet/link [post]

func (h *AuthHandler) LinkWallet(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	var req request.LinkWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid wallet request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Link wallet
	err := h.authService.LinkWallet(ctx, userUUID, req.Address, req.Type, req.Chain)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to link wallet"

		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		} else if err.Error() == "wallet already linked to another account" {
			status = http.StatusConflict
			message = "This wallet is already linked to another account"
		} else if err.Error() == "invalid wallet address format" {
			status = http.StatusBadRequest
			message = "Invalid wallet address format"
		}

		reqLogger.Error("Failed to link wallet", err, map[string]interface{}{
			"user_id": userUUID,
			"address": req.Address,
			"chain":   req.Chain,
		})

		ctx.JSON(status, response.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Get all user wallets
	wallets, err := h.authService.GetUserWallets(ctx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get user wallets", err, map[string]interface{}{
			"user_id": userUUID,
		})
	}

	// Create wallet responses
	var walletResponses []response.UserWalletResponse
	if err == nil {
		walletResponses = make([]response.UserWalletResponse, len(wallets))
		for i, wallet := range wallets {
			walletResponses[i] = response.UserWalletResponse{
				ID:        wallet.ID.String(),
				Address:   wallet.Address,
				Type:      wallet.Type,
				Chain:     wallet.Chain,
				IsDefault: wallet.IsDefault,
			}
		}
	}

	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Wallet linked successfully",
		Data: map[string]interface{}{
			"wallets": walletResponses,
		},
	})

	reqLogger.Info("Wallet linked successfully", map[string]interface{}{
		"user_id": userUUID,
		"address": req.Address,
		"chain":   req.Chain,
	})
}

// GetWallets returns all wallets for a user
// @Summary Get user wallets
// @Description Retrieve all blockchain wallets linked to the authenticated user
// @Tags wallet
// @Produce json
// @Security Bearer
// @Success 200 {object} response.SuccessResponse "Wallets retrieved successfully"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/wallet [get]
func (h *AuthHandler) GetWallets(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Get all user wallets
	wallets, err := h.authService.GetUserWallets(ctx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get user wallets", err, map[string]interface{}{
			"user_id": userUUID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve wallets",
		})
		return
	}

	// Create wallet responses
	walletResponses := make([]response.UserWalletResponse, len(wallets))
	for i, wallet := range wallets {
		walletResponses[i] = response.UserWalletResponse{
			ID:        wallet.ID.String(),
			Address:   wallet.Address,
			Type:      wallet.Type,
			Chain:     wallet.Chain,
			IsDefault: wallet.IsDefault,
		}
	}

	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Wallets retrieved successfully",
		Data: map[string]interface{}{
			"wallets": walletResponses,
		},
	})

	reqLogger.Debug("User wallets retrieved", map[string]interface{}{
		"user_id":      userUUID,
		"wallet_count": len(wallets),
	})
}

// GetUserDevices returns all active devices for the current user
// @Summary Get active devices
// @Description Retrieve all active devices/sessions for the authenticated user
// @Tags security
// @Produce json
// @Security Bearer
// @Success 200 {object} response.SuccessResponse "Active devices retrieved"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/security/devices [get]
func (h *AuthHandler) GetUserDevices(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Get active devices
	devices, err := h.authService.GetActiveDevices(ctx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get active devices", err, map[string]interface{}{
			"user_id": userUUID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to retrieve active devices",
		})
		return
	}

	// Convert to response objects
	deviceResponses := make([]response.DeviceResponse, len(devices))
	for i, device := range devices {
		deviceResponses[i] = response.DeviceResponse{
			SessionID:       device.SessionID.String(),
			Browser:         device.Browser,
			OperatingSystem: device.OperatingSystem,
			DeviceType:      device.DeviceType,
			IPAddress:       device.IPAddress,
			LoginType:       device.LoginType,
			LastUsed:        device.LastUsed,
			CreatedAt:       device.CreatedAt,
		}
	}

	// Return devices
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Active devices retrieved",
		Data: map[string]interface{}{
			"devices": deviceResponses,
		},
	})

	reqLogger.Debug("User devices retrieved", map[string]interface{}{
		"user_id":      userUUID,
		"device_count": len(devices),
	})
}

// RevokeDevice revokes a specific device/session
// @Summary Revoke device
// @Description Revoke a specific device/session for the authenticated user
// @Tags security
// @Accept json
// @Produce json
// @Security Bearer
// @Param revokeRequest body request.RevokeDeviceRequest true "Session ID to revoke"
// @Success 200 {object} response.SuccessResponse "Device revoked successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 403 {object} response.ErrorResponse "Session does not belong to user"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/security/devices/revoke [post]
func (h *AuthHandler) RevokeDevice(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Get session ID from request
	var req request.RevokeDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid revoke device request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Parse session ID
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid session ID",
		})
		return
	}

	// Prevent revoking the current session
	currentSessionID, _ := ctx.Get("session_id")
	if currentSessionID != nil && currentSessionID.(uuid.UUID) == sessionID {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Cannot revoke the current session. Use logout instead.",
		})
		return
	}

	// Revoke session
	err = h.authService.RevokeSession(ctx, userUUID, sessionID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to revoke device"

		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		} else if err.Error() == "session does not belong to user" {
			status = http.StatusForbidden
			message = "Session does not belong to user"
		}

		reqLogger.Error("Failed to revoke device", err, map[string]interface{}{
			"user_id":    userUUID,
			"session_id": sessionID,
		})

		ctx.JSON(status, response.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	// Get updated device list
	devices, _ := h.authService.GetActiveDevices(ctx, userUUID)
	var deviceResponses []response.DeviceResponse

	if devices != nil {
		deviceResponses = make([]response.DeviceResponse, len(devices))
		for i, device := range devices {
			deviceResponses[i] = response.DeviceResponse{
				SessionID:       device.SessionID.String(),
				Browser:         device.Browser,
				OperatingSystem: device.OperatingSystem,
				DeviceType:      device.DeviceType,
				IPAddress:       device.IPAddress,
				LoginType:       device.LoginType,
				LastUsed:        device.LastUsed,
				CreatedAt:       device.CreatedAt,
			}
		}
	}

	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Device revoked successfully",
		Data: map[string]interface{}{
			"devices": deviceResponses,
		},
	})

	reqLogger.Info("Device revoked successfully", map[string]interface{}{
		"user_id":    userUUID,
		"session_id": sessionID,
	})
}

// GetUserSecurityEvents returns security events for the user's account
// @Summary Get security events
// @Description Retrieve security events for the authenticated user's account
// @Tags security
// @Produce json
// @Security Bearer
// @Param type query string false "Filter by event type"
// @Param start_time query string false "Filter by start time (RFC3339 format)"
// @Param end_time query string false "Filter by end time (RFC3339 format)"
// @Success 200 {object} response.SuccessResponse "Security events retrieved"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/security/events [get]
func (h *AuthHandler) GetUserSecurityEvents(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Get optional filters
	eventType := ctx.Query("type")
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")

	var startTime, endTime time.Time
	var err error

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Success: false,
				Message: "Invalid start time format",
			})
			return
		}
	} else {
		// Default to 30 days ago
		startTime = time.Now().AddDate(0, 0, -30)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Success: false,
				Message: "Invalid end time format",
			})
			return
		}
	} else {
		// Default to now
		endTime = time.Now()
	}

	// Get events from security repository
	// Note: authService should expose this method or have a dedicated method
	securityRepo, ok := h.authService.(interface {
		GetSecurityEvents(ctx context.Context, userID uuid.UUID, eventType string, startTime, endTime time.Time) ([]domain.SecurityEvent, error)
	})
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Service does not support security events",
		})
		return
	}

	events, err := securityRepo.GetSecurityEvents(ctx, userUUID, eventType, startTime, endTime)
	if err != nil {
		reqLogger.Error("Failed to get security events", err, map[string]interface{}{
			"user_id": userUUID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to get security events",
		})
		return
	}

	// Convert to response objects
	eventResponses := make([]response.SecurityEventResponse, len(events))
	for i, event := range events {
		eventResponses[i] = response.SecurityEventResponse{
			ID:        event.ID.String(),
			EventType: event.EventType,
			IPAddress: event.IPAddress,
			UserAgent: event.UserAgent,
			Timestamp: event.Timestamp,
			Metadata:  event.Metadata,
		}
	}

	// Return events
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Security events retrieved",
		Data: map[string]interface{}{
			"events": eventResponses,
			"filters": map[string]interface{}{
				"start_time": startTime,
				"end_time":   endTime,
				"event_type": eventType,
			},
		},
	})

	reqLogger.Debug("Security events retrieved", map[string]interface{}{
		"user_id":     userUUID,
		"event_count": len(events),
		"start_time":  startTime,
		"end_time":    endTime,
	})
}

// SetupMFA initiates MFA setup for the user
// @Summary Setup MFA
// @Description Initialize multi-factor authentication for the authenticated user
// @Tags security
// @Produce json
// @Security Bearer
// @Success 200 {object} response.SuccessResponse "MFA setup initiated"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/security/mfa/setup [post]
func (h *AuthHandler) SetupMFA(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Set up MFA
	totpURI, err := h.authService.SetupMFA(ctx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to set up MFA", err, map[string]interface{}{
			"user_id": userUUID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to set up MFA: " + err.Error(),
		})
		return
	}

	// // Return success with TOTP URI
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "MFA setup initiated",
		Data: map[string]interface{}{
			"totp_uri": totpURI,
		},
	})

	reqLogger.Info("MFA setup initiated", map[string]interface{}{
		"user_id": userUUID,
	})
}

// VerifyMFA verifies an MFA code
// @Summary Verify MFA
// @Description Verify an MFA code for the authenticated user
// @Tags security
// @Accept json
// @Produce json
// @Security Bearer
// @Param mfaCode body request.VerifyMFARequest true "MFA code"
// @Success 200 {object} response.SuccessResponse "MFA code verified successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid MFA code"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/security/mfa/verify [post]
func (h *AuthHandler) VerifyMFA(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Convert to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Invalid user ID",
		})
		return
	}

	// Get code from request
	var req request.VerifyMFARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		reqLogger.Error("Invalid MFA verification request format", err, map[string]interface{}{
			"error": err.Error(),
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
		return
	}

	// Verify code
	valid, err := h.authService.VerifyMFA(ctx, userUUID, req.Code)
	if err != nil {
		reqLogger.Error("Failed to verify MFA code", err, map[string]interface{}{
			"user_id": userUUID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to verify MFA code: " + err.Error(),
		})
		return
	}

	if !valid {
		reqLogger.Warn("Invalid MFA code provided", map[string]interface{}{
			"user_id": userUUID,
		})
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid MFA code",
		})
		return
	}

	// Mark the session as MFA verified
	ctx.Set("mfa_verified", true)

	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "MFA code verified successfully",
	})

	reqLogger.Info("MFA code verified successfully", map[string]interface{}{
		"user_id": userUUID,
	})
}

// Logout logs out a user by revoking their session
// @Summary Logout
// @Description Logout the authenticated user by revoking their session
// @Tags authentication
// @Accept json
// @Produce json
// @Security Bearer
// @Param logoutRequest body request.LogoutRequest false "Session ID (optional, defaults to current session)"
// @Success 200 {object} response.SuccessResponse "Logged out successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Failure 500 {object} response.ErrorResponse "Internal server error"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(ctx *gin.Context) {
	// Extract request ID
	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)
	reqLogger.Debug("Processing login request")

	// Get authenticated user ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	// Get session ID from request
	var req request.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If no session ID provided, try to get from context
		sessionID, exists := ctx.Get("session_id")
		if !exists {
			reqLogger.Error("No session ID provided for logout", nil, nil)
			ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
				Success: false,
				Message: "Session ID is required",
			})
			return
		}

		// Use session ID from context
		sessionUUID, ok := sessionID.(uuid.UUID)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Success: false,
				Message: "Invalid session ID",
			})
			return
		}

		// Logout
		err = h.authService.Logout(ctx, sessionUUID)
		if err != nil {
			reqLogger.Error("Failed to logout", err, map[string]interface{}{
				"session_id": sessionUUID,
			})
			ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Success: false,
				Message: "Failed to logout: " + err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, response.SuccessResponse{
			Success: true,
			Message: "Logged out successfully",
		})
		return
	}

	// Parse session ID from request
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Success: false,
			Message: "Invalid session ID",
		})
		return
	}

	// Logout
	err = h.authService.Logout(ctx, sessionID)
	if err != nil {
		reqLogger.Error("Failed to logout", err, map[string]interface{}{
			"session_id": sessionID,
		})
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Success: false,
			Message: "Failed to logout: " + err.Error(),
		})
		return
	}

	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Logged out successfully",
	})

	reqLogger.Info("User logged out", map[string]interface{}{
		"user_id":    userID,
		"session_id": sessionID,
	})
}
