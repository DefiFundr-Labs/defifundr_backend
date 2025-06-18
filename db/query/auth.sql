-- name: CreateSession :one
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  user_agent,
  client_ip,
  last_used_at,
  web_oauth_client_id,
  oauth_access_token,
  oauth_id_token,
  user_login_type,
  mfa_verified,
  is_blocked,
  expires_at,
  created_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @refresh_token,
  @user_agent,
  @client_ip,
  COALESCE(@last_used_at, NOW()),
  @web_oauth_client_id,
  @oauth_access_token,
  @oauth_id_token,
  @user_login_type,
  COALESCE(@mfa_verified, FALSE),
  COALESCE(@is_blocked, FALSE),
  @expires_at,
  COALESCE(@created_at, NOW())
) RETURNING *;

-- name: GetSessionByID :one
SELECT * FROM sessions 
WHERE id = @id AND is_blocked = FALSE;

-- name: GetSessionsByUser :many
SELECT * FROM sessions 
WHERE user_id = @user_id AND is_blocked = FALSE
ORDER BY last_used_at DESC;

-- name: UpdateSessionLastUsed :exec
UPDATE sessions SET
  last_used_at = NOW()
WHERE id = @id;

-- name: RevokeSession :exec
UPDATE sessions SET
  is_blocked = TRUE
WHERE id = @id;

-- name: RevokeAllUserSessions :exec
UPDATE sessions SET
  is_blocked = TRUE
WHERE user_id = @user_id;

-- name: CleanupExpiredSessions :exec
DELETE FROM sessions 
WHERE expires_at < NOW() OR is_blocked = TRUE;

-- name: CreateUserDevice :one
INSERT INTO user_devices (
  id,
  user_id,
  device_token,
  platform,
  device_type,
  device_model,
  os_name,
  os_version,
  push_notification_token,
  is_active,
  is_verified,
  last_used_at,
  app_version,
  client_ip,
  expires_at,
  is_revoked,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @device_token,
  @platform,
  @device_type,
  @device_model,
  @os_name,
  @os_version,
  @push_notification_token,
  COALESCE(@is_active, TRUE),
  COALESCE(@is_verified, FALSE),
  COALESCE(@last_used_at, NOW()),
  @app_version,
  @client_ip,
  @expires_at,
  COALESCE(@is_revoked, FALSE),
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetUserDeviceByID :one
SELECT * FROM user_devices 
WHERE id = @id AND is_revoked = FALSE;

-- name: GetUserDevicesByUser :many
SELECT * FROM user_devices 
WHERE user_id = @user_id AND is_revoked = FALSE
ORDER BY last_used_at DESC;

-- name: GetUserDeviceByToken :one
SELECT * FROM user_devices 
WHERE device_token = @device_token AND is_revoked = FALSE;

-- name: UpdateUserDevice :one
UPDATE user_devices SET
  platform = COALESCE(@platform, platform),
  device_type = COALESCE(@device_type, device_type),
  device_model = COALESCE(@device_model, device_model),
  os_name = COALESCE(@os_name, os_name),
  os_version = COALESCE(@os_version, os_version),
  push_notification_token = COALESCE(@push_notification_token, push_notification_token),
  is_active = COALESCE(@is_active, is_active),
  is_verified = COALESCE(@is_verified, is_verified),
  last_used_at = NOW(),
  app_version = COALESCE(@app_version, app_version),
  client_ip = COALESCE(@client_ip, client_ip),
  expires_at = COALESCE(@expires_at, expires_at),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: RevokeUserDevice :exec
UPDATE user_devices SET
  is_revoked = TRUE,
  updated_at = NOW()
WHERE id = @id;

-- name: CreateSecurityEvent :one
INSERT INTO security_events (
  id,
  user_id,
  company_id,
  event_type,
  severity,
  ip_address,
  user_agent,
  metadata,
  created_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @company_id,
  @event_type,
  @severity,
  @ip_address,
  @user_agent,
  @metadata,
  COALESCE(@created_at, NOW())
) RETURNING *;

-- name: GetSecurityEventsByUser :many
SELECT * FROM security_events 
WHERE user_id = @user_id
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetSecurityEventsByCompany :many
SELECT * FROM security_events 
WHERE company_id = @company_id
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetSecurityEventsByType :many
SELECT * FROM security_events 
WHERE event_type = @event_type 
  AND (@user_id::uuid IS NULL OR user_id = @user_id)
  AND (@company_id::uuid IS NULL OR company_id = @company_id)
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;