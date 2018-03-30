package sqlxx_test

import (
	_ "database/sql"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

type CommentV2 struct {
	ID        int64
	Content   string
	UserID    int64
	User      *User
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

func (comment *CommentV2) WriteModel(mapper sqlxx.Mapper) error {
	return nil
}

func TestSchema_v2(t *testing.T) {
	env := setup(t, sqlxx.Cache(true))
	defer env.teardown()

	is := require.New(t)

	user := &User{}
	_ = user
	comment := &CommentV2{}
	builder := sqlxx.NewSchemaBuilder()
	comment.CreateSchema(builder)
	schema, err := builder.Create(env.driver, comment)

	is.NoError(err)
	is.NotNil(schema)
	spew.Dump(schema)

	s1, err := sqlxx.XGetSchema(env.driver, user)
	is.NoError(err)
	spew.Dump(s1)

	sqlxx.XGetSchema(env.driver, &Avatar{})

}
