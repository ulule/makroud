package sqlxx

import (
	"reflect"
	"sort"
	"strings"
)

// Model represents a database table.
type Model interface {
	TableName() string
}

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

// GetModelName returns name of the given model.
func GetModelName(model Model) string {
	return reflect.Indirect(reflect.ValueOf(model)).Type().Name()
}

// ----------------------------------------------------------------------------
// Columns
// ----------------------------------------------------------------------------

// Columns is a list of table columns.
type Columns []string

// Returns string representation of slice.
func (c Columns) String() string {
	sort.Strings(c)
	return strings.Join(c, ", ")
}

// ----------------------------------------------------------------------------
// Where clauses
// ----------------------------------------------------------------------------

// Conditions is a list of query conditions
type Conditions []string

// String returns conditions as AND query.
func (c Conditions) String() string {
	return strings.Join(c, " AND ")
}

// OR returns conditions as OR query.
func (c Conditions) OR() string {
	return strings.Join(c, " OR ")
}
