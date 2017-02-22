package sqlxx

import (
	"fmt"
	"reflect"

	"github.com/ulule/sqlxx/reflekt"
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
			columns = append(columns, f.ColumnPath)
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

// AssociationsByPath returns relations struct paths: Article.Author.Avatars
func (s Schema) AssociationsByPath() (map[string]Field, error) {
	return GetSchemaAssociations(s)
}

// GetSchemaAssociations returns flattened map of schema associations.
func GetSchemaAssociations(schema Schema) (map[string]Field, error) {
	var (
		err   error
		paths = map[string]Field{}
	)

	for _, f := range schema.Associations {
		if _, ok := paths[f.Name]; !ok {
			paths[f.Name] = f
		}

		schema, err = GetSchema(f.Association.Model)
		if err != nil {
			return nil, err
		}

		assocs, err := GetSchemaAssociations(schema)
		if err != nil {
			return nil, err
		}

		for _, assoc := range assocs {
			key := fmt.Sprintf("%s.%s", f.Name, assoc.Name)
			if _, ok := paths[key]; !ok {
				paths[key] = assoc
			}
		}
	}

	return paths, nil
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
	v := reflekt.GetIndirectValue(model)

	schema := Schema{
		Model:        model,
		ModelName:    reflekt.GetIndirectType(model).Name(),
		TableName:    model.TableName(),
		Fields:       map[string]Field{},
		Associations: map[string]Field{},
	}

	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)

		field, err := NewField(structField, model)
		if err != nil {
			return Schema{}, err
		}

		if field.IsExcluded {
			continue
		}

		if field.IsPrimaryKey {
			schema.PrimaryKeyField = field
		}

		if field.IsAssociation {
			schema.Associations[field.Name] = field
		} else {
			schema.Fields[field.Name] = field
		}
	}

	return schema, nil
}
