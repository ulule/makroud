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

// // GetFieldValues returns values for the given field for struct or slice.
// func GetFieldValues(out interface{}, name string) ([]interface{}, error) {
// 	if IsSlice(out) {
// 		return GetFieldValuesInSlice(out, name)
// 	}
//
// 	v, err := GetFieldValue(out, name)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return []interface{}{v}, nil
// }

// // GetFieldValuesInSlice returns values for the given field in slice of structs.
// func GetFieldValuesInSlice(slc interface{}, field string) ([]interface{}, error) {
// 	var (
// 		value  = reflect.ValueOf(slc).Elem()
// 		slcLen = value.Len()
// 		values []interface{}
// 	)
//
// 	for i := 0; i < slcLen; i++ {
// 		v, err := GetFieldValue(value.Index(i).Interface(), field)
// 		if err != nil {
// 			return nil, err
// 		}
// 		values = append(values, v)
// 	}
//
// 	return values, nil
// }

// ----------------------------------------------------------------------------
// Reflection
// ----------------------------------------------------------------------------

//
// // GetValues will extract given list of args from params using a sqlx mapper.
// // Example: Let's say we have this struct,
// //
// //     type User struct {
// //         IsActive bool
// //     }
// //
// // Using GetValues([]string{"is_active"}, &User{IsActive: true}, stmt.Mapper) will return map[is_active:true]
// func GetValues(args []string, params interface{}, mapper *reflectx.Mapper) (map[string]interface{}, error) {
// 	m, ok := params.(map[string]interface{})
// 	if ok {
// 		return getValuesWithMap(args, m)
// 	}
// 	return getValuesWithReflect(args, params, mapper)
// }
//
// // getValuesWithMap will filter given map with given list of args.
// func getValuesWithMap(args []string, params map[string]interface{}) (map[string]interface{}, error) {
// 	values := make(map[string]interface{})
//
// 	for i := range args {
// 		value, ok := params[args[i]]
// 		if !ok {
// 			return values, errors.Errorf("could not find name %s in %#v", args[i], params)
// 		}
// 		values[args[i]] = value
// 	}
//
// 	return values, nil
// }
//
// // getValuesWithReflect will analyze given struct to extract values from given list of args.
// func getValuesWithReflect(args []string, params interface{}, mapper *reflectx.Mapper) (map[string]interface{}, error) {
// 	values := make(map[string]interface{})
// 	source := GetIndirectValue(params)
//
// 	fields := mapper.TraversalsByName(source.Type(), args)
// 	for i, field := range fields {
// 		if len(field) == 0 {
// 			return values, errors.Errorf("could not find name %s in %#v", args[i], params)
// 		}
// 		value := reflectx.FieldByIndexesReadOnly(source, field)
// 		values[args[i]] = value.Interface()
// 	}
//
// 	return values, nil
// }
