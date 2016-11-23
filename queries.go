package sqlxx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// GetByParams executes a where with the given params and populates the given model.
func GetByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	return where(driver, out, params, true)
}

// FindByParams executes a where with the given params and populates the given models.
func FindByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	return where(driver, out, params, false)
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
func where(driver Driver, out interface{}, params map[string]interface{}, fetchOne bool) error {
	v := reflect.ValueOf(out).Interface().(Model)

	query, args, err := whereQuery(v, params, fetchOne)
	if err != nil {
		return err
	}

	if fetchOne {
		return driver.Get(&out, driver.Rebind(query), args...)
	}

	return driver.Select(&out, driver.Rebind(query), args...)
}
