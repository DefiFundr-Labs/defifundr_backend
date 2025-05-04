package repositories

import (
	"context"
	"errors"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/infrastructure/common/utils"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type OTPRepository struct {
	store db.Queries
}

func NewOtpRepository(store db.Queries) *OTPRepository {
	return &OTPRepository{store: store}
}

// CreateOTP creates a new OTP in the database
func (r *OTPRepository) CreateOTP(ctx context.Context, otp domain.OTPVerification) (*domain.OTPVerification, error) {

	hashedOtp := utils.Hash(otp.OTPCode)

	returnedData, err := r.store.CreateOTPVerification(ctx, db.CreateOTPVerificationParams{
		UserID:    pgtype.UUID{Bytes: otp.UserID, Valid: true},
		OtpCode:   otp.OTPCode,
		HashedOtp: hashedOtp,
	})

	if err != nil {
		return nil, err
	}

	return &domain.OTPVerification{
		ID:        returnedData.ID,
		UserID:    returnedData.UserID.Bytes,
		OTPCode:   returnedData.OtpCode,
		HashedOTP: returnedData.HashedOtp,
		Purpose:   domain.OTPPurpose(returnedData.Purpose),
	}, nil

}

// GetOTPByUserIDAndPurpose retrieves an OTP by user ID and purpose
func (r *OTPRepository) GetOTPByUserIDAndPurpose(ctx context.Context, userID uuid.UUID, purpose domain.OTPPurpose) (*domain.OTPVerification, error) {
	otpData, err := r.store.GetOTPVerificationByUserAndPurpose(ctx, db.GetOTPVerificationByUserAndPurposeParams{
		UserID:  pgtype.UUID{Bytes: userID, Valid: true},
		Purpose: db.OtpPurpose(purpose),
	})

	if err != nil {
		return nil, err
	}

	return &domain.OTPVerification{
		ID:        otpData.ID,
		UserID:    otpData.UserID.Bytes,
		OTPCode:   otpData.OtpCode,
		HashedOTP: otpData.HashedOtp,
		Purpose:   domain.OTPPurpose(otpData.Purpose),
	}, nil

}

// VerifyOTP verifies an OTP
func (r *OTPRepository) VerifyOTP(ctx context.Context, id uuid.UUID, code string) error {
	otpData, err := r.store.GetOTPVerificationByID(ctx, id)

	if err != nil {
		return err
	}

	if code != otpData.OtpCode {
		return errors.New("invalid OTP")
	}

	// Verify the OTP
	if otpData.HashedOtp != utils.Hash(otpData.OtpCode) {
		// new error
		return errors.New("invalid OTP")
	}

	// Check if the OTP has expired
	if otpData.ExpiresAt.Before(utils.GetCurrentTime()) {
		return errors.New("OTP has expired")
	}

	// Check if the OTP has been used
	if otpData.MaxAttempts <= otpData.AttemptsMade {
		return errors.New("max attempts reached")
	}

	// Mark the OTP as verified
	_, err = r.store.VerifyOTP(ctx, db.VerifyOTPParams{
		ID:      otpData.ID,
		OtpCode: otpData.OtpCode,
	})

	if err != nil {
		return err
	}

	return nil

}

// IncrementAttempts increments the number of attempts for an OTP
func (r *OTPRepository) IncrementAttempts(ctx context.Context, id uuid.UUID) error {

	_, err := r.store.UpdateOTPAttempts(ctx, id)

	if err != nil {
		return err
	}

	return nil
}
