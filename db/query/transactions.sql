-- name: CreateWalletTransaction :one
INSERT INTO wallet_transactions (
  id,
  wallet_address,
  transaction_hash,
  chain_id,
  block_number,
  from_address,
  to_address,
  token_address,
  token_symbol,
  amount,
  transaction_type,
  transaction_status,
  gas_price,
  gas_used,
  transaction_fee,
  reference_type,
  reference_id,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @wallet_address,
  @transaction_hash,
  @chain_id,
  @block_number,
  @from_address,
  @to_address,
  @token_address,
  @token_symbol,
  @amount,
  @transaction_type,
  COALESCE(@transaction_status, 'pending'),
  @gas_price,
  @gas_used,
  @transaction_fee,
  @reference_type,
  @reference_id,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetWalletTransactionByHash :one
SELECT wt.*, sn.name as network_name
FROM wallet_transactions wt
JOIN supported_networks sn ON wt.chain_id = sn.chain_id
WHERE wt.transaction_hash = @transaction_hash;

-- name: GetWalletTransactionsByAddress :many
SELECT wt.*, sn.name as network_name
FROM wallet_transactions wt
JOIN supported_networks sn ON wt.chain_id = sn.chain_id
WHERE wt.wallet_address = @wallet_address
ORDER BY wt.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetWalletTransactionsByNetwork :many
SELECT wt.*, sn.name as network_name
FROM wallet_transactions wt
JOIN supported_networks sn ON wt.chain_id = sn.chain_id
WHERE wt.chain_id = @chain_id
  AND (@wallet_address::text IS NULL OR wt.wallet_address = @wallet_address)
ORDER BY wt.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetUserTransactions :many
SELECT wt.*, sn.name as network_name
FROM wallet_transactions wt
JOIN supported_networks sn ON wt.chain_id = sn.chain_id
JOIN user_wallets uw ON wt.wallet_address = uw.wallet_address AND wt.chain_id = uw.chain_id
WHERE uw.user_id = @user_id
ORDER BY wt.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetCompanyTransactions :many
SELECT wt.*, sn.name as network_name
FROM wallet_transactions wt
JOIN supported_networks sn ON wt.chain_id = sn.chain_id
JOIN company_wallets cw ON wt.wallet_address = cw.wallet_address AND wt.chain_id = cw.chain_id
WHERE cw.company_id = @company_id
ORDER BY wt.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetTransactionsByToken :many
SELECT wt.*, sn.name as network_name
FROM wallet_transactions wt
JOIN supported_networks sn ON wt.chain_id = sn.chain_id
WHERE wt.token_symbol = @token_symbol
  AND (@user_id::uuid IS NULL OR wt.wallet_address IN (
    SELECT wallet_address FROM user_wallets WHERE user_id = @user_id
  ))
ORDER BY wt.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetPendingTransactions :many
SELECT wt.*, sn.name as network_name
FROM wallet_transactions wt
JOIN supported_networks sn ON wt.chain_id = sn.chain_id
WHERE wt.transaction_status = 'pending'
  AND (@user_id::uuid IS NULL OR wt.wallet_address IN (
    SELECT wallet_address FROM user_wallets WHERE user_id = @user_id
  ))
ORDER BY wt.created_at DESC;

-- name: UpdateWalletTransaction :one
UPDATE wallet_transactions SET
  block_number = COALESCE(@block_number, block_number),
  transaction_status = COALESCE(@transaction_status, transaction_status),
  gas_price = COALESCE(@gas_price, gas_price),
  gas_used = COALESCE(@gas_used, gas_used),
  transaction_fee = COALESCE(@transaction_fee, transaction_fee),
  updated_at = NOW()
WHERE transaction_hash = @transaction_hash
RETURNING *;

-- name: CreateFiatTransaction :one
INSERT INTO fiat_transactions (
  id,
  bank_account_id,
  transaction_reference,
  transaction_type,
  amount,
  currency,
  status,
  payment_provider,
  payment_method,
  provider_reference,
  provider_fee,
  reference_type,
  reference_id,
  metadata,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @bank_account_id,
  @transaction_reference,
  @transaction_type,
  @amount,
  @currency,
  COALESCE(@status, 'pending'),
  @payment_provider,
  @payment_method,
  @provider_reference,
  COALESCE(@provider_fee, 0),
  @reference_type,
  @reference_id,
  @metadata,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetFiatTransactionByReference :one
SELECT ft.*, ba.account_holder_name, ba.bank_name
FROM fiat_transactions ft
JOIN bank_accounts ba ON ft.bank_account_id = ba.id
WHERE ft.transaction_reference = @transaction_reference;

-- name: GetFiatTransactionsByBankAccount :many
SELECT ft.*, ba.account_holder_name, ba.bank_name
FROM fiat_transactions ft
JOIN bank_accounts ba ON ft.bank_account_id = ba.id
WHERE ft.bank_account_id = @bank_account_id
ORDER BY ft.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetUserFiatTransactions :many
SELECT ft.*, ba.account_holder_name, ba.bank_name
FROM fiat_transactions ft
JOIN bank_accounts ba ON ft.bank_account_id = ba.id
WHERE ba.user_id = @user_id
ORDER BY ft.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetCompanyFiatTransactions :many
SELECT ft.*, ba.account_holder_name, ba.bank_name
FROM fiat_transactions ft
JOIN bank_accounts ba ON ft.bank_account_id = ba.id
WHERE ba.company_id = @company_id
ORDER BY ft.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: UpdateFiatTransaction :one
UPDATE fiat_transactions SET
  status = COALESCE(@status, status),
  payment_provider = COALESCE(@payment_provider, payment_provider),
  payment_method = COALESCE(@payment_method, payment_method),
  provider_reference = COALESCE(@provider_reference, provider_reference),
  provider_fee = COALESCE(@provider_fee, provider_fee),
  metadata = COALESCE(@metadata, metadata),
  updated_at = NOW()
WHERE transaction_reference = @transaction_reference
RETURNING *;

-- name: CreateExchangeRate :one
INSERT INTO exchange_rates (
  id,
  base_currency,
  quote_currency,
  rate,
  source,
  timestamp,
  created_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @base_currency,
  @quote_currency,
  @rate,
  @source,
  @timestamp,
  COALESCE(@created_at, NOW())
) RETURNING *;

-- name: GetLatestExchangeRate :one
SELECT * FROM exchange_rates 
WHERE base_currency = @base_currency AND quote_currency = @quote_currency
ORDER BY timestamp DESC
LIMIT 1;

-- name: GetExchangeRateHistory :many
SELECT * FROM exchange_rates 
WHERE base_currency = @base_currency 
  AND quote_currency = @quote_currency
  AND timestamp >= @start_time
  AND timestamp <= @end_time
ORDER BY timestamp DESC;

-- name: GetAllLatestExchangeRates :many
SELECT DISTINCT ON (base_currency, quote_currency) *
FROM exchange_rates 
ORDER BY base_currency, quote_currency, timestamp DESC;