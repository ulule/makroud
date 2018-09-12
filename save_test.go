package sqlxx_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"

	"github.com/ulule/sqlxx"
)

func TestSave_Owl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		name := "Kika"
		featherColor := "white"
		favoriteFood := "Tomato"
		owl := &Owl{
			Name:         name,
			FeatherColor: featherColor,
			FavoriteFood: favoriteFood,
		}

		err := sqlxx.Save(ctx, driver, owl)
		is.NoError(err)
		is.NotEmpty(owl.ID)

		id := owl.ID

		query := loukoum.Select("*").From("ztp_owl").Where(loukoum.Condition("id").Equal(id))
		last := &Owl{}
		err = sqlxx.Exec(ctx, driver, query, last)
		is.NoError(err)
		is.Equal(name, last.Name)
		is.Equal(featherColor, last.FeatherColor)
		is.Equal(favoriteFood, last.FavoriteFood)

		favoriteFood = "Chocolate Cake"
		owl.FavoriteFood = favoriteFood

		err = sqlxx.Save(ctx, driver, owl)
		is.NoError(err)

		query = loukoum.Select("*").From("ztp_owl").Where(loukoum.Condition("id").Equal(id))
		last = &Owl{}
		err = sqlxx.Exec(ctx, driver, query, last)
		is.NoError(err)
		is.Equal(name, last.Name)
		is.Equal(featherColor, last.FeatherColor)
		is.Equal(favoriteFood, last.FavoriteFood)

	})
}

func TestSave_Meow(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{
			Name: "Hemlock",
		}

		err := sqlxx.Save(ctx, driver, cat)
		is.NoError(err)
		is.NotEmpty(cat.ID)

		t0 := time.Now()
		body := "meow"
		meow := &Meow{
			Body:  body,
			CatID: cat.ID,
		}

		err = sqlxx.Save(ctx, driver, meow)
		is.NoError(err)
		is.NotEmpty(meow.Hash)
		is.True(t0.Before(meow.CreatedAt))
		is.True(t0.Before(meow.UpdatedAt))
		is.True(time.Now().After(meow.CreatedAt))
		is.True(time.Now().After(meow.UpdatedAt))
		is.False(meow.DeletedAt.Valid)

		hash := meow.Hash

		query := loukoum.Select("*").From("ztp_meow").Where(loukoum.Condition("hash").Equal(hash))
		last := &Meow{}
		err = sqlxx.Exec(ctx, driver, query, last)
		is.NoError(err)
		is.Equal(body, last.Body)
		is.Equal(cat.ID, last.CatID)
		is.Equal(meow.CreatedAt.UnixNano(), last.CreatedAt.UnixNano())
		is.Equal(meow.UpdatedAt.UnixNano(), last.UpdatedAt.UnixNano())
		is.Equal(meow.DeletedAt.Valid, last.DeletedAt.Valid)
		is.Equal(meow.DeletedAt.Time.UnixNano(), last.DeletedAt.Time.UnixNano())

		t1 := time.Now()
		body = "meow meow!"
		meow.Body = body

		err = sqlxx.Save(ctx, driver, meow)
		is.NoError(err)
		is.True(t0.Before(meow.CreatedAt))
		is.True(t0.Before(meow.UpdatedAt))
		is.True(t1.After(meow.CreatedAt))
		is.True(t1.Before(meow.UpdatedAt))
		is.True(time.Now().After(meow.CreatedAt))
		is.True(time.Now().After(meow.UpdatedAt))

		query = loukoum.Select("*").From("ztp_meow").Where(loukoum.Condition("hash").Equal(hash))
		last = &Meow{}
		err = sqlxx.Exec(ctx, driver, query, last)
		is.NoError(err)
		is.Equal(body, last.Body)
		is.Equal(cat.ID, last.CatID)
		is.Equal(meow.CreatedAt.UnixNano(), last.CreatedAt.UnixNano())
		is.Equal(meow.UpdatedAt.UnixNano(), last.UpdatedAt.UnixNano())
		is.Equal(meow.DeletedAt.Valid, last.DeletedAt.Valid)
		is.Equal(meow.DeletedAt.Time.UnixNano(), last.DeletedAt.Time.UnixNano())

	})
}
