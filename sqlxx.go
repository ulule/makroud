package sqlxx

import (
	"fmt"
	"strings"
)

// GetByParams executes a WHERE with params and populates the given model
// instance with related data.
func GetByParams(driver Driver, model Model, params map[string]interface{}) error {
	schema, err := GetSchema(model)
	if err != nil {
		return err
	}

	columns := []string{}
	for _, column := range schema.Columns {
		columns = append(columns, column.PrefixedName)
	}

	wheres := []string{}
	for k := range params {
		wheres = append(wheres, fmt.Sprintf("%s.%s=:%s", model.TableName(), k, k))
	}

	_, err = driver.NamedQuery(fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(columns, ", "), model.TableName(), wheres), params)
	if err != nil {
		return err
	}

	return nil
}

// Preload preloads related fields.
func Preload(driver Driver, out Model, related ...string) error {
	return nil
}

// PreloadFuncs preloads with the given preloader functions.
func PreloadFuncs(driver Driver, out Model, preloaders ...Preloader) error {
	return nil
}
