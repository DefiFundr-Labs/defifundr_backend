package userhandler

import (
	"net/http"

	userdomain "github.com/demola234/defifundr/internal/features/user/domain"
	userport "github.com/demola234/defifundr/internal/features/user/port"
	userdto "github.com/demola234/defifundr/internal/features/user/dto"
	appErrors "github.com/demola234/defifundr/pkg/apperrors"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Handler handles user HTTP requests.
type Handler struct {
	service userport.UserService
}

// New creates a new user Handler.
func New(service userport.UserService) *Handler {
	return &Handler{service: service}
}

// GetProfile returns the authenticated user's profile.
func (h *Handler) GetProfile(ctx *gin.Context) {
	reqCtx, span := tracing.Tracer("user-handler").Start(ctx.Request.Context(), "GetProfile")
	defer span.End()

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, userdto.ErrorResponse{Message: "Unauthorized"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, userdto.ErrorResponse{Message: "Invalid user ID"})
		return
	}

	user, err := h.service.GetUserByID(reqCtx, userUUID)
	if err != nil {
		span.RecordError(err)
		if appErrors.IsAppError(err) {
			appErr := err.(*appErrors.AppError)
			if appErr.ErrorType == appErrors.ErrorTypeNotFound {
				ctx.JSON(http.StatusNotFound, userdto.ErrorResponse{Message: appErr.Error()})
				return
			}
			ctx.JSON(http.StatusBadRequest, userdto.ErrorResponse{Message: appErr.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, userdto.ErrorResponse{Message: appErrors.ErrInternalServer.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userdto.SuccessResponse{
		Success: true,
		Message: "User profile retrieved",
		Data: userdto.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
		},
	})
}

// UpdateProfile updates the authenticated user's profile.
func (h *Handler) UpdateProfile(ctx *gin.Context) {
	reqCtx, span := tracing.Tracer("user-handler").Start(ctx.Request.Context(), "UpdateProfile")
	defer span.End()

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, userdto.ErrorResponse{Message: "Unauthorized"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, userdto.ErrorResponse{Message: "Invalid user ID"})
		return
	}

	var req userdto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, userdto.ErrorResponse{Message: appErrors.ErrInvalidRequest.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, userdto.ErrorResponse{Message: err.Error()})
		return
	}

	currentUser, err := h.service.GetUserByID(reqCtx, userUUID)
	if err != nil {
		span.RecordError(err)
		ctx.JSON(http.StatusInternalServerError, userdto.ErrorResponse{Message: "Failed to retrieve user profile"})
		return
	}

	updatedUser := userdomain.User{
		ID:                  userUUID,
		Email:               currentUser.Email,
		AccountType:         currentUser.AccountType,
		PersonalAccountType: currentUser.PersonalAccountType,
		FirstName:           req.FirstName,
		LastName:            req.LastName,
		Nationality:         req.Nationality,
		Gender:              &req.Gender,
		ResidentialCountry:  &req.ResidentialCountry,
		JobRole:             &req.JobRole,
		EmploymentType:      &req.EmploymentType,
	}

	user, err := h.service.UpdateUser(reqCtx, updatedUser)
	if err != nil {
		span.RecordError(err)
		if appErrors.IsAppError(err) {
			ctx.JSON(http.StatusBadRequest, userdto.ErrorResponse{Message: err.(*appErrors.AppError).Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, userdto.ErrorResponse{Message: appErrors.ErrInternalServer.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userdto.SuccessResponse{
		Success: true,
		Message: "User profile updated",
		Data: userdto.UserResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

// ChangePassword changes the authenticated user's password.
func (h *Handler) ChangePassword(ctx *gin.Context) {
	reqCtx, span := tracing.Tracer("user-handler").Start(ctx.Request.Context(), "ChangePassword")
	defer span.End()

	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, userdto.ErrorResponse{Message: "Unauthorized"})
		return
	}
	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, userdto.ErrorResponse{Message: "Invalid user ID"})
		return
	}

	var req userdto.UpdateUserPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, userdto.ErrorResponse{Message: appErrors.ErrInvalidRequest.Error()})
		return
	}
	if err := req.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, userdto.ErrorResponse{Message: err.Error()})
		return
	}

	err := h.service.UpdatePassword(reqCtx, userUUID, req.OldPassword, req.NewPassword)
	if err != nil {
		span.RecordError(err)
		if appErrors.IsAppError(err) {
			appErr := err.(*appErrors.AppError)
			if appErr.ErrorType == appErrors.ErrorTypeUnauthorized {
				ctx.JSON(http.StatusUnauthorized, userdto.ErrorResponse{Message: appErr.Error()})
				return
			}
			ctx.JSON(http.StatusBadRequest, userdto.ErrorResponse{Message: appErr.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, userdto.ErrorResponse{Message: appErrors.ErrInternalServer.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userdto.SuccessResponse{Success: true, Message: "Password changed successfully"})
}
