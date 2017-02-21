package sqlxx

import (
	"fmt"
	"reflect"

	"github.com/ulule/sqlxx/reflekt"
)

// Schema is a model schema.
type Schema struct {
	Model        Model
	ModelName    string
	TableName    string
	PrimaryField Field
	Fields       map[string]Field
	Relations    map[string]Relation
}

// SetPrimaryField sets the given Field as schema primary key.
func (s *Schema) SetPrimaryField(f Field) {
	f.IsPrimary = true
	s.PrimaryField = f
	s.Fields[f.Name] = f
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

// RelationPaths returns relations struct paths: Article.Author.Avatars
func (s Schema) RelationPaths() map[string]Relation {
	return GetSchemaRelations(s)
}

// GetSchemaRelations returns flattened map of schema relations.
func GetSchemaRelations(schema Schema) map[string]Relation {
	paths := map[string]Relation{}

	for _, relation := range schema.Relations {
		paths[relation.Name] = relation

		rels := GetSchemaRelations(relation.Schema)
		for _, rel := range rels {
			paths[fmt.Sprintf("%s.%s", relation.Name, rel.Name)] = rel
		}
	}

	return paths
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
	var (
		err error
		v   = reflekt.GetIndirectValue(model)
	)

	schema := Schema{
		Model:     model,
		ModelName: reflekt.GetIndirectType(model).Name(),
		TableName: model.TableName(),
		Fields:    map[string]Field{},
		Relations: map[string]Relation{},
	}

	for i := 0; i < v.NumField(); i++ {
		var (
			structField = v.Type().Field(i)
			meta        = GetFieldMeta(structField, SupportedTags, TagsMapping)
		)

		if IsExcludedField(meta) {
			continue
		}

		if (meta.Type.Kind() == reflect.Struct) || (meta.Type.Kind() == reflect.Slice) {
			relationType := getRelationType(meta.Type)

			if _, ok := RelationTypes[relationType]; ok {
				schema.Relations[meta.Name], err = NewRelation(schema, model, meta, relationType)
				if err != nil {
					return Schema{}, err
				}

				continue
			}
		}

		field, err := NewField(model, meta)
		if err != nil {
			return Schema{}, err
		}

		if IsPrimaryKeyField(meta) {
			schema.SetPrimaryField(field)
		}

		schema.Fields[meta.Name] = field
	}

	return schema, nil
}
