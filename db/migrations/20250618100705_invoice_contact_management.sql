-- +goose Up
-- Invoice and Contract Management
CREATE TABLE invoices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_number VARCHAR(100) UNIQUE NOT NULL,
    issuer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    issue_date DATE NOT NULL,
    due_date DATE NOT NULL,
    total_amount DECIMAL(18, 6) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    payment_method VARCHAR(50),
    recipient_wallet_address VARCHAR(255),
    recipient_bank_account_id UUID REFERENCES bank_accounts(id),
    transaction_hash VARCHAR(255),
    payment_date TIMESTAMPTZ,
    rejection_reason TEXT,
    ipfs_hash VARCHAR(255),
    smart_contract_address VARCHAR(255),
    chain_id INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE invoice_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    invoice_id UUID NOT NULL REFERENCES invoices(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL DEFAULT 1,
    unit_price DECIMAL(18, 6) NOT NULL,
    amount DECIMAL(18, 6) NOT NULL,
    tax_rate DECIMAL(5, 2) DEFAULT 0,
    tax_amount DECIMAL(18, 6) DEFAULT 0,
    discount_percentage DECIMAL(5, 2) DEFAULT 0,
    discount_amount DECIMAL(18, 6) DEFAULT 0,
    total_amount DECIMAL(18, 6) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE contract_templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    template_type VARCHAR(50) NOT NULL,
    template_content JSONB NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE contracts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id UUID REFERENCES contract_templates(id),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID REFERENCES company_employees(id) ON DELETE CASCADE,
    freelancer_id UUID REFERENCES users(id) ON DELETE CASCADE,
    contract_title VARCHAR(255) NOT NULL,
    contract_type VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    status VARCHAR(50) DEFAULT 'draft',
    payment_terms JSONB,
    contract_document_url VARCHAR(255),
    ipfs_hash VARCHAR(255),
    smart_contract_address VARCHAR(255),
    chain_id INTEGER,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT chk_contract_party CHECK (
        (employee_id IS NOT NULL AND freelancer_id IS NULL) OR
        (employee_id IS NULL AND freelancer_id IS NOT NULL)
    )
);

CREATE TABLE payment_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    recipient_id UUID REFERENCES users(id) ON DELETE CASCADE,
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    request_title VARCHAR(255) NOT NULL,
    description TEXT,
    amount DECIMAL(18, 6) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    expiry_date TIMESTAMPTZ,
    payment_link VARCHAR(255),
    qr_code_url VARCHAR(255),
    payment_method VARCHAR(50),
    recipient_wallet_address VARCHAR(255),
    recipient_bank_account_id UUID REFERENCES bank_accounts(id),
    transaction_hash VARCHAR(255),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX idx_invoices_issuer_id ON invoices(issuer_id);
CREATE INDEX idx_invoices_recipient_id ON invoices(recipient_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_due_date ON invoices(due_date);
CREATE INDEX idx_invoice_items_invoice_id ON invoice_items(invoice_id);
CREATE INDEX idx_contract_templates_creator_id ON contract_templates(creator_id);
CREATE INDEX idx_contract_templates_company_id ON contract_templates(company_id);
CREATE INDEX idx_contracts_template_id ON contracts(template_id);
CREATE INDEX idx_contracts_company_id ON contracts(company_id);
CREATE INDEX idx_contracts_employee_id ON contracts(employee_id);
CREATE INDEX idx_contracts_freelancer_id ON contracts(freelancer_id);
CREATE INDEX idx_payment_requests_creator_id ON payment_requests(creator_id);
CREATE INDEX idx_payment_requests_recipient_id ON payment_requests(recipient_id);
CREATE INDEX idx_payment_requests_status ON payment_requests(status);

-- +goose Down
DROP TABLE IF EXISTS payment_requests CASCADE;
DROP TABLE IF EXISTS contracts CASCADE;
DROP TABLE IF EXISTS contract_templates CASCADE;
DROP TABLE IF EXISTS invoice_items CASCADE;
DROP TABLE IF EXISTS invoices CASCADE;