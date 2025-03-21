package services

import (
	"context"

	"github.com/defifundr/defifundr-backend/internal/core/domain"
	"github.com/defifundr/defifundr-backend/internal/core/ports"
	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo   ports.UserRepository
	walletRepo ports.WalletRepository
}

func NewUserService(userRepo ports.UserRepository, walletRepo ports.WalletRepository) ports.UserService {
	return &userService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
	}
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	// Get user wallets
	wallets, err := s.walletRepo.GetWalletsByUserID(ctx, id)
	if err == nil {
		user.Wallets = wallets
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, params domain.UpdateUserParams) (domain.User, error) {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return domain.User{}, domain.ErrUserNotFound
	}

	updatedUser, err := s.userRepo.UpdateUser(ctx, id, params)
	if err != nil {
		return domain.User{}, err
	}

	// Get user wallets
	wallets, err := s.walletRepo.GetWalletsByUserID(ctx, id)
	if err == nil {
		updatedUser.Wallets = wallets
	}

	return updatedUser, nil
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int32) ([]domain.User, error) {
	return s.userRepo.ListUsers(ctx, limit, offset)
}

// Auth service implementation
type authService struct {
	userRepo ports.UserRepository
}

func NewAuthService(userRepo ports.UserRepository) ports.AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) RegisterUser(ctx context.Context, params domain.CreateUserParams) (domain.User, error) {
	// Check if email already exists
	_, err := s.userRepo.GetUserByEmail(ctx, params.Email)
	if err == nil {
		return domain.User{}, domain.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return domain.User{}, err
	}

	// Create user with hashed password
	createParams := domain.CreateUserParams{
		Email:          params.Email,
		Password:       string(hashedPassword),
		Role:           params.Role,
		Name:           params.Name,
		JobRole:        params.JobRole,
		CompanyWebsite: params.CompanyWebsite,
		Country:        params.Country,
		EmploymentType: params.EmploymentType,
	}

	return s.userRepo.CreateUser(ctx, createParams)
}

func (s *authService) LoginUser(ctx context.Context, params domain.LoginUserParams) (string, error) {
	// Implementation for JWT token generation
	// This would include checking password hash and generating a token
	// For brevity, the implementation is simplified
	return "jwt-token", nil
}

func (s *authService) VerifyToken(ctx context.Context, token string) (domain.User, error) {
	// Implementation for token verification
	// For brevity, the implementation is simplified
	return domain.User{}, nil
}
