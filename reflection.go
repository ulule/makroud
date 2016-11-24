package sqlxx

import (
	"reflect"
	"strings"
)

// deferenceValue deferences the given value if it's a pointer or pointer to interface.
func deferenceValue(v reflect.Value) reflect.Value {
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

// getFieldTag returns field tag value.
func getFieldTag(structField reflect.StructField, name string) (map[string]string, error) {
	value := structField.Tag.Get(name)

	if len(value) == 0 {
		return nil, nil
	}

	results := map[string]string{}

	parts := strings.Split(value, " ")

	for _, part := range parts {
		splits := strings.Split(part, ":")
		results[splits[0]] = splits[1]
	}

	return results, nil
}
