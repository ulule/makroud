package sqlxx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

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

		fmt.Printf("%+v\n", schema)

	})
}
