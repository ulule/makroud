package sqlxx_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

// TODO: Finish and improve coverage

func TestSchema_Owl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		model := &Owl{}

		schema, err := sqlxx.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		is.Equal("ztp_owl", schema.TableName())
		is.Equal("id", schema.PrimaryKey().ColumnName())
		is.Equal("ztp_owl.id", schema.PrimaryKey().ColumnPath())
		is.False(schema.HasCreatedKey())
		is.False(schema.HasUpdatedKey())
		is.False(schema.HasDeletedKey())

		columns := schema.Columns()
		is.Len(columns, 5)
		is.Contains(columns, "id")
		is.Contains(columns, "name")
		is.Contains(columns, "feather_color")
		is.Contains(columns, "favorite_food")
		is.Contains(columns, "group_id")

		columns = schema.ColumnPaths()
		is.Len(columns, 5)
		is.Contains(columns, "ztp_owl.id")
		is.Contains(columns, "ztp_owl.name")
		is.Contains(columns, "ztp_owl.feather_color")
		is.Contains(columns, "ztp_owl.favorite_food")
		is.Contains(columns, "ztp_owl.group_id")

	})
}

func TestSchema_Cat(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		model := &Cat{}

		schema, err := sqlxx.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		is.Equal("ztp_cat", schema.TableName())
		is.Equal("id", schema.PrimaryKey().ColumnName())
		is.Equal("ztp_cat.id", schema.PrimaryKey().ColumnPath())

		is.True(schema.HasCreatedKey())
		is.True(schema.HasUpdatedKey())
		is.True(schema.HasDeletedKey())
		is.Equal("ztp_cat.created_at", schema.CreatedKeyPath())
		is.Equal("ztp_cat.updated_at", schema.UpdatedKeyPath())
		is.Equal("ztp_cat.deleted_at", schema.DeletedKeyPath())

		columns := schema.Columns()
		is.Len(columns, 5)
		is.Contains(columns, "id")
		is.Contains(columns, "name")
		is.Contains(columns, "created_at")
		is.Contains(columns, "updated_at")
		is.Contains(columns, "deleted_at")

		columns = schema.ColumnPaths()
		is.Len(columns, 5)
		is.Contains(columns, "ztp_cat.id")
		is.Contains(columns, "ztp_cat.name")
		is.Contains(columns, "ztp_cat.created_at")
		is.Contains(columns, "ztp_cat.updated_at")
		is.Contains(columns, "ztp_cat.deleted_at")

	})
}

func TestSchema_Meow(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		model := &Meow{}

		schema, err := sqlxx.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		is.Equal("ztp_meow", schema.TableName())
		is.Equal("hash", schema.PrimaryKey().ColumnName())
		is.Equal("ztp_meow.hash", schema.PrimaryKey().ColumnPath())

		is.True(schema.HasCreatedKey())
		is.True(schema.HasUpdatedKey())
		is.True(schema.HasDeletedKey())
		is.Equal("ztp_meow.created", schema.CreatedKeyPath())
		is.Equal("ztp_meow.updated", schema.UpdatedKeyPath())
		is.Equal("ztp_meow.deleted", schema.DeletedKeyPath())

		columns := schema.Columns()
		is.Len(columns, 5)
		is.Contains(columns, "hash")
		is.Contains(columns, "body")
		is.Contains(columns, "created")
		is.Contains(columns, "updated")
		is.Contains(columns, "deleted")

		columns = schema.ColumnPaths()
		is.Len(columns, 5)
		is.Contains(columns, "ztp_meow.hash")
		is.Contains(columns, "ztp_meow.body")
		is.Contains(columns, "ztp_meow.created")
		is.Contains(columns, "ztp_meow.updated")
		is.Contains(columns, "ztp_meow.deleted")

	})
}

func TestSchema_ExoChunk(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		model := &ExoChunk{}

		schema, err := sqlxx.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		// TODO: Finish

		// is.Equal("ztp_meow", schema.TableName())
		// is.Equal("hash", schema.PrimaryKey().ColumnName())
		// is.Equal("ztp_meow.hash", schema.PrimaryKey().ColumnPath())
		//
		// is.True(schema.HasCreatedKey())
		// is.True(schema.HasUpdatedKey())
		// is.True(schema.HasDeletedKey())
		// is.Equal("ztp_meow.created", schema.CreatedKeyPath())
		// is.Equal("ztp_meow.updated", schema.UpdatedKeyPath())
		// is.Equal("ztp_meow.deleted", schema.DeletedKeyPath())
		//
		// columns := schema.Columns()
		// is.Len(columns, 5)
		// is.Contains(columns, "hash")
		// is.Contains(columns, "body")
		// is.Contains(columns, "created")
		// is.Contains(columns, "updated")
		// is.Contains(columns, "deleted")
		//
		// columns = schema.ColumnPaths()
		// is.Len(columns, 5)
		// is.Contains(columns, "ztp_meow.hash")
		// is.Contains(columns, "ztp_meow.body")
		// is.Contains(columns, "ztp_meow.created")
		// is.Contains(columns, "ztp_meow.updated")
		// is.Contains(columns, "ztp_meow.deleted")

	})
}
