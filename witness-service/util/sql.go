// Copyright (c) 2024 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

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
