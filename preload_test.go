package sqlxx_test

import (
	"fmt"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

// ----------------------------------------------------------------------------
// Errors
// ----------------------------------------------------------------------------

func TestPreload_Error_Unaddressable(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	article := Article{}

	queries, err := sqlxx.PreloadWithQueries(env.driver, article, "Author")
	assert.Error(t, err)
	assert.Nil(t, queries)
}

func TestPreload_Error_UnknownRelation(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	article := env.Articles[0]

	queries, err := sqlxx.PreloadWithQueries(env.driver, &article, "Foo")
	assert.Error(t, err)
	assert.Nil(t, queries)
	assert.Zero(t, article.Author)
}

// ----------------------------------------------------------------------------
// Primary keys
// ----------------------------------------------------------------------------

func TestPreload_PrimaryKey_Null(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	category := env.createCategory("cat1", nil)

	queries, err := sqlxx.PreloadWithQueries(env.driver, &category, "User")
	assert.NoError(t, err)
	assert.Nil(t, queries)
	assert.Zero(t, category.User)

	category = env.createCategory("cat1", &env.Users[0].ID)

	queries, err = sqlxx.PreloadWithQueries(env.driver, &category, "User")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "FROM users WHERE users.id = ? LIMIT 1")
	assert.Len(t, queries[0].Args, 1)
	assert.Equal(t, category.UserID.Int64, queries[0].Args[0])
	assert.NotZero(t, category.UserID)
	assert.NotZero(t, category.User.ID)
}

// ----------------------------------------------------------------------------
// Single instance preloads
// ----------------------------------------------------------------------------

func TestPreload_Single_One_Level1(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	batman := env.createUser("batman")
	env.createUser("spiderman")
	article := env.createArticle(&batman)

	// Value
	queries, err := sqlxx.PreloadWithQueries(env.driver, &article, "Author")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "FROM users WHERE users.id = ? LIMIT 1")
	assert.Len(t, queries[0].Args, 1)
	assert.EqualValues(t, article.AuthorID, queries[0].Args[0])
	assert.NotZero(t, article.Author)
	assert.Equal(t, batman.ID, article.AuthorID)
	assert.Equal(t, batman.Username, article.Author.Username)

	// Pointer
	queries, err = sqlxx.PreloadWithQueries(env.driver, &article, "Reviewer")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "FROM users WHERE users.id = ? LIMIT 1")
	assert.Len(t, queries[0].Args, 1)
	assert.EqualValues(t, article.ReviewerID, queries[0].Args[0])
	assert.NotZero(t, article.Reviewer)
	assert.Equal(t, batman.ID, article.ReviewerID)
	assert.Equal(t, batman.Username, article.Reviewer.Username)
}

func TestPreload_Single_One_Level2(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	user := env.createUser("spiderman")
	article := env.createArticle(&user)

	queries, err := sqlxx.PreloadWithQueries(env.driver, &article, "Author", "Author.APIKey")
	assert.NoError(t, err)
	assert.NotNil(t, queries)

	authorQuery, ok := queries.ByTable("users")
	assert.True(t, ok)
	assert.Contains(t, authorQuery.Query, "WHERE users.id = ? LIMIT 1")
	assert.Len(t, authorQuery.Args, 1)

	apikeyQuery, ok := queries.ByTable("api_keys")
	assert.True(t, ok)
	assert.Contains(t, apikeyQuery.Query, "WHERE api_keys.id = ? LIMIT 1")
	assert.Len(t, apikeyQuery.Args, 1)

	assert.NotZero(t, article.Author)
	assert.NotZero(t, article.Author.APIKey)
	assert.Equal(t, user.ID, article.AuthorID)
	assert.Equal(t, user.Username, article.Author.Username)
	assert.NotZero(t, article.Author.APIKey.ID)
	assert.Equal(t, "spiderman-apikey", article.Author.APIKey.Key)
}

func TestPreload_Single_One_Level2_ValueAndPointer(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	user := env.createUser("spiderman")
	assert.NotEmpty(t, user)

	queries, err := sqlxx.PreloadWithQueries(env.driver, &user, "Avatar")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "FROM media WHERE media.id = ? LIMIT 1")
	assert.Len(t, queries[0].Args, 1)
	assert.NotNil(t, user.Avatar)
	assert.EqualValues(t, user.AvatarID.Int64, queries[0].Args[0])

	article := env.createArticle(&user)
	assert.NotEmpty(t, article)

	comment := env.createComment(&user, &article)
	assert.NotEmpty(t, comment)

	comments := []Comment{comment}

	queries, err = sqlxx.PreloadWithQueries(env.driver, &comments, "User", "User.Avatar")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 2)

	comment = comments[0]

	userQuery, ok := queries.ByTable("users")
	assert.True(t, ok)

	avatarQuery, ok := queries.ByTable("media")
	assert.True(t, ok)

	assert.Contains(t, userQuery.Query, "FROM users WHERE users.id IN (?)")
	assert.Len(t, userQuery.Args, 1)
	assert.EqualValues(t, comment.User.ID, userQuery.Args[0])

	assert.Contains(t, avatarQuery.Query, "WHERE media.id IN (?)")
	assert.Len(t, avatarQuery.Args, 1)
	assert.NotNil(t, user.Avatar)
	assert.EqualValues(t, comment.User.Avatar.ID, avatarQuery.Args[0])

	// Level 1 with Value
	assert.NotZero(t, comment.User)
	assert.Equal(t, user.ID, comment.UserID)
	assert.Equal(t, user.Username, comment.User.Username)

	// Level 2 with Pointer
	if comment.User.Avatar != nil {
		assert.Equal(t, user.Avatar.ID, comment.User.Avatar.ID)
		assert.Equal(t, user.Avatar.Path, comment.User.Avatar.Path)
	}
}

func TestPreload_Single_Many_Level1(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	user := env.createUser("wonderwoman")

	queries, err := sqlxx.PreloadWithQueries(env.driver, &user, "Avatars")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)

	avatarQuery, ok := queries.ByTable("avatars")
	assert.True(t, ok)
	assert.Contains(t, avatarQuery.Query, "WHERE avatars.user_id = ?")
	assert.Len(t, avatarQuery.Args, 1)
	assert.EqualValues(t, user.ID, avatarQuery.Args[0])

	assert.Len(t, user.Avatars, 5)
	for i, a := range user.Avatars {
		assert.NotZero(t, a.ID)
		assert.Equal(t, user.ID, a.UserID)
		assert.Equal(t, fmt.Sprintf("/avatars/wonderwoman-%d.png", i+1), a.Path)
	}
}

func TestPreload_Single_Many_One_Level2(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	user := env.createUser("wonderwoman")

	// Values

	queries, err := sqlxx.PreloadWithQueries(env.driver, &user, "Avatars", "Avatars.Filter")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 2)

	avatarQuery, ok := queries.ByTable("avatars")
	assert.True(t, ok)
	assert.Contains(t, avatarQuery.Query, "WHERE avatars.user_id = ?")
	assert.Len(t, avatarQuery.Args, 1)
	assert.EqualValues(t, user.ID, avatarQuery.Args[0])

	avatarFilterQuery, ok := queries.ByTable("avatar_filters")
	assert.True(t, ok)
	assert.Contains(t, avatarFilterQuery.Query, "avatar_filters.id IN (?, ?, ?, ?, ?)")
	assert.Len(t, avatarFilterQuery.Args, 5)

	assert.Len(t, user.Avatars, 5)
	for i, a := range user.Avatars {
		assert.NotZero(t, a.ID)
		assert.Equal(t, user.ID, a.UserID)
		assert.Equal(t, fmt.Sprintf("/avatars/wonderwoman-%d.png", i+1), a.Path)
		assert.NotZero(t, a.Filter)
		assert.NotZero(t, a.Filter.ID)
	}

	// Pointers

	queries, err = sqlxx.PreloadWithQueries(env.driver, &user, "Avatars", "Avatars.FilterPtr")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 2)

	avatarQuery, ok = queries.ByTable("avatars")
	assert.True(t, ok)
	assert.Contains(t, avatarQuery.Query, "WHERE avatars.user_id = ?")
	assert.Len(t, avatarQuery.Args, 1)
	assert.EqualValues(t, user.ID, avatarQuery.Args[0])

	avatarFilterQuery, ok = queries.ByTable("avatar_filters")
	assert.True(t, ok)
	assert.Contains(t, avatarFilterQuery.Query, "avatar_filters.id IN (?, ?, ?, ?, ?)")
	assert.Len(t, avatarFilterQuery.Args, 5)

	assert.Len(t, user.Avatars, 5)
	for i, a := range user.Avatars {
		assert.NotZero(t, a.ID)
		assert.Equal(t, user.ID, a.UserID)
		assert.Equal(t, fmt.Sprintf("/avatars/wonderwoman-%d.png", i+1), a.Path)
		assert.NotNil(t, a.FilterPtr)
		assert.NotZero(t, a.FilterPtr.ID)
	}
}

// ----------------------------------------------------------------------------
// Slice of instances preloads
// ----------------------------------------------------------------------------

func TestPreload_Slice_Level1_One(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	batman := env.createUser("batman")

	var articles []Article
	for i := 0; i < 5; i++ {
		articles = append(articles, env.createArticle(&batman))
	}

	// Value
	queries, err := sqlxx.PreloadWithQueries(env.driver, &articles, "Author")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "FROM users WHERE users.id IN (?)")
	assert.Len(t, queries[0].Args, 1)

	for _, a := range articles {
		assert.Equal(t, batman.ID, a.AuthorID)
		assert.Equal(t, batman.Username, a.Author.Username)
		assert.EqualValues(t, a.AuthorID, queries[0].Args[0])
	}

	// Pointer
	queries, err = sqlxx.PreloadWithQueries(env.driver, &articles, "Reviewer")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "FROM users WHERE users.id IN (?)")
	assert.Len(t, queries[0].Args, 1)

	for _, a := range articles {
		assert.Equal(t, batman.ID, a.ReviewerID)
		assert.Equal(t, batman.Username, a.Reviewer.Username)
		assert.EqualValues(t, a.ReviewerID, queries[0].Args[0])
	}
}

func TestPreload_Slice_Level1_One_DifferentPointerNull(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	batman := env.createUser("batman")
	robin := env.createUser("robin")
	catwoman := env.createUser("catwoman")
	users := []User{batman, robin, catwoman}

	// user_id => media_id
	avatars := map[int]int{}
	for _, user := range users {
		avatars[user.ID] = int(user.AvatarID.Int64)
	}

	queries, err := sqlxx.PreloadWithQueries(env.driver, &users, "Avatar")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "WHERE media.id IN (?, ?, ?)")
	assert.Len(t, queries[0].Args, 3)

	for _, user := range users {
		avatar := user.Avatar
		assert.NotNil(t, avatar)
		assert.Equal(t, avatar.ID, avatars[user.ID])
		assert.Contains(t, queries[0].Args, int64(avatar.ID))
	}
}

func TestPreload_Slice_Level1_One_Different(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	batman := env.createUser("batman")
	robin := env.createUser("robin")
	catwoman := env.createUser("catwoman")

	article1 := env.createArticle(&batman)
	article2 := env.createArticle(&robin)
	article3 := env.createArticle(&catwoman)
	articles := []Article{article1, article2, article3}

	queries, err := sqlxx.PreloadWithQueries(env.driver, &articles, "Author", "Reviewer")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 2)
	assert.Contains(t, queries[0].Query, "WHERE users.id IN")
	assert.Len(t, queries[0].Args, 3)

	assert.Equal(t, articles[0].AuthorID, batman.ID)
	assert.NotZero(t, articles[0].Author)
	assert.Equal(t, articles[0].ReviewerID, batman.ID)
	assert.NotZero(t, articles[0].Reviewer)

	assert.Equal(t, articles[1].AuthorID, robin.ID)
	assert.NotZero(t, articles[1].Author)
	assert.Equal(t, articles[1].ReviewerID, robin.ID)
	assert.NotZero(t, articles[1].Reviewer)

	assert.Equal(t, articles[2].AuthorID, catwoman.ID)
	assert.NotZero(t, articles[2].Author)
	assert.Equal(t, articles[2].ReviewerID, catwoman.ID)
	assert.NotZero(t, articles[2].Reviewer)

	assert.Equal(t, articles[0].Author, batman)
	assert.Equal(t, articles[1].Author, robin)
	assert.Equal(t, articles[2].Author, catwoman)

	assert.Equal(t, articles[0].Reviewer, &batman)
	assert.Equal(t, articles[1].Reviewer, &robin)
	assert.Equal(t, articles[2].Reviewer, &catwoman)
}

func TestPreload_Slice_Level1_Many(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	users := []User{}
	for i := 1; i < 6; i++ {
		users = append(users, env.createUser(fmt.Sprintf("user%d", i)))
	}

	for _, user := range users {
		assert.Zero(t, user.Avatars)
	}

	queries, err := sqlxx.PreloadWithQueries(env.driver, &users, "Avatars")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 1)
	assert.Contains(t, queries[0].Query, "avatars.user_id IN (?, ?, ?, ?, ?)")
	assert.Len(t, queries[0].Args, 5)

	for _, user := range users {
		assert.NotZero(t, user.Avatars)
		assert.Contains(t, queries[0].Args, int64(user.ID))

		for _, avatar := range user.Avatars {
			assert.NotZero(t, avatar.ID)
			assert.Equal(t, user.ID, avatar.UserID)
			assert.Equal(t, user.ID, avatar.UserID)
			assert.True(t, strings.HasPrefix(avatar.Path, fmt.Sprintf("/avatars/%s-", user.Username)))
		}
	}
}

func TestPreload_Slice_Level2_One_ValueAndPointer(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	user := env.createUser("spiderman")
	article := env.createArticle(&user)
	deadpool := env.createUser("deadpool")
	article2 := env.createArticle(&deadpool)
	articles := []Article{article, article2}

	queries, err := sqlxx.PreloadWithQueries(env.driver, &articles, "Author", "Author.APIKey")
	assert.NoError(t, err)
	assert.NotNil(t, queries)
	assert.Len(t, queries, 2)

	authorQuery, ok := queries.ByTable("users")
	assert.True(t, ok)
	assert.Contains(t, authorQuery.Query, "WHERE users.id IN (?, ?)")

	apikeyQuery, ok := queries.ByTable("api_keys")
	assert.True(t, ok)
	assert.Contains(t, apikeyQuery.Query, "WHERE api_keys.id IN (?, ?)")

	assert.Equal(t, user.ID, articles[0].Author.ID)
	assert.Equal(t, user.ID, articles[0].AuthorID)

	assert.Equal(t, user.Username, articles[0].Author.Username)
	assert.NotZero(t, articles[0].Author.APIKeyID)
	assert.Equal(t, "spiderman-apikey", articles[0].Author.APIKey.Key)

	assert.Equal(t, deadpool.ID, articles[1].Author.ID)
	assert.Equal(t, deadpool.ID, articles[1].AuthorID)
	assert.Equal(t, deadpool.Username, articles[1].Author.Username)

	assert.NotZero(t, articles[1].Author.APIKeyID)
	assert.Equal(t, "deadpool-apikey", articles[1].Author.APIKey.Key)
}
