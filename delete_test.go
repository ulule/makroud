package sqlxx_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/format"

	"github.com/ulule/sqlxx"
)

func TestDelete_DeleteOwl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		owl := &Owl{
			Name:         "Blake",
			FeatherColor: "brown",
			FavoriteFood: "Raspberry",
		}

		err := sqlxx.Save(driver, owl)
		is.NoError(err)
		id := owl.ID

		queries, err := sqlxx.DeleteWithQueries(driver, owl)
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		query := queries[0]
		expected := fmt.Sprintf("DELETE FROM ztp_owl WHERE (id = %s)", format.Int(owl.ID))
		is.Equal(expected, query.Raw)
		expected = "DELETE FROM ztp_owl WHERE (id = :arg_1)"
		is.Equal(expected, query.Query)
		is.Len(query.Args, 1)
		is.Equal(id, query.Args["arg_1"])

		check := loukoum.Select("COUNT(*)").From("ztp_owl").Where(loukoum.Condition("name").Equal("Blake"))
		count := -1
		err = sqlxx.Exec(driver, check, &count)
		is.NoError(err)
		is.NoError(err)
		is.Equal(0, count)

	})
}

func TestDelete_ArchiveOwl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		owl := &Owl{
			Name:         "Frosty",
			FeatherColor: "beige",
			FavoriteFood: "Wasabi",
		}

		err := sqlxx.Save(driver, owl)
		is.NoError(err)

		queries, err := sqlxx.ArchiveWithQueries(driver, owl)
		is.Error(err)
		is.Nil(queries)
		is.Equal(sqlxx.ErrSchemaDeletedKey, errors.Cause(err))

	})
}

func TestDelete_DeleteMeow(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		meow := &Meow{
			Body: "meow meow meow?",
		}

		err := sqlxx.Save(driver, meow)
		is.NoError(err)
		id := meow.Hash

		queries, err := sqlxx.DeleteWithQueries(driver, meow)
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		query := queries[0]
		expected := fmt.Sprintf("DELETE FROM ztp_meow WHERE (hash = %s)", format.String(meow.Hash))
		is.Equal(expected, query.Raw)
		expected = "DELETE FROM ztp_meow WHERE (hash = :arg_1)"
		is.Equal(expected, query.Query)
		is.Len(query.Args, 1)
		is.Equal(id, query.Args["arg_1"])

		check := loukoum.Select("COUNT(*)").From("ztp_meow").Where(loukoum.Condition("hash").Equal(id))
		count := -1
		err = sqlxx.Exec(driver, check, &count)
		is.NoError(err)
		is.NoError(err)
		is.Equal(0, count)

	})
}

func TestDelete_ArchiveMeow(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)
		meow := &Meow{
			Body: "meow! meow meow meow ?!",
		}

		err := sqlxx.Save(driver, meow)
		is.NoError(err)
		id := meow.Hash

		queries, err := sqlxx.ArchiveWithQueries(driver, meow)
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		query := queries[0]
		expected := fmt.Sprintf(
			"UPDATE ztp_meow SET deleted = NOW() WHERE (hash = %s) RETURNING deleted",
			format.String(meow.Hash),
		)
		is.Equal(expected, query.Raw)
		expected = "UPDATE ztp_meow SET deleted = NOW() WHERE (hash = :arg_1) RETURNING deleted"
		is.Equal(expected, query.Query)
		is.Len(query.Args, 1)
		is.Equal(id, query.Args["arg_1"])

		count := -1
		check := loukoum.Select("COUNT(*)").From("ztp_meow").
			Where(loukoum.Condition("hash").Equal(id))
		err = sqlxx.Exec(driver, check, &count)
		is.NoError(err)
		is.NoError(err)
		is.Equal(1, count)

		count = -1
		check = loukoum.Select("COUNT(*)").From("ztp_meow").
			Where(loukoum.Condition("hash").Equal(id)).And(loukoum.Condition("deleted").IsNull(true))
		err = sqlxx.Exec(driver, check, &count)
		is.NoError(err)
		is.NoError(err)
		is.Equal(0, count)

	})
}
