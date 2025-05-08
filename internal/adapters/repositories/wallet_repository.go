package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	db "github.com/demola234/defifundr/db/sqlc"
	"github.com/demola234/defifundr/internal/core/domain"
	"github.com/google/uuid"
	"github.com/demola234/defifundr/pkg/tracing"
	"github.com/jackc/pgx/v5/pgtype"
)

type WalletRepository struct {
	store db.Queries
}

func NewWalletRepository(store db.Queries) *WalletRepository {
	return &WalletRepository{
		store: store,
	}
}

// CreateWallet creates a new wallet for a user
func (r *WalletRepository) CreateWallet(ctx context.Context, wallet domain.UserWallet) error {
	ctx, span := tracing.Tracer("wallet-repository").Start(ctx, "CreateWallet")
	defer span.End()
	params := db.CreateUserWalletParams{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Address:   wallet.Address,
		Type:      wallet.Type,
		Chain:     wallet.Chain,
		IsDefault: wallet.IsDefault,
		CreatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	_, err := r.store.CreateUserWallet(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create user wallet: %w", err)
	}

	return nil
}

// GetWalletByAddress finds a wallet by its address
func (r *WalletRepository) GetWalletByAddress(ctx context.Context, address string) (*domain.UserWallet, error) {
	ctx, span := tracing.Tracer("wallet-repository").Start(ctx, "GetWalletByAddress")
	defer span.End()
	wallet, err := r.store.GetWalletByAddress(ctx, address)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get wallet by address: %w", err)
	}

	return mapDBWalletToDomain(wallet), nil
}

// GetWalletsByUserID gets all wallets for a user
func (r *WalletRepository) GetWalletsByUserID(ctx context.Context, userID uuid.UUID) ([]domain.UserWallet, error) {
	ctx, span := tracing.Tracer("wallet-repository").Start(ctx, "GetWalletsByUserID")
	defer span.End()
	wallets, err := r.store.GetWalletsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets by user ID: %w", err)
	}

	result := make([]domain.UserWallet, len(wallets))
	for i, wallet := range wallets {
		result[i] = *mapDBWalletToDomain(wallet)
	}

	return result, nil
}

// UpdateWallet updates a wallet
func (r *WalletRepository) UpdateWallet(ctx context.Context, wallet domain.UserWallet) error {
	ctx, span := tracing.Tracer("wallet-repository").Start(ctx, "UpdateWallet")
	defer span.End()
	params := db.UpdateUserWalletParams{
		ID:        wallet.ID,
		IsDefault: wallet.IsDefault,
		UpdatedAt: pgtype.Timestamp{Time: time.Now(), Valid: true},
	}

	_, err := r.store.UpdateUserWallet(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update wallet: %w", err)
	}

	return nil
}

// DeleteWallet deletes a wallet
func (r *WalletRepository) DeleteWallet(ctx context.Context, walletID uuid.UUID) error {
	err := r.store.DeleteUserWallet(ctx, walletID)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}

// Helper to map DB wallet to domain
func mapDBWalletToDomain(wallet db.UserWallets) *domain.UserWallet {
	return &domain.UserWallet{
		ID:        wallet.ID,
		UserID:    wallet.UserID,
		Address:   wallet.Address,
		Type:      wallet.Type,
		Chain:     wallet.Chain,
		IsDefault: wallet.IsDefault,
		CreatedAt: wallet.CreatedAt.Time,
		UpdatedAt: wallet.UpdatedAt.Time,
	}
}
