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
	DROP TABLE IF EXISTS avatars CASCADE;
	DROP TABLE IF EXISTS articles CASCADE;
`

var dbSchema = `CREATE TABLE users (
	id 				serial primary key not null,
	username 	    varchar(30) not null,
	is_active 		boolean default true,
    created_at 		timestamp default current_timestamp,
    updated_at 		timestamp default current_timestamp
);

CREATE TABLE avatars (
	id 				serial primary key not null,
	path 			varchar(255) not null,
	user_id 		integer references users(id),
    created_at 		timestamp default current_timestamp,
    updated_at 		timestamp default current_timestamp
);

CREATE TABLE articles (
	id 				serial primary key not null,
	title 			varchar(255) not null,
	author_id 		integer references users(id),
	is_published 	boolean default true,
    created_at 		timestamp default current_timestamp,
    updated_at 		timestamp default current_timestamp
);`

type TestData struct {
	User     User
	Avatars  []Avatar
	Articles []Article
}

type User struct {
	ID        int       `db:"id" sqlxx:"primary_key:true ignored:true"`
	Username  string    `db:"username"`
	IsActive  bool      `db:"is_active" sqlxx:"default:true"`
	CreatedAt time.Time `db:"created_at" sqlxx:"auto_now_add:true"`
	UpdatedAt time.Time `db:"updated_at" sqlxx:"default:now()"`

	// Avatars []Avatar
}

func (User) TableName() string {
	return "users"
}

type Avatar struct {
	ID        int       `db:"id" sqlxx:"primary_key:true"`
	Path      string    `db:"path"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (Avatar) TableName() string {
	return "avatars"
}

type Article struct {
	ID          int       `db:"id" sqlxx:"primary_key:true"`
	Title       string    `db:"title"`
	AuthorID    int       `db:"author_id"`
	IsPublished bool      `db:"is_published"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`

	Author User `sqlxx:"related:author_id"`
}

func (Article) TableName() string {
	return "articles"
}

func dbParam(param string) string {
	param = strings.ToUpper(param)

	if v := os.Getenv(fmt.Sprintf("DB_%s", param)); len(v) != 0 {
		return v
	}

	return dbDefaultParams[param]
}

func loadData(t *testing.T, driver Driver) *TestData {
	driver.MustExec("INSERT INTO users (username) VALUES ($1)", "jdoe")

	user := User{}
	require.NoError(t, driver.Get(&user, "SELECT * FROM users WHERE username=$1", "jdoe"))

	for i := 0; i < 5; i++ {
		driver.MustExec("INSERT INTO avatars (path, user_id) VALUES ($1, $2)", fmt.Sprintf("/avatars/%s-%d.png", user.Username, i), user.ID)

	}

	avatars := []Avatar{}
	require.NoError(t, driver.Select(&avatars, "SELECT * FROM avatars WHERE user_id=$1", user.ID))

	for i := 0; i < 5; i++ {
		driver.MustExec("INSERT INTO articles (title, author_id) VALUES ($1, $2)", fmt.Sprintf("Title #%d", i), user.ID)
	}

	articles := []Article{}
	require.NoError(t, driver.Select(&articles, "SELECT * FROM articles WHERE author_id=$1", user.ID))

	return &TestData{
		User:     user,
		Avatars:  avatars,
		Articles: articles,
	}
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
		dbx.MustExec(dropTables)
		require.NoError(t, db.Close())
	}
}
