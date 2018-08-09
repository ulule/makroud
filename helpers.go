package sqlxx

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/ulule/loukoum/builder"

	"github.com/ulule/sqlxx/reflectx"
)

// Exec will execute given query from a Loukoum builder.
// If an object is given, it will mutate it to match the row values.
func Exec(driver Driver, stmt builder.Builder, dest ...interface{}) error {
	start := time.Now()
	queries := Queries{NewQuery(stmt)}

	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	query, args := stmt.NamedQuery()

	err := exec(driver, query, args, dest...)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return nil
}

// RawExec will execute given query.
// If an object is given, it will mutate it to match the row values.
func RawExec(driver Driver, query string, dest ...interface{}) error {
	start := time.Now()
	queries := Queries{NewRawQuery(query)}

	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	err := exec(driver, query, nil, dest...)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return nil
}

func exec(driver Driver, query string, args map[string]interface{}, dest ...interface{}) error {
	named, err := driver.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot prepare statement")
	}
	defer driver.close(named, map[string]string{
		"query": query,
	})

	if len(dest) > 0 {
		if reflectx.IsSlice(dest[0]) {
			return execRows(driver, named, args, dest[0])
		} else {
			return execRow(driver, named, args, dest[0])
		}
	}

	_, err = named.Exec(args)
	return err
}

func execRows(driver Driver, named *sqlx.NamedStmt, args map[string]interface{}, dest interface{}) error {
	if !reflectx.IsPointer(dest) {
		return errors.Wrapf(ErrPointerRequired, "cannot execute query on %T", dest)
	}

	model, ok := reflectx.NewSliceValue(dest).(Model)
	if !ok {
		return named.Select(dest, args)
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return err
	}

	rows, err := named.Queryx(args)
	if err != nil {
		return err
	}

	list := reflectx.NewSlice(reflectx.GetSliceType(dest))

	for rows.Next() {
		mapper, err := ScanRows(rows)
		if err != nil {
			return err
		}

		row := reflectx.NewSliceValue(dest)
		err = schema.WriteModel(mapper, row.(Model))
		if err != nil {
			return err
		}

		reflectx.AppendSlice(list, row)
	}

	reflectx.CopySlice(dest, list)

	return nil
}

func execRow(driver Driver, named *sqlx.NamedStmt, args map[string]interface{}, dest interface{}) error {
	if !reflectx.IsPointer(dest) {
		return errors.Wrapf(ErrPointerRequired, "cannot execute query on %T", dest)
	}

	model, ok := dest.(Model)
	if !ok {
		return named.Get(dest, args)
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return err
	}

	row := named.QueryRow(args)
	if row == nil {
		return errors.Wrap(sql.ErrNoRows, "cannot obtain result from driver")
	}

	err = row.Err()
	if err != nil {
		return err
	}

	mapper, err := ScanRow(row)
	if err != nil && !IsErrNoRows(err) {
		return err
	}

	return schema.WriteModel(mapper, model)
}

// Count will execute given query to return a number from a aggregate function.
func Count(driver Driver, stmt builder.Builder) (int64, error) {
	count := int64(0)

	err := Exec(driver, stmt, &count)
	if IsErrNoRows(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FloatCount will execute given query to return a number (in float) from a aggregate function.
func FloatCount(driver Driver, stmt builder.Builder) (float64, error) {
	count := float64(0)

	err := Exec(driver, stmt, &count)
	if IsErrNoRows(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

// IsErrNoRows returns if given error is a "no rows" error.
func IsErrNoRows(err error) bool {
	return err != nil && errors.Cause(err) == sql.ErrNoRows
}
