package reflekt

import "reflect"

// Meta are low level field metadata.
type FieldMeta struct {
	Name  string
	Field reflect.StructField
	Type  reflect.Type
	Tags  Tags
}

// GetFieldMeta returns field reflect data.
func GetFieldMeta(field reflect.StructField, tags []string, tagsMapping map[string]string) FieldMeta {
	var (
		fieldName = field.Name
		fieldType = field.Type
	)

	if field.Type.Kind() == reflect.Ptr {
		fieldType = field.Type.Elem()
	}

	return FieldMeta{
		Name:  fieldName,
		Field: field,
		Type:  fieldType,
		Tags:  GetFieldTags(field, tags, tagsMapping),
	}
}
