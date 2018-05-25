package sqlxx

import (
	"time"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"

	"github.com/ulule/sqlxx/reflectx"
)

// Save saves the given instance.
func Save(driver Driver, model Model) error {
	_, err := SaveWithQueries(driver, model)
	return err
}

// SaveWithQueries saves the given instance and returns performed queries.
func SaveWithQueries(driver Driver, model Model) (Queries, error) {
	queries, err := save(driver, model)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute save")
	}
	return queries, nil
}

func save(driver Driver, model Model) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()
	queries := Queries{}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	values := loukoum.Map{}
	returning := []string{}

	pk := schema.PrimaryKey()
	id, hasPK := pk.ValueOpt(model)

	for _, column := range schema.fields {
		if column.IsPrimaryKey() {
			continue
		}

		value, err := reflectx.GetFieldValue(model, column.FieldName())
		if err != nil {
			return nil, err
		}

		if column.HasDefault() && reflectx.IsZero(value) && !hasPK {
			returning = append(returning, column.ColumnName())
		} else {
			values[column.ColumnName()] = value
		}
	}

	var builder builder.Builder

	if !hasPK {
		switch pk.Default() {
		case PrimaryKeyDBDefault:
			returning = append(returning, pk.ColumnName())

		case PrimaryKeyULIDDefault:
			ulid := GenerateULID(driver)
			mapper := map[string]interface{}{
				pk.ColumnName(): ulid,
			}
			values[pk.ColumnName()] = ulid
			schema.WriteModel(mapper, model)
		}

		builder = loukoum.Insert(schema.TableName()).Set(values).Returning(returning)
	} else {
		builder = loukoum.Update(schema.TableName()).Set(values).Where(loukoum.Condition(pk.ColumnName()).Equal(id))
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
	if len(mapper) > 0 {
		err = schema.WriteModel(mapper, model)
		if err != nil {
			return queries, err
		}
	}

	return queries, nil
}
