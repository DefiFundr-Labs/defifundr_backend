-- name: CreateUserWallet :one
INSERT INTO user_wallets (
    id, user_id, address, type, chain, is_default, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetWalletByAddress :one
SELECT * FROM user_wallets WHERE address = $1;

-- name: GetWalletsByUserID :many
SELECT * FROM user_wallets WHERE user_id = $1;

-- name: UpdateUserWallet :one
UPDATE user_wallets
SET is_default = $2, updated_at = $3
WHERE id = $1
RETURNING *;

-- name: DeleteUserWallet :exec
DELETE FROM user_wallets WHERE id = $1;
