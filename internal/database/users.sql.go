// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password, username, is_premium)
VALUES (
   $1,
   NOW(),
   NOW(),
   $2,
   $3,
   $4,
   $5
)
RETURNING id, created_at, updated_at, email, password, username, is_premium
`

type CreateUserParams struct {
	ID        uuid.UUID
	Email     string
	Password  string
	Username  string
	IsPremium bool
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Password,
		arg.Username,
		arg.IsPremium,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Username,
		&i.IsPremium,
	)
	return i, err
}

const getUserByIdentifier = `-- name: GetUserByIdentifier :one
SELECT id, created_at, updated_at, email, password, username, is_premium FROM users
WHERE email = $1 OR username = $2
`

type GetUserByIdentifierParams struct {
	Email    string
	Username string
}

func (q *Queries) GetUserByIdentifier(ctx context.Context, arg GetUserByIdentifierParams) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByIdentifier, arg.Email, arg.Username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Username,
		&i.IsPremium,
	)
	return i, err
}
