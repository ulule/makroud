package sqlxx

import (
	"reflect"

	"github.com/ulule/sqlxx/reflectx"
)

// ----------------------------------------------------------------------------
// Model
// ----------------------------------------------------------------------------

// ToModel converts the given instance to a Model instance.
func ToModel(instance interface{}) Model {
	value, ok := instance.(reflect.Type)
	if ok {
		if value.Kind() == reflect.Slice {
			value = reflectx.GetIndirectType(value.Elem())
		} else {
			value = reflectx.GetIndirectType(value)
		}

		model, ok := reflect.New(value).Elem().Interface().(Model)
		if ok {
			return model
		}

		return nil
	}

	out := reflect.Indirect(reflect.ValueOf(instance))

	// Single instance.
	if out.IsValid() && out.Kind() == reflect.Struct {
		model, ok := out.Interface().(Model)
		if ok {
			return model
		}
		return nil
	}

	// Slice of instances.
	if out.Kind() == reflect.Slice {
		// Slice of pointers
		if out.Type().Elem().Kind() == reflect.Ptr {
			model, ok := reflect.New(out.Type().Elem().Elem()).Interface().(Model)
			if ok {
				return model
			}

			return nil
		}

		// Slice of values
		model, ok := reflect.New(out.Type().Elem()).Interface().(Model)
		if ok {
			return model
		}

		return nil
	}

	// Try to convert it to model.
	model, ok := reflect.New(out.Type()).Interface().(Model)
	if ok {
		return model
	}

	return nil
}
