-- name: CreateCompany :one
INSERT INTO companies (
  id,
  owner_id,
  company_name,
  company_email,
  company_phone,
  company_size,
  company_industry,
  company_description,
  company_headquarters,
  company_logo,
  company_website,
  primary_contact_name,
  primary_contact_email,
  primary_contact_phone,
  company_address,
  company_city,
  company_postal_code,
  company_country,
  company_registration_number,
  registration_country,
  tax_id,
  incorporation_date,
  account_status,
  kyb_status,
  kyb_verified_at,
  kyb_verification_method,
  kyb_verification_provider,
  kyb_rejection_reason,
  legal_entity_type,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @owner_id,
  @company_name,
  @company_email,
  @company_phone,
  @company_size,
  @company_industry,
  @company_description,
  @company_headquarters,
  @company_logo,
  @company_website,
  @primary_contact_name,
  @primary_contact_email,
  @primary_contact_phone,
  @company_address,
  @company_city,
  @company_postal_code,
  @company_country,
  @company_registration_number,
  @registration_country,
  @tax_id,
  @incorporation_date,
  COALESCE(@account_status, 'pending'),
  COALESCE(@kyb_status, 'pending'),
  @kyb_verified_at,
  @kyb_verification_method,
  @kyb_verification_provider,
  @kyb_rejection_reason,
  @legal_entity_type,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetCompanyByID :one
SELECT * FROM companies WHERE id = @id;

-- name: GetCompaniesByOwner :many
SELECT * FROM companies 
WHERE owner_id = @owner_id
ORDER BY created_at DESC;

-- name: GetCompaniesAccessibleToUser :many
SELECT DISTINCT c.*
FROM companies c
LEFT JOIN company_users cu ON c.id = cu.company_id
WHERE c.owner_id = @user_id OR cu.user_id = @user_id
ORDER BY c.created_at DESC;

-- name: UpdateCompany :one
UPDATE companies SET
  company_name = COALESCE(@company_name, company_name),
  company_email = COALESCE(@company_email, company_email),
  company_phone = COALESCE(@company_phone, company_phone),
  company_size = COALESCE(@company_size, company_size),
  company_industry = COALESCE(@company_industry, company_industry),
  company_description = COALESCE(@company_description, company_description),
  company_headquarters = COALESCE(@company_headquarters, company_headquarters),
  company_logo = COALESCE(@company_logo, company_logo),
  company_website = COALESCE(@company_website, company_website),
  primary_contact_name = COALESCE(@primary_contact_name, primary_contact_name),
  primary_contact_email = COALESCE(@primary_contact_email, primary_contact_email),
  primary_contact_phone = COALESCE(@primary_contact_phone, primary_contact_phone),
  company_address = COALESCE(@company_address, company_address),
  company_city = COALESCE(@company_city, company_city),
  company_postal_code = COALESCE(@company_postal_code, company_postal_code),
  company_country = COALESCE(@company_country, company_country),
  company_registration_number = COALESCE(@company_registration_number, company_registration_number),
  registration_country = COALESCE(@registration_country, registration_country),
  tax_id = COALESCE(@tax_id, tax_id),
  incorporation_date = COALESCE(@incorporation_date, incorporation_date),
  account_status = COALESCE(@account_status, account_status),
  kyb_status = COALESCE(@kyb_status, kyb_status),
  kyb_verified_at = COALESCE(@kyb_verified_at, kyb_verified_at),
  kyb_verification_method = COALESCE(@kyb_verification_method, kyb_verification_method),
  kyb_verification_provider = COALESCE(@kyb_verification_provider, kyb_verification_provider),
  kyb_rejection_reason = COALESCE(@kyb_rejection_reason, kyb_rejection_reason),
  legal_entity_type = COALESCE(@legal_entity_type, legal_entity_type),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: CreateCompanyUser :one
INSERT INTO company_users (
  id,
  company_id,
  user_id,
  role,
  department,
  job_title,
  is_administrator,
  can_manage_payroll,
  can_manage_invoices,
  can_manage_employees,
  can_manage_company_settings,
  can_manage_bank_accounts,
  can_manage_wallets,
  permissions,
  is_active,
  added_by,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @user_id,
  @role,
  @department,
  @job_title,
  COALESCE(@is_administrator, FALSE),
  COALESCE(@can_manage_payroll, FALSE),
  COALESCE(@can_manage_invoices, FALSE),
  COALESCE(@can_manage_employees, FALSE),
  COALESCE(@can_manage_company_settings, FALSE),
  COALESCE(@can_manage_bank_accounts, FALSE),
  COALESCE(@can_manage_wallets, FALSE),
  @permissions,
  COALESCE(@is_active, TRUE),
  @added_by,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetCompanyUser :one
SELECT cu.*, u.email, u.account_status,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name
FROM company_users cu
JOIN users u ON cu.user_id = u.id
LEFT JOIN company_staff_profiles csp ON cu.id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE cu.company_id = @company_id AND cu.user_id = @user_id;

-- name: GetCompanyUsersByCompany :many
SELECT cu.*, u.email, u.account_status,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name
FROM company_users cu
JOIN users u ON cu.user_id = u.id
LEFT JOIN company_staff_profiles csp ON cu.id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE cu.company_id = @company_id AND cu.is_active = TRUE
ORDER BY cu.created_at DESC;

-- name: UpdateCompanyUser :one
UPDATE company_users SET
  role = COALESCE(@role, role),
  department = COALESCE(@department, department),
  job_title = COALESCE(@job_title, job_title),
  is_administrator = COALESCE(@is_administrator, is_administrator),
  can_manage_payroll = COALESCE(@can_manage_payroll, can_manage_payroll),
  can_manage_invoices = COALESCE(@can_manage_invoices, can_manage_invoices),
  can_manage_employees = COALESCE(@can_manage_employees, can_manage_employees),
  can_manage_company_settings = COALESCE(@can_manage_company_settings, can_manage_company_settings),
  can_manage_bank_accounts = COALESCE(@can_manage_bank_accounts, can_manage_bank_accounts),
  can_manage_wallets = COALESCE(@can_manage_wallets, can_manage_wallets),
  permissions = COALESCE(@permissions, permissions),
  is_active = COALESCE(@is_active, is_active),
  updated_at = NOW()
WHERE company_id = @company_id AND user_id = @user_id
RETURNING *;

-- name: CreateCompanyEmployee :one
INSERT INTO company_employees (
  id,
  company_id,
  user_id,
  employee_id,
  department,
  position,
  employment_status,
  employment_type,
  start_date,
  end_date,
  manager_id,
  salary_amount,
  salary_currency,
  salary_frequency,
  hourly_rate,
  payment_method,
  payment_split,
  tax_information,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @user_id,
  @employee_id,
  @department,
  @position,
  COALESCE(@employment_status, 'active'),
  @employment_type,
  @start_date,
  @end_date,
  @manager_id,
  @salary_amount,
  @salary_currency,
  @salary_frequency,
  @hourly_rate,
  @payment_method,
  @payment_split,
  @tax_information,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetCompanyEmployeeByID :one
SELECT ce.*, 
       u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name
FROM company_employees ce
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_staff_profiles csp ON ce.user_id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE ce.id = @id;

-- name: GetCompanyEmployeesByCompany :many
SELECT ce.*, 
       u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name
FROM company_employees ce
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_staff_profiles csp ON ce.user_id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE ce.company_id = @company_id
ORDER BY ce.created_at DESC;

-- name: UpdateCompanyEmployee :one
UPDATE company_employees SET
  employee_id = COALESCE(@employee_id, employee_id),
  department = COALESCE(@department, department),
  position = COALESCE(@position, position),
  employment_status = COALESCE(@employment_status, employment_status),
  employment_type = COALESCE(@employment_type, employment_type),
  start_date = COALESCE(@start_date, start_date),
  end_date = COALESCE(@end_date, end_date),
  manager_id = COALESCE(@manager_id, manager_id),
  salary_amount = COALESCE(@salary_amount, salary_amount),
  salary_currency = COALESCE(@salary_currency, salary_currency),
  salary_frequency = COALESCE(@salary_frequency, salary_frequency),
  hourly_rate = COALESCE(@hourly_rate, hourly_rate),
  payment_method = COALESCE(@payment_method, payment_method),
  payment_split = COALESCE(@payment_split, payment_split),
  tax_information = COALESCE(@tax_information, tax_information),
  updated_at = NOW()
WHERE id = @id
RETURNING *;