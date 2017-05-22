package sqlxx

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// GetByParams executes a where with the given params and populates the given model.
func GetByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	_, err := GetByParamsWithQueries(driver, out, params)
	return err
}

// GetByParamsWithQueries executes a where with the given params and populates the given model.
func GetByParamsWithQueries(driver Driver, out interface{}, params map[string]interface{}) (Queries, error) {
	start := time.Now()
	queries, err := where(driver, out, params, true)
	Log(driver, queries, time.Since(start))
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute select")
	}
	return queries, nil
}

// getByParams is a private wrapper which doesn't log statement.
func getByParams(driver Driver, out interface{}, params map[string]interface{}) (Queries, error) {
	return where(driver, out, params, true)
}

// FindByParams executes a where with the given params and populates the given models.
func FindByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	_, err := FindByParamsWithQueries(driver, out, params)
	return err
}

// FindByParamsWithQueries executes a where with the given params and populates the given models.
func FindByParamsWithQueries(driver Driver, out interface{}, params map[string]interface{}) (Queries, error) {
	start := time.Now()
	queries, err := where(driver, out, params, false)
	Log(driver, queries, time.Since(start))
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute select")
	}
	return queries, nil
}

// findByParams is a private wrapper which doesn't log statement.
func findByParams(driver Driver, out interface{}, params map[string]interface{}) (Queries, error) {
	return where(driver, out, params, false)
}

// whereQuery returns SQL where clause from model and params.
func whereQuery(driver Driver, model Model, params map[string]interface{},
	fetchOne bool) (string, []interface{}, error) {

	schema, err := GetSchema(driver, model)
	if err != nil {
		return "", nil, err
	}

	q := fmt.Sprintf(`SELECT %s FROM %s WHERE %s`,
		schema.ColumnPaths(),
		model.TableName(),
		schema.WhereColumnPaths(params),
	)

	if fetchOne {
		q = fmt.Sprintf(`%s LIMIT 1`, q)
	}

	query, args, err := sqlx.Named(q, params)
	if err != nil {
		return "", nil, err
	}

	return sqlx.In(query, args...)
}

// where executes a where clause.
func where(driver Driver, out interface{}, params map[string]interface{}, fetchOne bool) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	model := ToModel(out)

	query, args, err := whereQuery(driver, model, params, fetchOne)
	if err != nil {
		return nil, err
	}

	queries := Queries{{
		Query:  query,
		Args:   args,
		Params: params,
	}}

	if fetchOne {
		return queries, driver.Get(out, driver.Rebind(query), args...)
	}
	return queries, driver.Select(out, driver.Rebind(query), args...)
}
