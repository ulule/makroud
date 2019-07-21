package makroud

import (
	"context"
	"database/sql"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum/v3/builder"

	"github.com/ulule/makroud/reflectx"
)

// Exec will execute given query from a Loukoum builder.
// If an object is given, it will mutate it to match the row values.
func Exec(ctx context.Context, driver Driver, stmt builder.Builder, dest ...interface{}) error {
	if driver.hasLogger() {
		start := time.Now()
		query := NewQuery(stmt)

		defer func() {
			Log(ctx, driver, query, time.Since(start))
		}()
	}

	query, args := stmt.Query()

	err := exec(ctx, driver, query, args, dest...)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot execute query")
	}

	return nil
}

// RawExec will execute given query.
// If an object is given, it will mutate it to match the row values.
func RawExec(ctx context.Context, driver Driver, query string, dest ...interface{}) error {
	if driver.hasLogger() {
		start := time.Now()
		query := NewRawQuery(query)

		defer func() {
			Log(ctx, driver, query, time.Since(start))
		}()
	}

	err := exec(ctx, driver, query, nil, dest...)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot execute query")
	}

	return nil
}

// RawExecArgs will execute given query with given arguments.
// If an object is given, it will mutate it to match the row values.
func RawExecArgs(ctx context.Context, driver Driver, query string, args []interface{}, dest ...interface{}) error {
	if driver.hasLogger() {
		start := time.Now()
		query := Query{
			Raw:   query,
			Query: query,
			Args:  args,
		}

		defer func() {
			Log(ctx, driver, query, time.Since(start))
		}()
	}

	err := exec(ctx, driver, query, args, dest...)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot execute query")
	}

	return nil
}

func exec(ctx context.Context, driver Driver, query string, args []interface{}, dest ...interface{}) error {
	stmt, err := driver.Prepare(ctx, query)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot prepare statement")
	}
	defer driver.close(stmt, map[string]string{
		"query": query,
	})

	if len(dest) > 0 {
		if !reflectx.IsPointer(dest[0]) {
			return errors.Wrapf(ErrPointerRequired, "cannot execute query on %T", dest)
		}
		if reflectx.IsSlice(dest[0]) {
			return execRows(ctx, driver, stmt, args, dest[0])
		}
		return execRow(ctx, driver, stmt, args, dest[0])
	}

	return stmt.Exec(ctx, args)
}

func execRowsOnOthers(ctx context.Context, driver Driver,
	stmt Statement, args []interface{}, dest interface{}) error {

	element := reflectx.GetIndirectSliceType(dest)

	scannable := reflectx.IsScannable(element)
	if scannable {
		return execRowsOnScannable(ctx, driver, stmt, args, dest)
	}

	schemaless, err := GetSchemaless(driver, element)
	if err != nil {
		return err
	}

	rows, err := stmt.QueryRows(ctx, args)
	if err != nil {
		return err
	}

	base := reflectx.GetIndirectSliceType(dest)
	list := reflectx.GetIndirectValue(dest)

	for rows.Next() {
		val := reflectx.NewValue(base)

		err := schemaless.ScanRows(rows, val)
		if err != nil {
			return err
		}

		reflectx.AppendReflectSlice(list, val)
	}

	return nil
}

func execRowsOnScannable(ctx context.Context, driver Driver,
	stmt Statement, args []interface{}, dest interface{}) error {

	rows, err := stmt.QueryRows(ctx, args)
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	base := reflectx.GetIndirectSliceType(dest)
	list := reflectx.GetIndirectValue(dest)

	if len(columns) > 1 {
		return errors.Wrapf(ErrSliceOfScalarMultipleColumns,
			"cannot exec rows on slice of type %s with %d columns", reflectx.GetSliceType(dest), len(columns))
	}

	for rows.Next() {
		val := reflectx.NewValue(base)
		err = rows.Scan(val)
		if err != nil {
			return err
		}

		reflectx.AppendReflectSlice(list, val)
	}

	return nil
}

func execRows(ctx context.Context, driver Driver, stmt Statement, args []interface{}, dest interface{}) error {
	model, ok := reflectx.NewSliceValue(dest).(Model)
	if !ok {
		return execRowsOnOthers(ctx, driver, stmt, args, dest)
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return err
	}

	rows, err := stmt.QueryRows(ctx, args)
	if err != nil {
		return err
	}

	hasNext := false
	list := reflectx.NewReflectSlice(reflectx.GetSliceType(dest))

	for rows.Next() {
		hasNext = true
		model := reflectx.NewSliceValue(dest).(Model)

		err := schema.ScanRows(rows, model)
		if err != nil {
			return err
		}

		reflectx.AppendReflectSlice(list, model)
	}

	if hasNext {
		reflectx.SetReflectSlice(dest, list)
	}

	return nil
}

func execRow(ctx context.Context, driver Driver, stmt Statement, args []interface{}, dest interface{}) error {
	model, ok := reflectx.GetFlattenValue(dest).(Model)
	if !ok {
		return stmt.FindOne(ctx, dest, args)
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return err
	}

	row, err := stmt.QueryRow(ctx, args)
	if err != nil {
		return err
	}

	return schema.ScanRow(row, model)
}

// Count will execute given query to return a number from a aggregate function.
func Count(ctx context.Context, driver Driver, stmt builder.Builder) (int64, error) {
	count := int64(-1)

	err := Exec(ctx, driver, stmt, &count)
	if IsErrNoRows(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

// FloatCount will execute given query to return a number (in float) from a aggregate function.
func FloatCount(ctx context.Context, driver Driver, stmt builder.Builder) (float64, error) {
	count := float64(0)

	err := Exec(ctx, driver, stmt, &count)
	if IsErrNoRows(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

// IsErrNoRows returns if given error is a "no rows" error.
func IsErrNoRows(err error) bool {
	if err == nil {
		return false
	}
	err = errors.Cause(err)
	return err == sql.ErrNoRows || err == ErrNoRows
}

// toModel converts the given type to a Model instance.
func toModel(value reflect.Type) Model {
	if value.Kind() == reflect.Slice {
		value = reflectx.GetIndirectType(value.Elem())
	} else {
		value = reflectx.GetIndirectType(value)
	}

	model, ok := reflect.New(value).Elem().Interface().(Model)
	if ok {
		return model
	}

	return nil
}
