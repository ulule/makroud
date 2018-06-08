package sqlxx

import (
	"time"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
)

// Delete deletes the given instance.
func Delete(driver Driver, model Model) error {
	_, err := DeleteWithQueries(driver, model)
	return err
}

// DeleteWithQueries deletes the given instance and returns performed queries.
func DeleteWithQueries(driver Driver, model Model) (Queries, error) {
	queries, err := remove(driver, model)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute delete")
	}
	return queries, nil
}

// Archive archives the given instance.
func Archive(driver Driver, model Model) error {
	_, err := ArchiveWithQueries(driver, model)
	return err
}

// ArchiveWithQueries archives the given instance and returns performed queries.
func ArchiveWithQueries(driver Driver, model Model) (Queries, error) {
	queries, err := archive(driver, model)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute archive")
	}
	return queries, nil
}

func remove(driver Driver, model Model) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
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

	err = Exec(driver, builder)
	return queries, err
}

func archive(driver Driver, model Model) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()
	queries := Queries{}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	if !schema.HasDeletedKey() {
		return nil, errors.Wrapf(ErrDeletedKey, "sqlxx: %T doesn't support archive operation", model)
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

	mapper, err := ScanRow(row)
	if err != nil && !IsErrNoRows(err) {
		return queries, err
	}
	err = schema.WriteModel(mapper, model)
	if err != nil {
		return queries, err
	}

	return queries, err
}
