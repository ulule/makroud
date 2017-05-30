package sqlxx_test

import (
	_ "database/sql"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

type ProfileV2 struct {
	ID        int64
	UserID    int
	FirstName string
	LastName  string
}

func (ProfileV2) CreateSchema(builder sqlxx.SchemaBuilder) {
	builder.SetTableName("Profile", "profiles").
		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
		AddField("FirstName", "first_name").
		AddField("LastName", "last_name").
		AddField("UserID", "user_id", sqlxx.IsForeignKey("User"))
}

type UserV2 struct {
	ID        int64
	Username  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	Profile   ProfileV2
	Avatars   []AvatarV2
	KeyID     int64
	Key       UserKeyV2
}

func (UserV2) CreateSchema(builder sqlxx.SchemaBuilder) {
	builder.SetTableName("User", "users").
		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
		AddField("Username", "username").
		AddField("IsActive", "is_active", sqlxx.HasDefault("true")).
		AddField("CreatedAt", "created_at", sqlxx.HasDefault("NOW()")).
		AddField("UpdatedAt", "updated_at", sqlxx.HasDefault("NOW()")).
		AddField("DeletedAt", "deleted_at", sqlxx.IsArchiveKey()).
		AddField("KeyID", "key_id", sqlxx.IsForeignKey("UserKey")).
		AddAssociation("Avatars", "Avatar", sqlxx.AssociationTypeMany).
		AddAssociation("Profile", "Profile", sqlxx.AssociationTypeOne).
		AddAssociation("Key", "UserKey", sqlxx.AssociationTypeOne)
}

type UserKeyV2 struct {
	ID  int64
	Key string
}

func (UserKeyV2) CreateSchema(builder sqlxx.SchemaBuilder) {
	builder.SetTableName("UserKey", "keys").
		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
		AddField("Key", "payload")
}

type CommentV2 struct {
	ID        int64
	Content   string
	UserID    int64
	User      *UserV2
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (CommentV2) CreateSchema(builder sqlxx.SchemaBuilder) {
	builder.SetTableName("Comment", "comments").
		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
		AddField("Content", "content").
		AddField("CreatedAt", "created_at", sqlxx.HasDefault("NOW()")).
		AddField("UpdatedAt", "updated_at", sqlxx.HasDefault("NOW()")).
		AddField("UserID", "user_id", sqlxx.IsForeignKey("User")).
		AddAssociation("User", "User", sqlxx.AssociationTypeOne)
}

type AvatarV2 struct {
	ID        int64
	Path      string
	UserID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (AvatarV2) CreateSchema(builder sqlxx.SchemaBuilder) {
	builder.SetTableName("Avatar", "avatars").
		SetPrimaryKey("ID", "id", sqlxx.PrimaryKeyInteger).
		AddField("Path", "path").
		AddField("CreatedAt", "created_at", sqlxx.HasDefault("NOW()")).
		AddField("UpdatedAt", "updated_at", sqlxx.HasDefault("NOW()")).
		AddField("UserID", "user_id", sqlxx.IsForeignKey("User"))
}

func TestSchema_v2(t *testing.T) {
	env := setup(t, sqlxx.Cache(true))
	defer env.teardown()

	is := require.New(t)

	user := &UserV2{}
	_ = user
	// comment := &CommentV2{}
	// builder := sqlxx.NewSchemaBuilder()
	// comment.CreateSchema(builder)
	// schema, err := builder.Create(env.driver, comment)
	//
	// is.NoError(err)
	// is.NotNil(schema)
	// spew.Dump(schema)

	s1, err := sqlxx.XGetSchema(env.driver, user)
	is.NoError(err)
	spew.Dump(s1)
}
