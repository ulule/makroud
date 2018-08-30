package sqlxx

import (
	"fmt"

	"github.com/oklog/ulid"
	"github.com/pkg/errors"

	"github.com/ulule/sqlxx/reflectx"
)

// PKType define a primary key type.
type PKType uint8

// PrimaryKey types.
const (
	// PKUnknownType is an unknown primary key.
	PKUnknownType = PKType(iota)
	// PKIntegerType uses an integer as primary key.
	PKIntegerType
	// PrimaryKeyString uses a string as primary key.
	PKStringType
)

func (val PKType) String() string {
	switch val {
	case PKUnknownType:
		return ""
	case PKIntegerType:
		return "integer"
	case PKStringType:
		return "string"
	default:
		panic(fmt.Sprintf("sqlxx: unknown primary key type: %d", val))
	}
}

// IsCompatible returns if given foreign key is compatible with primary key.
func (val PKType) IsCompatible(key FKType) bool {
	switch val {
	case PKIntegerType:
		return key == FKIntegerType || key == FKOptionalIntegerType
	case PKStringType:
		return key == FKStringType || key == FKOptionalStringType
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
	pkType    PKType
	pkDefault PrimaryKeyDefault
}

// NewPrimaryKey creates a primary key from a field instance.
func NewPrimaryKey(field *Field) (*PrimaryKey, error) {
	pk := &PrimaryKey{
		Field:     *field,
		pkDefault: PrimaryKeyDBDefault,
	}

	switch reflectx.GetType(field.Type()) {
	case reflectx.Int64Type:
		pk.pkType = PKIntegerType
	case reflectx.StringType:
		pk.pkType = PKStringType
	default:
		return nil, errors.Errorf("cannot use '%s' as primary key type", field.Type().String())
	}

	if field.HasULID() {
		pk.pkDefault = PrimaryKeyULIDDefault
	}

	return pk, nil
}

// Type returns the primary key's type.
func (key PrimaryKey) Type() PKType {
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
	case PKIntegerType:
		id, err := reflectx.GetFieldValueInt64(model, key.FieldName())
		if err != nil || id == int64(0) {
			return int64(0), false
		}
		return id, true
	case PKStringType:
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
