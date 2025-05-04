-- +goose Up
CREATE TABLE user_wallets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    address TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    chain TEXT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE user_wallets;