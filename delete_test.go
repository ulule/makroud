package sqlxx_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"

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
		expected := fmt.Sprintf("DELETE FROM wp_owl WHERE (id = %d)", owl.ID)
		is.Equal(expected, query.Raw)
		expected = "DELETE FROM wp_owl WHERE (id = :arg_1)"
		is.Equal(expected, query.Query)
		is.Len(query.Args, 1)
		is.Equal(id, query.Args["arg_1"])

		check := loukoum.Select("COUNT(*)").From("wp_owl").Where(loukoum.Condition("name").Equal("Blake"))
		count := -1
		err = sqlxx.Fetch(driver, check, &count)
		is.NoError(err)
		is.NoError(err)
		is.Equal(0, count)

	})
}

// func TestDelete_ArchiveChunk(t *testing.T) {
// 	Setup(t)(func(driver sqlxx.Driver) {
// 		is := require.New(t)
// 		owl := &Owl{
// 			Name:         "Blake",
// 			FeatherColor: "brown",
// 			FavoriteFood: "Raspberry",
// 		}
//
// 		err := sqlxx.Save(driver, owl)
// 		is.NoError(err)
// 		id := owl.ID
//
// 		queries, err := sqlxx.DeleteWithQueries(driver, owl)
// 		is.NoError(err)
// 		is.NotNil(queries)
// 		is.Len(queries, 1)
// 		query := queries[0]
// 		expected := fmt.Sprintf("DELETE FROM wp_owl WHERE (id = %d)", owl.ID)
// 		is.Equal(expected, query.Raw)
// 		expected = "DELETE FROM wp_owl WHERE (id = :arg_1)"
// 		is.Equal(expected, query.Query)
// 		is.Len(query.Args, 1)
// 		is.Equal(id, query.Args["arg_1"])
//
// 		check := loukoum.Select("COUNT(*)").From("wp_owl").Where(loukoum.Condition("name").Equal("Blake"))
// 		count := -1
// 		err = sqlxx.Fetch(driver, check, &count)
// 		is.NoError(err)
// 		is.NoError(err)
// 		is.Equal(0, count)
//
// 	})
// }

//
// func TestDelete_ArchiveOwl(t *testing.T) {
// 	env := setup(t)
// 	defer env.teardown()
//
// 	is := require.New(t)
//
// 	user := &User{Username: "thoas"}
// 	err := sqlxx.Save(env.driver, user)
// 	is.NoError(err)
//
// 	queries, err := sqlxx.ArchiveWithQueries(env.driver, user)
// 	is.NoError(err)
// 	is.NotNil(queries)
// 	is.Len(queries, 1)
// 	is.Contains(queries[0].Query, "UPDATE users SET deleted_at = :deleted_at WHERE users.id = :id")
//
// 	query := `
// 		SELECT COUNT(*)
// 		FROM users
// 		WHERE username = :username
// 		AND deleted_at IS NULL
// 	`
// 	params := map[string]interface{}{
// 		"username": "thoas",
// 	}
//
// 	stmt, err := env.driver.PrepareNamed(query)
// 	is.NoError(err)
// 	is.NotNil(stmt)
//
// 	count := -1
// 	err = stmt.Get(&count, params)
// 	is.NoError(err)
// 	is.Equal(0, count)
// }
