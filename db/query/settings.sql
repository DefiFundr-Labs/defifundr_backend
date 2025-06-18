-- name: CreateSystemSetting :one
INSERT INTO system_settings (
  id,
  setting_key,
  setting_value,
  data_type,
  description,
  is_sensitive,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @setting_key,
  @setting_value,
  COALESCE(@data_type, 'string'),
  @description,
  COALESCE(@is_sensitive, FALSE),
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetSystemSetting :one
SELECT * FROM system_settings 
WHERE setting_key = @setting_key;

-- name: GetAllSystemSettings :many
SELECT * FROM system_settings 
ORDER BY setting_key;

-- name: GetPublicSystemSettings :many
SELECT setting_key, setting_value, data_type, description 
FROM system_settings 
WHERE is_sensitive = FALSE
ORDER BY setting_key;

-- name: UpdateSystemSetting :one
UPDATE system_settings SET
  setting_value = @setting_value,
  data_type = COALESCE(@data_type, data_type),
  description = COALESCE(@description, description),
  is_sensitive = COALESCE(@is_sensitive, is_sensitive),
  updated_at = NOW()
WHERE setting_key = @setting_key
RETURNING *;

-- name: DeleteSystemSetting :exec
DELETE FROM system_settings WHERE setting_key = @setting_key;

-- name: CreateCompanySetting :one
INSERT INTO company_settings (
  id,
  company_id,
  setting_key,
  setting_value,
  data_type,
  description,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @setting_key,
  @setting_value,
  COALESCE(@data_type, 'string'),
  @description,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetCompanySetting :one
SELECT * FROM company_settings 
WHERE company_id = @company_id AND setting_key = @setting_key;

-- name: GetCompanySettings :many
SELECT * FROM company_settings 
WHERE company_id = @company_id
ORDER BY setting_key;

-- name: UpdateCompanySetting :one
UPDATE company_settings SET
  setting_value = @setting_value,
  data_type = COALESCE(@data_type, data_type),
  description = COALESCE(@description, description),
  updated_at = NOW()
WHERE company_id = @company_id AND setting_key = @setting_key
RETURNING *;

-- name: DeleteCompanySetting :exec
DELETE FROM company_settings 
WHERE company_id = @company_id AND setting_key = @setting_key;

-- name: CreateUserSetting :one
INSERT INTO user_settings (
  id,
  user_id,
  setting_key,
  setting_value,
  data_type,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @setting_key,
  @setting_value,
  COALESCE(@data_type, 'string'),
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetUserSetting :one
SELECT * FROM user_settings 
WHERE user_id = @user_id AND setting_key = @setting_key;

-- name: GetUserSettings :many
SELECT * FROM user_settings 
WHERE user_id = @user_id
ORDER BY setting_key;

-- name: UpdateUserSetting :one
UPDATE user_settings SET
  setting_value = @setting_value,
  data_type = COALESCE(@data_type, data_type),
  updated_at = NOW()
WHERE user_id = @user_id AND setting_key = @setting_key
RETURNING *;

-- name: DeleteUserSetting :exec
DELETE FROM user_settings 
WHERE user_id = @user_id AND setting_key = @setting_key;

-- name: CreateFeatureFlag :one
INSERT INTO feature_flags (
  id,
  flag_key,
  description,
  is_enabled,
  rollout_percentage,
  conditions,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @flag_key,
  @description,
  COALESCE(@is_enabled, FALSE),
  COALESCE(@rollout_percentage, 0),
  @conditions,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetFeatureFlag :one
SELECT * FROM feature_flags WHERE flag_key = @flag_key;

-- name: GetAllFeatureFlags :many
SELECT * FROM feature_flags ORDER BY flag_key;

-- name: UpdateFeatureFlag :one
UPDATE feature_flags SET
  description = COALESCE(@description, description),
  is_enabled = COALESCE(@is_enabled, is_enabled),
  rollout_percentage = COALESCE(@rollout_percentage, rollout_percentage),
  conditions = COALESCE(@conditions, conditions),
  updated_at = NOW()
WHERE flag_key = @flag_key
RETURNING *;

-- name: CreateUserFeatureFlag :one
INSERT INTO user_feature_flags (
  id,
  user_id,
  flag_key,
  is_enabled,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @flag_key,
  @is_enabled,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetUserFeatureFlag :one
SELECT * FROM user_feature_flags 
WHERE user_id = @user_id AND flag_key = @flag_key;

-- name: GetUserFeatureFlags :many
SELECT uff.*, ff.description
FROM user_feature_flags uff
JOIN feature_flags ff ON uff.flag_key = ff.flag_key
WHERE uff.user_id = @user_id
ORDER BY uff.flag_key;

-- name: UpdateUserFeatureFlag :one
UPDATE user_feature_flags SET
  is_enabled = @is_enabled,
  updated_at = NOW()
WHERE user_id = @user_id AND flag_key = @flag_key
RETURNING *;

-- name: CreateCompanyFeatureFlag :one
INSERT INTO company_feature_flags (
  id,
  company_id,
  flag_key,
  is_enabled,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @flag_key,
  @is_enabled,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetCompanyFeatureFlag :one
SELECT * FROM company_feature_flags 
WHERE company_id = @company_id AND flag_key = @flag_key;

-- name: GetCompanyFeatureFlags :many
SELECT cff.*, ff.description
FROM company_feature_flags cff
JOIN feature_flags ff ON cff.flag_key = ff.flag_key
WHERE cff.company_id = @company_id
ORDER BY cff.flag_key;

-- name: UpdateCompanyFeatureFlag :one
UPDATE company_feature_flags SET
  is_enabled = @is_enabled,
  updated_at = NOW()
WHERE company_id = @company_id AND flag_key = @flag_key
RETURNING *;