package sqlxx_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"

	"github.com/ulule/sqlxx"
)

func TestDelete_DeleteOwl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		owl := &Owl{
			Name:         "Blake",
			FeatherColor: "brown",
			FavoriteFood: "Raspberry",
		}

		err := sqlxx.Save(ctx, driver, owl)
		is.NoError(err)
		is.NotEmpty(owl.ID)

		id := owl.ID

		err = sqlxx.Delete(ctx, driver, owl)
		is.NoError(err)

		query := loukoum.Select("COUNT(*)").From("ztp_owl").Where(loukoum.Condition("id").Equal(id))
		count, err := sqlxx.Count(ctx, driver, query)
		is.NoError(err)
		is.Equal(int64(0), count)

		query = loukoum.Select("COUNT(*)").From("ztp_owl").Where(loukoum.Condition("name").Equal("Blake"))
		count, err = sqlxx.Count(ctx, driver, query)
		is.NoError(err)
		is.NoError(err)
		is.Equal(int64(0), count)

	})
}

func TestDelete_ArchiveOwl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		owl := &Owl{
			Name:         "Frosty",
			FeatherColor: "beige",
			FavoriteFood: "Wasabi",
		}

		err := sqlxx.Save(ctx, driver, owl)
		is.NoError(err)
		is.NotEmpty(owl.ID)

		err = sqlxx.Archive(ctx, driver, owl)
		is.Error(err)
		is.Equal(sqlxx.ErrSchemaDeletedKey, errors.Cause(err))

	})
}

func TestDelete_DeleteMeow(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{
			Name: "Wolfram",
		}

		err := sqlxx.Save(ctx, driver, cat)
		is.NoError(err)
		is.NotEmpty(cat.ID)

		meow := &Meow{
			Body:  "meow meow meow?",
			CatID: cat.ID,
		}

		err = sqlxx.Save(ctx, driver, meow)
		is.NoError(err)
		is.NotEmpty(meow.Hash)

		id := meow.Hash

		err = sqlxx.Delete(ctx, driver, meow)
		is.NoError(err)

		query := loukoum.Select("COUNT(*)").From("ztp_meow").Where(loukoum.Condition("hash").Equal(id))
		count, err := sqlxx.Count(ctx, driver, query)
		is.NoError(err)
		is.Equal(int64(0), count)

	})
}

func TestDelete_ArchiveMeow(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{
			Name: "Wolfram",
		}

		err := sqlxx.Save(ctx, driver, cat)
		is.NoError(err)
		is.NotEmpty(cat.ID)

		meow := &Meow{
			Body:  "meow! meow meow meow ?!",
			CatID: cat.ID,
		}

		err = sqlxx.Save(ctx, driver, meow)
		is.NoError(err)
		is.NotEmpty(meow.Hash)

		id := meow.Hash

		err = sqlxx.Archive(ctx, driver, meow)
		is.NoError(err)

		query := loukoum.Select("COUNT(*)").From("ztp_meow").Where(loukoum.Condition("hash").Equal(id))
		count, err := sqlxx.Count(ctx, driver, query)
		is.NoError(err)
		is.Equal(int64(1), count)

		query = loukoum.Select("COUNT(*)").From("ztp_meow").
			Where(loukoum.Condition("hash").Equal(id)).And(loukoum.Condition("deleted").IsNull(true))
		count, err = sqlxx.Count(ctx, driver, query)
		is.NoError(err)
		is.Equal(int64(0), count)

	})
}
