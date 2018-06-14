package reflectx

import (
	"reflect"
)

// IsZero returns true if the given interface is a zero value or nil.
func IsZero(instance interface{}) bool {
	if instance == nil {
		return true
	}

	value, ok := instance.(reflect.Value)
	if !ok {
		value = reflect.ValueOf(instance)
	}
	if value.Kind() == reflect.Ptr && value.IsNil() {
		return true
	}

	value = reflect.Indirect(value)

	zero := reflect.Zero(value.Type())
	if value.Type().Comparable() && zero.Type().Comparable() {
		return value.Interface() == zero.Interface()
	}

	return reflect.DeepEqual(value.Interface(), zero.Interface())
}

// MakeZero returns a zero value for the given element.
func MakeZero(element interface{}) reflect.Value {
	t, ok := element.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(element)
	}

	return reflect.New(GetIndirectType(t)).Elem()
}
