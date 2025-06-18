-- name: CreateSupportedCountry :one
INSERT INTO supported_countries (
  id,
  country_code,
  country_name,
  region,
  currency_code,
  currency_symbol,
  is_active,
  is_high_risk,
  requires_enhanced_kyc,
  requires_enhanced_kyb,
  timezone,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @country_code,
  @country_name,
  @region,
  @currency_code,
  @currency_symbol,
  COALESCE(@is_active, TRUE),
  COALESCE(@is_high_risk, FALSE),
  COALESCE(@requires_enhanced_kyc, FALSE),
  COALESCE(@requires_enhanced_kyb, FALSE),
  @timezone,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetSupportedCountries :many
SELECT * FROM supported_countries 
WHERE is_active = TRUE
ORDER BY country_name;

-- name: GetSupportedCountryByCode :one
SELECT * FROM supported_countries 
WHERE country_code = @country_code AND is_active = TRUE;

-- name: CreateKYCDocument :one
INSERT INTO kyc_documents (
  id,
  user_id,
  country_id,
  document_type,
  document_number,
  document_country,
  issue_date,
  expiry_date,
  document_url,
  ipfs_hash,
  verification_status,
  verification_level,
  verification_notes,
  verified_by,
  verified_at,
  rejection_reason,
  metadata,
  meets_requirements,
  requirement_id,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @country_id,
  @document_type,
  @document_number,
  @document_country,
  @issue_date,
  @expiry_date,
  @document_url,
  @ipfs_hash,
  COALESCE(@verification_status, 'pending'),
  @verification_level,
  @verification_notes,
  @verified_by,
  @verified_at,
  @rejection_reason,
  @metadata,
  COALESCE(@meets_requirements, FALSE),
  @requirement_id,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetKYCDocumentsByUser :many
SELECT kd.*, sc.country_name, kcr.document_type as required_document_type
FROM kyc_documents kd
JOIN supported_countries sc ON kd.country_id = sc.id
LEFT JOIN kyc_country_requirements kcr ON kd.requirement_id = kcr.id
WHERE kd.user_id = @user_id
ORDER BY kd.created_at DESC;

-- name: GetKYCDocumentsByCountry :many
SELECT kd.*, sc.country_name, kcr.document_type as required_document_type
FROM kyc_documents kd
JOIN supported_countries sc ON kd.country_id = sc.id
LEFT JOIN kyc_country_requirements kcr ON kd.requirement_id = kcr.id
WHERE kd.user_id = @user_id AND kd.country_id = @country_id
ORDER BY kd.created_at DESC;

-- name: UpdateKYCDocument :one
UPDATE kyc_documents SET
  document_number = COALESCE(@document_number, document_number),
  document_country = COALESCE(@document_country, document_country),
  issue_date = COALESCE(@issue_date, issue_date),
  expiry_date = COALESCE(@expiry_date, expiry_date),
  document_url = COALESCE(@document_url, document_url),
  ipfs_hash = COALESCE(@ipfs_hash, ipfs_hash),
  verification_status = COALESCE(@verification_status, verification_status),
  verification_level = COALESCE(@verification_level, verification_level),
  verification_notes = COALESCE(@verification_notes, verification_notes),
  verified_by = COALESCE(@verified_by, verified_by),
  verified_at = COALESCE(@verified_at, verified_at),
  rejection_reason = COALESCE(@rejection_reason, rejection_reason),
  metadata = COALESCE(@metadata, metadata),
  meets_requirements = COALESCE(@meets_requirements, meets_requirements),
  updated_at = NOW()
WHERE id = @id
RETURNING *;

-- name: CreateKYBDocument :one
INSERT INTO kyb_documents (
  id,
  company_id,
  country_id,
  document_type,
  document_number,
  document_country,
  issue_date,
  expiry_date,
  document_url,
  ipfs_hash,
  verification_status,
  verification_level,
  verification_notes,
  verified_by,
  verified_at,
  rejection_reason,
  metadata,
  meets_requirements,
  requirement_id,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @country_id,
  @document_type,
  @document_number,
  @document_country,
  @issue_date,
  @expiry_date,
  @document_url,
  @ipfs_hash,
  COALESCE(@verification_status, 'pending'),
  @verification_level,
  @verification_notes,
  @verified_by,
  @verified_at,
  @rejection_reason,
  @metadata,
  COALESCE(@meets_requirements, FALSE),
  @requirement_id,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetKYBDocumentsByCompany :many
SELECT kbd.*, sc.country_name, kbcr.document_type as required_document_type
FROM kyb_documents kbd
JOIN supported_countries sc ON kbd.country_id = sc.id
LEFT JOIN kyb_country_requirements kbcr ON kbd.requirement_id = kbcr.id
WHERE kbd.company_id = @company_id
ORDER BY kbd.created_at DESC;

-- name: GetKYCCountryRequirements :many
SELECT * FROM kyc_country_requirements 
WHERE country_id = @country_id AND is_active = TRUE
ORDER BY document_type;

-- name: GetKYBCountryRequirements :many
SELECT * FROM kyb_country_requirements 
WHERE country_id = @country_id AND is_active = TRUE
ORDER BY document_type;

-- name: CreateUserCountryKYCStatus :one
INSERT INTO user_country_kyc_status (
  id,
  user_id,
  country_id,
  verification_status,
  verification_level,
  verification_date,
  expiry_date,
  rejection_reason,
  notes,
  risk_rating,
  restricted_features,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @user_id,
  @country_id,
  COALESCE(@verification_status, 'pending'),
  @verification_level,
  @verification_date,
  @expiry_date,
  @rejection_reason,
  @notes,
  @risk_rating,
  @restricted_features,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetUserKYCStatus :one
SELECT ukcs.*, sc.country_name
FROM user_country_kyc_status ukcs
JOIN supported_countries sc ON ukcs.country_id = sc.id
WHERE ukcs.user_id = @user_id AND ukcs.country_id = @country_id;

-- name: GetUserKYCStatusByCountryCode :one
SELECT ukcs.*, sc.country_name
FROM user_country_kyc_status ukcs
JOIN supported_countries sc ON ukcs.country_id = sc.id
WHERE ukcs.user_id = @user_id AND sc.country_code = @country_code;

-- name: UpdateUserKYCStatus :one
UPDATE user_country_kyc_status SET
  verification_status = COALESCE(@verification_status, verification_status),
  verification_level = COALESCE(@verification_level, verification_level),
  verification_date = COALESCE(@verification_date, verification_date),
  expiry_date = COALESCE(@expiry_date, expiry_date),
  rejection_reason = COALESCE(@rejection_reason, rejection_reason),
  notes = COALESCE(@notes, notes),
  risk_rating = COALESCE(@risk_rating, risk_rating),
  restricted_features = COALESCE(@restricted_features, restricted_features),
  updated_at = NOW()
WHERE user_id = @user_id AND country_id = @country_id
RETURNING *;

-- name: CreateCompanyCountryKYBStatus :one
INSERT INTO company_country_kyb_status (
  id,
  company_id,
  country_id,
  verification_status,
  verification_level,
  verification_date,
  expiry_date,
  rejection_reason,
  notes,
  risk_rating,
  restricted_features,
  created_at,
  updated_at
) VALUES (
  COALESCE(@id, uuid_generate_v4()),
  @company_id,
  @country_id,
  COALESCE(@verification_status, 'pending'),
  @verification_level,
  @verification_date,
  @expiry_date,
  @rejection_reason,
  @notes,
  @risk_rating,
  @restricted_features,
  COALESCE(@created_at, NOW()),
  COALESCE(@updated_at, NOW())
) RETURNING *;

-- name: GetCompanyKYBStatus :one
SELECT ckbs.*, sc.country_name
FROM company_country_kyb_status ckbs
JOIN supported_countries sc ON ckbs.country_id = sc.id
WHERE ckbs.company_id = @company_id AND ckbs.country_id = @country_id;

-- name: UpdateCompanyKYBStatus :one
UPDATE company_country_kyb_status SET
  verification_status = COALESCE(@verification_status, verification_status),
  verification_level = COALESCE(@verification_level, verification_level),
  verification_date = COALESCE(@verification_date, verification_date),
  expiry_date = COALESCE(@expiry_date, expiry_date),
  rejection_reason = COALESCE(@rejection_reason, rejection_reason),
  notes = COALESCE(@notes, notes),
  risk_rating = COALESCE(@risk_rating, risk_rating),
  restricted_features = COALESCE(@restricted_features, restricted_features),
  updated_at = NOW()
WHERE company_id = @company_id AND country_id = @country_id
RETURNING *;