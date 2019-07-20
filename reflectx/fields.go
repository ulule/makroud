package reflectx

import (
	"reflect"

	"github.com/pkg/errors"
)

// GetFields returns a list of field name.
func GetFields(element interface{}) ([]string, error) {
	value := GetIndirectValue(element)
	if value.Kind() != reflect.Struct {
		return nil, errors.New("makroud: cannot find fields on a non-struct interface")
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

// GetFieldsCount returns the number of exported fields for given type.
func GetFieldsCount(element reflect.Type) int {
	if element.Kind() != reflect.Struct {
		return 0
	}

	count := 0
	max := element.NumField()
	for i := 0; i < max; i++ {
		field := element.Field(i)
		// Ignore private or anonymous field...
		if field.PkgPath == "" {
			count++
		}
	}

	return count
}

// GetFieldByName returns the field in element with given name.
func GetFieldByName(element interface{}, name string) (reflect.StructField, bool) {
	value := GetIndirectValue(element)
	if value.Kind() != reflect.Struct {
		return reflect.StructField{}, false
	}
	return value.Type().FieldByName(name)
}

// GetFieldReflectTypeByName returns the field's type with given name.
func GetFieldReflectTypeByName(element interface{}, name string) (reflect.Type, error) {
	value, ok := GetFieldByName(element, name)
	if !ok {
		return nil, errors.Errorf("makroud: no such field %s in %T", name, element)
	}

	kind := value.Type
	for kind.Kind() == reflect.Ptr {
		kind = kind.Elem()
	}

	return kind, nil
}

// GetReflectFieldByIndexes returns a pointer of field value with the given traversal indexes.
func GetReflectFieldByIndexes(value reflect.Value, indexes []int) interface{} {
	for _, i := range indexes {
		value = reflect.Indirect(value).Field(i)
		// If this is a pointer and it's nil, allocate a new value and set it.
		if value.Kind() == reflect.Ptr && value.IsNil() {
			instance := reflect.New(GetIndirectType(value))
			value.Set(instance)
		}
		if value.Kind() == reflect.Map && value.IsNil() {
			value.Set(reflect.MakeMap(value.Type()))
		}
	}
	return value.Addr().Interface()
}

// GetFieldValueWithIndexes returns the field's value with given traversal indexes.
func GetFieldValueWithIndexes(value reflect.Value, indexes []int) (interface{}, error) {
	// Avoid calling Field on interface.
	if value.Kind() == reflect.Interface {
		value = reflect.ValueOf(value.Interface())
	}

	// Avoid calling Field on pointer.
	for value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}

	for _, i := range indexes {
		value = reflect.Indirect(value).Field(i)
		if !value.IsValid() {
			return nil, errors.Errorf("makroud: cannot find required field in %T", value)
		}
		if value.Kind() == reflect.Ptr && value.IsNil() {
			return nil, nil
		}
	}

	if !value.CanInterface() {
		return nil, errors.Errorf("makroud: cannot find required field in %T", value)
	}

	return value.Interface(), nil
}

// GetFieldValueWithName returns the field's value with given name.
func GetFieldValueWithName(value reflect.Value, name string) (interface{}, error) {
	// Avoid calling FieldByName on interface.
	if value.Kind() == reflect.Interface {
		value = reflect.ValueOf(value.Interface())
	}

	// Avoid calling FieldByName on pointer.
	for value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}

	// Avoid calling FieldByName on zero value
	if !value.IsValid() {
		return nil, errors.Errorf("makroud: no such field %s in %T", name, value)
	}

	field := value.FieldByName(name)
	if !field.IsValid() {
		return nil, errors.Errorf("makroud: no such field %s in %T", name, value)
	}
	if field.Kind() == reflect.Ptr && field.IsNil() {
		return nil, nil
	}
	if !field.CanInterface() {
		return nil, errors.Errorf("makroud: cannot find field %s in %T", name, value)
	}

	return field.Interface(), nil
}

// GetFieldValueInt64 returns int64 value for the given instance field.
func GetFieldValueInt64(instance interface{}, field string) (int64, error) {
	value, err := GetFieldValueWithName(GetIndirectValue(instance), field)
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
	value, err := GetFieldValueWithName(GetIndirectValue(instance), field)
	if err != nil {
		return 0, false, err
	}

	converted, ok := ToOptionalInt64(value)
	return converted, ok, nil
}

// GetFieldValueString returns string value for the given instance field.
func GetFieldValueString(instance interface{}, field string) (string, error) {
	value, err := GetFieldValueWithName(GetIndirectValue(instance), field)
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
	value, err := GetFieldValueWithName(GetIndirectValue(instance), field)
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
		return pushFieldValueOnPointer(instance, dest, output, dtype, strict)

	case reflect.Slice:
		return pushFieldValueOnSlice(instance, dest, output, strict)

	case reflect.Struct:
		return pushFieldValueOnStruct(instance, dest, output)

	default:
		return pushFieldValueOnDefault(instance, dest, output)
	}
}

// pushFieldValueOnPointer tries to push given value on pointer instance.
func pushFieldValueOnPointer(instance interface{}, dest reflect.Value, output reflect.Value,
	dtype reflect.Type, strict bool) error {

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
}

// pushFieldValueOnSlice tries to push given value on slice instance.
func pushFieldValueOnSlice(instance interface{}, dest reflect.Value, output reflect.Value, strict bool) error {
	return updateFieldValueOnSlice(dest, output, instance, strict)
}

// pushFieldValueOnStruct tries to push given value on struct instance.
func pushFieldValueOnStruct(instance interface{}, dest reflect.Value, output reflect.Value) error {

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

// pushFieldValueOnDefault tries to push given value on default instance.
func pushFieldValueOnDefault(instance interface{}, dest reflect.Value, output reflect.Value) error {

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
		return dest, errors.Errorf("makroud: no such field %s in %T", name, instance)
	}
	if !dest.CanSet() {
		return dest, errors.Errorf("makroud: cannot update field %s in %T", name, instance)
	}

	return dest, nil
}

// getOutputReflectValue returns a reflect value used to update a field in given instance.
func getOutputReflectValue(instance interface{}, name string, value interface{}) (output reflect.Value, err error) {
	output = reflect.ValueOf(value)
	if !output.IsValid() && value != nil {
		return output, errors.Errorf("makroud: cannot uses %T as value to update %s in %T", value, name, instance)
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
		return errors.Errorf("makroud: cannot use type %v to update type %v in %T", sotype, sdtype, instance)
	}

	AppendReflectSlice(dest, output)
	return nil
}

func setFieldValueOnStruct(dest reflect.Value, output reflect.Value, instance interface{}) error {
	otype := output.Type()
	dtype := dest.Type()

	// Verify that types are equals. Otherwise, returns a nice error.
	if dtype != otype {
		return errors.Errorf("makroud: cannot use type %v to update type %v in %T", otype, dtype, instance)
	}

	dest.Set(output)

	return nil
}

// nolint: gocyclo
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
		zero := MakeReflectZero(scan.Type().Elem())
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
			arg = MakeReflectZero(scan.Type())
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
