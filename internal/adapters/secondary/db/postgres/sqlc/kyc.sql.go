// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0
// source: kyc.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
)

const createKYC = `-- name: CreateKYC :one
INSERT INTO kyc (
  id,
  user_id,
  face_verification,
  identity_verification
) VALUES (
  uuid_generate_v4(), $1, $2, $3
) RETURNING id, user_id, face_verification, identity_verification, updated_at, created_at
`

type CreateKYCParams struct {
	UserID               uuid.UUID `json:"user_id"`
	FaceVerification     bool      `json:"face_verification"`
	IdentityVerification bool      `json:"identity_verification"`
}

func (q *Queries) CreateKYC(ctx context.Context, arg CreateKYCParams) (Kyc, error) {
	row := q.db.QueryRow(ctx, createKYC, arg.UserID, arg.FaceVerification, arg.IdentityVerification)
	var i Kyc
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FaceVerification,
		&i.IdentityVerification,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const getKYC = `-- name: GetKYC :one
SELECT id, user_id, face_verification, identity_verification, updated_at, created_at FROM kyc
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetKYC(ctx context.Context, userID uuid.UUID) (Kyc, error) {
	row := q.db.QueryRow(ctx, getKYC, userID)
	var i Kyc
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FaceVerification,
		&i.IdentityVerification,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateKYC = `-- name: UpdateKYC :one
UPDATE kyc
SET 
  face_verification = $2,
  identity_verification = $3,
  updated_at = now()
WHERE user_id = $1
RETURNING id, user_id, face_verification, identity_verification, updated_at, created_at
`

type UpdateKYCParams struct {
	UserID               uuid.UUID `json:"user_id"`
	FaceVerification     bool      `json:"face_verification"`
	IdentityVerification bool      `json:"identity_verification"`
}

func (q *Queries) UpdateKYC(ctx context.Context, arg UpdateKYCParams) (Kyc, error) {
	row := q.db.QueryRow(ctx, updateKYC, arg.UserID, arg.FaceVerification, arg.IdentityVerification)
	var i Kyc
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.FaceVerification,
		&i.IdentityVerification,
		&i.UpdatedAt,
		&i.CreatedAt,
	)
	return i, err
}
