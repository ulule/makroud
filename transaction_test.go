package sqlxx_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestTransaction_Commit(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	user := &UserV2{Username: "thaos"}

	queries, err := sqlxx.SaveWithQueries(env.driver, user)
	is.NoError(err)
	is.NotNil(queries)

	err = sqlxx.Transaction(env.driver, func(driver sqlxx.Driver) error {
		user.Username = "thoas"

		queries, err := sqlxx.SaveWithQueries(driver, user)
		is.NoError(err)
		is.NotNil(queries)

		return nil
	})
	is.NoError(err)

	record := &User{}
	queries, err = sqlxx.GetByParamsWithQueries(env.driver, record, map[string]interface{}{"id": user.ID})
	is.NoError(err)
	is.NotNil(queries)
	is.Equal("thoas", record.Username)
}

func TestTransaction_Rollback(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	user := &UserV2{Username: "thaos"}
	timeout := errors.New("tcp: read timeout on 10.0.3.11:7000")

	queries, err := sqlxx.SaveWithQueries(env.driver, user)
	is.NoError(err)
	is.NotNil(queries)

	err = sqlxx.Transaction(env.driver, func(driver sqlxx.Driver) error {
		user.Username = "thoas"

		queries, err := sqlxx.SaveWithQueries(driver, user)
		is.NoError(err)
		is.NotNil(queries)

		return timeout
	})
	is.Error(err)
	is.Equal(timeout, err)

	record := &User{}
	queries, err = sqlxx.GetByParamsWithQueries(env.driver, record, map[string]interface{}{"id": user.ID})
	is.NoError(err)
	is.NotNil(queries)
	is.Equal("thaos", record.Username)
}
