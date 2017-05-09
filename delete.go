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

	pkValue, err := GetFieldValue(out, pkField.Name)
	if err != nil {
		return nil, err
	}

	// GO TO HELL ZERO VALUES DELETION
	if IsZeroValue(pkValue) {
		return nil, fmt.Errorf("%v has no primary key, cannot be deleted", out)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = :%s", schema.TableName, pkField.ColumnName, pkField.ColumnName)

	queries := Queries{{Query: query}}

	_, err = driver.NamedExec(query, out)
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

	pkField := schema.PrimaryKeyField

	pkValue, err := GetFieldValue(out, pkField.Name)
	if err != nil {
		return nil, err
	}
	// GO TO HELL ZERO VALUES DELETION
	if IsZeroValue(pkValue) {
		return nil, fmt.Errorf("%v has no primary key, cannot be deleted", out)
	}

	field := schema.Fields[fieldName]

	now := time.Now()

	query := fmt.Sprintf("UPDATE %s SET %s = :%s WHERE %s = :%s",
		schema.TableName,
		field.ColumnName,
		field.ColumnName,
		pkField.ColumnName,
		pkField.ColumnName)

	m := map[string]interface{}{
		field.ColumnName:   now,
		pkField.ColumnName: pkValue,
	}

	queries := Queries{{Query: query, Params: m}}

	_, err = driver.NamedExec(query, m)
	if err != nil {
		return queries, err
	}

	return queries, nil
}
