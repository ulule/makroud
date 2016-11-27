package sqlxx

import "reflect"

// reflectValue returns the value that the interface v contains
// or that the pointer v points to.
func reflectValue(v reflect.Value) reflect.Value {
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

// reflectType returns type of the given interface.
func reflectType(itf interface{}) reflect.Type {
	typ := reflect.ValueOf(itf).Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

// reflectModel returns interface as a Model interface.
func reflectModel(itf interface{}) Model {
	value := reflectValue(reflect.ValueOf(itf))

	// Instance
	if value.IsValid() && value.Kind() == reflect.Struct {
		return value.Interface().(Model)
	}

	// Slice of models
	if value.Kind() == reflect.Slice {
		// Slice of model pointers
		if value.Type().Elem().Kind() == reflect.Ptr {
			return reflect.New(value.Type().Elem().Elem()).Interface().(Model)
		}

		// Slice of model values
		return reflect.New(value.Type().Elem()).Interface().(Model)
	}

	// Type
	if reflect.TypeOf(itf).Kind() == reflect.Ptr {
		typ := reflect.TypeOf(itf).Elem()

		// Struct
		if typ.Kind() == reflect.Struct {
			return reflect.New(typ).Interface().(Model)
		}

		// Slice
		return reflect.New(typ.Elem()).Interface().(Model)
	}

	return reflect.New(value.Type()).Interface().(Model)
}

// isZeroValue returns true if the given interface is a zero value.
func isZeroValue(itf interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(itf))
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}
