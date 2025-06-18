-- +goose Up
-- Tax Management
CREATE TABLE tax_rates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    country_id UUID NOT NULL REFERENCES supported_countries(id) ON DELETE CASCADE,
    region VARCHAR(100),
    tax_type VARCHAR(50) NOT NULL,
    rate DECIMAL(6, 3) NOT NULL,
    effective_date DATE NOT NULL,
    expiry_date DATE,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_tax_rate CHECK (rate >= 0 AND rate <= 100),
    CONSTRAINT chk_tax_dates CHECK (expiry_date IS NULL OR expiry_date > effective_date)
);

CREATE TABLE tax_calculations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    reference_type VARCHAR(50) NOT NULL,
    reference_id UUID NOT NULL,
    tax_rate_id UUID NOT NULL REFERENCES tax_rates(id) ON DELETE CASCADE,
    taxable_amount DECIMAL(18, 6) NOT NULL,
    tax_amount DECIMAL(18, 6) NOT NULL,
    calculation_date DATE NOT NULL,
    tax_period VARCHAR(20) NOT NULL,
    status VARCHAR(50) DEFAULT 'calculated',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_tax_calculation_amounts CHECK (
        taxable_amount >= 0 AND tax_amount >= 0
    ),
    CONSTRAINT chk_tax_calculation_entity CHECK (
        (user_id IS NOT NULL AND company_id IS NULL) OR
        (user_id IS NULL AND company_id IS NOT NULL)
    )
);

CREATE TABLE tax_documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    country_id UUID NOT NULL REFERENCES supported_countries(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    tax_year INTEGER NOT NULL,
    document_url VARCHAR(255),
    ipfs_hash VARCHAR(255),
    status VARCHAR(50) DEFAULT 'draft',
    expires_at DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_tax_document_entity CHECK (
        (user_id IS NOT NULL AND company_id IS NULL) OR
        (user_id IS NULL AND company_id IS NOT NULL)
    ),
    CONSTRAINT chk_tax_year CHECK (tax_year >= 2000 AND tax_year <= 2100)
);

-- Create indexes
CREATE INDEX idx_tax_rates_country_id ON tax_rates(country_id);
CREATE INDEX idx_tax_rates_tax_type ON tax_rates(tax_type);
CREATE INDEX idx_tax_rates_effective_date ON tax_rates(effective_date);
CREATE INDEX idx_tax_rates_expiry_date ON tax_rates(expiry_date);
CREATE INDEX idx_tax_calculations_user_id ON tax_calculations(user_id);
CREATE INDEX idx_tax_calculations_company_id ON tax_calculations(company_id);
CREATE INDEX idx_tax_calculations_reference ON tax_calculations(reference_type, reference_id);
CREATE INDEX idx_tax_calculations_tax_rate_id ON tax_calculations(tax_rate_id);
CREATE INDEX idx_tax_calculations_calculation_date ON tax_calculations(calculation_date);
CREATE INDEX idx_tax_documents_user_id ON tax_documents(user_id);
CREATE INDEX idx_tax_documents_company_id ON tax_documents(company_id);
CREATE INDEX idx_tax_documents_country_id ON tax_documents(country_id);
CREATE INDEX idx_tax_documents_tax_year ON tax_documents(tax_year);
CREATE INDEX idx_tax_documents_document_type ON tax_documents(document_type);

-- +goose Down
DROP TABLE IF EXISTS tax_documents CASCADE;
DROP TABLE IF EXISTS tax_calculations CASCADE;
DROP TABLE IF EXISTS tax_rates CASCADE;