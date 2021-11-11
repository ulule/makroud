package makroud

import (
	"context"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum/v3"
)

// Delete deletes the given instance.
func Delete(ctx context.Context, driver Driver, model Model) error {
	err := remove(ctx, driver, model)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot execute delete")
	}
	return nil
}

// Archive archives the given instance.
func Archive(ctx context.Context, driver Driver, model Model) error {
	err := archive(ctx, driver, model)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot execute archive")
	}
	return nil
}

func remove(ctx context.Context, driver Driver, model Model) error {
	if driver == nil {
		return errors.WithStack(ErrInvalidDriver)
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return err
	}

	pk := schema.PrimaryKey()
	id, err := pk.Value(model)
	if err != nil {
		return errors.Wrapf(err, "%T cannot be deleted", model)
	}

	q, args := loukoum.Delete(schema.TableName()).
		Where(loukoum.Condition(pk.ColumnName()).Equal(id)).
		Query()

	return driver.Exec(ctx, q, args...)
}

func archive(ctx context.Context, driver Driver, model Model) error {
	if driver == nil {
		return errors.WithStack(ErrInvalidDriver)
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return err
	}

	if !schema.HasDeletedKey() {
		return errors.Wrapf(ErrSchemaDeletedKey, "%T doesn't support archive operation", model)
	}

	pk := schema.PrimaryKey()
	id, err := pk.Value(model)
	if err != nil {
		return errors.Wrapf(err, "%T cannot be archived", model)
	}

	q, args := loukoum.Update(schema.TableName()).
		Set(loukoum.Pair(schema.DeletedKeyName(), loukoum.Raw("NOW()"))).
		Where(loukoum.Condition(pk.ColumnName()).Equal(id)).
		Returning(schema.DeletedKeyName()).
		Query()

	return driver.Exec(ctx, q, args...)
}
