package sqlxx

import (
	"reflect"
	"strings"
)

// Columns is a list of table columns.
type Columns []string

// Returns string representation of slice.
func (c Columns) String() string {
	return strings.Join(c, ", ")
}

// Schema is a model schema.
type Schema struct {
	TableName  string
	PrimaryKey Field
	Fields     map[string]Field
	Relations  map[string]Relation
}

// NewSchema returns a new Schema instance.
func NewSchema(model Model) *Schema {
	return &Schema{
		TableName: model.TableName(),
		Fields:    map[string]Field{},
		Relations: map[string]Relation{},
	}
}

// SetPrimaryKey sets the given Field as schema primary key.
func (s *Schema) SetPrimaryKey(f Field) {
	f.IsPrimary = true
	s.PrimaryKey = f
	s.Fields[f.Name] = f
}

// Columns returns schema columns.
func (s Schema) Columns() Columns {
	columns := Columns{}
	for _, f := range s.Fields {
		columns = append(columns, f.ColumnName)
	}
	return columns
}

// ColumnPaths returns schema column paths.
func (s Schema) ColumnPaths() Columns {
	columns := Columns{}
	for _, f := range s.Fields {
		columns = append(columns, f.ColumnPath())
	}
	return columns
}

// GetSchema returns model's table columns, extracted by reflection.
// The returned map is modelFieldName -> table_name.column_name
func GetSchema(model Model) (*Schema, error) {
	var err error

	schema := NewSchema(model)

	v := reflectValue(reflect.ValueOf(model))

	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		meta := makeMeta(structField)

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

		if v := meta.Tags.GetByKey(StructTagName, StructTagPrimaryKey); len(v) != 0 {
			schema.SetPrimaryKey(field)
			continue
		}

		schema.Fields[meta.Name] = field
	}

	return schema, nil
}

// getSchemaFromInterface returns Schema by reflecting model for the given interface.
func getSchemaFromInterface(out interface{}) (*Schema, error) {
	return GetSchema(reflectModel(out))
}
