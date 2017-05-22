package sqlxx

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

// Delete deletes the given instance.
func Delete(driver Driver, out interface{}) error {
	_, err := DeleteWithQueries(driver, out)
	return err
}

// DeleteWithQueries deletes the given instance and returns performed queries.
func DeleteWithQueries(driver Driver, out interface{}) (Queries, error) {
	queries, err := remove(driver, out)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute delete")
	}
	return queries, nil
}

// SoftDelete is an alias for Archive.
func SoftDelete(driver Driver, out interface{}, fieldName string) error {
	return Archive(driver, out, fieldName)
}

// SoftDeleteWithQueries is an alias for Archive.
func SoftDeleteWithQueries(driver Driver, out interface{}, fieldName string) (Queries, error) {
	return ArchiveWithQueries(driver, out, fieldName)
}

// Archive archives the given instance.
func Archive(driver Driver, out interface{}, fieldName string) error {
	_, err := ArchiveWithQueries(driver, out, fieldName)
	return err
}

// ArchiveWithQueries archives the given instance and returns performed queries.
func ArchiveWithQueries(driver Driver, out interface{}, fieldName string) (Queries, error) {
	queries, err := archive(driver, out, fieldName)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute archive")
	}
	return queries, nil
}

func remove(driver Driver, out interface{}) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()

	schema, err := GetSchema(driver, out)
	if err != nil {
		return nil, err
	}

	pkField := schema.PrimaryKeyField

	pk, err := GetFieldValueInt64(out, pkField.FieldName)
	if err != nil {
		return nil, err
	}

	if pk == int64(0) {
		return nil, errors.Errorf("%v cannot be deleted (no primary key)", out)
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE %s = :%s`,
		schema.TableName,
		pkField.ColumnPath(),
		pkField.ColumnName,
	)

	params := map[string]interface{}{
		pkField.ColumnName: pk,
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

func archive(driver Driver, out interface{}, fieldName string) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	start := time.Now()

	schema, err := GetSchema(driver, out)
	if err != nil {
		return nil, err
	}

	var (
		pkField        = schema.PrimaryKeyField
		deletedAtField = schema.Fields[fieldName]
		now            = time.Now().UTC()
	)

	pk, err := GetFieldValueInt64(out, pkField.FieldName)
	if err != nil {
		return nil, err
	}

	if pk == int64(0) {
		return nil, errors.Errorf("%v cannot be archived (no primary key)", out)
	}

	query := fmt.Sprintf(`UPDATE %s SET %s = :%s WHERE %s = :%s`,
		schema.TableName,
		deletedAtField.ColumnName,
		deletedAtField.ColumnName,
		pkField.ColumnPath(),
		pkField.ColumnName,
	)

	params := map[string]interface{}{
		deletedAtField.ColumnName: now,
		pkField.ColumnName:        pk,
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
