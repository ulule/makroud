package sqlxx

import (
	"github.com/jmoiron/sqlx"
)

// ----------------------------------------------------------------------------
// Mapper
// ----------------------------------------------------------------------------

// Mapper will be used to mutate a Model with row values.
type Mapper map[string]interface{}

// ScanRow will scan given sqlx.Row to created its Mapper.
func ScanRow(row *sqlx.Row) (Mapper, error) {
	mapper := map[string]interface{}{}
	err := row.MapScan(mapper)
	return mapper, err
}

// ScanRows will scan given sqlx.Rows to created its Mapper.
func ScanRows(rows *sqlx.Rows) (Mapper, error) {
	mapper := map[string]interface{}{}
	err := rows.MapScan(mapper)
	return mapper, err
}
