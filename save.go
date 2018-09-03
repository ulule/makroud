package sqlxx

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"

	"github.com/ulule/sqlxx/reflectx"
)

// Save saves the given instance.
func Save(ctx context.Context, driver Driver, model Model) error {
	_, err := SaveWithQueries(ctx, driver, model)
	return err
}

// SaveWithQueries saves the given instance and returns performed queries.
func SaveWithQueries(ctx context.Context, driver Driver, model Model) (Queries, error) {
	queries, err := save(ctx, driver, model)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute save")
	}
	return queries, nil
}

func save(ctx context.Context, driver Driver, model Model) (Queries, error) {
	if driver == nil {
		return nil, errors.WithStack(ErrInvalidDriver)
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

		} else if column.IsUpdatedKey() && hasPK {

			values[column.ColumnName()] = loukoum.Raw("NOW()")
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
			err = schema.WriteModel(mapper, model)
			if err != nil {
				return nil, err
			}
		}

		builder = loukoum.Insert(schema.TableName()).
			Set(values).
			Returning(returning)

	} else {

		builder = loukoum.Update(schema.TableName()).
			Set(values).
			Where(loukoum.Condition(pk.ColumnName()).Equal(id)).
			Returning(returning)

	}

	queries = append(queries, NewQuery(builder))
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	query, args := builder.NamedQuery()
	err = exec(ctx, driver, query, args, model)

	// Ignore no rows error if returning is empty.
	if IsErrNoRows(err) && len(returning) == 0 {
		return queries, nil
	}

	return queries, err
}
