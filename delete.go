package sqlxx

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
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

// SoftDeleteWithQueries is an alias for Archive.
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

func remove(driver Driver, model XModel) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()

	schema, err := XGetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	pk := schema.PrimaryKey()
	id, err := pk.Value(model)
	if err != nil {
		return nil, errors.Wrapf(err, "%T cannot be deleted", model)
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = :%s`,
		schema.TableName(),
		pk.ColumnPath(),
		pk.ColumnName(),
	)

	params := map[string]interface{}{
		pk.ColumnName(): id,
	}

	queries := Queries{{
		Query:  query,
		Params: params,
	}}

	// Log must be wrapped in a defered function so the duration computation is done when the function return a result.
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	_, err = driver.NamedExec(query, params)
	return queries, err
}

func archive(driver Driver, model XModel) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()

	schema, err := XGetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	if !schema.HasArchiveKey() {
		return nil, errors.New("model doesn't support archive operation")
	}

	pk := schema.PrimaryKey()
	id, err := pk.Value(model)
	if err != nil {
		return nil, err
	}

	archiveKey, archiveValue := schema.ArchiveKey()

	query := fmt.Sprintf(`UPDATE %s SET %s = :%s WHERE %s = :%s`,
		schema.TableName(),
		archiveKey,
		archiveKey,
		pk.ColumnPath(),
		pk.ColumnName(),
	)

	params := map[string]interface{}{
		archiveKey:      archiveValue,
		pk.ColumnName(): id,
	}

	queries := Queries{{
		Query:  query,
		Params: params,
	}}

	// Log must be wrapped in a defered function so the duration computation is done when the function return a result.
	defer func() {
		Log(driver, queries, time.Since(start))
	}()

	_, err = driver.NamedExec(query, params)
	return queries, err
}
