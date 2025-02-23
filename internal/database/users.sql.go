// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const changePassword = `-- name: ChangePassword :exec
UPDATE users
SET password = $2, updated_at = NOW()
WHERE id = $1
`

type ChangePasswordParams struct {
	ID       uuid.UUID
	Password string
}

func (q *Queries) ChangePassword(ctx context.Context, arg ChangePasswordParams) error {
	_, err := q.db.ExecContext(ctx, changePassword, arg.ID, arg.Password)
	return err
}

const changeUsername = `-- name: ChangeUsername :one
UPDATE users
SET username = $2, updated_at = NOW()
WHERE id = $1
RETURNING id, created_at, updated_at, email, password, username, is_premium, verification_code, is_verified
`

type ChangeUsernameParams struct {
	ID       uuid.UUID
	Username string
}

func (q *Queries) ChangeUsername(ctx context.Context, arg ChangeUsernameParams) (User, error) {
	row := q.db.QueryRowContext(ctx, changeUsername, arg.ID, arg.Username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Username,
		&i.IsPremium,
		&i.VerificationCode,
		&i.IsVerified,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, created_at, updated_at, email, password, username, is_premium, verification_code, is_verified FROM users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Email,
			&i.Password,
			&i.Username,
			&i.IsPremium,
			&i.VerificationCode,
			&i.IsVerified,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmailOrUsername = `-- name: GetUserByEmailOrUsername :one
SELECT id, created_at, updated_at, email, password, username, is_premium, verification_code, is_verified FROM users
WHERE email = $1 OR username = $2
`

type GetUserByEmailOrUsernameParams struct {
	Email    string
	Username string
}

func (q *Queries) GetUserByEmailOrUsername(ctx context.Context, arg GetUserByEmailOrUsernameParams) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmailOrUsername, arg.Email, arg.Username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Username,
		&i.IsPremium,
		&i.VerificationCode,
		&i.IsVerified,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, updated_at, email, password, username, is_premium, verification_code, is_verified FROM users
WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Username,
		&i.IsPremium,
		&i.VerificationCode,
		&i.IsVerified,
	)
	return i, err
}

const resetPassword = `-- name: ResetPassword :exec
UPDATE users
SET password = $2, updated_at = NOW()
WHERE id = $1
`

type ResetPasswordParams struct {
	ID       uuid.UUID
	Password string
}

func (q *Queries) ResetPassword(ctx context.Context, arg ResetPasswordParams) error {
	_, err := q.db.ExecContext(ctx, resetPassword, arg.ID, arg.Password)
	return err
}

const sendResetVerificationCode = `-- name: SendResetVerificationCode :exec
UPDATE users
SET verification_code = $2
WHERE id = $1
`

type SendResetVerificationCodeParams struct {
	ID               uuid.UUID
	VerificationCode int32
}

func (q *Queries) SendResetVerificationCode(ctx context.Context, arg SendResetVerificationCodeParams) error {
	_, err := q.db.ExecContext(ctx, sendResetVerificationCode, arg.ID, arg.VerificationCode)
	return err
}

const verifyVerificationCode = `-- name: VerifyVerificationCode :exec
UPDATE users 
SET verification_code = 0
WHERE id = $1
`

func (q *Queries) VerifyVerificationCode(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, verifyVerificationCode, id)
	return err
}
