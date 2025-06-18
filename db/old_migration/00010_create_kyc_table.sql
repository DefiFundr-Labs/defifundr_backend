-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE kyc (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  face_verification BOOLEAN NOT NULL,
  identity_verification BOOLEAN NOT NULL,
  verification_type VARCHAR(50) NOT NULL,
  verification_number VARCHAR(50) NOT NULL,
  verification_status VARCHAR(50) NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_kyc_user_id ON kyc(user_id);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS kyc;