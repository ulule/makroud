package reflectx

import (
	"database/sql/driver"
	"reflect"

	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
// Int64
// ----------------------------------------------------------------------------

var int64Type = reflect.TypeOf(int64(0))

// ToInt64 converts given value to int64.
func ToInt64(value interface{}) (int64, error) {
	cast, ok := value.(int64)
	if ok {
		return cast, nil
	}

	// sql.NullInt64 support
	valuer, ok := value.(driver.Valuer)
	if ok {
		v, err := valuer.Value()
		if err != nil || v == nil {
			return 0, errors.Wrap(err, "cannot convert to int64")
		}

		value = v
	}

	reflected := reflect.Indirect(reflect.ValueOf(value))

	if !reflected.IsValid() {
		return 0, errors.Errorf("invalid value: %v", value)
	}

	if !reflected.Type().ConvertibleTo(int64Type) {
		return 0, errors.Errorf("unable to convert %v to int64", reflected.Type())
	}

	return reflected.Convert(int64Type).Int(), nil
}

// ----------------------------------------------------------------------------
// String
// ----------------------------------------------------------------------------

var stringType = reflect.TypeOf("")

// ToString converts given value to string.
func ToString(value interface{}) (string, error) {
	cast, ok := value.(string)
	if ok {
		return cast, nil
	}

	// sql.NullString support
	valuer, ok := value.(driver.Valuer)
	if ok {
		v, err := valuer.Value()
		if err != nil || v == nil {
			return "", errors.Wrap(err, "cannot convert to string")
		}

		value = v
	}

	reflected := reflect.Indirect(reflect.ValueOf(value))

	if !reflected.IsValid() {
		return "", errors.Errorf("invalid value: %v", value)
	}

	if !reflected.Type().ConvertibleTo(stringType) {
		return "", errors.Errorf("unable to convert %v to string", reflected.Type())
	}

	return reflected.Convert(stringType).String(), nil
}
