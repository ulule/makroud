package sqlxx

import (
	"github.com/pkg/errors"
)

// Transaction will creates a transaction.
func Transaction(client *Client, handler func(client *Client) error) error {
	if client == nil {
		return ErrInvalidClient
	}

	tx, err := client.Beginx()
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot create a transaction")
	}

	session := client.copy(tx)

	err = handler(session)
	if err != nil {

		thr := tx.Rollback()
		if thr != nil {
			// TODO: Add an observer to collect this error.
			thr = errors.Wrap(err, "sqlxx: cannot rollback transaction")
			_ = thr
		}

		return err
	}

	err = tx.Commit()
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot commit transaction")
	}

	return nil
}
