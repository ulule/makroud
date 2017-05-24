package sqlxx_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestExec_InParams(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	env.createUser("batman")
	env.createUser("robin")
	env.createUser("catwoman")

	query := "UPDATE users SET is_active = false WHERE username IN (?);"
	params := []string{"batman", "robin", "catwoman"}
	queries, err := sqlxx.ExecInParamsWithQueries(env.driver, query, params)
	is.NoError(err)
	is.Len(queries, 1)
	is.Equal("UPDATE users SET is_active = false WHERE username IN (?, ?, ?);", queries[0].Query)
	is.Equal([]interface{}{"batman", "robin", "catwoman"}, queries[0].Args)

	query = "SELECT COUNT(*) FROM users WHERE is_active = :is_active"
	stmt, err := env.driver.PrepareNamed(query)
	is.NoError(err)
	is.NotNil(stmt)

	count := 0
	err = stmt.Get(&count, map[string]interface{}{
		"is_active": false,
	})
	is.NoError(err)
	is.Equal(3, count)
}

func TestFind_InParams(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	env.createUser("batman")
	env.createUser("robin")
	env.createUser("catwoman")

	query := "SELECT * FROM users WHERE is_active = true AND username IN (?);"
	users := &[]User{}
	params := []string{"batman", "robin", "catwoman"}
	queries, err := sqlxx.FindInParamsWithQueries(env.driver, users, query, params)
	is.NoError(err)
	is.Len(queries, 1)
	is.Equal("SELECT * FROM users WHERE is_active = true AND username IN (?, ?, ?);", queries[0].Query)
	is.Equal([]interface{}{"batman", "robin", "catwoman"}, queries[0].Args)

	is.Len(*users, 3)
	hasBatman := false
	hasRobin := false
	hasCatwoman := false
	for _, user := range *users {
		switch user.Username {
		case "batman":
			hasBatman = true
		case "robin":
			hasRobin = true
		case "catwoman":
			hasCatwoman = true
		default:
			is.FailNow("unexpected username", user.Username)
		}
	}
	is.True(hasBatman)
	is.True(hasRobin)
	is.True(hasCatwoman)
}

func TestExec_Simple(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	batman := env.createUser("batman")
	payload := struct {
		Username string
	}{
		Username: batman.Username,
	}

	query := "UPDATE users SET is_active = false WHERE username = :username;"
	err := sqlxx.Exec(env.driver, query, payload)
	is.NoError(err)

	user := &User{}
	err = sqlxx.GetByParams(env.driver, user, map[string]interface{}{
		"id": batman.ID,
	})
	is.NoError(err)
	is.False(user.IsActive)
	is.Equal("batman", user.Username)
}

func TestExec_Named(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	batman := env.createUser("batman")

	query := "UPDATE users SET is_active = false WHERE username = :username;"
	err := sqlxx.NamedExec(env.driver, query, map[string]interface{}{
		"username": batman.Username,
	})
	is.NoError(err)

	user := &User{}
	err = sqlxx.GetByParams(env.driver, user, map[string]interface{}{
		"id": batman.ID,
	})
	is.NoError(err)
	is.False(user.IsActive)
	is.Equal("batman", user.Username)
}

func TestExec_Sync(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	batman := env.createUser("baman")
	batman.IsActive = false
	batman.Username = "batman"

	t0 := time.Now()

	query := `
		UPDATE users
		SET is_active = :is_active,
		    username = :username,
		    updated_at = NOW()
	 	WHERE id = :id
		RETURNING updated_at;
	`
	err := sqlxx.Sync(env.driver, query, batman)
	is.NoError(err)
	is.True(t0.Before(batman.UpdatedAt))

	user := &User{}
	err = sqlxx.GetByParams(env.driver, user, map[string]interface{}{
		"id": batman.ID,
	})
	is.NoError(err)
	is.False(user.IsActive)
	is.Equal("batman", user.Username)
}
