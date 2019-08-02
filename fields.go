package makroud

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"

	"github.com/ulule/makroud/reflectx"
	"github.com/ulule/makroud/snaker"
)

// Field defines the column name, field and options from model.
//
// For example: If we have an User, we could have this primary key defined in User's schema.
//
//     Field {
//         ModelName:  User,
//         TableName:  users,
//         FieldName:  AvatarID,
//         ColumnName: avatar_id,
//         ColumnPath: users.avatar_id,
//     }
//
type Field struct {
	modelName       string
	tableName       string
	fieldName       string
	fieldIndex      []int
	columnPath      string
	columnName      string
	foreignKey      string
	relationName    string
	isPrimaryKey    bool
	isForeignKey    bool
	isAssociation   bool
	isExcluded      bool
	hasRelation     bool
	hasDefault      bool
	hasULID         bool
	hasUUIDV1       bool
	hasUUIDV4       bool
	isCreatedKey    bool
	isUpdatedKey    bool
	isDeletedKey    bool
	rtype           reflect.Type
	associationType AssociationType
}

// ModelName define the model name of this field.
func (field Field) ModelName() string {
	return field.modelName
}

// FieldName define the struct field name used for this field.
func (field Field) FieldName() string {
	return field.fieldName
}

// FieldIndex define the struct field index used for this field.
func (field Field) FieldIndex() []int {
	return field.fieldIndex
}

// TableName returns the model name's table name of this field.
func (field Field) TableName() string {
	return field.tableName
}

// ColumnPath returns the field's full column path.
func (field Field) ColumnPath() string {
	return field.columnPath
}

// ColumnName returns the field's column name.
func (field Field) ColumnName() string {
	return field.columnName
}

// IsPrimaryKey returns if the field is a primary key.
func (field Field) IsPrimaryKey() bool {
	return field.isPrimaryKey
}

// IsExcluded returns if the field excluded. (anonymous, private...)
func (field Field) IsExcluded() bool {
	return field.isExcluded
}

// IsForeignKey returns if the field is a foreign key.
func (field Field) IsForeignKey() bool {
	return field.isForeignKey
}

// ForeignKey returns the field foreign key's table name.
func (field Field) ForeignKey() string {
	return field.foreignKey
}

// HasRelation returns if the field has a explicit relation.
func (field Field) HasRelation() bool {
	return field.hasRelation
}

// RelationName returns the field relation name.
func (field Field) RelationName() string {
	return field.relationName
}

// IsAssociation returns if the field is an association.
func (field Field) IsAssociation() bool {
	return field.isAssociation
}

// IsAssociationType returns if the field has given association type.
func (field Field) IsAssociationType(kind AssociationType) bool {
	return field.isAssociation && field.associationType == kind
}

// HasDefault returns if the field has a default value and should be in returning statement.
func (field Field) HasDefault() bool {
	return field.hasDefault
}

// HasULID returns if the field has a ulid type for it's primary key.
func (field Field) HasULID() bool {
	return field.hasULID
}

// HasUUIDV1 returns if the field has a uuid v1 type for it's primary key.
func (field Field) HasUUIDV1() bool {
	return field.hasUUIDV1
}

// HasUUIDV4 returns if the field has a uuid v4 type for it's primary key.
func (field Field) HasUUIDV4() bool {
	return field.hasUUIDV4
}

// IsCreatedKey returns if the field is a created key.
func (field Field) IsCreatedKey() bool {
	return field.isCreatedKey
}

// IsUpdatedKey returns if the field is a updated key.
func (field Field) IsUpdatedKey() bool {
	return field.isUpdatedKey
}

// IsDeletedKey returns if the field is a deleted key.
func (field Field) IsDeletedKey() bool {
	return field.isDeletedKey
}

// Type returns the reflect's type of the field.
func (field Field) Type() reflect.Type {
	return field.rtype
}

// String returns a human readable version of current instance.
func (field Field) String() string {
	return DebugField(field)
}

// NewField creates a new field using given model and name.
func NewField(driver Driver, schema *Schema, model Model, name string, args ...ModelOpts) (*Field, error) {
	if schema == nil {
		return nil, errors.New("schema is required to generate a field instance")
	}

	opts := defaultModelOpts()
	if len(args) > 0 {
		opts = args[0]
	}

	field, ok := reflectx.GetFieldByName(model, name)
	if !ok {
		return nil, errors.Errorf("field '%s' not found in model", name)
	}

	tags := GetTags(field)

	rtype := field.Type
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	modelName := reflectx.GetIndirectTypeName(model)
	tableName := model.TableName()

	columnName := tags.GetByKey(TagName, TagKeyColumn)
	if columnName == "" {
		columnName = snaker.CamelToSnake(name)
	}

	columnPath := fmt.Sprintf("%s.%s", tableName, columnName)

	isPrimaryKey := tags.HasKey(TagName, TagKeyPrimaryKey)
	foreignKey := tags.GetByKey(TagName, TagKeyForeignKey)
	isForeignKey := foreignKey != ""
	isExcluded := tags.HasKey(TagName, TagKeyIgnored) || field.PkgPath != ""
	hasDefault := tags.HasKey(TagName, TagKeyDefault)
	hasULID := tags.GetByKey(TagName, TagKeyPrimaryKey) == TagKeyULID
	hasUUIDV1 := tags.GetByKey(TagName, TagKeyPrimaryKey) == TagKeyUUIDV1
	hasUUIDV4 := tags.GetByKey(TagName, TagKeyPrimaryKey) == TagKeyUUIDV4

	isCreatedKey := columnName == opts.CreatedKey
	isUpdatedKey := columnName == opts.UpdatedKey
	isDeletedKey := columnName == opts.DeletedKey

	hasDefault = hasDefault || isCreatedKey || isUpdatedKey

	instance := &Field{
		modelName:    modelName,
		tableName:    tableName,
		fieldName:    field.Name,
		fieldIndex:   field.Index,
		columnName:   columnName,
		columnPath:   columnPath,
		isPrimaryKey: isPrimaryKey,
		isForeignKey: isForeignKey,
		foreignKey:   foreignKey,
		isExcluded:   isExcluded,
		isCreatedKey: isCreatedKey,
		isUpdatedKey: isUpdatedKey,
		isDeletedKey: isDeletedKey,
		hasDefault:   hasDefault,
		hasULID:      hasULID,
		hasUUIDV1:    hasUUIDV1,
		hasUUIDV4:    hasUUIDV4,
		rtype:        rtype,
	}

	// Early return if the field type is not an association.
	reference := toModel(rtype)
	if reference == nil {
		return instance, nil
	}

	return getFieldAssocitationType(driver, instance, rtype, reference, tags)
}

func getFieldAssocitationType(driver Driver, instance *Field, rtype reflect.Type,
	reference Model, tags Tags) (*Field, error) {

	if instance.isPrimaryKey {
		return nil, errors.Errorf("field '%s' cannot be a primary key and a association", instance.fieldName)
	}

	associationType := AssociationTypeUndefined
	if rtype.Kind() == reflect.Struct {
		associationType = AssociationTypeOne
	}
	if rtype.Kind() == reflect.Slice {
		associationType = AssociationTypeMany
	}

	if associationType == AssociationTypeUndefined {
		return nil, errors.Errorf("unable to infer the association type for field '%s'", instance.fieldName)
	}

	relationName := tags.GetByKey(TagName, TagKeyRelation)
	hasRelation := relationName != ""

	instance.isAssociation = true
	instance.associationType = associationType
	instance.columnName = ""
	instance.columnPath = ""
	instance.isCreatedKey = false
	instance.isForeignKey = false
	instance.isUpdatedKey = false
	instance.isDeletedKey = false
	instance.hasDefault = false
	instance.hasULID = false
	instance.hasRelation = hasRelation
	instance.relationName = relationName

	return instance, nil
}
