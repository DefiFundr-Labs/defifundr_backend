package authhandler

import (
	"net/http"
	"strings"
	"time"

	"github.com/demola234/defifundr/infrastructure/common/logging"
	authdto "github.com/demola234/defifundr/internal/features/auth/dto"
	authport "github.com/demola234/defifundr/internal/features/auth/port"
	userdomain "github.com/demola234/defifundr/internal/features/user/domain"
	userport "github.com/demola234/defifundr/internal/features/user/port"
	appErrors "github.com/demola234/defifundr/pkg/apperrors"
	token "github.com/demola234/defifundr/pkg/token"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles authentication-related HTTP requests.
type Handler struct {
	service authport.AuthService
	logger  logging.Logger
}

// New creates a new auth Handler.
func New(service authport.AuthService, logger logging.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GetUserRepository exposes the user repository for use by middleware.
func (h *Handler) GetUserRepository() userport.UserRepository {
	return h.service.GetUserRepository()
}

// ────────────────────────────── Web3AuthLogin ──────────────────────────────

func (h *Handler) Web3AuthLogin(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "Web3AuthLogin")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	var req authdto.Web3AuthLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}
	if req.WebAuthToken == "" {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Web3Auth token is required"})
		return
	}

	user, session, err := h.service.AuthenticateWithWeb3(spanCtx, req.WebAuthToken, ctx.Request.UserAgent(), ctx.ClientIP())
	if err != nil {
		reqLogger.Error("Failed to authenticate with Web3Auth", err, nil)
		ctx.JSON(http.StatusUnauthorized, authdto.ErrorResponse{Success: false, Message: "Authentication failed: " + err.Error()})
		return
	}

	profileCompletion, _ := h.service.GetProfileCompletionStatus(spanCtx, user.ID)
	var completionData *authdto.ProfileCompletionResponse
	if profileCompletion != nil {
		completionData = &authdto.ProfileCompletionResponse{
			CompletionPercentage: profileCompletion.CompletionPercentage,
			MissingFields:        profileCompletion.MissingFields,
			RequiredActions:      profileCompletion.RequiredActions,
		}
	}

	wallets, _ := h.service.GetUserWallets(spanCtx, user.ID)
	walletResponses := make([]authdto.UserWalletResponse, len(wallets))
	for i, w := range wallets {
		walletResponses[i] = authdto.UserWalletResponse{ID: w.ID.String(), Address: w.Address, Type: w.Type, Chain: w.Chain, IsDefault: w.IsDefault}
	}

	profilePicture := ""
	if user.ProfilePicture != nil {
		profilePicture = *user.ProfilePicture
	}
	isNewUser := session.CreatedAt.Sub(user.CreatedAt) < time.Minute

	responseData := map[string]interface{}{
		"user": authdto.LoginUserResponse{
			ID: user.ID.String(), Email: user.Email, ProfilePicture: profilePicture,
			AccountType: user.AccountType, FirstName: user.FirstName, LastName: user.LastName,
			AuthProvider: user.AuthProvider, ProviderID: user.ProviderID,
			CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt,
		},
		"session": authdto.SessionResponse{
			ID: session.ID, UserID: user.ID, AccessToken: session.OAuthAccessToken,
			UserLoginType: session.UserLoginType, ExpiresAt: session.ExpiresAt, CreatedAt: session.CreatedAt,
		},
		"profile_completion": completionData,
		"wallets":            walletResponses,
	}
	if isNewUser {
		responseData["is_new_user"] = true
		responseData["onboarding_steps"] = []string{"complete_profile", "verify_email", "link_wallet"}
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{Success: true, Message: "Authentication successful", Data: responseData})
	reqLogger.Info("Web3Auth login successful", map[string]interface{}{"user_id": user.ID, "is_new_user": isNewUser})
}

// ────────────────────────────── RegisterUser ──────────────────────────────

func (h *Handler) RegisterUser(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "RegisterUser")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	var req authdto.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	user := userdomain.User{
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

	createdUser, err := h.service.RegisterUser(spanCtx, user, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to register user"
		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		} else if err.Error() == "email already registered" {
			status = http.StatusConflict
			message = "Email already registered"
		}
		reqLogger.Error("User registration failed", err, map[string]interface{}{"email": req.Email})
		ctx.JSON(status, authdto.ErrorResponse{Success: false, Message: message})
		return
	}

	session, err := h.service.CreateSession(spanCtx, createdUser.ID, ctx.Request.UserAgent(), ctx.ClientIP(), "", createdUser.Email, "email")
	if err != nil {
		reqLogger.Error("Failed to create session", err, nil)
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Registration successful but failed to create session"})
		return
	}

	ctx.JSON(http.StatusCreated, authdto.SuccessResponse{
		Success: true,
		Message: "User registered successfully",
		Data: map[string]interface{}{
			"user": authdto.LoginUserResponse{
				ID: createdUser.ID.String(), Email: createdUser.Email,
				FirstName: createdUser.FirstName, LastName: createdUser.LastName,
				AccountType: createdUser.AccountType, AuthProvider: createdUser.AuthProvider,
				CreatedAt: createdUser.CreatedAt, UpdatedAt: createdUser.UpdatedAt,
			},
			"session": authdto.SessionResponse{
				ID: session.ID, UserID: createdUser.ID, AccessToken: session.OAuthAccessToken,
				UserLoginType: session.UserLoginType, ExpiresAt: session.ExpiresAt, CreatedAt: session.CreatedAt,
			},
			"onboarding_steps": []string{"complete_profile", "verify_email"},
		},
	})
	reqLogger.Info("User registered successfully", map[string]interface{}{"user_id": createdUser.ID, "email": createdUser.Email})
}

// ────────────────────────────── Login ──────────────────────────────

func (h *Handler) Login(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "Login")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	var req authdto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	user := userdomain.User{Email: req.Email, AuthProvider: req.Provider, ProviderID: req.ProviderID, WebAuthToken: req.WebAuthToken}
	loggedInUser, err := h.service.Login(spanCtx, req.Email, user, req.Password)
	if err != nil {
		reqLogger.Error("Login failed", err, map[string]interface{}{"email": req.Email})
		ctx.JSON(http.StatusUnauthorized, authdto.ErrorResponse{Success: false, Message: "Invalid email or password"})
		return
	}

	session, err := h.service.CreateSession(spanCtx, loggedInUser.ID, ctx.Request.UserAgent(), ctx.ClientIP(), req.WebAuthToken, loggedInUser.Email, "login")
	if err != nil {
		reqLogger.Error("Failed to create session", err, nil)
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Login successful but failed to create session"})
		return
	}

	profilePicture := ""
	if loggedInUser.ProfilePicture != nil {
		profilePicture = *loggedInUser.ProfilePicture
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Login successful",
		Data: map[string]interface{}{
			"user": authdto.LoginUserResponse{
				ID: loggedInUser.ID.String(), Email: loggedInUser.Email,
				ProfilePicture: profilePicture, AccountType: loggedInUser.AccountType,
				FirstName: loggedInUser.FirstName, LastName: loggedInUser.LastName,
				AuthProvider: loggedInUser.AuthProvider, ProviderID: loggedInUser.ProviderID,
				CreatedAt: loggedInUser.CreatedAt, UpdatedAt: loggedInUser.UpdatedAt,
			},
			"session": authdto.SessionResponse{
				ID: session.ID, UserID: loggedInUser.ID, AccessToken: session.OAuthAccessToken,
				UserLoginType: req.Provider, ExpiresAt: session.ExpiresAt, CreatedAt: session.CreatedAt,
			},
		},
	})
}

// ────────────────────────────── RefreshToken ──────────────────────────────

func (h *Handler) RefreshToken(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "RefreshToken")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	var req authdto.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	session, accessToken, err := h.service.RefreshToken(spanCtx, req.RefreshToken, ctx.Request.UserAgent(), ctx.ClientIP())
	if err != nil {
		status := http.StatusUnauthorized
		message := "Invalid or expired refresh token"
		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		}
		reqLogger.Error("Failed to refresh token", err, nil)
		ctx.JSON(status, authdto.ErrorResponse{Success: false, Message: message})
		return
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Token refreshed successfully",
		Data: map[string]interface{}{
			"session": authdto.SessionResponse{
				ID: session.ID, UserID: session.UserID, AccessToken: accessToken,
				UserLoginType: session.UserLoginType, ExpiresAt: session.ExpiresAt, CreatedAt: session.CreatedAt,
			},
		},
	})
	reqLogger.Info("Token refreshed successfully", map[string]interface{}{"session_id": session.ID, "user_id": session.UserID})
}

// ────────────────────────────── UpdatePersonalDetails ──────────────────────────────

func (h *Handler) UpdatePersonalDetails(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "UpdatePersonalDetails")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	var req authdto.RegisterPersonalDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	currentUser, err := h.service.GetUserByID(spanCtx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get user", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to retrieve user data"})
		return
	}

	currentUser.FirstName = req.FirstName
	currentUser.LastName = req.LastName
	currentUser.Nationality = req.Nationality
	currentUser.PersonalAccountType = req.PersonalAccountType
	if req.PhoneNumber != "" {
		pn := req.PhoneNumber
		currentUser.PhoneNumber = &pn
	}

	updatedUser, err := h.service.RegisterPersonalDetails(spanCtx, *currentUser)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to update personal details"
		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		}
		reqLogger.Error("Failed to update personal details", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(status, authdto.ErrorResponse{Success: false, Message: message})
		return
	}

	profilePicture := ""
	if updatedUser.ProfilePicture != nil {
		profilePicture = *updatedUser.ProfilePicture
	}

	profileCompletion, _ := h.service.GetProfileCompletionStatus(spanCtx, updatedUser.ID)
	var completionData *authdto.ProfileCompletionResponse
	if profileCompletion != nil {
		completionData = &authdto.ProfileCompletionResponse{
			CompletionPercentage: profileCompletion.CompletionPercentage,
			MissingFields:        profileCompletion.MissingFields,
			RequiredActions:      profileCompletion.RequiredActions,
		}
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Personal details updated successfully",
		Data: map[string]interface{}{
			"user": authdto.LoginUserResponse{
				ID: updatedUser.ID.String(), Email: updatedUser.Email,
				ProfilePicture: profilePicture, AccountType: updatedUser.AccountType,
				FirstName: updatedUser.FirstName, LastName: updatedUser.LastName,
				Nationality: updatedUser.Nationality, PersonalAccountType: updatedUser.PersonalAccountType,
				CreatedAt: updatedUser.CreatedAt, UpdatedAt: updatedUser.UpdatedAt,
			},
			"profile_completion": completionData,
		},
	})
	reqLogger.Info("Personal details updated", map[string]interface{}{"user_id": updatedUser.ID})
}

// ────────────────────────────── UpdateAddressDetails ──────────────────────────────

func (h *Handler) UpdateAddressDetails(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "UpdateAddressDetails")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	var req authdto.RegisterAddressDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	currentUser, err := h.service.GetUserByID(spanCtx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get user", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to retrieve user data"})
		return
	}

	userAddress := req.UserAddress
	currentUser.UserAddress = &userAddress
	currentUser.City = req.City
	currentUser.PostalCode = req.PostalCode
	currentUser.ResidentialCountry = &req.Country

	updatedUser, err := h.service.RegisterAddressDetails(spanCtx, *currentUser)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to update address details"
		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		}
		reqLogger.Error("Failed to update address details", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(status, authdto.ErrorResponse{Success: false, Message: message})
		return
	}

	profileCompletion, _ := h.service.GetProfileCompletionStatus(spanCtx, updatedUser.ID)
	var completionData *authdto.ProfileCompletionResponse
	if profileCompletion != nil {
		completionData = &authdto.ProfileCompletionResponse{
			CompletionPercentage: profileCompletion.CompletionPercentage,
			MissingFields:        profileCompletion.MissingFields,
			RequiredActions:      profileCompletion.RequiredActions,
		}
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Address details updated successfully",
		Data: map[string]interface{}{
			"user":               updatedUser,
			"profile_completion": completionData,
		},
	})
	reqLogger.Info("Address details updated", map[string]interface{}{"user_id": updatedUser.ID})
}

// ────────────────────────────── UpdateBusinessDetails ──────────────────────────────

func (h *Handler) UpdateBusinessDetails(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "UpdateBusinessDetails")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	authPayload, exists := ctx.Get("authorization_payload")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, authdto.ErrorResponse{Success: false, Message: "Authorization payload not found"})
		return
	}
	payload := authPayload.(*token.Payload)

	var req authdto.RegisterBusinessDetailsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	companyInfo := userdomain.CompanyInfo{
		UserID:              payload.UserID,
		CompanyName:         &req.CompanyName,
		CompanySize:         &req.CompanySize,
		CompanyIndustry:     &req.CompanyIndustry,
		CompanyDescription:  &req.CompanyDescription,
		CompanyHeadquarters: &req.CompanyCountry,
		AccountType:         req.AccountType,
	}

	_, err := h.service.RegisterBusinessDetails(spanCtx, companyInfo)
	if err != nil {
		reqLogger.Error("Failed to update business details", err, map[string]interface{}{"user_id": payload.UserID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to update business details"})
		return
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{Success: true, Message: "Business details updated successfully"})
	reqLogger.Info("Business details updated", map[string]interface{}{"user_id": payload.UserID})
}

// ────────────────────────────── GetProfileCompletion ──────────────────────────────

func (h *Handler) GetProfileCompletion(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "GetProfileCompletion")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	profileCompletion, err := h.service.GetProfileCompletionStatus(spanCtx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get profile completion status", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to get profile completion status"})
		return
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Profile completion status retrieved",
		Data: authdto.ProfileCompletionResponse{
			CompletionPercentage: profileCompletion.CompletionPercentage,
			MissingFields:        profileCompletion.MissingFields,
			RequiredActions:      profileCompletion.RequiredActions,
		},
	})
}

// ────────────────────────────── LinkWallet ──────────────────────────────

func (h *Handler) LinkWallet(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "LinkWallet")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	var req authdto.LinkWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	err := h.service.LinkWallet(spanCtx, userUUID, req.Address, req.Type, req.Chain)
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
		reqLogger.Error("Failed to link wallet", err, map[string]interface{}{"user_id": userUUID, "address": req.Address})
		ctx.JSON(status, authdto.ErrorResponse{Success: false, Message: message})
		return
	}

	wallets, _ := h.service.GetUserWallets(spanCtx, userUUID)
	walletResponses := make([]authdto.UserWalletResponse, len(wallets))
	for i, w := range wallets {
		walletResponses[i] = authdto.UserWalletResponse{ID: w.ID.String(), Address: w.Address, Type: w.Type, Chain: w.Chain, IsDefault: w.IsDefault}
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Wallet linked successfully",
		Data:    map[string]interface{}{"wallets": walletResponses},
	})
	reqLogger.Info("Wallet linked successfully", map[string]interface{}{"user_id": userUUID, "address": req.Address})
}

// ────────────────────────────── GetWallets ──────────────────────────────

func (h *Handler) GetWallets(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "GetWallets")
	defer span.End()

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	wallets, err := h.service.GetUserWallets(spanCtx, userUUID)
	if err != nil {
		h.logger.Error("Failed to get user wallets", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to retrieve wallets"})
		return
	}

	walletResponses := make([]authdto.UserWalletResponse, len(wallets))
	for i, w := range wallets {
		walletResponses[i] = authdto.UserWalletResponse{ID: w.ID.String(), Address: w.Address, Type: w.Type, Chain: w.Chain, IsDefault: w.IsDefault}
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Wallets retrieved successfully",
		Data:    map[string]interface{}{"wallets": walletResponses},
	})
}

// ────────────────────────────── GetUserDevices ──────────────────────────────

func (h *Handler) GetUserDevices(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "GetUserDevices")
	defer span.End()

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	devices, err := h.service.GetActiveDevices(spanCtx, userUUID)
	if err != nil {
		h.logger.Error("Failed to get active devices", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to retrieve active devices"})
		return
	}

	deviceResponses := make([]authdto.DeviceResponse, len(devices))
	for i, d := range devices {
		deviceResponses[i] = authdto.DeviceResponse{
			SessionID:       d.SessionID.String(),
			Browser:         d.Browser,
			OperatingSystem: d.OperatingSystem,
			DeviceType:      d.DeviceType,
			IPAddress:       d.IPAddress,
			LoginType:       d.LoginType,
			LastUsed:        d.LastUsed,
			CreatedAt:       d.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Active devices retrieved",
		Data:    map[string]interface{}{"devices": deviceResponses},
	})
}

// ────────────────────────────── RevokeDevice ──────────────────────────────

func (h *Handler) RevokeDevice(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "RevokeDevice")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	var req authdto.RevokeDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid session ID"})
		return
	}

	currentSessionID, _ := ctx.Get("session_id")
	if currentSessionID != nil && currentSessionID.(uuid.UUID) == sessionID {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Cannot revoke the current session. Use logout instead."})
		return
	}

	if err := h.service.RevokeSession(spanCtx, userUUID, sessionID); err != nil {
		status := http.StatusInternalServerError
		message := "Failed to revoke device"
		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		} else if err.Error() == "session does not belong to user" {
			status = http.StatusForbidden
			message = "Session does not belong to user"
		}
		reqLogger.Error("Failed to revoke device", err, map[string]interface{}{"user_id": userUUID, "session_id": sessionID})
		ctx.JSON(status, authdto.ErrorResponse{Success: false, Message: message})
		return
	}

	devices, _ := h.service.GetActiveDevices(spanCtx, userUUID)
	deviceResponses := make([]authdto.DeviceResponse, len(devices))
	for i, d := range devices {
		deviceResponses[i] = authdto.DeviceResponse{
			SessionID: d.SessionID.String(), Browser: d.Browser, OperatingSystem: d.OperatingSystem,
			DeviceType: d.DeviceType, IPAddress: d.IPAddress, LoginType: d.LoginType,
			LastUsed: d.LastUsed, CreatedAt: d.CreatedAt,
		}
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Device revoked successfully",
		Data:    map[string]interface{}{"devices": deviceResponses},
	})
	reqLogger.Info("Device revoked", map[string]interface{}{"user_id": userUUID, "session_id": sessionID})
}

// ────────────────────────────── GetUserSecurityEvents ──────────────────────────────

func (h *Handler) GetUserSecurityEvents(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "GetUserSecurityEvents")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	// Parse optional time filters (for future use)
	startTimeStr := ctx.Query("start_time")
	endTimeStr := ctx.Query("end_time")
	_ = startTimeStr
	_ = endTimeStr

	events, err := h.service.GetUserSecurityEvents(spanCtx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to get security events", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to get security events"})
		return
	}

	eventResponses := make([]authdto.SecurityEventResponse, len(events))
	for i, e := range events {
		eventResponses[i] = authdto.SecurityEventResponse{
			ID:        e.ID.String(),
			EventType: e.EventType,
			IPAddress: e.IPAddress,
			UserAgent: e.UserAgent,
			Timestamp: e.Timestamp,
			Metadata:  e.Metadata,
		}
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "Security events retrieved",
		Data:    map[string]interface{}{"events": eventResponses},
	})
}

// ────────────────────────────── SetupMFA ──────────────────────────────

func (h *Handler) SetupMFA(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "SetupMFA")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	totpURI, err := h.service.SetupMFA(spanCtx, userUUID)
	if err != nil {
		reqLogger.Error("Failed to set up MFA", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to set up MFA: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{
		Success: true,
		Message: "MFA setup initiated",
		Data:    map[string]interface{}{"totp_uri": totpURI},
	})
	reqLogger.Info("MFA setup initiated", map[string]interface{}{"user_id": userUUID})
}

// ────────────────────────────── VerifyMFA ──────────────────────────────

func (h *Handler) VerifyMFA(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "VerifyMFA")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	userUUID, ok := getUserUUID(ctx)
	if !ok {
		return
	}

	var req authdto.VerifyMFARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	valid, err := h.service.VerifyMFA(spanCtx, userUUID, req.Code)
	if err != nil {
		reqLogger.Error("Failed to verify MFA code", err, map[string]interface{}{"user_id": userUUID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to verify MFA code: " + err.Error()})
		return
	}
	if !valid {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid MFA code"})
		return
	}

	ctx.Set("mfa_verified", true)
	ctx.JSON(http.StatusOK, authdto.SuccessResponse{Success: true, Message: "MFA code verified successfully"})
	reqLogger.Info("MFA code verified", map[string]interface{}{"user_id": userUUID})
}

// ────────────────────────────── Logout ──────────────────────────────

func (h *Handler) Logout(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "Logout")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	_, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, authdto.ErrorResponse{Success: false, Message: "Unauthorized"})
		return
	}

	var req authdto.LogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.SessionID == "" {
		// Fall back to session_id from context
		sessionID, exists := ctx.Get("session_id")
		if !exists {
			ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Session ID is required"})
			return
		}
		sessionUUID, ok := sessionID.(uuid.UUID)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Invalid session ID"})
			return
		}
		if err := h.service.Logout(spanCtx, sessionUUID); err != nil {
			reqLogger.Error("Failed to logout", err, map[string]interface{}{"session_id": sessionUUID})
			ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to logout: " + err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, authdto.SuccessResponse{Success: true, Message: "Logged out successfully"})
		return
	}

	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid session ID"})
		return
	}
	if err := h.service.Logout(spanCtx, sessionID); err != nil {
		reqLogger.Error("Failed to logout", err, map[string]interface{}{"session_id": sessionID})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to logout: " + err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, authdto.SuccessResponse{Success: true, Message: "Logged out successfully"})
	reqLogger.Info("User logged out", map[string]interface{}{"session_id": sessionID})
}

// ────────────────────────────── InitiatePasswordReset ──────────────────────────────

func (h *Handler) InitiatePasswordReset(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "InitiatePasswordReset")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	var req authdto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	if err := h.service.InitiatePasswordReset(spanCtx, req.Email); err != nil {
		if err.Error() == "password reset not available for OAuth accounts" {
			ctx.JSON(http.StatusForbidden, authdto.ErrorResponse{
				Success: false,
				Message: "Password reset is not available for OAuth accounts. Please use your social login provider.",
			})
			return
		}
		reqLogger.Error("Password reset initiation failed", err, map[string]interface{}{"email": req.Email})
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Failed to process password reset request"})
		return
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{Success: true, Message: "If this email exists, you will receive password reset instructions"})
	reqLogger.Info("Password reset initiated", map[string]interface{}{"email": req.Email})
}

// ────────────────────────────── VerifyResetOTP ──────────────────────────────

func (h *Handler) VerifyResetOTP(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "VerifyResetOTP")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	var req authdto.VerifyResetOTPRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	if err := h.service.VerifyResetOTP(spanCtx, req.Email, req.OTP); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "maximum attempts exceeded" {
			status = http.StatusTooManyRequests
		}
		reqLogger.Error("OTP verification failed", err, map[string]interface{}{"email": req.Email})
		ctx.JSON(status, authdto.ErrorResponse{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{Success: true, Message: "OTP verified successfully"})
	reqLogger.Info("OTP verified successfully", map[string]interface{}{"email": req.Email})
}

// ────────────────────────────── ResetPassword ──────────────────────────────

func (h *Handler) ResetPassword(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "ResetPassword")
	defer span.End()

	requestID, _ := ctx.Get("RequestID")
	reqLogger := h.logger.With("request_id", requestID)

	var req authdto.CompletePasswordResetRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, authdto.ErrorResponse{Success: false, Message: "Invalid request format: " + err.Error()})
		return
	}

	if err := h.service.ResetPassword(spanCtx, req.Email, req.OTP, req.NewPassword); err != nil {
		status := http.StatusBadRequest
		message := err.Error()
		switch {
		case strings.Contains(message, "OTP has expired"):
			status = http.StatusUnauthorized
		case strings.Contains(message, "maximum attempts exceeded"):
			status = http.StatusTooManyRequests
		case strings.Contains(message, "invalid OTP"):
			status = http.StatusUnauthorized
		case strings.Contains(message, "password must be"):
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
			message = "Failed to reset password"
		}
		reqLogger.Error("Password reset failed", err, map[string]interface{}{"email": req.Email})
		ctx.JSON(status, authdto.ErrorResponse{Success: false, Message: message})
		return
	}

	ctx.JSON(http.StatusOK, authdto.SuccessResponse{Success: true, Message: "Password reset successful. Please login with your new password."})
	reqLogger.Info("Password reset successful", map[string]interface{}{"email": req.Email})
}

// ────────────────────────────── helpers ──────────────────────────────

// getUserUUID extracts and validates the user_id from the gin context.
// It writes the appropriate error response and returns false on failure.
func getUserUUID(ctx *gin.Context) (uuid.UUID, bool) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, authdto.ErrorResponse{Success: false, Message: "Unauthorized"})
		return uuid.Nil, false
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, authdto.ErrorResponse{Success: false, Message: "Invalid user ID"})
		return uuid.Nil, false
	}
	return userUUID, true
}
