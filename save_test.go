package sqlxx_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"

	"github.com/ulule/sqlxx"
)

func TestSave_Owl(t *testing.T) {
	Setup(t, sqlxx.Cache(true))(func(driver sqlxx.Driver) {
		is := require.New(t)

		name := "Kika"
		featherColor := "white"
		favoriteFood := "Tomato"
		owl := &Owl{
			Name:         name,
			FeatherColor: featherColor,
			FavoriteFood: favoriteFood,
		}

		queries, err := sqlxx.SaveWithQueries(driver, owl)
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		query := queries[0]
		expected := fmt.Sprint(
			"INSERT INTO wp_owl (favorite_food, feather_color, name) ",
			"VALUES ('Tomato', 'white', 'Kika') RETURNING id",
		)
		is.Equal(expected, query.Raw)
		expected = fmt.Sprint(
			"INSERT INTO wp_owl (favorite_food, feather_color, name) ",
			"VALUES (:arg_1, :arg_2, :arg_3) RETURNING id",
		)
		is.Equal(expected, query.Query)
		is.Len(query.Args, 3)
		is.Equal(favoriteFood, query.Args["arg_1"])
		is.Equal(featherColor, query.Args["arg_2"])
		is.Equal(name, query.Args["arg_3"])

		favoriteFood = "Chocolate Cake"
		owl.FavoriteFood = favoriteFood
		id := owl.ID

		queries, err = sqlxx.SaveWithQueries(driver, owl)
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		query = queries[0]
		expected = fmt.Sprint(
			"UPDATE wp_owl SET favorite_food = 'Chocolate Cake', feather_color = 'white', ",
			"name = 'Kika' WHERE (id = 1)",
		)
		is.Equal(expected, query.Raw)
		expected = fmt.Sprint(
			"UPDATE wp_owl SET favorite_food = :arg_1, feather_color = :arg_2, ",
			"name = :arg_3 WHERE (id = :arg_4)",
		)
		is.Equal(expected, query.Query)
		is.Len(query.Args, 4)
		is.Equal(favoriteFood, query.Args["arg_1"])
		is.Equal(featherColor, query.Args["arg_2"])
		is.Equal(name, query.Args["arg_3"])
		is.Equal(id, query.Args["arg_4"])

		check := loukoum.Select("COUNT(*)").From("wp_owl").Where(loukoum.Condition("name").Equal("Kika"))
		count := -1
		err = sqlxx.Fetch(driver, check, &count)
		is.NoError(err)
		is.NoError(err)
		is.Equal(1, count)

	})
}

func TestSave_Meow(t *testing.T) {
	Setup(t, sqlxx.Cache(true))(func(driver sqlxx.Driver) {
		is := require.New(t)

		body := "meow"
		now := time.Now()
		meow := &Meow{
			Body: body,
		}

		queries, err := sqlxx.SaveWithQueries(driver, meow)
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		query := queries[0]
		expected := fmt.Sprint(
			"INSERT INTO wp_meow (body, deleted, hash) VALUES ('", body, "', NULL, '",
			meow.Hash, "') RETURNING created, updated",
		)
		is.Equal(expected, query.Raw)
		expected = fmt.Sprint(
			"INSERT INTO wp_meow (body, deleted, hash) VALUES (:arg_1, :arg_2, :arg_3) ",
			"RETURNING created, updated",
		)
		is.Equal(expected, query.Query)
		is.Len(query.Args, 3)
		is.Equal(body, query.Args["arg_1"])
		is.Equal(pq.NullTime{}, query.Args["arg_2"])
		is.Equal(meow.Hash, query.Args["arg_3"])
		is.False(meow.CreatedAt.IsZero())
		is.True(meow.CreatedAt.After(now))
		is.False(meow.UpdatedAt.IsZero())
		is.True(meow.UpdatedAt.After(now))

	})
}
