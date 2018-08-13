package sqlxx

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ulule/sqlxx/reflectx"
)

// FKType define a foreign key type.
type FKType uint8

// Foreign key types.
const (
	// FKUnknownType is an unknown foreign key.
	FKUnknownType = FKType(iota)
	// FKIntegerType uses an integer as foreign key.
	FKIntegerType
	// ForeignKeyString uses a string as foreign key.
	FKStringType
	// FKOptionalIntegerType uses an optional integer as foreign key.
	FKOptionalIntegerType
	// FKOptionalStringType uses an optional string as foreign key.
	FKOptionalStringType
)

// String returns a human readable version of current instance.
func (val FKType) String() string {
	switch val {
	case FKUnknownType:
		return ""
	case FKIntegerType:
		return "integer"
	case FKStringType:
		return "string"
	case FKOptionalIntegerType:
		return "option[integer]"
	case FKOptionalStringType:
		return "option[string]"
	default:
		panic(fmt.Sprintf("sqlxx: unknown foreign key type: %d", val))
	}
}

// IsCompatible returns if given primary key is compatible with foreign key.
func (val FKType) IsCompatible(key PKType) bool {
	switch val {
	case FKIntegerType, FKOptionalIntegerType:
		return key == PKIntegerType
	case FKStringType, FKOptionalStringType:
		return key == PKStringType
	default:
		return false
	}
}

// IsOptional returns if foreign key has an optional type.
func (val FKType) IsOptional() bool {
	return val == FKOptionalIntegerType || val == FKOptionalStringType
}

// ForeignKey is a composite object that define a foreign key for a model.
// This foreign key will be used later for Preload...
//
// For example: If we have an User, we could have this primary key defined in User's schema.
//
//     ForeignKey {
//         ModelName:  User,
//         TableName:  users,
//         FieldName:  AvatarID,
//         ColumnName: avatar_id,
//         ColumnPath: users.avatar_id,
//         Reference:  avatars,
//         Type:       int64,
//     }
//
type ForeignKey struct {
	Field
	fkTableName string
	fkType      FKType
}

// NewForeignKey creates a foreign key from a field instance.
func NewForeignKey(field *Field) (*ForeignKey, error) {
	pk := &ForeignKey{
		Field:       *field,
		fkTableName: field.ForeignKey(),
	}

	switch reflectx.GetType(field.Type()) {
	case reflectx.Int64Type:
		pk.fkType = FKIntegerType
	case reflectx.OptionalInt64Type:
		pk.fkType = FKOptionalIntegerType
	case reflectx.StringType:
		pk.fkType = FKStringType
	case reflectx.OptionalStringType:
		pk.fkType = FKOptionalStringType
	default:
		return nil, errors.Errorf("cannot use '%s' as foreign key type", field.Type().String())
	}

	return pk, nil
}

// Reference returns the foreign key's table name.
func (key ForeignKey) Reference() string {
	return key.fkTableName
}

// Type returns the foreign key's type.
func (key ForeignKey) Type() FKType {
	return key.fkType
}

// Reference defines a model relationship.
type Reference struct {
	Field
	isLocal bool
	local   ReferenceObject
	remote  ReferenceObject
}

// String returns a human readable version of current instance.
func (reference Reference) String() string {
	buffer := &bytes.Buffer{}
	debugReference(reference).write(buffer)
	return buffer.String()
}

// Local returns the local model.
func (reference Reference) Local() ReferenceObject {
	return reference.local
}

// Remote returns the remote model.
func (reference Reference) Remote() ReferenceObject {
	return reference.remote
}

// IsLocal returns if reference is local from the model, or from another model (remote).
func (reference Reference) IsLocal() bool {
	return reference.isLocal
}

// ReferenceObject defines a model used by Reference.
type ReferenceObject struct {
	schema       *Schema
	modelName    string
	tableName    string
	fieldName    string
	columnPath   string
	columnName   string
	isPrimaryKey bool
	isForeignKey bool
	pkType       PKType
	fkType       FKType
}

// Schema returns the reference schema.
func (object ReferenceObject) Schema() *Schema {
	return object.schema
}

// Model returns the reference model.
func (object ReferenceObject) Model() Model {
	return object.schema.Model()
}

// ModelName returns the model name of this reference.
func (object ReferenceObject) ModelName() string {
	return object.modelName
}

// TableName returns the table name of this reference.
func (object ReferenceObject) TableName() string {
	return object.tableName
}

// FieldName returns the field name of this reference.
func (object ReferenceObject) FieldName() string {
	return object.fieldName
}

// ColumnPath returns the full column path of this reference.
func (object ReferenceObject) ColumnPath() string {
	return object.columnPath
}

// ColumnName returns the column name of this reference.
func (object ReferenceObject) ColumnName() string {
	return object.columnName
}

// IsPrimaryKey returns if this reference is a primary key.
func (object ReferenceObject) IsPrimaryKey() bool {
	return object.isPrimaryKey
}

// IsForeignKey returns if this reference is a foreign key.
func (object ReferenceObject) IsForeignKey() bool {
	return object.isForeignKey
}

// PrimaryKeyType returns this reference primary key type.
func (object ReferenceObject) PrimaryKeyType() PKType {
	return object.pkType
}

// ForeignKeyType returns this reference foreign key type.
func (object ReferenceObject) ForeignKeyType() FKType {
	return object.fkType
}

// Columns returns this reference columns.
func (object ReferenceObject) Columns() []string {
	return object.schema.Columns().List()
}

// HasDeletedKey returns if an deleted key is defined for this reference.
func (object ReferenceObject) HasDeletedKey() bool {
	return object.schema.HasDeletedKey()
}

// DeletedKeyName returns reference schema deleted key column name.
func (object ReferenceObject) DeletedKeyName() string {
	return object.schema.DeletedKeyName()
}

// DeletedKeyPath returns reference schema deleted key column path.
func (object ReferenceObject) DeletedKeyPath() string {
	return object.schema.DeletedKeyPath()
}

// NewReference creates a reference from a field instance.
func NewReference(driver Driver, local *Schema, field *Field) (*Reference, error) {
	reference := ToModel(field.rtype)
	if reference == nil {
		return nil, errors.Errorf("invalid model: %s", field.rtype.String())
	}

	// TODO (novln): Allow circular reference...
	remote, err := GetSchema(driver, reference)
	if err != nil {
		return nil, err
	}

	switch field.associationType {
	case AssociationTypeOne:
		return newReferenceAsOne(driver, local, remote, field)

	case AssociationTypeMany:
		return newReferenceAsMany(driver, local, remote, field)

	default:
		panic("TODO")
		return nil, nil
	}
}

// Article.Author -> User
func newReferenceAsOne(driver Driver, local *Schema, remote *Schema, field *Field) (*Reference, error) {
	for _, element := range local.references {
		if element.ForeignKey() == remote.TableName() {
			target := remote.PrimaryKey()
			current := &Reference{
				Field:   *field,
				isLocal: true,
				local: ReferenceObject{
					schema:       local,
					modelName:    element.ModelName(),
					tableName:    element.TableName(),
					fieldName:    element.FieldName(),
					columnName:   element.ColumnName(),
					columnPath:   element.ColumnPath(),
					isForeignKey: true,
					fkType:       element.Type(),
				},
				remote: ReferenceObject{
					schema:       remote,
					modelName:    target.ModelName(),
					tableName:    target.TableName(),
					fieldName:    target.FieldName(),
					columnName:   target.ColumnName(),
					columnPath:   target.ColumnPath(),
					isPrimaryKey: true,
					pkType:       target.Type(),
				},
			}
			return current, nil
		}
	}

	for _, element := range remote.references {
		if element.ForeignKey() == local.TableName() {
			target := local.PrimaryKey()
			current := &Reference{
				Field:   *field,
				isLocal: false,
				local: ReferenceObject{
					schema:       local,
					modelName:    target.ModelName(),
					tableName:    target.TableName(),
					fieldName:    target.FieldName(),
					columnName:   target.ColumnName(),
					columnPath:   target.ColumnPath(),
					isPrimaryKey: true,
					pkType:       target.Type(),
				},
				remote: ReferenceObject{
					schema:       remote,
					modelName:    element.ModelName(),
					tableName:    element.TableName(),
					fieldName:    element.FieldName(),
					columnName:   element.ColumnName(),
					columnPath:   element.ColumnPath(),
					isForeignKey: true,
					fkType:       element.Type(),
				},
			}
			return current, nil
		}
	}

	return nil, errors.Errorf("cannot find foreign key for: %s.%s", field.ModelName(), field.FieldName())
}

// Article.Hashtags -> Tag
func newReferenceAsMany(driver Driver, local *Schema, remote *Schema, field *Field) (*Reference, error) {
	for _, element := range remote.references {
		if element.ForeignKey() == local.TableName() {
			target := local.PrimaryKey()
			current := &Reference{
				Field:   *field,
				isLocal: false,
				local: ReferenceObject{
					schema:       local,
					modelName:    target.ModelName(),
					tableName:    target.TableName(),
					fieldName:    target.FieldName(),
					columnName:   target.ColumnName(),
					columnPath:   target.ColumnPath(),
					isPrimaryKey: true,
					pkType:       target.Type(),
				},
				remote: ReferenceObject{
					schema:       remote,
					modelName:    element.ModelName(),
					tableName:    element.TableName(),
					fieldName:    element.FieldName(),
					columnName:   element.ColumnName(),
					columnPath:   element.ColumnPath(),
					isForeignKey: true,
					fkType:       element.Type(),
				},
			}
			return current, nil
		}
	}

	return nil, errors.Errorf("cannot find foreign key for: %s.%s", field.ModelName(), field.FieldName())
}
