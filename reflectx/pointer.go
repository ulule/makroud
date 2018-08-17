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

// IsPointer returns if given instance is a pointer.
func IsPointer(instance interface{}) bool {
	val, ok := instance.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(instance)
	}

	return val.Kind() == reflect.Ptr
}

// MakeReflectPointer makes a pointer from given reflect value.
func MakeReflectPointer(instance reflect.Value) reflect.Value {
	t := instance.Type()

	cp := reflect.New(t)
	cp.Elem().Set(instance)

	// Avoid double pointers...
	if t.Kind() == reflect.Ptr {
		return cp.Elem()
	}

	return cp
}

// CreateReflectPointer creates a reflect pointer from given value.
func CreateReflectPointer(instance interface{}) reflect.Value {
	return MakeReflectPointer(reflect.ValueOf(instance))
}
