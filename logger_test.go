package sqlxx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestLogger(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	logger := &logger{
		logs: make(chan string, 10),
	}

	driver, err := sqlxx.New(
		dbParamString(sqlxx.Host, "host", "PGHOST"),
		dbParamInt(sqlxx.Port, "port", "PGPORT"),
		dbParamString(sqlxx.User, "user", "PGUSER"),
		dbParamString(sqlxx.Password, "password", "PGPASSWORD"),
		dbParamString(sqlxx.Database, "name", "PGDATABASE"),
		sqlxx.Cache(false),
		sqlxx.WithLogger(logger),
	)
	is.NoError(err)
	is.NotNil(driver)

	user := &User{Username: "thoas", IsActive: false}

	err = sqlxx.Save(driver, user)
	is.NoError(err)
	log, err := logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "INSERT INTO users")
	is.Contains(log, "'thoas'")
	is.Contains(log, "now()")

	deletedAt := time.Date(2016, 06, 07, 21, 30, 28, 0, time.UTC)
	user = &User{Username: "novln", DeletedAt: &deletedAt}
	err = sqlxx.Save(driver, user)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "INSERT INTO users")
	is.Contains(log, "'novln'")
	is.Contains(log, "now()")
	is.Contains(log, "'2016-06-07T21:30:28Z'")

	user.CreatedAt = time.Date(2016, 02, 25, 07, 36, 17, 0, time.UTC)
	err = sqlxx.Save(driver, user)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "UPDATE users SET")
	is.Contains(log, "deleted_at = '2016-06-07T21:30:28Z'")
	is.Contains(log, "created_at = '2016-02-25T07:36:17Z'")
	is.Contains(log, "username = 'novln'")
	is.Contains(log, fmt.Sprintf("WHERE users.id = %d", user.ID))

	err = sqlxx.Archive(driver, user, "DeletedAt")
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "UPDATE users SET deleted_at = ")
	is.Contains(log, fmt.Sprintf("WHERE users.id = %d;", user.ID))

	err = sqlxx.Delete(driver, user)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Equal(fmt.Sprintf("DELETE FROM users WHERE users.id = %d;", user.ID),
		log)

	user = &User{}
	params := map[string]interface{}{"username": "thoas"}
	err = sqlxx.GetByParams(driver, user, params)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "SELECT users.")
	is.Contains(log, "FROM users WHERE users.username = 'thoas' LIMIT 1;")

	users := &[]User{}
	params = map[string]interface{}{"is_active": true}
	err = sqlxx.FindByParams(driver, users, params)
	is.NoError(err)
	log, err = logger.read()
	is.NoError(err)
	t.Log(log)
	is.Contains(log, "SELECT users.")
	is.Contains(log, "FROM users WHERE users.is_active = true;")

	batman := env.createUser("batman")
	err = sqlxx.Preload(driver, batman, "Avatars", "APIKey")
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
