package sqlxx

import (
	"time"

	"github.com/pkg/errors"
	lk "github.com/ulule/loukoum"
	lkb "github.com/ulule/loukoum/builder"
)

// Save saves the given instance.
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
	queries := Queries{}

	schema, err := XGetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	values := lk.Map{}
	returning := []string{}

	pk := schema.PrimaryKey()
	id, hasPK := pk.ValueOpt(model)

	for name, column := range schema.Fields {
		if column.IsPrimaryKey {
			continue
		}

		value, err := GetFieldValue(model, name)
		if err != nil {
			return nil, err
		}

		if column.HasDefault() && IsZero(value) && !hasPK {
			values[column.ColumnName()] = lk.Raw(column.Default())
			returning = append(returning, column.ColumnName())
		} else {
			values[column.ColumnName()] = value
		}
	}

	var builder lkb.Builder

	if !hasPK {
		returning = append(returning, pk.ColumnName())
		builder = lk.Insert(schema.TableName()).Set(values).Returning(returning)
	} else {
		builder = lk.Update(schema.TableName()).Set(values).Where(lk.Condition(pk.ColumnName()).Equal(id))
	}

	queries = append(queries, NewQuery(builder))
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	query, args := builder.NamedQuery()

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return queries, err
	}
	defer driver.close(stmt, map[string]string{
		"query": query,
	})

	row := stmt.QueryRow(args)
	if row == nil {
		return queries, errors.New("sqlxx: cannot obtain result from driver")
	}
	err = row.Err()
	if err != nil {
		return queries, err
	}

	mapper := map[string]interface{}{}
	err = row.MapScan(mapper)
	if err != nil && !IsErrNoRows(err) {
		return queries, err
	}

	err = model.WriteModel(mapper)
	if err != nil {
		return queries, err
	}

	return queries, nil
}
