package services

import (
	"context"
	"testing"

	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports/mocks"
	hashutils "github.com/demola234/defifundr/pkg/hash"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetUserByID(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	service := NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	expectedUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	// Mock expectations
	mockUserRepo.GetUserByIDReturns(expectedUser, nil)

	// Act
	result, err := service.GetUserByID(ctx, userID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedUser.ID, result.ID)
	assert.Equal(t, expectedUser.Email, result.Email)
	assert.Equal(t, expectedUser.FirstName, result.FirstName)
	assert.Equal(t, expectedUser.LastName, result.LastName)
	assert.Nil(t, result.Password) // Password should not be returned
	assert.Equal(t, 1, mockUserRepo.GetUserByIDCallCount())
}

func TestUserService_UpdateUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	service := NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	existingUser := &domain.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
	}

	updatedUser := domain.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "Updated",
		LastName:  "Name",
	}

	// Mock expectations
	mockUserRepo.GetUserByIDReturns(existingUser, nil)
	mockUserRepo.UpdateUserReturns(&updatedUser, nil)

	// Act
	result, err := service.UpdateUser(ctx, updatedUser)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedUser.ID, result.ID)
	assert.Equal(t, updatedUser.FirstName, result.FirstName)
	assert.Equal(t, updatedUser.LastName, result.LastName)
	assert.Nil(t, result.Password) // Password should not be returned
	assert.Equal(t, 1, mockUserRepo.GetUserByIDCallCount())
	assert.Equal(t, 1, mockUserRepo.UpdateUserCallCount())
}

func TestUserService_UpdatePassword(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	service := NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	oldPassword := "oldpassword"
	hashedOldPassword, _ := hashutils.HashPassword(oldPassword)
	newPassword := "newpassword"

	existingUser := &domain.User{
		ID:       userID,
		Password: &hashedOldPassword,
	}

	// Mock expectations
	mockUserRepo.GetUserByIDReturns(existingUser, nil)
	mockUserRepo.UpdatePasswordReturns(nil)

	// Act
	err := service.UpdatePassword(ctx, userID, oldPassword, newPassword)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, mockUserRepo.GetUserByIDCallCount())
	assert.Equal(t, 1, mockUserRepo.UpdatePasswordCallCount())
}

func TestUserService_UpdatePassword_InvalidOldPassword(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	service := NewUserService(mockUserRepo)

	ctx := context.Background()
	userID := uuid.New()
	oldPassword := "oldpassword"
	hashedOldPassword, _ := hashutils.HashPassword(oldPassword)
	wrongOldPassword := "wrongpassword"
	newPassword := "newpassword"

	existingUser := &domain.User{
		ID:       userID,
		Password: &hashedOldPassword,
	}

	// Mock expectations
	mockUserRepo.GetUserByIDReturns(existingUser, nil)

	// Act
	err := service.UpdatePassword(ctx, userID, wrongOldPassword, newPassword)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incorrect old password")
	assert.Equal(t, 1, mockUserRepo.GetUserByIDCallCount())
	assert.Equal(t, 0, mockUserRepo.UpdatePasswordCallCount()) // Should not be called
}
