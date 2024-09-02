package models

import (
	"database/sql"
	"encoding/json"

	"github.com/sqlc-dev/pqtype"
)

func NullStringToString(str sql.NullString) string {
	if !str.Valid {
		return ""
	}
	return str.String
}

func NullRawMessageToRawMessage(rawMessage pqtype.NullRawMessage) json.RawMessage {
	if !rawMessage.Valid {
		return json.RawMessage{}
	}
	return rawMessage.RawMessage
}
