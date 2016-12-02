package reflekt

import "reflect"

// ReflectValue returns the value that the interface v contains
// or that the pointer v points to.
func ReflectValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Interface && !v.IsNil() {
		if v.Elem().Kind() == reflect.Ptr && !v.Elem().IsNil() && v.Elem().Elem().Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	return v
}

// ReflectType returns type of the given interface.
func ReflectType(itf interface{}) reflect.Type {
	typ := reflect.ValueOf(itf).Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

// ReflectIndirectType returns indirect type for the given type.
func ReflectIndirectType(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Ptr {
		return ReflectIndirectType(typ.Elem())
	}

	return typ
}

// IsZeroValue returns true if the given interface is a zero value.
func IsZeroValue(itf interface{}) bool {
	v := reflect.ValueOf(itf)

	// Avoid call of reflect.Value.Interface on zero Value
	if !v.IsValid() {
		return true
	}

	return reflect.Indirect(v).Interface() == reflect.Zero(reflect.Indirect(v).Type()).Interface()
}

// IsSlice returns true if the given interface is a slice.
func IsSlice(itf interface{}) bool {
	return ReflectType(itf).Kind() == reflect.Slice
}

// MakeSlice takes a type and returns create a slice from.
func MakeSlice(itf interface{}) interface{} {
	sliceType := reflect.SliceOf(ReflectType(itf))

	slice := reflect.New(sliceType)
	slice.Elem().Set(reflect.MakeSlice(sliceType, 0, 0))

	return slice.Elem().Interface()
}
