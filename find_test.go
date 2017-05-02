package sqlxx_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestFind_GetPrimaryKeys(t *testing.T) {
	_, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	pks, err := sqlxx.GetPrimaryKeys(&fixtures.Articles, "ID")
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{1, 2, 3, 4, 5}, pks)

	pks, err = sqlxx.GetPrimaryKeys(&fixtures.Articles[0], "ID")
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{1}, pks)
}

func TestFind_GetByParams(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{}
	assert.NoError(t, sqlxx.GetByParams(db, &user, map[string]interface{}{"username": "jdoe", "is_active": true}))

	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "jdoe", user.Username)
	assert.True(t, user.IsActive)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)
}

func TestFind_FindByParams(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	users := []User{}
	assert.NoError(t, sqlxx.FindByParams(db, &users, map[string]interface{}{"is_active": true}))
	assert.Len(t, users, 1)

	user := users[0]
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "jdoe", user.Username)
	assert.True(t, user.IsActive)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)

	// SELEC IN
	users = []User{}
	assert.NoError(t, sqlxx.FindByParams(db, &users, map[string]interface{}{"is_active": true, "id": []int{1, 2, 3}}))
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "jdoe", user.Username)
}
