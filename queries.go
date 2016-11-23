package sqlxx

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// GetByParams executes a where with the given params and populates the given model.
func GetByParams(driver Driver, model Model, params map[string]interface{}) error {
	return where(driver, []Model{model}, params)
}

// FindByParams executes a where with the given params and populates the given models.
func FindByParams(driver Driver, models []Model, params map[string]interface{}) error {
	return where(driver, models, params)
}

// whereQuery returns SQL where clause from model and params.
func whereQuery(model Model, params map[string]interface{}, fetchOne bool) (string, []interface{}, error) {
	schema, err := GetSchema(model)
	if err != nil {
		return "", nil, err
	}

	columns := []string{}
	for _, column := range schema.Columns {
		columns = append(columns, column.PrefixedName)
	}

	wheres := []string{}
	for k := range params {
		wheres = append(wheres, fmt.Sprintf("%s.%s=:%s", model.TableName(), k, k))
	}

	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(columns, ", "),
		model.TableName(),
		strings.Join(wheres, ","))

	if fetchOne {
		q = fmt.Sprintf("%s LIMIT 1", q)
	}

	return sqlx.Named(q, params)
}

// where executes a where clause.
func where(driver Driver, models []Model, params map[string]interface{}) error {
	count := len(models)

	if count == 0 {
		return nil
	}

	fetchOne := count == 1

	query, args, err := whereQuery(models[0], params, fetchOne)
	if err != nil {
		return err
	}

	if fetchOne {
		return driver.Get(&models[0], driver.Rebind(query), args...)
	}

	return driver.Select(&models, driver.Rebind(query), args...)
}
