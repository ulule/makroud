package sqlxx

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// GetByParams executes a where with the given params and populates the given model.
func GetByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	_, err := where(driver, out, params, true)
	return err
}

// GetByParamsWithQueries executes a where with the given params and populates the given model.
func GetByParamsWithQueries(driver Driver, out interface{}, params map[string]interface{}) (Queries, error) {
	return where(driver, out, params, true)
}

// FindByParams executes a where with the given params and populates the given models.
func FindByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	_, err := where(driver, out, params, false)
	return err
}

// FindByParamsWithQueries executes a where with the given params and populates the given models.
func FindByParamsWithQueries(driver Driver, out interface{}, params map[string]interface{}) (Queries, error) {
	return where(driver, out, params, false)
}

// whereQuery returns SQL where clause from model and params.
func whereQuery(model Model, params map[string]interface{}, fetchOne bool) (string, []interface{}, error) {
	schema, err := GetSchema(model)
	if err != nil {
		return "", nil, err
	}

	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		schema.ColumnPaths(),
		model.TableName(),
		schema.WhereColumnPaths(params))

	if fetchOne {
		q = fmt.Sprintf("%s LIMIT 1", q)
	}

	query, args, err := sqlx.Named(q, params)
	if err != nil {
		return "", nil, err
	}

	return sqlx.In(query, args...)
}

// where executes a where clause.
func where(driver Driver, out interface{}, params map[string]interface{}, fetchOne bool) (Queries, error) {
	model := GetModelFromInterface(out)

	query, args, err := whereQuery(model, params, fetchOne)
	if err != nil {
		return nil, err
	}

	queries := Queries{{Query: query, Args: args, Params: params}}

	if fetchOne {
		return queries, driver.Get(out, driver.Rebind(query), args...)
	}

	return queries, driver.Select(out, driver.Rebind(query), args...)
}
