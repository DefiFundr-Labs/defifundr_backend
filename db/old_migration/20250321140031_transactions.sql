-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tx_hash VARCHAR(255) UNIQUE NOT NULL,
  transaction_pin_hash VARCHAR(255) NOT NULL,
  status VARCHAR(50) NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE UNIQUE INDEX idx_transactions_tx_hash ON transactions(tx_hash);

COMMENT ON COLUMN transactions.transaction_pin_hash IS 'hashed transaction pin';
COMMENT ON COLUMN transactions.status IS 'created, pending, not_found, failed';

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS transactions;