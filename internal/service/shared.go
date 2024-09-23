package service

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

func floatToString(f float32) string {
	return fmt.Sprintf("%.2f", f)
}

func intToString(i int) string {
	return fmt.Sprintf("%d", i)
}

func sqlNullStringToString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

func nullUuidToUuid(nuuid uuid.NullUUID) *uuid.UUID {
	if nuuid.Valid {
		return &nuuid.UUID
	}
	return nil
}
