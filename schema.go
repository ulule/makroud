package sqlxx

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/oleiade/reflections"
	"github.com/pkg/errors"
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
		names = append(names, f.FieldName)
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
func GetSchema(driver Driver, out interface{}) (Schema, error) {
	model := ToModel(out)

	if !driver.hasCache() {
		return newSchema(driver, model)
	}

	schema, found := driver.cache().GetSchema(model)
	if found {
		return schema, nil
	}

	schema, err := newSchema(driver, model)
	if err != nil {
		return schema, err
	}

	driver.cache().SetSchema(schema)

	return schema, nil
}

// newSchema returns model's table columns, extracted by reflection.
// The returned map is modelFieldName -> table_name.column_name
func newSchema(driver Driver, model Model) (Schema, error) {
	fmt.Printf("begin %T\n", model)
	defer fmt.Printf("end %T\n", model)
	schema := Schema{
		Model:        model,
		ModelName:    GetIndirectType(model).Name(),
		TableName:    model.TableName(),
		Fields:       map[string]Field{},
		Associations: map[string]Field{},
	}

	// TODO remove reflect here
	fields, err := reflections.Fields(model)
	if err != nil {
		return Schema{}, errors.Wrapf(err, "cannot use reflections to obtain %T fields", model)
	}

	for _, name := range fields {
		field, err := NewField(driver, &schema, model, name)
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
			schema.Fields[field.FieldName] = field
			continue
		}

		_, ok := schema.Associations[field.FieldName]
		if ok {
			continue
		}

		schema.Associations[field.FieldName] = field

		nextModel := field.ForeignKey.Reference.Model
		if field.IsAssociationTypeMany() {
			nextModel = field.ForeignKey.Model
		}

		nextSchema, err := GetSchema(driver, nextModel)
		if err != nil {
			return Schema{}, err
		}

		for k, v := range nextSchema.Associations {
			key := fmt.Sprintf("%s.%s", field.FieldName, k)
			_, ok := schema.Associations[key]
			if !ok {
				schema.Associations[key] = v
			}
		}
	}

	spew.Dump(schema)

	return schema, nil
}

// ----------------------------------------------------------------------------
// Columns
// ----------------------------------------------------------------------------

// Columns is a list of table columns.
type Columns []string

// Returns string representation of slice.
func (c Columns) String() string {
	sort.Strings(c)
	return strings.Join(c, ", ")
}

// GetColumns returns a comma-separated string representation of a model's table columns.
func GetColumns(driver Driver, model Model) (string, error) {
	schema, err := GetSchema(driver, model)
	if err != nil {
		return "", errors.Wrap(err, "sqlxx: cannot fetch schema informations")
	}

	columns := schema.ColumnPaths().String()
	return columns, nil
}

// ----------------------------------------------------------------------------
// Where clauses
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
