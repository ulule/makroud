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
	// PrimaryKeyUnknownType is an unknown primary key.
	PrimaryKeyUnknownType = PrimaryKeyType(iota)
	// PrimaryKeyIntegerType uses an integer as primary key.
	PrimaryKeyIntegerType
	// PrimaryKeyString uses a string as primary key.
	PrimaryKeyStringType
)

func (val PrimaryKeyType) String() string {
	switch val {
	case PrimaryKeyUnknownType:
		return ""
	case PrimaryKeyIntegerType:
		return "int64"
	case PrimaryKeyStringType:
		return "string"
	default:
		panic(fmt.Sprintf("sqlxx: unknown primary key type: %d", val))
	}
}

// Equals returns if given foreign key has the same type as primary key.
func (val PrimaryKeyType) Equals(key ForeignKeyType) bool {
	switch val {
	case PrimaryKeyIntegerType:
		return key == ForeignKeyIntegerType
	case PrimaryKeyStringType:
		return key == ForeignKeyStringType
	default:
		return false
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

func (e PrimaryKeyDefault) String() string {
	switch e {
	case PrimaryKeyDBDefault:
		return "db"
	case PrimaryKeyULIDDefault:
		return "ulid"
	default:
		panic(fmt.Sprintf("sqlxx: unknown primary key default types: %d", e))
	}
}

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
//         Type: int64,
//         Default: db,
//     }
//
type PrimaryKey struct {
	Field
	pkType    PrimaryKeyType
	pkDefault PrimaryKeyDefault
}

// NewPrimaryKey creates a primary key from a field instance.
func NewPrimaryKey(field *Field) (*PrimaryKey, error) {
	pk := &PrimaryKey{
		Field:     *field,
		pkDefault: PrimaryKeyDBDefault,
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
		id, err := reflectx.GetFieldValueInt64(model, key.FieldName())
		if err != nil || id == int64(0) {
			return int64(0), false
		}
		return id, true
	case PrimaryKeyStringType:
		id, err := reflectx.GetFieldValueString(model, key.FieldName())
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
