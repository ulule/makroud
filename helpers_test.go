package makroud_test

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum/v3"

	"github.com/ulule/makroud"
)

import "fmt"

func TestCount(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cats := []Cat{
			{Name: "Radio"},
			{Name: "Radish"},
			{Name: "Radium"},
			{Name: "Radix"},
			{Name: "Radman"},
			{Name: "Radmilla"},
		}

		for i := range cats {
			err := makroud.Save(ctx, driver, &cats[i])
			is.NoError(err)
		}

		query := loukoum.Select("COUNT(*)").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Rad%"))

		{
			count, err := makroud.Count(ctx, driver, query)
			is.NoError(err)
			is.Equal(int64(6), count)
		}
		{
			count, err := makroud.FloatCount(ctx, driver, query)
			is.NoError(err)
			is.Equal(float64(6), count)
		}

	})
}

func TestExec_List(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
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
			err := makroud.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		list := []string{}
		query := loukoum.Select("id").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Whi%"))

		err := makroud.Exec(ctx, driver, query, &list)
		is.NoError(err)

		is.Len(list, len(expected))
		for i := range expected {
			is.Contains(list, expected[i].ID)
		}

		list = []string{}
		query = loukoum.Select("id").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("XXXXX%"))

		err = makroud.Exec(ctx, driver, query, &list)
		is.NoError(err)
		is.NotNil(list)
		is.Empty(list)

		err = makroud.Exec(ctx, driver, query, []string{})
		is.Error(err)
		is.Equal(makroud.ErrPointerRequired, errors.Cause(err))

		expected = []*Cat{cat3, cat6}
		list = []string{}
		query = loukoum.Select("id").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Whi%")).
			And(loukoum.Condition("name").ILike("%ey"))

		err = makroud.Exec(ctx, driver, query, &list)
		is.NoError(err)

		is.Len(list, len(expected))
		for i := range expected {
			is.Contains(list, expected[i].ID)
		}

	})
}

func TestRawExec_List(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
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
			err := makroud.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		list := []string{}
		query := `SELECT id FROM ztp_cat WHERE name ILIKE 'Ver%'`
		err := makroud.RawExec(ctx, driver, query, &list)
		is.NoError(err)

		is.Len(list, len(expected))
		for i := range expected {
			is.Contains(list, expected[i].ID)
		}

		list = []string{}
		query = `SELECT id FROM ztp_cat WHERE name ILIKE 'XXXXX%'`
		err = makroud.RawExec(ctx, driver, query, &list)
		is.NoError(err)
		is.NotNil(list)
		is.Empty(list)

		err = makroud.RawExec(ctx, driver, query, []string{})
		is.Error(err)
		is.Equal(makroud.ErrPointerRequired, errors.Cause(err))

		expected = []*Cat{cat5, cat8}
		list = []string{}
		query = `SELECT id FROM ztp_cat WHERE name ILIKE 'Ve%' AND name ILIKE '%s'`
		err = makroud.RawExec(ctx, driver, query, &list)
		is.NoError(err)

		is.Len(list, len(expected))
		for i := range expected {
			is.Contains(list, expected[i].ID)
		}

	})
}

func TestExec_Fetch(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
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
			err := makroud.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		id := ""
		query := loukoum.Select("id").From("ztp_cat").
			Where(loukoum.Condition("name").Equal("Banker"))

		err := makroud.Exec(ctx, driver, query, &id)
		is.NoError(err)
		is.Equal(expected.ID, id)

		err = makroud.Exec(ctx, driver, query, id)
		is.Error(err)
		is.Equal(makroud.ErrPointerRequired, errors.Cause(err))

		id = ""
		query = loukoum.Select("id").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Ba%")).
			And(loukoum.Condition("name").ILike("%er"))

		err = makroud.Exec(ctx, driver, query, &id)
		is.NoError(err)
		is.Equal(expected.ID, id)

	})
}

func TestRawExec_Fetch(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
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
			err := makroud.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		id := ""
		query := `SELECT id FROM ztp_cat WHERE name = 'Calzone'`
		err := makroud.RawExec(ctx, driver, query, &id)
		is.NoError(err)
		is.Equal(expected.ID, id)

		err = makroud.RawExec(ctx, driver, query, id)
		is.Error(err)
		is.Equal(makroud.ErrPointerRequired, errors.Cause(err))

		id = ""
		query = `SELECT id FROM ztp_cat WHERE name ILIKE 'Cal%' AND name ILIKE '%ne'`
		err = makroud.RawExec(ctx, driver, query, &id)
		is.NoError(err)
		is.Equal(expected.ID, id)

	})
}

func TestRawExecArgs_Fetch(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{Name: "Ezekiel"}
		cat2 := &Cat{Name: "Edouard"}
		cat3 := &Cat{Name: "Eric"}
		cat4 := &Cat{Name: "Elliot"}

		cats := []*Cat{cat1, cat2, cat3, cat4}
		expected := cat3

		for i := range cats {
			err := makroud.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		id := ""
		query := `SELECT id FROM ztp_cat WHERE name = $1`
		err := makroud.RawExecArgs(ctx, driver, query, []interface{}{"Eric"}, &id)
		is.NoError(err)
		is.Equal(expected.ID, id)
	})
}

func TestExec_FetchModel(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{Name: "Afro"}
		cat2 := &Cat{Name: "Ajax"}
		cat3 := &Cat{Name: "Akbar"}
		cat4 := &Cat{Name: "Akiko"}

		cats := []*Cat{cat1, cat2, cat3, cat4}
		expected := cat4

		for i := range cats {
			err := makroud.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		query := loukoum.Select("*").From("ztp_cat").
			Where(loukoum.Condition("name").Equal("Akiko"))

		{
			result := &Cat{}
			err := makroud.Exec(ctx, driver, query, result)
			is.NoError(err)
			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
			is.Equal(expected.CreatedAt, result.CreatedAt)
			is.Equal(expected.UpdatedAt, result.UpdatedAt)
			is.Equal(expected.DeletedAt, result.DeletedAt)
		}

		{
			result := &Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
			is.Equal(expected.CreatedAt, result.CreatedAt)
			is.Equal(expected.UpdatedAt, result.UpdatedAt)
			is.Equal(expected.DeletedAt, result.DeletedAt)
		}

		{
			result := Cat{}
			err := makroud.Exec(ctx, driver, query, result)
			is.Error(err)
			is.Equal(makroud.ErrPointerRequired, errors.Cause(err))
		}

		query = loukoum.Select("*").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Ak%")).
			And(loukoum.Condition("name").ILike("%ko"))

		{
			result := &Cat{}
			err := makroud.Exec(ctx, driver, query, result)
			is.NoError(err)
			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
			is.Equal(expected.CreatedAt, result.CreatedAt)
			is.Equal(expected.UpdatedAt, result.UpdatedAt)
			is.Equal(expected.DeletedAt, result.DeletedAt)
		}

		{
			result := &Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
			is.Equal(expected.CreatedAt, result.CreatedAt)
			is.Equal(expected.UpdatedAt, result.UpdatedAt)
			is.Equal(expected.DeletedAt, result.DeletedAt)
		}

	})
}

func TestExec_ListModel(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat1 := &Cat{Name: "Amazon"}
		cat2 := &Cat{Name: "Amelia"}
		cat3 := &Cat{Name: "Amigo"}
		cat4 := &Cat{Name: "Amos"}

		cats := []*Cat{cat1, cat2, cat3, cat4}

		for i := range cats {
			err := makroud.Save(ctx, driver, cats[i])
			is.NoError(err)
		}

		query := loukoum.Select("*").From("ztp_cat").
			Where(loukoum.Condition("name").In("Amazon", "Amelia", "Amigo", "Amos"))

		{
			result := []Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(result, 4)
			is.Contains(cats, &result[0])
			is.Contains(cats, &result[1])
			is.Contains(cats, &result[2])
			is.Contains(cats, &result[3])
		}

		{
			result := []*Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(result, 4)
			is.Contains(cats, result[0])
			is.Contains(cats, result[1])
			is.Contains(cats, result[2])
			is.Contains(cats, result[3])
		}

		{
			result := &[]Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(*result, 4)
			is.Contains(cats, &(*result)[0])
			is.Contains(cats, &(*result)[1])
			is.Contains(cats, &(*result)[2])
			is.Contains(cats, &(*result)[3])
		}

		{
			result := &[]*Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(*result, 4)
			is.Contains(cats, (*result)[0])
			is.Contains(cats, (*result)[1])
			is.Contains(cats, (*result)[2])
			is.Contains(cats, (*result)[3])
		}

		{
			result := []Cat{}
			err := makroud.Exec(ctx, driver, query, result)
			is.Error(err)
			is.Equal(makroud.ErrPointerRequired, errors.Cause(err))
		}

		{
			result := []*Cat{}
			err := makroud.Exec(ctx, driver, query, result)
			is.Error(err)
			is.Equal(makroud.ErrPointerRequired, errors.Cause(err))
		}

		query = loukoum.Select("*").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Am%"))

		{
			result := []Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(result, 4)
			is.Contains(cats, &result[0])
			is.Contains(cats, &result[1])
			is.Contains(cats, &result[2])
			is.Contains(cats, &result[3])
		}

		{
			result := []*Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(result, 4)
			is.Contains(cats, result[0])
			is.Contains(cats, result[1])
			is.Contains(cats, result[2])
			is.Contains(cats, result[3])
		}

		{
			result := &[]Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(*result, 4)
			is.Contains(cats, &(*result)[0])
			is.Contains(cats, &(*result)[1])
			is.Contains(cats, &(*result)[2])
			is.Contains(cats, &(*result)[3])
		}

		{
			result := &[]*Cat{}
			err := makroud.Exec(ctx, driver, query, &result)
			is.NoError(err)
			is.Len(*result, 4)
			is.Contains(cats, (*result)[0])
			is.Contains(cats, (*result)[1])
			is.Contains(cats, (*result)[2])
			is.Contains(cats, (*result)[3])
		}

	})
}

func TestJoin_One(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		ctx := context.Background()

		GenerateZootopiaFixtures(ctx, driver, is)

		stmt := loukoum.Select(`ztp_cat.id "ztp_cat.id"`, `ztp_meow.*`, `ztp_cat.name "ztp_cat.name"`).
			From("ztp_meow").
			Join("ztp_cat", "ON ztp_meow.cat_id = ztp_cat.id", loukoum.InnerJoin).Limit(1)

		result := &Meow{}
		err := makroud.Exec(ctx, driver, stmt, result)
		is.NoError(err)
		is.NotNil(result.Cat)
		is.Equal(result.CatID, result.Cat.ID)
		is.NotZero(result.Cat.Name)
		is.Nil(result.Cat.Feeder)
	})
}

func TestJoin_OneMany(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		ctx := context.Background()

		GenerateZootopiaFixtures(ctx, driver, is)

		stmt := loukoum.Select(
			`ztp_cat.id "ztp_cat.id"`, `ztp_human.cat_id "ztp_cat.ztp_human.cat_id"`, `ztp_meow.*`,
			`ztp_cat.name "ztp_cat.name"`, `ztp_human.id "ztp_cat.ztp_human.id"`, `ztp_human.name "ztp_cat.ztp_human.name"`,
		).
			From("ztp_meow").
			Join("ztp_cat", "ON ztp_meow.cat_id = ztp_cat.id", loukoum.InnerJoin).
			Join("ztp_human", "ON ztp_human.cat_id = ztp_cat.id", loukoum.InnerJoin).
			Limit(1)

		result := &Meow{}
		err := makroud.Exec(ctx, driver, stmt, result)
		is.NoError(err)
		is.NotNil(result.Cat)
		is.Equal(result.CatID, result.Cat.ID)
		is.NotZero(result.Cat.Name)
		is.NotNil(result.Cat.Feeder)
		is.Equal(result.Cat.ID, result.Cat.Feeder.CatID.String)
	})
}

func TestJoin_Many(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)
		ctx := context.Background()

		GenerateZootopiaFixtures(ctx, driver, is)

		stmt := loukoum.Select(`ztp_cat.id "ztp_cat.id"`, `ztp_meow.*`, `ztp_cat.name "ztp_cat.name"`).
			From("ztp_meow").
			Join("ztp_cat", "ON ztp_meow.cat_id = ztp_cat.id", loukoum.InnerJoin)

		results := []Meow{}
		err := makroud.Exec(ctx, driver, stmt, &results)
		is.NoError(err)

		for i := range results {
			is.NotNil(results[i].Cat)
			is.Equal(results[i].CatID, results[i].Cat.ID)
			is.NotZero(results[i].Cat.Name)
		}
	})
}

func TestExec_ListPartial(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		type PartialMeow struct {
			Hash  string `db:"hash"`
			Body  string `db:"body"`
			CatID string `db:"cat_id"`
		}

		fixtures := GenerateZootopiaFixtures(ctx, driver, is)
		cat := fixtures.Cats[3]
		meow1 := fixtures.Meows[4]
		meow2 := fixtures.Meows[5]
		meow3 := fixtures.Meows[6]

		is.Equal(cat.ID, meow1.CatID)
		is.Equal(cat.ID, meow2.CatID)
		is.Equal(cat.ID, meow3.CatID)

		{
			stmt := loukoum.Select("hash").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []string{meow1.Hash, meow2.Hash, meow3.Hash}
			results := []string{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			is.Contains(expected, results[0])
			is.Contains(expected, results[1])
			is.Contains(expected, results[2])
		}

		{
			stmt := loukoum.Select("hash").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []string{meow1.Hash, meow2.Hash, meow3.Hash}
			results := []*string{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			is.Contains(expected, *results[0])
			is.Contains(expected, *results[1])
			is.Contains(expected, *results[2])
		}

		{
			stmt := loukoum.Select("created").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []int64{
				meow1.CreatedAt.UnixNano(),
				meow2.CreatedAt.UnixNano(),
				meow3.CreatedAt.UnixNano(),
			}
			results := []time.Time{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			is.Contains(expected, results[0].UnixNano())
			is.Contains(expected, results[1].UnixNano())
			is.Contains(expected, results[2].UnixNano())
		}

		{
			stmt := loukoum.Select("created").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []int64{
				meow1.CreatedAt.UnixNano(),
				meow2.CreatedAt.UnixNano(),
				meow3.CreatedAt.UnixNano(),
			}
			results := []*time.Time{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			is.Contains(expected, results[0].UnixNano())
			is.Contains(expected, results[1].UnixNano())
			is.Contains(expected, results[2].UnixNano())
		}

		{
			stmt := loukoum.Select("hash", "body").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			results := []string{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.Error(err)
			is.Equal(makroud.ErrSliceOfScalarMultipleColumns, errors.Cause(err))
		}

		// SANDBOX HERE

		results := []time.Time{}

		fmt.Printf("::3 %d / %d\n", len(results), cap(results))

		stmt := loukoum.Select("created").
			From("ztp_meow").
			Where(loukoum.Condition("cat_id").Equal(cat.ID))

		err := makroud.Exec(ctx, driver, stmt, &results)
		is.NoError(err)

		for i := range results {
			fmt.Printf("%+v\n", results[i])
		}

		fmt.Printf("::4 %d / %d\n", len(results), cap(results))

	})
}
