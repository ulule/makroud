package sqlxx

import (
	"github.com/pkg/errors"
)

// Transaction will creates a transaction.
func Transaction(driver Driver, handler func(driver Driver) error) error {
	if driver == nil {
		return errors.Wrap(ErrInvalidDriver, "sqlxx: cannot create a transaction")
	}

	tx, err := driver.Begin()
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
