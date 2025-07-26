package handlers

import (
	"net/http"
	"github.com/demola234/defifundr/pkg/tracing"

	"github.com/demola234/defifundr/internal/adapters/dto/request"
	"github.com/demola234/defifundr/infrastructure/common/logging"
	"github.com/demola234/defifundr/internal/adapters/dto/response"
	// "github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/demola234/defifundr/pkg/app_errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	userService ports.UserService
	logger      logging.Logger

}

// NewUserHandler creates a new user handler
func NewUserHandler(userService ports.UserService, logger logging.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Retrieve authenticated user's profile
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.UserResponse "User profile"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	// Start OTel span for this handler
	spanCtx, span := tracing.Tracer("user-handler").Start(ctx.Request.Context(), "GetProfile")
	defer span.End()

	// Use spanCtx for downstream calls
	ctxWithSpan := ctx.Copy()
	ctxWithSpan.Request = ctx.Request.WithContext(spanCtx)

	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Message: "Unauthorized",
		})
		return
	}

	// Convert user ID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: "Invalid user ID",
		})
		return
	}

	// Get user profile
	user, err := h.userService.GetUserByID(ctxWithSpan, userUUID)
	if err != nil {
		span.RecordError(err)

		errResponse := response.ErrorResponse{
			Message: appErrors.ErrInternalServer.Error(),
		}

		if appErrors.IsAppError(err) {
			appErr := err.(*appErrors.AppError)
			errResponse.Message = appErr.Error()

			if appErr.ErrorType == appErrors.ErrorTypeNotFound {
				ctx.JSON(http.StatusNotFound, errResponse)
				return
			}

			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	// Create response DTO
	userResponse := response.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "User profile retrieved",
		Data:    userResponse,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update authenticated user's profile information
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param profile body request.UpdateProfileRequest true "Profile data to update"
// @Success 200 {object} response.UserResponse "Updated user profile"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /users/profile [put]
// func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
// 	spanCtx, span := tracing.Tracer("user-handler").Start(ctx.Request.Context(), "UpdateProfile")
// 	defer span.End()
// 	ctxWithSpan := ctx.Copy()
// 	ctxWithSpan.Request = ctx.Request.WithContext(spanCtx)
// 	// Get user ID from context (set by auth middleware)
// 	userID, exists := ctx.Get("user_id")
// 	if !exists {
// 		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
// 			Message: "Unauthorized",
// 		})
// 		return
// 	}

// 	// Convert user ID to UUID
// 	userUUID, ok := userID.(uuid.UUID)
// 	if !ok {
// 		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
// 			Message: "Invalid user ID",
// 		})
// 		return
// 	}

// 	var req request.UpdateProfileRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
// 			Message: appErrors.ErrInvalidRequest.Error(),
// 			Success: false,
// 		})
// 		return
// 	}

// 	// Validate request data
// 	if err := req.Validate(); err != nil {
// 		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
// 			Message: appErrors.ErrInvalidRequest.Error(),
// 			Success: false,
// 		})
// 		return
// 	}

// 	// Get existing user
// 	currentUser, err := h.userService.GetUserByID(ctxWithSpan, userUUID)
// 	if err != nil {
// 		span.RecordError(err)
// 		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
// 			Message: "Failed to retrieve user profile",
// 		})
// 		return
// 	}

// 	// Update user object with new values
// 	updatedUser := domain.User{
// 		ID:                  userUUID,
// 		Email:               currentUser.Email, // Email cannot be changed
// 		FirstName:           req.FirstName,
// 		LastName:            req.LastName,
// 		AccountType:         currentUser.AccountType, // Account type cannot be changed
// 		PersonalAccountType: currentUser.PersonalAccountType,
// 		Nationality:         req.Nationality,
// 		Gender:              &req.Gender,
// 		ResidentialCountry:  &req.ResidentialCountry,
// 		JobRole:             &req.JobRole,
// 		EmploymentType:      &req.EmploymentType,
// 	}

// 	// Update user profile
// 	user, err := h.userService.UpdateUser(ctxWithSpan, updatedUser)
// 	if err != nil {
// 		span.RecordError(err)
// 		errResponse := response.ErrorResponse{
// 			Message: appErrors.ErrInternalServer.Error(),
// 		}

// 		if appErrors.IsAppError(err) {
// 			appErr := err.(*appErrors.AppError)
// 			errResponse.Message = appErr.Error()
// 			ctx.JSON(http.StatusBadRequest, errResponse)
// 			return
// 		}

// 		ctx.JSON(http.StatusInternalServerError, errResponse)
// 		return
// 	}

// 	// Create response DTO
// 	userResponse := response.UserResponse{
// 		ID:        user.ID.String(),
// 		Email:     user.Email,
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 	}

// 	ctx.JSON(http.StatusOK, response.SuccessResponse{
// 		Message: "User profile updated",
// 		Data:    userResponse,
// 	})
// }

// ChangePassword godoc
// @Summary Change user password
// @Description Change authenticated user's password
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param password body request.ChangePasswordRequest true "Password change data"
// @Success 200 {object} response.SuccessResponse "Password changed successfully"
// @Failure 400 {object} response.ErrorResponse "Invalid request"
// @Failure 401 {object} response.ErrorResponse "Unauthorized"
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("user-handler").Start(ctx.Request.Context(), "ChangePassword")
	defer span.End()
	ctxWithSpan := ctx.Copy()
	ctxWithSpan.Request = ctx.Request.WithContext(spanCtx)
	// Get user ID from context (set by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Message: "Unauthorized",
		})
		return
	}

	// Convert user ID to UUID
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Message: "Invalid user ID",
		})
		return
	}

	var req request.UpdateUserPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{
			Message: appErrors.ErrInvalidRequest.Error(),
			Success: false,
		})
		return
	}

	// Change password
	err := h.userService.UpdatePassword(ctxWithSpan, userUUID, req.OldPassword, req.NewPassword)
	if err != nil {
		span.RecordError(err)
		errResponse := response.ErrorResponse{
			Message: appErrors.ErrInternalServer.Error(),
		}

		if appErrors.IsAppError(err) {
			appErr := err.(*appErrors.AppError)
			errResponse.Message = appErr.Error()

			if appErr.ErrorType == appErrors.ErrorTypeUnauthorized {
				ctx.JSON(http.StatusUnauthorized, errResponse)
				return
			}

			ctx.JSON(http.StatusBadRequest, errResponse)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Message: "Password changed successfully",
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
// @Router /users/employee/profile/complete [put]
func (h *UserHandler) UpdateAddressDetails(ctx *gin.Context) {
	spanCtx, span := tracing.Tracer("auth-handler").Start(ctx.Request.Context(), "UpdateAddressDetails")
	defer span.End()
	ctxWithSpan := ctx.Copy()
	ctxWithSpan.Request = ctx.Request.WithContext(spanCtx)
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
	currentUser, err := h.userService.GetUserByID(ctx, userUUID)
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

	currentUser.PhoneNumber = req.PhoneNumber;


	// Update user with new details
	personalUserDetails, err := h.userService.GetPersonalUserByUserID(ctx, userUUID);

	personalUserDetails.UserAddress = &req.UserAddress
	personalUserDetails.UserCity = &req.City
	personalUserDetails.UserPostalCode = &req.PostalCode
	personalUserDetails.ResidentialCountry = &req.Country
	personalUserDetails.Nationality = &req.Nationality
	personalUserDetails.Gender = &req.Gender
	personalUserDetails.DateOfBirth = &req.DateOfBirth

	// Update user
	updatedUser, err := h.userService.UpdateUser(ctx, *currentUser)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to update phone number details"

		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		}

		reqLogger.Error("Failed to update phone number details", err, map[string]interface{}{
			"user_id": userUUID,
		})

		ctx.JSON(status, response.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}

	updatePersonalDetails, err := h.userService.UpdatePersonalUser(ctx, *personalUserDetails);
	if err != nil {
		status := http.StatusInternalServerError
		message := "Failed to update personal user details"

		if appErr, ok := err.(appErrors.AppError); ok {
			status = appErr.StatusCode()
			message = appErr.Error()
		}

		reqLogger.Error("Failed to update personal user details", err, map[string]interface{}{
			"user_id": userUUID,
		})

		ctx.JSON(status, response.ErrorResponse{
			Success: false,
			Message: message,
		})
		return
	}


	// Return success
	ctx.JSON(http.StatusOK, response.SuccessResponse{
		Success: true,
		Message: "Address details updated successfully",
		Data: map[string]interface{}{
			"user":               updatedUser,
			"personalDetails":		  updatePersonalDetails,
		},
	})

	reqLogger.Info("Address details updated successfully", map[string]interface{}{
		"user_id": updatedUser.ID,
	})
}