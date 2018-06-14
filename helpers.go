package sqlxx

import (
	"database/sql"
	"time"

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

	return exec(driver, query, args, dest...)
}

// RawExec will execute given query.
// If an object is given, it will mutate it to match the row values.
func RawExec(driver Driver, query string, dest ...interface{}) error {
	start := time.Now()
	queries := Queries{NewRawQuery(query)}

	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	return exec(driver, query, nil, dest...)
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
			err = named.Select(dest[0], args)
		} else {
			err = named.Get(dest[0], args)
		}
	} else {
		_, err = named.Exec(args)
	}
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return nil
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
