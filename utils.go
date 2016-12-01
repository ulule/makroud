package sqlxx

import "reflect"

// ----------------------------------------------------------------------------
// Reflecters
// ----------------------------------------------------------------------------

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

// reflectIndirectType returns indirect type for the given type.
func reflectIndirectType(typ reflect.Type) reflect.Type {
	if typ.Kind() == reflect.Ptr {
		return reflectIndirectType(typ.Elem())
	}

	return typ
}

// ----------------------------------------------------------------------------
// Checkers
// ----------------------------------------------------------------------------

// isZeroValue returns true if the given interface is a zero value.
func isZeroValue(itf interface{}) bool {
	v := reflect.ValueOf(itf)

	// Avoid call of reflect.Value.Interface on zero Value
	if !v.IsValid() {
		return true
	}

	return reflect.Indirect(v).Interface() == reflect.Zero(reflect.Indirect(v).Type()).Interface()
}

// isSlice returns true if the given interface is a slice.
func isSlice(itf interface{}) bool {
	return reflectType(itf).Kind() == reflect.Slice
}

// ----------------------------------------------------------------------------
// Builders
// ----------------------------------------------------------------------------

// makeModel returns model type.
func makeModel(typ reflect.Type) Model {
	if typ.Kind() == reflect.Slice {
		typ = reflectIndirectType(typ.Elem())
	} else {
		typ = reflectIndirectType(typ)
	}

	if model, isModel := reflect.New(typ).Elem().Interface().(Model); isModel {
		return model
	}

	return nil
}

// makeSlice takes a type and returns create a slice from.
func makeSlice(itf interface{}) interface{} {
	sliceType := reflect.SliceOf(reflectType(reflectModel(itf)))

	slice := reflect.New(sliceType)
	slice.Elem().Set(reflect.MakeSlice(sliceType, 0, 0))

	return slice.Elem().Interface()
}

// modelToInterface returns an interface from a model.
func modelToInterface(model Model, isMany bool) interface{} {
	if isMany {
		return reflect.New(reflect.TypeOf(makeSlice(model))).Interface()
	}

	return reflect.New(reflect.TypeOf(model)).Interface()
}
