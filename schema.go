package sqlxx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ulule/sqlxx/reflekt"
)

// ----------------------------------------------------------------------------
// Schema
// ----------------------------------------------------------------------------

// Schema is a model schema.
type Schema struct {
	ModelName    string
	TableName    string
	PrimaryField Field
	Fields       map[string]Field
	Relations    map[string]Relation
}

// NewSchema returns a new Schema instance.
func NewSchema(model Model) Schema {
	return Schema{
		ModelName: reflekt.ReflectType(model).Name(),
		TableName: model.TableName(),
		Fields:    map[string]Field{},
		Relations: map[string]Relation{},
	}
}

// SetPrimaryField sets the given Field as schema primary key.
func (s *Schema) SetPrimaryField(f Field) {
	f.IsPrimary = true
	s.PrimaryField = f
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
func (s Schema) WhereColumns(params map[string]interface{}) Conditions {
	return s.whereColumns(params, true)
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

// ----------------------------------------------------------------------------
// Schema API
// ----------------------------------------------------------------------------

// GetSchema returns model's table columns, extracted by reflection.
// The returned map is modelFieldName -> table_name.column_name
func GetSchema(model Model) (Schema, error) {
	var err error

	schema, found := cache.GetSchema(model)

	if found {
		return schema, nil
	}

	schema = NewSchema(model)

	v := reflekt.ReflectValue(model)

	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)

		meta := reflekt.GetFieldMeta(structField, SupportedTags, TagsMapping)

		if (meta.Type.Kind() == reflect.Struct) || (meta.Type.Kind() == reflect.Slice) {
			relationType := getRelationType(meta.Type)

			if _, ok := RelationTypes[relationType]; ok {
				schema.Relations[meta.Name], err = makeRelation(schema, model, meta, relationType)
				if err != nil {
					return Schema{}, err
				}

				continue
			}
		}

		field, err := makeField(model, meta)
		if err != nil {
			return Schema{}, err
		}

		if v := meta.Tags.GetByKey(StructTagName, StructTagPrimaryKey); len(v) != 0 {
			schema.SetPrimaryField(field)
			continue
		}

		schema.Fields[meta.Name] = field
	}

	cache.SetSchema(schema)

	return schema, nil
}

// GetSchemaFromInterface returns Schema by reflecting model for the given interface.
func GetSchemaFromInterface(out interface{}) (Schema, error) {
	return GetSchema(reflectModel(out))
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
// Columns
// ----------------------------------------------------------------------------

// Columns is a list of table columns.
type Columns []string

// Returns string representation of slice.
func (c Columns) String() string {
	return strings.Join(c, ", ")
}

// ----------------------------------------------------------------------------
// Conditions
// ----------------------------------------------------------------------------

// Conditions is a list of query conditions
type Conditions []string

// String returns conditions as AND query.
func (c Conditions) String() string {
	return strings.Join(c, " AND ")
}

// OR returns conditions as OR query.
func (c Conditions) OR() string {
	return strings.Join(c, " OR ")
}
