-- name: CreateSecurityEvent :one
INSERT INTO security_events (
    id, user_id, event_type, ip_address, user_agent, metadata, timestamp
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetRecentLoginEventsByUserID :many
SELECT * FROM security_events
WHERE user_id = $1 AND event_type = 'login'
ORDER BY timestamp DESC
LIMIT $2;

-- name: GetSecurityEventsByUserIDAndType :many
SELECT * FROM security_events
WHERE user_id = $1
  AND event_type = $2
  AND timestamp BETWEEN $3 AND $4
ORDER BY timestamp DESC;
