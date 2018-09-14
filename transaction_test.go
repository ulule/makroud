package makroud_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"

	"github.com/ulule/makroud"
)

func TestTransaction_Commit(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{Name: "Harlay"}

		err := makroud.Save(ctx, driver, cat)
		is.NoError(err)

		err = makroud.Transaction(driver, func(driver makroud.Driver) error {
			cat.Name = "Harley"
			err := makroud.Save(ctx, driver, cat)
			is.NoError(err)
			return nil
		})
		is.NoError(err)

		name := ""
		query := loukoum.Select("name").From("ztp_cat").Where(loukoum.Condition("id").Equal(cat.ID))
		err = makroud.Exec(ctx, driver, query, &name)
		is.NoError(err)
		is.Equal("Harley", name)

	})
}

func TestTransaction_Rollback(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{Name: "Gemmz"}
		timeout := errors.New("tcp: read timeout on 10.0.3.11:7000")

		err := makroud.Save(ctx, driver, cat)
		is.NoError(err)

		err = makroud.Transaction(driver, func(driver makroud.Driver) error {
			cat.Name = "Gemma"
			err := makroud.Save(ctx, driver, cat)
			is.NoError(err)
			return timeout
		})
		is.Error(err)
		is.Equal(timeout, err)

		name := ""
		query := loukoum.Select("name").From("ztp_cat").Where(loukoum.Condition("id").Equal(cat.ID))
		err = makroud.Exec(ctx, driver, query, &name)
		is.NoError(err)
		is.Equal("Gemmz", name)

	})
}
