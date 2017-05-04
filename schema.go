package sqlxx

import (
	"fmt"
	"reflect"

	"github.com/oleiade/reflections"
)

// Schema is a model schema.
type Schema struct {
	Model           Model
	ModelName       string
	TableName       string
	PrimaryKeyField Field
	Fields          map[string]Field
	Associations    map[string]Field
}

// FieldNames return all field names.
func (s Schema) FieldNames() []string {
	var names []string
	for _, f := range s.Fields {
		names = append(names, f.Name)
	}
	return names
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
func (s Schema) WhereColumns(params map[string]interface{}) Conditions {
	return s.whereColumns(params, false)
}

// WhereColumnPaths returns where clause with the given params with table prefix.
func (s Schema) WhereColumnPaths(params map[string]interface{}) Conditions {
	return s.whereColumns(params, true)
}

// whereColumns generates where clause for the given params.
func (s Schema) whereColumns(params map[string]interface{}, withTable bool) Conditions {
	wheres := Conditions{}
	for k, v := range params {
		column := k
		if withTable {
			column = fmt.Sprintf("%s.%s", s.TableName, k)
		}
		if reflect.Indirect(reflect.ValueOf(v)).Kind() == reflect.Slice {
			wheres = append(wheres, fmt.Sprintf("%s IN (:%s)", column, k))
		} else {
			wheres = append(wheres, fmt.Sprintf("%s = :%s", column, k))
		}
	}
	return wheres
}

// ----------------------------------------------------------------------------
// Initializers
// ----------------------------------------------------------------------------

// GetSchema returns the given schema from global cache
// If the given schema does not exists, returns false as bool.
func GetSchema(itf interface{}) (Schema, error) {
	var (
		err    error
		schema Schema
		model  = GetModelFromInterface(itf)
	)

	if cacheDisabled {
		return newSchema(model)
	}

	schema, found := cache.GetSchema(model)
	if found {
		return schema, nil
	}

	schema, err = newSchema(model)
	if err != nil {
		return schema, err
	}

	cache.SetSchema(schema)

	return schema, nil
}

// newSchema returns model's table columns, extracted by reflection.
// The returned map is modelFieldName -> table_name.column_name
func newSchema(model Model) (Schema, error) {
	schema := Schema{
		Model:        model,
		ModelName:    GetIndirectType(model).Name(),
		TableName:    model.TableName(),
		Fields:       map[string]Field{},
		Associations: map[string]Field{},
	}

	fields, err := reflections.Fields(model)
	if err != nil {
		return Schema{}, err
	}

	for _, name := range fields {
		field, err := NewField(&schema, model, name)
		if err != nil {
			return Schema{}, err
		}

		if field.IsExcluded {
			continue
		}

		if field.IsPrimaryKey {
			schema.PrimaryKeyField = field
		}

		if !field.IsAssociation {
			schema.Fields[field.Name] = field
			continue
		}

		if _, ok := schema.Associations[field.Name]; ok {
			continue
		}

		schema.Associations[field.Name] = field

		nextSchema, err := GetSchema(field.RelationModel())
		if err != nil {
			return Schema{}, err
		}

		for k, v := range nextSchema.Associations {
			key := fmt.Sprintf("%s.%s", field.Name, k)
			if _, ok := schema.Associations[key]; !ok {
				schema.Associations[key] = v
			}
		}
	}

	return schema, nil
}
