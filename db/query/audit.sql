-- name: CreateAuditLog :one
INSERT INTO audit_logs (
  id,
  user_id,
  company_id,
  action,
  entity_type,
  entity_id,
  previous_state,
  new_state,
  ip_address,
  user_agent,
  created_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @company_id,
  @action,
  @entity_type,
  @entity_id,
  @previous_state,
  @new_state,
  @ip_address,
  @user_agent,
  COALESCE(@created_at, NOW())
) RETURNING *;

-- name: GetAuditLogsByUser :many
SELECT al.*, u.email
FROM audit_logs al
LEFT JOIN users u ON al.user_id = u.id
WHERE al.user_id = @user_id
ORDER BY al.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetAuditLogsByCompany :many
SELECT al.*, u.email
FROM audit_logs al
LEFT JOIN users u ON al.user_id = u.id
WHERE al.company_id = @company_id
ORDER BY al.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetAuditLogsByEntity :many
SELECT al.*, u.email
FROM audit_logs al
LEFT JOIN users u ON al.user_id = u.id
WHERE al.entity_type = @entity_type AND al.entity_id = @entity_id
ORDER BY al.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetAuditLogsByAction :many
SELECT al.*, u.email
FROM audit_logs al
LEFT JOIN users u ON al.user_id = u.id
WHERE al.action = @action
  AND (@user_id::uuid IS NULL OR al.user_id = @user_id)
  AND (@company_id::uuid IS NULL OR al.company_id = @company_id)
  AND (@start_date::timestamptz IS NULL OR al.created_at >= @start_date)
  AND (@end_date::timestamptz IS NULL OR al.created_at <= @end_date)
ORDER BY al.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: SearchAuditLogs :many
SELECT al.*, u.email
FROM audit_logs al
LEFT JOIN users u ON al.user_id = u.id
WHERE (@user_id::uuid IS NULL OR al.user_id = @user_id)
  AND (@company_id::uuid IS NULL OR al.company_id = @company_id)
  AND (@entity_type::text IS NULL OR al.entity_type = @entity_type)
  AND (@action::text IS NULL OR al.action = @action)
  AND (@start_date::timestamptz IS NULL OR al.created_at >= @start_date)
  AND (@end_date::timestamptz IS NULL OR al.created_at <= @end_date)
ORDER BY al.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: CreateActivityLog :one
INSERT INTO activity_logs (
  id,
  user_id,
  activity_type,
  description,
  metadata,
  ip_address,
  user_agent,
  created_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @activity_type,
  @description,
  @metadata,
  @ip_address,
  @user_agent,
  COALESCE(@created_at, NOW())
) RETURNING *;

-- name: GetActivityLogsByUser :many
SELECT * FROM activity_logs 
WHERE user_id = @user_id
ORDER BY created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetActivityLogsByType :many
SELECT al.*, u.email
FROM activity_logs al
JOIN users u ON al.user_id = u.id
WHERE al.activity_type = @activity_type
  AND (@user_id::uuid IS NULL OR al.user_id = @user_id)
  AND (@start_date::timestamptz IS NULL OR al.created_at >= @start_date)
  AND (@end_date::timestamptz IS NULL OR al.created_at <= @end_date)
ORDER BY al.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetRecentActivity :many
SELECT al.*, u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name
FROM activity_logs al
JOIN users u ON al.user_id = u.id
LEFT JOIN company_staff_profiles csp ON u.id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE (@user_id::uuid IS NULL OR al.user_id = @user_id)
  AND al.created_at >= NOW() - INTERVAL '@hours hours'
ORDER BY al.created_at DESC
LIMIT @limit_val;

-- name: CleanupOldAuditLogs :exec
DELETE FROM audit_logs 
WHERE created_at < NOW() - INTERVAL '@days days';

-- name: CleanupOldActivityLogs :exec
DELETE FROM activity_logs 
WHERE created_at < NOW() - INTERVAL '@days days';

-- name: GetUserActivitySummary :one
SELECT 
  COUNT(*) as total_activities,
  COUNT(DISTINCT activity_type) as unique_activity_types,
  MAX(created_at) as last_activity,
  MIN(created_at) as first_activity
FROM activity_logs 
WHERE user_id = @user_id
  AND (@start_date::timestamptz IS NULL OR created_at >= @start_date)
  AND (@end_date::timestamptz IS NULL OR created_at <= @end_date);

-- name: GetCompanyAuditSummary :one
SELECT 
  COUNT(*) as total_audit_logs,
  COUNT(DISTINCT user_id) as unique_users,
  COUNT(DISTINCT action) as unique_actions,
  COUNT(DISTINCT entity_type) as unique_entity_types,
  MAX(created_at) as last_audit,
  MIN(created_at) as first_audit
FROM audit_logs 
WHERE company_id = @company_id
  AND (@start_date::timestamptz IS NULL OR created_at >= @start_date)
  AND (@end_date::timestamptz IS NULL OR created_at <= @end_date);