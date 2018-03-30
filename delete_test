package sqlxx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestDelete_Delete(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	user := &User{Username: "thoas"}
	err := sqlxx.Save(env.driver, user)
	is.NoError(err)

	queries, err := sqlxx.DeleteWithQueries(env.driver, user)
	is.NoError(err)
	is.NotNil(queries)
	is.Len(queries, 1)
	is.Contains(queries[0].Query, "DELETE FROM users WHERE users.id = :id")

	query := `
		SELECT COUNT(*)
		FROM users
		WHERE username = :username
	`
	params := map[string]interface{}{
		"username": "thoas",
	}

	stmt, err := env.driver.PrepareNamed(query)
	is.NoError(err)
	is.NotNil(stmt)

	count := -1
	err = stmt.Get(&count, params)
	is.NoError(err)
	is.Equal(0, count)
}

func TestDelete_Archive(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	user := &User{Username: "thoas"}
	err := sqlxx.Save(env.driver, user)
	is.NoError(err)

	queries, err := sqlxx.ArchiveWithQueries(env.driver, user)
	is.NoError(err)
	is.NotNil(queries)
	is.Len(queries, 1)
	is.Contains(queries[0].Query, "UPDATE users SET deleted_at = :deleted_at WHERE users.id = :id")

	query := `
		SELECT COUNT(*)
		FROM users
		WHERE username = :username
		AND deleted_at IS NULL
	`
	params := map[string]interface{}{
		"username": "thoas",
	}

	stmt, err := env.driver.PrepareNamed(query)
	is.NoError(err)
	is.NotNil(stmt)

	count := -1
	err = stmt.Get(&count, params)
	is.NoError(err)
	is.Equal(0, count)
}
