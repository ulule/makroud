package makroud

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/makroud/reflectx"
)

// Mapper will be used to mutate a Model with row values.
type Mapper map[string]interface{}

// Schema is a model schema.
type Schema struct {
	model        Model
	modelName    string
	tableName    string
	pk           PrimaryKey
	fields       map[string]Field
	references   map[string]ForeignKey
	associations map[string]Reference
	createdKey   *Field
	updatedKey   *Field
	deletedKey   *Field
}

// Model returns the schema model.
func (schema Schema) Model() Model {
	return schema.model
}

// ModelName returns the schema model name.
func (schema Schema) ModelName() string {
	return schema.modelName
}

// TableName returns the schema table name.
func (schema Schema) TableName() string {
	return schema.tableName
}

// PrimaryKey returns the schema primary key.
func (schema Schema) PrimaryKey() PrimaryKey {
	return schema.pk
}

// PrimaryKeyPath returns schema primary key column path.
func (schema Schema) PrimaryKeyPath() string {
	return schema.pk.ColumnPath()
}

// PrimaryKeyName returns schema primary key column name.
func (schema Schema) PrimaryKeyName() string {
	return schema.pk.ColumnName()
}

// HasCreatedKey returns if a created key is defined for current schema.
func (schema Schema) HasCreatedKey() bool {
	return schema.createdKey != nil
}

// CreatedKeyPath returns schema created key column path.
func (schema Schema) CreatedKeyPath() string {
	if schema.HasUpdatedKey() {
		return schema.createdKey.ColumnPath()
	}
	panic(fmt.Sprint("makroud: ", ErrSchemaCreatedKey))
}

// CreatedKeyName returns schema created key column name.
func (schema Schema) CreatedKeyName() string {
	if schema.HasUpdatedKey() {
		return schema.createdKey.ColumnName()
	}
	panic(fmt.Sprint("makroud: ", ErrSchemaCreatedKey))
}

// HasUpdatedKey returns if an updated key is defined for current schema.
func (schema Schema) HasUpdatedKey() bool {
	return schema.updatedKey != nil
}

// UpdatedKeyPath returns schema updated key column path.
func (schema Schema) UpdatedKeyPath() string {
	if schema.HasUpdatedKey() {
		return schema.updatedKey.ColumnPath()
	}
	panic(fmt.Sprint("makroud: ", ErrSchemaUpdatedKey))
}

// UpdatedKeyName returns schema deleted key column name.
func (schema Schema) UpdatedKeyName() string {
	if schema.HasUpdatedKey() {
		return schema.updatedKey.ColumnName()
	}
	panic(fmt.Sprint("makroud: ", ErrSchemaUpdatedKey))
}

// HasDeletedKey returns if a deleted key is defined for current schema.
func (schema Schema) HasDeletedKey() bool {
	return schema.deletedKey != nil
}

// DeletedKeyPath returns schema deleted key column path.
func (schema Schema) DeletedKeyPath() string {
	if schema.HasDeletedKey() {
		return schema.deletedKey.ColumnPath()
	}
	panic(fmt.Sprint("makroud: ", ErrSchemaDeletedKey))
}

// DeletedKeyName returns schema deleted key column name.
func (schema Schema) DeletedKeyName() string {
	if schema.HasDeletedKey() {
		return schema.deletedKey.ColumnName()
	}
	panic(fmt.Sprint("makroud: ", ErrSchemaDeletedKey))
}

// Columns returns schema columns without table prefix.
func (schema Schema) Columns() Columns {
	return schema.columns(false)
}

// ColumnPaths returns schema column with table prefix.
func (schema Schema) ColumnPaths() Columns {
	return schema.columns(true)
}

// columns generates column slice.
func (schema Schema) columns(withTable bool) Columns {
	columns := Columns{}
	if withTable {
		columns = append(columns, schema.pk.ColumnPath())
	} else {
		columns = append(columns, schema.pk.ColumnName())
	}
	for _, field := range schema.fields {
		if withTable {
			columns = append(columns, field.ColumnPath())
		} else {
			columns = append(columns, field.ColumnName())
		}
	}
	return columns
}

// HasColumn returns if a schema has a column or not.
func (schema Schema) HasColumn(column string) bool {
	if schema.pk.ColumnName() == column || schema.pk.ColumnPath() == column {
		return true
	}

	_, ok := schema.fields[column]
	if ok {
		return true
	}

	column = strings.TrimPrefix(column, fmt.Sprint(schema.TableName(), "."))
	_, ok = schema.fields[column]

	return ok
}

// nolint: gocyclo
func (schema Schema) getValues(value reflect.Value, columns []string, model Model) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	rest := make([]string, 0)
	associationsColumns := map[string]map[string]int{}

	for i, column := range columns {
		if schema.pk.ColumnName() == column || schema.pk.ColumnPath() == column {
			values[i] = reflectx.GetReflectFieldByIndexes(value, schema.pk.FieldIndex())
			continue
		}

		field, ok := schema.fields[column]
		if ok {
			values[i] = reflectx.GetReflectFieldByIndexes(value, field.FieldIndex())
			continue
		}

		column = strings.TrimPrefix(column, fmt.Sprint(schema.TableName(), "."))
		field, ok = schema.fields[column]
		if ok {
			values[i] = reflectx.GetReflectFieldByIndexes(value, field.FieldIndex())
			continue
		}

		// sorting associations columns in case of JOIN
		found := false
		for key, association := range schema.associations {
			trimed := strings.TrimPrefix(column, fmt.Sprint(association.Remote().TableName(), "."))
			if trimed != column {
				_, ok := associationsColumns[key]
				if !ok {
					associationsColumns[key] = map[string]int{}
				}

				// we keep the index in values for the scan
				associationsColumns[key][trimed] = i
				found = true
				break
			}
		}

		if found {
			continue
		}

		rest = append(rest, column)
	}

	if len(rest) > 0 {
		return nil, errors.Wrapf(ErrSchemaColumnRequired,
			"missing destination name %s in %T", strings.Join(rest, ", "), model)
	}

	for key, columns := range associationsColumns {
		// retrieve reflect field based on index
		model := reflectx.GetReflectFieldByIndexes(value, schema.associations[key].Field.FieldIndex())
		remote := schema.associations[key].Remote()

		associationValue := reflectx.GetIndirectValue(model)

		// columns concerned by the association
		rest := make([]string, 0, len(columns))
		for key := range columns {
			rest = append(rest, key)
		}

		associationValues, err := remote.Schema().getValues(associationValue, rest, remote.Model())
		if err != nil {
			return nil, err
		}

		for i := range associationValues {
			// index previous stored
			index := associationsColumns[key][rest[i]]

			values[index] = associationValues[i]
		}
	}

	return values, nil
}

// ScanRow executes a scan from given row into model.
func (schema Schema) ScanRow(row Row, model Model) error {
	columns, err := row.Columns()
	if err != nil {
		return err
	}

	value := reflectx.GetIndirectValue(model)
	if !reflectx.IsStruct(value) {
		return errors.Wrapf(ErrStructRequired, "cannot use mapper on %T", model)
	}

	values, err := schema.getValues(value, columns, model)
	if err != nil {
		return err
	}

	return row.Scan(values...)
}

// ScanRows executes a scan from current row into model.
func (schema Schema) ScanRows(rows Rows, model Model) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	value := reflectx.GetIndirectValue(model)
	if !reflectx.IsStruct(value) {
		return errors.Wrapf(ErrStructRequired, "cannot use mapper on %T", model)
	}

	values, err := schema.getValues(value, columns, model)
	if err != nil {
		return err
	}

	return rows.Scan(values...)
}

// ----------------------------------------------------------------------------
// Initializers
// ----------------------------------------------------------------------------

// GetSchema returns the schema from given model.
// If the schema does not exists, it returns an error.
func GetSchema(driver Driver, model Model) (*Schema, error) {
	return getSchema(driver, model, true)
}

// getSchema returns the schema from given model.
// If the schema does not exists, it returns an error.
// If throughout is true, it will execute a full scan of given model:
// this is a trick to allow circular import of model.
func getSchema(driver Driver, model Model, throughout bool) (*Schema, error) {
	if !driver.hasCache() {
		return newSchema(driver, model, throughout)
	}

	schema := driver.getCache().GetSchema(model)
	if schema != nil {
		return schema, nil
	}

	schema, err := newSchema(driver, model, throughout)
	if err != nil {
		return nil, err
	}

	if throughout {
		driver.getCache().SetSchema(schema)
	}

	return schema, nil
}

// defaultModelOpts returns the default model configuration.
func defaultModelOpts() ModelOpts {
	return ModelOpts{
		PrimaryKey: "id",
		CreatedKey: "created_at",
		UpdatedKey: "updated_at",
		DeletedKey: "deleted_at",
	}
}

// analyzeModelOpts analyzes given model to extract it's configuration.
func analyzeModelOpts(model Model) ModelOpts {
	opts := defaultModelOpts()

	mpk, ok := model.(interface {
		PrimaryKey() string
	})
	if ok {
		opts.PrimaryKey = mpk.PrimaryKey()
	}

	cpk, ok := model.(interface {
		CreatedKey() string
	})
	if ok {
		opts.CreatedKey = cpk.CreatedKey()
	}

	upk, ok := model.(interface {
		UpdatedKey() string
	})
	if ok {
		opts.UpdatedKey = upk.UpdatedKey()
	}

	dpk, ok := model.(interface {
		DeletedKey() string
	})
	if ok {
		opts.DeletedKey = dpk.DeletedKey()
	}

	return opts
}

// newSchema returns a schema from given model, extracted by reflection.
// The returned schema is a mapping of a model to table and columns.
// For example: Model.FieldName -> table_name.column_name
//
// If throughout is true, it will execute a full and complete scan of given model:
// this is a trick to allow circular import of model.
func newSchema(driver Driver, model Model, throughout bool) (*Schema, error) {
	fields, err := reflectx.GetFields(model)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot use reflections to obtain %T fields", model)
	}

	modelOpts := analyzeModelOpts(model)

	schema := &Schema{
		model:        reflectx.MakeZero(model).(Model),
		modelName:    reflectx.GetIndirectTypeName(model),
		tableName:    model.TableName(),
		fields:       map[string]Field{},
		references:   map[string]ForeignKey{},
		associations: map[string]Reference{},
	}

	relationships := map[string]*Field{}

	err = getSchemaFields(driver, schema, model, modelOpts, fields, relationships)
	if err != nil {
		return nil, err
	}

	if throughout {
		err = getSchemaAssociations(driver, schema, model, relationships)
		if err != nil {
			return nil, err
		}
	}

	err = inferSchemaPrimaryKey(model, modelOpts, schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func getSchemaFields(driver Driver, schema *Schema, model Model, modelOpts ModelOpts,
	fields []string, relationships map[string]*Field) error {

	for _, name := range fields {
		field, err := NewField(driver, schema, model, name, modelOpts)
		if err != nil {
			return err
		}

		if field.IsExcluded() {
			continue
		}

		err = inferSchemaTimeKey(model, modelOpts, schema, field)
		if err != nil {
			return err
		}

		if field.IsPrimaryKey() {
			err = handleSchemaPrimaryKey(schema, model, name, field)
			if err != nil {
				return err
			}
			continue
		}

		if field.IsForeignKey() {
			err = handleSchemaForeignKey(schema, model, name, field)
			if err != nil {
				return err
			}
		}

		if !field.IsAssociation() {
			schema.fields[field.ColumnName()] = *field
			continue
		}

		relationships[name] = field
	}

	return nil
}

func handleSchemaPrimaryKey(schema *Schema, model Model, name string, field *Field) error {
	pk, err := NewPrimaryKey(field)
	if err != nil {
		return errors.Wrapf(err, "cannot use '%s' as primary key for %T", name, model)
	}
	if schema.pk.TableName() != "" {
		return errors.Errorf("%T must have only one primary key", model)
	}
	schema.pk = *pk
	return nil
}

func handleSchemaForeignKey(schema *Schema, model Model, name string, field *Field) error {
	fk, err := NewForeignKey(field)
	if err != nil {
		return errors.Wrapf(err, "cannot use '%s' as foreign key for %T", name, model)
	}
	schema.references[fk.ColumnName()] = *fk
	return nil
}

func getSchemaAssociations(driver Driver, schema *Schema, model Model, relationships map[string]*Field) error {
	for name, field := range relationships {
		_, ok := schema.associations[field.FieldName()]
		if ok {
			continue
		}

		reference, err := NewReference(driver, schema, field)
		if err != nil {
			return errors.Wrapf(err, "cannot use '%s' as association for %T", name, model)
		}

		schema.associations[field.FieldName()] = *reference
	}
	return nil
}

func inferSchemaTimeKey(model Model, opts ModelOpts, schema *Schema, field *Field) error {
	if field.IsCreatedKey() {
		if schema.createdKey != nil {
			return errors.Errorf("%T must have only one created_at key", model)
		}
		schema.createdKey = field
	}

	if field.IsUpdatedKey() {
		if schema.updatedKey != nil {
			return errors.Errorf("%T must have only one updated_at key", model)
		}
		schema.updatedKey = field
	}

	if field.IsDeletedKey() {
		if schema.deletedKey != nil {
			return errors.Errorf("%T must have only one deleted_at key", model)
		}
		schema.deletedKey = field
	}

	return nil
}

func inferSchemaPrimaryKey(model Model, opts ModelOpts, schema *Schema) error {
	if schema.pk.TableName() != "" {
		return nil
	}
	for key, field := range schema.fields {
		if field.ColumnName() == opts.PrimaryKey {
			pk, err := NewPrimaryKey(&field)
			if err != nil {
				return errors.Wrapf(err, "cannot use primary key of %T", model)
			}
			schema.pk = *pk
			delete(schema.fields, key)
			return nil
		}
	}
	return errors.Errorf("%T must have a primary key", model)
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

// List returns table columns.
func (c Columns) List() []string {
	return c
}

// GetColumns returns a comma-separated string representation of a model's table columns.
func GetColumns(driver Driver, model Model) (Columns, error) {
	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, errors.Wrap(err, "makroud: cannot fetch schema informations")
	}

	columns := schema.ColumnPaths()
	return columns, nil
}
