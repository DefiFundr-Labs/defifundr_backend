-- name: CreateInvoice :one
INSERT INTO invoices (
  id,
  invoice_number,
  issuer_id,
  recipient_id,
  title,
  description,
  issue_date,
  due_date,
  total_amount,
  currency,
  status,
  payment_method,
  recipient_wallet_address,
  recipient_bank_account_id,
  transaction_hash,
  payment_date,
  rejection_reason,
  ipfs_hash,
  smart_contract_address,
  chain_id,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @invoice_number,
  @issuer_id,
  @recipient_id,
  @title,
  @description,
  @issue_date,
  @due_date,
  @total_amount,
  @currency,
  COALESCE(@status, 'draft'),
  @payment_method,
  @recipient_wallet_address,
  @recipient_bank_account_id,
  @transaction_hash,
  @payment_date,
  @rejection_reason,
  @ipfs_hash,
  @smart_contract_address,
  @chain_id,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetInvoiceByID :one
SELECT i.*, 
       issuer.email as issuer_email,
       COALESCE(issuer_profile.first_name, issuer_personal.first_name) as issuer_first_name,
       COALESCE(issuer_profile.last_name, issuer_personal.last_name) as issuer_last_name,
       recipient.company_name as recipient_company_name
FROM invoices i
JOIN users issuer ON i.issuer_id = issuer.id
JOIN companies recipient ON i.recipient_id = recipient.id
LEFT JOIN company_staff_profiles issuer_profile ON issuer.id = issuer_profile.id
LEFT JOIN personal_users issuer_personal ON issuer.id = issuer_personal.id
WHERE i.id = @id;

-- name: GetInvoiceByNumber :one
SELECT i.*, 
       issuer.email as issuer_email,
       COALESCE(issuer_profile.first_name, issuer_personal.first_name) as issuer_first_name,
       COALESCE(issuer_profile.last_name, issuer_personal.last_name) as issuer_last_name,
       recipient.company_name as recipient_company_name
FROM invoices i
JOIN users issuer ON i.issuer_id = issuer.id
JOIN companies recipient ON i.recipient_id = recipient.id
LEFT JOIN company_staff_profiles issuer_profile ON issuer.id = issuer_profile.id
LEFT JOIN personal_users issuer_personal ON issuer.id = issuer_personal.id
WHERE i.invoice_number = @invoice_number;

-- name: GetUserInvoices :many
SELECT i.*, 
       recipient.company_name as recipient_company_name
FROM invoices i
JOIN companies recipient ON i.recipient_id = recipient.id
WHERE i.issuer_id = @user_id
ORDER BY i.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetCompanyReceivedInvoices :many
SELECT i.*, 
       issuer.email as issuer_email,
       COALESCE(issuer_profile.first_name, issuer_personal.first_name) as issuer_first_name,
       COALESCE(issuer_profile.last_name, issuer_personal.last_name) as issuer_last_name
FROM invoices i
JOIN users issuer ON i.issuer_id = issuer.id
LEFT JOIN company_staff_profiles issuer_profile ON issuer.id = issuer_profile.id
LEFT JOIN personal_users issuer_personal ON issuer.id = issuer_personal.id
WHERE i.recipient_id = @company_id
ORDER BY i.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetCompanySentInvoices :many
SELECT i.*, 
       recipient.company_name as recipient_company_name
FROM invoices i
JOIN companies recipient ON i.recipient_id = recipient.id
JOIN company_users cu ON cu.company_id = recipient.id
WHERE cu.user_id = @user_id AND i.issuer_id = @user_id
ORDER BY i.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetInvoicesForApproval :many
SELECT i.*, 
       issuer.email as issuer_email,
       COALESCE(issuer_profile.first_name, issuer_personal.first_name) as issuer_first_name,
       COALESCE(issuer_profile.last_name, issuer_personal.last_name) as issuer_last_name
FROM invoices i
JOIN users issuer ON i.issuer_id = issuer.id
JOIN company_users cu ON cu.company_id = i.recipient_id
LEFT JOIN company_staff_profiles issuer_profile ON issuer.id = issuer_profile.id
LEFT JOIN personal_users issuer_personal ON issuer.id = issuer_personal.id
WHERE cu.user_id = @user_id 
  AND i.status = 'pending'
  AND (cu.can_manage_invoices = TRUE OR cu.is_administrator = TRUE)
ORDER BY i.created_at;

-- name: UpdateInvoice :one
UPDATE invoices SET
  title = COALESCE(@title, title),
  description = COALESCE(@description, description),
  issue_date = COALESCE(@issue_date, issue_date),
  due_date = COALESCE(@due_date, due_date),
  total_amount = COALESCE(@total_amount, total_amount),
  currency = COALESCE(@currency, currency),
  status = COALESCE(@status, status),
  payment_method = COALESCE(@payment_method, payment_method),
  recipient_wallet_address = COALESCE(@recipient_wallet_address, recipient_wallet_address),
  recipient_bank_account_id = COALESCE(@recipient_bank_account_id, recipient_bank_account_id),
  transaction_hash = COALESCE(@transaction_hash, transaction_hash),
  payment_date = COALESCE(@payment_date, payment_date),
  rejection_reason = COALESCE(@rejection_reason, rejection_reason),
  ipfs_hash = COALESCE(@ipfs_hash, ipfs_hash),
  smart_contract_address = COALESCE(@smart_contract_address, smart_contract_address),
  chain_id = COALESCE(@chain_id, chain_id),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: CreateInvoiceItem :one
INSERT INTO invoice_items (
  id,
  invoice_id,
  description,
  quantity,
  unit_price,
  amount,
  tax_rate,
  tax_amount,
  discount_percentage,
  discount_amount,
  total_amount,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @invoice_id,
  @description,
  COALESCE(@quantity, 1),
  @unit_price,
  @amount,
  COALESCE(@tax_rate, 0),
  COALESCE(@tax_amount, 0),
  COALESCE(@discount_percentage, 0),
  COALESCE(@discount_amount, 0),
  @total_amount,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetInvoiceItems :many
SELECT * FROM invoice_items 
WHERE invoice_id = @invoice_id
ORDER BY created_at;

-- name: UpdateInvoiceItem :one
UPDATE invoice_items SET
  description = COALESCE(@description, description),
  quantity = COALESCE(@quantity, quantity),
  unit_price = COALESCE(@unit_price, unit_price),
  amount = COALESCE(@amount, amount),
  tax_rate = COALESCE(@tax_rate, tax_rate),
  tax_amount = COALESCE(@tax_amount, tax_amount),
  discount_percentage = COALESCE(@discount_percentage, discount_percentage),
  discount_amount = COALESCE(@discount_amount, discount_amount),
  total_amount = COALESCE(@total_amount, total_amount),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: DeleteInvoiceItem :exec
DELETE FROM invoice_items WHERE id = @id;

-- name: SearchInvoices :many
SELECT i.*, 
       issuer.email as issuer_email,
       COALESCE(issuer_profile.first_name, issuer_personal.first_name) as issuer_first_name,
       COALESCE(issuer_profile.last_name, issuer_personal.last_name) as issuer_last_name,
       recipient.company_name as recipient_company_name
FROM invoices i
JOIN users issuer ON i.issuer_id = issuer.id
JOIN companies recipient ON i.recipient_id = recipient.id
LEFT JOIN company_staff_profiles issuer_profile ON issuer.id = issuer_profile.id
LEFT JOIN personal_users issuer_personal ON issuer.id = issuer_personal.id
WHERE (i.issuer_id = @user_id OR i.recipient_id IN (
    SELECT company_id FROM company_users WHERE user_id = @user_id
))
AND (
    i.invoice_number ILIKE '%' || @search_term || '%' OR
    i.title ILIKE '%' || @search_term || '%' OR
    recipient.company_name ILIKE '%' || @search_term || '%'
)
ORDER BY i.created_at DESC
LIMIT @limit_val OFFSET @offset_val;

-- name: GetOverdueInvoices :many
SELECT i.*, 
       issuer.email as issuer_email,
       COALESCE(issuer_profile.first_name, issuer_personal.first_name) as issuer_first_name,
       COALESCE(issuer_profile.last_name, issuer_personal.last_name) as issuer_last_name,
       recipient.company_name as recipient_company_name
FROM invoices i
JOIN users issuer ON i.issuer_id = issuer.id
JOIN companies recipient ON i.recipient_id = recipient.id
LEFT JOIN company_staff_profiles issuer_profile ON issuer.id = issuer_profile.id
LEFT JOIN personal_users issuer_personal ON issuer.id = issuer_personal.id
WHERE i.due_date < CURRENT_DATE 
  AND i.status NOT IN ('paid', 'cancelled')
  AND (@company_id::uuid IS NULL OR i.recipient_id = @company_id)
  AND (@user_id::uuid IS NULL OR i.issuer_id = @user_id)
ORDER BY i.due_date;

-- name: GetInvoicesByStatus :many
SELECT i.*, 
       issuer.email as issuer_email,
       COALESCE(issuer_profile.first_name, issuer_personal.first_name) as issuer_first_name,
       COALESCE(issuer_profile.last_name, issuer_personal.last_name) as issuer_last_name,
       recipient.company_name as recipient_company_name
FROM invoices i
JOIN users issuer ON i.issuer_id = issuer.id
JOIN companies recipient ON i.recipient_id = recipient.id
LEFT JOIN company_staff_profiles issuer_profile ON issuer.id = issuer_profile.id
LEFT JOIN personal_users issuer_personal ON issuer.id = issuer_personal.id
WHERE i.status = @status
  AND (@company_id::uuid IS NULL OR i.recipient_id = @company_id)
  AND (@user_id::uuid IS NULL OR i.issuer_id = @user_id)
ORDER BY i.created_at DESC
LIMIT @limit_val OFFSET @offset_val;