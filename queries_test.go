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

// ----------------------------------------------------------------------------
// Preloads
// ----------------------------------------------------------------------------

func TestPreload_IDZeroValue(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	article := Article{}
	is.NotNil(Preload(db, &article, "Author"))
}

func TestPreload_UnknownRelation(t *testing.T) {
	is := assert.New(t)

	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	article := fixtures.Articles[0]
	is.NotNil(Preload(db, &article, "Foo"))
}

func TestPreload_NullPrimaryKey(t *testing.T) {
	is := assert.New(t)

	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	category := createCategory(t, db, "cat1", nil)
	is.Nil(Preload(db, &category, "User"))

	category = createCategory(t, db, "cat1", &fixtures.User.ID)
	is.Nil(Preload(db, &category, "User"))
	is.NotZero(category.UserID)
	is.NotZero(category.User.ID)
}

func TestPreload_RelationPointer(t *testing.T) {
	is := assert.New(t)

	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	article := fixtures.Articles[0]
	is.Nil(Preload(db, &article, "Reviewer"))
	is.Equal(fixtures.User.ID, article.Reviewer.ID)
	is.Equal(fixtures.User.Username, article.Reviewer.Username)
}

// ----------------------------------------------------------------------------
// Preloads: OneToMany
// ----------------------------------------------------------------------------

func TestPreload_OneToMany_Level1(t *testing.T) {
	is := assert.New(t)

	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	// Instance

	article := fixtures.Articles[0]
	is.Nil(Preload(db, &article, "Author"))
	is.Equal(fixtures.User.ID, article.AuthorID)
	is.Equal(fixtures.User.Username, article.Author.Username)
}

func TestPreload_OneToMany_Level2(t *testing.T) {
	is := assert.New(t)

	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	// Instance

	article := fixtures.Articles[0]

	is.Nil(Preload(db, &article, "Author", "Author.APIKey"))

	is.NotZero(article.Author.APIKey.ID)
	is.Equal("this-is-my-scret-api-key", article.Author.APIKey.Key)

	// Slice

	articles := fixtures.Articles
	user := fixtures.User
	is.Nil(Preload(db, &articles, "Author", "Author.APIKey"))
	for _, article := range articles {
		is.Equal(user.ID, article.Author.ID)
		is.NotZero(article.Author.APIKeyID)
		is.Equal("this-is-my-scret-api-key", article.Author.APIKey.Key)
	}
}

// ----------------------------------------------------------------------------
// Preloads: ManyToOne
// ----------------------------------------------------------------------------

func TestPreload_ManyToOne_Level1(t *testing.T) {
	is := assert.New(t)

	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	// Instance

	user := fixtures.User
	is.Nil(Preload(db, &user, "Avatars"))
	is.Len(user.Avatars, 5)

	for i := 0; i < 5; i++ {
		is.Equal(i+1, user.Avatars[i].ID)
		is.Equal(user.ID, user.Avatars[i].UserID)
		is.Equal(fmt.Sprintf("/avatars/jdoe-%d.png", i), user.Avatars[i].Path)
	}

	// Slice

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
