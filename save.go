package makroud

import (
	"context"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"

	"github.com/ulule/makroud/reflectx"
)

// Save inserts or updates the given instance.
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

	err = generateSaveQuery(schema, model, hasPK, &returning, values)
	if err != nil {
		return err
	}

	builder, err := getSaveBuilder(driver, schema, model, pk, hasPK, id, &returning, values)
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

func generateSaveQuery(schema *Schema, model Model, hasPK bool, returning *[]string, values loukoum.Map) error {
	instance := reflectx.GetIndirectValue(model)
	for _, column := range schema.fields {
		if column.IsPrimaryKey() {
			continue
		}

		name := column.ColumnName()
		value, err := reflectx.GetFieldValueWithIndexes(instance, column.FieldIndex())
		if err != nil {
			return err
		}

		if !hasPK && column.HasDefault() && reflectx.IsZero(value) {

			(*returning) = append((*returning), name)

		} else if column.IsUpdatedKey() && hasPK {

			values[name] = loukoum.Raw("NOW()")
			(*returning) = append((*returning), name)

		} else {

			values[name] = value

		}
	}

	return nil
}

func getSaveBuilder(driver Driver, schema *Schema, model Model, pk PrimaryKey,
	hasPK bool, id interface{}, returning *[]string, values loukoum.Map) (builder.Builder, error) {

	if !hasPK {
		switch pk.Default() {
		case PrimaryKeyDBDefault:
			(*returning) = append((*returning), pk.ColumnName())

		case PrimaryKeyULIDDefault:
			ulid := GenerateULID(driver)
			values[pk.ColumnName()] = ulid
			(*returning) = append((*returning), pk.ColumnName())

		case PrimaryKeyUUIDV1Default:
			uuid := GenerateUUIDV1(driver)
			values[pk.ColumnName()] = uuid
			(*returning) = append((*returning), pk.ColumnName())

		case PrimaryKeyUUIDV4Default:
			uuid := GenerateUUIDV4(driver)
			values[pk.ColumnName()] = uuid
			(*returning) = append((*returning), pk.ColumnName())

		default:
			return nil, errors.Errorf("unsupported primary key type: %s", pk.Default())
		}

		builder := loukoum.Insert(schema.TableName()).
			Set(values).
			Returning((*returning))

		return builder, nil
	}

	builder := loukoum.Update(schema.TableName()).
		Set(values).
		Where(loukoum.Condition(pk.ColumnName()).Equal(id)).
		Returning((*returning))

	return builder, nil
}
