-- Create ENUM type for wallet status
CREATE TYPE user_wallet_addresses_statuses AS ENUM ('active', 'inactive', 'suspended', 'deleted');

-- Create wallets table
CREATE TABLE user_wallet_addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR NOT NULL,
    wallet_address VARCHAR UNIQUE NOT NULL,
    chain VARCHAR NOT NULL,
    status user_wallet_addresses_statuses NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users (username) ON DELETE CASCADE
);
