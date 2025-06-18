-- name: CreateSupportedNetwork :one
INSERT INTO supported_networks (
  id,
  name,
  chain_id,
  network_type,
  currency_symbol,
  block_explorer_url,
  rpc_url,
  is_evm_compatible,
  is_active,
  transaction_speed,
  average_block_time,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @name,
  @chain_id,
  @network_type,
  @currency_symbol,
  @block_explorer_url,
  @rpc_url,
  COALESCE(@is_evm_compatible, FALSE),
  COALESCE(@is_active, TRUE),
  @transaction_speed,
  @average_block_time,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetSupportedNetworks :many
SELECT * FROM supported_networks 
WHERE is_active = TRUE
ORDER BY chain_id;

-- name: GetSupportedNetworkByChainID :one
SELECT * FROM supported_networks 
WHERE chain_id = @chain_id AND is_active = TRUE;

-- name: GetMainnetNetworks :many
SELECT * FROM supported_networks 
WHERE network_type = 'mainnet' AND is_active = TRUE
ORDER BY chain_id;

-- name: GetTestnetNetworks :many
SELECT * FROM supported_networks 
WHERE network_type = 'testnet' AND is_active = TRUE
ORDER BY chain_id;

-- name: UpdateNetworkStatus :one
UPDATE supported_networks SET
  is_active = @is_active,
  updated_at = NOW()
WHERE chain_id = @chain_id
RETURNING *;

-- name: CreateSupportedToken :one
INSERT INTO supported_tokens (
  id,
  network_id,
  name,
  symbol,
  decimals,
  contract_address,
  token_type,
  logo_url,
  is_stablecoin,
  is_active,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @network_id,
  @name,
  @symbol,
  COALESCE(@decimals, 18),
  @contract_address,
  @token_type,
  @logo_url,
  COALESCE(@is_stablecoin, FALSE),
  COALESCE(@is_active, TRUE),
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetTokensByNetwork :many
SELECT st.*, sn.name as network_name, sn.chain_id
FROM supported_tokens st
JOIN supported_networks sn ON st.network_id = sn.id
WHERE sn.chain_id = @chain_id AND st.is_active = TRUE
ORDER BY st.symbol;

-- name: GetStablecoinsByNetwork :many
SELECT st.*, sn.name as network_name, sn.chain_id
FROM supported_tokens st
JOIN supported_networks sn ON st.network_id = sn.id
WHERE sn.chain_id = @chain_id AND st.is_stablecoin = TRUE AND st.is_active = TRUE
ORDER BY st.symbol;

-- name: SearchTokens :many
SELECT st.*, sn.name as network_name, sn.chain_id
FROM supported_tokens st
JOIN supported_networks sn ON st.network_id = sn.id
WHERE (st.name ILIKE '%' || @search_term || '%' OR st.symbol ILIKE '%' || @search_term || '%')
  AND st.is_active = TRUE
ORDER BY st.symbol
LIMIT @limit_val OFFSET @offset_val;

-- name: GetTokenByContract :one
SELECT st.*, sn.name as network_name, sn.chain_id
FROM supported_tokens st
JOIN supported_networks sn ON st.network_id = sn.id
WHERE st.contract_address = @contract_address AND sn.chain_id = @chain_id;

-- name: CreateUserWallet :one
INSERT INTO user_wallets (
  id,
  user_id,
  wallet_address,
  wallet_type,
  chain_id,
  is_default,
  is_verified,
  verification_method,
  verified_at,
  nickname,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @wallet_address,
  @wallet_type,
  @chain_id,
  COALESCE(@is_default, FALSE),
  COALESCE(@is_verified, FALSE),
  @verification_method,
  @verified_at,
  @nickname,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetUserWalletsByUser :many
SELECT uw.*, sn.name as network_name
FROM user_wallets uw
JOIN supported_networks sn ON uw.chain_id = sn.chain_id
WHERE uw.user_id = @user_id
ORDER BY uw.is_default DESC, uw.created_at DESC;

-- name: GetUserWalletsByNetwork :many
SELECT uw.*, sn.name as network_name
FROM user_wallets uw
JOIN supported_networks sn ON uw.chain_id = sn.chain_id
WHERE uw.user_id = @user_id AND uw.chain_id = @chain_id
ORDER BY uw.is_default DESC, uw.created_at DESC;

-- name: UpdateUserWallet :one
UPDATE user_wallets SET
  wallet_type = COALESCE(@wallet_type, wallet_type),
  is_default = COALESCE(@is_default, is_default),
  is_verified = COALESCE(@is_verified, is_verified),
  verification_method = COALESCE(@verification_method, verification_method),
  verified_at = COALESCE(@verified_at, verified_at),
  nickname = COALESCE(@nickname, nickname),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: SetUserWalletAsDefault :exec
UPDATE user_wallets SET
  is_default = CASE WHEN id = @wallet_id THEN TRUE ELSE FALSE END,
  updated_at = NOW()
WHERE user_id = @user_id AND chain_id = @chain_id;

-- name: DeleteUserWallet :exec
DELETE FROM user_wallets WHERE id = @id;