-- +goose Up
-- Country Management and Compliance
CREATE TABLE supported_countries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_code VARCHAR(3) UNIQUE NOT NULL,
    country_name VARCHAR(100) NOT NULL,
    region VARCHAR(50),
    currency_code VARCHAR(3),
    currency_symbol VARCHAR(5),
    is_active BOOLEAN DEFAULT TRUE,
    is_high_risk BOOLEAN DEFAULT FALSE,
    requires_enhanced_kyc BOOLEAN DEFAULT FALSE,
    requires_enhanced_kyb BOOLEAN DEFAULT FALSE,
    timezone VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE kyc_country_requirements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_id UUID NOT NULL REFERENCES supported_countries(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    is_required BOOLEAN DEFAULT TRUE,
    requirement_description TEXT,
    acceptable_document_formats VARCHAR(255),
    verification_level VARCHAR(50),
    additional_attributes JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE kyb_country_requirements (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_id UUID NOT NULL REFERENCES supported_countries(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    business_type VARCHAR(50),
    is_required BOOLEAN DEFAULT TRUE,
    requirement_description TEXT,
    acceptable_document_formats VARCHAR(255),
    verification_level VARCHAR(50),
    additional_attributes JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE compliance_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_id UUID NOT NULL REFERENCES supported_countries(id) ON DELETE CASCADE,
    rule_type VARCHAR(50) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    threshold_amount DECIMAL(18, 6),
    threshold_currency VARCHAR(3),
    rule_description TEXT,
    regulatory_reference VARCHAR(255),
    action_required VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_supported_countries_country_code ON supported_countries(country_code);
CREATE INDEX idx_supported_countries_is_active ON supported_countries(is_active);
CREATE INDEX idx_kyc_country_requirements_country_id ON kyc_country_requirements(country_id);
CREATE INDEX idx_kyc_country_requirements_document_type ON kyc_country_requirements(document_type);
CREATE INDEX idx_kyb_country_requirements_country_id ON kyb_country_requirements(country_id);
CREATE INDEX idx_kyb_country_requirements_document_type ON kyb_country_requirements(document_type);
CREATE INDEX idx_compliance_rules_country_id ON compliance_rules(country_id);
CREATE INDEX idx_compliance_rules_rule_type ON compliance_rules(rule_type);

-- +goose Down
DROP TABLE IF EXISTS compliance_rules CASCADE;
DROP TABLE IF EXISTS kyb_country_requirements CASCADE;
DROP TABLE IF EXISTS kyc_country_requirements CASCADE;
DROP TABLE IF EXISTS supported_countries CASCADE;