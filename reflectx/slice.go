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

// GetSliceType returns the type if slice.
func GetSliceType(instance interface{}) reflect.Type {
	return GetIndirectType(instance).Elem()
}

// NewSliceValue creates a new value for slice type.
func NewSliceValue(instance interface{}) interface{} {
	t := GetSliceType(instance)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return reflect.New(t).Interface()
}

// NewReflectSlice creates a new slice for type and returns it's pointer.
func NewReflectSlice(instance interface{}) reflect.Value {
	t, ok := instance.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(instance)
	}

	return reflect.New(reflect.SliceOf(t))
}

// AppendReflectSlice will append given element to  reflect slice.
func AppendReflectSlice(list reflect.Value, value interface{}) {
	target := list
	if list.Kind() == reflect.Ptr {
		target = list.Elem()
	}

	elem := target.Type().Elem()

	val, ok := value.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(value)
	}

	if elem.Kind() == reflect.Struct && val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if elem.Kind() == reflect.Ptr && val.Kind() == reflect.Struct && val.CanAddr() {
		val = val.Addr()
	}

	target.Set(reflect.Append(target, val))
}

// CopyReflectSlice will attach given reflect slice to the destination value.
func CopyReflectSlice(dest interface{}, list reflect.Value) {
	val, ok := dest.(reflect.Value)
	if !ok {
		val = reflect.ValueOf(dest)
	}
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for list.Kind() == reflect.Ptr {
		list = list.Elem()
	}

	val.Set(list)
}
