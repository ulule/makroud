package reflectx

import (
	"reflect"
)

// IsStruct returns true if the given instance is a struct.
func IsStruct(instance interface{}) bool {
	value, ok := instance.(reflect.Value)
	if !ok {
		value = reflect.ValueOf(instance)
	}

	return value.Kind() == reflect.Struct
}
