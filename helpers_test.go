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

		query := loukoum.Select(loukoum.Count("*")).From("ztp_cat").
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

		query = loukoum.Select(loukoum.Count("*")).From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Rod%"))

		{
			count, err := makroud.Count(ctx, driver, query)
			is.NoError(err)
			is.Equal(int64(0), count)
		}
		{
			count, err := makroud.FloatCount(ctx, driver, query)
			is.NoError(err)
			is.Equal(float64(0), count)
		}

		query = loukoum.Select("count").From("cats").
			With(loukoum.With("cats",
				loukoum.Select(loukoum.Count("*")).From("ztp_cat").
					Where(loukoum.Condition("name").ILike("Rod%")),
			)).
			Where(loukoum.Condition("count").GreaterThan(5))

		{
			count, err := makroud.Count(ctx, driver, query)
			is.NoError(err)
			is.Equal(int64(0), count)
		}
		{
			count, err := makroud.FloatCount(ctx, driver, query)
			is.NoError(err)
			is.Equal(float64(0), count)
		}

		query = loukoum.Select("*").From("ztp_cat").
			Where(loukoum.Condition("name").ILike("Rod%"))

		{
			count, err := makroud.Count(ctx, driver, query)
			is.Error(err)
			is.Equal(int64(0), count)
		}
		{
			count, err := makroud.FloatCount(ctx, driver, query)
			is.Error(err)
			is.Equal(float64(0), count)
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

// nolint: gocyclo
func TestExec_ListPartial(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		type PartialMeow struct {
			Hash  string `mk:"hash"`
			Body  string `mk:"body"`
			CatID string `mk:"cat_id"`
		}

		type PartialCat struct {
			Name string `makroud:"name"`
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
			for i := range results {
				is.Contains(expected, results[i])
			}
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
			for i := range results {
				is.Contains(expected, *results[i])
			}
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
			for i := range results {
				is.Contains(expected, results[i].UnixNano())
			}
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
			for i := range results {
				is.Contains(expected, results[i].UnixNano())
			}
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

		{
			stmt := loukoum.Select("hash", "body", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []PartialMeow{
				{
					Hash:  meow1.Hash,
					Body:  meow1.Body,
					CatID: meow1.CatID,
				},
				{
					Hash:  meow2.Hash,
					Body:  meow2.Body,
					CatID: meow2.CatID,
				},
				{
					Hash:  meow3.Hash,
					Body:  meow3.Body,
					CatID: meow3.CatID,
				},
			}

			results := []PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			for i := range results {
				is.Contains(expected, results[i])
			}
		}

		{
			stmt := loukoum.Select("hash", "body", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []*PartialMeow{
				{
					Hash:  meow1.Hash,
					Body:  meow1.Body,
					CatID: meow1.CatID,
				},
				{
					Hash:  meow2.Hash,
					Body:  meow2.Body,
					CatID: meow2.CatID,
				},
				{
					Hash:  meow3.Hash,
					Body:  meow3.Body,
					CatID: meow3.CatID,
				},
			}

			results := []*PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			for i := range results {
				is.Contains(expected, results[i])
			}
		}

		{
			type PartialMeow struct {
				Hash  *string `mk:"hash"`
				Body  *string `mk:"body"`
				CatID *string `mk:"cat_id"`
			}

			stmt := loukoum.Select("hash", "body", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []PartialMeow{
				{
					Hash:  &meow1.Hash,
					Body:  &meow1.Body,
					CatID: &meow1.CatID,
				},
				{
					Hash:  &meow2.Hash,
					Body:  &meow2.Body,
					CatID: &meow2.CatID,
				},
				{
					Hash:  &meow3.Hash,
					Body:  &meow3.Body,
					CatID: &meow3.CatID,
				},
			}

			results := []PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			for i := range results {
				is.Contains(expected, results[i])
			}
		}

		{
			type PartialMeow struct {
				Hash  string `db:"hash"`
				CatID string `db:"cat_id"`
			}

			stmt := loukoum.Select("hash", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []PartialMeow{
				{
					Hash:  meow1.Hash,
					CatID: meow1.CatID,
				},
				{
					Hash:  meow2.Hash,
					CatID: meow2.CatID,
				},
				{
					Hash:  meow3.Hash,
					CatID: meow3.CatID,
				},
			}

			results := []PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			for i := range results {
				is.Contains(expected, results[i])
			}
		}

		{
			type PartialMeow struct {
				Hash  string `mk:"hash"`
				CatID string `mk:"cat_id"`
			}

			stmt := loukoum.Select("hash", "body", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			results := []PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.Error(err)
			is.Equal(makroud.ErrSchemaColumnRequired, errors.Cause(err))
		}

		{
			type PartialMeow struct {
				Hash      string    `makroud:"hash"`
				Body      string    `makroud:"body"`
				CreatedAt time.Time `makroud:"created"`
			}

			stmt := loukoum.Select("hash", "created", "body").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []PartialMeow{
				{
					Hash:      meow1.Hash,
					CreatedAt: meow1.CreatedAt,
					Body:      meow1.Body,
				},
				{
					Hash:      meow2.Hash,
					CreatedAt: meow2.CreatedAt,
					Body:      meow2.Body,
				},
				{
					Hash:      meow3.Hash,
					CreatedAt: meow3.CreatedAt,
					Body:      meow3.Body,
				},
			}

			results := []PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			for i := range results {
				is.Contains(expected, results[i])
			}
		}

		{
			type PartialMeow struct {
				Hash      *string    `makroud:"hash"`
				Body      *string    `makroud:"body"`
				CreatedAt *time.Time `makroud:"created"`
			}

			stmt := loukoum.Select("hash", "created", "body").
				From("ztp_meow").
				Where(loukoum.Condition("cat_id").Equal(cat.ID))

			expected := []PartialMeow{
				{
					Hash:      &meow1.Hash,
					CreatedAt: &meow1.CreatedAt,
					Body:      &meow1.Body,
				},
				{
					Hash:      &meow2.Hash,
					CreatedAt: &meow2.CreatedAt,
					Body:      &meow2.Body,
				},
				{
					Hash:      &meow3.Hash,
					CreatedAt: &meow3.CreatedAt,
					Body:      &meow3.Body,
				},
			}

			results := []PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(len(expected), len(results))
			for i := range results {
				is.Contains(expected, results[i])
			}
		}

		{
			stmt := loukoum.Select("name").
				From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))

			expected := []PartialCat{
				{
					Name: cat.Name,
				},
			}

			results := []PartialCat{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(expected, results)
		}

		{
			stmt := loukoum.Select("name").
				From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))

			expected := []*PartialCat{
				{
					Name: cat.Name,
				},
			}

			results := []*PartialCat{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(expected, results)
		}

		{
			type PartialCat struct {
				Name string
			}

			stmt := loukoum.Select("name").
				From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))

			expected := []PartialCat{
				{
					Name: cat.Name,
				},
			}

			results := []PartialCat{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(expected, results)
		}

		{
			type PartialCat struct {
				Name string
			}

			stmt := loukoum.Select("name").
				From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))

			expected := []*PartialCat{
				{
					Name: cat.Name,
				},
			}

			results := []*PartialCat{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.NoError(err)
			is.Equal(expected, results)
		}
	})
}

// nolint: gocyclo
func TestExec_FetchPartial(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		type PartialMeow struct {
			Hash  string `mk:"hash"`
			Body  string `mk:"body"`
			CatID string `mk:"cat_id"`
		}

		type PartialCat struct {
			Name string `makroud:"name"`
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
			stmt := loukoum.Select("body").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow1.Hash))

			expected := meow1.Body
			var result string

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{
			stmt := loukoum.Select("cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow2.Hash))

			expected := &meow2.CatID
			var result *string

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.NotNil(result)
			is.Equal(*expected, *result)
		}

		{
			stmt := loukoum.Select("created").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow3.Hash))

			expected := meow3.CreatedAt
			var result time.Time

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)

		}

		{
			stmt := loukoum.Select("created").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow1.Hash))

			expected := &meow1.CreatedAt
			var result *time.Time

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected.UnixNano(), result.UnixNano())
		}

		{
			stmt := loukoum.Select("hash", "body").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow1.Hash))

			result := ""

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.Error(err)
			is.Equal(makroud.ErrSliceOfScalarMultipleColumns, errors.Cause(err))
		}

		{
			stmt := loukoum.Select("hash", "body", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow1.Hash))

			expected := PartialMeow{
				Hash:  meow1.Hash,
				Body:  meow1.Body,
				CatID: meow1.CatID,
			}

			result := PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{
			stmt := loukoum.Select("hash", "body", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow2.Hash))

			expected := &PartialMeow{
				Hash:  meow2.Hash,
				Body:  meow2.Body,
				CatID: meow2.CatID,
			}

			result := &PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{
			type PartialMeow struct {
				Hash  *string `mk:"hash"`
				Body  *string `mk:"body"`
				CatID *string `mk:"cat_id"`
			}

			stmt := loukoum.Select("hash", "body", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow3.Hash))

			expected := PartialMeow{
				Hash:  &meow3.Hash,
				Body:  &meow3.Body,
				CatID: &meow3.CatID,
			}

			result := PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{
			type PartialMeow struct {
				Hash  string `db:"hash"`
				CatID string `db:"cat_id"`
			}

			stmt := loukoum.Select("hash", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow1.Hash))

			expected := PartialMeow{
				Hash:  meow1.Hash,
				CatID: meow1.CatID,
			}

			result := PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{
			type PartialMeow struct {
				Hash  string `mk:"hash"`
				CatID string `mk:"cat_id"`
			}

			stmt := loukoum.Select("hash", "body", "cat_id").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow2.Hash))

			results := PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &results)
			is.Error(err)
			is.Equal(makroud.ErrSchemaColumnRequired, errors.Cause(err))
		}

		{
			type PartialMeow struct {
				Hash      string    `makroud:"hash"`
				Body      string    `makroud:"body"`
				CreatedAt time.Time `makroud:"created"`
			}

			stmt := loukoum.Select("hash", "created", "body").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow3.Hash))

			expected := PartialMeow{
				Hash:      meow3.Hash,
				CreatedAt: meow3.CreatedAt,
				Body:      meow3.Body,
			}

			result := PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{
			type PartialMeow struct {
				Hash      *string    `makroud:"hash"`
				Body      *string    `makroud:"body"`
				CreatedAt *time.Time `makroud:"created"`
			}

			stmt := loukoum.Select("hash", "created", "body").
				From("ztp_meow").
				Where(loukoum.Condition("hash").Equal(meow3.Hash))

			expected := &PartialMeow{
				Hash:      &meow3.Hash,
				CreatedAt: &meow3.CreatedAt,
				Body:      &meow3.Body,
			}

			result := &PartialMeow{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{
			stmt := loukoum.Select("name").
				From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))

			expected := PartialCat{
				Name: cat.Name,
			}

			result := PartialCat{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{
			stmt := loukoum.Select("name").
				From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))

			expected := &PartialCat{
				Name: cat.Name,
			}

			result := &PartialCat{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{

			type PartialCat struct {
				Name string
			}

			stmt := loukoum.Select("name").
				From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))

			expected := PartialCat{
				Name: cat.Name,
			}

			result := PartialCat{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}

		{
			type PartialCat struct {
				Name string
			}

			stmt := loukoum.Select("name").
				From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))

			expected := &PartialCat{
				Name: cat.Name,
			}

			result := &PartialCat{}

			err := makroud.Exec(ctx, driver, stmt, &result)
			is.NoError(err)
			is.Equal(expected, result)
		}
	})
}
