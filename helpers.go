package sqlxx

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"

	lkb "github.com/ulule/loukoum/builder"
)

// Exec will execute given query from a Loukoum builder.
// If an object is given, it will mutate it to match the row values.
func Exec(driver Driver, builder lkb.Builder, dest ...interface{}) error {
	start := time.Now()
	queries := Queries{NewQuery(builder)}

	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	query, args := builder.NamedQuery()

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot prepare statement")
	}
	defer driver.close(stmt, map[string]string{
		"query": query,
	})

	// TODO Handle single object.
	if len(dest) > 0 {
		err = stmt.Select(dest[0], args)
	} else {
		_, err = stmt.Exec(args)
	}
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return nil
}

// TODO Rename Sync to Save
// Sync will create or update a row from a Loukoum builder.
// If an object is given, it will mutate it to match the row values.
func Sync(driver Driver, builder lkb.Builder, dest interface{}) error {
	start := time.Now()
	queries := Queries{NewQuery(builder)}

	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	query, args := builder.NamedQuery()

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot prepare statement")
	}
	defer driver.close(stmt, map[string]string{
		"query": query,
	})

	err = stmt.Get(dest, args)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return nil
}

// TODO Merge Fetch and List
// Fetch returns an instance from a Loukoum builder.
func Fetch(driver Driver, builder lkb.Builder, dest interface{}) error {
	start := time.Now()
	queries := Queries{NewQuery(builder)}

	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	query, args := builder.NamedQuery()

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot prepare statement")
	}
	defer driver.close(stmt, map[string]string{
		"query": query,
	})

	err = stmt.Get(dest, args)
	if err != nil {
		return errors.Wrap(err, "sqlx: cannot execute query")
	}

	return nil
}

// List returns a slice of instances from a Loukoum builder.
func List(driver Driver, builder lkb.Builder, dest interface{}) error {
	start := time.Now()
	queries := Queries{NewQuery(builder)}

	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	query, args := builder.NamedQuery()

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot prepare statement")
	}
	defer driver.close(stmt, map[string]string{
		"query": query,
	})

	err = stmt.Select(dest, args)
	if err != nil && !IsErrNoRows(err) {
		return errors.Wrap(err, "sqlx: cannot execute query")
	}

	return nil
}

// Count will execute given query to return a number from a aggregate function.
func Count(driver Driver, builder lkb.Builder) (int64, error) {
	count := int64(0)

	err := Fetch(driver, builder, &count)
	if IsErrNoRows(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FloatCount will execute given query to return a number (in float) from a aggregate function.
func FloatCount(driver Driver, builder lkb.Builder) (float64, error) {
	count := float64(0)

	err := Fetch(driver, builder, &count)
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
