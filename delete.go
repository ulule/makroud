package sqlxx

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
)

// Delete deletes the given instance.
func Delete(ctx context.Context, driver Driver, model Model) error {
	_, err := DeleteWithQueries(ctx, driver, model)
	return err
}

// DeleteWithQueries deletes the given instance and returns performed queries.
func DeleteWithQueries(ctx context.Context, driver Driver, model Model) (Queries, error) {
	queries, err := remove(ctx, driver, model)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute delete")
	}
	return queries, nil
}

// Archive archives the given instance.
func Archive(ctx context.Context, driver Driver, model Model) error {
	_, err := ArchiveWithQueries(ctx, driver, model)
	return err
}

// ArchiveWithQueries archives the given instance and returns performed queries.
func ArchiveWithQueries(ctx context.Context, driver Driver, model Model) (Queries, error) {
	queries, err := archive(ctx, driver, model)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute archive")
	}
	return queries, nil
}

func remove(ctx context.Context, driver Driver, model Model) (Queries, error) {
	if driver == nil {
		return nil, errors.WithStack(ErrInvalidDriver)
	}

	start := time.Now()
	queries := Queries{}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	pk := schema.PrimaryKey()
	id, err := pk.Value(model)
	if err != nil {
		return queries, errors.Wrapf(err, "sqlxx: %T cannot be deleted", model)
	}

	builder := loukoum.Delete(schema.TableName()).
		Where(loukoum.Condition(pk.ColumnName()).Equal(id))

	queries = append(queries, NewQuery(builder))
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	err = Exec(ctx, driver, builder)
	return queries, err
}

func archive(ctx context.Context, driver Driver, model Model) (Queries, error) {
	if driver == nil {
		return nil, errors.WithStack(ErrInvalidDriver)
	}

	start := time.Now()
	queries := Queries{}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	if !schema.HasDeletedKey() {
		return nil, errors.Wrapf(ErrSchemaDeletedKey, "sqlxx: %T doesn't support archive operation", model)
	}

	pk := schema.PrimaryKey()
	id, err := pk.Value(model)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlxx: %T cannot be archived", model)
	}

	builder := loukoum.Update(schema.TableName()).
		Set(loukoum.Pair(schema.DeletedKeyName(), loukoum.Raw("NOW()"))).
		Where(loukoum.Condition(pk.ColumnName()).Equal(id)).
		Returning(schema.DeletedKeyName())

	queries = append(queries, NewQuery(builder))
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	query, args := builder.NamedQuery()
	err = exec(ctx, driver, query, args, model)
	return queries, err
}
