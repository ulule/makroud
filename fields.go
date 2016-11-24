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
type Tags map[string]string

// Field is a field.
type Field struct {
	TableName string
	Name      string
	Value     interface{}
	Tags      Tags
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
func newRelation(model Model, field string, typ RelationType) (Relation, error) {
	relatedField := Relation{Type: typ}

	relatedValue, err := reflections.GetField(model, field)
	if err != nil {
		return relatedField, err
	}

	related := relatedValue.(Model)

	relatedField.FK, err = newRelated(model, field)
	if err != nil {
		return relatedField, err
	}

	relatedField.FKReference, err = newForeignKeyField(related, "ID")
	if err != nil {
		return relatedField, err
	}

	return relatedField, nil
}

// newField returns full column name from model, field and tag.
func newField(model Model, field string) (Field, error) {
	tags, err := extractTags(model, field)
	if err != nil {
		return Field{}, err
	}

	column, _ := tags[SQLXStructTagName]

	if len(column) == 0 {
		column = snaker.CamelToSnake(field)
	}

	value, err := reflections.GetField(model, field)
	if err != nil {
		return Field{}, err
	}

	return Field{
		TableName: model.TableName(),
		Name:      column,
		Value:     value,
		Tags:      tags,
	}, nil
}

// newRelatedField creates a new related field.
func newRelatedField(model Model, field string) (Relation, error) {
	relatedField := Relation{}

	relatedValue, err := reflections.GetField(model, field)
	if err != nil {
		return relatedField, err
	}

	related := relatedValue.(Model)

	relatedField.FK, err = newRelated(model, field)
	if err != nil {
		return relatedField, err
	}

	relatedField.FKReference, err = newForeignKeyField(related, "ID")
	if err != nil {
		return relatedField, err
	}

	return relatedField, nil
}

func newRelated(model Model, field string) (Field, error) {
	f, err := newField(model, field)
	if err != nil {
		return Field{}, err
	}

	if _, ok := f.Tags[SQLXStructTagName]; !ok {
		f.Name = fmt.Sprintf("%s_id", f.Name)
	}

	return f, nil
}

// newForeignKeyField returns a foreign key field.
func newForeignKeyField(model Model, field string) (Field, error) {
	// Retrieve the model type
	reflectType := reflect.ValueOf(model).Type()

	// If it's a pointer, we must get the elem to avoid double pointer errors
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	// Then we can safely cast
	reflected := reflect.New(reflectType).Interface().(Model)

	f, err := newField(reflected, field)
	if err != nil {
		return Field{}, err
	}

	return f, nil
}

// extractTags return the struct tags (Map[key => value]) with sqlxx prefix
func extractTags(model Model, field string) (map[string]string, error) {
	tag, err := reflections.GetFieldTag(model, field, StructTagName)

	if err != nil {
		return nil, err
	}

	results := map[string]string{}

	column, err := reflections.GetFieldTag(model, field, SQLXStructTagName)

	if err != nil {
		return results, err
	}

	if len(column) > 0 {
		results[SQLXStructTagName] = column
	}

	if tag == "" {
		return results, err
	}

	parts := strings.Split(tag, " ")

	for _, part := range parts {
		splits := strings.Split(part, ":")
		results[splits[0]] = splits[1]
	}

	return results, err
}
