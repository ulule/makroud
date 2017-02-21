package sqlxx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ulule/sqlxx"
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
	relationType  sqlxx.RelationType
}

func TestGetSchemaInfiniteLoop(t *testing.T) {
	schema, err := sqlxx.GetSchema(User{})
	assert.Nil(t, err)
	assert.NotNil(t, schema)
}

func TestGetSchema(t *testing.T) {
	schema, err := sqlxx.GetSchema(Untagged{})
	assert.Nil(t, err)

	// Check unexported
	assert.NotContains(t, schema.FieldNames(), "unexportedField")

	// Check db excluded
	assert.NotContains(t, schema.FieldNames(), "DBExcludedField")

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
			sqlxx.RelationTypeOneToMany,
		},
		{
			"RelatedModelPtr",
			"related_model_ptr_id", "untagged.related_model_ptr_id",
			"custom_id", "related.custom_id",
			sqlxx.RelationTypeOneToMany,
		},
		{
			"ManyModel",
			"untagged_id", "many.untagged_id",
			"id", "untagged.id",
			sqlxx.RelationTypeManyToOne,
		},
		{
			"ManyModelPtr",
			"untagged_id", "many.untagged_id",
			"id", "untagged.id",
			sqlxx.RelationTypeManyToOne,
		},
		{
			"ManyModelPtrs",
			"untagged_id", "many.untagged_id",
			"id", "untagged.id",
			sqlxx.RelationTypeManyToOne,
		},
	})

	sqlxx.GetCache().Flush()

	schema, err = sqlxx.GetSchema(Tagged{})
	assert.Nil(t, err)

	// Check unexported
	assert.NotContains(t, schema.FieldNames(), "unexportedField")

	// Check db excluded
	assert.NotContains(t, schema.FieldNames(), "DBExcludedField")

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
			sqlxx.RelationTypeOneToMany,
		},
		{
			"RelatedModelPtr",
			"member_id", "tagged.member_id",
			"custom_id", "related.custom_id",
			sqlxx.RelationTypeOneToMany,
		},
		{
			"ManyModel",
			"tagged_id", "many.tagged_id",
			"id", "tagged.id",
			sqlxx.RelationTypeManyToOne,
		},
		{
			"ManyModelPtr",
			"tagged_id", "many.tagged_id",
			"id", "tagged.id",
			sqlxx.RelationTypeManyToOne,
		},
		{
			"ManyModelPtrs",
			"tagged_id", "many.tagged_id",
			"id", "tagged.id",
			sqlxx.RelationTypeManyToOne,
		},
	})
}

func TestSchemaRelationPaths(t *testing.T) {
	schema, err := sqlxx.GetSchema(Article{})
	assert.Nil(t, err)

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
		assert.True(t, ok)
		assert.Equal(t, r.modelName, relation.Schema.ModelName)
		assert.Equal(t, r.tableName, relation.Schema.TableName)
		assert.Equal(t, r.name, relation.Name)

	}
}

func TestGetSchema_PrimaryKeyField(t *testing.T) {
	// Implicit
	schema, err := sqlxx.GetSchema(ImplicitPrimaryKey{})
	assert.Nil(t, err)
	assert.Equal(t, schema.PrimaryField.Name, "ID")
	assert.Equal(t, schema.PrimaryField.ColumnName, "id")

	// Explicit
	schema, err = sqlxx.GetSchema(ExplicitPrimaryKey{})
	assert.Nil(t, err)
	assert.Equal(t, schema.PrimaryField.Name, "TadaID")
	assert.Equal(t, schema.PrimaryField.ColumnName, "tada_id")
}

func testFields(t *testing.T, schema sqlxx.Schema, results []fieldResultTest) {
	for _, r := range results {
		assert.Equal(t, r.columnName, schema.Fields[r.field].ColumnName)
		assert.Equal(t, r.columnPath, schema.Fields[r.field].ColumnPath())
	}
}

func testRelations(t *testing.T, schema sqlxx.Schema, results []relationResultTest) {
	for _, r := range results {
		var (
			relation = schema.Relations[r.fieldName]
			fk       = relation.FK
			ref      = relation.Reference
		)

		assert.Equal(t, r.fieldName, relation.Name)
		assert.Equal(t, r.fkColumnName, fk.ColumnName)
		assert.Equal(t, r.fkColumnPath, fk.ColumnPath())
		assert.Equal(t, r.refColumnName, ref.ColumnName)
		assert.Equal(t, r.refColumnPath, ref.ColumnPath())
		assert.Equal(t, r.relationType, relation.Type)
	}
}

// ----------------------------------------------------------------------------
// Test data
// ----------------------------------------------------------------------------

type ManyModel struct {
	ID   int `sqlxx:"primary_key:true"`
	Name string
}

func (ManyModel) TableName() string { return "many" }

type RelatedModel struct {
	ID   int `db:"custom_id" sqlxx:"primary_key:true"`
	Name string
}

func (RelatedModel) TableName() string { return "related" }

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

func (Untagged) TableName() string { return "untagged" }

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

func (Tagged) TableName() string { return "tagged" }

type ImplicitPrimaryKey struct {
	ID int
}

func (ImplicitPrimaryKey) TableName() string { return "implicitprimarykey" }

type ExplicitPrimaryKey struct {
	TadaID int `sqlxx:"primary_key"`
}

func (ExplicitPrimaryKey) TableName() string { return "explicitprimarykey" }
