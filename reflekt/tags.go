package reflekt

import (
	"reflect"
	"strings"
)

// TagProperty is a struct tag property
type TagProperty struct {
	Key   string
	Value string
}

// Tag is struct tag
type Tag struct {
	Name       string
	Properties []TagProperty
}

// Get returns value for the given property name.
func (t Tag) Get(key string) string {
	for _, p := range t.Properties {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}

// Tags a group of tag (usually for a struct field)
type Tags []Tag

// Get returns tag by name.
func (t Tags) Get(name string) *Tag {
	for _, tag := range t {
		if tag.Name == name {
			return &tag
		}
	}

	return nil
}

// Set sets the given tag into the slice.
func (t *Tags) Set(name string, property TagProperty) {
	var found bool

	tags := *t

	for i := range tags {
		if tags[i].Name == name {
			tags[i].Properties = append(tags[i].Properties, property)
			found = true
		}
	}

	if !found {
		tags = append(tags, Tag{
			Name:       name,
			Properties: []TagProperty{property},
		})
	}

	*t = tags
}

// GetByKey is a convenient shortcuts to get the value for a given tag key.
func (t Tags) GetByKey(name string, key string) string {
	if tag := t.Get(name); tag != nil {
		return tag.Get(key)
	}

	return ""
}

// GetFieldTags returns field tags
func GetFieldTags(field reflect.StructField, tagNames []string, propertyMapping map[string]string) Tags {
	rawTags := map[string]string{}

	for _, name := range tagNames {
		if _, ok := rawTags[name]; !ok {
			v := field.Tag.Get(name)
			if len(v) != 0 {
				rawTags[name] = v
			}
		}
	}

	tags := Tags{}

	for k, v := range rawTags {
		splits := strings.Split(v, ";")

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

			// Typically the case of sqlx tag that doesn't have key:value format (db:"field_name")
			if propertyKey, ok := propertyMapping[k]; ok {
				tags.Set(k, TagProperty{
					Key:   propertyKey,
					Value: strings.TrimSpace(splits[0]),
				})
				continue
			}

			// Typically, we have single property like "default", "ignored", etc.
			// To be consistent, we add true/false string values.
			if length == 1 {
				tags.Set(k, TagProperty{
					Key:   strings.TrimSpace(splits[0]),
					Value: "true",
				})
				continue
			}

			// Property named tag: key:value
			if length == 2 {
				tags.Set(k, TagProperty{
					Key:   strings.TrimSpace(splits[0]),
					Value: strings.TrimSpace(splits[1]),
				})
			}
		}
	}

	return tags
}
