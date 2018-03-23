package sqlxx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestFind_GetByParams(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	user := &User{}
	queries, err := sqlxx.GetByParamsWithQueries(env.driver, user, map[string]interface{}{
		"username": "lucius_fox", "is_active": true,
	})
	is.NoError(err)
	is.NotNil(queries)
	is.Len(queries, 1)
	is.Contains(queries[0].Query, "users.username = ?")
	is.Contains(queries[0].Query, "users.is_active = ?")
	is.Contains(queries[0].Args, user.Username)
	is.Contains(queries[0].Args, true)

	is.Equal(1, user.ID)
	is.Equal("lucius_fox", user.Username)
	is.True(user.IsActive)
	is.NotZero(user.CreatedAt)
	is.NotZero(user.UpdatedAt)
}

func TestFind_FindByParams(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	// Execute select WITHOUT clause 'IN'

	list := &UserList{}
	queries, err := sqlxx.FindByParamsWithQueries(env.driver, list, map[string]interface{}{
		"is_active": true,
	})
	is.NoError(err)
	is.NotNil(queries)
	is.Len(queries, 1)
	is.Contains(queries[0].Query, "users.is_active = ?")
	is.Contains(queries[0].Args, true)

	is.Len(list.users, 1)
	user := list.One()
	is.Equal(1, user.ID)
	is.Equal("lucius_fox", user.Username)
	is.True(user.IsActive)
	is.NotZero(user.CreatedAt)
	is.NotZero(user.UpdatedAt)

	// Execute select WITH clause 'IN'

	list = &UserList{}
	queries, err = sqlxx.FindByParamsWithQueries(env.driver, list, map[string]interface{}{
		"is_active": true, "id": []int{1, 2, 3},
	})
	is.NoError(err)
	is.NotNil(queries)
	is.Len(queries, 1)
	is.Contains(queries[0].Query, "users.is_active = ?")
	is.Contains(queries[0].Query, "users.id IN (?, ?, ?)")
	is.Contains(queries[0].Args, true)
	is.Contains(queries[0].Args, 1)
	is.Contains(queries[0].Args, 2)
	is.Contains(queries[0].Args, 3)

	is.Len(list.users, 1)
	user = list.One()
	is.Equal(1, user.ID)
	is.Equal("lucius_fox", user.Username)
	is.True(user.IsActive)
	is.NotZero(user.CreatedAt)
	is.NotZero(user.UpdatedAt)

}
