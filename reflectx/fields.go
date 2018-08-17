package reflectx

import (
	"database/sql"
	"reflect"

	"github.com/pkg/errors"
)

import "fmt"

// GetFields returns a list of field name.
func GetFields(element interface{}) ([]string, error) {
	dest := reflect.TypeOf(element)
	if dest == nil || (dest.Kind() != reflect.Ptr && dest.Kind() != reflect.Struct) {
		return nil, errors.New("sqlxx: cannot find fields on a non-struct interface")
	}

	value := reflect.ValueOf(element)
	if dest.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	count := value.Type().NumField()
	fields := []string{}
	for i := 0; i < count; i++ {
		field := value.Type().Field(i)
		// Ignore private or anonymous field...
		if field.PkgPath == "" {
			fields = append(fields, field.Name)
		}
	}

	return fields, nil
}

// GetFieldByName returns the field in element with given name.
func GetFieldByName(element interface{}, name string) (reflect.StructField, bool) {
	return reflect.Indirect(reflect.ValueOf(element)).Type().FieldByName(name)
}

// GetFieldValue returns the field's value with given name.
func GetFieldValue(element interface{}, name string) (interface{}, error) {
	value, ok := element.(reflect.Value)
	if !ok {
		value = reflect.Indirect(reflect.ValueOf(element))
	}

	if value.Kind() == reflect.Interface {
		value = reflect.ValueOf(value.Interface())
	}

	// Avoid calling FieldByName on pointer.
	value = reflect.Indirect(value)

	// Avoid calling FieldByName on zero value
	if !value.IsValid() {
		return nil, errors.Errorf("sqlxx: no such field %s in %T", name, element)
	}

	field := value.FieldByName(name)
	if !field.IsValid() {
		return nil, errors.Errorf("sqlxx: no such field %s in %T", name, element)
	}
	if field.Kind() == reflect.Ptr && field.IsNil() {
		return nil, nil
	}
	if !field.CanInterface() {
		return nil, errors.Errorf("sqlxx: cannot find field %s in %T", name, element)
	}

	return field.Interface(), nil
}

// GetFieldValueInt64 returns int64 value for the given instance field.
func GetFieldValueInt64(instance interface{}, field string) (int64, error) {
	value, err := GetFieldValue(instance, field)
	if err != nil {
		return 0, err
	}

	converted, err := ToInt64(value)
	if err != nil {
		return 0, err
	}

	return converted, nil
}

// GetFieldOptionalValueInt64 returns an optional int64 value for the given instance field.
func GetFieldOptionalValueInt64(instance interface{}, field string) (int64, bool, error) {
	value, err := GetFieldValue(instance, field)
	if err != nil {
		return 0, false, err
	}

	converted, ok := ToOptionalInt64(value)
	return converted, ok, nil
}

// GetFieldValueString returns string value for the given instance field.
func GetFieldValueString(instance interface{}, field string) (string, error) {
	value, err := GetFieldValue(instance, field)
	if err != nil {
		return "", err
	}

	converted, err := ToString(value)
	if err != nil {
		return "", err
	}

	return converted, nil
}

// GetFieldOptionalValueString returns an optional string value for the given instance field.
func GetFieldOptionalValueString(instance interface{}, field string) (string, bool, error) {
	value, err := GetFieldValue(instance, field)
	if err != nil {
		return "", false, err
	}

	converted, ok := ToOptionalString(value)
	return converted, ok, nil
}

// UpdateFieldValue updates the field's value with given name.
func UpdateFieldValue(instance interface{}, name string, value interface{}) error {
	v, ok := instance.(reflect.Value)
	if !ok {
		v = reflect.Indirect(reflect.ValueOf(instance))
	}

	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
	}

	dest := v.FieldByName(name)
	if !dest.IsValid() {
		return errors.Errorf("sqlxx: no such field %s in %T", name, instance)
	}
	if !dest.CanSet() {
		return errors.Errorf("sqlxx: cannot update field %s in %T", name, instance)
	}

	// Try scanner interface.
	scanner := reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	scan := dest
	if scan.CanAddr() && scan.Type().Kind() != reflect.Ptr {
		scan = scan.Addr()
	}
	if scan.Type().Implements(scanner) {

		// If scanner pointer is nil, allocate a zero value.
		if scan.IsNil() {
			zero := MakeZero(scan.Type().Elem())
			scan.Set(zero.Addr())
		}

		// If value is nil, creates a zero value.
		arg := reflect.Indirect(reflect.ValueOf(value))
		if value == nil {
			arg = reflect.ValueOf(MakeZero(scan.Type()))
		}

		// And try to use reflection to call Scan method.
		values := scan.MethodByName("Scan").Call([]reflect.Value{arg})
		if len(values) == 1 {
			okType := reflect.TypeOf((*error)(nil)).Elem() == values[0].Type()
			okValue := values[0].IsNil()
			if okType && okValue {
				return nil
			}
		}
	}

	// Otherwise, try to manually update field using reflection.
	output := reflect.Indirect(reflect.ValueOf(value))
	if !output.IsValid() {
		return errors.Errorf("sqlxx: cannot uses %T as value to update %s in %T", value, name, instance)
	}

	// If field's type is a pointer, create a pointer of the given value.
	if dest.Type().Kind() == reflect.Ptr {
		output = reflect.ValueOf(MakePointer(output.Interface()))
	}

	// Verify that types are equals. Otherwise, returns a nice error.
	if dest.Type() != output.Type() {
		return errors.Errorf("sqlxx: cannot use type %v to update type %v in %T", output.Type(), dest.Type(), instance)
	}

	dest.Set(output)

	return nil
}

// PushFieldValue updates the field's value with given name.
// TODO Better comments
func PushFieldValue(instance interface{}, name string, value interface{}) error {
	v, ok := instance.(reflect.Value)
	if !ok {
		v = reflect.Indirect(reflect.ValueOf(instance))
	}

	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
	}

	dest := v.FieldByName(name)
	if !dest.IsValid() {
		return errors.Errorf("sqlxx: no such field %s in %T", name, instance)
	}
	if !dest.CanSet() {
		return errors.Errorf("sqlxx: cannot update field %s in %T", name, instance)
	}

	// Otherwise, try to manually update field using reflection.
	output := reflect.Indirect(reflect.ValueOf(value))
	if !output.IsValid() {
		return errors.Errorf("sqlxx: cannot uses %T as value to update %s in %T", value, name, instance)
	}

	if dest.Kind() == reflect.Slice {
		fmt.Println(dest.Type().Kind())
		fmt.Println("::4")
		fmt.Println(output)
		AppendReflectSlice(dest, output)
		return nil
	}

	// If field's type is a pointer, create a pointer of the given value.
	if dest.Type().Kind() == reflect.Ptr {
		output = reflect.ValueOf(MakePointer(output.Interface()))
	}

	// Verify that types are equals. Otherwise, returns a nice error.
	if dest.Type() != output.Type() {
		return errors.Errorf("sqlxx: cannot use type %v to update type %v in %T", output.Type(), dest.Type(), instance)
	}

	dest.Set(output)

	return nil
}
