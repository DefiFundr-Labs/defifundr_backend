-- +goose Up
-- Create UUID extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Core User and Account Management
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone_number VARCHAR(50),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255),
    profile_picture_url VARCHAR(255),
    auth_provider VARCHAR(50),
    provider_id VARCHAR(255),
    email_verified BOOLEAN DEFAULT FALSE,
    email_verified_at TIMESTAMPTZ,
    phone_number_verified BOOLEAN DEFAULT FALSE,
    phone_number_verified_at TIMESTAMPTZ,
    account_type VARCHAR(50) NOT NULL,
    account_status VARCHAR(50) DEFAULT 'pending',
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    two_factor_method VARCHAR(50),
    user_login_type VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    last_login_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE personal_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nationality VARCHAR(255),
    residential_country VARCHAR(255),
    user_address VARCHAR(255),
    user_city VARCHAR(255),
    user_postal_code VARCHAR(255),
    gender VARCHAR(50),
    date_of_birth DATE,
    job_role VARCHAR(255),
    personal_account_type VARCHAR(50),
    employment_type VARCHAR(50),
    tax_id VARCHAR(255),
    default_payment_currency VARCHAR(50),
    default_payment_method VARCHAR(50),
    hourly_rate DECIMAL(18, 2),
    specialization VARCHAR(255),
    kyc_status VARCHAR(50) DEFAULT 'pending',
    kyc_verified_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_name VARCHAR(255) NOT NULL,
    company_email VARCHAR(255),
    company_phone VARCHAR(50),
    company_size VARCHAR(50),
    company_industry VARCHAR(255),
    company_description TEXT,
    company_headquarters VARCHAR(255),
    company_logo VARCHAR(255),
    company_website VARCHAR(255),
    primary_contact_name VARCHAR(255),
    primary_contact_email VARCHAR(255),
    primary_contact_phone VARCHAR(50),
    company_address VARCHAR(255),
    company_city VARCHAR(255),
    company_postal_code VARCHAR(255),
    company_country VARCHAR(255),
    company_registration_number VARCHAR(255),
    registration_country VARCHAR(255),
    tax_id VARCHAR(255),
    incorporation_date DATE,
    account_status VARCHAR(50) DEFAULT 'pending',
    kyb_status VARCHAR(50) DEFAULT 'pending',
    kyb_verified_at TIMESTAMPTZ,
    kyb_verification_method VARCHAR(50),
    kyb_verification_provider VARCHAR(255),
    kyb_rejection_reason TEXT,
    legal_entity_type VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE company_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    department VARCHAR(100),
    job_title VARCHAR(255),
    is_administrator BOOLEAN DEFAULT FALSE,
    can_manage_payroll BOOLEAN DEFAULT FALSE,
    can_manage_invoices BOOLEAN DEFAULT FALSE,
    can_manage_employees BOOLEAN DEFAULT FALSE,
    can_manage_company_settings BOOLEAN DEFAULT FALSE,
    can_manage_bank_accounts BOOLEAN DEFAULT FALSE,
    can_manage_wallets BOOLEAN DEFAULT FALSE,
    permissions JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    added_by UUID REFERENCES users(id),
    reports_to UUID REFERENCES company_users(id),
    hire_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, user_id)
);

CREATE TABLE company_employees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    employee_id VARCHAR(100),
    department VARCHAR(100),
    position VARCHAR(100),
    employment_status VARCHAR(50) DEFAULT 'active',
    employment_type VARCHAR(50),
    start_date DATE,
    end_date DATE,
    manager_id UUID REFERENCES company_users(id),
    salary_amount DECIMAL(18, 6),
    salary_currency VARCHAR(10),
    salary_frequency VARCHAR(50),
    hourly_rate DECIMAL(18, 6),
    payment_method VARCHAR(50),
    payment_split JSONB,
    tax_information JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, employee_id)
);

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_account_type ON users(account_type);
CREATE INDEX idx_users_account_status ON users(account_status);
CREATE INDEX idx_companies_owner_id ON companies(owner_id);
CREATE INDEX idx_company_users_company_id ON company_users(company_id);
CREATE INDEX idx_company_users_user_id ON company_users(user_id);
CREATE INDEX idx_company_employees_company_id ON company_employees(company_id);
CREATE INDEX idx_company_employees_user_id ON company_employees(user_id);

-- +goose Down
DROP TABLE IF EXISTS company_employees CASCADE;
DROP TABLE IF EXISTS company_users CASCADE;
DROP TABLE IF EXISTS companies CASCADE;
DROP TABLE IF EXISTS personal_users CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP EXTENSION IF EXISTS "uuid-ossp";