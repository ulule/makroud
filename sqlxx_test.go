package sqlxx

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

var dbDefaultParams = map[string]string{
	"USER":     "postgres",
	"PASSWORS": "",
	"HOST":     "localhost",
	"PORT":     "5432",
	"NAME":     "sqlxx_test",
}

var dropTables = `
	DROP TABLE IF EXISTS users CASCADE;
	DROP TABLE IF EXISTS api_keys CASCADE;
	DROP TABLE IF EXISTS profiles CASCADE;
	DROP TABLE IF EXISTS comments CASCADE;
	DROP TABLE IF EXISTS avatars CASCADE;
	DROP TABLE IF EXISTS categories CASCADE;
	DROP TABLE IF EXISTS articles CASCADE;
	DROP TABLE IF EXISTS articles_categories CASCADE;
`

var dbSchema = `CREATE TABLE api_keys (
	id 	serial primary key not null,
	key varchar(255) not null
);

CREATE TABLE users (
	id 				serial primary key not null,
	username 	    varchar(30) not null,
	is_active 		boolean default true,
	api_key_id		integer,
    created_at 		timestamp with time zone default current_timestamp,
    updated_at 		timestamp with time zone default current_timestamp,
    deleted_at 		timestamp with time zone
);

CREATE TABLE profiles (
	id 				serial primary key not null,
	user_id 		integer references users(id),
	first_name 		varchar(255) not null,
	last_name 		varchar(255) not null
);

CREATE TABLE avatars (
	id 				serial primary key not null,
	path 			varchar(255) not null,
	user_id 		integer references users(id),
    created_at 		timestamp with time zone default current_timestamp,
    updated_at 		timestamp with time zone default current_timestamp
);

CREATE TABLE articles (
	id 				serial primary key not null,
	title 			varchar(255) not null,
	author_id 		integer references users(id),
	is_published 	boolean default true,
    created_at 		timestamp with time zone default current_timestamp,
    updated_at 		timestamp with time zone default current_timestamp
);

CREATE TABLE comments (
	id 				serial primary key not null,
	user_id			integer references users(id),
	article_id		integer references articles(id),
	content			text,
    created_at 		timestamp with time zone default current_timestamp,
    updated_at 		timestamp with time zone default current_timestamp
);

CREATE TABLE categories (
	id 				serial primary key not null,
	name 			varchar(255) not null,
	user_id 		integer references users(id)
);

CREATE TABLE articles_categories (
	id 				serial primary key not null,
	article_id 		integer references articles(id),
	category_id 	integer references categories(id)
);`

type TestData struct {
	User               User
	APIKeys            []APIKey
	Profiles           []Profile
	Avatars            []Avatar
	Articles           []Article
	Categories         []Category
	ArticlesCategories []ArticleCategory
}

type APIKey struct {
	ID  int    `db:"id" sqlxx:"primary_key:true; ignored:true"`
	Key string `db:"key"`
}

func (APIKey) TableName() string { return "api_keys" }

type User struct {
	ID       int    `db:"id" sqlxx:"primary_key:true; ignored:true"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active" sqlxx:"default:true"`

	CreatedAt time.Time  `db:"created_at" sqlxx:"auto_now_add:true"`
	UpdatedAt time.Time  `db:"updated_at" sqlxx:"default:now()"`
	DeletedAt *time.Time `db:"deleted_at"`

	APIKeyID int `db:"api_key_id"`
	APIKey   APIKey

	Avatars  []Avatar
	Comments []Comment
	Profile  Profile
}

func (User) TableName() string { return "users" }

type Comment struct {
	ID        int `db:"id" sqlxx:"primary_key:true; ignored:true"`
	UserID    int `db:"user_id"`
	User      User
	ArticleID int `db:"article_id"`
	Article   Article
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at" sqlxx:"auto_now_add:true"`
	UpdatedAt time.Time `db:"updated_at" sqlxx:"default:now()"`
}

func (Comment) TableName() string { return "comments" }

type Profile struct {
	ID        int    `db:"id" sqlxx:"primary_key:true; ignored:true"`
	UserID    int    `db:"user_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

func (Profile) TableName() string { return "profiles" }

type Avatar struct {
	ID        int       `db:"id" sqlxx:"primary_key:true"`
	Path      string    `db:"path"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (Avatar) TableName() string { return "avatars" }

type Category struct {
	ID     int    `db:"id" sqlxx:"primary_key:true"`
	Name   string `db:"name"`
	UserID int    `db:"user_id"`
	User   User
}

func (Category) TableName() string { return "categories" }

type Article struct {
	ID          int       `db:"id" sqlxx:"primary_key:true"`
	Title       string    `db:"title"`
	AuthorID    int       `db:"author_id"`
	IsPublished bool      `db:"is_published"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Author      User
	// Categories  []Category
}

func (Article) TableName() string { return "articles" }

type ArticleCategory struct {
	ID         int `db:"id" sqlxx:"primary_key:true"`
	ArticleID  int `db:"article_id"`
	CategoryID int `db:"category_id"`
}

func (ArticleCategory) TableName() string { return "articles_categories" }

func dbParam(param string) string {
	param = strings.ToUpper(param)

	if v := os.Getenv(fmt.Sprintf("DB_%s", param)); len(v) != 0 {
		return v
	}

	return dbDefaultParams[param]
}

func loadData(t *testing.T, driver Driver) *TestData {
	// API Keys
	driver.MustExec("INSERT INTO api_keys (key) VALUES ($1)", "this-is-my-scret-api-key")
	apiKeys := []APIKey{}
	require.NoError(t, driver.Select(&apiKeys, "SELECT * FROM api_keys"))
	apiKey := apiKeys[0]

	// Users
	driver.MustExec("INSERT INTO users (username, api_key_id) VALUES ($1, $2)", "jdoe", apiKey.ID)
	user := User{}
	require.NoError(t, driver.Get(&user, "SELECT * FROM users WHERE username=$1", "jdoe"))

	// Avatars
	for i := 0; i < 5; i++ {
		driver.MustExec("INSERT INTO avatars (path, user_id) VALUES ($1, $2)", fmt.Sprintf("/avatars/%s-%d.png", user.Username, i), user.ID)
	}
	avatars := []Avatar{}
	require.NoError(t, driver.Select(&avatars, "SELECT * FROM avatars"))

	// Profiles
	driver.MustExec("INSERT INTO profiles (user_id, first_name, last_name) VALUES ($1, $2, $3)", user.ID, "John", "Doe")
	profiles := []Profile{}
	require.NoError(t, driver.Select(&profiles, "SELECT * FROM profiles"))

	// Categories
	for i := 0; i < 5; i++ {
		driver.MustExec("INSERT INTO categories (name, user_id) VALUES ($1, $2)", fmt.Sprintf("Category #%d", i), user.ID)
	}
	categories := []Category{}
	require.NoError(t, driver.Select(&categories, "SELECT * FROM categories"))

	// Articles
	for i := 0; i < 5; i++ {
		driver.MustExec("INSERT INTO articles (title, author_id) VALUES ($1, $2)", fmt.Sprintf("Title #%d", i), user.ID)
	}
	articles := []Article{}
	require.NoError(t, driver.Select(&articles, "SELECT * FROM articles"))

	// Articles <-> Categories
	for _, article := range articles {
		for _, category := range categories {
			driver.MustExec("INSERT INTO articles_categories (article_id, category_id) VALUES ($1, $2)", article.ID, category.ID)
		}
	}
	articlesCategories := []ArticleCategory{}
	require.NoError(t, driver.Select(&articlesCategories, "SELECT * FROM articles_categories"))

	return &TestData{
		APIKeys:            apiKeys,
		User:               user,
		Profiles:           profiles,
		Avatars:            avatars,
		Categories:         categories,
		Articles:           articles,
		ArticlesCategories: articlesCategories,
	}
}

func createUser(t *testing.T, driver Driver, username string) User {
	key := fmt.Sprintf("%s-apikey", username)

	driver.MustExec("INSERT INTO api_keys (key) VALUES ($1)", key)
	apiKey := APIKey{}
	require.NoError(t, driver.Get(&apiKey, "SELECT * FROM api_keys WHERE key=$1", key))

	driver.MustExec("INSERT INTO users (username, api_key_id) VALUES ($1, $2)", username, apiKey.ID)
	user := User{}
	require.NoError(t, driver.Get(&user, "SELECT * FROM users WHERE username=$1", username))

	for i := 1; i < 6; i++ {
		driver.MustExec("INSERT INTO avatars (path, user_id) VALUES ($1, $2)", fmt.Sprintf("/avatars/%s-%d.png", username, i), user.ID)
	}

	avatars := []Avatar{}
	require.NoError(t, driver.Select(&avatars, "SELECT * FROM avatars"))

	return user
}

func dbConnection(t *testing.T) (*sqlx.DB, *TestData, func()) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable;timezone=UTC",
		dbParam("user"),
		dbParam("password"),
		dbParam("host"),
		dbParam("port"),
		dbParam("name")))

	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "postgres")
	dbx.MustExec(dropTables)
	dbx.MustExec(dbSchema)

	return dbx, loadData(t, dbx), func() {
		if value := os.Getenv("KEEP_DB"); len(value) == 0 {
			dbx.MustExec(dropTables)
		}

		require.NoError(t, db.Close())
	}
}
