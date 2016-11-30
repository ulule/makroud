package sqlxx

import (
	"fmt"
	"reflect"

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

// makeForeignKeyField returns foreign key field.
func makeForeignKeyField(model Model, meta Meta) (Field, error) {
	field, err := makeField(model, meta)
	if err != nil {
		return Field{}, err
	}

	// Defaults to "fieldname_id"
	field.ColumnName = fmt.Sprintf("%s_id", field.ColumnName)

	// Get the SQLX one if any.
	if customName := field.Tags.GetByKey(SQLXStructTagName, "field"); len(customName) != 0 {
		field.ColumnName = customName
	}

	return field, nil
}

// makeReferenceField returns a foreign key reference field.
func makeReferenceField(referencedModel Model, name string) (Field, error) {
	reflectType := reflectType(referencedModel)

	reflected := reflect.New(reflectType).Interface().(Model)

	f, ok := reflectType.FieldByName(name)
	if !ok {
		return Field{}, fmt.Errorf("Field %s does not exist", name)
	}

	meta := Meta{Name: name, Field: f}

	field, err := makeField(reflected, meta)
	if err != nil {
		return Field{}, err
	}

	return field, nil
}
