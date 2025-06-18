-- +goose Up
-- KYC/KYB Management
CREATE TABLE kyc_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    country_id UUID NOT NULL REFERENCES supported_countries(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    document_number VARCHAR(100),
    document_country VARCHAR(100),
    issue_date DATE,
    expiry_date DATE,
    document_url VARCHAR(255),
    ipfs_hash VARCHAR(255),
    verification_status VARCHAR(50) DEFAULT 'pending',
    verification_level VARCHAR(50),
    verification_notes TEXT,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    rejection_reason TEXT,
    metadata JSONB,
    meets_requirements BOOLEAN DEFAULT FALSE,
    requirement_id UUID REFERENCES kyc_country_requirements(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE kyb_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    country_id UUID NOT NULL REFERENCES supported_countries(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    document_number VARCHAR(100),
    document_country VARCHAR(100),
    issue_date DATE,
    expiry_date DATE,
    document_url VARCHAR(255),
    ipfs_hash VARCHAR(255),
    verification_status VARCHAR(50) DEFAULT 'pending',
    verification_level VARCHAR(50),
    verification_notes TEXT,
    verified_by UUID REFERENCES users(id),
    verified_at TIMESTAMPTZ,
    rejection_reason TEXT,
    metadata JSONB,
    meets_requirements BOOLEAN DEFAULT FALSE,
    requirement_id UUID REFERENCES kyb_country_requirements(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE kyc_verification_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    verification_provider VARCHAR(100) NOT NULL,
    verification_reference VARCHAR(255),
    verification_status VARCHAR(50) DEFAULT 'pending',
    verification_result VARCHAR(50),
    response_data JSONB,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE kyb_verification_attempts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    verification_provider VARCHAR(100) NOT NULL,
    verification_reference VARCHAR(255),
    verification_status VARCHAR(50) DEFAULT 'pending',
    verification_result VARCHAR(50),
    response_data JSONB,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_country_kyc_status (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    country_id UUID NOT NULL REFERENCES supported_countries(id) ON DELETE CASCADE,
    verification_status VARCHAR(50) DEFAULT 'pending',
    verification_level VARCHAR(50),
    verification_date TIMESTAMPTZ,
    expiry_date TIMESTAMPTZ,
    rejection_reason TEXT,
    notes TEXT,
    risk_rating VARCHAR(20),
    restricted_features JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, country_id)
);

CREATE TABLE company_country_kyb_status (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    country_id UUID NOT NULL REFERENCES supported_countries(id) ON DELETE CASCADE,
    verification_status VARCHAR(50) DEFAULT 'pending',
    verification_level VARCHAR(50),
    verification_date TIMESTAMPTZ,
    expiry_date TIMESTAMPTZ,
    rejection_reason TEXT,
    notes TEXT,
    risk_rating VARCHAR(20),
    restricted_features JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(company_id, country_id)
);

-- Create indexes
CREATE INDEX idx_kyc_documents_user_id ON kyc_documents(user_id);
CREATE INDEX idx_kyc_documents_country_id ON kyc_documents(country_id);
CREATE INDEX idx_kyc_documents_verification_status ON kyc_documents(verification_status);
CREATE INDEX idx_kyb_documents_company_id ON kyb_documents(company_id);
CREATE INDEX idx_kyb_documents_country_id ON kyb_documents(country_id);
CREATE INDEX idx_kyb_documents_verification_status ON kyb_documents(verification_status);
CREATE INDEX idx_kyc_verification_attempts_user_id ON kyc_verification_attempts(user_id);
CREATE INDEX idx_kyb_verification_attempts_company_id ON kyb_verification_attempts(company_id);
CREATE INDEX idx_user_country_kyc_status_user_id ON user_country_kyc_status(user_id);
CREATE INDEX idx_company_country_kyb_status_company_id ON company_country_kyb_status(company_id);

-- +goose Down
DROP TABLE IF EXISTS company_country_kyb_status CASCADE;
DROP TABLE IF EXISTS user_country_kyc_status CASCADE;
DROP TABLE IF EXISTS kyb_verification_attempts CASCADE;
DROP TABLE IF EXISTS kyc_verification_attempts CASCADE;
DROP TABLE IF EXISTS kyb_documents CASCADE;
DROP TABLE IF EXISTS kyc_documents CASCADE;