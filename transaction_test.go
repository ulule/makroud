package makroud_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/ulule/loukoum/v3"

	"github.com/ulule/makroud"
)

func TestTransaction_Commit(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{Name: "Harlay"}

		err := makroud.Save(ctx, driver, cat)
		is.NoError(err)

		err = makroud.Transaction(ctx, driver, nil, func(tx makroud.Driver) error {
			cat.Name = "Harley"
			err := makroud.Save(ctx, tx, cat)
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

		err = makroud.Transaction(ctx, driver, nil, func(driver makroud.Driver) error {
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

func TestTransaction_ErrInvalidDriver(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		err := makroud.Transaction(ctx, nil, nil, func(tx makroud.Driver) error {
			return nil
		})
		is.Error(err)
		is.Equal(makroud.ErrInvalidDriver, errors.Cause(err))
	})
}

func TestTransaction_Nested(t *testing.T) {
	Setup(t, makroud.EnableSavepoint())(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{Name: "Sybil"}
		timeout := errors.New("tcp: read timeout on 10.0.3.11:7000")

		err := makroud.Save(ctx, driver, cat)
		is.NoError(err)

		getCatName := func(driver makroud.Driver) string {
			name := ""
			query := loukoum.Select("name").From("ztp_cat").
				Where(loukoum.Condition("id").Equal(cat.ID))
			err := makroud.Exec(ctx, driver, query, &name)
			is.NoError(err)
			return name
		}

		setCatName := func(driver makroud.Driver, name string) {
			query := loukoum.Update("ztp_cat").Set(loukoum.Map{"name": name}).
				Where(loukoum.Condition("id").Equal(cat.ID))
			err := makroud.Exec(ctx, driver, query)
			is.NoError(err)
		}

		// First transaction.
		err = makroud.Transaction(ctx, driver, nil, func(tx1 makroud.Driver) error {

			setCatName(tx1, "Sibyl")

			// Second transaction.
			err = makroud.Transaction(ctx, tx1, nil, func(tx2 makroud.Driver) error {

				is.Equal("Sibyl", getCatName(tx2))
				setCatName(tx2, "Sibil")

				// Third transaction.
				err = makroud.Transaction(ctx, tx2, nil, func(tx3 makroud.Driver) error {
					is.Equal("Sibil", getCatName(tx3))
					setCatName(tx3, "Sibilll")
					is.Equal("Sibilll", getCatName(tx3))
					return nil
				})
				is.NoError(err)
				is.Equal("Sibilll", getCatName(tx2))

				return timeout
			})
			is.Error(err)
			is.Equal(timeout, err)
			is.Equal("Sibyl", getCatName(tx1))

			return nil
		})
		is.NoError(err)
		is.Equal("Sibyl", getCatName(driver))
	})
}

func TestTransaction_IsolationLevel(t *testing.T) {
	Setup(t, makroud.EnableSavepoint())(func(driver makroud.Driver) {
		ctx := context.Background()
		is := require.New(t)

		cat := &Cat{Name: "Harlay"}

		err := makroud.Save(ctx, driver, cat)
		is.NoError(err)

		handler := func(driver makroud.Driver) error {
			cat.Name = "Harley"
			err := makroud.Save(ctx, driver, cat)
			return err
		}

		{
			opts := &makroud.TxOptions{
				Isolation: sql.LevelDefault,
			}
			err = makroud.Transaction(ctx, driver, opts, handler)
			is.NoError(err)
		}
		{
			opts := &makroud.TxOptions{
				Isolation: sql.LevelReadUncommitted,
			}
			err = makroud.Transaction(ctx, driver, opts, handler)
			is.NoError(err)
		}
		{
			opts := &makroud.TxOptions{
				Isolation: sql.LevelReadCommitted,
			}
			err = makroud.Transaction(ctx, driver, opts, handler)
			is.NoError(err)
		}
		{
			opts := &makroud.TxOptions{
				Isolation: sql.LevelWriteCommitted,
			}
			err = makroud.Transaction(ctx, driver, opts, handler)
			is.Error(err)
		}
		{
			opts := &makroud.TxOptions{
				Isolation: sql.LevelRepeatableRead,
			}
			err = makroud.Transaction(ctx, driver, opts, handler)
			is.NoError(err)
		}
		{
			opts := &makroud.TxOptions{
				Isolation: sql.LevelSnapshot,
			}
			err = makroud.Transaction(ctx, driver, opts, handler)
			is.Error(err)
		}
		{
			opts := &makroud.TxOptions{
				Isolation: sql.LevelSerializable,
			}
			err = makroud.Transaction(ctx, driver, opts, handler)
			is.NoError(err)
		}
		{
			opts := &makroud.TxOptions{
				Isolation: sql.LevelLinearizable,
			}
			err = makroud.Transaction(ctx, driver, opts, handler)
			is.Error(err)
		}
	})
}
