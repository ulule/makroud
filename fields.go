package sqlxx

import (
	"fmt"
	"reflect"

	"github.com/oleiade/reflections"
	"github.com/serenize/snaker"
)

// Struct tags
const (
	StructTagName     = "sqlxx"
	SQLXStructTagName = "db"
)

// SupportedTags are supported tags.
var SupportedTags = []string{
	StructTagName,
	SQLXStructTagName,
}

// RelationType is a field relation type.
type RelationType int

// Field types.
const (
	RelationTypeUnknown RelationType = iota
	RelationTypeOneToOne
	RelationTypeOneToMany
	RelationTypeManyToOne
	RelationTypeManyToMany
)

// RelationTypes are supported relations types.
var RelationTypes = map[RelationType]bool{
	RelationTypeOneToOne:   true,
	RelationTypeOneToMany:  true,
	RelationTypeManyToOne:  true,
	RelationTypeManyToMany: true,
}

// FieldMeta are low level field metadata.
type FieldMeta struct {
	Name  string
	Value reflect.Value
	Field reflect.StructField
	Type  reflect.Type
	Tags  Tags
}

func makeFieldMeta(structField reflect.StructField, value reflect.Value) FieldMeta {
	fieldName := structField.Name

	var structFieldType reflect.Type

	if structField.Type.Kind() == reflect.Ptr {
		structFieldType = structField.Type.Elem()
	} else {
		structFieldType = structField.Type
	}

	return FieldMeta{
		Name:  fieldName,
		Value: value,
		Field: structField,
		Type:  structFieldType,
		Tags:  makeTags(structField),
	}
}

// Field is a field.
type Field struct {
	TableName string
	Name      string
	Value     interface{}
	Tags      Tags
	Meta      FieldMeta
	IsPrimary bool
}

// HasValue returns true if value is a zero value.
func (f Field) HasValue() bool {
	return !isZeroValue(f.Value)
}

// PrefixedName returns the column name prefixed with the table name
func (f Field) PrefixedName() string {
	return fmt.Sprintf("%s.%s", f.TableName, f.Name)
}

// Relation represents an related field between two models.
type Relation struct {
	Type        RelationType
	FK          Field
	FKReference Field
}

// newRelatedField creates a new related field.
func newRelation(model Model, meta FieldMeta, typ RelationType) (Relation, error) {
	var err error

	relation := Relation{
		Type: typ,
	}

	related, err := reflections.GetField(model, meta.Name)
	if err != nil {
		return relation, err
	}

	relation.FK, err = newForeignKeyField(model, meta)
	if err != nil {
		return relation, err
	}

	relation.FKReference, err = newForeignKeyReferenceField(related.(Model), "ID")
	if err != nil {
		return relation, err
	}

	return relation, nil
}

// newField returns full column name from model, field and tag.
func newField(model Model, meta FieldMeta) (Field, error) {
	tags := makeTags(meta.Field)

	var name string

	if dbName := tags.GetByKey(SQLXStructTagName, "field"); len(dbName) != 0 {
		name = dbName
	} else {
		name = snaker.CamelToSnake(meta.Name)
	}

	v := reflectValue(meta.Value)

	var value interface{}
	if v.IsValid() {
		value = v.Interface()
	}

	return Field{
		TableName: model.TableName(),
		Name:      name,
		Tags:      tags,
		Value:     value,
	}, nil
}

// newForeignKeyField returns foreign key field.
func newForeignKeyField(model Model, meta FieldMeta) (Field, error) {
	field, err := newField(model, meta)
	if err != nil {
		return Field{}, err
	}

	// Defaults to "fieldname_id"
	field.Name = fmt.Sprintf("%s_id", field.Name)

	// Get the SQLX one if any.
	if customName := field.Tags.GetByKey(SQLXStructTagName, "field"); len(customName) != 0 {
		field.Name = customName
	}

	return field, nil
}

// newForeignKeyReferenceField returns a foreign key reference field.
func newForeignKeyReferenceField(referencedModel Model, name string) (Field, error) {
	reflectType := getReflectedType(referencedModel)

	reflected := reflect.New(reflectType).Interface().(Model)

	f, ok := reflectType.FieldByName(name)
	if !ok {
		return Field{}, fmt.Errorf("Field %s does not exist", name)
	}

	meta := FieldMeta{
		Name:  name,
		Field: f,
	}

	field, err := newField(reflected, meta)
	if err != nil {
		return Field{}, err
	}

	return field, nil
}
