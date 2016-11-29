package sqlxx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
}

func TestSave(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}
	require.NoError(t, Save(db, &user))

	is.NotZero(user.ID)
	is.Equal(true, user.IsActive)
	is.NotZero(user.UpdatedAt)

	user.Username = "gilles"
	require.NoError(t, Save(db, &user))

	m := map[string]interface{}{"username": "gilles"}

	query := `
	SELECT count(*)
	FROM %s
	WHERE username = :username
	`

	stmt, err := db.PrepareNamed(fmt.Sprintf(query, user.TableName()))
	require.NoError(t, err)

	var count int

	err = stmt.Get(&count, m)
	require.NoError(t, err)

	is.Equal(1, count)
}

func TestDelete(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}
	require.NoError(t, Save(db, &user))
	require.NoError(t, Delete(db, &user))

	m := map[string]interface{}{"username": "thoas"}

	query := `
	SELECT count(*)
	FROM %s
	WHERE username = :username
	`

	stmt, err := db.PrepareNamed(fmt.Sprintf(query, user.TableName()))
	require.NoError(t, err)

	var count int

	err = stmt.Get(&count, m)
	require.NoError(t, err)

	is.Equal(0, count)
}

func TestSoftDelete(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := User{Username: "thoas"}
	require.NoError(t, Save(db, &user))
	require.NoError(t, SoftDelete(db, &user, "DeletedAt"))

	m := map[string]interface{}{"username": "thoas"}

	query := `
	SELECT count(*)
	FROM %s
	WHERE username = :username
	AND deleted_at IS NULL
	`

	stmt, err := db.PrepareNamed(fmt.Sprintf(query, user.TableName()))
	require.NoError(t, err)

	var count int

	err = stmt.Get(&count, m)
	require.NoError(t, err)

	is.Equal(0, count)
}

func TestPreload(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	article := &Article{}
	// Test with invalid relations
	is.NotNil(Preload(db, article, "Foo"))

	// Test with valid relation
	is.Nil(Preload(db, article, "Author"))
	is.Nil(Preload(db, article, "Author.Avatars"))

}
