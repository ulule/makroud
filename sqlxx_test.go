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

	is.Equal(schema.Associations["User"].FK.PrefixedName, "foo.user_id")
	is.Equal(schema.Associations["User"].FKReference.PrefixedName, "users.id")

	is.Equal(schema.Associations["UserPtr"].FK.PrefixedName, "foo.user_ptr_id")
	is.Equal(schema.Associations["UserPtr"].FKReference.PrefixedName, "users.id")

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

	is.Equal(schema.Associations["User"].FK.PrefixedName, "foo.member_id")
	is.Equal(schema.Associations["User"].FKReference.PrefixedName, "users.custom_id")

	is.Equal(schema.Associations["UserPtr"].FK.PrefixedName, "foo.member_id")
	is.Equal(schema.Associations["UserPtr"].FKReference.PrefixedName, "users.custom_id")
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
