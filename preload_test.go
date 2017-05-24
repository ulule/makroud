package sqlxx_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

// ----------------------------------------------------------------------------
// Errors
// ----------------------------------------------------------------------------

func TestPreload_Error_Unaddressable(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	article := Article{}

	queries, err := sqlxx.PreloadWithQueries(env.driver, article, "Author")
	is.Error(err)
	is.Nil(queries)
}

func TestPreload_Error_UnknownRelation(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	article := env.Articles[0]

	queries, err := sqlxx.PreloadWithQueries(env.driver, &article, "Foo")
	is.Error(err)
	is.Nil(queries)
	is.Zero(article.Author)
}

// ----------------------------------------------------------------------------
// Primary keys
// ----------------------------------------------------------------------------

func TestPreload_PrimaryKey_Null(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	category := env.createCategory("cat1", nil)

	queries, err := sqlxx.PreloadWithQueries(env.driver, category, "User")
	is.NoError(err)
	is.Nil(queries)
	is.Zero(category.User)

	category = env.createCategory("cat1", &(env.Users[0].ID))

	queries, err = sqlxx.PreloadWithQueries(env.driver, category, "User")
	is.NoError(err)
	is.NotNil(queries)
	is.Len(queries, 1)
	is.Contains(queries[0].Query, "FROM users WHERE users.id = ? LIMIT 1")
	is.Len(queries[0].Args, 1)
	is.Equal(category.UserID.Int64, queries[0].Args[0])
	is.NotZero(category.UserID)
	is.NotZero(category.User.ID)
}

// ----------------------------------------------------------------------------
// Single instance preloads
// ----------------------------------------------------------------------------

func TestPreload_Single_One(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	var (
		user    = env.createUser("batman")
		article = env.createArticle(user)
	)

	// level 1
	{
		// Value

		is.Zero(article.Author)

		queries, err := sqlxx.PreloadWithQueries(env.driver, article, "Author")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)

		userQuery, ok := queries.ByTable("users")
		is.True(ok)
		is.Contains(userQuery.Query, "FROM users WHERE users.id = ? LIMIT 1")
		is.Len(userQuery.Args, 1)
		is.EqualValues(article.AuthorID, userQuery.Args[0])
		is.NotZero(article.Author)
		is.Equal(user.ID, article.AuthorID)
		is.Equal(user.Username, article.Author.Username)

		// Pointer

		is.Nil(article.Reviewer)

		queries, err = sqlxx.PreloadWithQueries(env.driver, article, "Reviewer")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)

		userQuery, ok = queries.ByTable("users")
		is.True(ok)
		is.Contains(userQuery.Query, "FROM users WHERE users.id = ? LIMIT 1")
		is.Len(userQuery.Args, 1)
		is.EqualValues(article.ReviewerID, userQuery.Args[0])
		is.NotZero(article.Reviewer)
		is.Equal(user.ID, article.ReviewerID)
		is.Equal(user.Username, article.Reviewer.Username)
	}

	// level 2
	{
		// Value

		article.Author = User{}

		queries, err := sqlxx.PreloadWithQueries(env.driver, article, "Author", "Author.APIKey")
		is.NoError(err)
		is.NotNil(queries)

		authorQuery, ok := queries.ByTable("users")
		is.True(ok)
		is.Contains(authorQuery.Query, "WHERE users.id = ? LIMIT 1")
		is.Len(authorQuery.Args, 1)

		apikeyQuery, ok := queries.ByTable("api_keys")
		is.True(ok)
		is.Contains(apikeyQuery.Query, "WHERE api_keys.id = ? LIMIT 1")
		is.Len(apikeyQuery.Args, 1)

		is.NotZero(article.Author)
		is.Equal(user.ID, article.AuthorID)
		is.Equal(user.Username, article.Author.Username)

		is.NotZero(article.Author.APIKey)
		is.NotZero(article.Author.APIKey.ID)
		is.Equal(fmt.Sprintf("%s-apikey", user.Username), article.Author.APIKey.Key)

		// Pointer

		article.Author = User{}

		queries, err = sqlxx.PreloadWithQueries(env.driver, article, "Author", "Author.APIKeyPtr")
		is.NoError(err)
		is.NotNil(queries)

		authorQuery, ok = queries.ByTable("users")
		is.True(ok)
		is.Contains(authorQuery.Query, "WHERE users.id = ? LIMIT 1")
		is.Len(authorQuery.Args, 1)

		apikeyQuery, ok = queries.ByTable("api_keys")
		is.True(ok)
		is.Contains(apikeyQuery.Query, "WHERE api_keys.id = ? LIMIT 1")
		is.Len(apikeyQuery.Args, 1)

		is.NotZero(article.Author)
		is.Equal(user.ID, article.AuthorID)
		is.Equal(user.Username, article.Author.Username)

		is.NotNil(article.Author.APIKeyPtr)
		is.NotZero(article.Author.APIKeyPtr.ID)
		is.Equal(fmt.Sprintf("%s-apikey", user.Username), article.Author.APIKeyPtr.Key)
	}

	// Level 3
	{
		// Value

		article.Author = User{}

		queries, err := sqlxx.PreloadWithQueries(env.driver, article, "Author", "Author.APIKey", "Author.APIKey.Partner")
		is.NoError(err)
		is.NotNil(queries)

		authorQuery, ok := queries.ByTable("users")
		is.True(ok)
		is.Contains(authorQuery.Query, "WHERE users.id = ? LIMIT 1")
		is.Len(authorQuery.Args, 1)

		apikeyQuery, ok := queries.ByTable("api_keys")
		is.True(ok)
		is.Contains(apikeyQuery.Query, "WHERE api_keys.id = ? LIMIT 1")
		is.Len(apikeyQuery.Args, 1)

		partnerQuery, ok := queries.ByTable("partners")
		is.True(ok)
		is.Contains(partnerQuery.Query, "WHERE partners.id = ? LIMIT 1")
		is.Len(partnerQuery.Args, 1)

		is.NotZero(article.Author)
		is.Equal(user.ID, article.AuthorID)
		is.Equal(user.Username, article.Author.Username)

		is.NotZero(article.Author.APIKey)
		is.NotZero(article.Author.APIKey.ID)
		is.Equal(fmt.Sprintf("%s-apikey", user.Username), article.Author.APIKey.Key)

		is.NotZero(article.Author.APIKey.Partner)
		is.NotZero(article.Author.APIKey.Partner.ID)
	}
}

func TestPreload_Single_Many(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	user := env.createUser("wonderwoman")

	// Level 1
	{
		queries, err := sqlxx.PreloadWithQueries(env.driver, user, "Avatars")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)

		avatarQuery, ok := queries.ByTable("avatars")
		is.True(ok)
		is.Contains(avatarQuery.Query, "WHERE avatars.user_id = ?")
		is.Len(avatarQuery.Args, 1)
		is.EqualValues(user.ID, avatarQuery.Args[0])

		is.Len(user.Avatars, 5)
		for i, a := range user.Avatars {
			is.NotZero(a.ID)
			is.Equal(user.ID, a.UserID)
			is.Equal(fmt.Sprintf("/avatars/wonderwoman-%d.png", i+1), a.Path)
		}
	}

	// Level 2 - One
	{
		user.Avatars = []Avatar{}

		// Values

		queries, err := sqlxx.PreloadWithQueries(env.driver, user, "Avatars", "Avatars.Filter")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 2)

		avatarQuery, ok := queries.ByTable("avatars")
		is.True(ok)
		is.Contains(avatarQuery.Query, "WHERE avatars.user_id = ?")
		is.Len(avatarQuery.Args, 1)
		is.EqualValues(user.ID, avatarQuery.Args[0])

		avatarFilterQuery, ok := queries.ByTable("avatar_filters")
		is.True(ok)
		is.Contains(avatarFilterQuery.Query, "avatar_filters.id IN (?, ?, ?, ?, ?)")
		is.Len(avatarFilterQuery.Args, 5)

		is.Len(user.Avatars, 5)
		for i, a := range user.Avatars {
			is.NotZero(a.ID)
			is.Equal(user.ID, a.UserID)
			is.Equal(fmt.Sprintf("/avatars/wonderwoman-%d.png", i+1), a.Path)
			is.NotZero(a.Filter)
			is.NotZero(a.Filter.ID)
		}

		// Pointers

		user.Avatars = []Avatar{}

		queries, err = sqlxx.PreloadWithQueries(env.driver, user, "Avatars", "Avatars.FilterPtr")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 2)

		avatarQuery, ok = queries.ByTable("avatars")
		is.True(ok)
		is.Contains(avatarQuery.Query, "WHERE avatars.user_id = ?")
		is.Len(avatarQuery.Args, 1)
		is.EqualValues(user.ID, avatarQuery.Args[0])

		avatarFilterQuery, ok = queries.ByTable("avatar_filters")
		is.True(ok)
		is.Contains(avatarFilterQuery.Query, "avatar_filters.id IN (?, ?, ?, ?, ?)")
		is.Len(avatarFilterQuery.Args, 5)

		is.Len(user.Avatars, 5)
		for i, a := range user.Avatars {
			is.NotZero(a.ID)
			is.Equal(user.ID, a.UserID)
			is.Equal(fmt.Sprintf("/avatars/wonderwoman-%d.png", i+1), a.Path)
			is.NotNil(a.FilterPtr)
			is.NotZero(a.FilterPtr.ID)
		}
	}
}

func TestPreload_Slice_One(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	user := env.createUser("batman")
	articles := []Article{}
	for i := 0; i < 5; i++ {
		article := env.createArticle(user)
		articles = append(articles, *article)
	}

	// Level 1
	{
		// Value

		for i := range articles {
			is.Zero(articles[i].Author)
		}

		queries, err := sqlxx.PreloadWithQueries(env.driver, &articles, "Author")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		is.Contains(queries[0].Query, "FROM users WHERE users.id IN (?)")
		is.Len(queries[0].Args, 1)

		for _, article := range articles {
			is.Equal(user.ID, article.AuthorID)
			is.Equal(user.Username, article.Author.Username)
			is.EqualValues(article.AuthorID, queries[0].Args[0])
		}

		// Pointer

		for i := range articles {
			is.Nil(articles[i].Reviewer)
		}

		queries, err = sqlxx.PreloadWithQueries(env.driver, &articles, "Reviewer")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		is.Contains(queries[0].Query, "FROM users WHERE users.id IN (?)")
		is.Len(queries[0].Args, 1)

		for _, a := range articles {
			is.Equal(user.ID, a.ReviewerID)
			is.Equal(user.Username, a.Reviewer.Username)
			is.EqualValues(a.ReviewerID, queries[0].Args[0])
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
			*batman,
			*robin,
			*catwoman,
		}

		// user_id => media_id
		avatars := map[int]int{}
		for _, user := range users {
			avatars[user.ID] = int(user.AvatarID.Int64)
		}

		queries, err := sqlxx.PreloadWithQueries(env.driver, &users, "Avatar")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)

		avatarQuery, ok := queries.ByTable("media")
		is.True(ok)
		is.Contains(avatarQuery.Query, "WHERE media.id IN (?, ?, ?)")
		is.Len(avatarQuery.Args, 3)

		for _, user := range users {
			avatar := user.Avatar
			is.NotNil(avatar)
			is.Equal(avatar.ID, avatars[user.ID])
			is.Contains(queries[0].Args, int64(avatar.ID))
		}
	}

	// Level 1 - Multiple paths (+ mix of value and pointers)
	{
		var (
			batman   = env.createUser("batman")
			robin    = env.createUser("robin")
			catwoman = env.createUser("catwoman")
			article1 = env.createArticle(batman)
			article2 = env.createArticle(robin)
			article3 = env.createArticle(catwoman)
		)

		articles := []Article{
			*article1,
			*article2,
			*article3,
		}

		queries, err := sqlxx.PreloadWithQueries(env.driver, &articles, "Author", "Reviewer")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 2)

		userQuery, ok := queries.ByTable("users")
		is.True(ok)
		is.Contains(userQuery.Query, "WHERE users.id IN")
		is.Len(userQuery.Args, 3)

		table := []struct {
			user    User
			article Article
		}{
			{*batman, articles[0]},
			{*robin, articles[1]},
			{*catwoman, articles[2]},
		}

		for _, tt := range table {
			is.Equal(tt.user.ID, tt.article.AuthorID)
			is.Equal(tt.user, tt.article.Author)
			is.Equal(tt.user.ID, tt.article.ReviewerID)
			is.Equal(tt.article.Reviewer, &tt.user)
		}
	}

	// Level 2 and more
	{

		var (
			projects = env.Projects
			project  = projects[0]
		)

		is.Len(projects, 1)
		is.Nil(project.Manager)

		queries, err := sqlxx.PreloadWithQueries(env.driver, &projects, "Manager", "Manager.User", "Manager.User.Avatar")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 3)

		managerQuery, ok := queries.ByTable("managers")
		is.True(ok)
		is.Contains(managerQuery.Query, "WHERE managers.id IN")
		is.Len(managerQuery.Args, 1)
		is.Equal(int64(env.Projects[0].ManagerID), managerQuery.Args[0])

		userQuery, ok := queries.ByTable("users")
		is.True(ok)
		is.Contains(userQuery.Query, "WHERE users.id IN")
		is.Len(userQuery.Args, 1)
		is.Equal(int64(env.Users[0].ID), userQuery.Args[0])

		mediaQuery, ok := queries.ByTable("media")
		is.True(ok)
		is.Contains(mediaQuery.Query, "WHERE media.id IN")
		is.Len(mediaQuery.Args, 1)
		is.Equal(int64(env.Medias[0].ID), mediaQuery.Args[0])

		for _, project := range projects {
			is.NotNil(project.Manager)
			is.Equal(env.Managers[0].ID, project.Manager.ID)

			is.NotNil(project.Manager.User)
			is.Equal(env.Users[0].ID, project.Manager.User.ID)

			is.NotNil(project.Manager.User.Avatar)
			is.Equal(env.Medias[0].ID, project.Manager.User.Avatar.ID)
		}
	}

}

func TestPreload_Slice_Many(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	users := []User{}
	for i := 1; i < 6; i++ {
		user := env.createUser(fmt.Sprintf("user%d", i))
		users = append(users, *user)
	}

	// Level 1
	{
		for _, user := range users {
			is.Nil(user.Avatars)
		}

		queries, err := sqlxx.PreloadWithQueries(env.driver, &users, "Avatars")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)

		avatarQuery, ok := queries.ByTable("avatars")
		is.True(ok)
		is.Contains(avatarQuery.Query, "avatars.user_id IN (?, ?, ?, ?, ?)")
		is.Len(avatarQuery.Args, 5)

		for _, user := range users {
			is.NotZero(user.Avatars)
			is.Contains(queries[0].Args, int64(user.ID))

			for _, avatar := range user.Avatars {
				is.NotZero(avatar.ID)
				is.Equal(user.ID, avatar.UserID)
				is.Equal(user.ID, avatar.UserID)
				is.True(strings.HasPrefix(avatar.Path, fmt.Sprintf("/avatars/%s-", user.Username)))
			}
		}
	}

	// Level 2
	{
		var (
			spiderman = env.createUser("spiderman")
			deadpool  = env.createUser("deadpool")
			article1  = env.createArticle(spiderman)
			article2  = env.createArticle(deadpool)
			articles  = []Article{*article1, *article2}
		)

		queries, err := sqlxx.PreloadWithQueries(env.driver, &articles, "Author", "Author.APIKey")
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 2)

		authorQuery, ok := queries.ByTable("users")
		is.True(ok)
		is.Contains(authorQuery.Query, "WHERE users.id IN (?, ?)")

		apikeyQuery, ok := queries.ByTable("api_keys")
		is.True(ok)
		is.Contains(apikeyQuery.Query, "WHERE api_keys.id IN (?, ?)")

		is.Equal(spiderman.ID, articles[0].Author.ID)
		is.Equal(spiderman.ID, articles[0].AuthorID)

		is.Equal(spiderman.Username, articles[0].Author.Username)
		is.NotZero(articles[0].Author.APIKeyID)
		is.Equal("spiderman-apikey", articles[0].Author.APIKey.Key)

		is.Equal(deadpool.ID, articles[1].Author.ID)
		is.Equal(deadpool.ID, articles[1].AuthorID)
		is.Equal(deadpool.Username, articles[1].Author.Username)

		is.NotZero(articles[1].Author.APIKeyID)
		is.Equal("deadpool-apikey", articles[1].Author.APIKey.Key)
	}
}
