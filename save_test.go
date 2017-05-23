package sqlxx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestSave_Save(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	user := &User{Username: "thoas"}

	queries, err := sqlxx.SaveWithQueries(env.driver, user)
	is.NoError(err)
	is.NotNil(queries)
	is.Len(queries, 1)
	is.Contains(queries[0].Query, "INSERT INTO users")
	is.Contains(queries[0].Query, ":username")
	is.Contains(queries[0].Params, "username")
	is.Equal("thoas", queries[0].Params["username"])

	is.NotZero(user.ID)
	is.Equal(true, user.IsActive)
	is.NotZero(user.UpdatedAt)

	user.Username = "gilles"

	queries, err = sqlxx.SaveWithQueries(env.driver, user)
	is.NoError(err)
	is.NotNil(queries)
	is.Len(queries, 1)
	is.Contains(queries[0].Query, "UPDATE users SET")
	is.Contains(queries[0].Query, "username = :username")

	query := `
		SELECT count(*)
		FROM users
		WHERE username = :username
	`
	params := map[string]interface{}{
		"username": "gilles",
	}

	stmt, err := env.driver.PrepareNamed(query)
	is.NoError(err)
	is.NotNil(stmt)

	count := -1
	err = stmt.Get(&count, params)
	is.NoError(err)
	is.Equal(1, count)
}
