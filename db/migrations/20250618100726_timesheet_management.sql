-- +goose Up
-- Timesheet Management
CREATE TABLE timesheets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES company_employees(id) ON DELETE CASCADE,
    period_id UUID REFERENCES payroll_periods(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'draft',
    total_hours DECIMAL(8, 2) DEFAULT 0,
    billable_hours DECIMAL(8, 2) DEFAULT 0,
    overtime_hours DECIMAL(8, 2) DEFAULT 0,
    hourly_rate DECIMAL(18, 6),
    rate_currency VARCHAR(10),
    total_amount DECIMAL(18, 6) DEFAULT 0,
    submitted_at TIMESTAMPTZ,
    approved_at TIMESTAMPTZ,
    approved_by UUID REFERENCES users(id),
    rejection_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE timesheet_entries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timesheet_id UUID NOT NULL REFERENCES timesheets(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    start_time TIME,
    end_time TIME,
    hours DECIMAL(5, 2) NOT NULL,
    is_billable BOOLEAN DEFAULT TRUE,
    is_overtime BOOLEAN DEFAULT FALSE,
    project VARCHAR(255),
    task VARCHAR(255),
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Add the foreign key reference from payroll_items to timesheets
ALTER TABLE payroll_items ADD CONSTRAINT fk_payroll_items_timesheet_id 
    FOREIGN KEY (timesheet_id) REFERENCES timesheets(id);

-- Create indexes
CREATE INDEX idx_timesheets_company_id ON timesheets(company_id);
CREATE INDEX idx_timesheets_employee_id ON timesheets(employee_id);
CREATE INDEX idx_timesheets_period_id ON timesheets(period_id);
CREATE INDEX idx_timesheets_status ON timesheets(status);
CREATE INDEX idx_timesheets_submitted_at ON timesheets(submitted_at);
CREATE INDEX idx_timesheet_entries_timesheet_id ON timesheet_entries(timesheet_id);
CREATE INDEX idx_timesheet_entries_date ON timesheet_entries(date);
CREATE INDEX idx_timesheet_entries_project ON timesheet_entries(project);

-- +goose Down
-- Remove the foreign key constraint from payroll_items
ALTER TABLE payroll_items DROP CONSTRAINT IF EXISTS fk_payroll_items_timesheet_id;

DROP TABLE IF EXISTS timesheet_entries CASCADE;
DROP TABLE IF EXISTS timesheets CASCADE;