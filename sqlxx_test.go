package sqlxx

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestToSnakeCase(t *testing.T) {
	var results = []struct {
		in  string
		out string
	}{
		{"FooBar", "foo_bar"},
		{"ID", "id"},
		{"UserID", "user_id"},
		{"User1234", "user1234"},
		{"blahBlah", "blah_blah"},
	}

	for _, tt := range results {
		s := toSnakeCase(tt.in)
		assert.Equal(t, s, tt.out)
	}
}

func TestTableColumns(t *testing.T) {
	is := assert.New(t)

	schema, err := GetModelSchema(StructWithoutTags{})
	is.NoError(err)

	is.Equal(map[string]string{
		"ID":                          "foo.id",
		"FirstName":                   "foo.first_name",
		"LastName":                    "foo.last_name",
		"ThisIsAVeryLongFieldName123": "foo.this_is_a_very_long_field_name123",
	}, schema.Columns)

	is.Equal(schema.Associations["User"], RelatedField{
		FK:          "foo.user_id",
		FKReference: "users.id",
	})

	is.Equal(schema.Associations["UserPtr"], RelatedField{
		FK:          "foo.user_ptr_id",
		FKReference: "users.id",
	})

	schema, err = GetModelSchema(StructWithTags{})
	is.NoError(err)

	is.Equal(map[string]string{
		"ID":                          "foo.public_id",
		"FirstName":                   "foo.firstname",
		"LastName":                    "foo.last_name",
		"ThisIsAVeryLongFieldName123": "foo.short_field",
	}, schema.Columns)

	is.Equal(schema.Associations["User"], RelatedField{
		FK:          "foo.member_id",
		FKReference: "users.custom_id",
	})

	is.Equal(schema.Associations["UserPtr"], RelatedField{
		FK:          "foo.member_id",
		FKReference: "users.custom_id",
	})
}

// ----------------------------------------------------------------------------
// Test data
// ----------------------------------------------------------------------------

type User struct {
	ID   int
	Name string
}

func (User) TableName() string {
	return "users"
}

type StructWithoutTags struct {
	ID                          int
	FirstName                   string
	LastName                    string
	ThisIsAVeryLongFieldName123 string
	User                        User
	UserPtr                     *User
}

func (StructWithoutTags) TableName() string {
	return "foo"
}

type StructWithTags struct {
	ID                          int    `db:"public_id"`
	FirstName                   string `db:"firstname"`
	LastName                    string
	ThisIsAVeryLongFieldName123 string `db:"short_field"`
	User                        User   `db:"member_id" sqlxx:"custom_id"`
	UserPtr                     *User  `db:"member_id" sqlxx:"custom_id"`
}

func (StructWithTags) TableName() string {
	return "foo"
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

func prepareDB(t *testing.T, driverName string) (*sqlx.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	return sqlx.NewDb(db, driverName), mock, func() {
		require.NoError(t, db.Close())
	}
}
