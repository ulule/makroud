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

// IsPointer returns if given instance is a pointer, and not a nil one.
func IsPointer(instance interface{}) bool {
	val, ok := instance.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(instance)
	}

	return val.Kind() == reflect.Ptr && !val.IsNil()
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

// GetReflectPointerType returns a reflect pointer from given value of first level.
//
// For example:
//
//  * Type "Foo" will returns "*Foo"
//  * Type "*Foo" will returns "*Foo"
//  * Type "**Foo" will returns "*Foo"
//
func GetReflectPointerType(instance reflect.Value) reflect.Type {
	return reflect.PtrTo(GetIndirectType(instance))
}
