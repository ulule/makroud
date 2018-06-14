package reflectx

import (
	"reflect"
)

// IsSlice returns true if the given instance is a slice.
func IsSlice(instance interface{}) bool {
	value, ok := instance.(reflect.Value)
	if !ok {
		value = reflect.ValueOf(instance)
	}

	return GetIndirectType(value.Type()).Kind() == reflect.Slice
}
