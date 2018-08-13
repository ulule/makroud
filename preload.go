package sqlxx

import (
	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"

	"github.com/ulule/sqlxx/reflectx"
)

// TODO:
//
//     |--------------|--------------|-----------------|----------------|
//     |    Source    |    Action    |    Reference    |     Status     |
//     |--------------|--------------|-----------------|----------------|
//     |      1       |      ->      |        1        |       Ok       |
//     |      1       |      <-      |        1        |       Ok       |
//     |      1?      |      ->      |        1        |       Ok       |
//     |      1       |      <-      |        1?       |       Ok       |
//     |      1       |      ->      |        N        |                |
//     |      1       |      <-      |        N        |                |
//     |      N       |      ->      |        1        |                |
//     |      N       |      <-      |        1        |                |
//     |      N       |      ->      |        N        |                |
//     |      N       |      <-      |        N        |                |
//     |--------------|--------------|-----------------|----------------|

// Preload preloads related fields.
func Preload(driver Driver, out interface{}, paths ...string) error {
	_, err := PreloadWithQueries(driver, out, paths...)
	return err
}

// PreloadWithQueries preloads related fields and returns performed queries.
func PreloadWithQueries(driver Driver, out interface{}, paths ...string) (Queries, error) {
	queries, err := preload(driver, out, paths...)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute preload")
	}
	return queries, nil
}

// preload preloads related fields.
func preload(driver Driver, dest interface{}, paths ...string) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	if !reflectx.IsPointer(dest) {
		return nil, errors.Wrapf(ErrPointerRequired, "cannot preload %T", dest)
	}

	if reflectx.IsSlice(dest) {
		panic("TODO")
	}

	return preloadOne(driver, dest, paths)
}

// preloadOne preload a single instance.
func preloadOne(driver Driver, dest interface{}, paths []string) (Queries, error) {
	model, ok := dest.(Model)
	if !ok {
		return nil, errors.Wrap(ErrPreloadInvalidSchema, "a model is required")
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		reference, ok := schema.associations[path]
		if !ok {
			return nil, errors.Errorf("'%s' is not a valid association", path)
		}

		if reference.IsAssociationType(AssociationTypeOne) {
			if reference.IsLocal() {

				remote := reference.Remote()
				local := reference.Local()

				err := preloadCheckLocalForeignKey(reference, local, remote)
				if err != nil {
					return nil, err
				}

				builder := loukoum.Select(remote.Columns()).From(remote.TableName()).Limit(1)
				if remote.HasDeletedKey() {
					builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
				}

				builder, preload, err := preloadAddLocalForeignKey(reference, local, remote, builder, dest)
				if err != nil {
					return nil, err
				}
				if !preload {
					return nil, nil
				}

				relation := reflectx.MakePointer(remote.Model())

				err = Exec(driver, builder, relation)
				if err != nil && !IsErrNoRows(err) {
					return nil, err
				}
				if IsErrNoRows(err) {
					err = errors.Wrapf(ErrPreloadInvalidModel, "foreign key '%s' has an invalid value: %s",
						local.FieldName(), err.Error())
					return nil, err
				}

				err = reflectx.UpdateFieldValue(model, reference.FieldName(), relation)
				if err != nil {
					return nil, err
				}

			} else {

				remote := reference.Remote()
				local := reference.Local()

				err := preloadCheckRemoteForeignKey(reference, local, remote)
				if err != nil {
					return nil, err
				}

				builder := loukoum.Select(remote.Columns()).From(remote.TableName()).Limit(1)
				if remote.HasDeletedKey() {
					builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
				}

				builder, preload, err := preloadAddRemoteForeignKey(reference, local, remote, builder, dest)
				if err != nil {
					return nil, err
				}
				if !preload {
					return nil, nil
				}

				relation := reflectx.MakePointer(remote.Model())

				err = Exec(driver, builder, relation)
				if err != nil && !IsErrNoRows(err) {
					return nil, err
				}
				if IsErrNoRows(err) {
					return nil, nil
				}

				err = reflectx.UpdateFieldValue(model, reference.FieldName(), relation)
				if err != nil {
					return nil, err
				}

			}
		} else {
			panic("TODO")
		}
	}

	// TODO
	return nil, nil
}

// func preloadSlice(driver Driver, dest interface{}) (Queries, error) {
// 	model, ok := reflectx.NewSliceValue(dest).(Model)
// 	if !ok {
// 		// TODO Better error handling
// 		panic("A slice of model is required")
// 	}
// }

func preloadCheckLocalForeignKey(reference Reference, local ReferenceObject, remote ReferenceObject) error {
	if !local.IsForeignKey() {
		return errors.Wrapf(ErrPreloadInvalidSchema,
			"association must have a local foreign key for: '%s'", reference.Type())
	}
	if !remote.IsPrimaryKey() {
		return errors.Wrapf(ErrPreloadInvalidSchema,
			"association must have a remote primary key for: '%s'", reference.Type())
	}
	if !local.ForeignKeyType().IsCompatible(remote.PrimaryKeyType()) {
		return errors.Wrapf(ErrPreloadInvalidSchema,
			"association must have a compatible primary key and foreign key for: '%s'", reference.Type())
	}
	return nil
}

func preloadCheckRemoteForeignKey(reference Reference, local ReferenceObject, remote ReferenceObject) error {
	if !local.IsPrimaryKey() {
		return errors.Wrapf(ErrPreloadInvalidSchema,
			"association must have a local primary key for: '%s'", reference.Type())
	}
	if !remote.IsForeignKey() {
		return errors.Wrapf(ErrPreloadInvalidSchema,
			"association must have a remote foreign key for: '%s'", reference.Type())
	}
	if !local.PrimaryKeyType().IsCompatible(remote.ForeignKeyType()) {
		return errors.Wrapf(ErrPreloadInvalidSchema,
			"association must have a compatible primary key and foreign key for: '%s'", reference.Type())
	}
	return nil
}

func preloadAddLocalForeignKey(reference Reference, local ReferenceObject, remote ReferenceObject,
	builder builder.Select, dest interface{}) (builder.Select, bool, error) {

	switch local.ForeignKeyType() {
	case FKStringType:

		fk, err := reflectx.GetFieldValueString(dest, local.FieldName())
		if err != nil {
			return builder, false, err
		}
		if fk == "" {
			return builder, false, errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))
		return builder, true, nil

	case FKIntegerType:

		fk, err := reflectx.GetFieldValueInt64(dest, local.FieldName())
		if err != nil {
			return builder, false, err
		}
		if fk == 0 {
			return builder, false, errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))
		return builder, true, nil

	case FKOptionalIntegerType:

		fk, ok, err := reflectx.GetFieldOptionalValueInt64(dest, local.FieldName())
		if err != nil {
			return builder, false, err
		}
		if !ok {
			return builder, false, nil
		}

		builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))
		return builder, true, nil

	case FKOptionalStringType:

		fk, ok, err := reflectx.GetFieldOptionalValueString(dest, local.FieldName())
		if err != nil {
			return builder, false, err
		}
		if !ok {
			return builder, false, nil
		}

		builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))
		return builder, true, nil

	default:
		return builder, false, errors.Errorf("'%s' is a unsupported foreign type for preload", reference.Type())
	}
}

func preloadAddRemoteForeignKey(reference Reference, local ReferenceObject, remote ReferenceObject,
	builder builder.Select, dest interface{}) (builder.Select, bool, error) {

	switch local.PrimaryKeyType() {
	case PKStringType:

		fk, err := reflectx.GetFieldValueString(dest, local.FieldName())
		if err != nil {
			return builder, false, err
		}
		if fk == "" {
			return builder, false, errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))
		return builder, true, nil

	case PKIntegerType:

		fk, err := reflectx.GetFieldValueInt64(dest, local.FieldName())
		if err != nil {
			return builder, false, err
		}
		if fk == 0 {
			return builder, false, errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))
		return builder, true, nil

	default:
		return builder, false, errors.Errorf("'%s' is a unsupported foreign type for preload", reference.Type())
	}
}
