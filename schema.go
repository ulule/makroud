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

	schema := &Schema{
		Fields:    map[string]Field{},
		Relations: map[string]Relation{},
	}

	v := reflectValue(reflect.ValueOf(model))

	for i := 0; i < v.NumField(); i++ {
		valueField := v.Field(i)
		structField := v.Type().Field(i)

		meta := makeMeta(structField, valueField)

		if (meta.Type.Kind() == reflect.Struct) || (meta.Type.Kind() == reflect.Slice) {
			relationType := getRelationType(meta.Type)

			if _, ok := RelationTypes[relationType]; ok {
				schema.Relations[meta.Name], err = newRelation(model, meta, relationType)
				if err != nil {
					return nil, err
				}

				continue
			}
		}

		field, err := newField(model, meta)
		if err != nil {
			return nil, err
		}

		if v := meta.Tags.GetByKey(StructTagName, "primary_key"); len(v) != 0 {
			schema.PrimaryKey = field
			field.IsPrimary = true
		}

		schema.Fields[meta.Name] = field
	}

	return schema, nil
}
