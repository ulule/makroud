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
	for _, f := range schema.Fields {
		columns = append(columns, f.PrefixedName())
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
	var (
		model Model
		typ   reflect.Type
	)

	value := reflect.ValueOf(out)

	if reflect.Indirect(value).Kind() == reflect.Slice {
		typ = value.Type().Elem().Elem()
	} else {
		typ = reflect.Indirect(value).Type()
	}

	model = reflect.New(typ).Interface().(Model)

	query, args, err := whereQuery(model, params, fetchOne)
	if err != nil {
		return err
	}

	if fetchOne {
		return driver.Get(value.Interface(), driver.Rebind(query), args...)
	}

	return driver.Select(value.Interface(), driver.Rebind(query), args...)
}
