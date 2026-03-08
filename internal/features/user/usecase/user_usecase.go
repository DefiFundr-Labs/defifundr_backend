package userusecase

import (
	"context"
	"errors"

	commons "github.com/demola234/defifundr/infrastructure/hash"
	userdomain "github.com/demola234/defifundr/internal/features/user/domain"
	userport "github.com/demola234/defifundr/internal/features/user/port"
	appErrors "github.com/demola234/defifundr/pkg/apperrors"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/google/uuid"
)

type userUseCase struct {
	userRepo userport.UserRepository
}

// New creates a new UserService.
func New(userRepo userport.UserRepository) userport.UserService {
	return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) GetUserByID(ctx context.Context, userID uuid.UUID) (*userdomain.User, error) {
	ctx, span := tracing.Tracer("user-usecase").Start(ctx, "GetUserByID")
	defer span.End()

	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return nil, appErrors.NewNotFoundError("user not found")
	}
	return user, nil
}

func (uc *userUseCase) UpdateUser(ctx context.Context, user userdomain.User) (*userdomain.User, error) {
	ctx, span := tracing.Tracer("user-usecase").Start(ctx, "UpdateUser")
	defer span.End()

	updated, err := uc.userRepo.UpdateUserPersonalDetails(ctx, user)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return updated, nil
}

func (uc *userUseCase) UpdatePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	ctx, span := tracing.Tracer("user-usecase").Start(ctx, "UpdatePassword")
	defer span.End()

	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		return appErrors.NewNotFoundError("user not found")
	}

	checked, err := commons.CheckPassword(oldPassword, user.PasswordHash)
	if err != nil || !checked {
		return appErrors.NewUnauthorizedError("invalid current password")
	}

	hashedNew, err := commons.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return uc.userRepo.UpdatePassword(ctx, userID, hashedNew)
}

func (uc *userUseCase) ResetUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	ctx, span := tracing.Tracer("user-usecase").Start(ctx, "ResetUserPassword")
	defer span.End()

	hashedNew, err := commons.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return uc.userRepo.UpdatePassword(ctx, userID, hashedNew)
}

func (uc *userUseCase) UpdateKYC(ctx context.Context, kyc userdomain.KYC) error {
	return errors.New("not implemented")
}
