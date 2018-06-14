package sqlxx

import (
	"reflect"

	"github.com/ulule/sqlxx/reflectx"
)

// ----------------------------------------------------------------------------
// Model
// ----------------------------------------------------------------------------

// ToModel converts the given instance to a Model instance.
func ToModel(itf interface{}) Model {
	typ, ok := itf.(reflect.Type)
	if ok {
		if typ.Kind() == reflect.Slice {
			typ = reflectx.GetIndirectType(typ.Elem())
		} else {
			typ = reflectx.GetIndirectType(typ)
		}

		model, ok := reflect.New(typ).Elem().Interface().(Model)
		if ok {
			return model
		}

		return nil
	}

	value := reflect.Indirect(reflect.ValueOf(itf))

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

// IsSlice returns true if the given interface is a slice.
func IsSlice(itf interface{}) bool {
	return reflectx.GetIndirectType(reflect.ValueOf(itf).Type()).Kind() == reflect.Slice
}

// MakePointer makes a copy of the given interface and returns a pointer.
func MakePointer(itf interface{}) interface{} {
	t := reflect.TypeOf(itf)

	cp := reflect.New(t)
	cp.Elem().Set(reflect.ValueOf(itf))

	// Avoid double pointers if itf is a pointer
	if t.Kind() == reflect.Ptr {
		return cp.Elem().Interface()
	}

	return cp.Interface()
}

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
