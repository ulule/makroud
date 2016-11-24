package sqlxx

import (
	"database/sql"
	"testing"
	"time"

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
		is.Equal(schema.Columns[r.field].TableName, r.table)
		is.Equal(schema.Columns[r.field].Name, r.name)
		is.Equal(schema.Columns[r.field].PrefixedName, r.prefixed)
	}

	is.Equal(schema.Associations["RelatedModel"].FK.PrefixedName, "foo.related_model_id")
	is.Equal(schema.Associations["RelatedModel"].FKReference.PrefixedName, "related.id")
	is.Equal(schema.Associations["RelatedModelPtr"].FK.PrefixedName, "foo.related_model_ptr_id")
	is.Equal(schema.Associations["RelatedModelPtr"].FKReference.PrefixedName, "related.id")

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
		is.Equal(schema.Columns[r.field].TableName, r.table)
		is.Equal(schema.Columns[r.field].Name, r.name)
		is.Equal(schema.Columns[r.field].PrefixedName, r.prefixed)
	}

	is.Equal(schema.Associations["RelatedModel"].FK.PrefixedName, "foo.member_id")
	is.Equal(schema.Associations["RelatedModel"].FKReference.PrefixedName, "related.custom_id")
	is.Equal(schema.Associations["RelatedModelPtr"].FK.PrefixedName, "foo.member_id")
	is.Equal(schema.Associations["RelatedModelPtr"].FKReference.PrefixedName, "related.custom_id")
}

func TestIsModel(t *testing.T) {
	is := assert.New(t)
	is.True(isModel(RelatedModel{}))
	is.True(isModel(&RelatedModel{}))
	is.True(isModel(User{}))
	is.True(isModel(&User{}))
	is.False(isModel(struct{ ID int }{1}))
	is.False(isModel(time.Time{}))
	is.False(isModel(8))
	is.False(isModel("hello"))
	is.False(isModel(sql.NullInt64{}))
}

// ----------------------------------------------------------------------------
// Test data
// ----------------------------------------------------------------------------

type RelatedModel struct {
	ID   int
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
	RelatedModel                RelatedModel  `db:"member_id" sqlxx:"related:custom_id"`
	RelatedModelPtr             *RelatedModel `db:"member_id" sqlxx:"related:custom_id"`
}

func (StructWithTags) TableName() string {
	return "foo"
}
