package makroud

import (
	"github.com/pkg/errors"
)

// ----------------------------------------------------------------------------
// Mapper
// ----------------------------------------------------------------------------

// Mapper will be used to mutate a Model with row values.
type Mapper map[string]interface{}

// ScanRow will scan given sqlx.Row to created its Mapper.
func ScanRow(row Row) (Mapper, error) {
	mapper := map[string]interface{}{}
	err := row.Write(mapper)
	if len(mapper) == 0 {
		return nil, errors.WithStack(ErrNoRows)
	}
	return mapper, err
}

// ScanRows will scan given sqlx.Rows to created its Mapper.
func ScanRows(rows Rows) (Mapper, error) {
	mapper := map[string]interface{}{}
	err := rows.Write(mapper)
	return mapper, err
}
