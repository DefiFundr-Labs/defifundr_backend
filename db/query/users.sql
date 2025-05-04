-- name: CreateUser :one
-- Creates a new user record and returns the created user
INSERT INTO users (
  id,
  email,
  password_hash,
  profile_picture,
  account_type,
  gender,
  personal_account_type,
  first_name,
  last_name,
  nationality,
  residential_country,
  job_role,
  company_name,
  company_size,
  company_industry,
  company_description,
  company_headquarters,
  auth_provider,
  provider_id,
  employee_type,
  company_website,
  employment_type,
  user_address,
  user_city,
  user_postal_code,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @email,
  @password_hash,
  COALESCE(@profile_picture, ''),
  @account_type,
  @gender,
  @personal_account_type,
  @first_name,
  @last_name,
  @nationality,
  @residential_country,
  @job_role,
  COALESCE(@company_name, ''),
  COALESCE(@company_size, ''),
  COALESCE(@company_industry, ''),
  COALESCE(@company_description, ''),
  COALESCE(@company_headquarters, ''),
  @auth_provider,
  @provider_id,
  @employee_type,
  COALESCE(@company_website, ''),
  COALESCE(@employment_type, ''),
  COALESCE(@user_address, ''),
  COALESCE(@user_city, ''),
  COALESCE(@user_postal_code, ''),
  COALESCE(@created_at, now()),
  COALESCE(@updated_at, now())
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1::uuid LIMIT 1;

-- name: GetUserByEmail :one
-- Retrieves a single user by their email address
SELECT * FROM users
WHERE email = $1
LIMIT 1;

-- name: CheckEmailExists :one
SELECT EXISTS (
        SELECT 1
        FROM users
        WHERE
            email = $1
        LIMIT 1
    ) AS exists;

-- name: ListUsers :many
-- Lists users with pagination support
SELECT *
FROM users
ORDER BY 
  CASE WHEN $3::text = 'ASC' THEN created_at END ASC,
  CASE WHEN $3::text = 'DESC' OR $3::text IS NULL THEN created_at END DESC
LIMIT $1
OFFSET $2;

-- name: ListUsersByAccountType :many
-- Lists users filtered by account type with pagination
SELECT *
FROM users
WHERE account_type = $3
ORDER BY 
  CASE WHEN $4::text = 'ASC' THEN created_at END ASC,
  CASE WHEN $4::text = 'DESC' OR $4::text IS NULL THEN created_at END DESC
LIMIT $1
OFFSET $2;

-- name: CountUsers :one
-- Counts the total number of users (useful for pagination)
SELECT COUNT(*) FROM users;

-- name: CountUsersByAccountType :one
-- Counts users filtered by account type
SELECT COUNT(*) FROM users
WHERE account_type = $1;



-- name: UpdateUser :one
-- Updates user details and returns the updated user
UPDATE users
SET
  email = COALESCE($2, email),
  profile_picture = $3,
  account_type = COALESCE($4, account_type),
  gender = $5,
  personal_account_type = COALESCE($6, personal_account_type),
  first_name = COALESCE($7, first_name),
  last_name = COALESCE($8, last_name),
  nationality = COALESCE($9, nationality),
  residential_country = $10,
  job_role = $11,
  company_website = $12,
  employment_type = $13,
  company_name = $14,
  company_headquarters = $15,
  company_size = $16,
  company_description = $17,
  company_industry = $18,
  auth_provider = COALESCE($19, auth_provider),
  provider_id = COALESCE($20, provider_id),
  user_address = COALESCE($21, user_address),
  user_city = COALESCE($22, user_city),
  user_postal_code = COALESCE($23, user_postal_code),
  updated_at = now()
WHERE id = $1
RETURNING *;


-- name: UpdateUserProfile :one
-- Updates a user's profile information
UPDATE users
SET
  profile_picture = COALESCE($2, profile_picture),
  first_name = COALESCE($3, first_name),
  last_name = COALESCE($4, last_name)
WHERE id = $1
RETURNING *;

-- name: UpdateUserPersonalDetails :one
-- Updates a user's personal details
UPDATE users
SET
  nationality = COALESCE($2, nationality),
  phone_number = COALESCE($3, phone_number),
  residential_country = COALESCE($4, residential_country),
  account_type = COALESCE($5, account_type),
  personal_account_type = COALESCE($6, personal_account_type),
  updated_at = now()
  WHERE id = $1
  RETURNING *;


-- name: UpdateUserAddress :one
-- Updates a user's address
UPDATE users
SET
  user_address = COALESCE($2, user_address),
  user_city = COALESCE($3, user_city),
  user_postal_code = COALESCE($4, user_postal_code),
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: UpdateUserCompanyDetails :one
-- Updates a user's company details
UPDATE users
SET
  company_name = COALESCE($2, company_name),
  company_headquarters = COALESCE($3, company_headquarters),
  company_size = COALESCE($4, company_size),
  company_industry = COALESCE($5, company_industry),
  company_description = COALESCE($6, company_description),
  company_headquarters = COALESCE($7, company_headquarters),
  account_type = COALESCE($8, account_type),
  updated_at = now()
WHERE id = $1
RETURNING *;


-- name: UpdateUserJobRole :one
-- Updates a user's job role
UPDATE users
SET
  job_role = COALESCE($2, job_role),
  updated_at = now()
WHERE id = $1
RETURNING *;


-- name: UpdateUserPassword :exec
-- Updates a user's password
UPDATE users
SET
  password_hash = $2,
  updated_at = now()
WHERE id = $1;

-- name: UpdateUserEmail :one
-- Updates a user's email address with validation that the new email is unique
UPDATE users
SET
  email = $2,
  updated_at = now()
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
-- Permanently deletes a user record
DELETE FROM users
WHERE id = $1;

-- name: SearchUsers :many
-- Searches for users by name, email, or nationality with pagination
SELECT *
FROM users
WHERE 
  (
    first_name ILIKE '%' || $3 || '%' OR
    last_name ILIKE '%' || $3 || '%' OR
    email ILIKE '%' || $3 || '%' OR
    nationality ILIKE '%' || $3 || '%'
  )
ORDER BY
  CASE WHEN $4::text = 'ASC' THEN created_at END ASC,
  CASE WHEN $4::text = 'DESC' OR $4::text IS NULL THEN created_at END DESC
LIMIT $1
OFFSET $2;

-- name: CountSearchUsers :one
-- Counts the number of users matching a search query
SELECT COUNT(*)
FROM users
WHERE 
  (
    first_name ILIKE '%' || $1 || '%' OR
    last_name ILIKE '%' || $1 || '%' OR
    email ILIKE '%' || $1 || '%' OR
    nationality ILIKE '%' || $1 || '%'
  );