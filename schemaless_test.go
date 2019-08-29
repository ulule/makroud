package makroud_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ulule/makroud"
	"github.com/ulule/makroud/reflectx"
)

func TestSchemaless_Hash(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)

		type Hash struct {
			Hash1 string `makroud:"hash_1"`
			Hash2 string `mk:"hash_2"`
			Hash3 string `db:"hash_3"`
		}

		hash1 := &Hash{}
		kind1 := reflectx.GetIndirectType(hash1)

		schemaless, err := makroud.GetSchemaless(driver, kind1)
		is.NoError(err)
		is.NotNil(schemaless)

		is.Equal(kind1, schemaless.Type())

		columns := schemaless.Columns()
		is.Len(columns, 3)
		is.Contains(columns, "hash_1")
		is.Contains(columns, "hash_2")
		is.Contains(columns, "hash_3")

		is.True(schemaless.HasColumn("hash_1"))
		is.True(schemaless.HasColumn("hash_2"))
		is.True(schemaless.HasColumn("hash_3"))
		is.False(schemaless.HasColumn("hash_4"))
		is.False(schemaless.HasColumn("hash_5"))

		hash2 := Hash{}
		kind2 := reflectx.GetIndirectType(hash2)

		schemaless, err = makroud.GetSchemaless(driver, kind1)
		is.NoError(err)
		is.NotNil(schemaless)

		is.Equal(kind1, kind2)
		is.Equal(kind1, schemaless.Type())

		key1, ok := schemaless.Key("hash_1")
		is.True(ok)
		is.NotEmpty(key1)
		is.NotEmpty(key1.FieldIndex())
		is.Equal("hash_1", key1.ColumnName())
		is.Equal("Hash1", key1.FieldName())

		key2, ok := schemaless.Key("hash_2")
		is.True(ok)
		is.NotEmpty(key2)
		is.NotEmpty(key2.FieldIndex())
		is.Equal("hash_2", key2.ColumnName())
		is.Equal("Hash2", key2.FieldName())

		key3, ok := schemaless.Key("hash_3")
		is.True(ok)
		is.NotEmpty(key3)
		is.NotEmpty(key3.FieldIndex())
		is.Equal("hash_3", key3.ColumnName())
		is.Equal("Hash3", key3.FieldName())

		key4, ok := schemaless.Key("foobar")
		is.False(ok)
		is.Empty(key4)

	})
}

func TestSchemaless_PartialHuman(t *testing.T) {
	Setup(t, makroud.Cache(false))(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		human1 := &Human{Name: "Flick"}
		err := makroud.Save(ctx, driver, human1)
		is.NoError(err)

		human2 := &Human{Name: "Aslanbek"}
		err = makroud.Save(ctx, driver, human2)
		is.NoError(err)

		human3 := &Human{Name: "Pepito"}
		err = makroud.Save(ctx, driver, human3)
		is.NoError(err)

		human4 := &Human{Name: "Oksana"}
		err = makroud.Save(ctx, driver, human4)
		is.NoError(err)

		type PartialHuman struct {
			Name  string `mk:"name"`
			Other int    `mk:"-"`
		}

		{
			schema, err := makroud.GetSchemaless(driver, reflectx.GetIndirectType(&PartialHuman{}))
			is.NoError(err)
			is.NotEmpty(schema)

			stmt, err := driver.Prepare(ctx, `SELECT name FROM ztp_human WHERE id = $1`)
			is.NoError(err)
			is.NotEmpty(stmt)

			row, err := stmt.QueryRow(ctx, human1.ID)
			is.NoError(err)
			is.NotEmpty(row)

			human := &PartialHuman{
				Name:  "",
				Other: 5432,
			}

			err = schema.ScanRow(row, human)
			is.NoError(err)
			is.Equal(human1.Name, human.Name)
			is.Equal(5432, human.Other)

			err = schema.ScanRow(row, human)
			is.Error(err)

			row, err = stmt.QueryRow(ctx, human2.ID)
			is.NoError(err)
			is.NotEmpty(row)

			err = schema.ScanRow(row, []string{})
			is.Error(err)
			is.Equal(makroud.ErrStructRequired, errors.Cause(err))

			err = stmt.Close()
			is.NoError(err)
		}

		{
			schema, err := makroud.GetSchemaless(driver, reflectx.GetIndirectType(&PartialHuman{}))
			is.NoError(err)
			is.NotEmpty(schema)

			stmt, err := driver.Prepare(ctx, `SELECT name FROM ztp_human`)
			is.NoError(err)
			is.NotEmpty(stmt)

			rows, err := stmt.QueryRows(ctx)
			is.NoError(err)
			is.NotEmpty(rows)

			expected := []string{
				human1.Name,
				human2.Name,
				human3.Name,
				human4.Name,
			}

			resp := map[string]bool{}

			human := &PartialHuman{
				Name:  "",
				Other: 5432,
			}

			for rows.Next() {
				err = schema.ScanRows(rows, human)
				is.NoError(err)
				is.Contains(expected, human.Name)
				is.Equal(5432, human.Other)
				resp[human.Name] = true
			}

			err = rows.Err()
			is.NoError(err)

			err = schema.ScanRows(rows, human)
			is.Error(err)

			err = rows.Close()
			is.NoError(err)

			for _, name := range expected {
				is.True(resp[name])
			}

			rows, err = stmt.QueryRows(ctx)
			is.NoError(err)
			is.NotEmpty(rows)

			err = schema.ScanRows(rows, []string{})
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
