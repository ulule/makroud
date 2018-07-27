package sqlxx

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/sqlxx/reflectx"
)

// ErrCreatedKey is returned when we cannot find a created key in given schema.
var ErrCreatedKey = fmt.Errorf("cannot find created key in schema")

// ErrUpdatedKey is returned when we cannot find a updated key in given schema.
var ErrUpdatedKey = fmt.Errorf("cannot find updated key in schema")

// ErrDeletedKey is returned when we cannot find a deleted key in given schema.
var ErrDeletedKey = fmt.Errorf("cannot find deleted key in schema")

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

// TableName returns the schema table name.
func (schema Schema) TableName() string {
	return schema.tableName
}

// PrimaryKey returns the schema primary key.
func (schema Schema) PrimaryKey() PrimaryKey {
	return schema.pk
}

// HasCreatedKey returns if an created key is defined for current schema.
func (schema Schema) HasCreatedKey() bool {
	return schema.createdKey != nil
}

// CreatedKeyPath returns schema created key column's name.
func (schema Schema) CreatedKeyPath() string {
	if schema.HasUpdatedKey() {
		return schema.createdKey.ColumnPath()
	}
	panic(fmt.Sprint("sqlxx: ", ErrCreatedKey))
}

// CreatedKeyName returns schema created key column's name.
func (schema Schema) CreatedKeyName() string {
	if schema.HasUpdatedKey() {
		return schema.createdKey.ColumnName()
	}
	panic(fmt.Sprint("sqlxx: ", ErrCreatedKey))
}

// HasUpdatedKey returns if an updated key is defined for current schema.
func (schema Schema) HasUpdatedKey() bool {
	return schema.updatedKey != nil
}

// UpdatedKeyPath returns schema updated key column's name.
func (schema Schema) UpdatedKeyPath() string {
	if schema.HasUpdatedKey() {
		return schema.updatedKey.ColumnPath()
	}
	panic(fmt.Sprint("sqlxx: ", ErrUpdatedKey))
}

// UpdatedKeyName returns schema deleted key column's name.
func (schema Schema) UpdatedKeyName() string {
	if schema.HasUpdatedKey() {
		return schema.updatedKey.ColumnName()
	}
	panic(fmt.Sprint("sqlxx: ", ErrUpdatedKey))
}

// HasDeletedKey returns if an deleted key is defined for current schema.
func (schema Schema) HasDeletedKey() bool {
	return schema.deletedKey != nil
}

// DeletedKeyPath returns schema deleted key column's name.
func (schema Schema) DeletedKeyPath() string {
	if schema.HasDeletedKey() {
		return schema.deletedKey.ColumnPath()
	}
	panic(fmt.Sprint("sqlxx: ", ErrDeletedKey))
}

// DeletedKeyName returns schema deleted key column's name.
func (schema Schema) DeletedKeyName() string {
	if schema.HasDeletedKey() {
		return schema.deletedKey.ColumnName()
	}
	panic(fmt.Sprint("sqlxx: ", ErrDeletedKey))
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

// WriteModel will try to updates given model from sqlx mapper.
func (schema Schema) WriteModel(mapper Mapper, model Model) error {
	if len(mapper) == 0 {
		return nil
	}
	for key, value := range mapper {
		if schema.pk.ColumnName() == key || schema.pk.ColumnPath() == key {
			err := reflectx.UpdateFieldValue(model, schema.pk.Field.FieldName(), value)
			if err != nil {
				return err
			}
			continue
		}

		field, ok := schema.fields[key]
		if ok {
			err := reflectx.UpdateFieldValue(model, field.FieldName(), value)
			if err != nil {
				return err
			}
			continue
		}

		key = strings.TrimPrefix(key, fmt.Sprint(schema.TableName(), "."))
		field, ok = schema.fields[key]
		if ok {
			err := reflectx.UpdateFieldValue(model, field.FieldName(), value)
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

func (schema Schema) writeField(model Model, field Field, value interface{}) error {
	return reflectx.UpdateFieldValue(model, field.FieldName(), value)
}

// ----------------------------------------------------------------------------
// Initializers
// ----------------------------------------------------------------------------

// GetSchema returns the given schema from global cache
// If the given schema does not exists, returns false as bool.
func GetSchema(driver Driver, model Model) (*Schema, error) {
	if !driver.hasCache() {
		return newSchema(driver, model)
	}

	schema := driver.cache().GetSchema(model)
	if schema != nil {
		return schema, nil
	}

	schema, err := newSchema(driver, model)
	if err != nil {
		return nil, err
	}

	driver.cache().SetSchema(schema)

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

// newSchema returns model's table columns, extracted by reflection.
// The returned map is Model.FieldName -> table_name.column_name
func newSchema(driver Driver, model Model) (*Schema, error) {
	fields, err := reflectx.GetFields(model)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot use reflections to obtain %T fields", model)
	}

	modelOpts := analyzeModelOpts(model)

	schema := &Schema{
		// Model:     model,
		modelName:    reflectx.GetIndirectTypeName(model),
		tableName:    model.TableName(),
		fields:       map[string]Field{},
		references:   map[string]ForeignKey{},
		associations: map[string]Reference{},
	}

	relationships := map[string]*Field{}

	for _, name := range fields {
		field, err := NewField(driver, schema, model, name, modelOpts)
		if err != nil {
			return nil, err
		}

		if field.IsExcluded() {
			continue
		}

		err = inferSchemaTimeKey(model, modelOpts, schema, field)
		if err != nil {
			return nil, err
		}

		if field.IsPrimaryKey() {
			pk, err := NewPrimaryKey(field)
			if err != nil {
				return nil, errors.Wrapf(err, "cannot use primary key of %T", model)
			}
			if schema.pk.TableName() != "" {
				return nil, errors.Errorf("sqlxx: %T must have only one primary key", model)
			}
			schema.pk = *pk
			continue
		}

		if field.IsForeignKey() {
			fk, err := NewForeignKey(field)
			if err != nil {
				return nil, errors.Wrapf(err, "cannot use '%s' as foreign key for %T", name, model)
			}
			schema.references[fk.ColumnName()] = *fk
		}

		if !field.IsAssociation() {
			schema.fields[field.ColumnName()] = *field
			continue
		}

		relationships[name] = field
	}

	for name, field := range relationships {
		fmt.Println(field)

		_, ok := schema.associations[field.FieldName()]
		if ok {
			continue
		}

		reference, err := NewReference(driver, schema, field)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot use '%s' as association for %T", name, model)
		}

		schema.associations[field.FieldName()] = *reference
		//
		// nextModel := field.ForeignKey.Reference.Model
		// if field.IsAssociationTypeMany() {
		// 	nextModel = field.ForeignKey.Model
		// }
		//
		// nextSchema, err := GetSchema(driver, nextModel)
		// if err != nil {
		// 	return Schema{}, err
		// }
		//
		// for k, v := range nextSchema.Associations {
		// 	key := fmt.Sprintf("%s.%s", field.FieldName, k)
		// 	_, ok := schema.Associations[key]
		// 	if !ok {
		// 		schema.Associations[key] = v
		// 	}
		// }
	}

	err = inferSchemaPrimaryKey(model, modelOpts, schema)
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func inferSchemaTimeKey(model Model, opts ModelOpts, schema *Schema, field *Field) error {
	if field.IsCreatedKey() {
		if schema.createdKey != nil {
			return errors.Errorf("sqlxx: %T must have only one created_at key", model)
		}
		schema.createdKey = field
	}

	if field.IsUpdatedKey() {
		if schema.updatedKey != nil {
			return errors.Errorf("sqlxx: %T must have only one updated_at key", model)
		}
		schema.updatedKey = field
	}

	if field.IsDeletedKey() {
		if schema.deletedKey != nil {
			return errors.Errorf("sqlxx: %T must have only one deleted_at key", model)
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
	return errors.Errorf("sqlxx: %T must have a primary key", model)
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
func GetColumns(driver Driver, model Model) (Columns, error) {
	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, errors.Wrap(err, "sqlxx: cannot fetch schema informations")
	}

	columns := schema.ColumnPaths()
	return columns, nil
}
