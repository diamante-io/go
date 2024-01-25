package history

import (
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
)

var _ driver.Valuer = (*AccountDataValue)(nil)
var _ sql.Scanner = (*AccountDataValue)(nil)

// Scan base64 decodes into an []byte
func (t *AccountDataValue) Scan(src interface{}) error {
	decoded, err := base64.StdEncoding.DecodeString(src.(string))
	if err != nil {
		return err
	}

	*t = decoded
	return nil
}

// Value implements driver.Valuer
func (value AccountDataValue) Value() (driver.Value, error) {
	return driver.Value([]uint8(base64.StdEncoding.EncodeToString(value))), nil
}

func (value AccountDataValue) Base64() string {
	return base64.StdEncoding.EncodeToString(value)
}
