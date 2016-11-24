package sqlxx

import "reflect"

// Schema is a model schema.
type Schema struct {
	Fields    map[string]Field
	Relations map[string]Relation
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
		var (
			structField     reflect.StructField
			structFieldType reflect.Type
			fieldName       string
		)

		// valueField := deference(v.Field(i))
		structField = v.Type().Field(i)
		fieldName = structField.Name

		if structField.Type.Kind() == reflect.Ptr {
			structFieldType = structField.Type.Elem()
		} else {
			structFieldType = structField.Type
		}

		if (structFieldType.Kind() == reflect.Struct) || (structFieldType.Kind() == reflect.Slice) {
			relationType := getFieldRelationType(structFieldType)

			if _, ok := RelationTypes[relationType]; ok {
				schema.Relations[fieldName], err = newRelation(model, fieldName, relationType)
				if err != nil {
					return nil, err
				}
				continue
			}
		}

		schema.Fields[fieldName], err = newField(model, fieldName)
		if err != nil {
			return nil, err
		}
	}

	return schema, nil
}
