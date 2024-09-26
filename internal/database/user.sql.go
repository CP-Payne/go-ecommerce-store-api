// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, email,name, hashed_password, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, email, hashed_password, created_at, updated_at
`

type CreateUserParams struct {
	ID             uuid.UUID
	Email          string
	Name           sql.NullString
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Name,
		arg.HashedPassword,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, name, email, hashed_password, created_at, updated_at FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserDetails = `-- name: GetUserDetails :one
SELECT id, email, name 
FROM users
WHERE id = $1
`

type GetUserDetailsRow struct {
	ID    uuid.UUID
	Email string
	Name  sql.NullString
}

func (q *Queries) GetUserDetails(ctx context.Context, id uuid.UUID) (GetUserDetailsRow, error) {
	row := q.db.QueryRowContext(ctx, getUserDetails, id)
	var i GetUserDetailsRow
	err := row.Scan(&i.ID, &i.Email, &i.Name)
	return i, err
}
