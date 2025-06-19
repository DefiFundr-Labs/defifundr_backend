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