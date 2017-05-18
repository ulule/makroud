package sqlxx

import (
	"fmt"
	"time"
)

// Delete deletes the given instance.
func Delete(driver Driver, out interface{}) error {
	_, err := remove(driver, out)
	return err
}

// DeleteWithQueries deletes the given instance and returns performed queries.
func DeleteWithQueries(driver Driver, out interface{}) (Queries, error) {
	return remove(driver, out)
}

// SoftDelete is an alias for Archive.
func SoftDelete(driver Driver, out interface{}, fieldName string) error {
	_, err := archive(driver, out, fieldName)
	return err
}

// SoftDeleteWithQueries is an alias for Archive.
func SoftDeleteWithQueries(driver Driver, out interface{}, fieldName string) (Queries, error) {
	return archive(driver, out, fieldName)
}

// Archive archives the given instance.
func Archive(driver Driver, out interface{}, fieldName string) error {
	_, err := archive(driver, out, fieldName)
	return err
}

// ArchiveWithQueries archives the given instance and returns performed queries.
func ArchiveWithQueries(driver Driver, out interface{}, fieldName string) (Queries, error) {
	return archive(driver, out, fieldName)
}

func remove(driver Driver, out interface{}) (Queries, error) {
	schema, err := GetSchema(out)
	if err != nil {
		return nil, err
	}

	pkField := schema.PrimaryKeyField

	pk, err := GetFieldValueInt64(out, pkField.FieldName)
	if err != nil {
		return nil, err
	}

	if pk == int64(0) {
		return nil, fmt.Errorf("%v cannot be deleted (no primary key)", out)
	}

	query := fmt.Sprintf(
		"DELETE FROM %s WHERE %s = :%s",
		schema.TableName,
		pkField.ColumnPath(),
		pkField.ColumnName)

	params := map[string]interface{}{
		pkField.ColumnName: pk,
	}

	queries := Queries{{
		Query:  query,
		Params: params,
	}}

	_, err = driver.NamedExec(query, params)
	if err != nil {
		return queries, err
	}

	return queries, nil
}

func archive(driver Driver, out interface{}, fieldName string) (Queries, error) {
	schema, err := GetSchema(out)
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
		return nil, fmt.Errorf("%v cannot be archived (no primary key)", out)
	}

	query := fmt.Sprintf(
		"UPDATE %s SET %s = :%s WHERE %s = :%s",
		schema.TableName,
		deletedAtField.ColumnName,
		deletedAtField.ColumnName,
		pkField.ColumnPath(),
		pkField.ColumnName)

	params := map[string]interface{}{
		deletedAtField.ColumnName: now,
		pkField.ColumnName:        pk,
	}

	queries := Queries{{
		Query:  query,
		Params: params,
	}}

	_, err = driver.NamedExec(query, params)
	if err != nil {
		return queries, err
	}

	return queries, nil
}
