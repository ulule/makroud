package reflekt

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
)

// ----------------------------------------------------------------------------
// Indirect helpers
// ----------------------------------------------------------------------------

// GetIndirectValue returns the value that the interface v contains or that the pointer v points to.
func GetIndirectValue(itf interface{}) reflect.Value {
	v, ok := itf.(reflect.Value)
	if !ok {
		v = reflect.ValueOf(itf)
	}

	v = reflect.Indirect(v)

	if v.Kind() == reflect.Interface && !v.IsNil() {
		v = v.Elem()
	}

	return v
}

// GetIndirectType returns indirect type for the given type.
func GetIndirectType(itf interface{}) reflect.Type {
	t, ok := itf.(reflect.Type)
	if !ok {
		t = reflect.TypeOf(itf)
	}

	if t.Kind() == reflect.Ptr {
		return GetIndirectType(t.Elem())
	}

	return t
}

// ----------------------------------------------------------------------------
// Cloners
// ----------------------------------------------------------------------------

// MakeSlice takes a type and returns create a slice from.
func MakeSlice(itf interface{}) interface{} {
	sliceType := reflect.SliceOf(GetIndirectType(itf))
	slice := reflect.New(sliceType)
	slice.Elem().Set(reflect.MakeSlice(sliceType, 0, 0))
	return slice.Elem().Interface()
}

// CloneType returns a new interface from a given interface.
func CloneType(itf interface{}, args ...reflect.Kind) interface{} {
	var kind reflect.Kind
	if len(args) > 0 {
		kind = args[0]
	}

	if kind == reflect.Slice {
		return reflect.New(reflect.TypeOf(MakeSlice(itf))).Interface()
	}

	return reflect.New(reflect.TypeOf(itf)).Interface()
}

// Copy makes a copy of the given interface.
func Copy(itf interface{}) interface{} {
	cp := reflect.New(GetIndirectType(reflect.TypeOf(itf)))
	cp.Elem().Set(reflect.ValueOf(itf))
	return cp.Interface()
}

// ----------------------------------------------------------------------------
// Struct field tags
// ----------------------------------------------------------------------------

// GetFieldValuesInSlice returns values for the given field in slice of structs.
func GetFieldValuesInSlice(slc interface{}, field string) ([]interface{}, error) {
	var (
		value  = reflect.ValueOf(slc).Elem()
		slcLen = value.Len()
		values []interface{}
	)

	for i := 0; i < slcLen; i++ {
		v, err := GetFieldValue(value.Index(i).Interface(), field)
		if err != nil {
			return nil, err
		}
		values = append(values, v)
	}

	return values, nil
}

// GetFieldValues returns values for the given field for struct or slice.
func GetFieldValues(out interface{}, name string) ([]interface{}, error) {
	if IsSlice(out) {
		return GetFieldValuesInSlice(out, name)
	}

	v, err := GetFieldValue(out, name)
	if err != nil {
		return nil, err
	}

	return []interface{}{v}, nil
}

// GetFieldValue returns the value
func GetFieldValue(itf interface{}, name string) (interface{}, error) {
	value, ok := itf.(reflect.Value)
	if !ok {
		value = GetIndirectValue(itf)
	}

	field := value.FieldByName(name)
	if !field.IsValid() {
		return nil, fmt.Errorf("No such field %s in %+v", name, itf)
	}

	return field.Interface(), nil
}

// SetFieldValue sets the provided value
func SetFieldValue(itf interface{}, name string, value interface{}) error {
	v, ok := itf.(reflect.Value)
	if !ok {
		v = GetIndirectValue(itf)
	}

	field := v.FieldByName(name)
	if !field.IsValid() {
		return fmt.Errorf("no such field %s in %+v", name, v.Interface())
	}

	if !field.CanSet() {
		return fmt.Errorf("cannot set %s field on %v%+v", name, v.Type().Name(), v.Interface())
	}

	fv := GetIndirectValue(value)
	if field.Type().Kind() == reflect.Ptr {
		fv = reflect.ValueOf(Copy(fv.Interface()))
	}

	if field.Type() != fv.Type() {
		return fmt.Errorf("provided value type %v didn't match field type %v", fv.Type(), field.Type())
	}

	field.Set(fv)

	return nil
}

// ----------------------------------------------------------------------------
// Checkers
// ----------------------------------------------------------------------------

// IsNullableType returns true if the given type is a nullable one.
func IsNullableType(t reflect.Type) bool {
	return t.ConvertibleTo(reflect.TypeOf((*driver.Valuer)(nil)).Elem())
}

// IsSlice returns true if the given interface is a slice.
func IsSlice(itf interface{}) bool {
	return GetIndirectType(reflect.ValueOf(itf).Type()).Kind() == reflect.Slice
}

// IsZeroValue returns true if the given interface is a zero value.
func IsZeroValue(itf interface{}) bool {
	v := reflect.ValueOf(itf)

	// Avoid call of reflect.Value.Interface on zero Value
	if !v.IsValid() {
		return true
	}

	return reflect.Indirect(v).Interface() == reflect.Zero(reflect.Indirect(v).Type()).Interface()
}

// ----------------------------------------------------------------------------
// Struct tags
// ----------------------------------------------------------------------------

// FieldTagProperty is a struct tag property
type FieldTagProperty struct {
	Key   string
	Value string
}

// String returns instance string
func (t FieldTagProperty) String() string {
	return fmt.Sprintf("%s:%v", t.Key, t.Value)
}

// FieldTag is struct tag
type FieldTag struct {
	Name       string
	Properties []FieldTagProperty
}

// String returns instance string
func (t FieldTag) String() string {
	var props []string
	for _, p := range t.Properties {
		props = append(props, fmt.Sprintf("%s", p))
	}
	return fmt.Sprintf("%s -- %s", t.Name, strings.Join(props, ", "))
}

// Get returns value for the given property name.
func (t FieldTag) Get(key string) string {
	for _, p := range t.Properties {
		if p.Key == key {
			return p.Value
		}
	}
	return ""
}

// FieldTags a group of tag (usually for a struct field)
type FieldTags []FieldTag

func (t FieldTags) String() string {
	var tags []string
	for _, tag := range tags {
		tags = append(tags, fmt.Sprintf("%s", tag))
	}
	return strings.Join(tags, "\n")
}

// Get returns tag by name.
func (t FieldTags) Get(name string) *FieldTag {
	for _, tag := range t {
		if tag.Name == name {
			return &tag
		}
	}
	return nil
}

// Set sets the given tag into the slice.
func (t *FieldTags) Set(name string, property FieldTagProperty) {
	var found bool

	tags := *t
	for i := range tags {
		if tags[i].Name == name {
			tags[i].Properties = append(tags[i].Properties, property)
			found = true
		}
	}

	if !found {
		tags = append(tags, FieldTag{
			Name:       name,
			Properties: []FieldTagProperty{property},
		})
	}

	*t = tags
}

// HasKey is a convenient shortcuts to check if a key is present.
func (t FieldTags) HasKey(name string, key string) bool {
	if tag := t.Get(name); tag != nil {
		return true
	}
	return false
}

// GetByKey is a convenient shortcuts to get the value for a given tag key.
func (t FieldTags) GetByKey(name string, key string) string {
	if tag := t.Get(name); tag != nil {
		return tag.Get(key)
	}
	return ""
}

// GetFieldTags returns field tags
func GetFieldTags(field reflect.StructField, tagNames []string, propertyMapping map[string]string) FieldTags {
	rawTags := map[string]string{}

	for _, name := range tagNames {
		if _, ok := rawTags[name]; !ok {
			v := field.Tag.Get(name)
			if len(v) != 0 {
				rawTags[name] = v
			}
		}
	}

	tags := FieldTags{}

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
				tags.Set(k, FieldTagProperty{
					Key:   propertyKey,
					Value: strings.TrimSpace(splits[0]),
				})
				continue
			}

			// Typically, we have single property like "default", "ignored", etc.
			// To be consistent, we add true/false string values.
			if length == 1 {
				tags.Set(k, FieldTagProperty{
					Key:   strings.TrimSpace(splits[0]),
					Value: "true",
				})
				continue
			}

			// Property named tag: key:value
			if length == 2 {
				tags.Set(k, FieldTagProperty{
					Key:   strings.TrimSpace(splits[0]),
					Value: strings.TrimSpace(splits[1]),
				})
			}
		}
	}

	return tags
}
