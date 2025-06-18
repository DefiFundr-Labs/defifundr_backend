-- name: CreateUser :one
INSERT INTO users (
  id,
  email,
  password_hash,
  auth_provider,
  provider_id,
  email_verified,
  email_verified_at,
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
  @email,
  @password_hash,
  @auth_provider,
  @provider_id,
  COALESCE(@email_verified, FALSE),
  @email_verified_at,
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
  password_hash = COALESCE(@password_hash, password_hash),
  auth_provider = COALESCE(@auth_provider, auth_provider),
  provider_id = COALESCE(@provider_id, provider_id),
  email_verified = COALESCE(@email_verified, email_verified),
  email_verified_at = COALESCE(@email_verified_at, email_verified_at),
  account_status = COALESCE(@account_status, account_status),
  two_factor_enabled = COALESCE(@two_factor_enabled, two_factor_enabled),
  two_factor_method = COALESCE(@two_factor_method, two_factor_method),
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

-- name: CreatePersonalUser :one
INSERT INTO personal_users (
  id,
  first_name,
  last_name,
  profile_picture,
  phone_number,
  phone_number_verified,
  phone_number_verified_at,
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
  @id,
  @first_name,
  @last_name,
  @profile_picture,
  @phone_number,
  COALESCE(@phone_number_verified, FALSE),
  @phone_number_verified_at,
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
SELECT pu.*, u.email, u.account_status
FROM personal_users pu
JOIN users u ON pu.id = u.id
WHERE pu.id = @id AND u.deleted_at IS NULL;

-- name: UpdatePersonalUser :one
UPDATE personal_users SET
  first_name = COALESCE(@first_name, first_name),
  last_name = COALESCE(@last_name, last_name),
  profile_picture = COALESCE(@profile_picture, profile_picture),
  phone_number = COALESCE(@phone_number, phone_number),
  phone_number_verified = COALESCE(@phone_number_verified, phone_number_verified),
  phone_number_verified_at = COALESCE(@phone_number_verified_at, phone_number_verified_at),
  nationality = COALESCE(@nationality, nationality),
  residential_country = COALESCE(@residential_country, residential_country),
  user_address = COALESCE(@user_address, user_address),
  user_city = COALESCE(@user_city, user_city),
  user_postal_code = COALESCE(@user_postal_code, user_postal_code),
  gender = COALESCE(@gender, gender),
  date_of_birth = COALESCE(@date_of_birth, date_of_birth),
  job_role = COALESCE(@job_role, job_role),
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