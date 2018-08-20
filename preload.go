package sqlxx

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"

	"github.com/ulule/sqlxx/reflectx"
)

import "fmt"

// TODO:
//
//     |--------------|--------------|-----------------|----------------|
//     |    Source    |    Action    |    Reference    |     Status     |
//     |--------------|--------------|-----------------|----------------|
//     |      1       |      ->      |        1        |       Ok       |
//     |      1       |      <-      |        1        |       Ok       |
//     |      1?      |      ->      |        1        |       Ok       |
//     |      1       |      <-      |        1?       |       Ok       |
//     |      1       |      ->      |        N        |       Ok       |
//     |      N       |      ->      |        1        |                |
//     |      N       |      <-      |        1        |                |
//     |      N       |      ->      |        N        |                |
//     |      N       |      <-      |        N        |                |
//     |--------------|--------------|-----------------|----------------|

// Preload preloads related fields.
func Preload(ctx context.Context, driver Driver, out interface{}, paths ...string) error {
	_, err := PreloadWithQueries(ctx, driver, out, paths...)
	return err
}

// PreloadWithQueries preloads related fields and returns performed queries.
func PreloadWithQueries(ctx context.Context, driver Driver, out interface{}, paths ...string) (Queries, error) {
	queries, err := preload(ctx, driver, out, paths...)
	if err != nil {
		return queries, errors.Wrap(err, "sqlxx: cannot execute preload")
	}
	return queries, nil
}

// preload preloads related fields.
func preload(ctx context.Context, driver Driver, dest interface{}, paths ...string) (Queries, error) {
	if driver == nil {
		return nil, ErrInvalidDriver
	}

	if !reflectx.IsPointer(dest) {
		return nil, errors.Wrapf(ErrPointerRequired, "cannot preload %T", dest)
	}

	if reflectx.IsSlice(dest) {
		return preloadMany(ctx, driver, dest, paths)
	}

	return preloadOne(ctx, driver, dest, paths)
}

// preloadOne preloads a single instance.
func preloadOne(ctx context.Context, driver Driver, dest interface{}, paths []string) (Queries, error) {
	model, ok := reflectx.GetFlattenValue(dest).(Model)
	if !ok {
		return nil, errors.Wrap(ErrPreloadInvalidSchema, "a model is required")
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	// TODO Refactor using preloadHandler
	preloader := preloadOneHandler{
		ctx:    ctx,
		driver: driver,
		model:  model,
		dest:   dest,
	}

	for _, path := range paths {
		reference, ok := schema.associations[path]
		if !ok {
			return nil, errors.Errorf("'%s' is not a valid association", path)
		}

		err := preloader.preload(reference)
		if err != nil {
			return nil, err
		}
	}

	// TODO
	return nil, nil
}

// preloadMany preloads a slice of instance.
func preloadMany(ctx context.Context, driver Driver, dest interface{}, paths []string) (Queries, error) {
	model, ok := reflectx.NewSliceValue(dest).(Model)
	if !ok {
		return nil, errors.Wrap(ErrPreloadInvalidSchema, "a model is required")
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	handler := preloadHandler{
		ctx:    ctx,
		driver: driver,
		model:  model,
	}

	for _, path := range paths {
		reference, ok := schema.associations[path]
		if !ok {
			return nil, errors.Errorf("'%s' is not a valid association", path)
		}

		err := handler.preload(reflectx.NewSlicePreloader(dest), reference)
		if err != nil {
			return nil, err
		}
	}

	// TODO
	return nil, nil
}

// Common version

type preloadHandler struct {
	ctx    context.Context
	driver Driver
	model  Model
}

func (handler *preloadHandler) preload(preloader reflectx.Preloader, reference Reference) error {
	if reference.IsAssociationType(AssociationTypeOne) {
		return handler.preloadOne(preloader, reference)
	} else {
		return handler.preloadMany(preloader, reference)
	}
}

func (handler *preloadHandler) preloadOne(preloader reflectx.Preloader, reference Reference) error {
	if reference.IsLocal() {
		return handler.preloadOneLocal(preloader, reference)
	} else {
		return handler.preloadOneRemote(preloader, reference)
	}
}

func (handler *preloadHandler) preloadOneLocal(preloader reflectx.Preloader, reference Reference) error {
	remote := reference.Remote()
	local := reference.Local()

	err := preloadCheckLocalForeignKey(reference, local, remote)
	if err != nil {
		return err
	}

	builder := loukoum.Select(remote.Columns()).From(remote.TableName())
	if remote.HasDeletedKey() {
		builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
	}

	switch local.ForeignKeyType() {
	case FKStringType:
		return handler.preloadOneLocalString(preloader, reference, builder)

	case FKIntegerType:
		return handler.preloadOneLocalInteger(preloader, reference, builder)

	case FKOptionalStringType:
		return handler.preloadOneLocalOptionalString(preloader, reference, builder)

	case FKOptionalIntegerType:
		return handler.preloadOneLocalOptionalInteger(preloader, reference, builder)

	default:
		return errors.Errorf("'%s' is a unsupported foreign key type for preload", reference.Type())
	}
}

func (handler *preloadHandler) preloadOneLocalString(preloader reflectx.Preloader,
	reference Reference, builder builder.Select) error {

	remote := reference.Remote()

	err := preloader.ForEach(func(element reflect.Value) error {
		pk, err := handler.fetchLocalForeignKeyString(reference, element.Interface())
		if err != nil {
			return err
		}
		return preloader.AddStringIndex(pk, element)
	})
	if err != nil {
		return err
	}

	list := preloader.StringIndexes()
	if len(list) == 0 {
		return nil
	}

	builder = builder.Where(loukoum.Condition(remote.ColumnPath()).In(list))

	err = preloader.OnExecute(reference.Type(), func(relation interface{}) error {
		err := Exec(handler.ctx, handler.driver, builder, relation)
		if err != nil && !IsErrNoRows(err) {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = preloader.OnUpdate(func(element interface{}) error {
		fk, err := reflectx.GetFieldValueString(element, remote.FieldName())
		if err != nil {
			return err
		}
		if fk == "" {
			return errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		return preloader.UpdateValueForStringIndex(reference.FieldName(), fk, element)
	})
	if err != nil {
		return err
	}

	return nil
}

func (handler *preloadHandler) preloadOneLocalInteger(preloader reflectx.Preloader,
	reference Reference, builder builder.Select) error {

	remote := reference.Remote()

	err := preloader.ForEach(func(element reflect.Value) error {
		pk, err := handler.fetchLocalForeignKeyInteger(reference, element.Interface())
		if err != nil {
			return err
		}
		return preloader.AddInt64Index(pk, element)
	})
	if err != nil {
		return err
	}

	list := preloader.Int64Indexes()
	if len(list) == 0 {
		return nil
	}

	builder = builder.Where(loukoum.Condition(remote.ColumnPath()).In(list))

	err = preloader.OnExecute(reference.Type(), func(relation interface{}) error {
		err := Exec(handler.ctx, handler.driver, builder, relation)
		if err != nil && !IsErrNoRows(err) {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = preloader.OnUpdate(func(element interface{}) error {
		fk, err := reflectx.GetFieldValueInt64(element, remote.FieldName())
		if err != nil {
			return err
		}
		if fk == 0 {
			return errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		return preloader.UpdateValueForInt64Index(reference.FieldName(), fk, element)
	})
	if err != nil {
		return err
	}

	return nil
}

func (handler *preloadHandler) preloadOneLocalOptionalString(preloader reflectx.Preloader,
	reference Reference, builder builder.Select) error {

	remote := reference.Remote()

	err := preloader.ForEach(func(element reflect.Value) error {
		pk, err := handler.fetchLocalForeignKeyOptionalString(reference, element.Interface())
		if err != nil {
			return err
		}
		if pk == "" {
			return nil
		}
		return preloader.AddStringIndex(pk, element)
	})
	if err != nil {
		return err
	}

	list := preloader.StringIndexes()
	if len(list) == 0 {
		return nil
	}

	builder = builder.Where(loukoum.Condition(remote.ColumnPath()).In(list))

	err = preloader.OnExecute(reference.Type(), func(relation interface{}) error {
		err := Exec(handler.ctx, handler.driver, builder, relation)
		if err != nil && !IsErrNoRows(err) {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = preloader.OnUpdate(func(element interface{}) error {
		fk, err := reflectx.GetFieldValueString(element, remote.FieldName())
		if err != nil {
			return err
		}
		if fk == "" {
			return errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		return preloader.UpdateValueForStringIndex(reference.FieldName(), fk, element)
	})
	if err != nil {
		return err
	}

	return nil
}

func (handler *preloadHandler) preloadOneLocalOptionalInteger(preloader reflectx.Preloader,
	reference Reference, builder builder.Select) error {

	remote := reference.Remote()

	err := preloader.ForEach(func(element reflect.Value) error {
		pk, err := handler.fetchLocalForeignKeyOptionalInteger(reference, element.Interface())
		if err != nil {
			return err
		}
		if pk == 0 {
			return nil
		}
		return preloader.AddInt64Index(pk, element)
	})
	if err != nil {
		return err
	}

	list := preloader.Int64Indexes()
	if len(list) == 0 {
		return nil
	}

	builder = builder.Where(loukoum.Condition(remote.ColumnPath()).In(list))

	err = preloader.OnExecute(reference.Type(), func(relation interface{}) error {
		err := Exec(handler.ctx, handler.driver, builder, relation)
		if err != nil && !IsErrNoRows(err) {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = preloader.OnUpdate(func(element interface{}) error {
		fk, err := reflectx.GetFieldValueInt64(element, remote.FieldName())
		if err != nil {
			return err
		}
		if fk == 0 {
			return errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		return preloader.UpdateValueForInt64Index(reference.FieldName(), fk, element)
	})
	if err != nil {
		return err
	}

	return nil
}

func (handler *preloadHandler) preloadOneRemote(preloader reflectx.Preloader, reference Reference) error {
	remote := reference.Remote()
	local := reference.Local()

	err := preloadCheckRemoteForeignKey(reference, local, remote)
	if err != nil {
		return err
	}

	builder := loukoum.Select(remote.Columns()).From(remote.TableName())
	if remote.HasDeletedKey() {
		builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
	}

	switch local.PrimaryKeyType() {
	case PKStringType:
		return handler.preloadOneRemoteString(preloader, reference, builder)

	case PKIntegerType:
		return handler.preloadOneRemoteInteger(preloader, reference, builder)

	default:
		return errors.Errorf("'%s' is a unsupported primary key type for preload", reference.Type())
	}
}

func (handler *preloadHandler) preloadOneRemoteString(preloader reflectx.Preloader,
	reference Reference, builder builder.Select) error {

	remote := reference.Remote()

	err := preloader.ForEach(func(element reflect.Value) error {
		pk, err := handler.fetchRemoteForeignKeyString(reference, element.Interface())
		if err != nil {
			return err
		}
		return preloader.AddStringIndex(pk, element)
	})
	if err != nil {
		return err
	}

	list := preloader.StringIndexes()
	if len(list) == 0 {
		return nil
	}

	builder = builder.Where(loukoum.Condition(remote.ColumnPath()).In(list))

	err = preloader.OnExecute(reference.Type(), func(relation interface{}) error {
		err := Exec(handler.ctx, handler.driver, builder, relation)
		if err != nil && !IsErrNoRows(err) {
			return err
		}
		fmt.Printf("555-1 %+v\n", relation)
		return nil
	})
	if err != nil {
		return err
	}

	err = preloader.OnUpdate(func(element interface{}) error {
		fk, err := reflectx.GetFieldValueString(element, remote.FieldName())
		if err != nil {
			return err
		}
		if fk == "" {
			return errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		return preloader.UpdateValueForStringIndex(reference.FieldName(), fk, element)
	})
	if err != nil {
		return err
	}

	return nil
}

func (handler *preloadHandler) preloadOneRemoteInteger(preloader reflectx.Preloader,
	reference Reference, builder builder.Select) error {

	remote := reference.Remote()

	err := preloader.ForEach(func(element reflect.Value) error {
		pk, err := handler.fetchRemoteForeignKeyInteger(reference, element.Interface())
		if err != nil {
			return err
		}
		return preloader.AddInt64Index(pk, element)
	})
	if err != nil {
		return err
	}

	list := preloader.Int64Indexes()
	if len(list) == 0 {
		return nil
	}

	builder = builder.Where(loukoum.Condition(remote.ColumnPath()).In(list))

	err = preloader.OnExecute(reference.Type(), func(relation interface{}) error {
		err := Exec(handler.ctx, handler.driver, builder, relation)
		if err != nil && !IsErrNoRows(err) {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	err = preloader.OnUpdate(func(element interface{}) error {
		fk, err := reflectx.GetFieldValueInt64(element, remote.FieldName())
		if err != nil {
			return err
		}
		if fk == 0 {
			return errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		return preloader.UpdateValueForInt64Index(reference.FieldName(), fk, element)
	})
	if err != nil {
		return err
	}

	return nil
}

func (handler *preloadHandler) preloadMany(preloader reflectx.Preloader, reference Reference) error {
	remote := reference.Remote()
	local := reference.Local()

	err := preloadCheckRemoteForeignKey(reference, local, remote)
	if err != nil {
		return err
	}

	builder := loukoum.Select(remote.Columns()).From(remote.TableName())
	if remote.HasDeletedKey() {
		builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
	}

	switch local.PrimaryKeyType() {
	case PKStringType:
		return handler.preloadManyString(preloader, reference, builder)

	case PKIntegerType:
		return handler.preloadManyInteger(preloader, reference, builder)

	default:
		return errors.Errorf("'%s' is a unsupported primary key type for preload", reference.Type())
	}
}

func (handler *preloadHandler) preloadManyString(preloader reflectx.Preloader,
	reference Reference, builder builder.Select) error {

	remote := reference.Remote()

	err := preloader.ForEach(func(element reflect.Value) error {
		pk, err := handler.fetchRemoteForeignKeyString(reference, element.Interface())
		if err != nil {
			return err
		}
		return preloader.AddStringIndex(pk, element)
	})
	if err != nil {
		return err
	}

	list := preloader.StringIndexes()
	if len(list) == 0 {
		return nil
	}

	builder = builder.Where(loukoum.Condition(remote.ColumnPath()).In(list))

	err = preloader.OnExecute(reference.Type(), func(relation interface{}) error {
		err := Exec(handler.ctx, handler.driver, builder, relation)
		if err != nil && !IsErrNoRows(err) {
			return err
		}
		fmt.Printf("555-2 %+v\n", relation)
		return nil
	})
	if err != nil {
		return err
	}

	err = preloader.OnUpdate(func(element interface{}) error {
		fk, err := reflectx.GetFieldValueString(element, remote.FieldName())
		if err != nil {
			return err
		}
		if fk == "" {
			return errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
		}

		return preloader.UpdateValueForStringIndex(reference.FieldName(), fk, element)
	})
	if err != nil {
		return err
	}

	return nil
}

func (handler *preloadHandler) preloadManyInteger(preloader reflectx.Preloader,
	reference Reference, builder builder.Select) error {

	panic("TODO")
}

func (handler *preloadHandler) fetchLocalForeignKeyString(reference Reference,
	value interface{}) (string, error) {

	local := reference.Local()

	if local.ForeignKeyType() != FKStringType {
		return "", errors.Errorf("invalid type: %s", local)
	}

	pk, err := reflectx.GetFieldValueString(value, local.FieldName())
	if err != nil {
		return "", err
	}
	if pk == "" {
		return "", errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
	}

	return pk, nil
}

func (handler *preloadHandler) fetchLocalForeignKeyInteger(reference Reference,
	value interface{}) (int64, error) {

	local := reference.Local()

	if local.ForeignKeyType() != FKIntegerType {
		return 0, errors.Errorf("invalid type: %s", local)
	}

	pk, err := reflectx.GetFieldValueInt64(value, local.FieldName())
	if err != nil {
		return 0, err
	}
	if pk == 0 {
		return 0, errors.Wrap(ErrPreloadInvalidModel, "foreign key has a zero value")
	}

	return pk, nil
}

func (handler *preloadHandler) fetchLocalForeignKeyOptionalString(reference Reference,
	value interface{}) (string, error) {

	local := reference.Local()

	if local.ForeignKeyType() != FKOptionalStringType {
		return "", errors.Errorf("invalid type: %s", local)
	}

	fk, ok, err := reflectx.GetFieldOptionalValueString(value, local.FieldName())
	if err != nil {
		return "", err
	}
	if !ok || fk == "" {
		return "", nil
	}

	return fk, nil
}

func (handler *preloadHandler) fetchLocalForeignKeyOptionalInteger(reference Reference,
	value interface{}) (int64, error) {

	local := reference.Local()

	if local.ForeignKeyType() != FKOptionalIntegerType {
		return 0, errors.Errorf("invalid type: %s", local)
	}

	fk, ok, err := reflectx.GetFieldOptionalValueInt64(value, local.FieldName())
	if err != nil {
		return 0, err
	}
	if !ok || fk == 0 {
		return 0, nil
	}

	return fk, nil
}

func (handler *preloadHandler) fetchRemoteForeignKeyString(reference Reference,
	value interface{}) (string, error) {

	local := reference.Local()

	if local.PrimaryKeyType() != PKStringType {
		return "", errors.Errorf("invalid type: %s", local)
	}

	pk, err := reflectx.GetFieldValueString(value, local.FieldName())
	if err != nil {
		return "", err
	}
	if pk == "" {
		return "", errors.Wrap(ErrPreloadInvalidModel, "primary key has a zero value")
	}

	return pk, nil
}

func (handler *preloadHandler) fetchRemoteForeignKeyInteger(reference Reference,
	value interface{}) (int64, error) {

	local := reference.Local()

	if local.PrimaryKeyType() != PKIntegerType {
		return 0, errors.Errorf("invalid type: %s", local)
	}

	pk, err := reflectx.GetFieldValueInt64(value, local.FieldName())
	if err != nil {
		return 0, err
	}
	if pk == 0 {
		return 0, errors.Wrap(ErrPreloadInvalidModel, "primary key has a zero value")
	}

	return pk, nil
}

// END Common

type preloadOneHandler struct {
	ctx    context.Context
	driver Driver
	model  Model
	dest   interface{}
}

func (preloader *preloadOneHandler) preload(reference Reference) error {
	if reference.IsAssociationType(AssociationTypeOne) {
		return preloader.preloadOne(reference)
	} else {
		return preloader.preloadMany(reference)
	}
}

func (preloader *preloadOneHandler) preloadOne(reference Reference) error {
	handler := &preloadHandler{
		ctx:    preloader.ctx,
		driver: preloader.driver,
		model:  preloader.model,
	}

	return handler.preload(reflectx.NewStructPreloader(preloader.dest), reference)
}

func (preloader *preloadOneHandler) preloadMany(reference Reference) error {
	ctx := preloader.ctx
	driver := preloader.driver
	dest := preloader.dest
	model := preloader.model
	remote := reference.Remote()
	local := reference.Local()

	err := preloadCheckRemoteForeignKey(reference, local, remote)
	if err != nil {
		return err
	}

	builder := loukoum.Select(remote.Columns()).From(remote.TableName())
	if remote.HasDeletedKey() {
		builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
	}

	builder, preload, err := preloadAddRemoteForeignKey(reference, local, remote, builder, dest)
	if err != nil {
		return err
	}
	if !preload {
		return nil
	}

	list := reflectx.NewReflectSlice(reflectx.GetSliceType(reference.Type()))
	relation := reflectx.MakePointer(list.Interface())

	err = Exec(ctx, driver, builder, relation)
	if err != nil && !IsErrNoRows(err) {
		return err
	}

	return reflectx.UpdateFieldValue(model, reference.FieldName(), relation)
}

func preloadCheckLocalForeignKey(reference Reference, local ReferenceObject, remote ReferenceObject) error {
	if !reference.IsLocal() {
		return errors.Wrapf(ErrPreloadInvalidSchema,
			"association must have a local reference for: '%s'", reference.Type())
	}
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
	if reference.IsLocal() {
		return errors.Wrapf(ErrPreloadInvalidSchema,
			"association cannot have a local reference for: '%s'", reference.Type())
	}
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

func preloadAddRemoteForeignKey(reference Reference, local ReferenceObject, remote ReferenceObject,
	builder builder.Select, dest interface{}) (builder.Select, bool, error) {

	switch local.PrimaryKeyType() {
	case PKStringType:

		pk, err := reflectx.GetFieldValueString(dest, local.FieldName())
		if err != nil {
			return builder, false, err
		}
		if pk == "" {
			return builder, false, errors.Wrap(ErrPreloadInvalidModel, "primary key has a zero value")
		}

		builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(pk))
		return builder, true, nil

	case PKIntegerType:

		pk, err := reflectx.GetFieldValueInt64(dest, local.FieldName())
		if err != nil {
			return builder, false, err
		}
		if pk == 0 {
			return builder, false, errors.Wrap(ErrPreloadInvalidModel, "primary key has a zero value")
		}

		builder = builder.Where(loukoum.Condition(remote.ColumnPath()).Equal(pk))
		return builder, true, nil

	default:
		return builder, false, errors.Errorf("'%s' is a unsupported primary key type for preload", reference.Type())
	}
}
