package sqlxx

import (
	"fmt"
	"strings"
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

	// SELEC IN
	users = []User{}
	require.NoError(t, FindByParams(db, &users, map[string]interface{}{"is_active": true, "id": []int{1, 2, 3}}))
	is.Equal(1, users[0].ID)
	is.Equal("jdoe", user.Username)

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

func TestPreload(t *testing.T) {
	is := assert.New(t)

	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	// Queries on zero values must fail

	article := &Article{}
	is.NotNil(Preload(db, article, "Author"))
	is.NotNil(Preload(db, article, "Author.Avatars"))

	// Invalid relations must fail

	article = &fixtures.Articles[0]
	user := &fixtures.User
	is.NotNil(Preload(db, article, "Foo"))

	// Single instance / first level / OneTo relation

	is.Nil(Preload(db, article, "Author"))
	is.Equal(fixtures.User.ID, article.AuthorID)
	is.Equal(fixtures.User.ID, article.Author.ID)
	is.Equal(fixtures.User.Username, article.Author.Username)

	// Single instance / first level / ManyTo relation

	is.Nil(Preload(db, user, "Avatars"))
	is.Len(user.Avatars, 5)
	for i := 0; i < 5; i++ {
		is.Equal(i+1, user.Avatars[i].ID)
		is.Equal(user.ID, user.Avatars[i].UserID)
		is.Equal(fmt.Sprintf("/avatars/jdoe-%d.png", i), user.Avatars[i].Path)
	}

	// Slice of instances / first level / OneTo relation

	articles := fixtures.Articles
	is.Nil(Preload(db, &articles, "Author"))
	for _, article := range articles {
		is.Equal(user.ID, article.Author.ID)
		is.Equal(user.ID, article.AuthorID)
		is.Equal(user.Username, article.Author.Username)
	}

	// Slice of instances / first level / ManyTo relation

	users := []User{}
	for i := 1; i < 6; i++ {
		users = append(users, createUser(t, db, fmt.Sprintf("user%d", i)))
	}

	for _, user := range users {
		is.Zero(user.Avatars)
	}

	is.Nil(Preload(db, &users, "Avatars"))

	for _, user := range users {
		is.NotZero(user.Avatars)
		for _, avatar := range user.Avatars {
			is.NotZero(avatar.ID)
			is.Equal(user.ID, avatar.UserID)
			is.Equal(user.ID, avatar.UserID)
			is.True(strings.HasPrefix(avatar.Path, fmt.Sprintf("/avatars/%s-", user.Username)))
		}
	}
}
