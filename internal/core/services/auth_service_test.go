package services

import (
	"context"
	"testing"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_RegisterUser(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	mockOTPRepo := new(mocks.FakeOTPRepository)
	mockSessionRepo := new(mocks.FakeSessionRepository)
	mockTokenMaker := new(mocks.FakeMaker)
	cfg := config.Config{
		RefreshTokenDuration: time.Hour * 24,
	}

	service := NewAuthService(mockUserRepo, mockOTPRepo, mockSessionRepo, mockTokenMaker, cfg)

	ctx := context.Background()
	user := domain.User{
		Email:     "test@example.com",
		FirstName: "Test",
		LastName:  "User",
		Password:  nil,
	}
	password := "password123"

	// Mock expectations
	mockUserRepo.GetUserByEmailReturns(nil, nil)
	mockUserRepo.CreateUserReturns(&user, nil)
	mockOTPRepo.CreateOTPReturns(&domain.OTPVerification{}, nil)

	// Act
	result, err := service.RegisterUser(ctx, user, password)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, 1, mockUserRepo.GetUserByEmailCallCount())
	assert.Equal(t, 1, mockUserRepo.CreateUserCallCount())
	assert.Equal(t, 1, mockOTPRepo.CreateOTPCallCount())
}

func TestAuthService_RegisterUser_EmailExists(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	mockOTPRepo := new(mocks.FakeOTPRepository)
	mockSessionRepo := new(mocks.FakeSessionRepository)
	mockTokenMaker := new(mocks.FakeMaker)
	cfg := config.Config{}

	service := NewAuthService(mockUserRepo, mockOTPRepo, mockSessionRepo, mockTokenMaker, cfg)

	ctx := context.Background()
	user := domain.User{
		Email: "existing@example.com",
	}
	password := "password123"

	existingUser := &domain.User{
		ID:    uuid.New(),
		Email: user.Email,
	}

	// Mock expectations
	mockUserRepo.GetUserByEmailReturns(existingUser, nil)

	// Act
	result, err := service.RegisterUser(ctx, user, password)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "user with this email already exists")
	assert.Equal(t, 1, mockUserRepo.GetUserByEmailCallCount())
}

func TestAuthService_Login(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	mockOTPRepo := new(mocks.FakeOTPRepository)
	mockSessionRepo := new(mocks.FakeSessionRepository)
	mockTokenMaker := new(mocks.FakeMaker)
	cfg := config.Config{
		RefreshTokenDuration: time.Hour * 24,
	}

	service := NewAuthService(mockUserRepo, mockOTPRepo, mockSessionRepo, mockTokenMaker, cfg)

	ctx := context.Background()
	email := "test@example.com"
	password := "password123"
	userAgent := "test-agent"
	clientIP := "127.0.0.1"

	user := &domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: &password,
	}

	session := &domain.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(cfg.RefreshTokenDuration),
	}

	// Mock expectations
	mockUserRepo.GetUserByEmailReturns(user, nil)
	mockSessionRepo.CreateSessionReturns(session, nil)

	// Act
	resultSession, resultUser, err := service.Login(ctx, email, password, userAgent, clientIP)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultSession)
	assert.NotNil(t, resultUser)
	assert.Equal(t, user.ID, resultUser.ID)
	assert.Equal(t, session.ID, resultSession.ID)
	assert.Equal(t, 1, mockUserRepo.GetUserByEmailCallCount())
	assert.Equal(t, 1, mockSessionRepo.CreateSessionCallCount())
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	mockOTPRepo := new(mocks.FakeOTPRepository)
	mockSessionRepo := new(mocks.FakeSessionRepository)
	mockTokenMaker := new(mocks.FakeMaker)
	cfg := config.Config{}

	service := NewAuthService(mockUserRepo, mockOTPRepo, mockSessionRepo, mockTokenMaker, cfg)

	ctx := context.Background()
	email := "test@example.com"
	correctPassword := "correctpassword"
	wrongPassword := "wrongpassword"
	userAgent := "test-agent"
	clientIP := "127.0.0.1"

	user := &domain.User{
		ID:       uuid.New(),
		Email:    email,
		Password: &correctPassword,
	}

	// Mock expectations
	mockUserRepo.GetUserByEmailReturns(user, nil)

	// Act
	resultSession, resultUser, err := service.Login(ctx, email, wrongPassword, userAgent, clientIP)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultSession)
	assert.Nil(t, resultUser)
	assert.Contains(t, err.Error(), "invalid credentials")
	assert.Equal(t, 1, mockUserRepo.GetUserByEmailCallCount())
}

func TestAuthService_VerifyEmail(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	mockOTPRepo := new(mocks.FakeOTPRepository)
	mockSessionRepo := new(mocks.FakeSessionRepository)
	mockTokenMaker := new(mocks.FakeMaker)
	cfg := config.Config{}

	service := NewAuthService(mockUserRepo, mockOTPRepo, mockSessionRepo, mockTokenMaker, cfg)

	ctx := context.Background()
	userID := uuid.New()
	code := "123456"

	// Mock expectations
	mockOTPRepo.VerifyOTPReturns(nil)

	// Act
	err := service.VerifyEmail(ctx, userID, code)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, mockOTPRepo.VerifyOTPCallCount())
}

func TestAuthService_RefreshToken(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	mockOTPRepo := new(mocks.FakeOTPRepository)
	mockSessionRepo := new(mocks.FakeSessionRepository)
	mockTokenMaker := new(mocks.FakeMaker)
	cfg := config.Config{
		AccessTokenDuration: time.Hour,
	}

	service := NewAuthService(mockUserRepo, mockOTPRepo, mockSessionRepo, mockTokenMaker, cfg)

	ctx := context.Background()
	refreshToken := "refresh-token"
	userAgent := "test-agent"
	clientIP := "127.0.0.1"

	user := &domain.User{
		ID:    uuid.New(),
		Email: "test@example.com",
	}

	session := &domain.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour),
	}

	// Mock expectations
	mockSessionRepo.GetSessionByRefreshTokenReturns(session, nil)
	mockUserRepo.GetUserByIDReturns(user, nil)
	mockTokenMaker.CreateTokenReturns("access-token", nil, nil)

	// Act
	resultSession, accessToken, err := service.RefreshToken(ctx, refreshToken, userAgent, clientIP)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resultSession)
	assert.NotEmpty(t, accessToken)
	assert.Equal(t, session.ID, resultSession.ID)
	assert.Equal(t, 1, mockSessionRepo.GetSessionByRefreshTokenCallCount())
	assert.Equal(t, 1, mockUserRepo.GetUserByIDCallCount())
	assert.Equal(t, 1, mockTokenMaker.CreateTokenCallCount())
}

func TestAuthService_RefreshToken_Expired(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	mockOTPRepo := new(mocks.FakeOTPRepository)
	mockSessionRepo := new(mocks.FakeSessionRepository)
	mockTokenMaker := new(mocks.FakeMaker)
	cfg := config.Config{}

	service := NewAuthService(mockUserRepo, mockOTPRepo, mockSessionRepo, mockTokenMaker, cfg)

	ctx := context.Background()
	refreshToken := "expired-token"
	userAgent := "test-agent"
	clientIP := "127.0.0.1"

	session := &domain.Session{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(-time.Hour), // Expired
	}

	// Mock expectations
	mockSessionRepo.GetSessionByRefreshTokenReturns(session, nil)

	// Act
	resultSession, accessToken, err := service.RefreshToken(ctx, refreshToken, userAgent, clientIP)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, resultSession)
	assert.Empty(t, accessToken)
	assert.Contains(t, err.Error(), "refresh token has expired")
	assert.Equal(t, 1, mockSessionRepo.GetSessionByRefreshTokenCallCount())
}

func TestAuthService_Logout(t *testing.T) {
	// Arrange
	mockUserRepo := new(mocks.FakeUserRepository)
	mockOTPRepo := new(mocks.FakeOTPRepository)
	mockSessionRepo := new(mocks.FakeSessionRepository)
	mockTokenMaker := new(mocks.FakeMaker)
	cfg := config.Config{}

	service := NewAuthService(mockUserRepo, mockOTPRepo, mockSessionRepo, mockTokenMaker, cfg)

	ctx := context.Background()
	sessionID := uuid.New()

	// Mock expectations
	mockSessionRepo.DeleteSessionReturns(nil)

	// Act
	err := service.Logout(ctx, sessionID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, 1, mockSessionRepo.DeleteSessionCallCount())
}
