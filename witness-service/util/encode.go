package util

import (
	"database/sql"
	"encoding/hex"
)

// EncodeToNullString encodes a byte array into a sql.NullString
func EncodeToNullString(bytes []byte) sql.NullString {
	if bytes == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: hex.EncodeToString(bytes),
		Valid:  true,
	}
}

// DecodeNullString decodes a sql.NullString to bytes
func DecodeNullString(str sql.NullString) ([]byte, error) {
	if !str.Valid {
		return nil, nil
	}
	return hex.DecodeString(str.String)
}
