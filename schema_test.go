package sqlxx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fieldResultTest struct {
	field      string
	tableName  string
	columnName string
	columnPath string
}

type relationResultTest struct {
	fieldName    string
	columnName   string
	columnPath   string
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

	schema, err = GetSchema(Article{})
	is.NoError(err)
	schema.RelationPaths()
}

func TestSchemaRelationPaths(t *testing.T) {
	is := assert.New(t)

	schema, err := GetSchema(Article{})
	is.NoError(err)

	relations := schema.RelationPaths()

	results := []struct {
		path      string
		modelName string
		tableName string
		name      string
	}{
		{"Author", "User", "users", "Author"},
		{"Author.Avatars", "Avatar", "avatars", "Avatars"},
	}

	for _, r := range results {
		relation, ok := relations[r.path]
		is.True(ok)
		is.Equal(r.modelName, relation.Schema.ModelName)
		is.Equal(r.tableName, relation.Schema.TableName)
		is.Equal(r.name, relation.Name)

	}
}

func testFields(t *testing.T, schema Schema, results []fieldResultTest) {
	is := assert.New(t)

	for _, r := range results {
		is.Equal(r.tableName, schema.Fields[r.field].TableName)
		is.Equal(r.columnName, schema.Fields[r.field].ColumnName)
		is.Equal(r.columnPath, schema.Fields[r.field].ColumnPath())
	}
}

func testRelations(t *testing.T, schema Schema, results []relationResultTest) {
	is := assert.New(t)

	for _, r := range results {
		relation := schema.Relations[r.fieldName]

		field := relation.FK
		if r.isReference {
			field = relation.Reference
			is.IsType(RelatedModel{}, relation.Model)
		}

		is.Equal(r.fieldName, relation.Name)
		is.Equal(r.columnName, field.ColumnName)
		is.Equal(r.columnPath, field.ColumnPath())
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
