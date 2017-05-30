package reflectx

import (
	"reflect"
)

// IsZero returns true if the given interface is a zero value or nil.
func IsZero(itf interface{}) bool {
	if itf == nil {
		return true
	}

	value, ok := itf.(reflect.Value)
	if !ok {
		value = reflect.Indirect(reflect.ValueOf(itf))
	}

	if value.Kind() == reflect.Ptr && value.IsNil() {
		return true
	}

	zero := reflect.Zero(value.Type())
	return value.Interface() == zero.Interface()
}
