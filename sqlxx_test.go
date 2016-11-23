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
	"USER":     "root",
	"PASSWORS": "",
	"HOST":     "localhost",
	"PORT":     "5432",
	"NAME":     "sqlxx_test",
}

var dbSchema = `
CREATE TABLE avatar (
	id serial primary key not null,
	path varchar(255) not null,
	user_id integer references user(id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

CREATE TABLE author (
	id serial primary key not null, 
	name varchar(30) not null,
	birthday timestamp with time zone,
	is_active boolean default true,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);

CREATE TABLE article (
	id serial primary key not null, 
	title varchar(255) not null,
	body text,
	author_id integer references user(id),
	is_published boolean default true,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp
);`

type Avatar struct {
	ID     int
	Path   string
	UserID int
}

type User struct {
	ID       int
	Name     string
	Birthday time.Time
	IsActive bool
	Avatars  []Avatar
}

type Article struct {
	ID          int
	Title       string
	AuthorID    int
	Author      User `sqlxx:"author_id"`
	IsPublished bool
}

func dbParam(param string) string {
	param = strings.ToUpper(param)

	if v := os.Getenv(fmt.Sprintf("DB_%s", param)); len(v) != 0 {
		return v
	}

	return dbDefaultParams[param]
}

func dbConnection(t *testing.T) (*sqlx.DB, func()) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=false;timezone=UTC",
		dbParam("user"),
		dbParam("password"),
		dbParam("host"),
		dbParam("port"),
		dbParam("name")))

	require.NoError(t, err)

	return sqlx.NewDb(db, "postgres"), func() {
		require.NoError(t, db.Close())
	}
}
