package makroud_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum/v3"

	"github.com/ulule/makroud"
)

func TestDelete_DeleteOwl(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		owl := &Owl{
			Name:         "Blake",
			FeatherColor: "brown",
			FavoriteFood: "Raspberry",
		}

		err := makroud.Save(ctx, driver, owl)
		is.NoError(err)
		is.NotEmpty(owl.ID)

		id := owl.ID

		err = makroud.Delete(ctx, driver, owl)
		is.NoError(err)

		query := loukoum.Select("COUNT(*)").From("ztp_owl").Where(loukoum.Condition("id").Equal(id))
		count, err := makroud.Count(ctx, driver, query)
		is.NoError(err)
		is.Equal(int64(0), count)

		query = loukoum.Select("COUNT(*)").From("ztp_owl").Where(loukoum.Condition("name").Equal("Blake"))
		count, err = makroud.Count(ctx, driver, query)
		is.NoError(err)
		is.NoError(err)
		is.Equal(int64(0), count)

	})
}

func TestDelete_ArchiveOwl(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		owl := &Owl{
			Name:         "Frosty",
			FeatherColor: "beige",
			FavoriteFood: "Wasabi",
		}

		err := makroud.Save(ctx, driver, owl)
		is.NoError(err)
		is.NotEmpty(owl.ID)

		err = makroud.Archive(ctx, driver, owl)
		is.Error(err)
		is.Equal(makroud.ErrSchemaDeletedKey, errors.Cause(err))

	})
}

func TestDelete_DeleteMeow(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{
			Name: "Wolfram",
		}

		err := makroud.Save(ctx, driver, cat)
		is.NoError(err)
		is.NotEmpty(cat.ID)

		meow := &Meow{
			Body:  "meow meow meow?",
			CatID: cat.ID,
		}

		err = makroud.Save(ctx, driver, meow)
		is.NoError(err)
		is.NotEmpty(meow.Hash)

		id := meow.Hash

		err = makroud.Delete(ctx, driver, meow)
		is.NoError(err)

		query := loukoum.Select("COUNT(*)").From("ztp_meow").Where(loukoum.Condition("hash").Equal(id))
		count, err := makroud.Count(ctx, driver, query)
		is.NoError(err)
		is.Equal(int64(0), count)

	})
}

func TestDelete_ArchiveMeow(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{
			Name: "Wolfram",
		}

		err := makroud.Save(ctx, driver, cat)
		is.NoError(err)
		is.NotEmpty(cat.ID)

		meow := &Meow{
			Body:  "meow! meow meow meow ?!",
			CatID: cat.ID,
		}

		err = makroud.Save(ctx, driver, meow)
		is.NoError(err)
		is.NotEmpty(meow.Hash)

		id := meow.Hash

		err = makroud.Archive(ctx, driver, meow)
		is.NoError(err)

		query := loukoum.Select("COUNT(*)").From("ztp_meow").Where(loukoum.Condition("hash").Equal(id))
		count, err := makroud.Count(ctx, driver, query)
		is.NoError(err)
		is.Equal(int64(1), count)

		query = loukoum.Select("COUNT(*)").From("ztp_meow").
			Where(loukoum.Condition("hash").Equal(id)).And(loukoum.Condition("deleted").IsNull(true))
		count, err = makroud.Count(ctx, driver, query)
		is.NoError(err)
		is.Equal(int64(0), count)

	})
}
