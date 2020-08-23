package query

import (
	"database/sql/driver"

	"errors"

	jsoniter "github.com/json-iterator/go"
)

//JSONRaw ...
type JSONRaw jsoniter.RawMessage

//Value ...
func (j JSONRaw) Value() (driver.Value, error) {
	byteArr := []byte(j)

	return driver.Value(byteArr), nil
}

// Scan ...
func (j *JSONRaw) Scan(src interface{}) error {
	var source []byte
	switch src.(type) {
	case string:
		source = []byte(src.(string))
	case []byte:
		source = src.([]byte)
	default:
		return errors.New("Incompatible type for SecureString")
	}

	*j = source

	return nil
}
