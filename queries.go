package sqlxx

import (
	"fmt"
	"strings"
)

// whereQuery returns SQL where clause from model and params.
func whereQuery(models []Model, params map[string]interface{}) (string, error) {
	count := len(models)
	if count == 0 {
		return "", nil
	}

	var (
		fetchOne = count == 1
		model    = models[0]
	)

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
	query, err := whereQuery(models, params)
	if err != nil {
		return err
	}

	_, err = driver.NamedQuery(query, params)
	if err != nil {
		return err
	}

	return nil
}
