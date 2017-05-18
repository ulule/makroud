package sqlxx_test

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/lib/pq"
	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

var dbDefaultOptions = map[string]sqlxx.Option{
	"USER":     sqlxx.User("postgres"),
	"PASSWORD": sqlxx.Password(""),
	"HOST":     sqlxx.Host("localhost"),
	"PORT":     sqlxx.Port(5432),
	"NAME":     sqlxx.Database("sqlxx_test"),
}

var dropTables = `
	DROP TABLE IF EXISTS users CASCADE;
	DROP TABLE IF EXISTS api_keys CASCADE;
	DROP TABLE IF EXISTS profiles CASCADE;
	DROP TABLE IF EXISTS comments CASCADE;
	DROP TABLE IF EXISTS avatar_filters CASCADE;
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
	id          serial primary key not null,
	username    varchar(30) not null,
	is_active   boolean default true,
	api_key_id  integer,
	avatar_id   integer,
	created_at  timestamp with time zone default current_timestamp,
	updated_at  timestamp with time zone default current_timestamp,
	deleted_at  timestamp with time zone
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

CREATE TABLE tags (
	id 		serial primary key not null,
	name 		varchar(255) not null
);

CREATE TABLE avatar_filters (
	id 		serial primary key not null,
	name 		varchar(255) not null
);

CREATE TABLE avatars (
	id 		serial primary key not null,
	path 		varchar(255) not null,
	user_id 	integer references users(id),
	filter_id	integer references avatar_filters(id),
	created_at 	timestamp with time zone default current_timestamp,
	updated_at 	timestamp with time zone default current_timestamp
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

	APIKeyID  int `db:"api_key_id"`
	APIKey    APIKey
	APIKeyPtr *APIKey `sqlxx:"fk:APIKeyID"`

	AvatarID sql.NullInt64 `db:"avatar_id"`
	Avatar   *Media

	Avatars []Avatar
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

type AvatarFilter struct {
	ID   int    `db:"id" sqlxx:"primary_key:true"`
	Name string `db:"name"`
}

func (AvatarFilter) TableName() string { return "avatar_filters" }

type Avatar struct {
	ID        int       `db:"id" sqlxx:"primary_key:true"`
	Path      string    `db:"path"`
	UserID    int       `db:"user_id"`
	FilterID  int       `db:"filter_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Filter    AvatarFilter
	FilterPtr *AvatarFilter `sqlxx:"fk:FilterID"`
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
// Loader
// ----------------------------------------------------------------------------

type environment struct {
	driver             *sqlxx.Client
	t                  *testing.T
	Users              []User
	APIKeys            []APIKey
	Profiles           []Profile
	AvatarFilters      []AvatarFilter
	Avatars            []Avatar
	Articles           []Article
	Categories         []Category
	Tags               []Tag
	ArticlesCategories []ArticleCategory
	Partners           []Partner
	Managers           []Manager
	Projects           []Project
	Medias             []Media
}

func (e *environment) load() {
	e.insertPartners()
	e.insertAPIKeys()
	e.insertMedias()
	e.insertUsers()
	e.insertManagers()
	e.insertProjects()
	e.insertAvatarFilters()
	e.insertAvatars()
	e.insertProfiles()
	e.insertCategories()
	e.insertTags()
	e.insertArticles()
}

func (e *environment) insertPartners() {
	e.driver.MustExec("INSERT INTO partners (name) VALUES ($1)", "Ulule")
	assert.NoError(e.t, e.driver.Select(&e.Partners, "SELECT * FROM partners"))
}

func (e *environment) insertAPIKeys() {
	e.driver.MustExec(
		"INSERT INTO api_keys (key, partner_id) VALUES ($1, $2)",
		"this-is-my-scret-api-key",
		e.Partners[0].ID)

	assert.NoError(e.t, e.driver.Select(&e.APIKeys, "SELECT * FROM api_keys"))
}

func (e *environment) insertMedias() {
	e.driver.MustExec(
		"INSERT INTO media (path) VALUES ($1)",
		"media/avatar.png")

	assert.NoError(e.t, e.driver.Select(&e.Medias, "SELECT * FROM media"))
}
func (e *environment) insertUsers() {
	e.driver.MustExec(
		"INSERT INTO users (username, api_key_id, avatar_id) VALUES ($1, $2, $3)",
		"jdoe",
		e.APIKeys[0].ID,
		e.Medias[0].ID)

	assert.NoError(e.t, e.driver.Select(&e.Users, "SELECT * FROM users WHERE username=$1", "jdoe"))
}

func (e *environment) insertManagers() {
	e.driver.MustExec(
		"INSERT INTO managers (name, user_id) VALUES ($1, $2)",
		"Super Owl",
		e.Users[0].ID)

	assert.NoError(e.t, e.driver.Select(&e.Managers, "SELECT * FROM managers"))
}

func (e *environment) insertProjects() {
	e.driver.MustExec(
		"INSERT INTO projects (name, manager_id, user_id) VALUES ($1, $2, $3)",
		"Super Project",
		e.Managers[0].ID,
		e.Users[0].ID)

	assert.NoError(e.t, e.driver.Select(&e.Projects, "SELECT * FROM projects"))
}

func (e *environment) insertAvatarFilters() {
	names := []string{"normal", "clarendon", "juno", "lark", "ludwig", "gingham", "valencia"}
	for _, name := range names {
		e.driver.MustExec("INSERT INTO avatar_filters (name) VALUES ($1)", name)
	}

	assert.NoError(e.t, e.driver.Select(&e.AvatarFilters, "SELECT * FROM avatar_filters"))
}

func (e *environment) insertAvatars() {
	for i := 0; i < 5; i++ {
		e.driver.MustExec(
			"INSERT INTO avatars (path, user_id, filter_id) VALUES ($1, $2, $3)",
			fmt.Sprintf("/avatars/%s-%d.png", e.Users[0].Username, i),
			e.Users[0].ID,
			e.AvatarFilters[0].ID)
	}

	assert.NoError(e.t, e.driver.Select(&e.Avatars, "SELECT * FROM avatars"))
}

func (e *environment) insertProfiles() {
	e.driver.MustExec(
		"INSERT INTO profiles (user_id, first_name, last_name) VALUES ($1, $2, $3)",
		e.Users[0].ID,
		"John",
		"Doe")

	assert.NoError(e.t, e.driver.Select(&e.Profiles, "SELECT * FROM profiles"))
}

func (e *environment) insertCategories() {
	for i := 0; i < 5; i++ {
		e.driver.MustExec(
			"INSERT INTO categories (name, user_id) VALUES ($1, $2)",
			fmt.Sprintf("Category #%d", i),
			e.Users[0].ID)
	}

	assert.NoError(e.t, e.driver.Select(&e.Categories, "SELECT * FROM categories"))
}

func (e *environment) insertTags() {
	e.driver.MustExec("INSERT INTO tags (name) VALUES ($1)", "Tag")
	assert.NoError(e.t, e.driver.Select(&e.Tags, "SELECT * FROM tags"))
}

func (e *environment) insertArticles() {
	for i := 0; i < 5; i++ {
		e.driver.MustExec(
			"INSERT INTO articles (title, author_id, reviewer_id, main_tag_id) VALUES ($1, $2, $3, $4)",
			fmt.Sprintf("Title #%d", i),
			e.Users[0].ID,
			e.Users[0].ID,
			e.Tags[0].ID)
	}

	assert.NoError(e.t, e.driver.Select(&e.Articles, "SELECT * FROM articles"))

	for _, article := range e.Articles {
		for _, category := range e.Categories {
			e.driver.MustExec(
				"INSERT INTO articles_categories (article_id, category_id) VALUES ($1, $2)",
				article.ID,
				category.ID)
		}
	}

	assert.NoError(e.t, e.driver.Select(&e.ArticlesCategories, "SELECT * FROM articles_categories"))
}

func (e *environment) createComment(user *User, article *Article) Comment {
	var id int

	err := e.driver.QueryRowx(
		"INSERT INTO comments (content, user_id, article_id) VALUES ($1, $2, $3) RETURNING id",
		"Lorem Ipsum",
		user.ID,
		article.ID).Scan(&id)

	assert.Nil(e.t, err)

	comment := Comment{}
	assert.NoError(e.t, e.driver.Get(&comment, "SELECT * FROM comments WHERE id = $1", id))

	return comment
}

func (e *environment) createArticle(user *User) Article {
	var id int

	err := e.driver.QueryRowx(
		"INSERT INTO articles (title, author_id, reviewer_id) VALUES ($1, $2, $3) RETURNING id",
		"Title",
		user.ID,
		user.ID).Scan(&id)

	assert.Nil(e.t, err)

	article := Article{}
	assert.NoError(e.t, e.driver.Get(&article, "SELECT * FROM articles WHERE id = $1", id))

	return article
}

func (e *environment) createUser(username string) User {
	var (
		key  = fmt.Sprintf("%s-apikey", username)
		name = fmt.Sprintf("%s-partner", username)
	)

	partner := Partner{}

	e.driver.MustExec(
		"INSERT INTO partners (name) VALUES ($1)",
		name)

	assert.NoError(e.t, e.driver.Get(&partner, "SELECT * FROM partners WHERE name = $1", name))

	media := Media{}

	e.driver.MustExec(
		"INSERT INTO media (path) VALUES ($1)",
		fmt.Sprintf("media/media-%s.png", username))

	assert.NoError(e.t, e.driver.Get(&media, "SELECT * FROM media ORDER BY id DESC LIMIT 1"))

	apiKey := APIKey{}

	e.driver.MustExec(
		"INSERT INTO api_keys (key, partner_id) VALUES ($1, $2)",
		key,
		partner.ID)

	assert.NoError(e.t, e.driver.Get(&apiKey, "SELECT * FROM api_keys WHERE key = $1", key))

	user := User{}

	e.driver.MustExec(
		"INSERT INTO users (username, api_key_id, avatar_id) VALUES ($1, $2, $3)",
		username,
		apiKey.ID,
		media.ID)

	assert.NoError(e.t, e.driver.Get(&user, "SELECT * FROM users WHERE username=$1", username))

	for i := 1; i < 6; i++ {
		e.driver.MustExec(
			"INSERT INTO avatars (path, user_id, filter_id) VALUES ($1, $2, $3)",
			fmt.Sprintf("/avatars/%s-%d.png", username, i),
			user.ID,
			i)
	}

	return user
}

func (e *environment) createCategory(name string, userID *int) Category {
	e.driver.MustExec("INSERT INTO categories (name) VALUES ($1)", name)
	if userID != nil {
		e.driver.MustExec("UPDATE categories SET user_id=$1 WHERE name=$2", *userID, name)
	}

	category := Category{}
	assert.NoError(e.t, e.driver.Get(&category, "SELECT * FROM categories WHERE name=$1", name))

	return category
}

func (e *environment) teardown() {
	value := os.Getenv("KEEP_DB")
	if len(value) == 0 {
		e.driver.MustExec(dropTables)
	}

	assert.NoError(e.t, e.driver.Close())
}

func dbParamString(option func(string) sqlxx.Option, param string, env ...string) sqlxx.Option {
	param = strings.ToUpper(param)
	v := os.Getenv(fmt.Sprintf("DB_%s", param))
	if len(v) != 0 {
		return option(v)
	}
	for i := range env {
		v = os.Getenv(env[i])
		if len(v) != 0 {
			return option(v)
		}
	}
	return dbDefaultOptions[param]
}

func dbParamInt(option func(int) sqlxx.Option, param string, env ...string) sqlxx.Option {
	param = strings.ToUpper(param)
	v := os.Getenv(fmt.Sprintf("DB_%s", param))
	n, err := strconv.Atoi(v)
	if err == nil {
		return option(n)
	}
	for i := range env {
		v = os.Getenv(env[i])
		n, err = strconv.Atoi(v)
		if err == nil {
			return option(n)
		}
	}
	return dbDefaultOptions[param]
}

func setup(t *testing.T) *environment {

	db, err := sqlxx.New(
		dbParamString(sqlxx.Host, "host", "PGHOST"),
		dbParamInt(sqlxx.Port, "port", "PGPORT"),
		dbParamString(sqlxx.User, "user", "PGUSER"),
		dbParamString(sqlxx.Password, "password", "PGPASSWORD"),
		dbParamString(sqlxx.Database, "name", "PGDATABASE"),
	)
	assert.NoError(t, err)
	assert.NotNil(t, db)

	db.MustExec(dropTables)
	db.MustExec(dbSchema)

	env := &environment{t: t, driver: db}
	env.load()

	return env
}
