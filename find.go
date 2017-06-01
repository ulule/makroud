package sqlxx

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// GetByParams executes a where with the given params and populates the given model.
func GetByParams(driver Driver, model XModel, params map[string]interface{}) error {
	_, err := GetByParamsWithQueries(driver, model, params)
	return err
}

// GetByParamsWithQueries executes a where with the given params and populates the given model.
func GetByParamsWithQueries(driver Driver, model XModel, params map[string]interface{}) (Queries, error) {
	start := time.Now()
	queries, err := getByParams(driver, model, params)
	Log(driver, queries, time.Since(start))
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute select")
	}
	return queries, nil
}

// getByParams is a private wrapper which doesn't log statement.
func getByParams(driver Driver, model XModel, params map[string]interface{}) (Queries, error) {
	return where(driver, model, params, true, func(query string, args []interface{}) error {
		row := driver.QueryRowx(driver.Rebind(query), args...)
		if row == nil {
			return errors.New("cannot obtain result from driver")
		}
		err := row.Err()
		if err != nil {
			return err
		}

		mapper, err := ScanRow(row)
		if err != nil {
			return err
		}

		err = model.WriteModel(mapper)
		if err != nil {
			return err
		}

		return nil
	})
}

// FindByParams executes a where with the given params and populates the given models.
func FindByParams(driver Driver, models XModels, params map[string]interface{}) error {
	_, err := FindByParamsWithQueries(driver, models, params)
	return err
}

// FindByParamsWithQueries executes a where with the given params and populates the given models.
func FindByParamsWithQueries(driver Driver, models XModels, params map[string]interface{}) (Queries, error) {
	start := time.Now()
	queries, err := findByParams(driver, models, params)
	Log(driver, queries, time.Since(start))
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute select")
	}
	return queries, nil
}

// findByParams is a private wrapper which doesn't log statement.
func findByParams(driver Driver, models XModels, params map[string]interface{}) (Queries, error) {
	return where(driver, models.Model(), params, false, func(query string, args []interface{}) error {

		rows, err := driver.Queryx(driver.Rebind(query), args...)
		if rows == nil {
			return errors.New("cannot obtain results from driver")
		}
		if err != nil {
			return err
		}
		defer driver.close(rows)
		err = rows.Err()
		if err != nil {
			return err
		}

		for rows.Next() {

			mapper, err := ScanRows(rows)
			if err != nil {
				return err
			}

			err = models.Append(mapper)
			if err != nil {
				return err
			}

		}

		err = rows.Err()
		if err != nil {
			return err
		}

		return nil
	})
}

// whereQuery returns SQL where clause from model and params.
func whereQuery(driver Driver, model XModel, params map[string]interface{},
	fetchOne bool) (string, []interface{}, error) {

	schema, err := XGetSchema(driver, model)
	if err != nil {
		return "", nil, err
	}

	statement := fmt.Sprintf(`SELECT %s FROM %s WHERE %s`,
		schema.ColumnPaths(),
		schema.TableName(),
		schema.WhereColumnPaths(params),
	)

	if fetchOne {
		statement = fmt.Sprint(statement, ` LIMIT 1`)
	}

	query, args, err := sqlx.Named(statement, params)
	if err != nil {
		return "", nil, err
	}

	return sqlx.In(query, args...)
}

// where executes a where clause.
func where(driver Driver, model XModel, params map[string]interface{}, fetchOne bool,
	callback func(query string, args []interface{}) error) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	query, args, err := whereQuery(driver, model, params, fetchOne)
	if err != nil {
		return nil, err
	}

	queries := Queries{{
		Query:  query,
		Args:   args,
		Params: params,
	}}

	return queries, callback(query, args)
}
