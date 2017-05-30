package sqlxx_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
	"github.com/ulule/sqlxx/reflectx"
)

func TestUtils_MakePointer(t *testing.T) {
	is := require.New(t)

	type embedType struct {
		value int
	}

	type anyType struct {
		value    int
		embed    embedType
		embedPtr *embedType
	}

	any := anyType{value: 1}
	anyPtr := &anyType{value: 1}

	results := []interface{}{
		sqlxx.MakePointer(any),
		sqlxx.MakePointer(anyPtr),
	}

	for _, r := range results {
		is.Equal(1, r.(*anyType).value)
		is.Equal(reflect.ValueOf(r).Kind(), reflect.Ptr)
		is.Equal(reflect.ValueOf(r).Type().Elem(), reflect.TypeOf(anyType{}))
	}

	anyWithEmbed := anyType{value: 1, embed: embedType{value: 2}}
	anyWithEmbedPtr := anyType{value: 1, embedPtr: &embedType{value: 2}}

	results = []interface{}{
		sqlxx.MakePointer(anyWithEmbed.embed),
		sqlxx.MakePointer(anyWithEmbedPtr.embedPtr),
	}

	for _, r := range results {
		is.Equal(2, r.(*embedType).value)
		is.Equal(reflect.ValueOf(r).Kind(), reflect.Ptr)
		is.Equal(reflect.ValueOf(r).Type().Elem(), reflect.TypeOf(embedType{}))
	}
}

func TestUtils_IsZero(t *testing.T) {
	is := require.New(t)

	type user struct {
		Name   *string
		Fk     sql.NullInt64
		FkPtr  *sql.NullInt64
		Active bool
	}

	name := "thoas"
	empty := ""
	scenarios := []struct {
		value    user
		field    string
		expected bool
	}{
		{
			value:    user{},
			field:    "Name",
			expected: true,
		},
		{
			value:    user{Name: &empty},
			field:    "Name",
			expected: true,
		},
		{
			value:    user{Name: &name},
			field:    "Name",
			expected: false,
		},
		{
			value:    user{},
			field:    "FkPtr",
			expected: true,
		},
		{
			value:    user{FkPtr: &sql.NullInt64{}},
			field:    "FkPtr",
			expected: true,
		},
		{
			value:    user{FkPtr: &sql.NullInt64{Valid: true, Int64: 64}},
			field:    "FkPtr",
			expected: false,
		},
		{
			value:    user{Fk: sql.NullInt64{}},
			field:    "Fk",
			expected: true,
		},
		{
			value:    user{Fk: sql.NullInt64{Valid: true}},
			field:    "Fk",
			expected: false,
		},
		{
			value:    user{},
			field:    "Active",
			expected: true,
		},
		{
			value:    user{Active: false},
			field:    "Active",
			expected: true,
		},
		{
			value:    user{Active: true},
			field:    "Active",
			expected: false,
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #%d", (i + 1))

		field, err := reflectx.GetFieldValue(scenario.value, scenario.field)
		is.NoError(err)

		isZero := reflectx.IsZero(field)
		is.Equal(scenario.expected, isZero, message)
	}

}
