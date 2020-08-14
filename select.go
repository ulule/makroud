package makroud

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ulule/loukoum/v3"
	"github.com/ulule/loukoum/v3/builder"
	"github.com/ulule/loukoum/v3/stmt"

	"github.com/ulule/makroud/reflectx"
)

// Select retrieves the given instance using given arguments as criteria.
// This method accepts loukoum's stmt.Order, stmt.Offet, stmt.Limit and stmt.Expression as arguments.
// For unsupported statement, they will be ignored.
func Select(ctx context.Context, driver Driver, dest interface{}, args ...interface{}) error {
	if !reflectx.IsPointer(dest) {
		return errors.Wrapf(ErrPointerRequired, "makroud: cannot execute query on %T", dest)
	}
	if reflectx.IsSlice(dest) {
		return selectRows(ctx, driver, dest, args)
	}
	return selectRow(ctx, driver, dest, args)
}

func selectRow(ctx context.Context, driver Driver, dest interface{}, args []interface{}) error {
	model, ok := reflectx.GetFlattenValue(dest).(Model)
	if !ok {
		return errors.Wrapf(ErrModelRequired, "makroud: cannot execute query on %T", dest)
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return errors.Wrapf(err, "makroud: cannot fetch schema informations on %T", dest)
	}

	columns := schema.ColumnPaths()

	query, parsed := parseSelectArgs(loukoum.Select(columns.List()).From(model.TableName()), args)
	if !parsed.hasLimit {
		query = query.Limit(1)
	}
	if !parsed.hasOrder {
		query = query.OrderBy(loukoum.Order(schema.PrimaryKeyName()))
	}
	if schema.HasDeletedKey() {
		query = query.Where(loukoum.Condition(schema.DeletedKeyName()).IsNull(true))
	}

	return Exec(ctx, driver, query, dest)
}

func selectRows(ctx context.Context, driver Driver, dest interface{}, args []interface{}) error {
	model, ok := reflectx.NewSliceValue(dest).(Model)
	if !ok {
		return errors.Wrapf(ErrModelRequired, "makroud: cannot execute query on %T", dest)
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return errors.Wrapf(err, "makroud: cannot fetch schema informations on %T", dest)
	}

	columns := schema.ColumnPaths()

	query, parsed := parseSelectArgs(loukoum.Select(columns.List()).From(model.TableName()), args)
	if !parsed.hasOrder {
		query = query.OrderBy(loukoum.Order(schema.PrimaryKeyName()))
	}
	if schema.HasDeletedKey() {
		query = query.Where(loukoum.Condition(schema.DeletedKeyName()).IsNull(true))
	}

	return Exec(ctx, driver, query, dest)
}

type parsedSelectArgs struct {
	hasLimit      bool
	hasOffset     bool
	hasOrder      bool
	hasExpression bool
}

func parseSelectArgs(query builder.Select, args []interface{}) (builder.Select, parsedSelectArgs) {
	result := parsedSelectArgs{}
	for i := range args {
		switch v := args[i].(type) {
		case stmt.Limit:
			result.hasLimit = true
			query = query.Limit(v)
		case stmt.Offset:
			result.hasOffset = true
			query = query.Offset(v)
		case stmt.Order:
			result.hasOrder = true
			query = query.OrderBy(v)
		case stmt.Expression:
			result.hasExpression = true
			query = query.Where(v)
		}
	}
	return query, result
}
