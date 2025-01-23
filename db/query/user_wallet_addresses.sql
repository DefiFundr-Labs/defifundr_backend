-- name: CreateUserWallet :one

INSERT INTO user_wallet_addresses (
    user_id,
    wallet_address,
    chain,
    status
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserWallets :many

SELECT * FROM user_wallet_addresses
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetWalletById :one

SELECT * FROM user_wallet_addresses WHERE id = $1 AND user_id = $2 LIMIT 1;

-- name: GetWalletByAddress :one

SELECT * FROM user_wallet_addresses WHERE wallet_address = $1 AND user_id = $2 LIMIT 1;

-- name: UpdateUserWalletStatus :one

UPDATE user_wallet_addresses
SET status = $3, updated_at = now()
WHERE wallet_address = $1 AND user_id = $2
RETURNING *;

-- name: SoftDeleteUserWallet :one

UPDATE user_wallet_addresses
SET deleted_at = now(), updated_at = now(), status = 'deleted'
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: HardDeleteUserWallet :one

DELETE FROM user_wallet_addresses WHERE wallet_address = $1 AND user_id = $2 RETURNING *;

-- name: CheckWalletExists :one

SELECT EXISTS (
    SELECT 1 FROM user_wallet_addresses WHERE wallet_address = $1 AND user_id = $2 LIMIT 1
);
