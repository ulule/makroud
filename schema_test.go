package sqlxx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fieldResultTest struct {
	field    string
	table    string
	name     string
	prefixed string
}

type relationResultTest struct {
	fieldName    string
	name         string
	prefixedName string
	tableName    string
	relationType RelationType
	isReference  bool
}

func TestGetSchema(t *testing.T) {
	is := assert.New(t)

	schema, err := GetSchema(StructWithoutTags{})
	is.NoError(err)

	testFields(t, schema, []fieldResultTest{
		{"ID", "foo", "id", "foo.id"},
		{"FirstName", "foo", "first_name", "foo.first_name"},
		{"LastName", "foo", "last_name", "foo.last_name"},
		{"ThisIsAVeryLongFieldName123", "foo", "this_is_a_very_long_field_name123", "foo.this_is_a_very_long_field_name123"},
	})

	testRelations(t, schema, []relationResultTest{
		{"RelatedModel", "related_model_id", "foo.related_model_id", "foo", RelationTypeOneToMany, false},
		{"RelatedModelPtr", "related_model_ptr_id", "foo.related_model_ptr_id", "foo", RelationTypeOneToMany, false},
		{"RelatedModel", "custom_id", "related.custom_id", "related", RelationTypeOneToMany, true},
		{"RelatedModelPtr", "custom_id", "related.custom_id", "related", RelationTypeOneToMany, true},
		{"RelatedSlice", "custom_id", "related.custom_id", "related", RelationTypeManyToOne, true},
		{"RelatedSlicePtr", "custom_id", "related.custom_id", "related", RelationTypeManyToOne, true},
		{"RelatedPtrSlice", "custom_id", "related.custom_id", "related", RelationTypeManyToOne, true},
	})

	schema, err = GetSchema(StructWithTags{})
	is.NoError(err)

	testFields(t, schema, []fieldResultTest{
		{"ID", "foo", "public_id", "foo.public_id"},
		{"FirstName", "foo", "firstname", "foo.firstname"},
		{"LastName", "foo", "last_name", "foo.last_name"},
		{"ThisIsAVeryLongFieldName123", "foo", "short_field", "foo.short_field"},
	})

	testRelations(t, schema, []relationResultTest{
		{"RelatedModel", "member_id", "foo.member_id", "foo", RelationTypeOneToMany, false},
		{"RelatedModelPtr", "member_id", "foo.member_id", "foo", RelationTypeOneToMany, false},
		{"RelatedModel", "custom_id", "related.custom_id", "related", RelationTypeOneToMany, true},
		{"RelatedModelPtr", "custom_id", "related.custom_id", "related", RelationTypeOneToMany, true},
		{"RelatedSlice", "custom_id", "related.custom_id", "related", RelationTypeManyToOne, true},
		{"RelatedSlicePtr", "custom_id", "related.custom_id", "related", RelationTypeManyToOne, true},
		{"RelatedPtrSlice", "custom_id", "related.custom_id", "related", RelationTypeManyToOne, true},
	})
}

func testFields(t *testing.T, schema *Schema, results []fieldResultTest) {
	is := assert.New(t)

	for _, r := range results {
		is.Equal(r.table, schema.Fields[r.field].TableName)
		is.Equal(r.name, schema.Fields[r.field].Name)
		is.Equal(r.prefixed, schema.Fields[r.field].PrefixedName())
	}
}

func testRelations(t *testing.T, schema *Schema, results []relationResultTest) {
	is := assert.New(t)

	for _, r := range results {
		relation := schema.Relations[r.fieldName]

		field := relation.FK
		if r.isReference {
			field = relation.FKReference
		}

		is.Equal(r.name, field.Name)
		is.Equal(r.prefixedName, field.PrefixedName())
		is.Equal(r.tableName, field.TableName)
		is.Equal(r.relationType, relation.Type)
	}
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
	RelatedSlice                []RelatedModel
	RelatedSlicePtr             *[]RelatedModel
	RelatedPtrSlice             []*RelatedModel
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
	RelatedSlice                []RelatedModel
	RelatedSlicePtr             *[]RelatedModel
	RelatedPtrSlice             []*RelatedModel
}

func (StructWithTags) TableName() string {
	return "foo"
}
