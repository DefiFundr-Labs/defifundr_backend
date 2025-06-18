-- +goose Up
-- Functions and Triggers

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS '
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_invoice_number()
RETURNS TRIGGER AS '
BEGIN
    IF NEW.invoice_number IS NULL THEN
        NEW.invoice_number := ''INV-'' || EXTRACT(YEAR FROM NOW()) || ''-'' || 
                             LPAD(NEXTVAL(''invoice_number_seq'')::TEXT, 6, ''0'');
    END IF;
    RETURN NEW;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd

-- Create sequence for invoice numbers
CREATE SEQUENCE IF NOT EXISTS invoice_number_seq START 1;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION validate_payment_split()
RETURNS TRIGGER AS '
DECLARE
    split_total DECIMAL(5,2) := 0;
    split_item JSONB;
BEGIN
    IF NEW.payment_split IS NOT NULL THEN
        FOR split_item IN SELECT jsonb_array_elements(NEW.payment_split)
        LOOP
            split_total := split_total + (split_item->>''percentage'')::DECIMAL(5,2);
        END LOOP;
        
        IF split_total != 100.00 THEN
            RAISE EXCEPTION ''Payment split percentages must total 100, got %'', split_total;
        END IF;
    END IF;
    RETURN NEW;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_timesheet_totals()
RETURNS TRIGGER AS '
DECLARE
    timesheet_record RECORD;
BEGIN
    IF TG_OP = ''DELETE'' THEN
        SELECT id INTO timesheet_record FROM timesheets WHERE id = OLD.timesheet_id;
    ELSE
        SELECT id INTO timesheet_record FROM timesheets WHERE id = NEW.timesheet_id;
    END IF;
    
    UPDATE timesheets SET
        total_hours = (
            SELECT COALESCE(SUM(hours), 0) 
            FROM timesheet_entries 
            WHERE timesheet_id = timesheet_record.id
        ),
        billable_hours = (
            SELECT COALESCE(SUM(hours), 0) 
            FROM timesheet_entries 
            WHERE timesheet_id = timesheet_record.id AND is_billable = true
        ),
        overtime_hours = (
            SELECT COALESCE(SUM(hours), 0) 
            FROM timesheet_entries 
            WHERE timesheet_id = timesheet_record.id AND is_overtime = true
        ),
        updated_at = NOW()
    WHERE id = timesheet_record.id;
    
    RETURN COALESCE(NEW, OLD);
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd

-- Create updated_at triggers for all tables that have updated_at columns
CREATE TRIGGER trigger_users_updated_at BEFORE UPDATE ON users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_personal_users_updated_at BEFORE UPDATE ON personal_users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_companies_updated_at BEFORE UPDATE ON companies 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_company_users_updated_at BEFORE UPDATE ON company_users 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_company_staff_profiles_updated_at BEFORE UPDATE ON company_staff_profiles 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_company_employees_updated_at BEFORE UPDATE ON company_employees 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_user_devices_updated_at BEFORE UPDATE ON user_devices 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_supported_countries_updated_at BEFORE UPDATE ON supported_countries 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_kyc_country_requirements_updated_at BEFORE UPDATE ON kyc_country_requirements 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_kyb_country_requirements_updated_at BEFORE UPDATE ON kyb_country_requirements 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_compliance_rules_updated_at BEFORE UPDATE ON compliance_rules 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_kyc_documents_updated_at BEFORE UPDATE ON kyc_documents 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_kyb_documents_updated_at BEFORE UPDATE ON kyb_documents 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_supported_networks_updated_at BEFORE UPDATE ON supported_networks 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_supported_tokens_updated_at BEFORE UPDATE ON supported_tokens 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_user_wallets_updated_at BEFORE UPDATE ON user_wallets 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_company_wallets_updated_at BEFORE UPDATE ON company_wallets 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_bank_accounts_updated_at BEFORE UPDATE ON bank_accounts 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_payroll_periods_updated_at BEFORE UPDATE ON payroll_periods 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_payrolls_updated_at BEFORE UPDATE ON payrolls 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_payroll_items_updated_at BEFORE UPDATE ON payroll_items 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_invoices_updated_at BEFORE UPDATE ON invoices 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_invoice_items_updated_at BEFORE UPDATE ON invoice_items 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_contracts_updated_at BEFORE UPDATE ON contracts 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_contract_templates_updated_at BEFORE UPDATE ON contract_templates 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_payment_requests_updated_at BEFORE UPDATE ON payment_requests 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_timesheets_updated_at BEFORE UPDATE ON timesheets 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_timesheet_entries_updated_at BEFORE UPDATE ON timesheet_entries 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_wallet_transactions_updated_at BEFORE UPDATE ON wallet_transactions 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trigger_fiat_transactions_updated_at BEFORE UPDATE ON fiat_transactions 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create invoice number generation trigger
CREATE TRIGGER trigger_generate_invoice_number BEFORE INSERT ON invoices 
    FOR EACH ROW EXECUTE FUNCTION generate_invoice_number();

-- Create payment split validation triggers
CREATE TRIGGER trigger_validate_payroll_payment_split BEFORE INSERT OR UPDATE ON payroll_items 
    FOR EACH ROW EXECUTE FUNCTION validate_payment_split();

CREATE TRIGGER trigger_validate_employee_payment_split BEFORE INSERT OR UPDATE ON company_employees 
    FOR EACH ROW EXECUTE FUNCTION validate_payment_split();

-- Create timesheet totals update triggers
CREATE TRIGGER trigger_update_timesheet_totals_insert AFTER INSERT ON timesheet_entries 
    FOR EACH ROW EXECUTE FUNCTION update_timesheet_totals();

CREATE TRIGGER trigger_update_timesheet_totals_update AFTER UPDATE ON timesheet_entries 
    FOR EACH ROW EXECUTE FUNCTION update_timesheet_totals();

CREATE TRIGGER trigger_update_timesheet_totals_delete AFTER DELETE ON timesheet_entries 
    FOR EACH ROW EXECUTE FUNCTION update_timesheet_totals();

-- +goose Down
-- Drop all triggers
DROP TRIGGER IF EXISTS trigger_update_timesheet_totals_delete ON timesheet_entries;
DROP TRIGGER IF EXISTS trigger_update_timesheet_totals_update ON timesheet_entries;
DROP TRIGGER IF EXISTS trigger_update_timesheet_totals_insert ON timesheet_entries;
DROP TRIGGER IF EXISTS trigger_validate_employee_payment_split ON company_employees;
DROP TRIGGER IF EXISTS trigger_validate_payroll_payment_split ON payroll_items;
DROP TRIGGER IF EXISTS trigger_generate_invoice_number ON invoices;
DROP TRIGGER IF EXISTS trigger_fiat_transactions_updated_at ON fiat_transactions;
DROP TRIGGER IF EXISTS trigger_wallet_transactions_updated_at ON wallet_transactions;
DROP TRIGGER IF EXISTS trigger_timesheet_entries_updated_at ON timesheet_entries;
DROP TRIGGER IF EXISTS trigger_timesheets_updated_at ON timesheets;
DROP TRIGGER IF EXISTS trigger_payment_requests_updated_at ON payment_requests;
DROP TRIGGER IF EXISTS trigger_contract_templates_updated_at ON contract_templates;
DROP TRIGGER IF EXISTS trigger_contracts_updated_at ON contracts;
DROP TRIGGER IF EXISTS trigger_invoice_items_updated_at ON invoice_items;
DROP TRIGGER IF EXISTS trigger_invoices_updated_at ON invoices;
DROP TRIGGER IF EXISTS trigger_payroll_items_updated_at ON payroll_items;
DROP TRIGGER IF EXISTS trigger_payrolls_updated_at ON payrolls;
DROP TRIGGER IF EXISTS trigger_payroll_periods_updated_at ON payroll_periods;
DROP TRIGGER IF EXISTS trigger_bank_accounts_updated_at ON bank_accounts;
DROP TRIGGER IF EXISTS trigger_company_wallets_updated_at ON company_wallets;
DROP TRIGGER IF EXISTS trigger_user_wallets_updated_at ON user_wallets;
DROP TRIGGER IF EXISTS trigger_supported_tokens_updated_at ON supported_tokens;
DROP TRIGGER IF EXISTS trigger_supported_networks_updated_at ON supported_networks;
DROP TRIGGER IF EXISTS trigger_kyb_documents_updated_at ON kyb_documents;
DROP TRIGGER IF EXISTS trigger_kyc_documents_updated_at ON kyc_documents;
DROP TRIGGER IF EXISTS trigger_compliance_rules_updated_at ON compliance_rules;
DROP TRIGGER IF EXISTS trigger_kyb_country_requirements_updated_at ON kyb_country_requirements;
DROP TRIGGER IF EXISTS trigger_kyc_country_requirements_updated_at ON kyc_country_requirements;
DROP TRIGGER IF EXISTS trigger_supported_countries_updated_at ON supported_countries;
DROP TRIGGER IF EXISTS trigger_user_devices_updated_at ON user_devices;
DROP TRIGGER IF EXISTS trigger_company_employees_updated_at ON company_employees;
DROP TRIGGER IF EXISTS trigger_company_staff_profiles_updated_at ON company_staff_profiles;
DROP TRIGGER IF EXISTS trigger_company_users_updated_at ON company_users;
DROP TRIGGER IF EXISTS trigger_companies_updated_at ON companies;
DROP TRIGGER IF EXISTS trigger_personal_users_updated_at ON personal_users;
DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;

-- Drop functions
DROP FUNCTION IF EXISTS update_timesheet_totals();
DROP FUNCTION IF EXISTS validate_payment_split();
DROP FUNCTION IF EXISTS generate_invoice_number();
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop sequence
DROP SEQUENCE IF EXISTS invoice_number_seq;