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
	TableName    string
	Name         string
	PrefixedName string
	Value        interface{}
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

		// Columns

		tag, err := reflections.GetFieldTag(model, field, SQLXStructTagName)
		if err != nil {
			return nil, err
		}

		col, err := newColumn(model, field, tag, value, false, false)
		if err != nil {
			return nil, err
		}

		schema.Columns[field] = col
	}

	return schema, nil
}

func extractTags(model Model, field string) (map[string]string, error) {
	tag, err := reflections.GetFieldTag(model, field, StructTagName)

	if err != nil {
		return nil, err
	}

	results := map[string]string{}

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

	dbTag, err := reflections.GetFieldTag(model, field, SQLXStructTagName)

	if err != nil {
		return relatedField, err
	}

	tags, err := extractTags(model, field)

	if err != nil {
		return relatedField, err
	}

	tag, _ := tags["related"]

	related := relatedValue.(Model)

	relatedField.FK, err = newColumn(model, field, dbTag, relatedValue, true, false)
	if err != nil {
		return relatedField, err
	}

	relatedField.FKReference, err = newColumn(related, field, tag, relatedValue, true, true)
	if err != nil {
		return relatedField, err
	}

	return relatedField, nil
}

// newColumn returns full column name from model, field and tag.
func newColumn(model Model, field string, tag string, value interface{}, isRelated bool, isReference bool) (Column, error) {
	// Retrieve the model type
	reflectType := reflect.ValueOf(model).Type()

	// If it's a pointer, we must get the elem to avoid double pointer errors
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	// Then we can safely cast
	reflected := reflect.New(reflectType).Interface().(Model)

	hasTag := len(tag) > 0

	// Build column name from tag or field
	column := tag
	if !hasTag {
		column = snaker.CamelToSnake(field)
	}

	// It's not a related field, early return
	if !isRelated {
		return Column{
			TableName:    reflected.TableName(),
			Name:         column,
			PrefixedName: fmt.Sprintf("%s.%s", reflected.TableName(), column),
			Value:        value,
		}, nil
	}

	// Reference primary key fields are "id" and "field_id"
	if isReference {
		column = "id"

		if hasTag {
			column = tag
		}

		return Column{
			TableName:    reflected.TableName(),
			Name:         column,
			PrefixedName: fmt.Sprintf("%s.%s", reflected.TableName(), column),
			Value:        value,
		}, nil
	}

	// It's a foreign key
	column = fmt.Sprintf("%s_id", column)
	if hasTag {
		column = tag
	}

	return Column{
		TableName:    reflected.TableName(),
		Name:         column,
		PrefixedName: fmt.Sprintf("%s.%s", reflected.TableName(), column),
		Value:        value,
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
