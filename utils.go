package sqlxx

import (
	"reflect"

	"github.com/pkg/errors"

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

// SetFieldValue sets the provided value
func SetFieldValue(itf interface{}, name string, value interface{}) error {
	v, ok := itf.(reflect.Value)
	if !ok {
		v = reflect.Indirect(reflect.ValueOf(itf))
	}

	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
	}

	field := v.FieldByName(name)
	if !field.IsValid() {
		return errors.Errorf("no such field %s in %+v", name, v.Interface())
	}

	if !field.CanSet() {
		return errors.Errorf("cannot set %s field on %v%+v", name, v.Type().Name(), v.Interface())
	}

	fv := reflect.Indirect(reflect.ValueOf(value))
	if !fv.IsValid() {
		return nil
	}

	if field.Type().Kind() == reflect.Ptr {
		fv = reflect.ValueOf(MakePointer(fv.Interface()))
	}

	if field.Type() != fv.Type() {
		return errors.Errorf("provided value type %v didn't match field type %v", fv.Type(), field.Type())
	}

	field.Set(fv)

	return nil
}

// ----------------------------------------------------------------------------
// Reflection
// ----------------------------------------------------------------------------

// IsSlice returns true if the given interface is a slice.
func IsSlice(itf interface{}) bool {
	return reflectx.GetIndirectType(reflect.ValueOf(itf).Type()).Kind() == reflect.Slice
}

// GetZero returns a zero value for the given element.
func GetZero(element interface{}) reflect.Value {
	t, ok := element.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(element)
	}

	return reflect.New(reflectx.GetIndirectType(t)).Elem()
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
