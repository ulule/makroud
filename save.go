package makroud

import (
	"context"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"

	"github.com/ulule/makroud/reflectx"
)

// Save saves the given instance.
func Save(ctx context.Context, driver Driver, model Model) error {
	err := save(ctx, driver, model)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot execute save")
	}
	return nil
}

func save(ctx context.Context, driver Driver, model Model) error {
	if driver == nil {
		return errors.WithStack(ErrInvalidDriver)
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return err
	}

	values := loukoum.Map{}
	returning := []string{}

	pk := schema.PrimaryKey()
	id, hasPK := pk.ValueOpt(model)

	err = generateSaveQuery(schema, model, hasPK, &returning, &values)
	if err != nil {
		return err
	}

	builder, err := getSaveBuilder(driver, schema, model, pk, hasPK, id, &returning, &values)
	if err != nil {
		return err
	}

	err = Exec(ctx, driver, builder, model)

	// Ignore no rows error if returning is empty.
	if IsErrNoRows(err) && len(returning) == 0 {
		return nil
	}

	return err
}

func generateSaveQuery(schema *Schema, model Model, hasPK bool, returning *[]string, values *loukoum.Map) error {
	for _, column := range schema.fields {
		if column.IsPrimaryKey() {
			continue
		}

		value, err := reflectx.GetFieldValue(model, column.FieldName())
		if err != nil {
			return err
		}

		if column.HasDefault() && reflectx.IsZero(value) && !hasPK {

			(*returning) = append((*returning), column.ColumnName())

		} else if column.IsUpdatedKey() && hasPK {

			(*values)[column.ColumnName()] = loukoum.Raw("NOW()")
			(*returning) = append((*returning), column.ColumnName())

		} else {

			(*values)[column.ColumnName()] = value

		}
	}

	return nil
}

func getSaveBuilder(driver Driver, schema *Schema, model Model, pk PrimaryKey,
	hasPK bool, id interface{}, returning *[]string, values *loukoum.Map) (builder.Builder, error) {

	if !hasPK {
		switch pk.Default() {
		case PrimaryKeyDBDefault:
			(*returning) = append((*returning), pk.ColumnName())

		case PrimaryKeyULIDDefault:
			ulid := GenerateULID(driver)
			mapper := map[string]interface{}{
				pk.ColumnName(): ulid,
			}
			(*values)[pk.ColumnName()] = ulid
			err := schema.WriteModel(mapper, model)
			if err != nil {
				return nil, err
			}
		}

		builder := loukoum.Insert(schema.TableName()).
			Set((*values)).
			Returning((*returning))

		return builder, nil
	}

	builder := loukoum.Update(schema.TableName()).
		Set((*values)).
		Where(loukoum.Condition(pk.ColumnName()).Equal(id)).
		Returning((*returning))

	return builder, nil
}
