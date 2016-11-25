package sqlxx

import "reflect"

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
func getReflectedType(entity interface{}) reflect.Type {
	typ := reflect.ValueOf(entity).Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}

// getReflectedValue returns reflected value of the given entity.
func getReflectedValue(entity interface{}) reflect.Value {
	if reflect.TypeOf(entity).Kind() == reflect.Ptr {
		return reflect.ValueOf(entity).Elem()
	}

	return reflect.ValueOf(entity)
}

// isZeroValue returns true if the given interface is a zero value.
func isZeroValue(entity interface{}) bool {
	return !reflect.Indirect(reflect.ValueOf(entity)).IsValid()
}
