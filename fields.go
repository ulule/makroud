package sqlxx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oleiade/reflections"
	"github.com/serenize/snaker"
)

// ----------------------------------------------------------------------------
// Tag
// ----------------------------------------------------------------------------

// Tag is a field tag.
type Tag map[string]string

// Get returns value for the given key or zero value.
func (t Tag) Get(key string) string {
	v, _ := t[key]
	return v
}

// ----------------------------------------------------------------------------
// Tags
// ----------------------------------------------------------------------------

// Tags are field tags.
type Tags map[string]Tag

// Get returns the given tag.
func (t Tags) Get(name string) (Tag, error) {
	tag, ok := t[name]
	if !ok {
		return nil, fmt.Errorf("tag %s does not exist", name)
	}

	return tag, nil
}

// GetByKey is a convenient shortcuts to get the value for a given tag key.
func (t Tags) GetByKey(name string, key string) string {
	if tag, err := t.Get(name); err == nil {
		if v := tag.Get(key); len(v) != 0 {
			return v
		}
	}

	return ""
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

func getFieldTags(structField reflect.StructField, names ...string) map[string]string {
	tags := map[string]string{}

	for _, name := range names {
		if _, ok := tags[name]; !ok {
			tags[name] = structField.Tag.Get(name)
		}
	}

	return tags
}

// ----------------------------------------------------------------------------
// Meta
// ----------------------------------------------------------------------------

// Meta are low level field metadata.
type Meta struct {
	Name  string
	Value reflect.Value
	Field reflect.StructField
	Type  reflect.Type
	Tags  Tags
}

func makeMeta(structField reflect.StructField, value reflect.Value) Meta {
	fieldName := structField.Name

	var structFieldType reflect.Type

	if structField.Type.Kind() == reflect.Ptr {
		structFieldType = structField.Type.Elem()
	} else {
		structFieldType = structField.Type
	}

	return Meta{
		Name:  fieldName,
		Value: value,
		Field: structField,
		Type:  structFieldType,
		Tags:  makeTags(structField),
	}
}

// ----------------------------------------------------------------------------
// Field
// ----------------------------------------------------------------------------

// Field is a field.
type Field struct {
	TableName string
	Name      string
	Value     interface{}
	Tags      Tags
	Meta      Meta
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

// newField returns full column name from model, field and tag.
func newField(model Model, meta Meta) (Field, error) {
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
func newForeignKeyField(model Model, meta Meta) (Field, error) {
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
	reflectType := reflectType(referencedModel)

	reflected := reflect.New(reflectType).Interface().(Model)

	f, ok := reflectType.FieldByName(name)
	if !ok {
		return Field{}, fmt.Errorf("Field %s does not exist", name)
	}

	meta := Meta{
		Name:  name,
		Field: f,
	}

	field, err := newField(reflected, meta)
	if err != nil {
		return Field{}, err
	}

	return field, nil
}

// ----------------------------------------------------------------------------
// Relation
// ----------------------------------------------------------------------------

// Relation represents an related field between two models.
type Relation struct {
	Type        RelationType
	FK          Field
	FKReference Field
}

// newRelatedField creates a new related field.
func newRelation(model Model, meta Meta, typ RelationType) (Relation, error) {
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

// getRelationType returns RelationType for the given reflect.Type.
func getRelationType(typ reflect.Type) RelationType {
	if typ.Kind() == reflect.Slice {
		if _, isModel := reflect.New(typ.Elem()).Interface().(Model); isModel {
			return RelationTypeManyToOne
		}

		return RelationTypeUnknown
	}

	if _, isModel := reflect.New(typ).Interface().(Model); isModel {
		return RelationTypeOneToMany
	}

	return RelationTypeUnknown
}
