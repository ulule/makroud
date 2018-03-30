package sqlxx

import (
	"time"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
)

// Delete deletes the given instance.
func Delete(driver Driver, model XModel) error {
	_, err := DeleteWithQueries(driver, model)
	return err
}

// DeleteWithQueries deletes the given instance and returns performed queries.
func DeleteWithQueries(driver Driver, model XModel) (Queries, error) {
	queries, err := remove(driver, model)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute delete")
	}
	return queries, nil
}

// SoftDelete is an alias for Archive.
func SoftDelete(driver Driver, model XModel) error {
	return Archive(driver, model)
}

// SoftDeleteWithQueries is an alias for ArchiveWithQueries.
func SoftDeleteWithQueries(driver Driver, model XModel) (Queries, error) {
	return ArchiveWithQueries(driver, model)
}

// Archive archives the given instance.
func Archive(driver Driver, model XModel) error {
	_, err := ArchiveWithQueries(driver, model)
	return err
}

// ArchiveWithQueries archives the given instance and returns performed queries.
func ArchiveWithQueries(driver Driver, model XModel) (Queries, error) {
	queries, err := archive(driver, model)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute archive")
	}
	return queries, nil
}

// TODO Queries ???
func remove(driver Driver, model XModel) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()
	queries := Queries{}

	schema, err := XGetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	pk := schema.PrimaryKey()
	id, err := pk.Value(model)
	if err != nil {
		return queries, errors.Wrapf(err, "sqlxx: %T cannot be deleted", model)
	}

	builder := loukoum.Delete(schema.TableName()).
		Where(loukoum.Condition(pk.ColumnPath()).Equal(id))

	queries = append(queries, NewQuery(builder))
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	err = Exec(driver, builder)
	return queries, err
}

// TODO Queries ???
func archive(driver Driver, model XModel) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()
	queries := Queries{}

	schema, err := XGetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	if !schema.HasArchiveKey() {
		return nil, errors.Wrapf(err, "sqlxx: %T doesn't support archive operation", model)
	}

	pk := schema.PrimaryKey()
	id, err := pk.Value(model)
	if err != nil {
		return nil, errors.Wrapf(err, "sqlxx: %T cannot be archived", model)
	}

	builder := loukoum.Update(schema.TableName()).
		Set(loukoum.Pair(schema.ArchiveKeyPath(), schema.ArchiveKeyValue())).
		Where(loukoum.Condition(pk.ColumnPath()).Equal(id))

	queries = append(queries, NewQuery(builder))
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	err = Exec(driver, builder, model)
	return queries, err
}
