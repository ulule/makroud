package makroud

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

// TxOptions is an alias for sql.TxOptions to reduce import leak.
// This alias allows the use of makroud.TxOptions and sql.TxOptions seamlessly.
type TxOptions = sql.TxOptions

// List of supported isolation level for a postgres transaction.
const (
	LevelDefault         = sql.LevelDefault
	LevelReadUncommitted = sql.LevelReadUncommitted
	LevelReadCommitted   = sql.LevelReadCommitted
	LevelRepeatableRead  = sql.LevelRepeatableRead
	LevelSerializable    = sql.LevelSerializable
)

// Transaction will creates a transaction.
func Transaction(ctx context.Context, driver Driver, opts *TxOptions,
	handler func(driver Driver) error) error {

	if driver == nil {
		return errors.Wrap(ErrInvalidDriver, "makroud: cannot create a transaction")
	}

	tx, err := driver.BeginContext(ctx, opts)
	if err != nil {
		return err
	}

	err = handler(tx)
	if err != nil {

		thr := tx.Rollback()
		if thr != nil {
			// TODO (novln): Add an observer to collect this error.
			_ = thr
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
