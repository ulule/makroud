package sqlxx

import (
	"fmt"
	"reflect"
	"strings"
)

// Meta are low level field metadata.
type Meta struct {
	Name  string
	Field reflect.StructField
	Type  reflect.Type
	Tags  Tags
}

// makeMeta returns field reflect data.
func makeMeta(field reflect.StructField) Meta {
	var (
		fieldName = field.Name
		fieldType = field.Type
	)

	if field.Type.Kind() == reflect.Ptr {
		fieldType = field.Type.Elem()
	}

	return Meta{
		Name:  fieldName,
		Field: field,
		Type:  fieldType,
		Tags:  makeTags(field),
	}
}

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
			length := len(splits)

			if length == 0 {
				continue
			}

			// format: db:"field_name" -> "field" -> "field_name"
			if k == SQLXStructTagName {
				tags[k]["field"] = strings.TrimSpace(splits[0])
				continue
			}

			// Typically, we have single property like "default", "ignored", etc.
			// To be consistent, we add true/false string values.
			if length == 1 {
				tags[k][strings.TrimSpace(splits[0])] = "true"
				continue
			}

			// Typical key / value
			if length == 2 {
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
