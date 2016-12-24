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
	is.Zero(article.Author)
}

func TestPreload_NullPrimaryKey(t *testing.T) {
	is := assert.New(t)

	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	category := createCategory(t, db, "cat1", nil)
	is.Nil(Preload(db, &category, "User"))
	is.Zero(category.User)

	category = createCategory(t, db, "cat1", &fixtures.User.ID)
	is.Nil(Preload(db, &category, "User"))
	is.NotZero(category.UserID)
	is.NotZero(category.User.ID)
}

// ----------------------------------------------------------------------------
// Preloads: OneToMany
// ----------------------------------------------------------------------------

func TestPreload_OneToMany_Level1(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	batman := createUser(t, db, "batman")
	article := createArticle(t, db, &batman)

	//
	// Instance
	//

	// Value
	is.Nil(Preload(db, &article, "Author"))
	is.NotZero(article.Author)
	is.Equal(batman.ID, article.AuthorID)
	is.Equal(batman.Username, article.Author.Username)

	// Pointer
	is.Nil(Preload(db, &article, "Reviewer"))
	is.NotZero(article.Reviewer)
	is.Equal(batman.ID, article.ReviewerID)
	is.Equal(batman.Username, article.Reviewer.Username)

	//
	// Slice
	//

	var articles []Article
	for i := 0; i < 5; i++ {
		articles = append(articles, createArticle(t, db, &batman))
	}

	// Value
	is.Nil(Preload(db, &articles, "Author"))
	for _, a := range articles {
		is.Equal(batman.ID, a.AuthorID)
		is.Equal(batman.Username, a.Author.Username)
	}

	// Pointer
	is.Nil(Preload(db, &articles, "Reviewer"))
	for _, a := range articles {
		is.Equal(batman.ID, a.ReviewerID)
		is.Equal(batman.Username, a.Reviewer.Username)
	}
}

func TestPreload_OneToMany_Level1_Different(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	batman := createUser(t, db, "batman")
	robin := createUser(t, db, "robin")
	catwoman := createUser(t, db, "catwoman")
	article1 := createArticle(t, db, &batman)
	article2 := createArticle(t, db, &robin)
	article3 := createArticle(t, db, &catwoman)

	articles := []Article{
		article1,
		article2,
		article3,
	}

	is.Nil(Preload(db, &articles, "Author", "Reviewer"))
	is.Equal(articles[0].AuthorID, batman.ID)
	is.NotZero(articles[0].Author)
	is.Equal(articles[0].ReviewerID, batman.ID)
	is.NotZero(articles[0].Reviewer)
	is.Equal(articles[1].AuthorID, robin.ID)
	is.NotZero(articles[1].Author)
	is.Equal(articles[1].ReviewerID, robin.ID)
	is.NotZero(articles[1].Reviewer)
	is.Equal(articles[2].AuthorID, catwoman.ID)
	is.NotZero(articles[2].Author)
	is.Equal(articles[2].ReviewerID, catwoman.ID)
	is.NotZero(articles[2].Reviewer)

	is.Equal(articles[0].Author, batman)
	is.Equal(articles[1].Author, robin)
	is.Equal(articles[2].Author, catwoman)

	is.Equal(articles[0].Reviewer, &batman)
	is.Equal(articles[1].Reviewer, &robin)
	is.Equal(articles[2].Reviewer, &catwoman)
}

func TestPreload_OneToMany_Level2(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := createUser(t, db, "spiderman")

	//
	// Instance
	//

	article := createArticle(t, db, &user)
	is.Nil(Preload(db, &article, "Author", "Author.APIKey"))
	is.NotZero(article.Author)
	is.NotZero(article.Author.APIKey)
	is.Equal(user.ID, article.AuthorID)
	is.Equal(user.Username, article.Author.Username)
	is.NotZero(article.Author.APIKey.ID)
	is.Equal("spiderman-apikey", article.Author.APIKey.Key)
}

func TestPreload_OneToMany_Level2_Multiple(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := createUser(t, db, "spiderman")
	article := createArticle(t, db, &user)

	deadpool := createUser(t, db, "deadpool")
	article2 := createArticle(t, db, &deadpool)

	articles := []Article{article, article2}

	is.Nil(Preload(db, &articles, "Author", "Author.APIKey"))

	is.Equal(user.ID, articles[0].Author.ID)
	is.Equal(user.ID, articles[0].AuthorID)
	is.Equal(user.Username, articles[0].Author.Username)
	is.NotZero(articles[0].Author.APIKeyID)
	is.Equal("spiderman-apikey", articles[0].Author.APIKey.Key)

	is.Equal(deadpool.ID, articles[1].Author.ID)
	is.Equal(deadpool.ID, articles[1].AuthorID)
	is.Equal(deadpool.Username, articles[1].Author.Username)
	is.NotZero(articles[1].Author.APIKeyID)
	is.Equal("deadpool-apikey", articles[1].Author.APIKey.Key)
}

// ----------------------------------------------------------------------------
// Preloads: ManyToOne
// ----------------------------------------------------------------------------

func TestPreload_ManyToOne_Level1(t *testing.T) {
	is := assert.New(t)

	db, _, shutdown := dbConnection(t)
	defer shutdown()

	//
	// Instance
	//

	user := createUser(t, db, "wonderwoman")
	is.Nil(Preload(db, &user, "Avatars"))
	is.Len(user.Avatars, 5)

	for i, a := range user.Avatars {
		is.NotZero(a.ID)
		is.Equal(user.ID, a.UserID)
		is.Equal(fmt.Sprintf("/avatars/wonderwoman-%d.png", i+1), a.Path)
	}

	//
	// Slice
	//

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
