package sqlxx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestSchema_GetSchema_InfiniteLoop(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	schema, err := sqlxx.GetSchema(env.driver, User{})
	is.NoError(err)
	is.NotNil(schema)
}

func TestSchema_FieldNames(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	schema, err := sqlxx.GetSchema(env.driver, Untagged{})
	is.NoError(err)
	is.NotNil(schema)
	is.NotContains(schema.FieldNames(), "unexportedField")
	is.NotContains(schema.FieldNames(), "DBExcludedField")
}

func TestSchema_GetColumns(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	expected := "untagged.first_name, untagged.id, untagged.last_name, untagged.this_is_a_very_long_field_name123"
	columns, err := sqlxx.GetColumns(env.driver, Untagged{})
	is.NoError(err)
	is.NotNil(columns)
	is.Equal(expected, columns)
}

func TestSchema_Fields(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

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
		schema, err := sqlxx.GetSchema(env.driver, modelResults.model)
		is.NoError(err)

		for _, tt := range modelResults.results {
			f, ok := schema.Fields[tt.name]
			is.True(ok)
			is.Equal(tt.model, f.ModelName)
			is.Equal(tt.table, f.TableName)
			is.Equal(tt.name, f.FieldName)
			is.Equal(tt.column, f.ColumnName)
			is.Equal(tt.columnPath, f.ColumnPath())
		}
	}
}

func TestSchema_Associations(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	schema, err := sqlxx.GetSchema(env.driver, Article{})
	is.NoError(err)
	is.NotNil(schema)

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
		is.True(ok, tt.path)
		is.Equal(tt.fk.ModelName, f.ForeignKey.ModelName)
		is.Equal(tt.fk.TableName, f.ForeignKey.TableName)
		is.Equal(tt.fk.FieldName, f.ForeignKey.FieldName)
		is.Equal(tt.fk.ColumnName, f.ForeignKey.ColumnName)
		is.Equal(tt.fk.AssociationFieldName, f.ForeignKey.AssociationFieldName)
		is.Equal(tt.fk.Reference.ModelName, f.ForeignKey.Reference.ModelName)
		is.Equal(tt.fk.Reference.TableName, f.ForeignKey.Reference.TableName)
		is.Equal(tt.fk.Reference.FieldName, f.ForeignKey.Reference.FieldName)
		is.Equal(tt.fk.Reference.ColumnName, f.ForeignKey.Reference.ColumnName)
	}
}

func TestSchema_PrimaryKeyField(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	// Implicit
	schema, err := sqlxx.GetSchema(env.driver, ImplicitPrimaryKey{})
	is.NoError(err)
	is.Equal(schema.PrimaryKeyField.FieldName, "ID")
	is.Equal(schema.PrimaryKeyField.ColumnName, "id")

	// Explicit
	schema, err = sqlxx.GetSchema(env.driver, ExplicitPrimaryKey{})
	is.NoError(err)
	is.Equal(schema.PrimaryKeyField.FieldName, "TadaID")
	is.Equal(schema.PrimaryKeyField.ColumnName, "tada_id")
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
