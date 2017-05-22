package sqlxx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestSave_Save(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	username := "thoas"
	createdAt := time.Date(2016, 17, 6, 23, 10, 02, 0, time.UTC)
	isActive := false
	user := &User{Username: username, IsActive: isActive, CreatedAt: createdAt}

	queries, err := sqlxx.SaveWithQueries(env.driver, user)
	is.NoError(err)
	is.NotNil(queries)
	is.Len(queries, 1)
	is.Contains(queries[0].Query, "INSERT INTO users")
	is.Contains(queries[0].Query, ":created_at")
	is.Equal(createdAt, queries[0].Params["created_at"])
	is.Contains(queries[0].Query, ":username")
	is.Equal(username, queries[0].Params["username"])
	_, ok := queries[0].Params["is_active"]
	is.False(ok)
	is.NotContains(queries[0].Query, ":is_active")
	_, ok = queries[0].Params["updated_at"]
	is.False(ok)
	is.NotContains(queries[0].Query, ":updated_at")
	is.NotZero(user.ID)
	is.Equal(true, user.IsActive)
	is.NotZero(user.UpdatedAt)

	user.Username = "gilles"

	queries, err = sqlxx.SaveWithQueries(env.driver, user)
	is.NoError(err)
	is.Contains(queries[0].Query, "UPDATE users SET")
	is.Contains(queries[0].Query, "username = :username")
	is.Equal("gilles", queries[0].Params["username"])

	m := map[string]interface{}{"username": "gilles"}

	query := `
	SELECT count(*)
	FROM %s
	WHERE username = :username
	`

	stmt, err := env.driver.PrepareNamed(fmt.Sprintf(query, user.TableName()))
	is.NoError(err)

	var count int
	err = stmt.Get(&count, m)
	is.NoError(err)
	is.Equal(1, count)
}
