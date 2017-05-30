package sqlxx

import (
	"reflect"
	"strings"
)

const (
	// TagName defines sqlxx tag namespace.
	TagName = "sqlxx"
	// TagNameAlt defines sqlx tag namespace.
	TagNameAlt = "db"
)

// TagsList is a list of supported tags, this include sqlx and sqlxx one.
var TagsList = []string{
	TagName,
	TagNameAlt,
}

// TagsMapper is a mapper that convert sqlx tag to sqlxx one...
var TagsMapper = map[string]string{
	TagNameAlt: TagKeyColumn,
}

// Tag modifiers
const (
	TagKeyIgnored    = "-"
	TagKeyDefault    = "default"
	TagKeyColumn     = "column"
	TagKeyForeignKey = "fk"
	TagKeyPrimaryKey = "pk"
)

// TagProperty is a struct tag property.
type TagProperty struct {
	key   string
	value string
}

// Key returns tag property key.
func (prop TagProperty) Key() string {
	return prop.key
}

// Value returns tag property value.
func (prop TagProperty) Value() string {
	return prop.value
}

// String returns a human readable version of current instance.
func (prop TagProperty) String() string {
	return DebugTagProperty(prop)
}

// Tag is tag defined in a model.
type Tag struct {
	name       string
	properties []TagProperty
}

// Name returns tag name.
func (tag Tag) Name() string {
	return tag.name
}

// Properties returns tag properties.
func (tag Tag) Properties() []TagProperty {
	return tag.properties
}

// String returns a human readable version of current instance.
func (tag Tag) String() string {
	return DebugTag(tag)
}

// Get returns value for the given property name.
func (tag Tag) Get(key string) string {
	for _, property := range tag.properties {
		if property.key == key {
			return property.value
		}
	}
	return ""
}

// Tags is a group of tag defined in a model.
type Tags []Tag

// String returns a human readable version of current instance.
func (tags Tags) String() string {
	return DebugTags(tags)
}

// Get returns tag by name.
func (tags Tags) Get(name string) *Tag {
	for i := range tags {
		if tags[i].name == name {
			return &tags[i]
		}
	}
	return nil
}

// Set sets the given tag into the slice.
func (tags *Tags) Set(name string, property TagProperty) {
	copy := *tags

	found := false
	for i := range copy {
		if copy[i].name == name {
			copy[i].properties = append(copy[i].properties, property)
			found = true
		}
	}

	if !found {
		copy = append(copy, Tag{
			name: name,
			properties: []TagProperty{
				property,
			},
		})
	}

	*tags = copy
}

// HasKey is a convenient shortcuts to check if a key is present.
func (tags Tags) HasKey(name string, key string) bool {
	tag := tags.Get(name)
	if tag != nil {
		return tag.Get(key) != ""
	}
	return false
}

// GetByKey is a convenient shortcuts to get the value for a given tag key.
func (tags Tags) GetByKey(name string, key string) string {
	tag := tags.Get(name)
	if tag != nil {
		return tag.Get(key)
	}
	return ""
}

// TagsAnalyzerOption is a functional option to configure TagsAnalyzerOptions.
type TagsAnalyzerOption func(*TagsAnalyzerOptions)

// TagsAnalyzerOptions defines the tags analyzer behavior.
type TagsAnalyzerOptions struct {
	// Name defines the default tag name.
	Name string
	// Collector defines what tags should be analyzed.
	Collector []string
	// Mapper defines how to convert a supported tag to the default one.
	Mapper map[string]string
}

// GetTags returns field tags.
func GetTags(field reflect.StructField, args ...TagsAnalyzerOption) Tags {
	options := &TagsAnalyzerOptions{
		Name:      TagName,
		Collector: TagsList,
		Mapper:    TagsMapper,
	}

	for i := range args {
		args[i](options)
	}

	list := map[string]string{}

	for _, name := range options.Collector {
		_, ok := list[name]
		if !ok {
			v := field.Tag.Get(name)
			if len(v) != 0 {
				list[name] = v
			}
		}
	}

	tags := Tags{}
	for k, v := range list {
		splits := strings.Split(v, ",")

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

			// Typically the case of sqlx tag that doesn't have key:value format (db:"column_name")
			k, ok := options.Mapper[k]
			if ok {
				v := strings.TrimSpace(splits[0])
				if v != TagKeyIgnored {
					tags.Set(options.Name, TagProperty{
						key:   k,
						value: v,
					})
				} else {
					tags.Set(options.Name, TagProperty{
						key:   TagKeyIgnored,
						value: "true",
					})
				}
				continue
			}

			// Typically, we have single property like "default", "ignored", etc...
			// To be consistent, we add true/false string values.
			if length == 1 {
				tags.Set(options.Name, TagProperty{
					key:   strings.TrimSpace(splits[0]),
					value: "true",
				})
				continue
			}

			// Property named tag: key:value
			if length == 2 {
				tags.Set(options.Name, TagProperty{
					key:   strings.TrimSpace(splits[0]),
					value: strings.TrimSpace(splits[1]),
				})
				continue
			}
		}
	}

	return tags
}
