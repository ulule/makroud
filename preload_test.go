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

func TestPreload_Single_One(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	var (
		user    = env.createUser("batman")
		article = env.createArticle(&user)
	)

	// level 1
	{
		// Value

		assert.Zero(t, article.Author)

		queries, err := sqlxx.PreloadWithQueries(env.driver, &article, "Author")
		assert.NoError(t, err)
		assert.NotNil(t, queries)
		assert.Len(t, queries, 1)

		userQuery, ok := queries.ByTable("users")
		assert.True(t, ok)
		assert.Contains(t, userQuery.Query, "FROM users WHERE users.id = ? LIMIT 1")
		assert.Len(t, userQuery.Args, 1)
		assert.EqualValues(t, article.AuthorID, userQuery.Args[0])
		assert.NotZero(t, article.Author)
		assert.Equal(t, user.ID, article.AuthorID)
		assert.Equal(t, user.Username, article.Author.Username)

		// Pointer

		assert.Nil(t, article.Reviewer)

		queries, err = sqlxx.PreloadWithQueries(env.driver, &article, "Reviewer")
		assert.NoError(t, err)
		assert.NotNil(t, queries)
		assert.Len(t, queries, 1)

		userQuery, ok = queries.ByTable("users")
		assert.True(t, ok)
		assert.Contains(t, userQuery.Query, "FROM users WHERE users.id = ? LIMIT 1")
		assert.Len(t, userQuery.Args, 1)
		assert.EqualValues(t, article.ReviewerID, userQuery.Args[0])
		assert.NotZero(t, article.Reviewer)
		assert.Equal(t, user.ID, article.ReviewerID)
		assert.Equal(t, user.Username, article.Reviewer.Username)
	}

	// level 2
	{
		// Value

		article.Author = User{}

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
		assert.Equal(t, user.ID, article.AuthorID)
		assert.Equal(t, user.Username, article.Author.Username)

		assert.NotZero(t, article.Author.APIKey)
		assert.NotZero(t, article.Author.APIKey.ID)
		assert.Equal(t, fmt.Sprintf("%s-apikey", user.Username), article.Author.APIKey.Key)

		// Pointer

		article.Author = User{}

		queries, err = sqlxx.PreloadWithQueries(env.driver, &article, "Author", "Author.APIKeyPtr")
		assert.NoError(t, err)
		assert.NotNil(t, queries)

		authorQuery, ok = queries.ByTable("users")
		assert.True(t, ok)
		assert.Contains(t, authorQuery.Query, "WHERE users.id = ? LIMIT 1")
		assert.Len(t, authorQuery.Args, 1)

		apikeyQuery, ok = queries.ByTable("api_keys")
		assert.True(t, ok)
		assert.Contains(t, apikeyQuery.Query, "WHERE api_keys.id = ? LIMIT 1")
		assert.Len(t, apikeyQuery.Args, 1)

		assert.NotZero(t, article.Author)
		assert.Equal(t, user.ID, article.AuthorID)
		assert.Equal(t, user.Username, article.Author.Username)

		assert.NotNil(t, article.Author.APIKeyPtr)
		assert.NotZero(t, article.Author.APIKeyPtr.ID)
		assert.Equal(t, fmt.Sprintf("%s-apikey", user.Username), article.Author.APIKeyPtr.Key)
	}

	// Level 3
	{
		// Value

		article.Author = User{}

		queries, err := sqlxx.PreloadWithQueries(env.driver, &article, "Author", "Author.APIKey", "Author.APIKey.Partner")
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

		partnerQuery, ok := queries.ByTable("partners")
		assert.True(t, ok)
		assert.Contains(t, partnerQuery.Query, "WHERE partners.id = ? LIMIT 1")
		assert.Len(t, partnerQuery.Args, 1)

		assert.NotZero(t, article.Author)
		assert.Equal(t, user.ID, article.AuthorID)
		assert.Equal(t, user.Username, article.Author.Username)

		assert.NotZero(t, article.Author.APIKey)
		assert.NotZero(t, article.Author.APIKey.ID)
		assert.Equal(t, fmt.Sprintf("%s-apikey", user.Username), article.Author.APIKey.Key)

		assert.NotZero(t, article.Author.APIKey.Partner)
		assert.NotZero(t, article.Author.APIKey.Partner.ID)
	}
}

func TestPreload_Single_Many(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	user := env.createUser("wonderwoman")

	// Level 1
	{
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

	// Level 2 - One
	{
		user.Avatars = []Avatar{}

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

		user.Avatars = []Avatar{}

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
}

func TestPreload_Slice_One(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	var (
		user     User
		articles []Article
	)

	user = env.createUser("batman")
	for i := 0; i < 5; i++ {
		articles = append(articles, env.createArticle(&user))
	}

	// Level 1
	{
		// Value

		for i := range articles {
			assert.Zero(t, articles[i].Author)
		}

		queries, err := sqlxx.PreloadWithQueries(env.driver, &articles, "Author")
		assert.NoError(t, err)
		assert.NotNil(t, queries)
		assert.Len(t, queries, 1)
		assert.Contains(t, queries[0].Query, "FROM users WHERE users.id IN (?)")
		assert.Len(t, queries[0].Args, 1)

		for _, a := range articles {
			assert.Equal(t, user.ID, a.AuthorID)
			assert.Equal(t, user.Username, a.Author.Username)
			assert.EqualValues(t, a.AuthorID, queries[0].Args[0])
		}

		// Pointer

		for i := range articles {
			assert.Nil(t, articles[i].Reviewer)
		}

		queries, err = sqlxx.PreloadWithQueries(env.driver, &articles, "Reviewer")
		assert.NoError(t, err)
		assert.NotNil(t, queries)
		assert.Len(t, queries, 1)
		assert.Contains(t, queries[0].Query, "FROM users WHERE users.id IN (?)")
		assert.Len(t, queries[0].Args, 1)

		for _, a := range articles {
			assert.Equal(t, user.ID, a.ReviewerID)
			assert.Equal(t, user.Username, a.Reviewer.Username)
			assert.EqualValues(t, a.ReviewerID, queries[0].Args[0])
		}
	}

	// Level 1 + Int conversion
	{
		var (
			batman   = env.createUser("batman")
			robin    = env.createUser("robin")
			catwoman = env.createUser("catwoman")
		)

		users := []User{
			batman,
			robin,
			catwoman,
		}

		// user_id => media_id
		avatars := map[int]int{}
		for _, user := range users {
			avatars[user.ID] = int(user.AvatarID.Int64)
		}

		queries, err := sqlxx.PreloadWithQueries(env.driver, &users, "Avatar")
		assert.NoError(t, err)
		assert.NotNil(t, queries)
		assert.Len(t, queries, 1)

		avatarQuery, ok := queries.ByTable("media")
		assert.True(t, ok)
		assert.Contains(t, avatarQuery.Query, "WHERE media.id IN (?, ?, ?)")
		assert.Len(t, avatarQuery.Args, 3)

		for _, user := range users {
			avatar := user.Avatar
			assert.NotNil(t, avatar)
			assert.Equal(t, avatar.ID, avatars[user.ID])
			assert.Contains(t, queries[0].Args, int64(avatar.ID))
		}
	}

	// Level 1 - Multiple paths (+ mix of value and pointers)
	{
		var (
			batman   = env.createUser("batman")
			robin    = env.createUser("robin")
			catwoman = env.createUser("catwoman")
			article1 = env.createArticle(&batman)
			article2 = env.createArticle(&robin)
			article3 = env.createArticle(&catwoman)
		)

		articles := []*Article{
			&article1,
			&article2,
			&article3,
		}

		queries, err := sqlxx.PreloadWithQueries(env.driver, &articles, "Author", "Reviewer")
		assert.NoError(t, err)
		assert.NotNil(t, queries)
		assert.Len(t, queries, 2)

		userQuery, ok := queries.ByTable("users")
		assert.True(t, ok)
		assert.Contains(t, userQuery.Query, "WHERE users.id IN")
		assert.Len(t, userQuery.Args, 3)

		table := []struct {
			user    User
			article Article
		}{
			{batman, article1},
			{robin, article2},
			{catwoman, article3},
		}

		for _, tt := range table {
			assert.Equal(t, tt.article.AuthorID, tt.user.ID)
			assert.Equal(t, tt.article.Author, tt.user)
			assert.Equal(t, tt.article.ReviewerID, tt.user.ID)
			assert.Equal(t, tt.article.Reviewer, &tt.user)
		}
	}

	// Level 2 and more
	{

		var (
			projects = env.Projects
			project  = projects[0]
		)

		assert.Len(t, projects, 1)
		assert.Nil(t, project.Manager)

		queries, err := sqlxx.PreloadWithQueries(env.driver, &projects, "Manager", "Manager.User", "Manager.User.Avatar")
		assert.NoError(t, err)
		assert.NotNil(t, queries)
		assert.Len(t, queries, 3)

		managerQuery, ok := queries.ByTable("managers")
		assert.True(t, ok)
		assert.Contains(t, managerQuery.Query, "WHERE managers.id IN")
		assert.Len(t, managerQuery.Args, 1)
		assert.Equal(t, int64(env.Projects[0].ManagerID), managerQuery.Args[0])

		userQuery, ok := queries.ByTable("users")
		assert.True(t, ok)
		assert.Contains(t, userQuery.Query, "WHERE users.id IN")
		assert.Len(t, userQuery.Args, 1)
		assert.Equal(t, int64(env.Users[0].ID), userQuery.Args[0])

		mediaQuery, ok := queries.ByTable("media")
		assert.True(t, ok)
		assert.Contains(t, mediaQuery.Query, "WHERE media.id IN")
		assert.Len(t, mediaQuery.Args, 1)
		assert.Equal(t, int64(env.Medias[0].ID), mediaQuery.Args[0])

		for _, project := range projects {
			assert.NotNil(t, project.Manager)
			assert.Equal(t, env.Managers[0].ID, project.Manager.ID)

			assert.NotNil(t, project.Manager.User)
			assert.Equal(t, env.Users[0].ID, project.Manager.User.ID)

			assert.NotNil(t, project.Manager.User.Avatar)
			assert.Equal(t, env.Medias[0].ID, project.Manager.User.Avatar.ID)
		}
	}

}

func TestPreload_Slice_Many(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	users := []User{}
	for i := 1; i < 6; i++ {
		users = append(users, env.createUser(fmt.Sprintf("user%d", i)))
	}

	// Level 1
	{
		for _, user := range users {
			assert.Nil(t, user.Avatars)
		}

		queries, err := sqlxx.PreloadWithQueries(env.driver, &users, "Avatars")
		assert.NoError(t, err)
		assert.NotNil(t, queries)
		assert.Len(t, queries, 1)

		avatarQuery, ok := queries.ByTable("avatars")
		assert.True(t, ok)
		assert.Contains(t, avatarQuery.Query, "avatars.user_id IN (?, ?, ?, ?, ?)")
		assert.Len(t, avatarQuery.Args, 5)

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

	// Level 2
	{
		var (
			user     = env.createUser("spiderman")
			deadpool = env.createUser("deadpool")
			article1 = env.createArticle(&user)
			article2 = env.createArticle(&deadpool)
			articles = []Article{article1, article2}
		)

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
}
