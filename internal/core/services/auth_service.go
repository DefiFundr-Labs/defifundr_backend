package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/demola234/defifundr/config"
	"github.com/demola234/defifundr/infrastructure/common/utils"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/demola234/defifundr/internal/core/ports"
	"github.com/demola234/defifundr/pkg/random"
	tokenMaker "github.com/demola234/defifundr/pkg/token_maker"
	"github.com/google/uuid"
)

type authService struct {
	userRepo    ports.UserRepository
	otpRepo     ports.OTPRepository
	sessionRepo ports.SessionRepository
	tokenMaker  tokenMaker.Maker
	config      config.Config
}

// NewAuthService creates a new instance of authService.
func NewAuthService(
	userRepo ports.UserRepository,
	otpRepo ports.OTPRepository,
	sessionRepo ports.SessionRepository,
	tokenMaker tokenMaker.Maker,
	config config.Config,
) ports.AuthService {
	return &authService{
		userRepo:    userRepo,
		otpRepo:     otpRepo,
		sessionRepo: sessionRepo,
		tokenMaker:  tokenMaker,
		config:      config,
	}
}

// GenerateOTP implements ports.AuthService.
func (a *authService) GenerateOTP(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose, contactMethod string) (*domain.OTPVerification, error) {
	otpCode := random.RandomOtp()
	hashedOtp := utils.Hash(otpCode)

	otp, err := a.otpRepo.CreateOTP(ctx, domain.OTPVerification{
		ID:            uuid.New(),
		UserID:        userID,
		OTPCode:       otpCode,
		HashedOTP:     hashedOtp,
		Purpose:       purpose,
		ContactMethod: contactMethod,
		ExpiresAt:     time.Now().Add(time.Minute * 5),
		IsVerified:    false,
	})

	if err != nil {
		return nil, err
	}

	// Note: Commented out as the emailService might not be initialized yet
	// a.emailService.SendVerificationEmail(ctx, user.Email, user.FirstName, otp.OTPCode)

	return otp, nil
}

// Login implements ports.AuthService.
func (a *authService) Login(ctx context.Context, email string, password string, userAgent string, clientIP string) (*domain.Session, *domain.User, error) {
	// Check if user with the provided email exists
	user, err := a.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user with email %s: %w", email, err)
	}

	if user == nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// Verify password - replace with your password verification function
	if !verifyPassword(password, *user.Password) {
		return nil, nil, errors.New("invalid credentials")
	}

	// Create a new session
	sessionID := domain.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		UserAgent:    userAgent,
		ClientIP:     clientIP,
		RefreshToken: random.RandomString(32),
		ExpiresAt:    time.Now().Add(a.config.RefreshTokenDuration),
		CreatedAt:    time.Now(),
	}
	session, err := a.sessionRepo.CreateSession(ctx, sessionID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, user, nil
}

// Logout implements ports.AuthService.
func (a *authService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	return a.sessionRepo.DeleteSession(ctx, sessionID)
}

// RefreshToken implements ports.AuthService.
func (a *authService) RefreshToken(ctx context.Context, refreshToken string, userAgent string, clientIP string) (*domain.Session, string, error) {
	// Get session by refresh token
	session, err := a.sessionRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get session by refresh token: %w", err)
	}

	// Check if session exists
	if session == nil {
		return nil, "", errors.New("invalid refresh token")
	}

	// Check if session is blocked
	if session.IsBlocked {
		return nil, "", errors.New("session is blocked")
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		return nil, "", errors.New("refresh token has expired")
	}

	// Get the user
	user, err := a.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	// Generate a new access token
	accessToken, _, err := a.tokenMaker.CreateToken(
		user.Email,
		user.ID.String(),
		a.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create access token: %w", err)
	}

	// Optional: Refresh the session with a new refresh token
	// For security reasons, many implementations rotate refresh tokens
	if a.config.RotateRefreshTokens {
		// Delete the old session
		err = a.sessionRepo.DeleteSession(ctx, session.ID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to delete old session: %w", err)
		}

		// Create a new session

		newSessionID := domain.Session{
			ID:           uuid.New(),
			UserID:       user.ID,
			RefreshToken: random.RandomString(32),
			UserAgent:    userAgent,
			ClientIP:     clientIP,
			IsBlocked:    false,
			ExpiresAt:    time.Now().Add(a.config.RefreshTokenDuration),
			CreatedAt:    time.Now(),
		}

		session, err = a.sessionRepo.CreateSession(ctx, newSessionID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to create new session: %w", err)
		}
	}

	return session, accessToken, nil
}

// RegisterUser implements ports.AuthService.
func (a *authService) RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error) {
	// Check if user with the same email already exists
	existingUser, err := a.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create the user
	userDetails, err := a.userRepo.CreateUser(ctx, domain.User{
		ID:                  uuid.New(),
		Email:               user.Email,
		Password:            &password,
		ProfilePicture:      user.ProfilePicture,
		AccountType:         user.AccountType,
		Gender:              user.Gender,
		PersonalAccountType: user.PersonalAccountType,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Nationality:         user.Nationality,
		ResidentialCountry:  user.ResidentialCountry,
		JobRole:             user.JobRole,
		CompanyWebsite:      user.CompanyWebsite,
		EmploymentType:      user.EmploymentType,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate OTP for email verification
	_, err = a.GenerateOTP(ctx, userDetails.ID, domain.OTPPurposeEmailVerification, userDetails.Email)
	if err != nil {
		// We still return the user even if OTP generation fails
		return userDetails, nil
	}

	// Note: Commented out as the emailService might not be initialized yet
	// a.emailService.SendVerificationEmail(ctx, user.Email, user.FirstName, otp.OTPCode)

	return userDetails, nil
}

// VerifyEmail implements ports.AuthService.
func (a *authService) VerifyEmail(ctx context.Context, userID uuid.UUID, code string) error {
	return a.VerifyOTP(ctx, userID, domain.OTPPurposeEmailVerification, code)
}

// VerifyOTP implements ports.AuthService.
func (a *authService) VerifyOTP(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose, code string) error {
	return a.otpRepo.VerifyOTP(ctx, userID, code)
}

// Helper function to verify password
func verifyPassword(plainPassword, hashedPassword string) bool {
	// Replace with your actual password verification logic
	// This might use bcrypt or another hashing algorithm
	return plainPassword == hashedPassword // This is just a placeholder!
}
