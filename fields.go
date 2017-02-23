package sqlxx

import (
	"fmt"
	"strings"

	"reflect"

	"github.com/serenize/snaker"
	"github.com/ulule/sqlxx/reflekt"
)

// Field is a field.
type Field struct {
	// The reflect.StructField instance
	StructField reflect.StructField
	// The field name
	Name string
	// The field struct tags
	Tags reflekt.FieldTags
	// Does this field is a primary key?
	IsPrimaryKey bool
	// Does this field is excluded? (anonymous, private, non-sql...)
	IsExcluded bool

	// Model is the zero-valued field's model used to generate schema from.
	Model Model
	// Model name that contains this field
	ModelName string
	// Table name of the model that contains this field
	TableName string
	// The database column name
	ColumnName string

	// Is this field an association? (preload)
	IsAssociation bool
	// The association struct instance
	Association *Association
}

// ColumnPath returns database full column path.
func (f Field) ColumnPath() string {
	return fmt.Sprintf("%s.%s", f.TableName, f.ColumnName)
}

// IsForeignKey returns true if the given fields looks like a foreign key or
// had been explicitly set as foreign key field.
func (f Field) IsForeignKey() bool {
	if f.Tags.HasKey(StructTagName, StructTagForeignKey) {
		return true
	}

	// Typically MyFieldID/MyFieldPK
	if len(f.Name) > len(PrimaryKeyFieldName) && strings.HasSuffix(f.Name, PrimaryKeyFieldName) {
		return true
	}

	return false
}

// ReverseAssociation reverses association (used for AssociationTypeMany).
func (f *Field) ReverseAssociation() {
	// User.Avatars -> user_id
	var (
		tableName      = f.TableName
		assocTableName = f.Association.TableName
	)

	f.Association.TableName = tableName
	f.TableName = assocTableName
	f.ColumnName = fmt.Sprintf("%s_%s", snaker.CamelToSnake(f.ModelName), f.Association.PrimaryKeyField.ColumnName)
}

// NewField returns full column name from model, field and tag.
func NewField(structField reflect.StructField, model Model) (Field, error) {
	var (
		err                 error
		name                = structField.Name
		tags                = reflekt.GetFieldTags(structField, SupportedTags, TagsMapping)
		modelName           = GetModelName(model)
		tableName           = model.TableName()
		columnName          = snaker.CamelToSnake(name)
		columnNameOverrided = false
		isPrimaryKey        = name == PrimaryKeyFieldName || len(tags.GetByKey(StructTagName, StructTagPrimaryKey)) != 0
		isExcluded          = IsExcludedField(structField, tags)
	)

	if v := tags.GetByKey(SQLXStructTagName, StructTagSQLXField); v != "" {
		columnName = v
		columnNameOverrided = true
	}

	field := Field{
		StructField:  structField,
		Name:         name,
		Tags:         tags,
		IsPrimaryKey: isPrimaryKey,
		IsExcluded:   isExcluded,
		Model:        model,
		ModelName:    modelName,
		TableName:    tableName,
		ColumnName:   columnName,
	}

	field.Association, field.IsAssociation, err = NewAssociation(structField)
	if err != nil {
		return field, err
	}

	if field.IsAssociation {
		if !columnNameOverrided {
			field.ColumnName = fmt.Sprintf("%s_id", columnName)
		}

		if field.Association.Type == AssociationTypeMany {
			field.ReverseAssociation()
		}
	}

	return field, nil
}

// IsExcludedField returns true if field must be excluded from schema.
func IsExcludedField(f reflect.StructField, tags reflekt.FieldTags) bool {
	if f.PkgPath != "" {
		return true
	}

	// Skip db:"-"
	if tag := tags.GetByKey(SQLXStructTagName, StructTagSQLXField); tag == "-" {
		return true
	}

	return false
}
