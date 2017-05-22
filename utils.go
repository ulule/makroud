package sqlxx

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
// Model
// ----------------------------------------------------------------------------

// ToModel converts the given instance to a Model instance.
func ToModel(itf interface{}) Model {
	typ, ok := itf.(reflect.Type)
	if ok {
		if typ.Kind() == reflect.Slice {
			typ = GetIndirectType(typ.Elem())
		} else {
			typ = GetIndirectType(typ)
		}

		model, ok := reflect.New(typ).Elem().Interface().(Model)
		if ok {
			return model
		}

		return nil
	}

	value := reflect.Indirect(reflect.ValueOf(itf))

	// Single instance
	if value.IsValid() && value.Kind() == reflect.Struct {
		return value.Interface().(Model)
	}

	// Slice of instances
	if value.Kind() == reflect.Slice {
		// Slice of pointers
		if value.Type().Elem().Kind() == reflect.Ptr {
			return reflect.New(value.Type().Elem().Elem()).Interface().(Model)
		}

		// Slice of values
		return reflect.New(value.Type().Elem()).Interface().(Model)
	}

	return reflect.New(value.Type()).Interface().(Model)
}

// ----------------------------------------------------------------------------
// Int64
// ----------------------------------------------------------------------------

var int64Type = reflect.TypeOf(int64(0))

// IntToInt64 converts given int to int64.
func IntToInt64(value interface{}) (int64, error) {
	if cast, ok := value.(int64); ok {
		return cast, nil
	}

	// sql.NullInt* support
	if valuer, ok := value.(driver.Valuer); ok {
		v, err := valuer.Value()
		if err != nil || v == nil {
			return 0, errors.Wrap(err, "cannot convert to int64")
		}

		value = v
	}

	reflected := reflect.Indirect(reflect.ValueOf(value))

	if !reflected.IsValid() {
		return 0, errors.Errorf("invalid value: %v", value)
	}

	if !reflected.Type().ConvertibleTo(int64Type) {
		return 0, errors.Errorf("unable to convert %v to int64", reflected.Type())
	}

	return reflected.Convert(int64Type).Int(), nil
}

// ----------------------------------------------------------------------------
// Field values
// ----------------------------------------------------------------------------

// GetFieldValue returns the value
func GetFieldValue(itf interface{}, name string) (interface{}, error) {
	value, ok := itf.(reflect.Value)
	if !ok {
		value = reflect.Indirect(reflect.ValueOf(itf))
	}

	if value.Kind() == reflect.Interface {
		value = reflect.ValueOf(value.Interface())
	}

	// Avoid calling FieldByName on ptr
	value = reflect.Indirect(value)

	// Avoid calling FieldByName on zero value
	if !value.IsValid() {
		return nil, nil
	}

	field := value.FieldByName(name)
	if !field.IsValid() {
		return nil, errors.Errorf("no such field %s in %+v", name, itf)
	}

	if field.Kind() == reflect.Ptr && field.IsNil() {
		return nil, nil
	}

	return field.Interface(), nil
}

// GetFieldValueInt64 returns int64 value for the given instance field.
func GetFieldValueInt64(instance interface{}, field string) (int64, error) {
	value, err := GetFieldValue(instance, field)
	if err != nil {
		return 0, err
	}

	converted, err := IntToInt64(value)
	if err != nil {
		return 0, err
	}

	return converted, nil
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

// SetFieldValue sets the provided value
func SetFieldValue(itf interface{}, name string, value interface{}) error {
	v, ok := itf.(reflect.Value)
	if !ok {
		v = reflect.Indirect(reflect.ValueOf(itf))
	}

	if v.Kind() == reflect.Interface {
		v = reflect.ValueOf(v.Interface())
	}

	field := v.FieldByName(name)
	if !field.IsValid() {
		return errors.Errorf("no such field %s in %+v", name, v.Interface())
	}

	if !field.CanSet() {
		return errors.Errorf("cannot set %s field on %v%+v", name, v.Type().Name(), v.Interface())
	}

	fv := reflect.Indirect(reflect.ValueOf(value))
	if !fv.IsValid() {
		return nil
	}

	if field.Type().Kind() == reflect.Ptr {
		fv = reflect.ValueOf(MakePointer(fv.Interface()))
	}

	if field.Type() != fv.Type() {
		return errors.Errorf("provided value type %v didn't match field type %v", fv.Type(), field.Type())
	}

	field.Set(fv)

	return nil
}

// ----------------------------------------------------------------------------
// Reflection
// ----------------------------------------------------------------------------

// IsSlice returns true if the given interface is a slice.
func IsSlice(itf interface{}) bool {
	return GetIndirectType(reflect.ValueOf(itf).Type()).Kind() == reflect.Slice
}

// IsZero returns true if the given interface is a zero value or nil.
func IsZero(itf interface{}) bool {
	if itf == nil {
		return true
	}

	value, ok := itf.(reflect.Value)
	if !ok {
		value = reflect.Indirect(reflect.ValueOf(itf))
	}

	if value.Kind() == reflect.Ptr && value.IsNil() {
		return true
	}

	zero := reflect.Zero(value.Type())
	return value.Interface() == zero.Interface()
}

// GetIndirectType returns indirect type for the given type.
func GetIndirectType(itf interface{}) reflect.Type {
	var (
		t  reflect.Type
		ok bool
	)

	for {
		t, ok = itf.(reflect.Type)
		if !ok {
			t = reflect.TypeOf(itf)
		}

		if t.Kind() != reflect.Ptr {
			break
		}

		itf = t.Elem()
	}

	return t
}

// MakePointer makes a copy of the given interface and returns a pointer.
func MakePointer(itf interface{}) interface{} {
	t := reflect.TypeOf(itf)

	cp := reflect.New(t)
	cp.Elem().Set(reflect.ValueOf(itf))

	// Avoid double pointers if itf is a pointer
	if t.Kind() == reflect.Ptr {
		return cp.Elem().Interface()
	}

	return cp.Interface()
}

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
	tag := t.Get(name)
	return tag != nil
}

// GetByKey is a convenient shortcuts to get the value for a given tag key.
func (t FieldTags) GetByKey(name string, key string) string {
	tag := t.Get(name)
	if tag != nil {
		return tag.Get(key)
	}
	return ""
}

// GetFieldTags returns field tags
func GetFieldTags(field reflect.StructField, tagNames []string, propertyMapping map[string]string) FieldTags {
	rawTags := map[string]string{}

	for _, name := range tagNames {
		_, ok := rawTags[name]
		if !ok {
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
