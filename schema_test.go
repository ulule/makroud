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

type associationResultTest struct {
	field                 string
	columnName            string
	columnPath            string
	associationColumnName string
	associationColumnPath string
	associationType       sqlxx.AssociationType
}

func TestGetSchema(t *testing.T) {
	schema, err := sqlxx.GetSchema(Untagged{})
	assert.Nil(t, err)

	// Check unexported
	assert.NotContains(t, schema.FieldNames(), "unexportedField")

	// Check db excluded
	assert.NotContains(t, schema.FieldNames(), "DBExcludedField")

	testFields(t, schema, []fieldResultTest{
		{"ID", "id", "untagged.id"},
		{"FirstName", "first_name", "untagged.first_name"},
		{"LastName", "last_name", "untagged.last_name"},
		{"ThisIsAVeryLongFieldName123", "this_is_a_very_long_field_name123", "untagged.this_is_a_very_long_field_name123"},
	})

	testRelations(t, schema, []associationResultTest{
		{"RelatedModel", "related_model_id", "untagged.related_model_id", "custom_id", "related.custom_id", sqlxx.AssociationTypeOne},
		{"RelatedModelPtr", "related_model_ptr_id", "untagged.related_model_ptr_id", "custom_id", "related.custom_id", sqlxx.AssociationTypeOne},
		{"ManyModel", "untagged_id", "many.untagged_id", "id", "untagged.id", sqlxx.AssociationTypeMany},
		{"ManyModelPtr", "untagged_id", "many.untagged_id", "id", "untagged.id", sqlxx.AssociationTypeMany},
		{"ManyModelPtrs", "untagged_id", "many.untagged_id", "id", "untagged.id", sqlxx.AssociationTypeMany},
	})

	schema, err = sqlxx.GetSchema(Tagged{})
	assert.Nil(t, err)

	// Check unexported
	assert.NotContains(t, schema.FieldNames(), "unexportedField")

	// Check db excluded
	assert.NotContains(t, schema.FieldNames(), "DBExcludedField")

	testFields(t, schema, []fieldResultTest{
		{"ID", "public_id", "tagged.public_id"},
		{"FirstName", "firstname", "tagged.firstname"},
		{"LastName", "last_name", "tagged.last_name"},
		{"ThisIsAVeryLongFieldName123", "short_field", "tagged.short_field"},
	})

	testRelations(t, schema, []associationResultTest{
		{"RelatedModel", "member_id", "tagged.member_id", "custom_id", "related.custom_id", sqlxx.AssociationTypeOne},
		{"RelatedModelPtr", "member_id", "tagged.member_id", "custom_id", "related.custom_id", sqlxx.AssociationTypeOne},
		{"ManyModel", "tagged_id", "many.tagged_id", "id", "tagged.id", sqlxx.AssociationTypeMany},
		{"ManyModelPtr", "tagged_id", "many.tagged_id", "id", "tagged.id", sqlxx.AssociationTypeMany},
		{"ManyModelPtrs", "tagged_id", "many.tagged_id", "id", "tagged.id", sqlxx.AssociationTypeMany},
	})
}

func TestGetSchema_InfiniteLoop(t *testing.T) {
	schema, err := sqlxx.GetSchema(User{})
	assert.Nil(t, err)
	assert.NotNil(t, schema)
}

func TestSchema_AssociationsByPath(t *testing.T) {
	schema, err := sqlxx.GetSchema(Article{})
	assert.Nil(t, err)

	fields, err := schema.AssociationsByPath()
	assert.Nil(t, err)

	results := []struct {
		path      string
		modelName string
		tableName string
		name      string
	}{
		{"Author", "User", "users", "ID"},
		{"Author.Avatars", "Avatar", "avatars", "ID"},
	}

	for _, tt := range results {
		f, ok := fields[tt.path]
		assert.True(t, ok, tt.path)
		assert.Equal(t, tt.modelName, f.Association.ModelName)
		assert.Equal(t, tt.tableName, f.Association.TableName)
		assert.Equal(t, tt.name, f.Association.FieldName)
	}
}

func TestSchema_PrimaryKeyField(t *testing.T) {
	// Implicit
	schema, err := sqlxx.GetSchema(ImplicitPrimaryKey{})
	assert.Nil(t, err)
	assert.Equal(t, schema.PrimaryKeyField.Name, "ID")
	assert.Equal(t, schema.PrimaryKeyField.ColumnName, "id")

	// Explicit
	schema, err = sqlxx.GetSchema(ExplicitPrimaryKey{})
	assert.Nil(t, err)
	assert.Equal(t, schema.PrimaryKeyField.Name, "TadaID")
	assert.Equal(t, schema.PrimaryKeyField.ColumnName, "tada_id")
}

func testFields(t *testing.T, schema sqlxx.Schema, results []fieldResultTest) {
	for _, r := range results {
		assert.Equal(t, r.columnName, schema.Fields[r.field].ColumnName)
		assert.Equal(t, r.columnPath, schema.Fields[r.field].ColumnPath())
	}
}

func testRelations(t *testing.T, schema sqlxx.Schema, results []associationResultTest) {
	for _, r := range results {
		var (
			f     = schema.Associations[r.field]
			assoc = f.Association
		)

		assert.Equal(t, r.field, f.Name)
		assert.Equal(t, r.columnName, f.ColumnName)
		assert.Equal(t, r.columnPath, f.ColumnPath())
		assert.Equal(t, r.associationColumnName, assoc.ColumnName)
		assert.Equal(t, r.associationColumnPath, assoc.ColumnPath())
		assert.Equal(t, r.associationType, assoc.Type)
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
