package sqlxx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oleiade/reflections"
	"github.com/serenize/snaker"
)

// Schema is a model schema.
type Schema struct {
	Columns      map[string]Column
	Associations map[string]RelatedField
}

// Column is a database column
type Column struct {
	TableName string
	Name      string
	Value     interface{}
	Tags      map[string]string
}

// PrefixedName returns the column name prefixed with the table name
func (c Column) PrefixedName() string {
	return fmt.Sprintf("%s.%s", c.TableName, c.Name)
}

// RelatedField represents an related field between two models.
type RelatedField struct {
	FK          Column
	FKReference Column
}

// GetSchema returns model's table columns, extracted by reflection.
// The returned map is modelFieldName -> table_name.column_name
func GetSchema(model Model) (*Schema, error) {
	fields, err := reflections.Fields(model)
	if err != nil {
		return nil, err
	}

	schema := &Schema{
		Columns:      map[string]Column{},
		Associations: map[string]RelatedField{},
	}

	for _, field := range fields {
		value, err := reflections.GetField(model, field)

		if err != nil {
			return nil, err
		}

		// Associations

		if isModel(value) {
			relatedField, err := newRelatedField(model, field)
			if err != nil {
				return nil, err
			}
			schema.Associations[field] = relatedField
			continue
		}

		// TODO: handle slice of models here

		col, err := newColumn(model, field)
		if err != nil {
			return nil, err
		}

		schema.Columns[field] = col
	}

	return schema, nil
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

// newRelatedField creates a new related field.
func newRelatedField(model Model, field string) (RelatedField, error) {
	relatedField := RelatedField{}

	relatedValue, err := reflections.GetField(model, field)
	if err != nil {
		return relatedField, err
	}

	related := relatedValue.(Model)

	relatedField.FK, err = newRelatedColumn(model, field)

	if err != nil {
		return relatedField, err
	}

	relatedField.FKReference, err = newForeignColumn(related, "ID")

	if err != nil {
		return relatedField, err
	}

	return relatedField, nil
}

func newRelatedColumn(model Model, field string) (Column, error) {
	column, err := newColumn(model, field)

	if err != nil {
		return Column{}, err
	}

	columnName, _ := column.Tags[SQLXStructTagName]

	if len(columnName) == 0 {
		columnName := fmt.Sprintf("%s_id", column.Name)
		column.Name = columnName
	}

	return column, nil
}

func newForeignColumn(model Model, field string) (Column, error) {
	// Retrieve the model type
	reflectType := reflect.ValueOf(model).Type()

	// If it's a pointer, we must get the elem to avoid double pointer errors
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	// Then we can safely cast
	reflected := reflect.New(reflectType).Interface().(Model)

	column, err := newColumn(reflected, field)

	if err != nil {
		return Column{}, err
	}

	return column, nil
}

// newColumn returns full column name from model, field and tag.
func newColumn(model Model, field string) (Column, error) {
	tags, err := extractTags(model, field)

	if err != nil {
		return Column{}, err
	}

	column, _ := tags[SQLXStructTagName]

	if len(column) == 0 {
		column = snaker.CamelToSnake(field)
	}

	value, err := reflections.GetField(model, field)

	if err != nil {
		return Column{}, err
	}

	return Column{
		TableName: model.TableName(),
		Name:      column,
		Value:     value,
		Tags:      tags,
	}, nil
}

// isModel returns true if the given reflect value is sqlxx.Model.
func isModel(value interface{}) bool {
	kind := reflect.TypeOf(value).Kind()

	if !(kind == reflect.Struct || kind == reflect.Ptr) {
		return false
	}

	typ := reflect.ValueOf(value).Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	_, ok := reflect.New(typ).Interface().(Model)

	return ok
}
