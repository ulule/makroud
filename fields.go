package sqlxx

import (
	"fmt"

	"github.com/serenize/snaker"
)

// Field is a field.
type Field struct {
	// Struct field name.
	Name string
	// Struct field metadata (reflect data).
	Meta Meta
	// Struct field tags.
	Tags Tags
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
func makeField(model Model, meta Meta) (Field, error) {
	tags := makeTags(meta.Field)

	var columnName string

	if dbName := tags.GetByKey(SQLXStructTagName, "field"); len(dbName) != 0 {
		columnName = dbName
	} else {
		columnName = snaker.CamelToSnake(meta.Name)
	}

	return Field{
		Name:       meta.Name,
		Meta:       meta,
		Tags:       tags,
		TableName:  model.TableName(),
		ColumnName: columnName,
	}, nil
}
