package reflectx

import (
	"reflect"

	"github.com/pkg/errors"
)

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
	return PushFieldValue(instance, name, value, true)
}

// PushFieldValue updates the field value identified by given name using a lazy-type push mechanism.
// This mechanism means that it will try to match structs with pointers, slices or literals,
// and try to replace it's value.
// If strict mode is disabled, it will even append given value if field is a slice.
func PushFieldValue(instance interface{}, name string, value interface{}, strict bool) error {
	// Find field identified by given name.
	dest, err := getDestinationReflectValue(instance, name)
	if err != nil {
		return err
	}

	// Wraps given value in a reflect wrapper.
	output, err := getOutputReflectValue(instance, name, value)
	if err != nil {
		return err
	}

	// If output is nil, do nothing...
	if output.Kind() == reflect.Ptr && output.IsNil() || !output.IsValid() {
		return nil
	}

	dtype := dest.Type()
	dkind := dtype.Kind()

	switch dkind {
	case reflect.Ptr:
		elem := dtype.Elem()

		switch elem.Kind() {
		case reflect.Slice:

			// If destination field is a slice, initialize it if it's nil and append given value on slice.
			if dest.IsNil() {
				list := NewReflectSlice(elem.Elem())
				dest.Set(list)
			}

			return updateFieldValueOnSlice(dest, output, instance, strict)

		case reflect.Struct:

			// Try to use scanner interface if available.
			if tryScannerOnFieldValue(dest, output) {
				return nil
			}

			// If destination field is a struct, create a pointer of the given value.
			output = MakeReflectPointer(output)

			return setFieldValueOnStruct(dest, output, instance)

		case reflect.Ptr:

			// Try to use scanner interface if available.
			if tryScannerOnFieldValue(dest, output) {
				return nil
			}

			// If destination field is a pointer, just forward given value.
			return setFieldValueOnStruct(dest, output, instance)

		default:

			// Try to use scanner interface if available.
			if tryScannerOnFieldValue(dest, output) {
				return nil
			}

			// If value is not a pointer, create a pointer of the given value.
			output = MakeReflectPointer(output)

			return setFieldValueOnStruct(dest, output, instance)
		}

	case reflect.Slice:
		return updateFieldValueOnSlice(dest, output, instance, strict)

	case reflect.Struct:

		// Try to use scanner interface if available.
		if tryScannerOnFieldValue(dest, output) {
			return nil
		}

		// If value is a pointer, forward it's indirect value.
		if output.Kind() == reflect.Ptr && !output.IsNil() {
			output = output.Elem()
		}

		return setFieldValueOnStruct(dest, output, instance)

	default:

		// Try to use scanner interface if available.
		if tryScannerOnFieldValue(dest, output) {
			return nil
		}

		// If value is a pointer, forward it's indirect value.
		if output.Kind() == reflect.Ptr && !output.IsNil() {
			output = output.Elem()
		}

		return setFieldValueOnStruct(dest, output, instance)
	}
}

// getDestinationReflectValue returns field value identified by given name from instance.
func getDestinationReflectValue(instance interface{}, name string) (dest reflect.Value, err error) {
	// Get a indirect value from given instance.
	v := GetIndirectValue(instance)

	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
	}

	// Find field identified by given name.
	dest = v.FieldByName(name)
	if !dest.IsValid() {
		return dest, errors.Errorf("sqlxx: no such field %s in %T", name, instance)
	}
	if !dest.CanSet() {
		return dest, errors.Errorf("sqlxx: cannot update field %s in %T", name, instance)
	}

	return dest, nil
}

// getOutputReflectValue returns a reflect value used to update a field in given instance.
func getOutputReflectValue(instance interface{}, name string, value interface{}) (output reflect.Value, err error) {
	output = reflect.ValueOf(value)
	if !output.IsValid() && value != nil {
		return output, errors.Errorf("sqlxx: cannot uses %T as value to update %s in %T", value, name, instance)
	}
	return output, nil
}

func appendFieldValueOnSlice(dest reflect.Value, output reflect.Value, instance interface{}) error {
	sdtype := dest.Type()
	sotype := output.Type()

	// Get indirect type from value: *Foobar -> Foobar
	otype := sotype
	for otype.Kind() == reflect.Ptr {
		otype = otype.Elem()
	}

	// Get indirect type from slice: *[]Foobar -> []Foobar
	dtype := sdtype
	for dtype.Kind() == reflect.Ptr {
		dtype = dtype.Elem()
	}

	// Get indirect type from slice subtype: []*Foobar -> []Foobar
	dtype = dtype.Elem()
	for dtype.Kind() == reflect.Ptr {
		dtype = dtype.Elem()
	}

	// Verify that types are equals. Otherwise, returns a nice error.
	if dtype != otype {
		return errors.Errorf("sqlxx: cannot use type %v to update type %v in %T", sotype, sdtype, instance)
	}

	AppendReflectSlice(dest, output)
	return nil
}

func setFieldValueOnStruct(dest reflect.Value, output reflect.Value, instance interface{}) error {
	otype := output.Type()
	dtype := dest.Type()

	// Verify that types are equals. Otherwise, returns a nice error.
	if dtype != otype {
		return errors.Errorf("sqlxx: cannot use type %v to update type %v in %T", otype, dtype, instance)
	}

	dest.Set(output)

	return nil
}

func tryScannerOnFieldValue(dest reflect.Value, output reflect.Value) bool {
	// Try scanner interface.
	scan := dest
	if scan.CanAddr() && scan.Type().Kind() != reflect.Ptr {
		scan = scan.Addr()
	}
	if !scan.Type().Implements(scannerType) {
		return false
	}

	// If scanner pointer is nil, allocate a zero value.
	if scan.IsNil() {
		zero := MakeZero(scan.Type().Elem())
		scan.Set(zero.Addr())
	}

	// If output is direct valid value for destination type, bypass scan method.
	dtype := dest.Type()
	if dtype.Kind() == reflect.Ptr {
		dtype = dtype.Elem()
	}
	otype := output.Type()
	if otype.Kind() == reflect.Ptr {
		otype = otype.Elem()
	}
	if otype == dtype {
		// Try to match dest and output type.
		if dest.Kind() != reflect.Ptr && output.Kind() == reflect.Ptr {
			output = output.Elem()
		}
		if dest.Kind() == reflect.Ptr && output.Kind() != reflect.Ptr {
			output = MakeReflectPointer(output)
		}
		dest.Set(output)
		return true
	}

	// If value is a pointer, get indirect value (or a zero value if nil).
	arg := output
	if arg.Kind() == reflect.Ptr {
		if arg.IsNil() {
			arg = MakeZero(scan.Type())
		} else {
			arg = output.Elem()
		}
	}

	// And try to use reflection to call Scan method.
	values := scan.MethodByName("Scan").Call([]reflect.Value{arg})
	if len(values) == 1 {
		okType := reflect.TypeOf((*error)(nil)).Elem() == values[0].Type()
		okValue := values[0].IsNil()
		if okType && okValue {
			return true
		}
	}

	return false
}

func updateFieldValueOnSlice(dest reflect.Value, output reflect.Value, instance interface{}, strict bool) error {

	// Try to use append mechanism if strict mode is disabled.
	if !strict {
		otype := output.Type()
		for otype.Kind() == reflect.Ptr {
			otype = otype.Elem()
		}
		if otype.Kind() != reflect.Slice {
			return appendFieldValueOnSlice(dest, output, instance)
		}
	}

	// Try to match dest and output type.
	if dest.Kind() != reflect.Ptr && output.Kind() == reflect.Ptr {
		output = output.Elem()
	}
	if dest.Kind() == reflect.Ptr && output.Kind() != reflect.Ptr {
		output = MakeReflectPointer(output)
	}

	return setFieldValueOnStruct(dest, output, instance)
}
