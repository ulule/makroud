package sqlxx

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type XField struct {
	// Model name that contains this field
	ModelName string
	// Table name of the model that contains this field
	TableName string
	// The field name
	FieldName string
	// The database column name
	ColumnName string
	// Does this field is a primary key?
	IsPrimaryKey bool
	// IsForeignKey define if field should behave like a foreign key
	IsForeignKey bool
	defValue     string
	// IsArchiveKey define if field can be used as an archive key.
	IsArchiveKey bool
	// ArchiveValue is a generator to create a new value for the archive key.
	ArchiveValue func() interface{}
}

// HasDefault define if a default statement is defined in model.
func (e XField) HasDefault() bool {
	return e.defValue != ""
}

// Default define the default statement to use if value in model is undefined.
func (e XField) Default() string {
	return e.defValue
}

// FieldOption is used to define field configuration.
type FieldOption interface {
	apply(*fieldOp)
}

// fieldOption is a private implementation of FieldOption using a callback to configure a Field.
type fieldOption func(*fieldOp)

func (o fieldOption) apply(op *fieldOp) {
	o(op)
}

func IsForeignKey(reference string) FieldOption {
	return fieldOption(func(op *fieldOp) {
		op.IsForeignKey = true
		op.ReferenceName = reference
	})
}

func IsArchiveKey() FieldOption {
	return fieldOption(func(op *fieldOp) {
		op.IsArchiveKey = true
		op.ArchiveValue = func() interface{} {
			return time.Now()
		}
	})
}

// HasDefault define a default value for a field.
func HasDefault(value string) FieldOption {
	return fieldOption(func(op *fieldOp) {
		op.Default = value
	})
}

// NOTE
// ArchiveKey -> define a generator so we can use boolean or time.Time, for example...

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
type XPrimaryKey struct {
	modelName string
	tableName string
	pkName    string
	pkColumn  string
	pkType    PrimaryKeyType
}

// ModelName define the model name of this primary key.
func (e XPrimaryKey) ModelName() string {
	return e.modelName
}

// FieldName define the struct field name used as primary key.
func (e XPrimaryKey) FieldName() string {
	return e.pkName
}

// TableName returns the primary key's table name.
func (e XPrimaryKey) TableName() string {
	return e.tableName
}

// ColumnPath returns the primary key's full column path.
func (e XPrimaryKey) ColumnPath() string {
	return fmt.Sprintf("%s.%s", e.tableName, e.pkColumn)
}

// ColumnName returns the primary key's column path.
func (e XPrimaryKey) ColumnName() string {
	return e.pkColumn
}

// Type returns the primary key's type.
func (e XPrimaryKey) Type() PrimaryKeyType {
	return e.pkType
}

func (e XPrimaryKey) Value(model XModel) (interface{}, error) {
	id, ok := e.ValueOpt(model)
	if !ok {
		return nil, errors.New("invalid pk value")
	}
	return id, nil
}

func (e XPrimaryKey) ValueOpt(model XModel) (interface{}, bool) {
	switch e.pkType {
	case PrimaryKeyInteger:
		id, err := GetFieldValueInt64(model, e.pkName)
		if err != nil {
			return 0, false
		}
		if id == int64(0) {
			return 0, false
		}
		return id, true
	default:
		return nil, false
	}
}

// ForeignKey is a composite object that define a foreign key for a model.
//
// For example: If we have an User with an Avatar, we could have this foreign key defined in User's schema.
//
//     ForeignKey {
//         ModelName: User,
//         TableName: users,
//         FieldName: AvatarID,
//         ColumnName: avatar_id,
//         ReferenceName: Avatar,
//     }
//
type XForeignKey struct {
	// ModelName define the model name of this foreign key.
	ModelName string
	// TableName define the database table.
	TableName string
	// FieldName define the struct field name used as foreign key.
	FieldName string
	// ColumnName define the database column name of this foreign key.
	ColumnName string
	// ReferenceName define the reference model name.
	ReferenceName string
}

type XAssociationReference struct {
	// ModelName define the model name of this association.
	ModelName string
	// FieldName define the struct field name used to contains the value(s) of this association.
	FieldName string
	// Source define the field used to ... TODO
	Source XField
	// Reference define the field used to ... TODO
	Reference XField
	// Type define the association type. (ie: has once, has many, etc...)
	Type AssociationType
}
