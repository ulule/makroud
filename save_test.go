package sqlxx_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/format"

	"github.com/ulule/sqlxx"
)

func TestSave_Owl(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
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
			"INSERT INTO ztp_owl (favorite_food, feather_color, group_id, name) ",
			"VALUES ('Tomato', 'white', NULL, 'Kika') RETURNING id",
		)
		is.Equal(expected, query.Raw)
		expected = fmt.Sprint(
			"INSERT INTO ztp_owl (favorite_food, feather_color, group_id, name) ",
			"VALUES (:arg_1, :arg_2, :arg_3, :arg_4) RETURNING id",
		)
		is.Equal(expected, query.Query)
		is.Len(query.Args, 4)
		is.Equal(favoriteFood, query.Args["arg_1"])
		is.Equal(featherColor, query.Args["arg_2"])
		is.Equal(sql.NullInt64{}, query.Args["arg_3"])
		is.Equal(name, query.Args["arg_4"])

		favoriteFood = "Chocolate Cake"
		owl.FavoriteFood = favoriteFood
		id := owl.ID

		queries, err = sqlxx.SaveWithQueries(driver, owl)
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		query = queries[0]
		expected = fmt.Sprint(
			"UPDATE ztp_owl SET favorite_food = 'Chocolate Cake', feather_color = 'white', ",
			"group_id = NULL, name = 'Kika' WHERE (id = 1)",
		)
		is.Equal(expected, query.Raw)
		expected = fmt.Sprint(
			"UPDATE ztp_owl SET favorite_food = :arg_1, feather_color = :arg_2, ",
			"group_id = :arg_3, name = :arg_4 WHERE (id = :arg_5)",
		)
		is.Equal(expected, query.Query)
		is.Len(query.Args, 5)
		is.Equal(favoriteFood, query.Args["arg_1"])
		is.Equal(featherColor, query.Args["arg_2"])
		is.Equal(sql.NullInt64{}, query.Args["arg_3"])
		is.Equal(name, query.Args["arg_4"])
		is.Equal(id, query.Args["arg_5"])

		check := loukoum.Select("COUNT(*)").From("ztp_owl").Where(loukoum.Condition("name").Equal("Kika"))
		count := -1
		err = sqlxx.Exec(driver, check, &count)
		is.NoError(err)
		is.NoError(err)
		is.Equal(1, count)

	})
}

func TestSave_Meow(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)

		cat := &Cat{
			Name: "Hemlock",
		}
		err := sqlxx.Save(driver, cat)
		is.NoError(err)

		t0 := time.Now()
		body := "meow"
		meow := &Meow{
			Body:  body,
			CatID: cat.ID,
		}

		queries, err := sqlxx.SaveWithQueries(driver, meow)
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		query := queries[0]
		expected := fmt.Sprint(
			"INSERT INTO ztp_meow (body, cat_id, deleted, hash) VALUES (", format.String(body), ", ",
			format.String(cat.ID), ", NULL, ",
			format.String(meow.Hash), ") RETURNING created, updated",
		)
		is.Equal(expected, query.Raw)
		expected = fmt.Sprint(
			"INSERT INTO ztp_meow (body, cat_id, deleted, hash) VALUES (:arg_1, :arg_2, :arg_3, :arg_4) ",
			"RETURNING created, updated",
		)
		is.Equal(expected, query.Query)
		is.Len(query.Args, 4)
		is.Equal(body, query.Args["arg_1"])
		is.Equal(cat.ID, query.Args["arg_2"])
		is.Equal(pq.NullTime{}, query.Args["arg_3"])
		is.Equal(meow.Hash, query.Args["arg_4"])
		is.False(meow.CreatedAt.IsZero())
		is.True(meow.CreatedAt.After(t0))
		is.False(meow.UpdatedAt.IsZero())
		is.True(meow.UpdatedAt.After(t0))

		t1 := time.Now()
		body = "meow meow!"
		meow.Body = body

		queries, err = sqlxx.SaveWithQueries(driver, meow)
		is.NoError(err)
		is.NotNil(queries)
		is.Len(queries, 1)
		query = queries[0]
		expected = fmt.Sprint(
			"UPDATE ztp_meow SET body = ", format.String(body), ", cat_id = ", format.String(cat.ID),
			", created = ", format.Time(meow.CreatedAt), ", deleted = NULL, updated = NOW() WHERE (hash = ",
			format.String(meow.Hash), ") RETURNING updated",
		)
		is.Equal(expected, query.Raw)
		expected = fmt.Sprint(
			"UPDATE ztp_meow SET body = :arg_1, cat_id = :arg_2, created = :arg_3, deleted = :arg_4, updated = NOW() ",
			"WHERE (hash = :arg_5) RETURNING updated",
		)
		is.Equal(expected, query.Query)
		is.Len(query.Args, 5)
		is.Equal(body, query.Args["arg_1"])
		is.Equal(cat.ID, query.Args["arg_2"])
		is.Equal(meow.CreatedAt, query.Args["arg_3"])
		is.Equal(pq.NullTime{}, query.Args["arg_4"])
		is.Equal(meow.Hash, query.Args["arg_5"])
		is.False(meow.CreatedAt.IsZero())
		is.True(meow.CreatedAt.After(t0))
		is.True(meow.CreatedAt.Before(t1))
		is.False(meow.UpdatedAt.IsZero())
		is.True(meow.UpdatedAt.After(t0))
		is.True(meow.UpdatedAt.After(t1))

	})
}
