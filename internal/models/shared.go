package models

import "database/sql"

func NullStringToString(str sql.NullString) string {
	if !str.Valid {
		return ""
	}
	return str.String
}
