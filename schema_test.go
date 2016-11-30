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
	isMany       bool
}

func TestGetSchema(t *testing.T) {
	is := assert.New(t)

	schema, err := GetSchema(StructWithoutTags{})
	is.NoError(err)

	testFields(t, schema, []fieldResultTest{
		{"ID", "untagged", "id", "untagged.id"},
		{"FirstName", "untagged", "first_name", "untagged.first_name"},
		{"LastName", "untagged", "last_name", "untagged.last_name"},
		{"ThisIsAVeryLongFieldName123", "untagged", "this_is_a_very_long_field_name123", "untagged.this_is_a_very_long_field_name123"},
	})

	testRelations(t, schema, []relationResultTest{
		{"RelatedModel", "related_model_id", "untagged.related_model_id", "untagged", RelationTypeOneToMany, false, false},
		{"RelatedModelPtr", "related_model_ptr_id", "untagged.related_model_ptr_id", "untagged", RelationTypeOneToMany, false, false},
		{"RelatedModel", "custom_id", "related.custom_id", "related", RelationTypeOneToMany, true, false},
		{"RelatedModelPtr", "custom_id", "related.custom_id", "related", RelationTypeOneToMany, true, false},
		{"ManyModel", "id", "many.id", "many", RelationTypeManyToOne, true, true},
		{"ManyModelPtr", "id", "many.id", "many", RelationTypeManyToOne, true, true},
		{"ManyModelPtrs", "id", "many.id", "many", RelationTypeManyToOne, true, true},
	})

	cache.Flush()

	schema, err = GetSchema(StructWithTags{})
	is.NoError(err)

	testFields(t, schema, []fieldResultTest{
		{"ID", "tagged", "public_id", "tagged.public_id"},
		{"FirstName", "tagged", "firstname", "tagged.firstname"},
		{"LastName", "tagged", "last_name", "tagged.last_name"},
		{"ThisIsAVeryLongFieldName123", "tagged", "short_field", "tagged.short_field"},
	})

	testRelations(t, schema, []relationResultTest{
		{"RelatedModel", "member_id", "tagged.member_id", "tagged", RelationTypeOneToMany, false, false},
		{"RelatedModelPtr", "member_id", "tagged.member_id", "tagged", RelationTypeOneToMany, false, false},
		{"RelatedModel", "custom_id", "related.custom_id", "related", RelationTypeOneToMany, true, false},
		{"RelatedModelPtr", "custom_id", "related.custom_id", "related", RelationTypeOneToMany, true, false},
		{"ManyModel", "id", "many.id", "many", RelationTypeManyToOne, true, true},
		{"ManyModelPtr", "id", "many.id", "many", RelationTypeManyToOne, true, true},
		{"ManyModelPtrs", "id", "many.id", "many", RelationTypeManyToOne, true, true},
	})
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

			if !r.isMany {
				is.IsType(RelatedModel{}, relation.Model)
			} else {
				is.IsType(ManyModel{}, relation.Model)
			}
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

type ManyModel struct {
	ID   int `sqlxx:"primary_key:true"`
	Name string
}

func (ManyModel) TableName() string {
	return "many"
}

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

	RelatedModel    RelatedModel
	RelatedModelPtr *RelatedModel

	ManyModel     []ManyModel
	ManyModelPtr  *[]ManyModel
	ManyModelPtrs []*ManyModel
}

func (StructWithoutTags) TableName() string {
	return "untagged"
}

type StructWithTags struct {
	ID                          int    `db:"public_id"`
	FirstName                   string `db:"firstname"`
	LastName                    string
	ThisIsAVeryLongFieldName123 string        `db:"short_field"`
	RelatedModel                RelatedModel  `db:"member_id"`
	RelatedModelPtr             *RelatedModel `db:"member_id"`

	ManyModel     []ManyModel
	ManyModelPtr  *[]ManyModel
	ManyModelPtrs []*ManyModel
}

func (StructWithTags) TableName() string {
	return "tagged"
}
