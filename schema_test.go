package makroud_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/pkg/errors"
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

		is.True(schema.HasColumn("id"))
		is.True(schema.HasColumn("name"))
		is.True(schema.HasColumn("feather_color"))
		is.True(schema.HasColumn("favorite_food"))
		is.True(schema.HasColumn("group_id"))
		is.False(schema.HasColumn("human_id"))
		is.False(schema.HasColumn("eyesight"))

		columns = schema.ColumnPaths()
		is.Len(columns, 5)
		is.Contains(columns, "ztp_owl.id")
		is.Contains(columns, "ztp_owl.name")
		is.Contains(columns, "ztp_owl.feather_color")
		is.Contains(columns, "ztp_owl.favorite_food")
		is.Contains(columns, "ztp_owl.group_id")

		is.True(schema.HasColumn("ztp_owl.id"))
		is.True(schema.HasColumn("ztp_owl.name"))
		is.True(schema.HasColumn("ztp_owl.feather_color"))
		is.True(schema.HasColumn("ztp_owl.favorite_food"))
		is.True(schema.HasColumn("ztp_owl.group_id"))
		is.False(schema.HasColumn("ztp_owl.human_id"))
		is.False(schema.HasColumn("ztp_owl.eyesight"))

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

		is.True(schema.HasColumn("id"))
		is.True(schema.HasColumn("name"))
		is.True(schema.HasColumn("created_at"))
		is.True(schema.HasColumn("updated_at"))
		is.True(schema.HasColumn("deleted_at"))
		is.False(schema.HasColumn("human_id"))
		is.False(schema.HasColumn("favorite_food"))

		columns = schema.ColumnPaths()
		is.Len(columns, 5)
		is.Contains(columns, "ztp_cat.id")
		is.Contains(columns, "ztp_cat.name")
		is.Contains(columns, "ztp_cat.created_at")
		is.Contains(columns, "ztp_cat.updated_at")
		is.Contains(columns, "ztp_cat.deleted_at")

		is.True(schema.HasColumn("ztp_cat.id"))
		is.True(schema.HasColumn("ztp_cat.name"))
		is.True(schema.HasColumn("ztp_cat.created_at"))
		is.True(schema.HasColumn("ztp_cat.updated_at"))
		is.True(schema.HasColumn("ztp_cat.deleted_at"))
		is.False(schema.HasColumn("ztp_cat.human_id"))
		is.False(schema.HasColumn("ztp_cat.favorite_food"))
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

func TestSchema_Human(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		human1 := &Human{Name: "Aurélien"}
		err := makroud.Save(ctx, driver, human1)
		is.NoError(err)

		human2 := &Human{Name: "Carl"}
		err = makroud.Save(ctx, driver, human2)
		is.NoError(err)

		human3 := &Human{Name: "Clémence"}
		err = makroud.Save(ctx, driver, human3)
		is.NoError(err)

		human4 := &Human{Name: "Renzo"}
		err = makroud.Save(ctx, driver, human4)
		is.NoError(err)

		verifyHuman := func(expected *Human, actual *Human) {
			is.Equal(expected.ID, actual.ID)
			is.Equal(expected.Name, actual.Name)
			is.Equal(expected.CreatedAt.UnixNano(), actual.CreatedAt.UnixNano())
			is.Equal(expected.UpdatedAt.UnixNano(), actual.UpdatedAt.UnixNano())
			is.Equal(expected.DeletedAt.Valid, actual.DeletedAt.Valid)
			is.Equal(expected.DeletedAt.Time.UnixNano(), actual.DeletedAt.Time.UnixNano())
			is.Equal(expected.CatID.Valid, actual.CatID.Valid)
			is.Equal(expected.CatID.String, actual.CatID.String)
		}

		{
			schema, err := makroud.GetSchema(driver, &Human{})
			is.NoError(err)
			is.NotEmpty(schema)

			stmt, err := driver.Prepare(ctx, `SELECT * FROM ztp_human WHERE id = $1`)
			is.NoError(err)
			is.NotEmpty(stmt)

			row, err := stmt.QueryRow(ctx, human3.ID)
			is.NoError(err)
			is.NotEmpty(row)

			human := &Human{}
			err = schema.ScanRow(row, human)
			is.NoError(err)
			verifyHuman(human3, human)

			err = schema.ScanRow(row, human)
			is.Error(err)

			row, err = stmt.QueryRow(ctx, human4.ID)
			is.NoError(err)
			is.NotEmpty(row)

			var model makroud.Model

			err = schema.ScanRow(row, model)
			is.Error(err)
			is.Equal(makroud.ErrStructRequired, errors.Cause(err))

			err = stmt.Close()
			is.NoError(err)
		}

		{
			schema, err := makroud.GetSchema(driver, &Human{})
			is.NoError(err)
			is.NotEmpty(schema)

			stmt, err := driver.Prepare(ctx, `SELECT * FROM ztp_human`)
			is.NoError(err)
			is.NotEmpty(stmt)

			rows, err := stmt.QueryRows(ctx)
			is.NoError(err)
			is.NotEmpty(rows)

			expected := []string{
				human1.ID,
				human2.ID,
				human3.ID,
				human4.ID,
			}

			resp := map[string]*Human{}

			for rows.Next() {
				human := &Human{}
				err = schema.ScanRows(rows, human)
				is.NoError(err)
				resp[human.ID] = human
			}

			err = rows.Err()
			is.NoError(err)

			err = schema.ScanRows(rows, &Human{})
			is.Error(err)

			err = rows.Close()
			is.NoError(err)

			for _, id := range expected {
				is.NotEmpty(resp[id])
				human := resp[id]
				switch id {
				case human1.ID:
					verifyHuman(human1, human)
				case human2.ID:
					verifyHuman(human2, human)
				case human3.ID:
					verifyHuman(human3, human)
				case human4.ID:
					verifyHuman(human4, human)
				}
			}

			rows, err = stmt.QueryRows(ctx)
			is.NoError(err)
			is.NotEmpty(rows)

			var model makroud.Model

			err = schema.ScanRows(rows, model)
			is.Error(err)
			is.Equal(makroud.ErrStructRequired, errors.Cause(err))

			err = rows.Err()
			is.NoError(err)

			err = rows.Close()
			is.NoError(err)

			err = stmt.Close()
			is.NoError(err)
		}
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
