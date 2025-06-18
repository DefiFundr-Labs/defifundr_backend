-- +goose Up
-- Wallet and Payment Infrastructure
CREATE TABLE supported_networks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    chain_id INTEGER UNIQUE NOT NULL,
    network_type VARCHAR(50) NOT NULL,
    currency_symbol VARCHAR(10) NOT NULL,
    block_explorer_url VARCHAR(255),
    rpc_url VARCHAR(255),
    is_evm_compatible BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    transaction_speed VARCHAR(50),
    average_block_time INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE supported_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    network_id UUID NOT NULL REFERENCES supported_networks(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(20) NOT NULL,
    decimals INTEGER NOT NULL DEFAULT 18,
    contract_address VARCHAR(255),
    token_type VARCHAR(50) NOT NULL,
    logo_url VARCHAR(255),
    is_stablecoin BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(network_id, symbol, contract_address)
);

CREATE TABLE user_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    wallet_address VARCHAR(255) NOT NULL,
    wallet_type VARCHAR(50) NOT NULL,
    chain_id INTEGER NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_method VARCHAR(50),
    verified_at TIMESTAMPTZ,
    nickname VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, wallet_address, chain_id)
);

CREATE TABLE company_wallets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    wallet_address VARCHAR(255) NOT NULL,
    wallet_type VARCHAR(50) NOT NULL,
    chain_id INTEGER NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    multisig_config JSONB,
    required_approvals INTEGER DEFAULT 1,
    wallet_name VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, wallet_address, chain_id)
);

CREATE TABLE bank_accounts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    account_number VARCHAR(255) NOT NULL,
    account_holder_name VARCHAR(255) NOT NULL,
    bank_name VARCHAR(255) NOT NULL,
    bank_code VARCHAR(100),
    routing_number VARCHAR(100),
    swift_code VARCHAR(50),
    iban VARCHAR(100),
    account_type VARCHAR(50) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    country VARCHAR(100) NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    is_verified BOOLEAN DEFAULT FALSE,
    verification_method VARCHAR(50),
    verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_bank_account_owner CHECK (
        (user_id IS NOT NULL AND company_id IS NULL) OR
        (user_id IS NULL AND company_id IS NOT NULL)
    )
);

-- Create indexes
CREATE INDEX idx_supported_networks_chain_id ON supported_networks(chain_id);
CREATE INDEX idx_supported_networks_is_active ON supported_networks(is_active);
CREATE INDEX idx_supported_tokens_network_id ON supported_tokens(network_id);
CREATE INDEX idx_supported_tokens_symbol ON supported_tokens(symbol);
CREATE INDEX idx_user_wallets_user_id ON user_wallets(user_id);
CREATE INDEX idx_user_wallets_wallet_address ON user_wallets(wallet_address);
CREATE INDEX idx_user_wallets_chain_id ON user_wallets(chain_id);
CREATE INDEX idx_company_wallets_company_id ON company_wallets(company_id);
CREATE INDEX idx_company_wallets_wallet_address ON company_wallets(wallet_address);
CREATE INDEX idx_bank_accounts_user_id ON bank_accounts(user_id);
CREATE INDEX idx_bank_accounts_company_id ON bank_accounts(company_id);
CREATE INDEX idx_bank_accounts_currency ON bank_accounts(currency);

-- +goose Down
DROP TABLE IF EXISTS bank_accounts CASCADE;
DROP TABLE IF EXISTS company_wallets CASCADE;
DROP TABLE IF EXISTS user_wallets CASCADE;
DROP TABLE IF EXISTS supported_tokens CASCADE;
DROP TABLE IF EXISTS supported_networks CASCADE;