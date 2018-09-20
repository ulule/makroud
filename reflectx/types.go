package reflectx

import (
	"database/sql"
	"database/sql/driver"
	"reflect"

	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

var (
	int64Type      = reflect.TypeOf(int64(0))
	nullInt64Type  = reflect.TypeOf(sql.NullInt64{})
	stringType     = reflect.TypeOf("")
	nullStringType = reflect.TypeOf(sql.NullString{})
	scannerType    = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
)

// ----------------------------------------------------------------------------
// Int64
// ----------------------------------------------------------------------------

// ToInt64 converts given value to int64.
// In case of a optional value (ie: sql.NullInt64), please use ToOptionalInt64.
func ToInt64(value interface{}) (int64, error) {
	cast, ok := value.(int64)
	if ok {
		return cast, nil
	}

	// For sql.NullInt64 and sql.NullFloat64 support.
	valuer, ok := value.(driver.Valuer)
	if ok {
		v, err := valuer.Value()
		if err != nil {
			return 0, errors.Wrap(err, "cannot convert to int64")
		}
		if v == nil {
			return 0, errors.Errorf("cannot convert to int64")
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

// ToOptionalInt64 try to converts given value to int64.
func ToOptionalInt64(value interface{}) (int64, bool) {
	v, err := ToInt64(value)
	return v, err == nil
}

// ----------------------------------------------------------------------------
// String
// ----------------------------------------------------------------------------

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
		if err != nil {
			return "", errors.Wrap(err, "cannot convert to string")
		}
		if v == nil {
			return "", errors.Errorf("cannot convert to string")
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

// ToOptionalString try to converts given value to string.
func ToOptionalString(value interface{}) (string, bool) {
	v, err := ToString(value)
	return v, err == nil
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

// Type defines high level types used by reflectx. It's a subset of go types.
type Type uint8

const (
	// UnsupportedType is an unsupported type.
	UnsupportedType = Type(iota)
	// Int64Type uses an integer.
	Int64Type
	// StringType uses a string.
	StringType
	// OptionalInt64Type uses an optional integer.
	OptionalInt64Type
	// OptionalStringType uses an optional string.
	OptionalStringType
)

// GetType returns high level type from given reflect type.
func GetType(value reflect.Type) Type {
	indirect := GetIndirectType(value)
	switch indirect.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return Int64Type

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return Int64Type

	case reflect.String:
		return StringType

	case reflect.Struct:
		if indirect == nullStringType {
			return OptionalStringType
		}
		if indirect == nullInt64Type {
			return OptionalInt64Type
		}
	}
	return UnsupportedType
}
