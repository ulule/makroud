package sqlxx

import (
	"database/sql/driver"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/ulule/sqlxx/reflekt"
)

// GetPrimaryKeys returns primary keys for the given interface.
func GetPrimaryKeys(out interface{}, name string) ([]interface{}, error) {
	var (
		value  = reflekt.GetIndirectValue(GetModelFromInterface(out))
		isNull = reflekt.IsNullableType(value.FieldByName(name).Type())
	)

	pks, err := reflekt.GetFieldValues(out, name)
	if err != nil {
		return nil, err
	}

	var values []interface{}
	for i := range pks {
		if !isNull {
			if reflekt.IsZeroValue(pks[i]) {
				return nil, fmt.Errorf("Cannot perform query on zero value (%s=%v)", name, pks[i])
			}
			values = append(values, pks[i])
		} else {
			valuer := reflekt.Copy(pks[i]).(driver.Valuer)
			if v, err := valuer.Value(); err == nil && v != nil {
				values = append(values, v)
			}
		}
	}

	return values, nil
}

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
func where(driver Driver, out interface{}, params map[string]interface{}, fetchOne bool) error {
	model := GetModelFromInterface(out)

	query, args, err := whereQuery(model, params, fetchOne)
	if err != nil {
		return err
	}

	if fetchOne {
		return driver.Get(out, driver.Rebind(query), args...)
	}

	return driver.Select(out, driver.Rebind(query), args...)
}
