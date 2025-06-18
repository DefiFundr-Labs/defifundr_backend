-- name: CreateNotificationTemplate :one
INSERT INTO notification_templates (
  id,
  template_name,
  template_type,
  subject,
  content,
  variables,
  is_active,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @template_name,
  @template_type,
  @subject,
  @content,
  @variables,
  COALESCE(@is_active, TRUE),
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetNotificationTemplateByName :one
SELECT * FROM notification_templates 
WHERE template_name = @template_name AND is_active = TRUE;

-- name: GetNotificationTemplatesByType :many
SELECT * FROM notification_templates 
WHERE template_type = @template_type AND is_active = TRUE
ORDER BY template_name;

-- name: CreateNotification :one
INSERT INTO notifications (
  id,
  user_id,
  template_id,
  notification_type,
  title,
  content,
  reference_type,
  reference_id,
  is_read,
  read_at,
  delivery_status,
  priority,
  created_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @template_id,
  @notification_type,
  @title,
  @content,
  @reference_type,
  @reference_id,
  COALESCE(@is_read, FALSE),
  @read_at,
  COALESCE(@delivery_status, 'pending'),
  COALESCE(@priority, 'normal'),
  COALESCE(@created_at, NOW())
) RETURNING *;

-- name: GetUserNotifications :many
SELECT n.*, nt.template_name
FROM notifications n
LEFT JOIN notification_templates nt ON n.template_id = nt.id
WHERE n.user_id = @user_id
ORDER BY n.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetUnreadNotifications :many
SELECT n.*, nt.template_name
FROM notifications n
LEFT JOIN notification_templates nt ON n.template_id = nt.id
WHERE n.user_id = @user_id AND n.is_read = FALSE
ORDER BY n.created_at DESC;

-- name: GetNotificationsByType :many
SELECT n.*, nt.template_name
FROM notifications n
LEFT JOIN notification_templates nt ON n.template_id = nt.id
WHERE n.user_id = @user_id AND n.notification_type = @notification_type
ORDER BY n.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: MarkNotificationAsRead :one
UPDATE notifications SET
  is_read = TRUE,
  read_at = NOW()
WHERE id = @id AND user_id = @user_id
RETURNING *;

-- name: MarkAllNotificationsAsRead :exec
UPDATE notifications SET
  is_read = TRUE,
  read_at = NOW()
WHERE user_id = @user_id AND is_read = FALSE;

-- name: DeleteNotification :exec
DELETE FROM notifications 
WHERE id = @id AND user_id = @user_id;

-- name: GetNotificationCount :one
SELECT 
  COUNT(*) as total_count,
  COUNT(*) FILTER (WHERE is_read = FALSE) as unread_count
FROM notifications 
WHERE user_id = @user_id;

-- name: CreateRole :one
INSERT INTO roles (
  id,
  company_id,
  role_name,
  description,
  is_system_role,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @role_name,
  @description,
  COALESCE(@is_system_role, FALSE),
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetRolesByCompany :many
SELECT * FROM roles 
WHERE company_id = @company_id OR is_system_role = TRUE
ORDER BY is_system_role DESC, role_name;

-- name: GetRoleByID :one
SELECT * FROM roles WHERE id = @id;

-- name: UpdateRole :one
UPDATE roles SET
  role_name = COALESCE(@role_name, role_name),
  description = COALESCE(@description, description),
  updated_at = NOW()
WHERE id = @id AND is_system_role = FALSE
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles 
WHERE id = @id AND is_system_role = FALSE;

-- name: CreatePermission :one
INSERT INTO permissions (
  id,
  permission_key,
  description,
  category,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @permission_key,
  @description,
  @category,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetAllPermissions :many
SELECT * FROM permissions 
ORDER BY category, permission_key;

-- name: GetPermissionsByCategory :many
SELECT * FROM permissions 
WHERE category = @category
ORDER BY permission_key;

-- name: CreateRolePermission :one
INSERT INTO role_permissions (
  id,
  role_id,
  permission_id,
  created_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @role_id,
  @permission_id,
  COALESCE(@created_at, NOW())
) RETURNING *;

-- name: GetRolePermissions :many
SELECT rp.*, p.permission_key, p.description, p.category
FROM role_permissions rp
JOIN permissions p ON rp.permission_id = p.id
WHERE rp.role_id = @role_id
ORDER BY p.category, p.permission_key;

-- name: DeleteRolePermission :exec
DELETE FROM role_permissions 
WHERE role_id = @role_id AND permission_id = @permission_id;

-- name: CreateUserRole :one
INSERT INTO user_roles (
  id,
  user_id,
  role_id,
  company_id,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @role_id,
  @company_id,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetUserRoles :many
SELECT ur.*, r.role_name, r.description, c.company_name
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
LEFT JOIN companies c ON ur.company_id = c.id
WHERE ur.user_id = @user_id
ORDER BY c.company_name, r.role_name;

-- name: GetUserRolesByCompany :many
SELECT ur.*, r.role_name, r.description
FROM user_roles ur
JOIN roles r ON ur.role_id = r.id
WHERE ur.user_id = @user_id AND ur.company_id = @company_id
ORDER BY r.role_name;

-- name: GetUsersWithRole :many
SELECT ur.*, u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name
FROM user_roles ur
JOIN users u ON ur.user_id = u.id
LEFT JOIN company_staff_profiles csp ON u.id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE ur.role_id = @role_id
  AND (@company_id::uuid IS NULL OR ur.company_id = @company_id)
ORDER BY u.email;

-- name: DeleteUserRole :exec
DELETE FROM user_roles 
WHERE user_id = @user_id AND role_id = @role_id 
  AND (@company_id::uuid IS NULL OR company_id = @company_id);

-- name: GetUserPermissions :many
SELECT DISTINCT p.permission_key, p.description, p.category
FROM user_roles ur
JOIN role_permissions rp ON ur.role_id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE ur.user_id = @user_id
  AND (@company_id::uuid IS NULL OR ur.company_id = @company_id)
ORDER BY p.category, p.permission_key;