package sqlxx

import (
	"fmt"
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

// Columns returns schema columns without table prefix.
func (s Schema) Columns() Columns {
	return s.columns(false)
}

// ColumnPaths returns schema column with table prefix.
func (s Schema) ColumnPaths() Columns {
	return s.columns(true)
}

// columns generates column slice.
func (s Schema) columns(withTable bool) Columns {
	columns := Columns{}

	for _, f := range s.Fields {
		if withTable {
			columns = append(columns, f.ColumnPath())
		} else {
			columns = append(columns, f.ColumnName)
		}
	}

	return columns
}

// WhereColumns returns where clause with the given params without table prefix.
func (s Schema) WhereColumns(params map[string]interface{}) Columns {
	return s.whereColumns(params, true)
}

// WhereColumnPaths returns where clause with the given params with table prefix.
func (s Schema) WhereColumnPaths(params map[string]interface{}) Columns {
	return s.whereColumns(params, true)
}

// whereColumns generates where clause for the given params.
func (s Schema) whereColumns(params map[string]interface{}, withTable bool) Columns {
	wheres := Columns{}

	for k := range params {
		column := k
		if withTable {
			column = fmt.Sprintf("%s.%s", s.TableName, k)
		}

		wheres = append(wheres, fmt.Sprintf("%s=:%s", column, k))
	}

	return wheres
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
