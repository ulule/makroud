package sqlxx_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestSchema_GetSchema_InfiniteLoop(t *testing.T) {
	schema, err := sqlxx.GetSchema(User{})
	assert.Nil(t, err)
	assert.NotNil(t, schema)
}

func TestSchema_FieldNames(t *testing.T) {
	schema, err := sqlxx.GetSchema(Untagged{})
	assert.Nil(t, err)

	assert.NotContains(t, schema.FieldNames(), "unexportedField")
	assert.NotContains(t, schema.FieldNames(), "DBExcludedField")
}

func TestSchema_Fields(t *testing.T) {
	type Results struct {
		name       string
		model      string
		table      string
		column     string
		columnPath string
	}

	type ModelResults struct {
		model   interface{}
		results []Results
	}

	untaggedModelResults := ModelResults{
		model: Untagged{},
		results: []Results{
			{"ID", "Untagged", "untagged", "id", "untagged.id"},
			{"FirstName", "Untagged", "untagged", "first_name", "untagged.first_name"},
			{"LastName", "Untagged", "untagged", "last_name", "untagged.last_name"},
			{"ThisIsAVeryLongFieldName123", "Untagged", "untagged", "this_is_a_very_long_field_name123", "untagged.this_is_a_very_long_field_name123"},
		},
	}

	taggedModelResults := ModelResults{
		model: Tagged{},
		results: []Results{
			{"ID", "Tagged", "tagged", "public_id", "tagged.public_id"},
			{"FirstName", "Tagged", "tagged", "firstname", "tagged.firstname"},
			{"LastName", "Tagged", "tagged", "last_name", "tagged.last_name"},
			{"ThisIsAVeryLongFieldName123", "Tagged", "tagged", "short_field", "tagged.short_field"},
		},
	}

	for _, modelResults := range []ModelResults{untaggedModelResults, taggedModelResults} {
		schema, err := sqlxx.GetSchema(modelResults.model)
		assert.Nil(t, err)

		for _, tt := range modelResults.results {
			f, ok := schema.Fields[tt.name]
			assert.True(t, ok)
			assert.Equal(t, tt.model, f.ModelName)
			assert.Equal(t, tt.table, f.TableName)
			assert.Equal(t, tt.name, f.FieldName)
			assert.Equal(t, tt.column, f.ColumnName)
			assert.Equal(t, tt.columnPath, f.ColumnPath())
		}
	}
}

func TestSchema_Associations(t *testing.T) {
	schema, err := sqlxx.GetSchema(Article{})
	assert.Nil(t, err)

	results := []struct {
		path string
		fk   *sqlxx.ForeignKey
	}{
		{
			path: "Author",
			fk: &sqlxx.ForeignKey{
				ModelName:            "Article",
				TableName:            "articles",
				FieldName:            "AuthorID",
				ColumnName:           "author_id",
				AssociationFieldName: "Author",
				Reference: &sqlxx.ForeignKey{
					ModelName:  "User",
					TableName:  "users",
					FieldName:  "ID",
					ColumnName: "id",
				},
			},
		},
		{
			path: "Author.Avatars",
			fk: &sqlxx.ForeignKey{
				ModelName:            "Avatar",
				TableName:            "avatars",
				FieldName:            "UserID",
				ColumnName:           "user_id",
				AssociationFieldName: "User",
				Reference: &sqlxx.ForeignKey{
					ModelName:            "User",
					TableName:            "users",
					FieldName:            "ID",
					ColumnName:           "id",
					AssociationFieldName: "Avatars",
				},
			},
		},
		{
			path: "Author.APIKey",
			fk: &sqlxx.ForeignKey{
				ModelName:            "User",
				TableName:            "users",
				FieldName:            "APIKeyID",
				ColumnName:           "api_key_id",
				AssociationFieldName: "APIKey",
				Reference: &sqlxx.ForeignKey{
					ModelName:  "APIKey",
					TableName:  "api_keys",
					FieldName:  "ID",
					ColumnName: "id",
				},
			},
		},
		{
			path: "Author.APIKey.Partner",
			fk: &sqlxx.ForeignKey{
				ModelName:            "APIKey",
				TableName:            "api_keys",
				FieldName:            "PartnerID",
				ColumnName:           "partner_id",
				AssociationFieldName: "Partner",
				Reference: &sqlxx.ForeignKey{
					ModelName:  "Partner",
					TableName:  "partners",
					FieldName:  "ID",
					ColumnName: "id",
				},
			},
		},
	}

	for _, tt := range results {
		f, ok := schema.Associations[tt.path]
		assert.True(t, ok, tt.path)
		assert.Equal(t, tt.fk.ModelName, f.ForeignKey.ModelName)
		assert.Equal(t, tt.fk.TableName, f.ForeignKey.TableName)
		assert.Equal(t, tt.fk.FieldName, f.ForeignKey.FieldName)
		assert.Equal(t, tt.fk.ColumnName, f.ForeignKey.ColumnName)
		assert.Equal(t, tt.fk.AssociationFieldName, f.ForeignKey.AssociationFieldName)
		assert.Equal(t, tt.fk.Reference.ModelName, f.ForeignKey.Reference.ModelName)
		assert.Equal(t, tt.fk.Reference.TableName, f.ForeignKey.Reference.TableName)
		assert.Equal(t, tt.fk.Reference.FieldName, f.ForeignKey.Reference.FieldName)
		assert.Equal(t, tt.fk.Reference.ColumnName, f.ForeignKey.Reference.ColumnName)
	}
}

func TestSchema_PrimaryKeyField(t *testing.T) {
	// Implicit
	schema, err := sqlxx.GetSchema(ImplicitPrimaryKey{})
	assert.Nil(t, err)
	assert.Equal(t, schema.PrimaryKeyField.FieldName, "ID")
	assert.Equal(t, schema.PrimaryKeyField.ColumnName, "id")

	// Explicit
	schema, err = sqlxx.GetSchema(ExplicitPrimaryKey{})
	assert.Nil(t, err)
	assert.Equal(t, schema.PrimaryKeyField.FieldName, "TadaID")
	assert.Equal(t, schema.PrimaryKeyField.ColumnName, "tada_id")
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
