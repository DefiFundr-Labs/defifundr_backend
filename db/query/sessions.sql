-- name: CreateSession :one
-- Creates a new session and returns the created session record
INSERT INTO sessions (
  id,
  user_id,
  refresh_token,
  user_agent,
  web_oauth_client_id,
  oauth_access_token,
  oauth_id_token,
  user_login_type,
  last_used_at,
  mfa_enabled,
  client_ip,
  is_blocked,
  expires_at,
  created_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, COALESCE($14, now())
) RETURNING *;

-- name: GetSessionByID :one
-- Retrieves a session by its ID
SELECT * FROM sessions
WHERE id = $1
LIMIT 1;

-- name: GetSessionByRefreshToken :one
-- Retrieves a session by refresh token
SELECT * FROM sessions
WHERE refresh_token = $1
LIMIT 1;

-- name: GetSessionsByUserID :many
-- Retrieves all sessions for a specific user
SELECT * FROM sessions
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: GetActiveSessions :many
-- Retrieves active (non-expired, non-blocked) sessions with pagination
SELECT * FROM sessions
WHERE is_blocked = false AND expires_at > now()
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: GetActiveSessionsByUserID :many
-- Retrieves active sessions for a specific user
SELECT * FROM sessions
WHERE user_id = $1 AND is_blocked = false AND expires_at > now()
ORDER BY created_at DESC;

-- name: UpdateSession :one
-- Updates session details
UPDATE sessions
SET
  user_agent = COALESCE($2, user_agent),
  web_oauth_client_id = $3,
  oauth_access_token = $4,
  oauth_id_token = $5,
  mfa_enabled = COALESCE($6, mfa_enabled),
  client_ip = COALESCE($7, client_ip),
  is_blocked = COALESCE($8, is_blocked),
  expires_at = COALESCE($9, expires_at),
  last_used_at = COALESCE($10, last_used_at)
WHERE id = $1
RETURNING *;

-- name: UpdateRefreshToken :one
-- Updates just the refresh token of a session
UPDATE sessions
SET refresh_token = $2
WHERE id = $1
RETURNING *;

-- name: BlockSession :exec
-- Blocks a session (marks it as invalid)
UPDATE sessions
SET is_blocked = true
WHERE id = $1;

-- name: BlockAllUserSessions :exec
-- Blocks all sessions for a specific user
UPDATE sessions
SET is_blocked = true
WHERE user_id = $1;

-- name: BlockExpiredSessions :exec
-- Blocks all expired sessions
UPDATE sessions
SET is_blocked = true
WHERE expires_at <= now() AND is_blocked = false;

-- name: DeleteSession :exec
-- Deletes a session by its ID
DELETE FROM sessions
WHERE id = $1;

-- name: DeleteSessionsByUserID :exec
-- Deletes all sessions for a specific user
DELETE FROM sessions
WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
-- Cleans up expired sessions that are older than the specified date
DELETE FROM sessions
WHERE expires_at < $1;

-- name: CountActiveSessions :one
-- Counts the number of active sessions
SELECT COUNT(*) FROM sessions
WHERE is_blocked = false AND expires_at > now();

-- name: CountActiveSessionsByUserID :one
-- Counts the number of active sessions for a specific user
SELECT COUNT(*) FROM sessions
WHERE user_id = $1 AND is_blocked = false AND expires_at > now();

-- name: UpdateSessionRefreshToken :one
UPDATE sessions
SET 
  refresh_token = $2,
  last_used_at = $3
WHERE id = $1
RETURNING *;