package sqlxx

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Save saves the model and populate it to the database
func Save(driver Driver, model XModel) error {
	_, err := SaveWithQueries(driver, model)
	return err
}

// SaveWithQueries saves the given instance and returns performed queries.
func SaveWithQueries(driver Driver, model XModel) (Queries, error) {
	queries, err := save(driver, model)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute save")
	}
	return queries, nil
}

func save(driver Driver, model XModel) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()

	query := ""
	params := make(map[string]interface{})

	schema, err := XGetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	returning := []string{}
	columns := []string{}
	values := []string{}

	pk := schema.PrimaryKey()
	id, hasPK := pk.ValueOpt(model)

	for name, column := range schema.Fields {
		if column.IsPrimaryKey {
			continue
		}

		columns = append(columns, column.ColumnName)
		fv, err := GetFieldValue(model, name)
		if err != nil {
			return nil, err
		}

		value := ""
		if column.HasDefault() && IsZero(fv) && !hasPK {
			value = column.Default()
			returning = append(returning, column.ColumnName)
		} else {
			value = fmt.Sprintf(":%s", column.ColumnName)
			params[column.ColumnName] = fv
		}

		values = append(values, value)
	}

	if !hasPK {

		query = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
			schema.TableName(),
			strings.Join(columns, ", "),
			strings.Join(values, ", "),
		)

		returning = append(returning, pk.ColumnName())

	} else {

		updates := []string{}
		for i := range columns {
			updates = append(updates, fmt.Sprintf("%s = %s", columns[i], values[i]))
		}

		params[pk.ColumnName()] = id
		query = fmt.Sprintf(`UPDATE %s SET %s WHERE %s = :%s`,
			schema.TableName(),
			strings.Join(updates, ", "),
			pk.ColumnPath(),
			pk.ColumnName(),
		)
	}

	if len(returning) > 0 {
		query = fmt.Sprintf(`%s RETURNING %s`, query, strings.Join(returning, ", "))
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

	row := stmt.QueryRow(params)
	if row == nil {
		return queries, errors.New("cannot obtain result from driver")
	}
	err = row.Err()
	if err != nil {
		return queries, err
	}

	mapper := map[string]interface{}{}
	err = row.MapScan(mapper)
	if err != nil && err != sql.ErrNoRows {
		return queries, err
	}

	err = model.WriteModel(mapper)
	if err != nil {
		return queries, err
	}

	return queries, nil
}
