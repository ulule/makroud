package sqlxx

import (
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
