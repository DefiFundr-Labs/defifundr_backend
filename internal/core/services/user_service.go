package services

import (
	"context"
	"github.com/demola234/defifundr/pkg/tracing"
	"errors"
	"fmt"
	"time"

	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	utils "github.com/demola234/defifundr/pkg/hash"

	"github.com/google/uuid"
)

type userService struct {
	userRepo ports.UserRepository
}

// UpdateKYC implements ports.UserService.
func (u *userService) UpdateKYC(ctx context.Context, kyc domain.KYC) error {
	ctx, span := tracing.Tracer("user-service").Start(ctx, "UpdateKYC")
	defer span.End()
	panic("unimplemented")
}

// NewUserService creates a new instance of userService.
func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetUserByID implements ports.UserService.
func (u *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	ctx, span := tracing.Tracer("user-service").Start(ctx, "GetUserByID")
	defer span.End()
	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("failed to get user with ID %s: %w", userID, err)
	}

	// Remove sensitive information
	user.Password = nil

	return user, nil
}

// UpdateUser implements ports.UserService.
func (u *userService) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	ctx, span := tracing.Tracer("user-service").Start(ctx, "UpdateUser")
	defer span.End()
	// Verify user exists
	existingUser, err := u.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with ID %s: %w", user.ID, err)
	}

	// Preserve fields that shouldn't be updated through this method
	user.Password = existingUser.Password
	user.CreatedAt = existingUser.CreatedAt
	user.UpdatedAt = time.Now()

	// Update the user
	updatedUser, err := u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Remove sensitive information
	updatedUser.Password = nil

	return updatedUser, nil
}

// UpdatePassword implements ports.UserService.
func (u *userService) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	ctx, span := tracing.Tracer("user-service").Start(ctx, "UpdatePassword")
	defer span.End()
	// Get the user
	user, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user with ID %s: %w", userID, err)
	}

	// Verify old password
	if err := utils.CheckPassword(oldPassword, *user.Password); err != nil {
		return errors.New("incorrect old password")
	}

	// Validate new password (length, complexity, etc.)
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update the password
	err = u.userRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// Helper function to validate password
func validatePassword(password string) error {
	// Implement password validation logic
	// For example: minimum length, required characters, etc.
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Add more validation as needed

	return nil
}

// ResetUserPassword updates password without requiring old password (for password reset flow)
func (u *userService) ResetUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Get the user to verify it exists
	_, err := u.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user with ID %s: %w", userID, err)
	}

	// Validate new password
	if err := validatePassword(newPassword); err != nil {
		return err
	}

	// Hash the new password
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update the password
	err = u.userRepo.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}