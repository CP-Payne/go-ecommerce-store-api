// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID
	Name           sql.NullString
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
