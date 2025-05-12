-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2025-05-12T16:19:10.975Z

CREATE TABLE "users" (
  "id" UUID PRIMARY KEY,
  "email" VARCHAR(255) UNIQUE,
  "password_hash" VARCHAR(255),
  "auth_provider" VARCHAR(50),
  "provider_id" VARCHAR(255),
  "email_verified" BOOLEAN,
  "email_verified_at" TIMESTAMPTZ,
  "account_type" VARCHAR(50),
  "account_status" VARCHAR(50),
  "two_factor_enabled" BOOLEAN,
  "two_factor_method" VARCHAR(50),
  "user_login_type" VARCHAR(50),
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ,
  "last_login_at" TIMESTAMPTZ,
  "deleted_at" TIMESTAMPTZ
);

CREATE TABLE "personal_users" (
  "id" UUID PRIMARY KEY,
  "first_name" VARCHAR(255),
  "last_name" VARCHAR(255),
  "profile_picture" VARCHAR(255),
  "phone_number" VARCHAR(50),
  "phone_number_verified" BOOLEAN,
  "phone_number_verified_at" TIMESTAMPTZ,
  "nationality" VARCHAR(255),
  "residential_country" VARCHAR(255),
  "user_address" VARCHAR(255),
  "user_city" VARCHAR(255),
  "user_postal_code" VARCHAR(255),
  "gender" VARCHAR(50),
  "date_of_birth" DATE,
  "job_role" VARCHAR(255),
  "personal_account_type" VARCHAR(50),
  "employment_type" VARCHAR(50),
  "tax_id" VARCHAR(255),
  "default_payment_currency" VARCHAR(50),
  "default_payment_method" VARCHAR(50),
  "hourly_rate" DECIMAL(18,2),
  "specialization" VARCHAR(255),
  "kyc_status" VARCHAR(50),
  "kyc_verified_at" TIMESTAMPTZ,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "companies" (
  "id" UUID PRIMARY KEY,
  "owner_id" UUID,
  "company_name" VARCHAR(255),
  "company_email" VARCHAR(255),
  "company_phone" VARCHAR(50),
  "company_size" VARCHAR(50),
  "company_industry" VARCHAR(255),
  "company_description" TEXT,
  "company_headquarters" VARCHAR(255),
  "company_logo" VARCHAR(255),
  "company_website" VARCHAR(255),
  "primary_contact_name" VARCHAR(255),
  "primary_contact_email" VARCHAR(255),
  "primary_contact_phone" VARCHAR(50),
  "company_address" VARCHAR(255),
  "company_city" VARCHAR(255),
  "company_postal_code" VARCHAR(255),
  "company_country" VARCHAR(255),
  "company_registration_number" VARCHAR(255),
  "registration_country" VARCHAR(255),
  "tax_id" VARCHAR(255),
  "incorporation_date" DATE,
  "account_status" VARCHAR(50),
  "kyb_status" VARCHAR(50),
  "kyb_verified_at" TIMESTAMPTZ,
  "kyb_verification_method" VARCHAR(50),
  "kyb_verification_provider" VARCHAR(255),
  "kyb_rejection_reason" TEXT,
  "legal_entity_type" VARCHAR(50),
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "company_users" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "user_id" UUID,
  "role" VARCHAR(50),
  "department" VARCHAR(100),
  "job_title" VARCHAR(255),
  "is_administrator" BOOLEAN,
  "can_manage_payroll" BOOLEAN,
  "can_manage_invoices" BOOLEAN,
  "can_manage_employees" BOOLEAN,
  "can_manage_company_settings" BOOLEAN,
  "can_manage_bank_accounts" BOOLEAN,
  "can_manage_wallets" BOOLEAN,
  "permissions" JSONB,
  "is_active" BOOLEAN,
  "added_by" UUID,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "company_staff_profiles" (
  "id" UUID PRIMARY KEY,
  "first_name" VARCHAR(255),
  "last_name" VARCHAR(255),
  "profile_picture" VARCHAR(255),
  "phone_number" VARCHAR(50),
  "email" VARCHAR(255),
  "department" VARCHAR(100),
  "job_title" VARCHAR(255),
  "reports_to" UUID,
  "hire_date" DATE,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "company_employees" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "user_id" UUID,
  "employee_id" VARCHAR(100),
  "department" VARCHAR(100),
  "position" VARCHAR(100),
  "employment_status" VARCHAR(50),
  "employment_type" VARCHAR(50),
  "start_date" DATE,
  "end_date" DATE,
  "manager_id" UUID,
  "salary_amount" DECIMAL(18,6),
  "salary_currency" VARCHAR(10),
  "salary_frequency" VARCHAR(50),
  "hourly_rate" DECIMAL(18,6),
  "payment_method" VARCHAR(50),
  "payment_split" JSONB,
  "tax_information" JSONB,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "sessions" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "refresh_token" VARCHAR(1024),
  "user_agent" TEXT,
  "client_ip" VARCHAR(45),
  "last_used_at" TIMESTAMPTZ,
  "web_oauth_client_id" TEXT,
  "oauth_access_token" TEXT,
  "oauth_id_token" TEXT,
  "user_login_type" VARCHAR(100),
  "mfa_verified" BOOLEAN,
  "is_blocked" BOOLEAN,
  "expires_at" TIMESTAMPTZ,
  "created_at" TIMESTAMPTZ
);

CREATE TABLE "user_devices" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "device_token" VARCHAR(255),
  "platform" VARCHAR(50),
  "device_type" VARCHAR(100),
  "device_model" VARCHAR(100),
  "os_name" VARCHAR(50),
  "os_version" VARCHAR(50),
  "push_notification_token" VARCHAR(255),
  "is_active" BOOLEAN,
  "is_verified" BOOLEAN,
  "last_used_at" TIMESTAMPTZ,
  "app_version" VARCHAR(50),
  "client_ip" VARCHAR(45),
  "expires_at" TIMESTAMPTZ,
  "is_revoked" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "security_events" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "company_id" UUID,
  "event_type" VARCHAR(100),
  "severity" VARCHAR(50),
  "ip_address" VARCHAR(45),
  "user_agent" TEXT,
  "metadata" JSONB,
  "created_at" TIMESTAMPTZ
);

CREATE TABLE "supported_countries" (
  "id" UUID PRIMARY KEY,
  "country_code" VARCHAR(3),
  "country_name" VARCHAR(100),
  "region" VARCHAR(50),
  "currency_code" VARCHAR(3),
  "currency_symbol" VARCHAR(5),
  "is_active" BOOLEAN,
  "is_high_risk" BOOLEAN,
  "requires_enhanced_kyc" BOOLEAN,
  "requires_enhanced_kyb" BOOLEAN,
  "timezone" VARCHAR(50),
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "kyc_country_requirements" (
  "id" UUID PRIMARY KEY,
  "country_id" UUID,
  "document_type" VARCHAR(50),
  "is_required" BOOLEAN,
  "requirement_description" TEXT,
  "acceptable_document_formats" VARCHAR(255),
  "verification_level" VARCHAR(50),
  "additional_attributes" JSONB,
  "is_active" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "kyb_country_requirements" (
  "id" UUID PRIMARY KEY,
  "country_id" UUID,
  "document_type" VARCHAR(50),
  "business_type" VARCHAR(50),
  "is_required" BOOLEAN,
  "requirement_description" TEXT,
  "acceptable_document_formats" VARCHAR(255),
  "verification_level" VARCHAR(50),
  "additional_attributes" JSONB,
  "is_active" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "compliance_rules" (
  "id" UUID PRIMARY KEY,
  "country_id" UUID,
  "rule_type" VARCHAR(50),
  "entity_type" VARCHAR(50),
  "threshold_amount" DECIMAL(18,6),
  "threshold_currency" VARCHAR(3),
  "rule_description" TEXT,
  "regulatory_reference" VARCHAR(255),
  "action_required" VARCHAR(50),
  "is_active" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "kyc_documents" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "country_id" UUID,
  "document_type" VARCHAR(50),
  "document_number" VARCHAR(100),
  "document_country" VARCHAR(100),
  "issue_date" DATE,
  "expiry_date" DATE,
  "document_url" VARCHAR(255),
  "ipfs_hash" VARCHAR(255),
  "verification_status" VARCHAR(50),
  "verification_level" VARCHAR(50),
  "verification_notes" TEXT,
  "verified_by" UUID,
  "verified_at" TIMESTAMPTZ,
  "rejection_reason" TEXT,
  "metadata" JSONB,
  "meets_requirements" BOOLEAN,
  "requirement_id" UUID,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "kyb_documents" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "country_id" UUID,
  "document_type" VARCHAR(50),
  "document_number" VARCHAR(100),
  "document_country" VARCHAR(100),
  "issue_date" DATE,
  "expiry_date" DATE,
  "document_url" VARCHAR(255),
  "ipfs_hash" VARCHAR(255),
  "verification_status" VARCHAR(50),
  "verification_level" VARCHAR(50),
  "verification_notes" TEXT,
  "verified_by" UUID,
  "verified_at" TIMESTAMPTZ,
  "rejection_reason" TEXT,
  "metadata" JSONB,
  "meets_requirements" BOOLEAN,
  "requirement_id" UUID,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "kyc_verification_attempts" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "verification_provider" VARCHAR(100),
  "verification_reference" VARCHAR(255),
  "verification_status" VARCHAR(50),
  "verification_result" VARCHAR(50),
  "response_data" JSONB,
  "error_message" TEXT,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "kyb_verification_attempts" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "verification_provider" VARCHAR(100),
  "verification_reference" VARCHAR(255),
  "verification_status" VARCHAR(50),
  "verification_result" VARCHAR(50),
  "response_data" JSONB,
  "error_message" TEXT,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "user_country_kyc_status" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "country_id" UUID,
  "verification_status" VARCHAR(50),
  "verification_level" VARCHAR(50),
  "verification_date" TIMESTAMPTZ,
  "expiry_date" TIMESTAMPTZ,
  "rejection_reason" TEXT,
  "notes" TEXT,
  "risk_rating" VARCHAR(20),
  "restricted_features" JSONB,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "company_country_kyb_status" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "country_id" UUID,
  "verification_status" VARCHAR(50),
  "verification_level" VARCHAR(50),
  "verification_date" TIMESTAMPTZ,
  "expiry_date" TIMESTAMPTZ,
  "rejection_reason" TEXT,
  "notes" TEXT,
  "risk_rating" VARCHAR(20),
  "restricted_features" JSONB,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "supported_networks" (
  "id" UUID PRIMARY KEY,
  "name" VARCHAR(100),
  "chain_id" INTEGER,
  "network_type" VARCHAR(50),
  "currency_symbol" VARCHAR(10),
  "block_explorer_url" VARCHAR(255),
  "rpc_url" VARCHAR(255),
  "is_evm_compatible" BOOLEAN,
  "is_active" BOOLEAN,
  "transaction_speed" VARCHAR(50),
  "average_block_time" INTEGER,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "supported_tokens" (
  "id" UUID PRIMARY KEY,
  "network_id" UUID,
  "name" VARCHAR(100),
  "symbol" VARCHAR(20),
  "decimals" INTEGER,
  "contract_address" VARCHAR(255),
  "token_type" VARCHAR(50),
  "logo_url" VARCHAR(255),
  "is_stablecoin" BOOLEAN,
  "is_active" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "user_wallets" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "wallet_address" VARCHAR(255),
  "wallet_type" VARCHAR(50),
  "chain_id" INTEGER,
  "is_default" BOOLEAN,
  "is_verified" BOOLEAN,
  "verification_method" VARCHAR(50),
  "verified_at" TIMESTAMPTZ,
  "nickname" VARCHAR(100),
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "company_wallets" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "wallet_address" VARCHAR(255),
  "wallet_type" VARCHAR(50),
  "chain_id" INTEGER,
  "is_default" BOOLEAN,
  "multisig_config" JSONB,
  "required_approvals" INTEGER,
  "wallet_name" VARCHAR(100),
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "bank_accounts" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "company_id" UUID,
  "account_number" VARCHAR(255),
  "account_holder_name" VARCHAR(255),
  "bank_name" VARCHAR(255),
  "bank_code" VARCHAR(100),
  "routing_number" VARCHAR(100),
  "swift_code" VARCHAR(50),
  "iban" VARCHAR(100),
  "account_type" VARCHAR(50),
  "currency" VARCHAR(3),
  "country" VARCHAR(100),
  "is_default" BOOLEAN,
  "is_verified" BOOLEAN,
  "verification_method" VARCHAR(50),
  "verified_at" TIMESTAMPTZ,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "payroll_periods" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "period_name" VARCHAR(100),
  "frequency" VARCHAR(50),
  "start_date" DATE,
  "end_date" DATE,
  "payment_date" DATE,
  "status" VARCHAR(50),
  "is_recurring" BOOLEAN,
  "next_period_id" UUID,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "payrolls" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "period_id" UUID,
  "name" VARCHAR(255),
  "description" TEXT,
  "total_amount" DECIMAL(24,6),
  "base_currency" VARCHAR(10),
  "status" VARCHAR(50),
  "execution_type" VARCHAR(50),
  "scheduled_execution_time" TIMESTAMPTZ,
  "executed_at" TIMESTAMPTZ,
  "smart_contract_address" VARCHAR(255),
  "chain_id" INTEGER,
  "transaction_hash" VARCHAR(255),
  "created_by" UUID,
  "approved_by" UUID,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "payroll_items" (
  "id" UUID PRIMARY KEY,
  "payroll_id" UUID,
  "employee_id" UUID,
  "base_amount" DECIMAL(18,6),
  "base_currency" VARCHAR(10),
  "payment_amount" DECIMAL(18,6),
  "payment_currency" VARCHAR(10),
  "exchange_rate" DECIMAL(24,12),
  "payment_method" VARCHAR(50),
  "payment_split" JSONB,
  "status" VARCHAR(50),
  "transaction_hash" VARCHAR(255),
  "recipient_wallet_address" VARCHAR(255),
  "recipient_bank_account_id" UUID,
  "notes" TEXT,
  "timesheet_id" UUID,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "invoices" (
  "id" UUID PRIMARY KEY,
  "invoice_number" VARCHAR(100),
  "issuer_id" UUID,
  "recipient_id" UUID,
  "title" VARCHAR(255),
  "description" TEXT,
  "issue_date" DATE,
  "due_date" DATE,
  "total_amount" DECIMAL(18,6),
  "currency" VARCHAR(10),
  "status" VARCHAR(50),
  "payment_method" VARCHAR(50),
  "recipient_wallet_address" VARCHAR(255),
  "recipient_bank_account_id" UUID,
  "transaction_hash" VARCHAR(255),
  "payment_date" TIMESTAMPTZ,
  "rejection_reason" TEXT,
  "ipfs_hash" VARCHAR(255),
  "smart_contract_address" VARCHAR(255),
  "chain_id" INTEGER,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "invoice_items" (
  "id" UUID PRIMARY KEY,
  "invoice_id" UUID,
  "description" TEXT,
  "quantity" DECIMAL(10,2),
  "unit_price" DECIMAL(18,6),
  "amount" DECIMAL(18,6),
  "tax_rate" DECIMAL(5,2),
  "tax_amount" DECIMAL(18,6),
  "discount_percentage" DECIMAL(5,2),
  "discount_amount" DECIMAL(18,6),
  "total_amount" DECIMAL(18,6),
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "contracts" (
  "id" UUID PRIMARY KEY,
  "template_id" UUID,
  "company_id" UUID,
  "employee_id" UUID,
  "freelancer_id" UUID,
  "contract_title" VARCHAR(255),
  "contract_type" VARCHAR(50),
  "start_date" DATE,
  "end_date" DATE,
  "status" VARCHAR(50),
  "payment_terms" JSONB,
  "contract_document_url" VARCHAR(255),
  "ipfs_hash" VARCHAR(255),
  "smart_contract_address" VARCHAR(255),
  "chain_id" INTEGER,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "contract_templates" (
  "id" UUID PRIMARY KEY,
  "creator_id" UUID,
  "company_id" UUID,
  "name" VARCHAR(255),
  "description" TEXT,
  "template_type" VARCHAR(50),
  "template_content" JSONB,
  "is_public" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "payment_requests" (
  "id" UUID PRIMARY KEY,
  "creator_id" UUID,
  "recipient_id" UUID,
  "company_id" UUID,
  "request_title" VARCHAR(255),
  "description" TEXT,
  "amount" DECIMAL(18,6),
  "currency" VARCHAR(10),
  "status" VARCHAR(50),
  "expiry_date" TIMESTAMPTZ,
  "payment_link" VARCHAR(255),
  "qr_code_url" VARCHAR(255),
  "payment_method" VARCHAR(50),
  "recipient_wallet_address" VARCHAR(255),
  "recipient_bank_account_id" UUID,
  "transaction_hash" VARCHAR(255),
  "paid_at" TIMESTAMPTZ,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "timesheets" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "employee_id" UUID,
  "period_id" UUID,
  "status" VARCHAR(50),
  "total_hours" DECIMAL(8,2),
  "billable_hours" DECIMAL(8,2),
  "overtime_hours" DECIMAL(8,2),
  "hourly_rate" DECIMAL(18,6),
  "rate_currency" VARCHAR(10),
  "total_amount" DECIMAL(18,6),
  "submitted_at" TIMESTAMPTZ,
  "approved_at" TIMESTAMPTZ,
  "approved_by" UUID,
  "rejection_reason" TEXT,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "timesheet_entries" (
  "id" UUID PRIMARY KEY,
  "timesheet_id" UUID,
  "date" DATE,
  "start_time" TIME,
  "end_time" TIME,
  "hours" DECIMAL(5,2),
  "is_billable" BOOLEAN,
  "is_overtime" BOOLEAN,
  "project" VARCHAR(255),
  "task" VARCHAR(255),
  "description" TEXT,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "wallet_transactions" (
  "id" UUID PRIMARY KEY,
  "wallet_address" VARCHAR(255),
  "transaction_hash" VARCHAR(255),
  "chain_id" INTEGER,
  "block_number" BIGINT,
  "from_address" VARCHAR(255),
  "to_address" VARCHAR(255),
  "token_address" VARCHAR(255),
  "token_symbol" VARCHAR(20),
  "amount" DECIMAL(36,18),
  "transaction_type" VARCHAR(50),
  "transaction_status" VARCHAR(50),
  "gas_price" DECIMAL(36,18),
  "gas_used" BIGINT,
  "transaction_fee" DECIMAL(36,18),
  "reference_type" VARCHAR(50),
  "reference_id" UUID,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "fiat_transactions" (
  "id" UUID PRIMARY KEY,
  "bank_account_id" UUID,
  "transaction_reference" VARCHAR(255),
  "transaction_type" VARCHAR(50),
  "amount" DECIMAL(18,6),
  "currency" VARCHAR(3),
  "status" VARCHAR(50),
  "payment_provider" VARCHAR(100),
  "payment_method" VARCHAR(50),
  "provider_reference" VARCHAR(255),
  "provider_fee" DECIMAL(18,6),
  "reference_type" VARCHAR(50),
  "reference_id" UUID,
  "metadata" JSONB,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "exchange_rates" (
  "id" UUID PRIMARY KEY,
  "base_currency" VARCHAR(10),
  "quote_currency" VARCHAR(10),
  "rate" DECIMAL(24,12),
  "source" VARCHAR(100),
  "timestamp" TIMESTAMPTZ,
  "created_at" TIMESTAMPTZ
);

CREATE TABLE "leave_types" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "name" VARCHAR(100),
  "description" TEXT,
  "is_paid" BOOLEAN,
  "accrual_rate" DECIMAL(5,2),
  "accrual_period" VARCHAR(20),
  "maximum_balance" DECIMAL(5,2),
  "is_active" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "leave_balances" (
  "id" UUID PRIMARY KEY,
  "employee_id" UUID,
  "leave_type_id" UUID,
  "balance" DECIMAL(6,2),
  "accrued" DECIMAL(6,2),
  "used" DECIMAL(6,2),
  "last_accrual_date" DATE,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "leave_requests" (
  "id" UUID PRIMARY KEY,
  "employee_id" UUID,
  "leave_type_id" UUID,
  "start_date" DATE,
  "end_date" DATE,
  "days" DECIMAL(5,2),
  "reason" TEXT,
  "status" VARCHAR(50),
  "approved_by" UUID,
  "approved_at" TIMESTAMPTZ,
  "rejected_reason" TEXT,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "expense_categories" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "name" VARCHAR(100),
  "description" TEXT,
  "is_active" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "expenses" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "company_id" UUID,
  "category_id" UUID,
  "amount" DECIMAL(18,6),
  "currency" VARCHAR(10),
  "expense_date" DATE,
  "description" TEXT,
  "receipt_url" VARCHAR(255),
  "ipfs_hash" VARCHAR(255),
  "status" VARCHAR(50),
  "payment_transaction_id" UUID,
  "approved_by" UUID,
  "approved_at" TIMESTAMPTZ,
  "rejected_reason" TEXT,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "tax_rates" (
  "id" UUID PRIMARY KEY,
  "country_id" UUID,
  "region" VARCHAR(100),
  "tax_type" VARCHAR(50),
  "rate" DECIMAL(6,3),
  "effective_date" DATE,
  "expiry_date" DATE,
  "description" TEXT,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "tax_calculations" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "company_id" UUID,
  "reference_type" VARCHAR(50),
  "reference_id" UUID,
  "tax_rate_id" UUID,
  "taxable_amount" DECIMAL(18,6),
  "tax_amount" DECIMAL(18,6),
  "calculation_date" DATE,
  "tax_period" VARCHAR(20),
  "status" VARCHAR(50),
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "tax_documents" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "company_id" UUID,
  "country_id" UUID,
  "document_type" VARCHAR(50),
  "tax_year" INTEGER,
  "document_url" VARCHAR(255),
  "ipfs_hash" VARCHAR(255),
  "status" VARCHAR(50),
  "expires_at" DATE,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "notification_templates" (
  "id" UUID PRIMARY KEY,
  "template_name" VARCHAR(100),
  "template_type" VARCHAR(50),
  "subject" VARCHAR(255),
  "content" TEXT,
  "variables" JSONB,
  "is_active" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "notifications" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "template_id" UUID,
  "notification_type" VARCHAR(50),
  "title" VARCHAR(255),
  "content" TEXT,
  "reference_type" VARCHAR(50),
  "reference_id" UUID,
  "is_read" BOOLEAN,
  "read_at" TIMESTAMPTZ,
  "delivery_status" VARCHAR(50),
  "priority" VARCHAR(20),
  "created_at" TIMESTAMPTZ
);

CREATE TABLE "roles" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "role_name" VARCHAR(100),
  "description" TEXT,
  "is_system_role" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "permissions" (
  "id" UUID PRIMARY KEY,
  "permission_key" VARCHAR(100),
  "description" TEXT,
  "category" VARCHAR(100),
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "role_permissions" (
  "id" UUID PRIMARY KEY,
  "role_id" UUID,
  "permission_id" UUID,
  "created_at" TIMESTAMPTZ
);

CREATE TABLE "user_roles" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "role_id" UUID,
  "company_id" UUID,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "system_settings" (
  "id" UUID PRIMARY KEY,
  "setting_key" VARCHAR(100),
  "setting_value" TEXT,
  "data_type" VARCHAR(50),
  "description" TEXT,
  "is_sensitive" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "company_settings" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "setting_key" VARCHAR(100),
  "setting_value" TEXT,
  "data_type" VARCHAR(50),
  "description" TEXT,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "user_settings" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "setting_key" VARCHAR(100),
  "setting_value" TEXT,
  "data_type" VARCHAR(50),
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "audit_logs" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "company_id" UUID,
  "action" VARCHAR(100),
  "entity_type" VARCHAR(100),
  "entity_id" UUID,
  "previous_state" JSONB,
  "new_state" JSONB,
  "ip_address" VARCHAR(45),
  "user_agent" TEXT,
  "created_at" TIMESTAMPTZ
);

CREATE TABLE "activity_logs" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "activity_type" VARCHAR(100),
  "description" TEXT,
  "metadata" JSONB,
  "ip_address" VARCHAR(45),
  "user_agent" TEXT,
  "created_at" TIMESTAMPTZ
);

CREATE TABLE "api_keys" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "company_id" UUID,
  "api_key_hash" VARCHAR(255),
  "name" VARCHAR(100),
  "permissions" JSONB,
  "expires_at" TIMESTAMPTZ,
  "is_active" BOOLEAN,
  "last_used_at" TIMESTAMPTZ,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "integration_connections" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "company_id" UUID,
  "integration_type" VARCHAR(100),
  "provider" VARCHAR(100),
  "access_token" TEXT,
  "refresh_token" TEXT,
  "token_expires_at" TIMESTAMPTZ,
  "connection_data" JSONB,
  "is_active" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "webhooks" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "company_id" UUID,
  "webhook_url" VARCHAR(255),
  "event_types" VARCHAR[],
  "secret_key" VARCHAR(255),
  "description" TEXT,
  "is_active" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "feature_flags" (
  "id" UUID PRIMARY KEY,
  "flag_key" VARCHAR(100),
  "description" TEXT,
  "is_enabled" BOOLEAN,
  "rollout_percentage" INTEGER,
  "conditions" JSONB,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "user_feature_flags" (
  "id" UUID PRIMARY KEY,
  "user_id" UUID,
  "flag_key" VARCHAR(100),
  "is_enabled" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

CREATE TABLE "company_feature_flags" (
  "id" UUID PRIMARY KEY,
  "company_id" UUID,
  "flag_key" VARCHAR(100),
  "is_enabled" BOOLEAN,
  "created_at" TIMESTAMPTZ,
  "updated_at" TIMESTAMPTZ
);

ALTER TABLE "personal_users" ADD FOREIGN KEY ("id") REFERENCES "users" ("id");

ALTER TABLE "companies" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id");

ALTER TABLE "company_users" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "company_users" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "company_users" ADD FOREIGN KEY ("added_by") REFERENCES "users" ("id");

ALTER TABLE "company_staff_profiles" ADD FOREIGN KEY ("id") REFERENCES "company_users" ("id");

ALTER TABLE "company_staff_profiles" ADD FOREIGN KEY ("reports_to") REFERENCES "company_users" ("id");

ALTER TABLE "company_employees" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "company_employees" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "company_employees" ADD FOREIGN KEY ("manager_id") REFERENCES "company_users" ("id");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_devices" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "security_events" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "security_events" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "kyc_country_requirements" ADD FOREIGN KEY ("country_id") REFERENCES "supported_countries" ("id");

ALTER TABLE "kyb_country_requirements" ADD FOREIGN KEY ("country_id") REFERENCES "supported_countries" ("id");

ALTER TABLE "compliance_rules" ADD FOREIGN KEY ("country_id") REFERENCES "supported_countries" ("id");

ALTER TABLE "kyc_documents" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "kyc_documents" ADD FOREIGN KEY ("country_id") REFERENCES "supported_countries" ("id");

ALTER TABLE "kyc_documents" ADD FOREIGN KEY ("verified_by") REFERENCES "users" ("id");

ALTER TABLE "kyc_documents" ADD FOREIGN KEY ("requirement_id") REFERENCES "kyc_country_requirements" ("id");

ALTER TABLE "kyb_documents" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "kyb_documents" ADD FOREIGN KEY ("country_id") REFERENCES "supported_countries" ("id");

ALTER TABLE "kyb_documents" ADD FOREIGN KEY ("verified_by") REFERENCES "users" ("id");

ALTER TABLE "kyb_documents" ADD FOREIGN KEY ("requirement_id") REFERENCES "kyb_country_requirements" ("id");

ALTER TABLE "kyc_verification_attempts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "kyb_verification_attempts" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "user_country_kyc_status" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_country_kyc_status" ADD FOREIGN KEY ("country_id") REFERENCES "supported_countries" ("id");

ALTER TABLE "company_country_kyb_status" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "company_country_kyb_status" ADD FOREIGN KEY ("country_id") REFERENCES "supported_countries" ("id");

ALTER TABLE "supported_tokens" ADD FOREIGN KEY ("network_id") REFERENCES "supported_networks" ("id");

ALTER TABLE "user_wallets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "company_wallets" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "bank_accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "bank_accounts" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "payroll_periods" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "payroll_periods" ADD FOREIGN KEY ("next_period_id") REFERENCES "payroll_periods" ("id");

ALTER TABLE "payrolls" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "payrolls" ADD FOREIGN KEY ("period_id") REFERENCES "payroll_periods" ("id");

ALTER TABLE "payrolls" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "payrolls" ADD FOREIGN KEY ("approved_by") REFERENCES "users" ("id");

ALTER TABLE "payroll_items" ADD FOREIGN KEY ("payroll_id") REFERENCES "payrolls" ("id");

ALTER TABLE "payroll_items" ADD FOREIGN KEY ("employee_id") REFERENCES "company_employees" ("id");

ALTER TABLE "payroll_items" ADD FOREIGN KEY ("recipient_bank_account_id") REFERENCES "bank_accounts" ("id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("issuer_id") REFERENCES "users" ("id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("recipient_id") REFERENCES "companies" ("id");

ALTER TABLE "invoices" ADD FOREIGN KEY ("recipient_bank_account_id") REFERENCES "bank_accounts" ("id");

ALTER TABLE "invoice_items" ADD FOREIGN KEY ("invoice_id") REFERENCES "invoices" ("id");

ALTER TABLE "contracts" ADD FOREIGN KEY ("template_id") REFERENCES "contract_templates" ("id");

ALTER TABLE "contracts" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "contracts" ADD FOREIGN KEY ("employee_id") REFERENCES "company_employees" ("id");

ALTER TABLE "contracts" ADD FOREIGN KEY ("freelancer_id") REFERENCES "users" ("id");

ALTER TABLE "contract_templates" ADD FOREIGN KEY ("creator_id") REFERENCES "users" ("id");

ALTER TABLE "contract_templates" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "payment_requests" ADD FOREIGN KEY ("creator_id") REFERENCES "users" ("id");

ALTER TABLE "payment_requests" ADD FOREIGN KEY ("recipient_id") REFERENCES "users" ("id");

ALTER TABLE "payment_requests" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "payment_requests" ADD FOREIGN KEY ("recipient_bank_account_id") REFERENCES "bank_accounts" ("id");

ALTER TABLE "timesheets" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "timesheets" ADD FOREIGN KEY ("employee_id") REFERENCES "company_employees" ("id");

ALTER TABLE "timesheets" ADD FOREIGN KEY ("period_id") REFERENCES "payroll_periods" ("id");

ALTER TABLE "timesheets" ADD FOREIGN KEY ("approved_by") REFERENCES "users" ("id");

ALTER TABLE "payroll_items" ADD FOREIGN KEY ("timesheet_id") REFERENCES "timesheets" ("id");

ALTER TABLE "timesheet_entries" ADD FOREIGN KEY ("timesheet_id") REFERENCES "timesheets" ("id");

ALTER TABLE "fiat_transactions" ADD FOREIGN KEY ("bank_account_id") REFERENCES "bank_accounts" ("id");

ALTER TABLE "leave_types" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "leave_balances" ADD FOREIGN KEY ("employee_id") REFERENCES "company_employees" ("id");

ALTER TABLE "leave_balances" ADD FOREIGN KEY ("leave_type_id") REFERENCES "leave_types" ("id");

ALTER TABLE "leave_requests" ADD FOREIGN KEY ("employee_id") REFERENCES "company_employees" ("id");

ALTER TABLE "leave_requests" ADD FOREIGN KEY ("leave_type_id") REFERENCES "leave_types" ("id");

ALTER TABLE "leave_requests" ADD FOREIGN KEY ("approved_by") REFERENCES "users" ("id");

ALTER TABLE "expense_categories" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "expenses" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "expenses" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "expenses" ADD FOREIGN KEY ("category_id") REFERENCES "expense_categories" ("id");

ALTER TABLE "expenses" ADD FOREIGN KEY ("approved_by") REFERENCES "users" ("id");

ALTER TABLE "tax_rates" ADD FOREIGN KEY ("country_id") REFERENCES "supported_countries" ("id");

ALTER TABLE "tax_calculations" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tax_calculations" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "tax_calculations" ADD FOREIGN KEY ("tax_rate_id") REFERENCES "tax_rates" ("id");

ALTER TABLE "tax_documents" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tax_documents" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "tax_documents" ADD FOREIGN KEY ("country_id") REFERENCES "supported_countries" ("id");

ALTER TABLE "notifications" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "notifications" ADD FOREIGN KEY ("template_id") REFERENCES "notification_templates" ("id");

ALTER TABLE "roles" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "role_permissions" ADD FOREIGN KEY ("permission_id") REFERENCES "permissions" ("id");

ALTER TABLE "user_roles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "user_roles" ADD FOREIGN KEY ("role_id") REFERENCES "roles" ("id");

ALTER TABLE "user_roles" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "company_settings" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "user_settings" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "audit_logs" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "activity_logs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "api_keys" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "api_keys" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "integration_connections" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "integration_connections" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "webhooks" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "webhooks" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "user_feature_flags" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "company_feature_flags" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");
