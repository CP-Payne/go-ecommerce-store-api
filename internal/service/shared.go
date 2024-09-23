package service

import (
	"database/sql"
	"fmt"
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
