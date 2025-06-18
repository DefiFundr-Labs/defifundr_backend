-- +goose Up
-- Payroll Management
CREATE TABLE payroll_periods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    period_name VARCHAR(100) NOT NULL,
    frequency VARCHAR(50) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    payment_date DATE NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    is_recurring BOOLEAN DEFAULT FALSE,
    next_period_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Add self-referencing foreign key for next_period_id
ALTER TABLE payroll_periods ADD CONSTRAINT fk_next_period_id 
    FOREIGN KEY (next_period_id) REFERENCES payroll_periods(id);

CREATE TABLE payrolls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    period_id UUID NOT NULL REFERENCES payroll_periods(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    total_amount DECIMAL(24, 6) DEFAULT 0,
    base_currency VARCHAR(10) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft',
    execution_type VARCHAR(50) DEFAULT 'manual',
    scheduled_execution_time TIMESTAMPTZ,
    executed_at TIMESTAMPTZ,
    smart_contract_address VARCHAR(255),
    chain_id INTEGER,
    transaction_hash VARCHAR(255),
    created_by UUID NOT NULL REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE payroll_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    payroll_id UUID NOT NULL REFERENCES payrolls(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES company_employees(id) ON DELETE CASCADE,
    base_amount DECIMAL(18, 6) NOT NULL,
    base_currency VARCHAR(10) NOT NULL,
    payment_amount DECIMAL(18, 6) NOT NULL,
    payment_currency VARCHAR(10) NOT NULL,
    exchange_rate DECIMAL(24, 12) DEFAULT 1,
    payment_method VARCHAR(50) NOT NULL,
    payment_split JSONB,
    status VARCHAR(50) DEFAULT 'pending',
    transaction_hash VARCHAR(255),
    recipient_wallet_address VARCHAR(255),
    recipient_bank_account_id UUID REFERENCES bank_accounts(id),
    notes TEXT,
    timesheet_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Create indexes
CREATE INDEX idx_payroll_periods_company_id ON payroll_periods(company_id);
CREATE INDEX idx_payroll_periods_start_date ON payroll_periods(start_date);
CREATE INDEX idx_payroll_periods_end_date ON payroll_periods(end_date);
CREATE INDEX idx_payroll_periods_status ON payroll_periods(status);
CREATE INDEX idx_payrolls_company_id ON payrolls(company_id);
CREATE INDEX idx_payrolls_period_id ON payrolls(period_id);
CREATE INDEX idx_payrolls_status ON payrolls(status);
CREATE INDEX idx_payrolls_created_by ON payrolls(created_by);
CREATE INDEX idx_payroll_items_payroll_id ON payroll_items(payroll_id);
CREATE INDEX idx_payroll_items_employee_id ON payroll_items(employee_id);
CREATE INDEX idx_payroll_items_status ON payroll_items(status);
CREATE INDEX idx_payroll_items_timesheet_id ON payroll_items(timesheet_id);

-- +goose Down
DROP TABLE IF EXISTS payroll_items CASCADE;
DROP TABLE IF EXISTS payrolls CASCADE;
DROP TABLE IF EXISTS payroll_periods CASCADE;