package sqlxx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSchema(t *testing.T) {
	is := assert.New(t)

	schema, err := GetSchema(StructWithoutTags{})
	is.NoError(err)

	var results = []struct {
		field    string
		table    string
		name     string
		prefixed string
	}{
		{"ID", "foo", "id", "foo.id"},
		{"FirstName", "foo", "first_name", "foo.first_name"},
		{"LastName", "foo", "last_name", "foo.last_name"},
		{"ThisIsAVeryLongFieldName123", "foo", "this_is_a_very_long_field_name123", "foo.this_is_a_very_long_field_name123"},
	}

	for _, r := range results {
		is.Equal(r.table, schema.Fields[r.field].TableName)
		is.Equal(r.name, schema.Fields[r.field].Name)
		is.Equal(r.prefixed, schema.Fields[r.field].PrefixedName())
	}

	is.Equal("foo.related_model_id", schema.Relations["RelatedModel"].FK.PrefixedName())
	is.Equal("related.custom_id", schema.Relations["RelatedModel"].FKReference.PrefixedName())
	is.Equal("foo.related_model_ptr_id", schema.Relations["RelatedModelPtr"].FK.PrefixedName())
	is.Equal("related.custom_id", schema.Relations["RelatedModelPtr"].FKReference.PrefixedName())

	schema, err = GetSchema(StructWithTags{})
	is.NoError(err)

	results = []struct {
		field    string
		table    string
		name     string
		prefixed string
	}{
		{"ID", "foo", "public_id", "foo.public_id"},
		{"FirstName", "foo", "firstname", "foo.firstname"},
		{"LastName", "foo", "last_name", "foo.last_name"},
		{"ThisIsAVeryLongFieldName123", "foo", "short_field", "foo.short_field"},
	}

	for _, r := range results {
		is.Equal(r.table, schema.Fields[r.field].TableName)
		is.Equal(r.name, schema.Fields[r.field].Name)
		is.Equal(r.prefixed, schema.Fields[r.field].PrefixedName())
	}

	is.Equal("foo.member_id", schema.Relations["RelatedModel"].FK.PrefixedName())
	is.Equal("related.custom_id", schema.Relations["RelatedModel"].FKReference.PrefixedName())
	is.Equal("foo.member_id", schema.Relations["RelatedModelPtr"].FK.PrefixedName())
	is.Equal("related.custom_id", schema.Relations["RelatedModelPtr"].FKReference.PrefixedName())
}

// ----------------------------------------------------------------------------
// Test data
// ----------------------------------------------------------------------------

type RelatedModel struct {
	ID   int `db:"custom_id" sqlxx:"primary_key:true"`
	Name string
}

func (RelatedModel) TableName() string {
	return "related"
}

type StructWithoutTags struct {
	ID                          int
	FirstName                   string
	LastName                    string
	ThisIsAVeryLongFieldName123 string
	RelatedModel                RelatedModel
	RelatedModelPtr             *RelatedModel
}

func (StructWithoutTags) TableName() string {
	return "foo"
}

type StructWithTags struct {
	ID                          int    `db:"public_id"`
	FirstName                   string `db:"firstname"`
	LastName                    string
	ThisIsAVeryLongFieldName123 string        `db:"short_field"`
	RelatedModel                RelatedModel  `db:"member_id"`
	RelatedModelPtr             *RelatedModel `db:"member_id"`
}

func (StructWithTags) TableName() string {
	return "foo"
}
