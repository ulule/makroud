package sqlxx

import (
	"reflect"

	"github.com/ulule/sqlxx/reflekt"
)

// GetModelFromInterface returns interface as a Model interface.
func GetModelFromInterface(itf interface{}) Model {
	value := reflekt.GetIndirectValue(itf)

	// Single instance
	if value.IsValid() && value.Kind() == reflect.Struct {
		return value.Interface().(Model)
	}

	// Slice of instances
	if value.Kind() == reflect.Slice {
		// Slice of pointers
		if value.Type().Elem().Kind() == reflect.Ptr {
			return reflect.New(value.Type().Elem().Elem()).Interface().(Model)
		}
		// Slice of values
		return reflect.New(value.Type().Elem()).Interface().(Model)
	}

	return reflect.New(value.Type()).Interface().(Model)
}

// GetModelFromType returns model type.
func GetModelFromType(typ reflect.Type) Model {
	if typ.Kind() == reflect.Slice {
		typ = reflekt.GetIndirectType(typ.Elem())
	} else {
		typ = reflekt.GetIndirectType(typ)
	}

	if model, isModel := reflect.New(typ).Elem().Interface().(Model); isModel {
		return model
	}

	return nil
}
