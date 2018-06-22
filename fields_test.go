package sqlxx_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestFields_Owl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		schema := &sqlxx.Schema{}
		model := Owl{}

		field, err := sqlxx.NewField(driver, schema, model, "tracking")
		is.NoError(err)
		is.Equal("Owl", field.ModelName())
		is.Equal("wp_owl", field.TableName())
		is.Equal("tracking", field.ColumnName())
		is.Equal("wp_owl.tracking", field.ColumnPath())
		is.True(field.IsExcluded())
		is.False(field.IsPrimaryKey())
		is.Equal(reflect.Bool, field.Type().Kind())

		field, err = sqlxx.NewField(driver, schema, model, "FavoriteFood")
		is.NoError(err)
		is.Equal("Owl", field.ModelName())
		is.Equal("wp_owl", field.TableName())
		is.Equal("favorite_food", field.ColumnName())
		is.Equal("wp_owl.favorite_food", field.ColumnPath())
		is.False(field.IsExcluded())
		is.False(field.IsPrimaryKey())
		is.Equal(reflect.String, field.Type().Kind())

		field, err = sqlxx.NewField(driver, schema, model, "FeatherColor")
		is.NoError(err)
		is.Equal("Owl", field.ModelName())
		is.Equal("wp_owl", field.TableName())
		is.Equal("feather_color", field.ColumnName())
		is.Equal("wp_owl.feather_color", field.ColumnPath())
		is.False(field.IsExcluded())
		is.False(field.IsPrimaryKey())
		is.Equal(reflect.String, field.Type().Kind())

		field, err = sqlxx.NewField(driver, schema, model, "Name")
		is.NoError(err)
		is.Equal("Owl", field.ModelName())
		is.Equal("wp_owl", field.TableName())
		is.Equal("name", field.ColumnName())
		is.Equal("wp_owl.name", field.ColumnPath())
		is.False(field.IsExcluded())
		is.False(field.IsPrimaryKey())
		is.Equal(reflect.String, field.Type().Kind())

		field, err = sqlxx.NewField(driver, schema, model, "ID")
		is.NoError(err)
		is.Equal("Owl", field.ModelName())
		is.Equal("wp_owl", field.TableName())
		is.Equal("id", field.ColumnName())
		is.Equal("wp_owl.id", field.ColumnPath())
		is.False(field.IsExcluded())
		is.True(field.IsPrimaryKey())
		is.Equal(reflect.Int64, field.Type().Kind())

	})
}
