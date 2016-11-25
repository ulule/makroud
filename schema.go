package sqlxx

import "reflect"

// Schema is a model schema.
type Schema struct {
	PrimaryKey Field
	Fields     map[string]Field
	Relations  map[string]Relation
}

// GetSchema returns model's table columns, extracted by reflection.
// The returned map is modelFieldName -> table_name.column_name
func GetSchema(model Model) (*Schema, error) {
	var err error

	v := reflect.ValueOf(model)

	schema := &Schema{
		Fields:    map[string]Field{},
		Relations: map[string]Relation{},
	}

	v = deferenceValue(v)

	for i := 0; i < v.NumField(); i++ {
		valueField := deferenceValue(v.Field(i))
		structField := v.Type().Field(i)
		fieldMeta := makeFieldMeta(structField, valueField)

		if (fieldMeta.Type.Kind() == reflect.Struct) || (fieldMeta.Type.Kind() == reflect.Slice) {
			relationType := getFieldRelationType(fieldMeta.Type)

			if _, ok := RelationTypes[relationType]; ok {
				schema.Relations[fieldMeta.Name], err = newRelation(model, fieldMeta, relationType)
				if err != nil {
					return nil, err
				}

				continue
			}
		}

		schema.Fields[fieldMeta.Name], err = newField(model, fieldMeta)
		if err != nil {
			return nil, err
		}
	}

	return schema, nil
}
