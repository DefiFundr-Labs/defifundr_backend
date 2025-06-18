-- name: CreatePayrollPeriod :one
INSERT INTO payroll_periods (
  id,
  company_id,
  period_name,
  frequency,
  start_date,
  end_date,
  payment_date,
  status,
  is_recurring,
  next_period_id,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @period_name,
  @frequency,
  @start_date,
  @end_date,
  @payment_date,
  COALESCE(@status, 'draft'),
  COALESCE(@is_recurring, FALSE),
  @next_period_id,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetPayrollPeriodByID :one
SELECT * FROM payroll_periods WHERE id = @id;

-- name: GetPayrollPeriodsByCompany :many
SELECT * FROM payroll_periods 
WHERE company_id = @company_id
ORDER BY start_date DESC;

-- name: GetActivePayrollPeriods :many
SELECT * FROM payroll_periods 
WHERE company_id = @company_id AND status = 'active'
ORDER BY start_date DESC;

-- name: UpdatePayrollPeriod :one
UPDATE payroll_periods SET
  period_name = COALESCE(@period_name, period_name),
  frequency = COALESCE(@frequency, frequency),
  start_date = COALESCE(@start_date, start_date),
  end_date = COALESCE(@end_date, end_date),
  payment_date = COALESCE(@payment_date, payment_date),
  status = COALESCE(@status, status),
  is_recurring = COALESCE(@is_recurring, is_recurring),
  next_period_id = COALESCE(@next_period_id, next_period_id),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: CreatePayroll :one
INSERT INTO payrolls (
  id,
  company_id,
  period_id,
  name,
  description,
  total_amount,
  base_currency,
  status,
  execution_type,
  scheduled_execution_time,
  executed_at,
  smart_contract_address,
  chain_id,
  transaction_hash,
  created_by,
  approved_by,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @period_id,
  @name,
  @description,
  COALESCE(@total_amount, 0),
  @base_currency,
  COALESCE(@status, 'draft'),
  COALESCE(@execution_type, 'manual'),
  @scheduled_execution_time,
  @executed_at,
  @smart_contract_address,
  @chain_id,
  @transaction_hash,
  @created_by,
  @approved_by,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetPayrollByID :one
SELECT p.*, pp.period_name, pp.start_date, pp.end_date
FROM payrolls p
JOIN payroll_periods pp ON p.period_id = pp.id
WHERE p.id = @id;

-- name: GetPayrollsByCompany :many
SELECT p.*, pp.period_name, pp.start_date, pp.end_date
FROM payrolls p
JOIN payroll_periods pp ON p.period_id = pp.id
WHERE p.company_id = @company_id
ORDER BY p.created_at DESC;

-- name: GetPayrollsByStatus :many
SELECT p.*, pp.period_name, pp.start_date, pp.end_date
FROM payrolls p
JOIN payroll_periods pp ON p.period_id = pp.id
WHERE p.company_id = @company_id AND p.status = @status
ORDER BY p.created_at DESC;

-- name: UpdatePayroll :one
UPDATE payrolls SET
  name = COALESCE(@name, name),
  description = COALESCE(@description, description),
  total_amount = COALESCE(@total_amount, total_amount),
  base_currency = COALESCE(@base_currency, base_currency),
  status = COALESCE(@status, status),
  execution_type = COALESCE(@execution_type, execution_type),
  scheduled_execution_time = COALESCE(@scheduled_execution_time, scheduled_execution_time),
  executed_at = COALESCE(@executed_at, executed_at),
  smart_contract_address = COALESCE(@smart_contract_address, smart_contract_address),
  chain_id = COALESCE(@chain_id, chain_id),
  transaction_hash = COALESCE(@transaction_hash, transaction_hash),
  approved_by = COALESCE(@approved_by, approved_by),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: CreatePayrollItem :one
INSERT INTO payroll_items (
  id,
  payroll_id,
  employee_id,
  base_amount,
  base_currency,
  payment_amount,
  payment_currency,
  exchange_rate,
  payment_method,
  payment_split,
  status,
  transaction_hash,
  recipient_wallet_address,
  recipient_bank_account_id,
  notes,
  timesheet_id,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @payroll_id,
  @employee_id,
  @base_amount,
  @base_currency,
  @payment_amount,
  @payment_currency,
  COALESCE(@exchange_rate, 1),
  @payment_method,
  @payment_split,
  COALESCE(@status, 'pending'),
  @transaction_hash,
  @recipient_wallet_address,
  @recipient_bank_account_id,
  @notes,
  @timesheet_id,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetPayrollItemByID :one
SELECT pi.*, 
       ce.employee_id, ce.position,
       u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name
FROM payroll_items pi
JOIN company_employees ce ON pi.employee_id = ce.id
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_staff_profiles csp ON ce.user_id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE pi.id = @id;

-- name: GetPayrollItemsByPayroll :many
SELECT pi.*, 
       ce.employee_id, ce.position,
       u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name
FROM payroll_items pi
JOIN company_employees ce ON pi.employee_id = ce.id
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_staff_profiles csp ON ce.user_id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE pi.payroll_id = @payroll_id
ORDER BY pi.created_at;

-- name: GetPayrollItemsByEmployee :many
SELECT pi.*, p.name as payroll_name, pp.period_name
FROM payroll_items pi
JOIN payrolls p ON pi.payroll_id = p.id
JOIN payroll_periods pp ON p.period_id = pp.id
WHERE pi.employee_id = @employee_id
ORDER BY pi.created_at DESC;

-- name: GetUserPayrollHistory :many
SELECT pi.*, p.name as payroll_name, pp.period_name, c.company_name
FROM payroll_items pi
JOIN payrolls p ON pi.payroll_id = p.id
JOIN payroll_periods pp ON p.period_id = pp.id
JOIN companies c ON p.company_id = c.id
JOIN company_employees ce ON pi.employee_id = ce.id
WHERE ce.user_id = @user_id
ORDER BY pi.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: UpdatePayrollItem :one
UPDATE payroll_items SET
  base_amount = COALESCE(@base_amount, base_amount),
  base_currency = COALESCE(@base_currency, base_currency),
  payment_amount = COALESCE(@payment_amount, payment_amount),
  payment_currency = COALESCE(@payment_currency, payment_currency),
  exchange_rate = COALESCE(@exchange_rate, exchange_rate),
  payment_method = COALESCE(@payment_method, payment_method),
  payment_split = COALESCE(@payment_split, payment_split),
  status = COALESCE(@status, status),
  transaction_hash = COALESCE(@transaction_hash, transaction_hash),
  recipient_wallet_address = COALESCE(@recipient_wallet_address, recipient_wallet_address),
  recipient_bank_account_id = COALESCE(@recipient_bank_account_id, recipient_bank_account_id),
  notes = COALESCE(@notes, notes),
  timesheet_id = COALESCE(@timesheet_id, timesheet_id),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: GetPayrollApprovalQueue :many
SELECT pi.*, 
       ce.employee_id, ce.position,
       u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name,
       p.name as payroll_name
FROM payroll_items pi
JOIN company_employees ce ON pi.employee_id = ce.id
JOIN payrolls p ON pi.payroll_id = p.id
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_staff_profiles csp ON ce.user_id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE p.company_id = @company_id AND pi.status = 'pending'
ORDER BY pi.created_at;

-- name: GetUpcomingPayments :many
SELECT pi.*, 
       ce.employee_id, ce.position,
       u.email,
       COALESCE(csp.first_name, pu.first_name) as first_name,
       COALESCE(csp.last_name, pu.last_name) as last_name,
       p.name as payroll_name,
       pp.payment_date
FROM payroll_items pi
JOIN company_employees ce ON pi.employee_id = ce.id
JOIN payrolls p ON pi.payroll_id = p.id
JOIN payroll_periods pp ON p.period_id = pp.id
LEFT JOIN users u ON ce.user_id = u.id
LEFT JOIN company_staff_profiles csp ON ce.user_id = csp.id
LEFT JOIN personal_users pu ON u.id = pu.id
WHERE ce.user_id = @user_id 
  AND pi.status IN ('approved', 'processing')
  AND pp.payment_date >= CURRENT_DATE
ORDER BY pp.payment_date;