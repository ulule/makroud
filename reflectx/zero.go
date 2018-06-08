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

// MakeZero returns a zero value for the given element.
func MakeZero(element interface{}) reflect.Value {
	t, ok := element.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(element)
	}

	return reflect.New(GetIndirectType(t)).Elem()
}
