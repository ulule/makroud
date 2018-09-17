package makroud_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/makroud"
)

func TestSchema_Owl(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		model := &Owl{}

		schema, err := makroud.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		is.IsType(*model, schema.Model())
		is.Equal("Owl", schema.ModelName())
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
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		model := &Cat{}

		schema, err := makroud.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		is.IsType(*model, schema.Model())
		is.Equal("Cat", schema.ModelName())
		is.Equal("ztp_cat", schema.TableName())
		is.Equal("id", schema.PrimaryKey().ColumnName())
		is.Equal("ztp_cat.id", schema.PrimaryKey().ColumnPath())

		is.True(schema.HasCreatedKey())
		is.True(schema.HasUpdatedKey())
		is.True(schema.HasDeletedKey())
		is.Equal("created_at", schema.CreatedKeyName())
		is.Equal("updated_at", schema.UpdatedKeyName())
		is.Equal("deleted_at", schema.DeletedKeyName())
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
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		model := &Meow{}

		schema, err := makroud.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		is.IsType(*model, schema.Model())
		is.Equal("Meow", schema.ModelName())
		is.Equal("ztp_meow", schema.TableName())
		is.Equal("hash", schema.PrimaryKey().ColumnName())
		is.Equal("ztp_meow.hash", schema.PrimaryKey().ColumnPath())

		is.True(schema.HasCreatedKey())
		is.True(schema.HasUpdatedKey())
		is.True(schema.HasDeletedKey())
		is.Equal("created", schema.CreatedKeyName())
		is.Equal("updated", schema.UpdatedKeyName())
		is.Equal("deleted", schema.DeletedKeyName())
		is.Equal("ztp_meow.created", schema.CreatedKeyPath())
		is.Equal("ztp_meow.updated", schema.UpdatedKeyPath())
		is.Equal("ztp_meow.deleted", schema.DeletedKeyPath())

		columns := schema.Columns()
		is.Len(columns, 6)
		is.Contains(columns, "hash")
		is.Contains(columns, "body")
		is.Contains(columns, "cat_id")
		is.Contains(columns, "created")
		is.Contains(columns, "updated")
		is.Contains(columns, "deleted")

		columns = schema.ColumnPaths()
		is.Len(columns, 6)
		is.Contains(columns, "ztp_meow.hash")
		is.Contains(columns, "ztp_meow.body")
		is.Contains(columns, "ztp_meow.cat_id")
		is.Contains(columns, "ztp_meow.created")
		is.Contains(columns, "ztp_meow.updated")
		is.Contains(columns, "ztp_meow.deleted")

	})
}

func TestSchema_ExoChunk(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		model := &ExoChunk{}

		schema, err := makroud.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		is.IsType(*model, schema.Model())
		is.Equal("ExoChunk", schema.ModelName())
		is.Equal("exo_chunk", schema.TableName())
		is.Equal("hash", schema.PrimaryKey().ColumnName())
		is.Equal("exo_chunk.hash", schema.PrimaryKey().ColumnPath())

		is.False(schema.HasCreatedKey())
		is.False(schema.HasUpdatedKey())
		is.False(schema.HasDeletedKey())

		columns := schema.Columns()
		is.Len(columns, 6)
		is.Contains(columns, "hash")
		is.Contains(columns, "bytes")
		is.Contains(columns, "organization_id")
		is.Contains(columns, "user_id")
		is.Contains(columns, "mode_id")
		is.Contains(columns, "file_id")

		columns = schema.ColumnPaths()
		is.Len(columns, 6)
		is.Contains(columns, "exo_chunk.hash")
		is.Contains(columns, "exo_chunk.bytes")
		is.Contains(columns, "exo_chunk.organization_id")
		is.Contains(columns, "exo_chunk.user_id")
		is.Contains(columns, "exo_chunk.mode_id")
		is.Contains(columns, "exo_chunk.file_id")

	})
}

func TestColumns_Owl(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		model := &Owl{}

		columns, err := makroud.GetColumns(driver, model)
		is.NoError(err)
		is.NotEmpty(columns)

		is.Contains(columns, "ztp_owl.favorite_food")
		is.Contains(columns, "ztp_owl.feather_color")
		is.Contains(columns, "ztp_owl.group_id")
		is.Contains(columns, "ztp_owl.id")
		is.Contains(columns, "ztp_owl.name")

		is.Equal(fmt.Sprint("ztp_owl.favorite_food, ztp_owl.feather_color,",
			" ztp_owl.group_id, ztp_owl.id, ztp_owl.name"), columns.String())

	})
}

func TestColumns_Cat(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		model := &Cat{}

		columns, err := makroud.GetColumns(driver, model)
		is.NoError(err)
		is.NotEmpty(columns)

		is.Contains(columns, "ztp_cat.created_at")
		is.Contains(columns, "ztp_cat.deleted_at")
		is.Contains(columns, "ztp_cat.id")
		is.Contains(columns, "ztp_cat.name")
		is.Contains(columns, "ztp_cat.updated_at")

		is.Equal(fmt.Sprint("ztp_cat.created_at, ztp_cat.deleted_at, ztp_cat.id,",
			" ztp_cat.name, ztp_cat.updated_at"), columns.String())

	})
}

func TestColumns_Meow(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		model := &Meow{}

		columns, err := makroud.GetColumns(driver, model)
		is.NoError(err)
		is.NotEmpty(columns)

		is.Contains(columns, "ztp_meow.body")
		is.Contains(columns, "ztp_meow.cat_id")
		is.Contains(columns, "ztp_meow.created")
		is.Contains(columns, "ztp_meow.deleted")
		is.Contains(columns, "ztp_meow.hash")
		is.Contains(columns, "ztp_meow.updated")

		is.Equal(fmt.Sprint("ztp_meow.body, ztp_meow.cat_id, ztp_meow.created, ztp_meow.deleted,",
			" ztp_meow.hash, ztp_meow.updated"), columns.String())

	})
}

func TestColumns_ExoChunk(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		model := &ExoChunk{}

		columns, err := makroud.GetColumns(driver, model)
		is.NoError(err)
		is.NotEmpty(columns)

		is.Contains(columns, "exo_chunk.bytes")
		is.Contains(columns, "exo_chunk.file_id")
		is.Contains(columns, "exo_chunk.hash")
		is.Contains(columns, "exo_chunk.mode_id")
		is.Contains(columns, "exo_chunk.organization_id")
		is.Contains(columns, "exo_chunk.user_id")

		is.Equal(fmt.Sprint("exo_chunk.bytes, exo_chunk.file_id, exo_chunk.hash, exo_chunk.mode_id, ",
			"exo_chunk.organization_id, exo_chunk.user_id"), columns.String())

	})
}
