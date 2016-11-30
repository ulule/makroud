package sqlxx

import "reflect"

// Meta are low level field metadata.
type Meta struct {
	Name  string
	Field reflect.StructField
	Type  reflect.Type
	Tags  Tags
}

// makeMeta returns field reflect data.
func makeMeta(field reflect.StructField) Meta {
	var (
		fieldName = field.Name
		fieldType = field.Type
	)

	if field.Type.Kind() == reflect.Ptr {
		fieldType = field.Type.Elem()
	}

	return Meta{
		Name:  fieldName,
		Field: field,
		Type:  fieldType,
		Tags:  makeTags(field),
	}
}
