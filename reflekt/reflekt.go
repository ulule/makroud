package reflekt

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

// ----------------------------------------------------------------------------
// Indirect helpers
// ----------------------------------------------------------------------------

// GetIndirectValue returns the value that the interface v contains or that the pointer v points to.
func GetIndirectValue(itf interface{}) reflect.Value {
	v, ok := itf.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(itf)
	}

	v = reflect.Indirect(v)

	if v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}

	return v
}

// GetIndirectType returns indirect type for the given type.
func GetIndirectType(itf interface{}) reflect.Type {
	t, ok := itf.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(itf)
	}

	if t.Kind() == reflect.Ptr {
		return GetIndirectType(t.Elem())
	}

	return t
}

// ----------------------------------------------------------------------------
// Cloners
// ----------------------------------------------------------------------------

// MakeSlice takes a type and returns create a slice from.
func MakeSlice(itf interface{}) interface{} {
	sliceType := reflect.SliceOf(GetIndirectType(itf))
	slice := reflect.New(sliceType)
	slice.Elem().Set(reflect.MakeSlice(sliceType, 0, 0))
	return slice.Elem().Interface()
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
	cp := reflect.New(GetIndirectType(reflect.TypeOf(itf)))
	cp.Elem().Set(reflect.ValueOf(itf))
	return cp.Interface()
}

// ----------------------------------------------------------------------------
// Struct fields
// ----------------------------------------------------------------------------

// GetFieldValuesInSlice returns values for the given field in slice of structs.
func GetFieldValuesInSlice(slc interface{}, field string) ([]interface{}, error) {
	var (
		value  = reflect.ValueOf(slc).Elem()
		slcLen = value.Len()
		values []interface{}
	)

	for i := 0; i < slcLen; i++ {
		v, err := GetFieldValue(value.Index(i).Interface(), field)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}

	return values, nil
}

// GetFieldValues returns values for the given field for struct or slice.
func GetFieldValues(out interface{}, name string) ([]interface{}, error) {
	if IsSlice(out) {
		return GetFieldValuesInSlice(out, name)
	}

	v, err := GetFieldValue(out, name)
	if err != nil {
		return nil, err
	}

	return []interface{}{v}, nil
}

// GetFieldValue returns the value
func GetFieldValue(itf interface{}, name string) (interface{}, error) {
	value, ok := itf.(reflect.Value)
	if !ok {
		value = GetIndirectValue(itf)
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
		v = GetIndirectValue(itf)
	}

	field := v.FieldByName(name)
	if !field.IsValid() {
		return fmt.Errorf("no such field %s in %+v", name, v.Interface())
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set %s field on %v%+v", name, v.Type().Name(), v.Interface())
	}

	fv := GetIndirectValue(value)
	if field.Type().Kind() == reflect.Ptr {
		fv = reflect.ValueOf(Copy(fv.Interface()))
	}

	if field.Type() != fv.Type() {
		return fmt.Errorf("provided value type %v didn't match field type %v", fv.Type(), field.Type())
	}

	field.Set(fv)

	return nil
}

// ----------------------------------------------------------------------------
// Checkers
// ----------------------------------------------------------------------------

// IsNullableType returns true if the given type is a nullable one.
func IsNullableType(t reflect.Type) bool {
	return t.ConvertibleTo(reflect.TypeOf((*driver.Valuer)(nil)).Elem())
}

// IsSlice returns true if the given interface is a slice.
func IsSlice(itf interface{}) bool {
	return GetIndirectType(reflect.ValueOf(itf).Type()).Kind() == reflect.Slice
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
