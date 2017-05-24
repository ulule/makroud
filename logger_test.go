package sqlxx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestLogger(t *testing.T) {
	logger := &logger{
		logs: make(chan string, 10),
	}

	env := setup(t, sqlxx.WithLogger(logger))
	defer env.teardown()

	is := require.New(t)

	user := &User{Username: "thoas", IsActive: false}
	err := sqlxx.Save(env.driver, user)
	is.NoError(err)
	log, err := logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "INSERT INTO users")
	is.Contains(log, "'thoas'")
	is.Contains(log, "now()")

	deletedAt := time.Date(2016, 06, 07, 21, 30, 28, 0, time.UTC)
	user = &User{Username: "novln", DeletedAt: &deletedAt}
	err = sqlxx.Save(env.driver, user)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "INSERT INTO users")
	is.Contains(log, "'novln'")
	is.Contains(log, "now()")
	is.Contains(log, "'2016-06-07T21:30:28Z'")

	user.CreatedAt = time.Date(2016, 02, 25, 07, 36, 17, 0, time.UTC)
	err = sqlxx.Save(env.driver, user)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "UPDATE users SET")
	is.Contains(log, "deleted_at = '2016-06-07T21:30:28Z'")
	is.Contains(log, "created_at = '2016-02-25T07:36:17Z'")
	is.Contains(log, "username = 'novln'")
	is.Contains(log, fmt.Sprintf("WHERE users.id = %d", user.ID))

	err = sqlxx.Archive(env.driver, user, "DeletedAt")
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "UPDATE users SET deleted_at = ")
	is.Contains(log, fmt.Sprintf("WHERE users.id = %d;", user.ID))

	err = sqlxx.Delete(env.driver, user)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Equal(fmt.Sprintf("DELETE FROM users WHERE users.id = %d;", user.ID), log)

	user = &User{}
	params := map[string]interface{}{"username": "thoas"}
	err = sqlxx.GetByParams(env.driver, user, params)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "SELECT users.")
	is.Contains(log, "FROM users WHERE users.username = 'thoas' LIMIT 1;")

	users := &[]User{}
	params = map[string]interface{}{"is_active": true}
	err = sqlxx.FindByParams(env.driver, users, params)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "SELECT users.")
	is.Contains(log, "FROM users WHERE users.is_active = true;")

	batman := env.createUser("batman")
	err = sqlxx.Preload(env.driver, batman, "Avatars", "APIKey")
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "SELECT avatars.")
	is.Contains(log, fmt.Sprintf("FROM avatars WHERE avatars.user_id = %d;", batman.ID))
	is.Contains(log, "SELECT api_keys.")
	is.Contains(log, fmt.Sprintf("FROM api_keys WHERE api_keys.id = %d LIMIT 1;", batman.APIKeyID))

	log, err = logger.read()
	is.Equal(ErrLogTimeout, err)
	is.Equal("", log)

	query := `
		UPDATE users
		   SET is_active = true
		   WHERE username IN (?);
	`
	list := []string{"batman", "robin", "catwoman"}
	err = sqlxx.ExecInParams(env.driver, query, list)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Equal("UPDATE users SET is_active = true WHERE username IN ('batman', 'robin', 'catwoman');", log)

	query = `SELECT * FROM users WHERE is_active = true AND username IN (?);`
	list = []string{"batman", "robin", "catwoman", "joker"}
	users = &[]User{}
	err = sqlxx.FindInParams(env.driver, users, query, list)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Equal(fmt.Sprint("SELECT * FROM users WHERE is_active = true AND ",
		"username IN ('batman', 'robin', 'catwoman', 'joker');"), log)

	query = `UPDATE users SET is_active = false WHERE username = :username;`
	user = &User{Username: "novln"}
	err = sqlxx.Exec(env.driver, query, user)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Equal("UPDATE users SET is_active = false WHERE username = 'novln';", log)

	query = `UPDATE users SET is_active = true WHERE username = :username;`
	params = map[string]interface{}{
		"username": "novln",
	}
	err = sqlxx.NamedExec(env.driver, query, params)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Equal("UPDATE users SET is_active = true WHERE username = 'novln';", log)

	query = `
		UPDATE users
		SET username = :username,
		    is_active = :is_active,
		    updated_at = NOW()
		WHERE id = :id
		RETURNING updated_at;
	`
	catwoman := env.createUser("catwman")
	catwoman.Username = "catwoman"
	err = sqlxx.Sync(env.driver, query, catwoman)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Equal(fmt.Sprint("UPDATE users SET username = 'catwoman', is_active = true, updated_at = NOW() WHERE id = ",
		catwoman.ID, " RETURNING updated_at;"), log)
}

type logger struct {
	logs chan string
}

func (e *logger) Log(query string, duration time.Duration) {
	e.logs <- query
}

var ErrLogTimeout = fmt.Errorf("logger timeout")

func (e *logger) read() (string, error) {
	select {
	case log := <-e.logs:
		return log, nil
	case <-time.After(500 * time.Millisecond):
		return "", ErrLogTimeout
	}
}
