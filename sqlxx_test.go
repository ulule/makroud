package sqlxx_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

var dbDefaultOptions = map[string]sqlxx.Option{
	"USER":     sqlxx.User("postgres"),
	"PASSWORD": sqlxx.Password(""),
	"HOST":     sqlxx.Host("localhost"),
	"PORT":     sqlxx.Port(5432),
	"NAME":     sqlxx.Database("sqlxx_test"),
}

// ----------------------------------------------------------------------------
// Miscellaneous models
// ----------------------------------------------------------------------------

type Elements struct {
	Air     string `db:"air"`
	Fire    string `sqlxx:"column:fire"`
	Water   string `sqlxx:"-"`
	Earth   string `sqlxx:"column:earth,default"`
	Fifth   string
	enabled bool
}

func (Elements) TableName() string {
	return "rune_elements"
}

// ----------------------------------------------------------------------------
// Object storage application
// ----------------------------------------------------------------------------

// type ExoFile struct {
// 	// Columns
// 	ID   string `sqlxx:"column:id,pk:ulid"`
// 	Name string `sqlxx:"column:name"`
// 	Path string `sqlxx:"column:path"`
// 	// Relationships
// 	Chunk []ExoChunk
// }
//
// func (ExoFile) TableName() string {
// 	return "exo_file"
// }
//
// type ExoFileChunk struct {
// 	// Columns
// 	ID      string `sqlxx:"column:id,pk:ulid"`
// 	FileID  string `sqlxx:"column:file_id,fk:exo_file"`
// 	ChunkID string `sqlxx:"column:chunk_id,fk:exo_chunk"`
// }
//
// func (ExoFileChunk) TableName() string {
// 	return "exo_file_chunk"
// }

type ExoChunk struct {
	// Columns
	Hash   string `sqlxx:"column:hash,pk:ulid"`
	Bytes  string `sqlxx:"column:bytes"`
	ModeID int64  `sqlxx:"column:mode_id,fk:exo_chunk_mode"`
	// Relationships
	Signature *ExoChunkSignature
	Mode      *ExoChunkMode
}

func (ExoChunk) TableName() string {
	return "exo_chunk"
}

type ExoChunkSignature struct {
	ID      string `sqlxx:"column:id,pk:ulid"`
	ChunkID string `sqlxx:"column:chunk_id,fk:exo_chunk"`
	Bytes   string `sqlxx:"column:bytes"`
}

func (ExoChunkSignature) TableName() string {
	return "exo_chunk_signature"
}

type ExoChunkMode struct {
	ID   int64  `sqlxx:"column:id,pk"`
	Mode string `sqlxx:"column:mode"`
}

func (ExoChunkMode) TableName() string {
	return "exo_chunk_mode"
}

// ----------------------------------------------------------------------------
// Zootopia
// ----------------------------------------------------------------------------

type Group struct {
	// Columns
	ID   int64  `sqlxx:"column:id,pk"`
	Name string `sqlxx:"column:name"`
}

func (Group) TableName() string {
	return "ztp_group"
}

type Center struct {
	// Columns
	ID   string `sqlxx:"column:id"`
	Name string `sqlxx:"column:name"`
	Area string `sqlxx:"column:area"`
}

func (Center) TableName() string {
	return "ztp_center"
}

type Owl struct {
	// Columns
	ID           int64         `sqlxx:"column:id,pk"`
	Name         string        `sqlxx:"column:name"`
	FeatherColor string        `sqlxx:"column:feather_color"`
	FavoriteFood string        `sqlxx:"column:favorite_food"`
	GroupID      sql.NullInt64 `sqlxx:"column:group_id,fk:ztp_group"`
	// Relationships
	Group    *Group
	Packages []Package
}

func (Owl) TableName() string {
	return "ztp_owl"
}

type Package struct {
	// Columns
	ID            string        `sqlxx:"column:id"`
	Status        string        `sqlxx:"column:status"`
	SenderID      string        `sqlxx:"column:sender_id,fk:ztp_center"`
	ReceiverID    string        `sqlxx:"column:receiver_id,fk:ztp_center"`
	TransporterID sql.NullInt64 `sqlxx:"column:transporter_id,fk:ztp_owl"`
	// Relationships
	Sender   *Center
	Receiver *Center
}

func (Package) TableName() string {
	return "ztp_package"
}

type Cat struct {
	// Columns
	ID        string      `sqlxx:"column:id,pk:ulid"`
	Name      string      `sqlxx:"column:name"`
	CreatedAt time.Time   `sqlxx:"column:created_at,default"`
	UpdatedAt time.Time   `sqlxx:"column:updated_at,default"`
	DeletedAt pq.NullTime `sqlxx:"column:deleted_at"`
	// Relationships
	Owner *Human
	Meows []*Meow
}

func (Cat) TableName() string {
	return "ztp_cat"
}

type Meow struct {
	// Columns
	Hash      string      `sqlxx:"column:hash,pk:ulid"`
	Body      string      `sqlxx:"column:body"`
	CatID     string      `sqlxx:"column:cat_id,fk:ztp_cat"`
	CreatedAt time.Time   `sqlxx:"column:created"`
	UpdatedAt time.Time   `sqlxx:"column:updated"`
	DeletedAt pq.NullTime `sqlxx:"column:deleted"`
}

func (Meow) TableName() string {
	return "ztp_meow"
}

func (Meow) CreatedKey() string {
	return "created"
}

func (Meow) UpdatedKey() string {
	return "updated"
}

func (Meow) DeletedKey() string {
	return "deleted"
}

type Human struct {
	// Columns
	ID        string         `sqlxx:"column:id,pk:ulid"`
	Name      string         `sqlxx:"column:name"`
	CreatedAt time.Time      `sqlxx:"column:created_at,default"`
	UpdatedAt time.Time      `sqlxx:"column:updated_at,default"`
	DeletedAt pq.NullTime    `sqlxx:"column:deleted_at"`
	CatID     sql.NullString `sqlxx:"column:cat_id,fk:ztp_cat"`
}

func (Human) TableName() string {
	return "ztp_human"
}

// ----------------------------------------------------------------------------
// Loader
// ----------------------------------------------------------------------------

type environment struct {
	driver *sqlxx.Client
	is     *require.Assertions
}

func (e *environment) startup(ctx context.Context) {
	DropTables(ctx, e.driver)
	CreateTables(ctx, e.driver)
}

func (e *environment) shutdown(ctx context.Context) {
	value := os.Getenv("DB_KEEP")
	if len(value) == 0 {
		DropTables(ctx, e.driver)
	}
	e.is.NoError(e.driver.Close())
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

type SetupCallback func(handler SetupHandler)

type SetupHandler func(driver sqlxx.Driver)

func Setup(t require.TestingT, options ...sqlxx.Option) SetupCallback {
	is := require.New(t)
	ctx := context.Background()
	opts := []sqlxx.Option{
		dbParamString(sqlxx.Host, "host", "PGHOST"),
		dbParamInt(sqlxx.Port, "port", "PGPORT"),
		dbParamString(sqlxx.User, "user", "PGUSER"),
		dbParamString(sqlxx.Password, "password", "PGPASSWORD"),
		dbParamString(sqlxx.Database, "name", "PGDATABASE"),
		sqlxx.Cache(true),
	}
	opts = append(opts, options...)

	db, err := sqlxx.New(opts...)
	is.NoError(err)
	is.NotNil(db)

	env := &environment{
		is:     is,
		driver: db,
	}

	return func(handler SetupHandler) {
		env.startup(ctx)
		handler(db)
		env.shutdown(ctx)
	}
}

func DropTables(ctx context.Context, db *sqlxx.Client) {
	db.MustExec(ctx, `
		-- Simple schema
		DROP TABLE IF EXISTS ztp_human CASCADE;
		DROP TABLE IF EXISTS ztp_package CASCADE;
		DROP TABLE IF EXISTS ztp_owl CASCADE;
		DROP TABLE IF EXISTS ztp_cat CASCADE;
		DROP TABLE IF EXISTS ztp_meow CASCADE;
		DROP TABLE IF EXISTS ztp_group CASCADE;
		DROP TABLE IF EXISTS ztp_center CASCADE;

		-- Object storage application
		DROP TABLE IF EXISTS exo_chunk_signature CASCADE;
		DROP TABLE IF EXISTS exo_chunk CASCADE;
		DROP TABLE IF EXISTS exo_chunk_mode CASCADE;

	`)
}

func CreateTables(ctx context.Context, db *sqlxx.Client) {
	db.MustExec(ctx, `

		--
		-- Zootopia schema
		--

		CREATE TABLE ztp_group (
			id                SERIAL PRIMARY KEY NOT NULL,
			name              VARCHAR(255) NOT NULL
		);
		CREATE TABLE ztp_center (
			id                VARCHAR(32) PRIMARY KEY NOT NULL DEFAULT md5(random()::text),
			name              VARCHAR(255) NOT NULL,
			area              VARCHAR(255) NOT NULL
		);
		CREATE TABLE ztp_owl (
			id                SERIAL PRIMARY KEY NOT NULL,
			name              VARCHAR(255) NOT NULL,
			feather_color     VARCHAR(255) NOT NULL,
			favorite_food     VARCHAR(255) NOT NULL,
			group_id          INTEGER REFERENCES ztp_group(id)
		);
		CREATE TABLE ztp_cat (
			id                VARCHAR(26) PRIMARY KEY NOT NULL,
			name              VARCHAR(255) NOT NULL,
			created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at        TIMESTAMP WITH TIME ZONE
		);
		CREATE TABLE ztp_meow (
			hash              VARCHAR(26) PRIMARY KEY NOT NULL,
			body              VARCHAR(2048) NOT NULL,
			cat_id            VARCHAR(26) NOT NULL REFERENCES ztp_cat(id),
			created           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted           TIMESTAMP WITH TIME ZONE
		);
		CREATE TABLE ztp_human (
			id                VARCHAR(26) PRIMARY KEY NOT NULL,
			name              VARCHAR(255) NOT NULL,
			cat_id            VARCHAR(26) REFERENCES ztp_cat(id),
			created_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at        TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			deleted_at        TIMESTAMP WITH TIME ZONE
		);
		CREATE TABLE ztp_package (
			id                VARCHAR(32) PRIMARY KEY NOT NULL DEFAULT md5(random()::text),
			status            VARCHAR(255) NOT NULL,
			sender_id         VARCHAR(32) NOT NULL REFERENCES ztp_center(id),
			receiver_id       VARCHAR(32) NOT NULL REFERENCES ztp_center(id),
			transporter_id    INTEGER REFERENCES ztp_owl(id)
		);

		--
		-- Object storage application
		--

		CREATE TABLE exo_chunk_mode (
			id              SERIAL PRIMARY KEY NOT NULL,
			mode            VARCHAR(255) NOT NULL
		);
		CREATE TABLE exo_chunk (
			hash            VARCHAR(26) PRIMARY KEY NOT NULL,
			bytes           VARCHAR(2048) NOT NULL,
			mode_id         INTEGER NOT NULL REFERENCES exo_chunk_mode(id) ON DELETE RESTRICT
		);
		CREATE TABLE exo_chunk_signature (
			id               VARCHAR(26) PRIMARY KEY NOT NULL,
			chunk_id         VARCHAR(26) NOT NULL REFERENCES exo_chunk(hash),
			bytes            VARCHAR(2048) NOT NULL
		);

		--
		-- Application schema
		--

		-- TODO

	`)
}
