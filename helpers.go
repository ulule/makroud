package sqlxx

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// ExecInParams will execute query with given array of parameters.
func ExecInParams(driver Driver, query string, data interface{}) error {
	_, err := ExecInParamsWithQueries(driver, query, data)
	return err
}

// ExecInParams will execute query with given array of parameters and returns performed queries.
func ExecInParamsWithQueries(driver Driver, query string, data interface{}) (Queries, error) {
	start := time.Now()

	queries := Queries{{
		Query: query,
		Args:  []interface{}{data},
	}}

	// Log must be wrapped in a defered function so the duration computation is done when the function return a result.
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	fullquery, fulldata, err := sqlx.In(query, data)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot prepare statement")
	}

	queries[0].Query = fullquery
	queries[0].Args = fulldata

	fullquery = driver.Rebind(fullquery)
	_, err = driver.Exec(fullquery, fulldata...)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute statement")
	}

	return queries, nil
}

// FindInParams will find every rows that matches given array of parameters.
func FindInParams(driver Driver, out interface{}, query string, data interface{}) error {
	_, err := FindInParamsWithQueries(driver, out, query, data)
	return err
}

// FindInParamsWithQueries will find every rows that matches given array of parameters and returns performed queries.
func FindInParamsWithQueries(driver Driver, out interface{}, query string, data interface{}) (Queries, error) {
	start := time.Now()

	queries := Queries{{
		Query: query,
		Args:  []interface{}{data},
	}}

	// Log must be wrapped in a defered function so the duration computation is done when the function return a result.
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	fullquery, fulldata, err := sqlx.In(query, data)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot prepare statement")
	}

	queries[0].Query = fullquery
	queries[0].Args = fulldata

	fullquery = driver.Rebind(fullquery)
	err = driver.Select(out, fullquery, fulldata...)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return queries, nil
}

// Exec will execute given query.
func Exec(driver Driver, query string, out interface{}) error {
	start := time.Now()

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot prepare statement")
	}
	defer driver.close(stmt)

	params, err := GetValues(stmt.Params, out, stmt.Stmt.Mapper)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot prepare statement")
	}

	queries := Queries{{
		Query:  query,
		Params: params,
	}}

	// Log must be wrapped in a defered function so the duration computation is done when the function return a result.
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	_, err = stmt.Exec(out)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return nil
}

// NamedExec will execute given query.
func NamedExec(driver Driver, query string, params map[string]interface{}) error {
	start := time.Now()

	queries := Queries{{
		Query:  query,
		Params: params,
	}}

	// Log must be wrapped in a defered function so the duration computation is done when the function return a result.
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	_, err := driver.NamedExec(query, params)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return nil
}

// Sync will create or update a row from given query.
// It will also mutate the given object to match the row values.
func Sync(driver Driver, query string, out interface{}) error {
	start := time.Now()

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot prepare statement")
	}
	defer driver.close(stmt)

	params, err := GetValues(stmt.Params, out, stmt.Stmt.Mapper)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot prepare statement")
	}

	queries := Queries{{
		Query:  query,
		Params: params,
	}}

	// Log must be wrapped in a defered function so the duration computation is done when the function return a result.
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	// run query with out values
	// then override out with returned value
	err = stmt.Get(out, out)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute query")
	}

	return nil
}