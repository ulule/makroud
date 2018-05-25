package sqlxx

import (
	"fmt"
	"reflect"

	"github.com/oklog/ulid"
	"github.com/pkg/errors"

	"github.com/ulule/sqlxx/reflectx"
)

// PrimaryKeyType define a primary key type.
type PrimaryKeyType uint8

// PrimaryKey types.
const (
	// PrimaryKeyIntegerType uses an integer as primary key.
	PrimaryKeyIntegerType = PrimaryKeyType(iota)
	// PrimaryKeyString uses a string as primary key.
	PrimaryKeyStringType
)

func (e PrimaryKeyType) String() string {
	switch e {
	case PrimaryKeyIntegerType:
		return "int64"
	case PrimaryKeyStringType:
		return "string"
	default:
		panic(fmt.Sprintf("sqlxx: unknown primary key type: %d", e))
	}
}

// PrimaryKeyDefault define how primary key value is generated.
type PrimaryKeyDefault uint8

// PrimaryKey default types.
const (
	// PrimaryKeyDBDefault uses internal db mechanism to define primary key value.
	PrimaryKeyDBDefault = PrimaryKeyDefault(iota)
	// PrimaryKeyULIDDefault uses a ulid generator to define primary key value.
	PrimaryKeyULIDDefault
)

// TODO Add unit test

// PrimaryKey is a composite object that define a primary key for a model.
//
// For example: If we have an User, we could have this primary key defined in User's schema.
//
//     PrimaryKey {
//         ModelName: User,
//         TableName: users,
//         FieldName: ID,
//         ColumnName: id,
//         ColumnPath: users.id,
//         Type: integer,
//     }
//
type PrimaryKey struct {
	modelName    string
	tableName    string
	pkName       string
	pkColumnName string
	pkColumnPath string
	pkType       PrimaryKeyType
	pkDefault    PrimaryKeyDefault
}

// NewPrimaryKey creates a primary key from a field instance.
func NewPrimaryKey(field *Field) (*PrimaryKey, error) {
	pk := &PrimaryKey{
		modelName:    field.modelName,
		tableName:    field.tableName,
		pkName:       field.fieldName,
		pkColumnName: field.columnName,
		pkColumnPath: field.columnPath,
		pkDefault:    PrimaryKeyDBDefault,
	}

	switch field.Type().Kind() {
	case reflect.String:
		pk.pkType = PrimaryKeyStringType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		pk.pkType = PrimaryKeyIntegerType
	default:
		return nil, errors.Errorf("cannot use '%s' as primary key type", field.Type().String())
	}

	if field.HasULID() {
		pk.pkDefault = PrimaryKeyULIDDefault
	}

	return pk, nil
}

// ModelName define the model name of this primary key.
func (key PrimaryKey) ModelName() string {
	return key.modelName
}

// FieldName define the struct field name used as primary key.
func (key PrimaryKey) FieldName() string {
	return key.pkName
}

// TableName returns the primary key's table name.
func (key PrimaryKey) TableName() string {
	return key.tableName
}

// ColumnPath returns the primary key's full column path.
func (key PrimaryKey) ColumnPath() string {
	return key.pkColumnPath
}

// ColumnName returns the primary key's column name.
func (key PrimaryKey) ColumnName() string {
	return key.pkColumnName
}

// Type returns the primary key's type.
func (key PrimaryKey) Type() PrimaryKeyType {
	return key.pkType
}

// Default returns the primary key's default mechanism.
func (key PrimaryKey) Default() PrimaryKeyDefault {
	return key.pkDefault
}

// Value returns the primary key's value, or an error if undefined.
func (key PrimaryKey) Value(model Model) (interface{}, error) {
	id, ok := key.ValueOpt(model)
	if !ok {
		return nil, errors.New("invalid pk value")
	}
	return id, nil
}

// ValueOpt may returns the primary key's value, if defined.
func (key PrimaryKey) ValueOpt(model Model) (interface{}, bool) {
	switch key.pkType {
	case PrimaryKeyIntegerType:
		id, err := reflectx.GetFieldValueInt64(model, key.pkName)
		if err != nil || id == int64(0) {
			return int64(0), false
		}
		return id, true
	case PrimaryKeyStringType:
		id, err := reflectx.GetFieldValueString(model, key.pkName)
		if err != nil || id == "" {
			return "", false
		}
		return id, true
	default:
		return nil, false
	}
}

// GenerateULID generates a new ulid.
func GenerateULID(driver Driver) string {
	return ulid.MustNew(ulid.Now(), driver.entropy()).String()
}
