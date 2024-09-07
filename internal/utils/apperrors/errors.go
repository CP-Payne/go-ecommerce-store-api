package apperrors

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrConflict  = errors.New("conflict")
	ErrInternal  = errors.New("internal error")
	ErrNotFound  = errors.New("resource not found")
	ErrAuthCode  = errors.New("auth code")
	ErrParseUUID = errors.New("could not parse UUID")
)

func IsPqError(err error, code pq.ErrorCode) bool {
	var pqErr *pq.Error
	return errors.As(err, &pqErr) && pqErr.Code == code
}

// IsUniqueViolation checks for unique violation errors
func IsUniqueViolation(err error) bool {
	return IsPqError(err, "23505")
}

// IsNoRowsError checks if the error is a "no rows in result set" error
func IsNoRowsError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
