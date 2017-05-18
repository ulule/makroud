package sqlxx

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Save saves the model and populate it to the database
func Save(driver Driver, out interface{}) error {
	_, err := SaveWithQueries(driver, out)
	return err
}

// SaveWithQueries saves the given instance and returns performed queries.
func SaveWithQueries(driver Driver, out interface{}) (Queries, error) {
	queries, err := save(driver, out)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute save")
	}
	return queries, nil
}

func save(driver Driver, out interface{}) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	schema, err := GetSchema(driver, out)
	if err != nil {
		return nil, err
	}

	var (
		columns        = []string{}
		ignoredColumns = []string{}
		values         = []string{}
		query          string
	)

	for _, column := range schema.Fields {
		var (
			isIgnored    bool
			hasDefault   bool
			defaultValue string
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
			if hasDefault {
				values = append(values, defaultValue)
			} else {
				values = append(values, fmt.Sprintf(":%s", column.ColumnName))
			}
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
		Query: query,
	}}

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return queries, err
	}
	defer stmt.Close()

	err = stmt.Get(out, out)
	return queries, err
}
