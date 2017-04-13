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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var dbDefaultParams = map[string]string{
	"USER":     "postgres",
	"PASSWORD": "",
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
	DROP TABLE IF EXISTS tags CASCADE;
	DROP TABLE IF EXISTS articles CASCADE;
	DROP TABLE IF EXISTS articles_categories CASCADE;
	DROP TABLE IF EXISTS partners CASCADE;
	DROP TABLE IF EXISTS media CASCADE;
	DROP TABLE IF EXISTS projects CASCADE;
	DROP TABLE IF EXISTS managers CASCADE;
`

var dbSchema = `
CREATE TABLE api_keys (
	id 		serial primary key not null,
	partner_id 	integer,
	key 		varchar(255) not null
);

CREATE TABLE partners (
	id 		serial primary key not null,
	name		varchar(255) not null
);

CREATE TABLE managers (
	id 		serial primary key not null,
	name		varchar(255) not null,
	user_id 	integer
);

CREATE TABLE projects (
	id 		serial primary key not null,
	name		varchar(255) not null,
	manager_id 	integer,
	user_id 	integer
);

CREATE TABLE users (
	id 		serial primary key not null,
	username 	 varchar(30) not null,
	is_active 	boolean default true,
	api_key_id	integer,
	avatar_id	integer,
	created_at 	timestamp with time zone default current_timestamp,
	updated_at 	timestamp with time zone default current_timestamp,
	deleted_at 	timestamp with time zone
);

CREATE TABLE profiles (
	id 		serial primary key not null,
	user_id		integer references users(id),
	first_name 	varchar(255) not null,
	last_name 	varchar(255) not null
);

CREATE TABLE media (
	id		serial primary key not null,
	path 		varchar(255) not null,
	created_at	timestamp with time zone default current_timestamp,
	updated_at	timestamp with time zone default current_timestamp
);

CREATE TABLE avatars (
	id 		serial primary key not null,
	path 		varchar(255) not null,
	user_id 	integer references users(id),
	created_at 	timestamp with time zone default current_timestamp,
	updated_at 	timestamp with time zone default current_timestamp
);

CREATE TABLE tags (
	id 		serial primary key not null,
	name 		varchar(255) not null
);

CREATE TABLE articles (
	id 		serial primary key not null,
	title 		varchar(255) not null,
	author_id 	integer references users(id),
	reviewer_id 	integer references users(id),
	main_tag_id 	integer references tags(id),
	is_published 	boolean default true,
	created_at 	timestamp with time zone default current_timestamp,
	updated_at 	timestamp with time zone default current_timestamp
);

CREATE TABLE comments (
	id 		serial primary key not null,
	user_id		integer references users(id),
	article_id	integer references articles(id),
	content		text,
	created_at 	timestamp with time zone default current_timestamp,
	updated_at 	timestamp with time zone default current_timestamp
);

CREATE TABLE categories (
	id 		serial primary key not null,
	name 		varchar(255) not null,
	user_id 	integer references users(id)
);

CREATE TABLE articles_categories (
	id 		serial primary key not null,
	article_id 	integer references articles(id),
	category_id 	integer references categories(id)
);`

type TestData struct {
	User               User
	APIKeys            []APIKey
	Profiles           []Profile
	Avatars            []Avatar
	Articles           []Article
	Categories         []Category
	Tags               []Tag
	ArticlesCategories []ArticleCategory
	Partners           []Partner
	Managers           []Manager
	Projects           []Project
}

type Partner struct {
	ID   int    `db:"id" sqlxx:"primary_key:true; ignored:true"`
	Name string `db:"name"`
}

func (Partner) TableName() string { return "partners" }

type Manager struct {
	ID     int    `db:"id" sqlxx:"primary_key:true; ignored:true"`
	Name   string `db:"name"`
	UserID int    `db:"user_id"`
	User   *User
}

func (Manager) TableName() string { return "managers" }

type Project struct {
	ID        int    `db:"id" sqlxx:"primary_key:true; ignored:true"`
	Name      string `db:"name"`
	ManagerID int    `db:"manager_id"`
	UserID    int    `db:"user_id"`
	Manager   *Manager
	User      *User
}

func (Project) TableName() string { return "projects" }

type APIKey struct {
	ID        int    `db:"id" sqlxx:"primary_key:true; ignored:true"`
	Key       string `db:"key"`
	Partner   Partner
	PartnerID int `db:"partner_id"`
}

func (APIKey) TableName() string { return "api_keys" }

type Media struct {
	ID        int       `db:"id" sqlxx:"primary_key:true"`
	Path      string    `db:"path"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (Media) TableName() string { return "media" }

type User struct {
	ID       int    `db:"id" sqlxx:"primary_key:true; ignored:true"`
	Username string `db:"username"`
	IsActive bool   `db:"is_active" sqlxx:"default:true"`

	CreatedAt time.Time  `db:"created_at" sqlxx:"auto_now_add:true"`
	UpdatedAt time.Time  `db:"updated_at" sqlxx:"default:now()"`
	DeletedAt *time.Time `db:"deleted_at"`

	APIKeyID int `db:"api_key_id"`
	APIKey   APIKey
	AvatarID sql.NullInt64 `db:"avatar_id"`
	Avatar   *Media

	Avatars []Avatar
	// Comments []Comment
	Profile Profile
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
	ID     int           `db:"id" sqlxx:"primary_key:true"`
	Name   string        `db:"name"`
	UserID sql.NullInt64 `db:"user_id"`
	User   User
}

func (Category) TableName() string { return "categories" }

// This model has a different ID type.
type Tag struct {
	ID   uint   `db:"id"`
	Name string `db:"name"`
}

func (Tag) TableName() string { return "tags" }

type Article struct {
	ID          int       `db:"id" sqlxx:"primary_key:true"`
	Title       string    `db:"title"`
	AuthorID    int       `db:"author_id"`
	ReviewerID  int       `db:"reviewer_id"`
	IsPublished bool      `db:"is_published"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Author      User
	Reviewer    *User
	MainTagID   sql.NullInt64 `db:"main_tag_id"`
	MainTag     *Tag
}

func (Article) TableName() string { return "articles" }

type ArticleCategory struct {
	ID         int `db:"id" sqlxx:"primary_key:true"`
	ArticleID  int `db:"article_id"`
	CategoryID int `db:"category_id"`
}

func (ArticleCategory) TableName() string { return "articles_categories" }

// ----------------------------------------------------------------------------
// Models with different IDs
// ----------------------------------------------------------------------------

func dbParam(param string) string {
	param = strings.ToUpper(param)

	if v := os.Getenv(fmt.Sprintf("DB_%s", param)); len(v) != 0 {
		return v
	}

	return dbDefaultParams[param]
}

func loadData(t *testing.T, driver Driver) *TestData {
	// Partners
	driver.MustExec("INSERT INTO partners (name) VALUES ($1)", "Ulule")
	partners := []Partner{}
	require.NoError(t, driver.Select(&partners, "SELECT * FROM partners"))
	partner := partners[0]

	// API Keys
	driver.MustExec("INSERT INTO api_keys (key, partner_id) VALUES ($1, $2)", "this-is-my-scret-api-key", partner.ID)
	apiKeys := []APIKey{}
	require.NoError(t, driver.Select(&apiKeys, "SELECT * FROM api_keys"))
	apiKey := apiKeys[0]

	// Media
	driver.MustExec("INSERT INTO media (path) VALUES ($1)", "media/avatar.png")
	media := Media{}
	require.NoError(t, driver.Get(&media, "SELECT * FROM media LIMIT 1"))

	// Users
	driver.MustExec("INSERT INTO users (username, api_key_id, avatar_id) VALUES ($1, $2, $3)", "jdoe", apiKey.ID, media.ID)
	user := User{}
	require.NoError(t, driver.Get(&user, "SELECT * FROM users WHERE username=$1", "jdoe"))

	// Managers
	driver.MustExec("INSERT INTO managers (name, user_id) VALUES ($1, $2)", "Super Owl", user.ID)
	managers := []Manager{}
	require.NoError(t, driver.Select(&managers, "SELECT * FROM managers"))
	manager := managers[0]

	// Projects
	driver.MustExec("INSERT INTO projects (name, manager_id, user_id) VALUES ($1, $2, $3)", "Super Project", manager.ID, user.ID)
	projects := []Project{}
	require.NoError(t, driver.Select(&projects, "SELECT * FROM projects"))

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

	// Tags
	driver.MustExec("INSERT INTO tags (name) VALUES ($1)", "Tag")
	tags := []Tag{}
	require.NoError(t, driver.Select(&tags, "SELECT * FROM tags"))
	tag := tags[0]

	// Articles
	for i := 0; i < 5; i++ {
		driver.MustExec("INSERT INTO articles (title, author_id, reviewer_id, main_tag_id) VALUES ($1, $2, $3, $4)", fmt.Sprintf("Title #%d", i), user.ID, user.ID, tag.ID)
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
		Tags:               tags,
		Articles:           articles,
		ArticlesCategories: articlesCategories,
		Partners:           partners,
		Managers:           managers,
		Projects:           projects,
	}
}

func createArticle(t *testing.T, driver Driver, user *User) Article {
	is := assert.New(t)

	var id int
	err := driver.QueryRowx("INSERT INTO articles (title, author_id, reviewer_id) VALUES ($1, $2, $3) RETURNING id", "Title", user.ID, user.ID).Scan(&id)
	is.Nil(err)

	article := Article{}
	require.NoError(t, driver.Get(&article, "SELECT * FROM articles WHERE id = $1", id))

	return article
}

func createComment(t *testing.T, driver Driver, user *User, article *Article) Comment {
	var id int
	err := driver.QueryRowx("INSERT INTO comments (content, user_id, article_id) VALUES ($1, $2, $3) RETURNING id", "Lorem Ipsum", user.ID, article.ID).Scan(&id)
	require.Nil(t, err)

	comment := Comment{}
	require.NoError(t, driver.Get(&comment, "SELECT * FROM comments WHERE id = $1", id))

	return comment
}

func createUser(t *testing.T, driver Driver, username string) User {
	key := fmt.Sprintf("%s-apikey", username)

	name := fmt.Sprintf("%s-partner", username)

	driver.MustExec("INSERT INTO partners (name) VALUES ($1)", name)
	partner := Partner{}
	require.NoError(t, driver.Get(&partner, "SELECT * FROM partners WHERE name = $1", name))

	driver.MustExec("INSERT INTO media (path) VALUES ($1)", fmt.Sprintf("media/media-%s.png", username))
	media := Media{}
	require.NoError(t, driver.Get(&media, "SELECT * FROM media ORDER BY id DESC LIMIT 1"))

	driver.MustExec("INSERT INTO api_keys (key, partner_id) VALUES ($1, $2)", key, partner.ID)
	apiKey := APIKey{}
	require.NoError(t, driver.Get(&apiKey, "SELECT * FROM api_keys WHERE key = $1", key))

	driver.MustExec("INSERT INTO users (username, api_key_id, avatar_id) VALUES ($1, $2, $3)", username, apiKey.ID, media.ID)
	user := User{}
	require.NoError(t, driver.Get(&user, "SELECT * FROM users WHERE username=$1", username))

	for i := 1; i < 6; i++ {
		driver.MustExec("INSERT INTO avatars (path, user_id) VALUES ($1, $2)", fmt.Sprintf("/avatars/%s-%d.png", username, i), user.ID)
	}

	avatars := []Avatar{}
	require.NoError(t, driver.Select(&avatars, "SELECT * FROM avatars"))

	return user
}

func createCategory(t *testing.T, driver Driver, name string, userID *int) Category {
	driver.MustExec("INSERT INTO categories (name) VALUES ($1)", name)

	if userID != nil {
		driver.MustExec("UPDATE categories SET user_id=$1 WHERE name=$2", *userID, name)
	}

	category := Category{}
	require.NoError(t, driver.Get(&category, "SELECT * FROM categories WHERE name=$1", name))

	return category
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
