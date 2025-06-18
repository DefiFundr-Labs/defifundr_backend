-- name: CreateTimesheet :one
INSERT INTO timesheets (
  id,
  company_id,
  employee_id,
  period_id,
  status,
  total_hours,
  billable_hours,
  overtime_hours,
  hourly_rate,
  rate_currency,
  total_amount,
  submitted_at,
  approved_at,
  approved_by,
  rejection_reason,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @employee_id,
  @period_id,
  COALESCE(@status, 'draft'),
  COALESCE(@total_hours, 0),
  COALESCE(@billable_hours, 0),
  COALESCE(@overtime_hours, 0),
  @hourly_rate,
  @rate_currency,
  COALESCE(@total_amount, 0),
  @submitted_at,
  @approved_at,
  @approved_by,
  @rejection_reason,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetTimesheetByID :one
SELECT t.*, 
       ce.employee_id, ce.position,
       u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name,
       pp.period_name, pp.start_date, pp.end_date
FROM timesheets t
JOIN company_employees ce ON t.employee_id = ce.id
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_staff_profiles csp ON ce.user_id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
LEFT JOIN payroll_periods pp ON t.period_id = pp.id
WHERE t.id = @id;

-- name: GetTimesheetsByEmployee :many
SELECT t.*, 
       pp.period_name, pp.start_date, pp.end_date
FROM timesheets t
LEFT JOIN payroll_periods pp ON t.period_id = pp.id
WHERE t.employee_id = @employee_id
ORDER BY t.created_at DESC;

-- name: GetUserTimesheets :many
SELECT t.*, 
       ce.employee_id,
       c.company_name,
       pp.period_name, pp.start_date, pp.end_date
FROM timesheets t
JOIN company_employees ce ON t.employee_id = ce.id
JOIN companies c ON t.company_id = c.id
LEFT JOIN payroll_periods pp ON t.period_id = pp.id
WHERE ce.user_id = @user_id
ORDER BY t.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetCompanyTimesheets :many
SELECT t.*, 
       ce.employee_id, ce.position,
       u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name,
       pp.period_name, pp.start_date, pp.end_date
FROM timesheets t
JOIN company_employees ce ON t.employee_id = ce.id
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_staff_profiles csp ON ce.user_id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
LEFT JOIN payroll_periods pp ON t.period_id = pp.id
WHERE t.company_id = @company_id
ORDER BY t.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetTimesheetsByStatus :many
SELECT t.*, 
       ce.employee_id, ce.position,
       u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name,
       pp.period_name, pp.start_date, pp.end_date
FROM timesheets t
JOIN company_employees ce ON t.employee_id = ce.id
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_staff_profiles csp ON ce.user_id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
LEFT JOIN payroll_periods pp ON t.period_id = pp.id
WHERE t.company_id = @company_id AND t.status = @status
ORDER BY t.created_at DESC;

-- name: UpdateTimesheet :one
UPDATE timesheets SET
  period_id = COALESCE(@period_id, period_id),
  status = COALESCE(@status, status),
  total_hours = COALESCE(@total_hours, total_hours),
  billable_hours = COALESCE(@billable_hours, billable_hours),
  overtime_hours = COALESCE(@overtime_hours, overtime_hours),
  hourly_rate = COALESCE(@hourly_rate, hourly_rate),
  rate_currency = COALESCE(@rate_currency, rate_currency),
  total_amount = COALESCE(@total_amount, total_amount),
  submitted_at = COALESCE(@submitted_at, submitted_at),
  approved_at = COALESCE(@approved_at, approved_at),
  approved_by = COALESCE(@approved_by, approved_by),
  rejection_reason = COALESCE(@rejection_reason, rejection_reason),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: CreateTimesheetEntry :one
INSERT INTO timesheet_entries (
  id,
  timesheet_id,
  date,
  start_time,
  end_time,
  hours,
  is_billable,
  is_overtime,
  project,
  task,
  description,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @timesheet_id,
  @date,
  @start_time,
  @end_time,
  @hours,
  COALESCE(@is_billable, TRUE),
  COALESCE(@is_overtime, FALSE),
  @project,
  @task,
  @description,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetTimesheetEntries :many
SELECT * FROM timesheet_entries 
WHERE timesheet_id = @timesheet_id
ORDER BY date, start_time;

-- name: GetTimesheetEntriesByDate :many
SELECT * FROM timesheet_entries 
WHERE timesheet_id = @timesheet_id 
  AND date >= @start_date 
  AND date <= @end_date
ORDER BY date, start_time;

-- name: UpdateTimesheetEntry :one
UPDATE timesheet_entries SET
  date = COALESCE(@date, date),
  start_time = COALESCE(@start_time, start_time),
  end_time = COALESCE(@end_time, end_time),
  hours = COALESCE(@hours, hours),
  is_billable = COALESCE(@is_billable, is_billable),
  is_overtime = COALESCE(@is_overtime, is_overtime),
  project = COALESCE(@project, project),
  task = COALESCE(@task, task),
  description = COALESCE(@description, description),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteTimesheetEntry :exec
DELETE FROM timesheet_entries WHERE id = @id;

-- name: GetTimesheetProjects :many
SELECT DISTINCT project 
FROM timesheet_entries 
WHERE timesheet_id IN (
  SELECT id FROM timesheets WHERE company_id = @company_id
)
AND project IS NOT NULL
ORDER BY project;

-- name: GetEmployeeTimesheetSummary :one
SELECT 
  COUNT(*) as total_timesheets,
  SUM(total_hours) as total_hours,
  SUM(billable_hours) as total_billable_hours,
  SUM(overtime_hours) as total_overtime_hours,
  SUM(total_amount) as total_amount
FROM timesheets 
WHERE employee_id = @employee_id
  AND (@start_date::date IS NULL OR created_at >= @start_date)
  AND (@end_date::date IS NULL OR created_at <= @end_date);

-- name: GetCompanyTimesheetSummary :one
SELECT 
  COUNT(*) as total_timesheets,
  SUM(total_hours) as total_hours,
  SUM(billable_hours) as total_billable_hours,
  SUM(overtime_hours) as total_overtime_hours,
  SUM(total_amount) as total_amount
FROM timesheets 
WHERE company_id = @company_id
  AND (@start_date::date IS NULL OR created_at >= @start_date)
  AND (@end_date::date IS NULL OR created_at <= @end_date);