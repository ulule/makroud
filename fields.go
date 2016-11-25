package sqlxx

import (
	"fmt"
	"reflect"
	"strings"

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

// Tags are field tags.
type Tags map[string]map[string]string

// SQLXFieldName returns SQLX field name.
func (t Tags) SQLXFieldName() string {
	var name string

	// Check if we have a "db:field_name". If so, use it.
	if values, ok := t[SQLXStructTagName]; ok && values != nil {
		if fieldName, ok := values["field"]; ok {
			name = fieldName
		}
	}

	return name
}

// makeTags returns field tags formatted.
func makeTags(structField reflect.StructField) Tags {
	tags := Tags{}

	rawTags := getFieldTags(structField, SupportedTags...)

	for k, v := range rawTags {
		splits := strings.Split(v, ";")

		tags[k] = map[string]string{}

		// Properties
		vals := []string{}
		for _, s := range splits {
			if len(s) != 0 {
				vals = append(vals, strings.TrimSpace(s))
			}
		}

		// Key / value
		for _, v := range vals {
			splits = strings.Split(v, ":")

			if len(splits) == 0 {
				continue
			}

			// format: db:"field_name" -> "field" -> "field_name"
			if k == SQLXStructTagName {
				tags[k]["field"] = strings.TrimSpace(splits[0])
				continue
			}

			if len(splits) >= 2 {
				tags[k][strings.TrimSpace(splits[0])] = strings.TrimSpace(splits[1])
			}
		}
	}

	return tags
}

// FieldMeta are low level field metadata.
type FieldMeta struct {
	Name  string
	Value reflect.Value
	Field reflect.StructField
	Type  reflect.Type
	Tags  Tags
}

// Field is a field.
type Field struct {
	TableName string
	Name      string
	Value     interface{}
	Tags      Tags
	Meta      FieldMeta
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
	}
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

	// Defaults to snakecase version of field name.
	name := snaker.CamelToSnake(meta.Name)

	// Get the SQLX one if any.
	if customName := tags.SQLXFieldName(); len(customName) != 0 {
		name = customName
	}

	return Field{
		TableName: model.TableName(),
		Name:      name,
		Tags:      tags,
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
	if customName := field.Tags.SQLXFieldName(); len(customName) != 0 {
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
