package sqlxx

import (
	"reflect"
	"sort"
	"strings"

	"github.com/ulule/sqlxx/reflekt"
)

// ----------------------------------------------------------------------------
// Model
// ----------------------------------------------------------------------------

// Model represents a database table.
type Model interface {
	TableName() string
}

// GetModelFromInterface returns interface as a Model interface.
func GetModelFromInterface(itf interface{}) Model {
	value := reflekt.GetIndirectValue(itf)

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

// GetModelFromType returns model type.
func GetModelFromType(typ reflect.Type) Model {
	if typ.Kind() == reflect.Slice {
		typ = reflekt.GetIndirectType(typ.Elem())
	} else {
		typ = reflekt.GetIndirectType(typ)
	}

	if model, isModel := reflect.New(typ).Elem().Interface().(Model); isModel {
		return model
	}

	return nil
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
