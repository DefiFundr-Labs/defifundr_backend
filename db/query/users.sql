-- ================================
-- USERS TABLE QUERIES
-- ================================

-- name: CreateUser :one
INSERT INTO users (
  id,
  first_name,
  last_name,
  phone_number,
  email,
  password_hash,
  profile_picture_url,
  auth_provider,
  provider_id,
  email_verified,
  email_verified_at,
  phone_number_verified,
  phone_number_verified_at,
  account_type,
  account_status,
  two_factor_enabled,
  two_factor_method,
  user_login_type,
  created_at,
  updated_at,
  last_login_at,
  deleted_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @first_name,
  @last_name,
  @phone_number,
  @email,
  @password_hash,
  @profile_picture_url,
  @auth_provider,
  @provider_id,
  COALESCE(@email_verified, FALSE),
  @email_verified_at,
  COALESCE(@phone_number_verified, FALSE),
  @phone_number_verified_at,
  @account_type,
  COALESCE(@account_status, 'pending'),
  COALESCE(@two_factor_enabled, FALSE),
  @two_factor_method,
  @user_login_type,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW()),
  @last_login_at,
  @deleted_at
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = @id AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM users 
WHERE email = @email AND deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE users SET
  first_name = COALESCE(@first_name, first_name),
  last_name = COALESCE(@last_name, last_name),
  phone_number = COALESCE(@phone_number, phone_number),
  password_hash = COALESCE(@password_hash, password_hash),
  profile_picture_url = COALESCE(@profile_picture_url, profile_picture_url),
  auth_provider = COALESCE(@auth_provider, auth_provider),
  provider_id = COALESCE(@provider_id, provider_id),
  email_verified = COALESCE(@email_verified, email_verified),
  email_verified_at = COALESCE(@email_verified_at, email_verified_at),
  phone_number_verified = COALESCE(@phone_number_verified, phone_number_verified),
  phone_number_verified_at = COALESCE(@phone_number_verified_at, phone_number_verified_at),
  account_status = COALESCE(@account_status, account_status),
  two_factor_enabled = COALESCE(@two_factor_enabled, two_factor_enabled),
  two_factor_method = COALESCE(@two_factor_method, two_factor_method),
  user_login_type = COALESCE(@user_login_type, user_login_type),
  last_login_at = COALESCE(@last_login_at, last_login_at),
  updated_at = NOW()
WHERE id = @id AND deleted_at IS NULL
RETURNING *;

-- name: UpdateUserLoginTime :exec
UPDATE users SET
  last_login_at = NOW(),
  updated_at = NOW()
WHERE id = @id;

-- name: SoftDeleteUser :exec
UPDATE users SET
  deleted_at = NOW(),
  updated_at = NOW()
WHERE id = @id;

-- name: ListUsers :many
SELECT * FROM users 
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: SearchUsersByEmail :many
SELECT * FROM users 
WHERE email ILIKE '%' || @search_term || '%'
  AND deleted_at IS NULL
ORDER BY email
LIMIT @limit_val OFFSET @offset_val;

-- name: GetUsersByAccountType :many
SELECT * FROM users 
WHERE account_type = @account_type AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: CountUsers :one
SELECT COUNT(*) FROM users WHERE deleted_at IS NULL;

-- name: CountUsersByAccountType :one
SELECT COUNT(*) FROM users 
WHERE account_type = @account_type AND deleted_at IS NULL;

-- ================================
-- PERSONAL_USERS TABLE QUERIES
-- ================================

-- name: CreatePersonalUser :one
INSERT INTO personal_users (
  id,
  user_id,
  nationality,
  residential_country,
  user_address,
  user_city,
  user_postal_code,
  gender,
  date_of_birth,
  job_role,
  personal_account_type,
  employment_type,
  tax_id,
  default_payment_currency,
  default_payment_method,
  hourly_rate,
  specialization,
  kyc_status,
  kyc_verified_at,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @nationality,
  @residential_country,
  @user_address,
  @user_city,
  @user_postal_code,
  @gender,
  @date_of_birth,
  @job_role,
  @personal_account_type,
  @employment_type,
  @tax_id,
  @default_payment_currency,
  @default_payment_method,
  @hourly_rate,
  @specialization,
  COALESCE(@kyc_status, 'pending'),
  @kyc_verified_at,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetPersonalUserByID :one
SELECT * FROM personal_users WHERE id = @id;

-- name: GetPersonalUserByUserID :one
SELECT * FROM personal_users WHERE user_id = @user_id;

-- name: GetPersonalUserWithUserDetails :one
SELECT 
  pu.*,
  u.first_name,
  u.last_name,
  u.email,
  u.phone_number,
  u.profile_picture_url,
  u.account_status,
  u.email_verified,
  u.phone_number_verified
FROM personal_users pu
JOIN users u ON pu.user_id = u.id
WHERE pu.id = @id AND u.deleted_at IS NULL;

-- name: UpdatePersonalUser :one
UPDATE personal_users SET
  nationality = COALESCE(@nationality, nationality),
  residential_country = COALESCE(@residential_country, residential_country),
  user_address = COALESCE(@user_address, user_address),
  user_city = COALESCE(@user_city, user_city),
  user_postal_code = COALESCE(@user_postal_code, user_postal_code),
  gender = COALESCE(@gender, gender),
  date_of_birth = COALESCE(@date_of_birth, date_of_birth),
  job_role = COALESCE(@job_role, job_role),
  personal_account_type = COALESCE(@personal_account_type, personal_account_type),
  employment_type = COALESCE(@employment_type, employment_type),
  tax_id = COALESCE(@tax_id, tax_id),
  default_payment_currency = COALESCE(@default_payment_currency, default_payment_currency),
  default_payment_method = COALESCE(@default_payment_method, default_payment_method),
  hourly_rate = COALESCE(@hourly_rate, hourly_rate),
  specialization = COALESCE(@specialization, specialization),
  kyc_status = COALESCE(@kyc_status, kyc_status),
  kyc_verified_at = COALESCE(@kyc_verified_at, kyc_verified_at),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeletePersonalUser :exec
DELETE FROM personal_users WHERE id = @id;

-- name: ListPersonalUsers :many
SELECT 
  pu.*,
  u.first_name,
  u.last_name,
  u.email,
  u.account_status
FROM personal_users pu
JOIN users u ON pu.user_id = u.id
WHERE u.deleted_at IS NULL
ORDER BY pu.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetPersonalUsersByKYCStatus :many
SELECT * FROM personal_users 
WHERE kyc_status = @kyc_status
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- ================================
-- COMPANIES TABLE QUERIES
-- ================================

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

-- name: GetCompanyByOwnerID :one
SELECT * FROM companies WHERE owner_id = @owner_id;

-- name: GetCompanyWithOwnerDetails :one
SELECT 
  c.*,
  u.first_name as owner_first_name,
  u.last_name as owner_last_name,
  u.email as owner_email
FROM companies c
JOIN users u ON c.owner_id = u.id
WHERE c.id = @id AND u.deleted_at IS NULL;

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
  legal_entity_type = COALESCE(@legal_entity_type, legal_entity_type),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: UpdateCompanyKYB :one
UPDATE companies SET
  kyb_status = @kyb_status,
  kyb_verified_at = @kyb_verified_at,
  kyb_verification_method = @kyb_verification_method,
  kyb_verification_provider = @kyb_verification_provider,
  kyb_rejection_reason = @kyb_rejection_reason,
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteCompany :exec
DELETE FROM companies WHERE id = @id;

-- name: ListCompanies :many
SELECT 
  c.*,
  u.first_name as owner_first_name,
  u.last_name as owner_last_name,
  u.email as owner_email
FROM companies c
JOIN users u ON c.owner_id = u.id
WHERE u.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetCompaniesByKYBStatus :many
SELECT * FROM companies 
WHERE kyb_status = @kyb_status
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: SearchCompaniesByName :many
SELECT * FROM companies 
WHERE company_name ILIKE '%' || @search_term || '%'
ORDER BY company_name
LIMIT @limit_val OFFSET @offset_val;

-- name: GetCompaniesByIndustry :many
SELECT * FROM companies 
WHERE company_industry = @company_industry
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- ================================
-- COMPANY_USERS TABLE QUERIES
-- ================================

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
  reports_to,
  hire_date,
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
  @reports_to,
  @hire_date,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetCompanyUserByID :one
SELECT * FROM company_users WHERE id = @id;

-- name: GetCompanyUserByCompanyAndUser :one
SELECT * FROM company_users 
WHERE company_id = @company_id AND user_id = @user_id;

-- name: GetCompanyUserWithDetails :one
SELECT 
  cu.*,
  u.first_name,
  u.last_name,
  u.email,
  u.phone_number,
  c.company_name
FROM company_users cu
JOIN users u ON cu.user_id = u.id
JOIN companies c ON cu.company_id = c.id
WHERE cu.id = @id AND u.deleted_at IS NULL;

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
  reports_to = COALESCE(@reports_to, reports_to),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeactivateCompanyUser :one
UPDATE company_users SET
  is_active = FALSE,
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteCompanyUser :exec
DELETE FROM company_users WHERE id = @id;

-- name: ListCompanyUsers :many
SELECT 
  cu.*,
  u.first_name,
  u.last_name,
  u.email,
  u.phone_number
FROM company_users cu
JOIN users u ON cu.user_id = u.id
WHERE cu.company_id = @company_id AND u.deleted_at IS NULL
ORDER BY cu.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetCompanyAdministrators :many
SELECT 
  cu.*,
  u.first_name,
  u.last_name,
  u.email
FROM company_users cu
JOIN users u ON cu.user_id = u.id
WHERE cu.company_id = @company_id 
  AND cu.is_administrator = TRUE 
  AND cu.is_active = TRUE
  AND u.deleted_at IS NULL;

-- name: GetUserCompanies :many
SELECT 
  cu.*,
  c.company_name,
  c.company_logo,
  c.company_industry
FROM company_users cu
JOIN companies c ON cu.company_id = c.id
WHERE cu.user_id = @user_id AND cu.is_active = TRUE
ORDER BY cu.created_at DESC;

-- ================================
-- COMPANY_EMPLOYEES TABLE QUERIES
-- ================================

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
SELECT * FROM company_employees WHERE id = @id;

-- name: GetCompanyEmployeeByEmployeeID :one
SELECT * FROM company_employees 
WHERE company_id = @company_id AND employee_id = @employee_id;

-- name: GetCompanyEmployeeWithDetails :one
SELECT 
  ce.*,
  u.first_name,
  u.last_name,
  u.email,
  u.phone_number,
  c.company_name,
  m.first_name as manager_first_name,
  m.last_name as manager_last_name
FROM company_employees ce
JOIN companies c ON ce.company_id = c.id
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_users cu ON ce.manager_id = cu.id
LEFT JOIN users m ON cu.user_id = m.id
WHERE ce.id = @id;

-- name: UpdateCompanyEmployee :one
UPDATE company_employees SET
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

-- name: UpdateEmployeeStatus :one
UPDATE company_employees SET
  employment_status = @employment_status,
  end_date = COALESCE(@end_date, end_date),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteCompanyEmployee :exec
DELETE FROM company_employees WHERE id = @id;

-- name: ListCompanyEmployees :many
SELECT 
  ce.*,
  u.first_name,
  u.last_name,
  u.email
FROM company_employees ce
LEFT JOIN users u ON ce.user_id = u.id
WHERE ce.company_id = @company_id
ORDER BY ce.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetActiveEmployees :many
SELECT 
  ce.*,
  u.first_name,
  u.last_name,
  u.email
FROM company_employees ce
LEFT JOIN users u ON ce.user_id = u.id
WHERE ce.company_id = @company_id 
  AND ce.employment_status = 'active'
ORDER BY ce.start_date DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetEmployeesByDepartment :many
SELECT 
  ce.*,
  u.first_name,
  u.last_name,
  u.email
FROM company_employees ce
LEFT JOIN users u ON ce.user_id = u.id
WHERE ce.company_id = @company_id 
  AND ce.department = @department
ORDER BY ce.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetEmployeesByManager :many
SELECT 
  ce.*,
  u.first_name,
  u.last_name,
  u.email
FROM company_employees ce
LEFT JOIN users u ON ce.user_id = u.id
WHERE ce.manager_id = @manager_id
ORDER BY ce.start_date DESC;

-- name: CountCompanyEmployees :one
SELECT COUNT(*) FROM company_employees 
WHERE company_id = @company_id;

-- name: CountActiveEmployees :one
SELECT COUNT(*) FROM company_employees 
WHERE company_id = @company_id AND employment_status = 'active';