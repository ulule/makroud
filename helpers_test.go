package sqlxx_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"

	"github.com/ulule/sqlxx"
)

func TestCount(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cats := []Cat{
			Cat{Name: "Radio"},
			Cat{Name: "Radish"},
			Cat{Name: "Radium"},
			Cat{Name: "Radix"},
			Cat{Name: "Radman"},
			Cat{Name: "Radmilla"},
		}

		for i := range cats {
			err := sqlxx.Save(ctx, driver, &cats[i])
			is.NoError(err)
		}

		query := loukoum.Select("COUNT(*)").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Rad%"))

		{
			count, err := sqlxx.Count(ctx, driver, query)
			is.NoError(err)
			is.Equal(int64(6), count)
		}
		{
			count, err := sqlxx.FloatCount(ctx, driver, query)
			is.NoError(err)
			is.Equal(float64(6), count)
		}

	})
}

func TestExec_List(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{Name: "Wheezie"}
		cat2 := &Cat{Name: "Whimsy"}
		cat3 := &Cat{Name: "Whiskey"}
		cat4 := &Cat{Name: "Whisper"}
		cat5 := &Cat{Name: "Whitman"}
		cat6 := &Cat{Name: "Whitney"}
		cat7 := &Cat{Name: "Whistle"}
		cat8 := &Cat{Name: "Wheaton"}

		cats := []*Cat{cat1, cat2, cat3, cat4, cat5, cat6, cat7, cat8}
		expected := []*Cat{cat2, cat3, cat4, cat5, cat6, cat7}

		for i := range cats {
			err := sqlxx.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		list := []string{}
		query := loukoum.Select("id").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Whi%"))

		err := sqlxx.Exec(ctx, driver, query, &list)
		is.NoError(err)

		is.Len(list, len(expected))
		for i := range expected {
			is.Contains(list, expected[i].ID)
		}

		err = sqlxx.Exec(ctx, driver, query, []string{})
		is.Error(err)
		is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))

	})
}

func TestRawExec_List(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{Name: "Venice"}
		cat2 := &Cat{Name: "Vera"}
		cat3 := &Cat{Name: "Vermont"}
		cat4 := &Cat{Name: "Vermouth"}
		cat5 := &Cat{Name: "Versailles"}
		cat6 := &Cat{Name: "Vernetta"}
		cat7 := &Cat{Name: "Vertigo"}
		cat8 := &Cat{Name: "Venus"}

		cats := []*Cat{cat1, cat2, cat3, cat4, cat5, cat6, cat7, cat8}
		expected := []*Cat{cat2, cat3, cat4, cat5, cat6, cat7}

		for i := range cats {
			err := sqlxx.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		list := []string{}
		query := `SELECT id FROM ztp_cat WHERE name ILIKE 'Ver%'`
		err := sqlxx.RawExec(ctx, driver, query, &list)
		is.NoError(err)

		is.Len(list, len(expected))
		for i := range expected {
			is.Contains(list, expected[i].ID)
		}

		err = sqlxx.RawExec(ctx, driver, query, []string{})
		is.Error(err)
		is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))

	})
}

func TestExec_Fetch(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{Name: "Bambino"}
		cat2 := &Cat{Name: "Banana"}
		cat3 := &Cat{Name: "Bandit"}
		cat4 := &Cat{Name: "Bangle"}
		cat5 := &Cat{Name: "Banjo"}
		cat6 := &Cat{Name: "Banker"}
		cat7 := &Cat{Name: "Banshee"}
		cat8 := &Cat{Name: "Baron"}

		cats := []*Cat{cat1, cat2, cat3, cat4, cat5, cat6, cat7, cat8}
		expected := cat6

		for i := range cats {
			err := sqlxx.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		id := ""
		query := loukoum.Select("id").From("ztp_cat").
			Where(loukoum.Condition("name").Equal("Banker"))

		err := sqlxx.Exec(ctx, driver, query, &id)
		is.NoError(err)
		is.Equal(expected.ID, id)

		err = sqlxx.Exec(ctx, driver, query, id)
		is.Error(err)
		is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))

	})
}

func TestRawExec_Fetch(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{Name: "Cake"}
		cat2 := &Cat{Name: "Calvin"}
		cat3 := &Cat{Name: "Calypso"}
		cat4 := &Cat{Name: "Calzone"}
		cat5 := &Cat{Name: "Cambridge"}
		cat6 := &Cat{Name: "Cameo"}
		cat7 := &Cat{Name: "Campbell"}
		cat8 := &Cat{Name: "Cannes"}

		cats := []*Cat{cat1, cat2, cat3, cat4, cat5, cat6, cat7, cat8}
		expected := cat4

		for i := range cats {
			err := sqlxx.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		id := ""
		query := `SELECT id FROM ztp_cat WHERE name = 'Calzone'`
		err := sqlxx.RawExec(ctx, driver, query, &id)
		is.NoError(err)
		is.Equal(expected.ID, id)

		err = sqlxx.RawExec(ctx, driver, query, id)
		is.Error(err)
		is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))

	})
}

func TestExec_FetchModel(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{Name: "Afro"}
		cat2 := &Cat{Name: "Ajax"}
		cat3 := &Cat{Name: "Akbar"}
		cat4 := &Cat{Name: "Akiko"}

		cats := []*Cat{cat1, cat2, cat3, cat4}
		expected := cat4

		for i := range cats {
			err := sqlxx.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		query := loukoum.Select("*").From("ztp_cat").
			Where(loukoum.Condition("name").Equal("Akiko"))

		{
			result := &Cat{}
			err := sqlxx.Exec(ctx, driver, query, result)
			is.NoError(err)
			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
			is.Equal(expected.CreatedAt, result.CreatedAt)
			is.Equal(expected.UpdatedAt, result.UpdatedAt)
			is.Equal(expected.DeletedAt, result.DeletedAt)
		}

		{
			result := &Cat{}
			err := sqlxx.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
			is.Equal(expected.CreatedAt, result.CreatedAt)
			is.Equal(expected.UpdatedAt, result.UpdatedAt)
			is.Equal(expected.DeletedAt, result.DeletedAt)
		}

		{
			result := Cat{}
			err := sqlxx.Exec(ctx, driver, query, result)
			is.Error(err)
			is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))
		}

	})
}

func TestExec_ListModel(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{Name: "Amazon"}
		cat2 := &Cat{Name: "Amelia"}
		cat3 := &Cat{Name: "Amigo"}
		cat4 := &Cat{Name: "Amos"}

		cats := []*Cat{cat1, cat2, cat3, cat4}

		for i := range cats {
			err := sqlxx.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		query := loukoum.Select("*").From("ztp_cat").
			Where(loukoum.Condition("name").In("Amazon", "Amelia", "Amigo", "Amos"))

		{
			result := []Cat{}
			err := sqlxx.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(result, 4)
			is.Contains(cats, &result[0])
			is.Contains(cats, &result[1])
			is.Contains(cats, &result[2])
			is.Contains(cats, &result[3])
		}

		{
			result := []*Cat{}
			err := sqlxx.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(result, 4)
			is.Contains(cats, result[0])
			is.Contains(cats, result[1])
			is.Contains(cats, result[2])
			is.Contains(cats, result[3])
		}

		{
			result := &[]Cat{}
			err := sqlxx.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(*result, 4)
			is.Contains(cats, &(*result)[0])
			is.Contains(cats, &(*result)[1])
			is.Contains(cats, &(*result)[2])
			is.Contains(cats, &(*result)[3])
		}

		{
			result := &[]*Cat{}
			err := sqlxx.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(*result, 4)
			is.Contains(cats, (*result)[0])
			is.Contains(cats, (*result)[1])
			is.Contains(cats, (*result)[2])
			is.Contains(cats, (*result)[3])
		}

		{
			result := []Cat{}
			err := sqlxx.Exec(ctx, driver, query, result)
			is.Error(err)
			is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))
		}

		{
			result := []*Cat{}
			err := sqlxx.Exec(ctx, driver, query, result)
			is.Error(err)
			is.Equal(sqlxx.ErrPointerRequired, errors.Cause(err))
		}

	})
}
