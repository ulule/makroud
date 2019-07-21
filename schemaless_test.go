package makroud_test

import (
	"testing"

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
