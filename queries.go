package sqlxx

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
)

// Delete deletes the model in the database
func Delete(driver Driver, out Model) error {
	schema, err := GetSchema(out)

	if err != nil {
		return err
	}

	primaryColumn := schema.PrimaryColumn

	// GO TO HELL ZERO VALUES DELETION
	if !primaryColumn.HasValue() {
		return fmt.Errorf("%v has no primary key, cannot be deleted", out)
	}

	wheres := []string{fmt.Sprintf("%s = :%s", primaryColumn.Name, primaryColumn.Name)}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s",
		out.TableName(),
		strings.Join(wheres, ", "))

	_, err = driver.NamedExec(query, out)

	if err != nil {
		return err
	}

	return nil
}

// Save saves the model and populate it to the database
func Save(driver Driver, out Model) error {
	schema, err := GetSchema(out)

	if err != nil {
		return err
	}

	columns := []string{}
	values := []string{}
	ignoredColumns := []string{}

	for _, column := range schema.Columns {
		_, isIgnored := column.Tags["ignored"]
		defaultValue, hasDefault := column.Tags["default"]

		if isIgnored || hasDefault {
			ignoredColumns = append(ignoredColumns, column.Name)
		}

		if !isIgnored {
			columns = append(columns, column.Name)

			if hasDefault {
				values = append(values, defaultValue)
			} else {
				values = append(values, fmt.Sprintf(":%s", column.Name))
			}
		}
	}

	var query string

	primaryColumn := schema.PrimaryColumn

	if !primaryColumn.HasValue() {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			out.TableName(),
			strings.Join(columns, ", "),
			strings.Join(values, ", "))

	} else {
		updates := []string{}

		for i := range columns {
			updates = append(updates, fmt.Sprintf("%s = %s", columns[i], values[i]))
		}

		wheres := []string{fmt.Sprintf("%s = :%s", primaryColumn.Name, primaryColumn.Name)}

		query = fmt.Sprintf("UPDATE %s SET %s WHERE %s",
			out.TableName(),
			strings.Join(updates, ", "),
			strings.Join(wheres, ", "))
	}

	if len(ignoredColumns) > 0 {
		query = fmt.Sprintf("%s RETURNING %s", query, strings.Join(ignoredColumns, ", "))
	}

	stmt, err := driver.PrepareNamed(query)

	if err != nil {
		return err
	}

	err = stmt.Get(out, out)

	if err != nil {
		return err
	}

	return nil
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

	columns := []string{}
	for _, column := range schema.Columns {
		columns = append(columns, column.PrefixedName())
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
