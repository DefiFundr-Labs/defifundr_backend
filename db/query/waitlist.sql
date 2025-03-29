-- name: CreateWaitlistEntry :one
-- Creates a new waitlist entry and returns the created entry
INSERT INTO waitlist (
  id,
  email,
  full_name,
  referral_code,
  referral_source,
  status,
  signup_date,
  invited_date,
  registered_date,
  metadata,
  created_at,
  updated_at
) VALUES (
  COALESCE($1, uuid_generate_v4()),
  $2,
  $3,
  $4,
  $5,
  COALESCE($6, 'waiting'),
  COALESCE($7, now()),
  $8,
  $9,
  $10,
  COALESCE($11, now()),
  COALESCE($12, now())
) RETURNING *;

-- name: GetWaitlistEntryByEmail :one
-- Retrieves a single waitlist entry by email address
SELECT * FROM waitlist
WHERE email = $1
LIMIT 1;

-- name: GetWaitlistEntryByID :one
-- Retrieves a single waitlist entry by ID
SELECT * FROM waitlist 
WHERE id = $1 OR id::text = $1 
LIMIT 1;

-- name: GetWaitlistEntryByReferralCode :one
-- Retrieves a single waitlist entry by referral code
SELECT * FROM waitlist
WHERE referral_code = $1
LIMIT 1;

-- name: ListWaitlistEntries :many
-- Lists waitlist entries with pagination and filtering support
SELECT *
FROM waitlist
WHERE 
  ($1::text IS NULL OR status = $1) AND
  ($2::text IS NULL OR referral_source = $2)
ORDER BY 
  CASE WHEN $5::text = 'signup_date_asc' THEN signup_date END ASC,
  CASE WHEN $5::text = 'signup_date_desc' OR $5::text IS NULL THEN signup_date END DESC
LIMIT $3
OFFSET $4;

-- name: CountWaitlistEntries :one
-- Counts the total number of waitlist entries matching filters
SELECT COUNT(*) FROM waitlist
WHERE 
  ($1::text IS NULL OR status = $1) AND
  ($2::text IS NULL OR referral_source = $2);

-- name: UpdateWaitlistEntryStatus :exec
-- Updates the status of a waitlist entry
UPDATE waitlist
SET
  status = $2,
  invited_date = CASE WHEN $2 = 'invited' AND invited_date IS NULL THEN now() ELSE invited_date END,
  registered_date = CASE WHEN $2 = 'registered' AND registered_date IS NULL THEN now() ELSE registered_date END,
  updated_at = now()
WHERE id = $1;

-- name: ExportWaitlistEntries :many
-- Retrieves all waitlist entries for export
SELECT *
FROM waitlist
ORDER BY signup_date DESC;

-- name: DeleteWaitlistEntry :exec
-- Permanently deletes a waitlist entry
DELETE FROM waitlist
WHERE id = $1;

-- name: GetWaitlistPosition :one
-- Gets the position of an entry in the waitlist (by signup date)
SELECT COUNT(*) 
FROM waitlist AS w1
WHERE w1.status = 'waiting' AND w1.signup_date <= (
  SELECT w2.signup_date FROM waitlist AS w2 WHERE w2.id = $1
);

-- name: SearchWaitlist :many
-- Searches for waitlist entries by email or name with pagination
SELECT *
FROM waitlist
WHERE 
  (
    email ILIKE '%' || $1 || '%' OR
    full_name ILIKE '%' || $1 || '%'
  )
ORDER BY signup_date DESC
LIMIT $2
OFFSET $3;

-- name: CountSearchWaitlist :one
-- Counts the number of waitlist entries matching a search query
SELECT COUNT(*)
FROM waitlist
WHERE 
  (
    email ILIKE '%' || $1 || '%' OR
    full_name ILIKE '%' || $1 || '%'
  );

-- name: GetWaitlistStatsBySource :many
-- Gets waitlist statistics grouped by referral source
SELECT 
  referral_source, 
  COUNT(*) as count
FROM waitlist
GROUP BY referral_source
ORDER BY count DESC;

-- name: GetWaitlistStatsByStatus :many
-- Gets waitlist statistics grouped by status
SELECT 
  status, 
  COUNT(*) as count
FROM waitlist
GROUP BY status
ORDER BY count DESC;