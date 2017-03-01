package sqlxx

import (
	"fmt"
	"time"

	"github.com/ulule/sqlxx/reflekt"
)

// Delete deletes the model in the database
func Delete(driver Driver, out interface{}) error {
	schema, err := GetSchema(out)
	if err != nil {
		return err
	}

	pkField := schema.PrimaryKeyField
	pkValue, _ := reflekt.GetFieldValue(out, pkField.Name)

	// GO TO HELL ZERO VALUES DELETION
	if reflekt.IsZeroValue(pkValue) {
		return fmt.Errorf("%v has no primary key, cannot be deleted", out)
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s = :%s",
		schema.TableName,
		pkField.ColumnName,
		pkField.ColumnName)

	_, err = driver.NamedExec(query, out)
	if err != nil {
		return err
	}

	return nil
}

// SoftDelete is an alias for Archive
func SoftDelete(driver Driver, out interface{}, fieldName string) error {
	return Archive(driver, out, fieldName)
}

// Archive archives the model in the database
func Archive(driver Driver, out interface{}, fieldName string) error {

	schema, err := GetSchema(out)
	if err != nil {
		return err
	}

	pkField := schema.PrimaryKeyField
	pkValue, err := reflekt.GetFieldValue(out, pkField.Name)

	// GO TO HELL ZERO VALUES DELETION
	if reflekt.IsZeroValue(pkValue) {
		return fmt.Errorf("%v has no primary key, cannot be deleted", out)
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

	_, err = driver.NamedExec(query, m)
	if err != nil {
		return err
	}

	return nil

}
