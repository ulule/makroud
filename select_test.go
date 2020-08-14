package makroud_test

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum/v3"

	"github.com/ulule/makroud"
)

func TestSelect_Row(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		humans := []Human{
			{Name: "Ethel"},
			{Name: "Marcela"},
			{Name: "Juliette"},
			{Name: "Sarah"},
			{Name: "Valens"},
			{Name: "Marin"},
		}

		for i := range humans {
			err := makroud.Save(ctx, driver, &humans[i])
			is.NoError(err)
		}

		sort.Slice(humans, func(i, j int) bool {
			return humans[i].ID < humans[j].ID
		})

		{
			expected := humans[0]

			result := &Human{}
			err := makroud.Select(ctx, driver, result)
			is.NoError(err)

			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
		}
		{
			expected := humans[5]

			result := &Human{}
			err := makroud.Select(ctx, driver, result, loukoum.Order("id", loukoum.Desc))
			is.NoError(err)

			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
		}
		{
			expected := humans[3]

			result := &Human{}
			err := makroud.Select(ctx, driver, result,
				loukoum.Order("id", loukoum.Desc), loukoum.Offset(2))
			is.NoError(err)

			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
		}
		{
			expected := humans[2]

			result := &Human{}
			err := makroud.Select(ctx, driver, result,
				loukoum.Condition("name").Like(fmt.Sprint(expected.Name[0:2], "%")),
			)
			is.NoError(err)

			is.Equal(expected.ID, result.ID)
			is.Equal(expected.Name, result.Name)
		}
	})
}

func TestSelect_Rows(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		humans := []Human{
			{Name: "Erwann"},
			{Name: "Aymeric"},
			{Name: "Enora"},
			{Name: "Loïc"},
			{Name: "Armelle"},
			{Name: "Yüna"},
		}

		for i := range humans {
			err := makroud.Save(ctx, driver, &humans[i])
			is.NoError(err)
		}

		sort.Slice(humans, func(i, j int) bool {
			return humans[i].ID < humans[j].ID
		})

		{
			result := &[]Human{}
			err := makroud.Select(ctx, driver, result)
			is.NoError(err)
			is.Len(*result, len(humans))

			for i := range *result {
				expected := humans[i]

				is.Equal(expected.ID, (*result)[i].ID)
				is.Equal(expected.Name, (*result)[i].Name)
			}
		}
		{
			result := &[]Human{}
			err := makroud.Select(ctx, driver, result, loukoum.Order("id", loukoum.Desc))
			is.NoError(err)
			is.Len(*result, len(humans))

			for i := range *result {
				idx := len(humans) - (1 + i)
				expected := humans[idx]

				is.Equal(expected.ID, (*result)[i].ID)
				is.Equal(expected.Name, (*result)[i].Name)
			}
		}
		{

			result := &[]Human{}
			err := makroud.Select(ctx, driver, result,
				loukoum.Order("id", loukoum.Desc), loukoum.Offset(2))
			is.NoError(err)
			is.Len(*result, len(humans)-2)

			for i := range *result {
				idx := len(humans) - (3 + i)
				expected := humans[idx]

				is.Equal(expected.ID, (*result)[i].ID)
				is.Equal(expected.Name, (*result)[i].Name)
			}
		}
		{

			result := &[]Human{}
			err := makroud.Select(ctx, driver, result,
				loukoum.Order("id", loukoum.Desc), loukoum.Offset(2), loukoum.Limit(2))
			is.NoError(err)
			is.Len(*result, 2)

			for i := range *result {
				idx := len(humans) - (3 + i)
				expected := humans[idx]

				is.Equal(expected.ID, (*result)[i].ID)
				is.Equal(expected.Name, (*result)[i].Name)
			}
		}
		{
			expected := humans[2]

			result := &[]Human{}
			err := makroud.Select(ctx, driver, result,
				loukoum.Condition("name").Like(fmt.Sprint(expected.Name[0:2], "%")),
			)
			is.NoError(err)
			is.Len(*result, 1)

			is.Equal(expected.ID, (*result)[0].ID)
			is.Equal(expected.Name, (*result)[0].Name)
		}
	})
}
