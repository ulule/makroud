package sqlxx_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum"

	"github.com/ulule/sqlxx"
)

func TestTransaction_Commit(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{Name: "Harlay"}

		err := sqlxx.Save(ctx, driver, cat)
		is.NoError(err)

		err = sqlxx.Transaction(driver, func(driver sqlxx.Driver) error {
			cat.Name = "Harley"
			err := sqlxx.Save(ctx, driver, cat)
			is.NoError(err)
			return nil
		})
		is.NoError(err)

		name := ""
		query := loukoum.Select("name").From("ztp_cat").Where(loukoum.Condition("id").Equal(cat.ID))
		err = sqlxx.Exec(ctx, driver, query, &name)
		is.NoError(err)
		is.Equal("Harley", name)

	})
}

func TestTransaction_Rollback(t *testing.T) {
	Setup(t)(func(driver sqlxx.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{Name: "Gemmz"}
		timeout := errors.New("tcp: read timeout on 10.0.3.11:7000")

		err := sqlxx.Save(ctx, driver, cat)
		is.NoError(err)

		err = sqlxx.Transaction(driver, func(driver sqlxx.Driver) error {
			cat.Name = "Gemma"
			err := sqlxx.Save(ctx, driver, cat)
			is.NoError(err)
			return timeout
		})
		is.Error(err)
		is.Equal(timeout, err)

		name := ""
		query := loukoum.Select("name").From("ztp_cat").Where(loukoum.Condition("id").Equal(cat.ID))
		err = sqlxx.Exec(ctx, driver, query, &name)
		is.NoError(err)
		is.Equal("Gemmz", name)

	})
}
