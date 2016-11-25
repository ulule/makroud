package sqlxx

import (
	"fmt"
	"reflect"
	"strings"
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
		if v := tag.Get(key); len(name) != 0 {
			return v
		}
	}
	return ""
}

// ----------------------------------------------------------------------------
// Initializers
// ----------------------------------------------------------------------------

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
