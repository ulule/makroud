package sqlxx

import (
	"fmt"
	"reflect"

	"github.com/serenize/snaker"
	"github.com/ulule/sqlxx/reflekt"
)

// ----------------------------------------------------------------------------
// Field
// ----------------------------------------------------------------------------

// Field is a field.
type Field struct {
	// Struct field name.
	Name string
	// Struct field metadata (reflect data).
	Meta FieldMeta
	// Struct field tags.
	Tags reflekt.FieldTags
	// TableName is the database table name.
	TableName string
	// ColumnName is the database column name.
	ColumnName string
	// Is a primary key?
	IsPrimary bool
}

// ColumnPath returns the column name prefixed with the table name.
func (f Field) ColumnPath() string {
	return fmt.Sprintf("%s.%s", f.TableName, f.ColumnName)
}

// makeField returns full column name from model, field and tag.
func makeField(model Model, meta FieldMeta) (Field, error) {
	var columnName string

	if dbName := meta.Tags.GetByKey(SQLXStructTagName, "field"); len(dbName) != 0 {
		columnName = dbName
	} else {
		columnName = snaker.CamelToSnake(meta.Name)
	}

	return Field{
		Name:       meta.Name,
		Meta:       meta,
		Tags:       meta.Tags,
		TableName:  model.TableName(),
		ColumnName: columnName,
	}, nil
}

// IsExcludedField returns true if field must be excluded from schema.
func IsExcludedField(meta FieldMeta) bool {
	// Skip unexported fields
	if len(meta.Field.PkgPath) != 0 {
		return true
	}

	// Skip db:"-"
	if f := meta.Tags.GetByKey(SQLXStructTagName, "field"); f == "-" {
		return true
	}

	return false
}

// IsPrimaryKeyField returns true if field is a primary key field.
func IsPrimaryKeyField(meta FieldMeta) bool {
	return (meta.Name == PrimaryKeyFieldName || len(meta.Tags.GetByKey(StructTagName, StructTagPrimaryKey)) != 0)
}

// ----------------------------------------------------------------------------
// Field meta
// ----------------------------------------------------------------------------

// FieldMeta are low level field metadata.
type FieldMeta struct {
	Name  string
	Field reflect.StructField
	Type  reflect.Type
	Tags  reflekt.FieldTags
}

// GetFieldMeta returns field reflect data.
func GetFieldMeta(field reflect.StructField, tags []string, tagsMapping map[string]string) FieldMeta {
	var (
		fieldName = field.Name
		fieldType = field.Type
	)

	if field.Type.Kind() == reflect.Ptr {
		fieldType = field.Type.Elem()
	}

	return FieldMeta{
		Name:  fieldName,
		Field: field,
		Type:  fieldType,
		Tags:  reflekt.GetFieldTags(field, tags, tagsMapping),
	}
}
