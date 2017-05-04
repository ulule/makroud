package sqlxx

import (
	"fmt"
	"strings"
)

// Save saves the model and populate it to the database
func Save(driver Driver, out interface{}) error {
	schema, err := GetSchema(out)
	if err != nil {
		return err
	}

	var (
		columns        = []string{}
		ignoredColumns = []string{}
		values         = []string{}
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

	var query string

	pkField := schema.PrimaryKeyField

	pkValue, err := GetFieldValue(out, pkField.Name)
	if err != nil {
		return err
	}

	if IsZeroValue(pkValue) {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			schema.TableName,
			strings.Join(columns, ", "),
			strings.Join(values, ", "))
	} else {
		updates := []string{}

		for i := range columns {
			updates = append(updates, fmt.Sprintf("%s = %s", columns[i], values[i]))
		}

		query = fmt.Sprintf("UPDATE %s SET %s WHERE %s = :%s",
			schema.TableName,
			strings.Join(updates, ", "),
			pkField.ColumnName,
			pkField.ColumnName)
	}

	if len(ignoredColumns) > 0 {
		query = fmt.Sprintf("%s RETURNING %s", query, strings.Join(ignoredColumns, ", "))
	}

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.Get(out, out)
	if err != nil {
		return err
	}

	return nil
}
