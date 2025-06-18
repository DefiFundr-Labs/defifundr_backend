-- +goose Up
-- Transaction Management
CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    wallet_address VARCHAR(255) NOT NULL,
    transaction_hash VARCHAR(255) UNIQUE NOT NULL,
    chain_id INTEGER NOT NULL,
    block_number BIGINT,
    from_address VARCHAR(255) NOT NULL,
    to_address VARCHAR(255) NOT NULL,
    token_address VARCHAR(255),
    token_symbol VARCHAR(20),
    amount DECIMAL(36, 18) NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    transaction_status VARCHAR(50) DEFAULT 'pending',
    gas_price DECIMAL(36, 18),
    gas_used BIGINT,
    transaction_fee DECIMAL(36, 18),
    reference_type VARCHAR(50),
    reference_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE fiat_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bank_account_id UUID NOT NULL REFERENCES bank_accounts(id) ON DELETE CASCADE,
    transaction_reference VARCHAR(255) UNIQUE NOT NULL,
    transaction_type VARCHAR(50) NOT NULL,
    amount DECIMAL(18, 6) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    payment_provider VARCHAR(100),
    payment_method VARCHAR(50),
    provider_reference VARCHAR(255),
    provider_fee DECIMAL(18, 6) DEFAULT 0,
    reference_type VARCHAR(50),
    reference_id UUID,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE exchange_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    base_currency VARCHAR(10) NOT NULL,
    quote_currency VARCHAR(10) NOT NULL,
    rate DECIMAL(24, 12) NOT NULL,
    source VARCHAR(100) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(base_currency, quote_currency, timestamp, source)
);

-- Create indexes
CREATE INDEX idx_wallet_transactions_wallet_address ON wallet_transactions(wallet_address);
CREATE INDEX idx_wallet_transactions_transaction_hash ON wallet_transactions(transaction_hash);
CREATE INDEX idx_wallet_transactions_chain_id ON wallet_transactions(chain_id);
CREATE INDEX idx_wallet_transactions_from_address ON wallet_transactions(from_address);
CREATE INDEX idx_wallet_transactions_to_address ON wallet_transactions(to_address);
CREATE INDEX idx_wallet_transactions_transaction_type ON wallet_transactions(transaction_type);
CREATE INDEX idx_wallet_transactions_reference ON wallet_transactions(reference_type, reference_id);
CREATE INDEX idx_wallet_transactions_created_at ON wallet_transactions(created_at);

CREATE INDEX idx_fiat_transactions_bank_account_id ON fiat_transactions(bank_account_id);
CREATE INDEX idx_fiat_transactions_transaction_reference ON fiat_transactions(transaction_reference);
CREATE INDEX idx_fiat_transactions_transaction_type ON fiat_transactions(transaction_type);
CREATE INDEX idx_fiat_transactions_status ON fiat_transactions(status);
CREATE INDEX idx_fiat_transactions_reference ON fiat_transactions(reference_type, reference_id);
CREATE INDEX idx_fiat_transactions_created_at ON fiat_transactions(created_at);

CREATE INDEX idx_exchange_rates_base_currency ON exchange_rates(base_currency);
CREATE INDEX idx_exchange_rates_quote_currency ON exchange_rates(quote_currency);
CREATE INDEX idx_exchange_rates_timestamp ON exchange_rates(timestamp);
CREATE INDEX idx_exchange_rates_base_quote ON exchange_rates(base_currency, quote_currency);

-- +goose Down
DROP TABLE IF EXISTS exchange_rates CASCADE;
DROP TABLE IF EXISTS fiat_transactions CASCADE;
DROP TABLE IF EXISTS wallet_transactions CASCADE;