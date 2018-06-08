package sqlxx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

import "time"
import "github.com/lib/pq"
import "github.com/ulule/amber/registry"

// import "reflect"
// import "github.com/ulule/sqlxx/reflectx"

func TestSchema_Owl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		model := &Owl{}

		schema, err := sqlxx.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		is.Equal("wp_owl", schema.TableName())
		is.Equal("id", schema.PrimaryKey().ColumnName())
		is.Equal("wp_owl.id", schema.PrimaryKey().ColumnPath())
		is.False(schema.HasCreatedKey())
		is.False(schema.HasUpdatedKey())
		is.False(schema.HasDeletedKey())

		columns := schema.Columns()
		is.Len(columns, 4)
		is.Contains(columns, "id")
		is.Contains(columns, "name")
		is.Contains(columns, "feather_color")
		is.Contains(columns, "favorite_food")

		columns = schema.ColumnPaths()
		is.Len(columns, 4)
		is.Contains(columns, "wp_owl.id")
		is.Contains(columns, "wp_owl.name")
		is.Contains(columns, "wp_owl.feather_color")
		is.Contains(columns, "wp_owl.favorite_food")

		// TODO REMOVE
		fmt.Printf("%+v\n", schema)

	})
}

func TestSchema_Cat(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		model := &Cat{}

		schema, err := sqlxx.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		// TODO REMOVE
		fmt.Printf("%+v\n", schema)

		is.Equal("wp_cat", schema.TableName())
		is.Equal("id", schema.PrimaryKey().ColumnName())
		is.Equal("wp_cat.id", schema.PrimaryKey().ColumnPath())

		is.True(schema.HasCreatedKey())
		is.True(schema.HasUpdatedKey())
		is.True(schema.HasDeletedKey())
		is.Equal("wp_cat.created_at", schema.CreatedKeyPath())
		is.Equal("wp_cat.updated_at", schema.UpdatedKeyPath())
		is.Equal("wp_cat.deleted_at", schema.DeletedKeyPath())

		columns := schema.Columns()
		is.Len(columns, 5)
		is.Contains(columns, "id")
		is.Contains(columns, "name")
		is.Contains(columns, "created_at")
		is.Contains(columns, "updated_at")
		is.Contains(columns, "deleted_at")

		columns = schema.ColumnPaths()
		is.Len(columns, 5)
		is.Contains(columns, "wp_cat.id")
		is.Contains(columns, "wp_cat.name")
		is.Contains(columns, "wp_cat.created_at")
		is.Contains(columns, "wp_cat.updated_at")
		is.Contains(columns, "wp_cat.deleted_at")

	})
}

func TestSchema_Meow(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		model := &Meow{}

		schema, err := sqlxx.GetSchema(driver, model)
		is.NoError(err)
		is.NotNil(schema)

		is.Equal("wp_meow", schema.TableName())
		is.Equal("hash", schema.PrimaryKey().ColumnName())
		is.Equal("wp_meow.hash", schema.PrimaryKey().ColumnPath())

		is.True(schema.HasCreatedKey())
		is.True(schema.HasUpdatedKey())
		is.True(schema.HasDeletedKey())
		is.Equal("wp_meow.created", schema.CreatedKeyPath())
		is.Equal("wp_meow.updated", schema.UpdatedKeyPath())
		is.Equal("wp_meow.deleted", schema.DeletedKeyPath())

		columns := schema.Columns()
		is.Len(columns, 5)
		is.Contains(columns, "hash")
		is.Contains(columns, "body")
		is.Contains(columns, "created")
		is.Contains(columns, "updated")
		is.Contains(columns, "deleted")

		columns = schema.ColumnPaths()
		is.Len(columns, 5)
		is.Contains(columns, "wp_meow.hash")
		is.Contains(columns, "wp_meow.body")
		is.Contains(columns, "wp_meow.created")
		is.Contains(columns, "wp_meow.updated")
		is.Contains(columns, "wp_meow.deleted")

	})
}

var model *Meow

func BenchmarkSchemaWriteModel1(b *testing.B) {
	var m *Meow
	Setup(b)(func(driver sqlxx.Driver) {
		is := require.New(b)

		err := sqlxx.Save(driver, &Meow{
			Body: "meow!",
		})
		is.NoError(err)

		stmt, err := driver.PrepareNamed("SELECT * FROM wp_meow LIMIT 1;")
		is.NoError(err)

		h := registry.NewHistogram("duration")

		for n := 0; n < b.N; n++ {

			type Meow3 struct {
				Hash      string      `db:"hash"`
				Body      string      `db:"body"`
				CreatedAt time.Time   `db:"created"`
				UpdatedAt time.Time   `db:"updated"`
				DeletedAt pq.NullTime `db:"deleted"`
			}

			model := &Meow3{}
			row := stmt.QueryRowx(map[string]interface{}{})
			t0 := time.Now()
			err = row.StructScan(model)
			h.Observe(int64(time.Since(t0)))
			is.NoError(err)
		}

		fmt.Println(time.Duration(h.Mean()))
	})
	model = m
}

func BenchmarkSchemaWriteModel2(b *testing.B) {
	var m *Meow
	Setup(b)(func(driver sqlxx.Driver) {
		is := require.New(b)
		schema, err := sqlxx.GetSchema(driver, &Meow{})
		is.NoError(err)

		h := registry.NewHistogram("duration")

		for n := 0; n < b.N; n++ {

			model := &Meow{}

			//hash := "01CEXVDPT7YAAAV9J3CWWAH8KP"
			body := "meow meow!"
			// created := time.Now()
			// updated := time.Now()
			// deleted := time.Now()

			mapper := sqlxx.Mapper{
				// "hash": hash,
				"body": body,
				// "created": created,
				// "updated": updated,
				// "deleted": deleted,
			}

			t0 := time.Now()
			schema.WriteModel(mapper, model)
			h.Observe(int64(time.Since(t0)))

			m = model
		}

		fmt.Println(time.Duration(h.Mean()))
	})
	model = m
}
