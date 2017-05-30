package sqlxx

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/sqlxx/reflectx"
)

// Schema is a model schema.
type Schema struct {
	model      Model
	modelName  string
	tableName  string
	pk         PrimaryKey
	fields     map[string]Field
	createdKey *Field
	updatedKey *Field
	deletedKey *Field
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
		return schema.updatedKey.ColumnPath()
	}
	return ""
}

// HasUpdatedKey returns if an updated key is defined for current schema.
func (schema Schema) HasUpdatedKey() bool {
	return schema.updatedKey != nil
}

// UpdateKeyPath returns schema updated key column's name.
func (schema Schema) UpdateKeyPath() string {
	if schema.HasUpdatedKey() {
		return schema.updatedKey.ColumnPath()
	}
	return ""
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
	return ""
}

// DeletedKeyValue returns schema deleted key's value.
func (schema Schema) DeletedKeyValue() interface{} {
	if schema.HasDeletedKey() {
		//return schema.deletedKey.ArchiveValue()
		return nil
	}
	return nil
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
func (schema Schema) WriteModel(mapper map[string]interface{}, model Model) error {
	for key, value := range mapper {
		if schema.pk.ColumnName() == key || schema.pk.ColumnPath() == key {
			err := reflectx.SetFieldValue(model, schema.pk.FieldName(), value)
			if err != nil {
				return err
			}
			continue
		}

		field, ok := schema.fields[key]
		if ok {
			err := reflectx.SetFieldValue(model, field.FieldName(), value)
			if err != nil {
				return err
			}
			continue
		}

		key = strings.TrimPrefix(key, fmt.Sprint(schema.TableName(), "."))
		field, ok = schema.fields[key]
		if ok {
			err := reflectx.SetFieldValue(model, field.FieldName(), value)
			if err != nil {
				return err
			}
			continue
		}
	}
	return nil
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

// newSchema returns model's table columns, extracted by reflection.
// The returned map is Model.FieldName -> table_name.column_name
func newSchema(driver Driver, model Model) (*Schema, error) {
	fields, err := reflectx.GetFields(model)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot use reflections to obtain %T fields", model)
	}

	schema := &Schema{
		// Model:     model,
		modelName: reflectx.GetIndirectTypeName(model),
		tableName: model.TableName(),
		fields:    map[string]Field{},
		// Associations: map[string]Field{},
	}

	for _, name := range fields {
		field, err := NewField(driver, schema, model, name)
		if err != nil {
			return nil, err
		}

		fmt.Println(field)

		if field.IsExcluded() {
			continue
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

		if !field.IsAssociation() {
			schema.fields[field.FieldName()] = *field
			continue
		}

		//
		// _, ok := schema.Associations[field.FieldName]
		// if ok {
		// 	continue
		// }

		// schema.Associations[field.FieldName] = field
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

	if schema.pk.TableName() == "" {
		return nil, errors.Errorf("sqlxx: %T must have a primary key", model)
	}

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
func GetColumns(driver Driver, model Model) (Columns, error) {
	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, errors.Wrap(err, "sqlxx: cannot fetch schema informations")
	}

	columns := schema.ColumnPaths()
	return columns, nil
}
