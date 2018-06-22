package sqlxx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"

	"github.com/ulule/sqlxx"
)

func TestExec_List(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
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
			err := sqlxx.Save(driver, cats[i])
			is.NoError(err)
		}

		list := []string{}
		query := loukoum.Select("id").From("wp_cat").
			Where(loukoum.Condition("name").ILike("Whi%"))

		err := sqlxx.Exec(driver, query, &list)
		is.NoError(err)

		is.Len(list, len(expected))
		for i := range expected {
			is.Contains(list, expected[i].ID)
		}

	})
}

func TestRawExec_List(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
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
			err := sqlxx.Save(driver, cats[i])
			is.NoError(err)
		}

		list := []string{}
		query := `SELECT id FROM wp_cat WHERE name ILIKE 'Ver%'`
		err := sqlxx.RawExec(driver, query, &list)
		is.NoError(err)

		is.Len(list, len(expected))
		for i := range expected {
			is.Contains(list, expected[i].ID)
		}

	})
}

func TestExec_Fetch(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
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
			err := sqlxx.Save(driver, cats[i])
			is.NoError(err)
		}

		id := ""
		query := loukoum.Select("id").From("wp_cat").
			Where(loukoum.Condition("name").Equal("Banker"))

		err := sqlxx.Exec(driver, query, &id)
		is.NoError(err)

		is.Equal(expected.ID, id)

	})
}

func TestRawExec_Fetch(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
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
			err := sqlxx.Save(driver, cats[i])
			is.NoError(err)
		}

		id := ""
		query := `SELECT id FROM wp_cat WHERE name = 'Calzone'`
		err := sqlxx.RawExec(driver, query, &id)
		is.NoError(err)

		is.Equal(expected.ID, id)

	})
}

func TestExec_FetchModel(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)

		cat1 := &Cat{Name: "Afro"}
		cat2 := &Cat{Name: "Ajax"}
		cat3 := &Cat{Name: "Akbar"}
		cat4 := &Cat{Name: "Akiko"}

		cats := []*Cat{cat1, cat2, cat3, cat4}
		expected := cat4

		for i := range cats {
			err := sqlxx.Save(driver, cats[i])
			is.NoError(err)
		}

		result := &Cat{}
		query := loukoum.Select("*").From("wp_cat").
			Where(loukoum.Condition("name").Equal("Akiko"))

		err := sqlxx.Exec(driver, query, result)
		is.NoError(err)

		is.Equal(expected.ID, result.ID)
		is.Equal(expected.Name, result.Name)
		is.Equal(expected.CreatedAt, result.CreatedAt)
		is.Equal(expected.UpdatedAt, result.UpdatedAt)
		is.Equal(expected.DeletedAt, result.DeletedAt)

	})
}

func TestExec_ListModel(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		is := require.New(t)

		cat1 := &Cat{Name: "Amazon"}
		cat2 := &Cat{Name: "Amelia"}
		cat3 := &Cat{Name: "Amigo"}
		cat4 := &Cat{Name: "Amos"}

		cats := []*Cat{cat1, cat2, cat3, cat4}

		for i := range cats {
			err := sqlxx.Save(driver, cats[i])
			is.NoError(err)
		}

		result := []Cat{}
		query := loukoum.Select("*").From("wp_cat").
			Where(loukoum.Condition("name").In("Amazon", "Amelia", "Amigo", "Amos"))

		err := sqlxx.Exec(driver, query, &result)
		is.NoError(err)
		is.Len(result, 4)

		for i := range result {
			is.Contains(cats, &(result[i]))
		}

	})
}
