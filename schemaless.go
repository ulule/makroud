package makroud

import (
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/makroud/reflectx"
	"github.com/ulule/makroud/snaker"
)

// Schemaless is a light version of schema for structs that are not models.
// It's a simple mapper from column name to struct field, without primary key, associations and so on...
type Schemaless struct {
	rtype reflect.Type
	keys  map[string]SchemalessKey
}

// Type returns the reflect's type of the schema.
func (schema Schemaless) Type() reflect.Type {
	return schema.rtype
}

// Columns returns schema columns.
func (schema Schemaless) Columns() Columns {
	columns := Columns{}
	for _, key := range schema.keys {
		columns = append(columns, key.ColumnName())
	}
	return columns
}

// HasColumn returns if a schema has a column or not.
func (schema Schemaless) HasColumn(column string) bool {
	_, ok := schema.keys[column]
	return ok
}

// Key returns the SchemalessKey instance for given column.
func (schema Schemaless) Key(column string) (SchemalessKey, bool) {
	key, ok := schema.keys[column]
	return key, ok
}

// ScanRows executes a scan from current row into schemaless instance.
func (schema Schemaless) ScanRows(rows Rows, val interface{}) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	value := reflectx.GetIndirectValue(val)
	if !reflectx.IsStruct(value) {
		return errors.Wrapf(ErrStructRequired, "cannot use mapper on %T", val)
	}

	values, err := schema.getValues(value, columns, val)
	if err != nil {
		return err
	}

	return rows.Scan(values...)
}

func (schema Schemaless) getValues(value reflect.Value, columns []string, val interface{}) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	missing := make([]string, 0)

	for i, column := range columns {
		key, ok := schema.keys[column]
		if !ok {
			missing = append(missing, column)
			continue
		}

		values[i] = reflectx.GetReflectFieldByIndexes(value, key.FieldIndex())
	}

	if len(missing) > 0 {
		return nil, errors.Wrapf(ErrSchemaColumnRequired,
			"missing destination name %s in %T", strings.Join(missing, ", "), val)
	}

	return values, nil
}

// SchemalessKey is a light version of schema field.
// It contains the column name and the struct field information.
type SchemalessKey struct {
	columnName string
	fieldName  string
	fieldIndex []int
}

// ColumnName returns the column name for this schemaless key.
func (key SchemalessKey) ColumnName() string {
	return key.columnName
}

// FieldName define the struct field name used for this schemaless key.
func (key SchemalessKey) FieldName() string {
	return key.fieldName
}

// FieldIndex define the struct field index used for this schemaless key.
func (key SchemalessKey) FieldIndex() []int {
	return key.fieldIndex
}

// ----------------------------------------------------------------------------
// Initializers
// ----------------------------------------------------------------------------

// GetSchemaless returns the schema information from given type that are not models.
// If no information could be gathered, it returns an error.
func GetSchemaless(driver Driver, value reflect.Type) (*Schemaless, error) {
	if !driver.hasCache() {
		return newSchemaless(driver, value)
	}

	schema := driver.getCache().GetSchemaless(value)
	if schema != nil {
		return schema, nil
	}

	schema, err := newSchemaless(driver, value)
	if err != nil {
		return nil, err
	}

	driver.getCache().SetSchemaless(schema)

	return schema, nil
}

// newSchemaless returns a schemaless from given type, extracted by reflection.
// The returned schemaless is a mapping of a struct to columns.
// For example: Type.FieldName -> column_name
//
// If you need a mapping with a database table, please use a Schema instead of a Schemaless instance.
// You'll have better features such as primary key, foreign key, associations and so on...
func newSchemaless(driver Driver, value reflect.Type) (*Schemaless, error) {
	fields, err := reflectx.GetFields(value)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot use reflections to obtain %s fields", value.String())
	}

	schema := &Schemaless{
		rtype: value,
		keys:  map[string]SchemalessKey{},
	}

	for _, name := range fields {
		field, ok := reflectx.GetFieldByName(value, name)
		if !ok {
			return nil, errors.Errorf("field '%s' not found in model", name)
		}

		tags := GetTags(field, NewOnlyColumnTagsAnalyzerOption())

		isExcluded := tags.HasKey(TagName, TagKeyIgnored) || field.PkgPath != ""
		if isExcluded {
			continue
		}

		columnName := tags.GetByKey(TagName, TagKeyColumn)
		if columnName == "" {
			columnName = snaker.CamelToSnake(name)
		}

		key := SchemalessKey{
			columnName: columnName,
			fieldName:  field.Name,
			fieldIndex: field.Index,
		}

		schema.keys[key.ColumnName()] = key
	}

	return schema, nil
}
