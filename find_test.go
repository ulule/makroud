package sqlxx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPrimaryKeys(t *testing.T) {
	is := assert.New(t)

	_, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	pks, err := GetPrimaryKeys(&fixtures.Articles, "ID")
	is.Nil(err)
	is.Equal([]interface{}{1, 2, 3, 4, 5}, pks)

	pks, err = GetPrimaryKeys(&fixtures.Articles[0], "ID")
	is.Nil(err)
	is.Equal([]interface{}{1}, pks)
}

func TestGetByParams(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{}
	require.NoError(t, GetByParams(db, &user, map[string]interface{}{"username": "jdoe", "is_active": true}))

	is.Equal(1, user.ID)
	is.Equal("jdoe", user.Username)
	is.True(user.IsActive)
	is.NotZero(user.CreatedAt)
	is.NotZero(user.UpdatedAt)
}

func TestFindByParams(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	users := []User{}
	require.NoError(t, FindByParams(db, &users, map[string]interface{}{"is_active": true}))

	is.Len(users, 1)

	user := users[0]
	is.Equal(1, user.ID)
	is.Equal("jdoe", user.Username)
	is.True(user.IsActive)
	is.NotZero(user.CreatedAt)
	is.NotZero(user.UpdatedAt)

	// SELEC IN
	users = []User{}
	require.NoError(t, FindByParams(db, &users, map[string]interface{}{"is_active": true, "id": []int{1, 2, 3}}))
	is.Equal(1, users[0].ID)
	is.Equal("jdoe", user.Username)

}
