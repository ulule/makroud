package sqlxx

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Save saves the model and populate it to the database
func Save(driver Driver, out Model) error {
	_, err := SaveWithQueries(driver, out)
	return err
}

// SaveWithQueries saves the given instance and returns performed queries.
func SaveWithQueries(driver Driver, out Model) (Queries, error) {
	queries, err := save(driver, out)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute save")
	}
	return queries, nil
}

func save(driver Driver, out Model) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()

	schema, err := GetSchema(driver, out)
	if err != nil {
		return nil, err
	}

	var (
		columns        = []string{}
		ignoredColumns = []string{}
		values         = []string{}
		params         = make(map[string]interface{})
		query          string
	)

	// TODO Bug with PK
	for name, column := range schema.Fields {
		var (
			isIgnored    bool
			hasDefault   bool
			defaultValue string
			value        string
		)

		tag := column.Tags.Get(StructTagName)
		if tag != nil {
			isIgnored = len(tag.Get(StructTagIgnored)) != 0
			defaultValue = tag.Get(StructTagDefault)
			hasDefault = len(defaultValue) != 0
		}

		if isIgnored || hasDefault {
			ignoredColumns = append(ignoredColumns, column.ColumnName)
		}

		if !isIgnored {
			columns = append(columns, column.ColumnName)

			fv, err := GetFieldValue(out, name)
			if err != nil {
				return nil, err
			}

			if hasDefault && IsZero(fv) {
				value = defaultValue
			} else {
				value = fmt.Sprintf(":%s", column.ColumnName)
				params[column.ColumnName] = fv
			}

			values = append(values, value)
		}
	}

	pkField := schema.PrimaryKeyField

	pk, err := GetFieldValueInt64(out, pkField.FieldName)
	if err != nil {
		return nil, err
	}

	if pk == int64(0) {
		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
			schema.TableName,
			strings.Join(columns, ", "),
			strings.Join(values, ", "),
		)
	} else {
		updates := []string{}
		for i := range columns {
			updates = append(updates, fmt.Sprintf("%s = %s", columns[i], values[i]))
		}

		params[pkField.ColumnName] = pk
		query = fmt.Sprintf(`UPDATE %s SET %s WHERE %s = :%s`,
			schema.TableName,
			strings.Join(updates, ", "),
			pkField.ColumnPath(),
			pkField.ColumnName,
		)
	}

	if len(ignoredColumns) > 0 {
		query = fmt.Sprintf(`%s RETURNING %s`, query, strings.Join(ignoredColumns, ", "))
	}

	queries := Queries{{
		Query:  query,
		Params: params,
	}}

	// Log must be wrapped in a defered function so the duration computation is done when the function return a result.
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return queries, err
	}
	defer driver.close(stmt)

	err = stmt.Get(out, out)
	return queries, err
}
