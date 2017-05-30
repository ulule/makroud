package reflectx

import (
	"reflect"
)

// GetIndirectType returns indirect type for the given element.
func GetIndirectType(element interface{}) reflect.Type {
	value := element
	for {
		t, ok := value.(reflect.Type)
		if !ok {
			t = reflect.TypeOf(value)
		}
		if t.Kind() != reflect.Ptr {
			return t
		}
		value = t.Elem()
	}
}

// GetIndirectTypeName returns indirect type name for the given element.
func GetIndirectTypeName(element interface{}) string {
	return GetIndirectType(element).Name()
}

// GetIndirectValue returns indirect value for the given element.
func GetIndirectValue(element interface{}) reflect.Value {
	v := reflect.ValueOf(element)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}
