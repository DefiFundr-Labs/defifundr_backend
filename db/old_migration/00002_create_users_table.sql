-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash VARCHAR(255),
  profile_picture VARCHAR(255) DEFAULT '',
  account_type VARCHAR(50) NOT NULL, -- business, personal
  gender VARCHAR(50) DEFAULT '',
  personal_account_type VARCHAR(50) NOT NULL, -- contractor, freelancer, employee
  phone_number VARCHAR(50) DEFAULT '',
  phone_number_verified BOOLEAN DEFAULT false,
  phone_number_verified_at TIMESTAMPTZ,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  nationality VARCHAR(255) NOT NULL,
  residential_country VARCHAR(255),
  job_role VARCHAR(255),
  company_name VARCHAR(255) DEFAULT '',
  company_size VARCHAR(255) DEFAULT '',
  company_industry VARCHAR(255) DEFAULT '',
  company_description VARCHAR(255) DEFAULT '',
  company_headquarters VARCHAR(255) DEFAULT '',
  user_address VARCHAR(255) DEFAULT '',
  user_city VARCHAR(255) DEFAULT '',
  user_postal_code VARCHAR(255) DEFAULT '',
  employee_type VARCHAR(255) DEFAULT '',
  auth_provider VARCHAR(255),
  provider_id VARCHAR(255) NOT NULL,
  company_website VARCHAR(255),
  employment_type VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_users_email ON users(email);

COMMENT ON COLUMN users.account_type IS 'business, personal';
COMMENT ON COLUMN users.personal_account_type IS 'contractor, freelancer, employee';

-- +goose Down
DROP TABLE IF EXISTS users;
