package sqlxx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fieldResultTest struct {
	field      string
	columnName string
	columnPath string
}

type relationResultTest struct {
	fieldName     string
	fkColumnName  string
	fkColumnPath  string
	refColumnName string
	refColumnPath string
	relationType  RelationType
}

func TestGetSchemaInfiniteLoop(t *testing.T) {
	is := assert.New(t)

	schema, err := GetSchema(User{})
	is.NoError(err)
	is.NotNil(schema)
}

func TestGetSchema(t *testing.T) {
	is := assert.New(t)

	schema, err := GetSchema(Untagged{})
	is.NoError(err)

	// Check unexported
	is.NotContains(schema.FieldNames(), "unexportedField")

	// Check db excluded
	is.NotContains(schema.FieldNames(), "DBExcludedField")

	testFields(t, schema, []fieldResultTest{
		{
			"ID",
			"id", "untagged.id",
		},
		{
			"FirstName",
			"first_name", "untagged.first_name",
		},
		{
			"LastName",
			"last_name", "untagged.last_name",
		},
		{
			"ThisIsAVeryLongFieldName123",
			"this_is_a_very_long_field_name123", "untagged.this_is_a_very_long_field_name123",
		},
	})

	testRelations(t, schema, []relationResultTest{
		{
			"RelatedModel",
			"related_model_id", "untagged.related_model_id",
			"custom_id", "related.custom_id",
			RelationTypeOneToMany,
		},
		{
			"RelatedModelPtr",
			"related_model_ptr_id", "untagged.related_model_ptr_id",
			"custom_id", "related.custom_id",
			RelationTypeOneToMany,
		},
		{
			"ManyModel",
			"untagged_id", "many.untagged_id",
			"id", "untagged.id",
			RelationTypeManyToOne,
		},
		{
			"ManyModelPtr",
			"untagged_id", "many.untagged_id",
			"id", "untagged.id",
			RelationTypeManyToOne,
		},
		{
			"ManyModelPtrs",
			"untagged_id", "many.untagged_id",
			"id", "untagged.id",
			RelationTypeManyToOne,
		},
	})

	cache.Flush()

	schema, err = GetSchema(Tagged{})
	is.NoError(err)

	// Check unexported
	is.NotContains(schema.FieldNames(), "unexportedField")

	// Check db excluded
	is.NotContains(schema.FieldNames(), "DBExcludedField")

	testFields(t, schema, []fieldResultTest{
		{
			"ID",
			"public_id", "tagged.public_id",
		},
		{
			"FirstName",
			"firstname", "tagged.firstname",
		},
		{
			"LastName",
			"last_name", "tagged.last_name",
		},
		{
			"ThisIsAVeryLongFieldName123",
			"short_field", "tagged.short_field",
		},
	})

	testRelations(t, schema, []relationResultTest{
		{
			"RelatedModel",
			"member_id", "tagged.member_id",
			"custom_id", "related.custom_id",
			RelationTypeOneToMany,
		},
		{
			"RelatedModelPtr",
			"member_id", "tagged.member_id",
			"custom_id", "related.custom_id",
			RelationTypeOneToMany,
		},
		{
			"ManyModel",
			"tagged_id", "many.tagged_id",
			"id", "tagged.id",
			RelationTypeManyToOne,
		},
		{
			"ManyModelPtr",
			"tagged_id", "many.tagged_id",
			"id", "tagged.id",
			RelationTypeManyToOne,
		},
		{
			"ManyModelPtrs",
			"tagged_id", "many.tagged_id",
			"id", "tagged.id",
			RelationTypeManyToOne,
		},
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
		is.Equal(r.columnName, schema.Fields[r.field].ColumnName)
		is.Equal(r.columnPath, schema.Fields[r.field].ColumnPath())
	}
}

func testRelations(t *testing.T, schema Schema, results []relationResultTest) {
	is := assert.New(t)

	for _, r := range results {
		relation := schema.Relations[r.fieldName]

		fk := relation.FK
		ref := relation.Reference

		is.Equal(r.fieldName, relation.Name)
		is.Equal(r.fkColumnName, fk.ColumnName)
		is.Equal(r.fkColumnPath, fk.ColumnPath())
		is.Equal(r.refColumnName, ref.ColumnName)
		is.Equal(r.refColumnPath, ref.ColumnPath())
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

type Untagged struct {
	ID                          int
	FirstName                   string
	LastName                    string
	ThisIsAVeryLongFieldName123 string

	// Skip it
	unexportedField string

	// Skip db:"-"
	DBExcludedField int `db:"-"`

	RelatedModel    RelatedModel
	RelatedModelPtr *RelatedModel

	ManyModel     []ManyModel
	ManyModelPtr  *[]ManyModel
	ManyModelPtrs []*ManyModel
}

func (Untagged) TableName() string {
	return "untagged"
}

type Tagged struct {
	ID                          int    `db:"public_id"`
	FirstName                   string `db:"firstname"`
	LastName                    string
	ThisIsAVeryLongFieldName123 string        `db:"short_field"`
	RelatedModel                RelatedModel  `db:"member_id"`
	RelatedModelPtr             *RelatedModel `db:"member_id"`

	// Skip it
	unexportedField string

	// Skip db:"-"
	DBExcludedField int `db:"-"`

	ManyModel     []ManyModel
	ManyModelPtr  *[]ManyModel
	ManyModelPtrs []*ManyModel
}

func (Tagged) TableName() string {
	return "tagged"
}
