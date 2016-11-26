package sqlxx

import "reflect"

// reflectValue returns the value that the interface v contains
// or that the pointer v points to.
func reflectValue(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Interface && !v.IsNil() {
		if v.Elem().Kind() == reflect.Ptr && !v.Elem().IsNil() && v.Elem().Elem().Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	return v
}

// reflectType returns type of the given interface.
func reflectType(itf interface{}) reflect.Type {
	typ := reflect.ValueOf(itf).Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

// reflectModel returns interface as a Model interface.
func reflectModel(itf interface{}) Model {
	v := reflect.Indirect(reflect.ValueOf(itf))

	if v.IsValid() && v.Kind() != reflect.Slice {
		return reflect.ValueOf(itf).Interface().(Model)
	}

	var typ reflect.Type

	if reflect.Indirect(v).Kind() == reflect.Slice {
		typ = v.Type().Elem()
	} else {
		typ = reflect.Indirect(v).Type()
	}

	return reflect.New(typ).Interface().(Model)
}

// getFieldType returns the field type for the given value.
func getFieldRelationType(typ reflect.Type) RelationType {
	if typ.Kind() == reflect.Slice {
		if _, isModel := reflect.New(typ.Elem()).Interface().(Model); isModel {
			return RelationTypeManyToOne
		}
		return RelationTypeUnknown
	}

	if _, isModel := reflect.New(typ).Interface().(Model); isModel {
		return RelationTypeOneToMany
	}

	return RelationTypeUnknown
}

func getFieldTags(structField reflect.StructField, names ...string) map[string]string {
	tags := map[string]string{}

	for _, name := range names {
		if _, ok := tags[name]; !ok {
			tags[name] = structField.Tag.Get(name)
		}
	}

	return tags
}

// getType returns type.
func getReflectedType(itf interface{}) reflect.Type {
	typ := reflect.ValueOf(itf).Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

// getReflectedValue returns reflected value of the given itf.
func getReflectedValue(itf interface{}) reflect.Value {
	if reflect.TypeOf(itf).Kind() == reflect.Ptr {
		return reflect.ValueOf(itf).Elem()
	}

	return reflect.ValueOf(itf)
}

// isZeroValue returns true if the given interface is a zero value.
func isZeroValue(itf interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(itf))
	return v.Interface() == reflect.Zero(v.Type()).Interface()
}

// getStructFields returns struct fields of value.
func getStructFields(v reflect.Value) []reflect.StructField {
	if v.Kind() != reflect.Struct {
		return nil
	}

	fields := []reflect.StructField{}

	for i := 0; i < v.NumField(); i++ {
		fields = append(fields, v.Type().Field(i))
	}

	return fields
}
