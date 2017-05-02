package sqlxx

import (
	"fmt"
	"reflect"
)

// IntToInt64 converts given int to int64.
func IntToInt64(value interface{}) (int64, error) {
	var (
		int64Type = reflect.TypeOf(int64(0))
		v         = reflect.Indirect(reflect.ValueOf(value))
	)

	if !v.IsValid() {
		return 0, fmt.Errorf("invalid value: %v", value)
	}

	if !v.Type().ConvertibleTo(int64Type) {
		return 0, fmt.Errorf("unable to convert %v to int64", v.Type())
	}

	return v.Convert(int64Type).Int(), nil
}
