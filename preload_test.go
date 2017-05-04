package sqlxx_test

import (
	"fmt"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestPreload_Unaddressable(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	article := Article{}
	assert.NotNil(t, sqlxx.Preload(db, article, "Author"))
}

func TestPreload_UnknownRelation(t *testing.T) {
	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	article := fixtures.Articles[0]
	assert.Error(t, sqlxx.Preload(db, &article, "Foo"))
	assert.Zero(t, article.Author)
}

func TestPreload_NullPrimaryKey(t *testing.T) {
	db, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	category := createCategory(t, db, "cat1", nil)
	assert.NoError(t, sqlxx.Preload(db, &category, "User"))
	assert.Zero(t, category.User)

	category = createCategory(t, db, "cat1", &fixtures.User.ID)
	assert.NoError(t, sqlxx.Preload(db, &category, "User"))
	assert.NotZero(t, category.UserID)
	assert.NotZero(t, category.User.ID)
}

// ----------------------------------------------------------------------------
// Preloads: OneToOne
// ----------------------------------------------------------------------------

func TestPreload_OneToOne_Level1(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	batman := createUser(t, db, "batman")
	article := createArticle(t, db, &batman)

	//
	// Instance
	//

	// Value
	assert.Nil(t, sqlxx.Preload(db, &article, "Author"))
	assert.NotZero(t, article.Author)
	assert.Equal(t, batman.ID, article.AuthorID)
	assert.Equal(t, batman.Username, article.Author.Username)

	// Pointer
	assert.Nil(t, sqlxx.Preload(db, &article, "Reviewer"))
	assert.NotZero(t, article.Reviewer)
	assert.Equal(t, batman.ID, article.ReviewerID)
	assert.Equal(t, batman.Username, article.Reviewer.Username)
}

func TestPreload_ManyToOne_Level1_Same(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	batman := createUser(t, db, "batman")

	var articles []Article
	for i := 0; i < 5; i++ {
		articles = append(articles, createArticle(t, db, &batman))
	}

	// Value
	assert.Nil(t, sqlxx.Preload(db, &articles, "Author"))
	for _, a := range articles {
		assert.Equal(t, batman.ID, a.AuthorID)
		assert.Equal(t, batman.Username, a.Author.Username)
	}

	// Pointer
	assert.Nil(t, sqlxx.Preload(db, &articles, "Reviewer"))
	for _, a := range articles {
		assert.Equal(t, batman.ID, a.ReviewerID)
		assert.Equal(t, batman.Username, a.Reviewer.Username)
	}
}

func TestPreload_ManyToOne_Level1_Different_Pointer_Null(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	batman := createUser(t, db, "batman")
	robin := createUser(t, db, "robin")
	catwoman := createUser(t, db, "catwoman")

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

	assert.Nil(t, sqlxx.Preload(db, &users, "Avatar"))

	for i, _ := range users {
		assert.NotNil(t, users[i].Avatar)
		assert.Equal(t, users[i].Avatar.ID, avatars[users[i].ID])
	}
}

func TestPreload_ManyToOne_Level1_Different(t *testing.T) {
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

	assert.Nil(t, sqlxx.Preload(db, &articles, "Author", "Reviewer"))
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

func TestPreload_OneToOne_Level2(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := createUser(t, db, "spiderman")

	article := createArticle(t, db, &user)
	assert.Nil(t, sqlxx.Preload(db, &article, "Author", "Author.APIKey"))
	assert.NotZero(t, article.Author)
	assert.NotZero(t, article.Author.APIKey)
	assert.Equal(t, user.ID, article.AuthorID)
	assert.Equal(t, user.Username, article.Author.Username)
	assert.NotZero(t, article.Author.APIKey.ID)
	assert.Equal(t, "spiderman-apikey", article.Author.APIKey.Key)
}

func TestPreload_OneToOne_Level2_MultipleEither(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := createUser(t, db, "spiderman")
	assert.NotEmpty(t, user)
	assert.Nil(t, sqlxx.Preload(db, &user, "Avatar"))
	assert.NotNil(t, user.Avatar)

	article := createArticle(t, db, &user)
	assert.NotEmpty(t, article)

	comment := createComment(t, db, &user, &article)
	assert.NotEmpty(t, comment)

	comments := []Comment{comment}

	// Preload
	assert.Nil(t, sqlxx.Preload(db, &comments, "User", "User.Avatar"))
	comment = comments[0]

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

func TestPreload_ManyToOne_Level2_Multiple(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := createUser(t, db, "spiderman")
	article := createArticle(t, db, &user)

	deadpool := createUser(t, db, "deadpool")
	article2 := createArticle(t, db, &deadpool)

	articles := []Article{article, article2}

	assert.Nil(t, sqlxx.Preload(db, &articles, "Author", "Author.APIKey"))

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

// ----------------------------------------------------------------------------
// Preloads: ToMany
// ----------------------------------------------------------------------------

func TestPreload_OneToMany_Level1_Simple(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	user := createUser(t, db, "wonderwoman")
	assert.Nil(t, sqlxx.Preload(db, &user, "Avatars"))
	assert.Len(t, user.Avatars, 5)

	for i, a := range user.Avatars {
		assert.NotZero(t, a.ID)
		assert.Equal(t, user.ID, a.UserID)
		assert.Equal(t, fmt.Sprintf("/avatars/wonderwoman-%d.png", i+1), a.Path)
	}
}

func TestPreload_ManyToMany_Level1(t *testing.T) {
	db, _, shutdown := dbConnection(t)
	defer shutdown()

	users := []User{}
	for i := 1; i < 6; i++ {
		users = append(users, createUser(t, db, fmt.Sprintf("user%d", i)))
	}

	for _, user := range users {
		assert.Zero(t, user.Avatars)
	}

	assert.Nil(t, sqlxx.Preload(db, &users, "Avatars"))

	for _, user := range users {
		assert.NotZero(t, user.Avatars)
		for _, avatar := range user.Avatars {
			assert.NotZero(t, avatar.ID)
			assert.Equal(t, user.ID, avatar.UserID)
			assert.Equal(t, user.ID, avatar.UserID)
			assert.True(t, strings.HasPrefix(avatar.Path, fmt.Sprintf("/avatars/%s-", user.Username)))
		}
	}
}
