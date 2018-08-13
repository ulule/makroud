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
	v, ok := element.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(element)
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

// GetFlattenValue returns a zero value from GetFlattenType.
//
// For example:
//
//  * Type "Foo" will returns "Foo"...
//  * Type "*Foo" will returns "*Foo"...
//  * Type "**Foo" will returns "*Foo"...
//  * Type "***Foo" will returns "*Foo"...
//  * Type "[]Foo" will returns "[]Foo"...
//  * Type "*[]Foo" will returns "*[]Foo"...
//  * Type "**[]Foo" will returns "*[]Foo"...
//  * Type "***[]Foo" will returns "*[]Foo"...
//
func GetFlattenValue(element interface{}) interface{} {
	t, ok := element.(reflect.Value)
	if !ok {
		t = reflect.ValueOf(element)
	}
	if t.Kind() != reflect.Ptr {
		return t.Interface()
	}
	for {
		e := t.Elem()
		if e.Kind() != reflect.Ptr {
			return t.Interface()
		}
		t = e
	}
}
