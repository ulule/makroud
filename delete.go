package sqlxx

import (
	"fmt"
	"time"
)

// Delete deletes the model in the database
func Delete(driver Driver, out interface{}) (Queries, error) {
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

// SoftDelete is an alias for Archive
func SoftDelete(driver Driver, out interface{}, fieldName string) (Queries, error) {
	return Archive(driver, out, fieldName)
}

// Archive archives the model in the database
func Archive(driver Driver, out interface{}, fieldName string) (Queries, error) {
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
