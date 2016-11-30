package sqlxx

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/oleiade/reflections"
)

// SoftDelete soft deletes the model in the database
func SoftDelete(driver Driver, out interface{}, fieldName string) error {
	schema, err := GetSchemaFromInterface(out)
	if err != nil {
		return err
	}

	pkField := schema.PrimaryField
	pkValue, err := reflections.GetField(out, pkField.Name)

	// GO TO HELL ZERO VALUES DELETION
	if isZeroValue(pkValue) {
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

// Delete deletes the model in the database
func Delete(driver Driver, out interface{}) error {
	schema, err := GetSchemaFromInterface(out)
	if err != nil {
		return err
	}

	pkField := schema.PrimaryField
	pkValue, _ := reflections.GetField(out, pkField.Name)

	// GO TO HELL ZERO VALUES DELETION
	if isZeroValue(pkValue) {
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

// Save saves the model and populate it to the database
func Save(driver Driver, out interface{}) error {
	schema, err := GetSchemaFromInterface(out)
	if err != nil {
		return err
	}

	var (
		columns        = []string{}
		ignoredColumns = []string{}
		values         = []string{}
	)

	for _, column := range schema.Fields {
		var (
			isIgnored    bool
			hasDefault   bool
			defaultValue string
		)

		if tag, err := column.Tags.Get(StructTagName); err == nil {
			isIgnored = len(tag.Get(StructTagIgnored)) != 0
			defaultValue = tag.Get(StructTagDefault)
			hasDefault = len(defaultValue) != 0
		}

		if isIgnored || hasDefault {
			ignoredColumns = append(ignoredColumns, column.ColumnName)
		}

		if !isIgnored {
			columns = append(columns, column.ColumnName)

			if hasDefault {
				values = append(values, defaultValue)
			} else {
				values = append(values, fmt.Sprintf(":%s", column.ColumnName))
			}
		}
	}

	var query string

	pkField := schema.PrimaryField
	pkValue, _ := reflections.GetField(out, pkField.Name)

	if isZeroValue(pkValue) {
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			schema.TableName,
			strings.Join(columns, ", "),
			strings.Join(values, ", "))
	} else {
		updates := []string{}

		for i := range columns {
			updates = append(updates, fmt.Sprintf("%s = %s", columns[i], values[i]))
		}

		query = fmt.Sprintf("UPDATE %s SET %s WHERE %s = :%s",
			schema.TableName,
			strings.Join(updates, ", "),
			pkField.ColumnName,
			pkField.ColumnName)
	}

	if len(ignoredColumns) > 0 {
		query = fmt.Sprintf("%s RETURNING %s", query, strings.Join(ignoredColumns, ", "))
	}

	stmt, err := driver.PrepareNamed(query)
	if err != nil {
		return err
	}

	err = stmt.Get(out, out)
	if err != nil {
		return err
	}

	return nil
}

// Preload preloads related fields.
func Preload(driver Driver, out interface{}, fields ...string) error {
	schema, err := GetSchemaFromInterface(out)
	if err != nil {
		return err
	}

	pk, err := reflections.GetField(out, schema.PrimaryField.Name)
	if err != nil {
		return err
	}

	_, err = GetRelationQueries(schema, []interface{}{pk}, fields...)
	if err != nil {
		return err
	}

	return nil
}

// GetByParams executes a where with the given params and populates the given model.
func GetByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	return where(driver, out, params, true)
}

// FindByParams executes a where with the given params and populates the given models.
func FindByParams(driver Driver, out interface{}, params map[string]interface{}) error {
	return where(driver, out, params, false)
}

// whereQuery returns SQL where clause from model and params.
func whereQuery(model Model, params map[string]interface{}, fetchOne bool) (string, []interface{}, error) {
	schema, err := GetSchema(model)
	if err != nil {
		return "", nil, err
	}

	q := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		schema.ColumnPaths(),
		model.TableName(),
		schema.WhereColumnPaths(params))

	if fetchOne {
		q = fmt.Sprintf("%s LIMIT 1", q)
	}

	query, args, err := sqlx.Named(q, params)
	if err != nil {
		return "", nil, err
	}

	return sqlx.In(query, args...)
}

// where executes a where clause.
func where(driver Driver, out interface{}, params map[string]interface{}, fetchOne bool) error {
	model := reflectModel(out)

	query, args, err := whereQuery(model, params, fetchOne)
	if err != nil {
		return err
	}

	if fetchOne {
		return driver.Get(out, driver.Rebind(query), args...)
	}

	return driver.Select(out, driver.Rebind(query), args...)
}
