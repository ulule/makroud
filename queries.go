package sqlxx

import (
	"fmt"
	"strings"
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
func whereQuery(model Model, params map[string]interface{}, fetchOne bool) (string, error) {
	schema, err := GetSchema(model)
	if err != nil {
		return "", err
	}

	columns := []string{}
	for _, column := range schema.Columns {
		columns = append(columns, column.PrefixedName)
	}

	wheres := []string{}
	for k := range params {
		wheres = append(wheres, fmt.Sprintf("%s.%s=:%s", model.TableName(), k, k))
	}

	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(columns, ", "), model.TableName(), wheres)

	if fetchOne {
		q = fmt.Sprintf("%s LIMIT 1", q)
	}

	return q, nil
}

// where executes a where clause.
func where(driver Driver, models []Model, params map[string]interface{}) error {
	count := len(models)

	if count == 0 {
		return nil
	}

	fetchOne := count == 1

	query, err := whereQuery(models[0], params, fetchOne)
	if err != nil {
		return err
	}

	if fetchOne {
		if err = driver.Get(&models[0], query, params); err != nil {
			return err
		}
		return nil
	}

	// TODO Find with driver.Selec()

	return nil
}
