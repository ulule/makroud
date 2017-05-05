package sqlxx_test

import (
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestFind_GetByParams(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{}

	queries, err := sqlxx.GetByParams(db, &user, map[string]interface{}{"username": "jdoe", "is_active": true})
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "users.username = ?")
	assert.Contains(t, queries[0].Query, "users.is_active = ?")
	assert.Contains(t, queries[0].Args, user.Username)
	assert.Contains(t, queries[0].Args, true)

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

	queries, err := sqlxx.FindByParams(db, &users, map[string]interface{}{"is_active": true})
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "users.is_active = ?")
	assert.Contains(t, queries[0].Args, true)
	assert.Len(t, users, 1)

	user := users[0]
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "jdoe", user.Username)
	assert.True(t, user.IsActive)
	assert.NotZero(t, user.CreatedAt)
	assert.NotZero(t, user.UpdatedAt)

	// SELEC IN
	users = []User{}
	queries, err = sqlxx.FindByParams(db, &users, map[string]interface{}{"is_active": true, "id": []int{1, 2, 3}})
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "users.is_active = ?")
	assert.Contains(t, queries[0].Query, "users.id IN (?, ?, ?)")
	assert.Contains(t, queries[0].Args, true)
	assert.Contains(t, queries[0].Args, 1)
	assert.Contains(t, queries[0].Args, 2)
	assert.Contains(t, queries[0].Args, 3)

	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "jdoe", user.Username)
}
