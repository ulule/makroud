package reflectx

import (
	"reflect"
)

// MakePointer makes a copy of the given interface and returns a pointer.
func MakePointer(instance interface{}) interface{} {
	t := reflect.TypeOf(instance)

	cp := reflect.New(t)
	cp.Elem().Set(reflect.ValueOf(instance))

	// Avoid double pointers...
	if t.Kind() == reflect.Ptr {
		return cp.Elem().Interface()
	}

	return cp.Interface()
}
