package sqlxx

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
	//"github.com/ulule/loukoum/builder"

	"github.com/ulule/sqlxx/reflectx"
)

// TODO:
//
//     |--------------|--------------|-----------------|
//     |    Source    |    Action    |    Reference    |
//     |--------------|--------------|-----------------|
//     |      1       |      ->      |        1        |
//     |      1       |      <-      |        1        |
//     |    Source    |    Action    |    Reference    |
//     |    Source    |    Action    |    Reference    |
//     |    Source    |    Action    |    Reference    |

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
func preload(driver Driver, out interface{}, paths ...string) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	// if !reflect.Indirect(reflect.ValueOf(out)).CanAddr() {
	// 	return nil, errors.New("model instance must be addressable (pointer required)")
	// }

	return preloadOne(driver, out, paths)
}

func preloadOne(driver Driver, dest interface{}, paths []string) (Queries, error) {
	model, ok := dest.(Model)
	if !ok {
		// TODO Better error handling
		panic("A model is required")
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

				if !local.IsForeignKey() {
					// TODO Better error handling
					return nil, errors.Errorf("association must have a local foreign key")
				}
				if !remote.IsPrimaryKey() {
					// TODO Better error handling
					return nil, errors.Errorf("association must have a remote primary key")
				}
				if !local.ForeignKeyType().Equals(remote.PrimaryKeyType()) {
					// TODO Better error handling
					return nil, errors.Errorf("association must have a the same type between primary and foreign key")
				}

				builder := loukoum.Select(remote.Columns()).From(remote.TableName()).Limit(1)
				if remote.HasDeletedKey() {
					builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
				}

				fmt.Println(remote.ForeignKeyType())

				switch local.ForeignKeyType() {
				case ForeignKeyStringType:

					fk, err := reflectx.GetFieldValueString(dest, local.FieldName())
					if err != nil {
						return nil, err
					}
					if fk == "" {
						// TODO Better error handling
						return nil, errors.Errorf("foreign key has a zero value")
					}

					builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))

				case ForeignKeyIntegerType:

					fk, err := reflectx.GetFieldValueInt64(dest, local.FieldName())
					if err != nil {
						return nil, err
					}
					if fk == 0 {
						// TODO Better error handling
						return nil, errors.Errorf("foreign key has a zero value")
					}

					builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))

				default:
					return nil, errors.Errorf("'%s' is a unsupported foreign type for preload", reference.Type())
				}

				relation := reflectx.MakePointer(remote.Model())

				// TODO Handle no rows result
				err = Exec(driver, builder, relation)
				if err != nil {
					return nil, err
				}

				err := reflectx.UpdateFieldValue(model, reference.FieldName(), relation)
				if err != nil {
					return nil, err
				}

			} else {

				remote := reference.Remote()
				local := reference.Local()

				if !local.IsPrimaryKey() {
					// TODO Better error handling
					return nil, errors.Errorf("association must have a local primary key")
				}
				if !remote.IsForeignKey() {
					// TODO Better error handling
					return nil, errors.Errorf("association must have a remote foreign key")
				}
				if !local.PrimaryKeyType().Equals(remote.ForeignKeyType()) {
					// TODO Better error handling
					return nil, errors.Errorf("association must have a the same type between primary and foreign key")
				}

				builder := loukoum.Select(remote.Columns()).From(remote.TableName()).Limit(1)
				if remote.HasDeletedKey() {
					builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
				}

				switch local.PrimaryKeyType() {
				case PrimaryKeyStringType:

					fk, err := reflectx.GetFieldValueString(dest, local.FieldName())
					if err != nil {
						return nil, err
					}
					if fk == "" {
						// TODO Better error handling
						return nil, errors.Errorf("foreign key has a zero value")
					}

					builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))

				case PrimaryKeyIntegerType:

					fk, err := reflectx.GetFieldValueInt64(dest, local.FieldName())
					if err != nil {
						return nil, err
					}
					if fk == 0 {
						// TODO Better error handling
						return nil, errors.Errorf("foreign key has a zero value")
					}

					builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(fk))

				default:
					return nil, errors.Errorf("'%s' is a unsupported foreign type for preload", reference.Type())
				}

				relation := reflectx.MakePointer(remote.Model())

				// TODO Handle no rows result
				err = Exec(driver, builder, relation)
				if err != nil {
					return nil, err
				}

				err := reflectx.UpdateFieldValue(model, reference.FieldName(), relation)
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
