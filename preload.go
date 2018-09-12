package sqlxx

import (
	"context"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"

	"github.com/ulule/sqlxx/reflectx"
)

// Preload preloads related fields.
func Preload(ctx context.Context, driver Driver, out interface{}, paths ...string) error {
	_, err := preload(ctx, driver, out, paths...)
	if err != nil {
		return errors.Wrap(err, "sqlxx: cannot execute preload")
	}
	return nil
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
		return nil, errors.WithStack(ErrInvalidDriver)
	}

	if !reflectx.IsPointer(dest) {
		return nil, errors.Wrapf(ErrPointerRequired, "cannot preload %T", dest)
	}

	groups := getPreloadGroupPath(paths)
	if len(groups) == 0 {
		return nil, nil
	}

	queries := Queries{}

	for i, paths := range groups {
		if i == 0 {

			// Execute a preload of first level.
			q, err := executePreloadHandler(ctx, driver, dest, paths)
			if err != nil {
				return nil, err
			}
			queries = append(queries, q...)

		} else {

			// Otherwise, execute a preload with a walker for other levels.
			q, err := executePreloadWalker(ctx, driver, dest, paths)
			if err != nil {
				return nil, err
			}
			queries = append(queries, q...)

		}
	}

	return queries, nil
}

type preloadGroupPath map[string]preloadPath

type preloadPath []string

func (p preloadPath) Level() int {
	return len(p)
}

func (p preloadPath) Path() string {
	return strings.Join(p, ".")
}

func (p preloadPath) Parent() string {
	l := len(p) - 1
	return strings.Join(p[0:l], ".")
}

func (p preloadPath) Name() string {
	l := len(p) - 1
	if l < 0 {
		return ""
	}
	return p[l]
}

func getPreloadGroupPath(paths []string) []preloadGroupPath {
	groups := []preloadGroupPath{}

	for i := range paths {

		// Clean up preload path and create a slice of it's level.
		// "User.Profile.Avatar" -> ["User", "Profile", "Avatar"]
		path := strings.Trim(paths[i], ".")
		list := strings.Split(path, ".")

		// Increase levels slice length if required.
		for len(list) > len(groups) {
			groups = append(groups, preloadGroupPath{})
		}

		// Fill preload path to ensure we have a sequential walkthrough.
		// [0] -> "User"
		// [1] -> "User.Profile"
		// [2] -> "User.Profile.Avatar"
		for i := 0; i < len(list); i++ {
			n := i + 1
			path := preloadPath(list[0:n])
			groups[i][path.Path()] = path
		}
	}

	return groups
}

// getPreloadHandler returns a preload handler for first level.
func getPreloadHandler(ctx context.Context, driver Driver, dest interface{}) (*preloadHandler, error) {
	getModel := func(dest interface{}) (Model, error) {
		model, ok := reflectx.GetFlattenValue(dest).(Model)
		if !ok {
			return nil, errors.Wrap(ErrPreloadInvalidSchema, "a model is required")
		}
		return model, nil
	}

	if reflectx.IsSlice(dest) {
		getModel = func(dest interface{}) (Model, error) {
			model, ok := reflectx.NewSliceValue(dest).(Model)
			if !ok {
				return nil, errors.Wrap(ErrPreloadInvalidSchema, "a model is required")
			}
			return model, nil
		}
	}

	model, err := getModel(dest)
	if err != nil {
		return nil, err
	}

	schema, err := GetSchema(driver, model)
	if err != nil {
		return nil, err
	}

	handler := &preloadHandler{
		ctx:    ctx,
		driver: driver,
		model:  model,
		schema: schema,
		dest:   dest,
	}

	return handler, nil
}

// executePreloadHandler will executes a preload on first level.
// If you need to execute a preload on the second level (and/or after),
// please use executePreloadWalker instead.
func executePreloadHandler(ctx context.Context, driver Driver,
	dest interface{}, paths map[string]preloadPath) (queries Queries, err error) {

	handler, err := getPreloadHandler(ctx, driver, dest)
	if err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	associations := handler.schema.associations

	wg.Add(len(paths))
	for _, path := range paths {

		if path.Level() != 1 {
			return nil, errors.Wrapf(ErrPreloadInvalidPath, "cannot execute preload of '%s'", path.Path())
		}

		reference, ok := associations[path.Name()]
		if !ok {
			return nil, errors.Wrapf(ErrPreloadInvalidPath, "'%s' is not a valid association", path.Path())
		}

		go func(path preloadPath) {

			perr := handler.preload(reference)
			if perr != nil {
				mutex.Lock()
				defer mutex.Unlock()
				if err != nil {
					err = perr
				}
			}

			wg.Done()

		}(path)
	}

	wg.Wait()

	// TODO Handle queries...
	return queries, err
}

// executePreloadWalker will executes a preload from the second level and beyond...
func executePreloadWalker(ctx context.Context, driver Driver,
	dest interface{}, paths map[string]preloadPath) (queries Queries, err error) {

	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	wg.Add(len(paths))
	for _, path := range paths {
		go func(path preloadPath) {

			walker := reflectx.NewWalker(dest)
			defer walker.Close()

			werr := walker.Find(path.Parent(), func(values interface{}) error {
				q, perr := preload(ctx, driver, values, path.Name())
				if perr != nil {
					return perr
				}

				mutex.Lock()
				defer mutex.Unlock()

				queries = append(queries, q...)
				return nil
			})

			if werr != nil {
				mutex.Lock()
				defer mutex.Unlock()
				if err != nil {
					err = werr
				}
			}

			wg.Done()

		}(path)
	}

	wg.Wait()

	return queries, err
}

type preloadHandler struct {
	ctx    context.Context
	driver Driver
	model  Model
	schema *Schema
	dest   interface{}
}

func (handler *preloadHandler) preload(reference Reference) error {
	if reference.IsAssociationType(AssociationTypeOne) {
		return handler.preloadOne(reference)
	}
	return handler.preloadMany(reference)
}

func (handler *preloadHandler) preloadOne(reference Reference) error {
	if reference.IsLocal() {
		return handler.preloadOneLocal(reference)
	}
	return handler.preloadOneRemote(reference)
}

func (handler *preloadHandler) preloadOneLocal(reference Reference) error {
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
		preloader := reflectx.NewStringPreloader(reference.FieldName(), reference.Type(), handler.dest)
		return handler.preloadString(preloader, reference, builder,
			getPreloadForEachCallbackLocalString(preloader, reference))

	case FKIntegerType:
		preloader := reflectx.NewIntegerPreloader(reference.FieldName(), reference.Type(), handler.dest)
		return handler.preloadInteger(preloader, reference, builder,
			getPreloadForEachCallbackLocalInteger(preloader, reference))

	case FKOptionalStringType:
		preloader := reflectx.NewStringPreloader(reference.FieldName(), reference.Type(), handler.dest)
		return handler.preloadString(preloader, reference, builder,
			getPreloadForEachCallbackLocalOptionalString(preloader, reference))

	case FKOptionalIntegerType:
		preloader := reflectx.NewIntegerPreloader(reference.FieldName(), reference.Type(), handler.dest)
		return handler.preloadInteger(preloader, reference, builder,
			getPreloadForEachCallbackLocalOptionalInteger(preloader, reference))

	default:
		return errors.Errorf("'%s' is a unsupported foreign key type for preload", reference.Type())
	}
}

func (handler *preloadHandler) preloadOneRemote(reference Reference) error {
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
		preloader := reflectx.NewStringPreloader(reference.FieldName(), reference.Type(), handler.dest)
		return handler.preloadString(preloader, reference, builder,
			getPreloadForEachCallbackRemoteString(preloader, reference))

	case PKIntegerType:
		preloader := reflectx.NewIntegerPreloader(reference.FieldName(), reference.Type(), handler.dest)
		return handler.preloadInteger(preloader, reference, builder,
			getPreloadForEachCallbackRemoteInteger(preloader, reference))

	default:
		return errors.Errorf("'%s' is a unsupported primary key type for preload", reference.Type())
	}
}

func (handler *preloadHandler) preloadMany(reference Reference) error {
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
		preloader := reflectx.NewStringPreloader(reference.FieldName(), reference.Type(), handler.dest)
		return handler.preloadString(preloader, reference, builder,
			getPreloadForEachCallbackRemoteString(preloader, reference))

	case PKIntegerType:
		preloader := reflectx.NewIntegerPreloader(reference.FieldName(), reference.Type(), handler.dest)
		return handler.preloadInteger(preloader, reference, builder,
			getPreloadForEachCallbackRemoteInteger(preloader, reference))

	default:
		return errors.Errorf("'%s' is a unsupported primary key type for preload", reference.Type())
	}
}

func getPreloadForEachCallbackRemoteString(preloader *reflectx.StringPreloader,
	reference Reference) func(element reflectx.PreloadValue) error {

	return func(element reflectx.PreloadValue) error {
		pk, err := preloadFetchRemoteForeignKeyString(reference, element.Unwrap())
		if err != nil {
			return err
		}
		return preloader.AddIndex(pk, element)
	}
}

func getPreloadForEachCallbackLocalString(preloader *reflectx.StringPreloader,
	reference Reference) func(element reflectx.PreloadValue) error {

	return func(element reflectx.PreloadValue) error {
		pk, err := preloadFetchLocalForeignKeyString(reference, element.Unwrap())
		if err != nil {
			return err
		}
		return preloader.AddIndex(pk, element)
	}
}

func getPreloadForEachCallbackLocalOptionalString(preloader *reflectx.StringPreloader,
	reference Reference) func(element reflectx.PreloadValue) error {

	return func(element reflectx.PreloadValue) error {
		pk, err := preloadFetchLocalForeignKeyOptionalString(reference, element.Unwrap())
		if err != nil {
			return err
		}
		if pk == "" {
			return nil
		}
		return preloader.AddIndex(pk, element)
	}
}

func getPreloadForEachCallbackRemoteInteger(preloader *reflectx.IntegerPreloader,
	reference Reference) func(element reflectx.PreloadValue) error {

	return func(element reflectx.PreloadValue) error {
		pk, err := preloadFetchRemoteForeignKeyInteger(reference, element.Unwrap())
		if err != nil {
			return err
		}
		return preloader.AddIndex(pk, element)
	}
}

func getPreloadForEachCallbackLocalInteger(preloader *reflectx.IntegerPreloader,
	reference Reference) func(element reflectx.PreloadValue) error {

	return func(element reflectx.PreloadValue) error {
		pk, err := preloadFetchLocalForeignKeyInteger(reference, element.Unwrap())
		if err != nil {
			return err
		}
		return preloader.AddIndex(pk, element)
	}
}

func getPreloadForEachCallbackLocalOptionalInteger(preloader *reflectx.IntegerPreloader,
	reference Reference) func(element reflectx.PreloadValue) error {

	return func(element reflectx.PreloadValue) error {
		pk, err := preloadFetchLocalForeignKeyOptionalInteger(reference, element.Unwrap())
		if err != nil {
			return err
		}
		if pk == 0 {
			return nil
		}
		return preloader.AddIndex(pk, element)
	}
}

func (handler *preloadHandler) preloadString(preloader *reflectx.StringPreloader, reference Reference,
	builder builder.Select, preloadCallback func(element reflectx.PreloadValue) error) error {

	remote := reference.Remote()

	err := preloader.ForEach(preloadCallback)
	if err != nil {
		return err
	}

	list := preloader.Indexes()
	if len(list) == 0 {
		return nil
	}

	builder = builder.Where(loukoum.Condition(remote.ColumnPath()).In(list))

	err = preloader.OnExecute(func(relation interface{}) error {
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

		return preloader.UpdateValueOnIndex(fk, element)
	})
	if err != nil {
		return err
	}

	return nil
}

func (handler *preloadHandler) preloadInteger(preloader *reflectx.IntegerPreloader,
	reference Reference, builder builder.Select, preloadCallback func(element reflectx.PreloadValue) error) error {

	remote := reference.Remote()

	err := preloader.ForEach(preloadCallback)
	if err != nil {
		return err
	}

	list := preloader.Indexes()
	if len(list) == 0 {
		return nil
	}

	builder = builder.Where(loukoum.Condition(remote.ColumnPath()).In(list))

	err = preloader.OnExecute(func(relation interface{}) error {
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

		return preloader.UpdateValueOnIndex(fk, element)
	})
	if err != nil {
		return err
	}

	return nil
}

func preloadFetchLocalForeignKeyString(reference Reference, value interface{}) (string, error) {

	local := reference.Local()

	if local.ForeignKeyType() != FKStringType {
		return "", errors.Errorf("invalid type: %s", local.ForeignKeyType())
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

func preloadFetchLocalForeignKeyInteger(reference Reference, value interface{}) (int64, error) {

	local := reference.Local()

	if local.ForeignKeyType() != FKIntegerType {
		return 0, errors.Errorf("invalid type: %s", local.ForeignKeyType())
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

func preloadFetchLocalForeignKeyOptionalString(reference Reference,
	value interface{}) (string, error) {

	local := reference.Local()

	if local.ForeignKeyType() != FKOptionalStringType {
		return "", errors.Errorf("invalid type: %s", local.ForeignKeyType())
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

func preloadFetchLocalForeignKeyOptionalInteger(reference Reference, value interface{}) (int64, error) {

	local := reference.Local()

	if local.ForeignKeyType() != FKOptionalIntegerType {
		return 0, errors.Errorf("invalid type: %s", local.ForeignKeyType())
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

func preloadFetchRemoteForeignKeyString(reference Reference, value interface{}) (string, error) {

	local := reference.Local()

	if local.PrimaryKeyType() != PKStringType {
		return "", errors.Errorf("invalid type: %s", local.ForeignKeyType())
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

func preloadFetchRemoteForeignKeyInteger(reference Reference,
	value interface{}) (int64, error) {

	local := reference.Local()

	if local.PrimaryKeyType() != PKIntegerType {
		return 0, errors.Errorf("invalid type: %s", local.ForeignKeyType())
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
