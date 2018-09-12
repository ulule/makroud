package reflectx

import (
	"reflect"
)

// GetIndirectType returns indirect type for the given element.
func GetIndirectType(element interface{}) reflect.Type {
	value := element

	v, ok := value.(reflect.Value)
	if ok {
		value = v.Type()
	}

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

// GetFlattenValue returns a zero value from a flattened indirect type.
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
	return GetFlattenReflectValue(element).Interface()
}

// GetFlattenReflectValue returns a reflect value from a flattened indirect type.
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
func GetFlattenReflectValue(element interface{}) reflect.Value {
	v, ok := element.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(element)
	}
	if v.Kind() != reflect.Ptr {
		return v
	}
	for {
		e := v.Elem()
		if e.Kind() != reflect.Ptr {
			return v
		}
		v = e
	}
}
