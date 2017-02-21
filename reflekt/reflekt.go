package reflekt

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// ReflectValue returns the value that the interface v contains
// or that the pointer v points to.
func ReflectValue(itf interface{}) reflect.Value {
	v := reflect.ValueOf(itf)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}

	return v
}

// ReflectType returns type of the given interface.
func ReflectType(itf interface{}) reflect.Type {
	return ReflectIndirectType(reflect.ValueOf(itf).Type())
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

// GetFieldValues returns values for the given field for struct or slice.
func GetFieldValues(out interface{}, name string) ([]interface{}, error) {
	var values []interface{}

	if IsSlice(out) {
		value := reflect.ValueOf(out).Elem()

		for i := 0; i < value.Len(); i++ {
			item := value.Index(i).Interface()

			v, err := GetFieldValue(item, name)
			if err != nil {
				return nil, err
			}

			values = append(values, v)
		}

		return values, nil
	}

	v, err := GetFieldValue(out, name)
	if err != nil {
		return nil, err
	}

	values = append(values, v)

	return values, nil
}

// CloneType returns a new interface from a given interface.
func CloneType(itf interface{}, args ...reflect.Kind) interface{} {
	var kind reflect.Kind

	if len(args) > 0 {
		kind = args[0]
	}

	if kind == reflect.Slice {
		return reflect.New(reflect.TypeOf(MakeSlice(itf))).Interface()
	}

	return reflect.New(reflect.TypeOf(itf)).Interface()
}

// Copy makes a copy of the given interface.
func Copy(itf interface{}) interface{} {
	cp := reflect.New(ReflectType(itf))
	cp.Elem().Set(reflect.ValueOf(itf))
	return cp.Interface()
}

// GetFieldValue returns the value
func GetFieldValue(itf interface{}, name string) (interface{}, error) {
	value, ok := itf.(reflect.Value)
	if !ok {
		value = ReflectValue(itf)
	}

	field := value.FieldByName(name)

	if !field.IsValid() {
		return nil, fmt.Errorf("No such field %s in %+v", name, itf)
	}

	return field.Interface(), nil
}

// SetFieldValue sets the provided value
func SetFieldValue(itf interface{}, name string, value interface{}) error {
	v, ok := itf.(reflect.Value)
	if !ok {
		v = ReflectValue(itf)
	}

	field := v.FieldByName(name)

	if !field.IsValid() {
		return fmt.Errorf("no such field %s in %+v", name, v.Interface())
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set %s field on %v%+v", name, v.Type().Name(), v.Interface())
	}

	fv := ReflectValue(value)

	if field.Type().Kind() == reflect.Ptr {
		fv = reflect.ValueOf(Copy(fv.Interface()))
	}

	if field.Type() != fv.Type() {
		return fmt.Errorf("provided value type %v didn't match field type %v", fv.Type(), field.Type())
	}

	field.Set(fv)

	return nil
}

// IsNullableType returns true if the given type is a nullable one.
func IsNullableType(t reflect.Type) bool {
	return t.ConvertibleTo(reflect.TypeOf((*driver.Valuer)(nil)).Elem())
}
