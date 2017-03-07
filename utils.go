package sqlxx

import (
	"fmt"
	"reflect"

	"github.com/ulule/sqlxx/reflekt"
)

// InterfaceToModel returns interface as a Model interface.
func InterfaceToModel(itf interface{}) Model {
	var (
		value = reflekt.ReflectValue(itf)
		kind  = value.Kind()
	)

	// Single instance
	if value.IsValid() && kind == reflect.Struct {
		return value.Interface().(Model)
	}

	// Slice of instances
	if kind == reflect.Slice {
		// Slice of pointers
		if value.Type().Elem().Kind() == reflect.Ptr {
			return reflect.New(value.Type().Elem().Elem()).Interface().(Model)
		}
		// Slice of values
		return reflect.New(value.Type().Elem()).Interface().(Model)
	}

	return reflect.New(value.Type()).Interface().(Model)
}

// TypeToModel returns model type.
func TypeToModel(typ reflect.Type) Model {
	if typ.Kind() == reflect.Slice {
		typ = reflekt.ReflectIndirectType(typ.Elem())
	} else {
		typ = reflekt.ReflectIndirectType(typ)
	}

	if model, isModel := reflect.New(typ).Elem().Interface().(Model); isModel {
		return model
	}

	return nil
}

// InterfaceToSchema returns Schema by reflecting model for the given interface.
func InterfaceToSchema(out interface{}) (Schema, error) {
	return GetSchema(InterfaceToModel(out))
}

// IntToInt64 converts given int to int64.
func IntToInt64(value interface{}) (int64, error) {
	var (
		int64Type = reflect.TypeOf(int64(0))
		v         = reflect.Indirect(reflect.ValueOf(value))
	)

	if !v.IsValid() {
		return 0, fmt.Errorf("invalid value: %v", value)
	}

	if !v.Type().ConvertibleTo(int64Type) {
		return 0, fmt.Errorf("unable to convert %v to int64", v.Type())
	}

	return v.Convert(int64Type).Int(), nil
}
