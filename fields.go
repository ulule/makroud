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
	// The database columnn path
	ColumnPath string

	// Is this field a foreign key?
	IsForeignKey bool
	// Is this field an association? (preload)
	IsAssociation bool
	// The association struct instance
	Association *Association
}

// NewField returns full column name from model, field and tag.
func NewField(structField reflect.StructField, model Model) (Field, error) {
	var (
		name         = structField.Name
		tags         = reflekt.GetFieldTags(structField, SupportedTags, TagsMapping)
		modelName    = GetModelName(model)
		tableName    = model.TableName()
		columnName   = snaker.CamelToSnake(name)
		isPrimaryKey = name == PrimaryKeyFieldName || len(tags.GetByKey(StructTagName, StructTagPrimaryKey)) != 0
		isExcluded   = IsExcludedField(structField, tags)
	)

	if v := tags.GetByKey(SQLXStructTagName, StructTagSQLXField); v != "" {
		columnName = v
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
		ColumnPath:   fmt.Sprintf("%s.%s", tableName, columnName),
	}

	association, isAssociation, err := NewAssociation(structField)
	if err != nil {
		return field, err
	}

	if isAssociation {
		field.IsAssociation = true
		field.Association = association
	}

	if IsForeignKey(field) {
		field.IsForeignKey = true
	}

	return field, nil
}

// IsForeignKey returns true if the given fields looks like a foreign key or
// had been explicitly set as foreign key field.
func IsForeignKey(f Field) bool {
	if f.Tags.HasKey(StructTagName, StructTagForeignKey) {
		return true
	}

	// Typically MyFieldID/MyFieldPK
	if len(f.Name) > len(PrimaryKeyFieldName) && strings.HasSuffix(f.Name, PrimaryKeyFieldName) {
		return true
	}

	return false
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
