package makroud

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/ulule/loukoum/v3"
	"github.com/ulule/loukoum/v3/builder"

	"github.com/ulule/makroud/reflectx"
)

// Preload preloads related fields.
func Preload(ctx context.Context, driver Driver, out interface{}, handlers ...PreloadHandler) error {
	err := preload(ctx, driver, out, handlers...)
	if err != nil {
		return errors.Wrap(err, "makroud: cannot execute preload")
	}
	return nil
}

// preload preloads related fields.
func preload(ctx context.Context, driver Driver, dest interface{}, handlers ...PreloadHandler) error {
	if driver == nil {
		return errors.WithStack(ErrInvalidDriver)
	}

	if !reflectx.IsPointer(dest) && !reflectx.IsSlice(dest) {
		return errors.Wrapf(ErrPointerOrSliceRequired, "cannot preload %T", dest)
	}

	groups := getPreloadGroupOperations(handlers)
	if len(groups) == 0 {
		return nil
	}

	for i, group := range groups {
		if i == 0 {
			// Execute a preload of first level.
			err := executePreloadHandler(ctx, driver, dest, group)
			if err != nil {
				return err
			}
		} else {
			// Otherwise, execute a preload with a walker for other levels.
			err := executePreloadWalker(ctx, driver, dest, group)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// PreloadHandler defines what resources should be preloaded.
type PreloadHandler struct {
	field    string
	unscoped bool
	callback func(query builder.Select) builder.Select
}

// WithPreloadField returns a handler that preload a field.
func WithPreloadField(field string) PreloadHandler {
	return PreloadHandler{
		field:    field,
		unscoped: false,
		callback: func(query builder.Select) builder.Select {
			return query
		},
	}
}

// WithPreloadCallback returns a handler that preload a field with conditions.
func WithPreloadCallback(field string, callback func(query builder.Select) builder.Select) PreloadHandler {
	return PreloadHandler{
		field:    field,
		unscoped: false,
		callback: callback,
	}
}

// WithUnscopedPreload unscopes given preload handler.
func WithUnscopedPreload(handler PreloadHandler) PreloadHandler {
	handler.unscoped = true
	return handler
}

// preloadGroupOperation defines, for a preload level, what resources should be preloaded.
type preloadGroupOperation map[string]preloadOperation

// preloadOperation defines what and how a resource should be preloaded.
type preloadOperation struct {
	levels  []string
	handler PreloadHandler
}

// Level returns the preload level.
func (o preloadOperation) Level() int {
	return len(o.levels)
}

// Unscoped returns if preload is unscoped.
func (o preloadOperation) Unscoped() bool {
	return o.handler.unscoped
}

// Path returns the preload full path.
func (o preloadOperation) Path() string {
	return o.handler.field
}

// Callback returns the preload conditions to execute for this operation.
func (o preloadOperation) Callback() func(query builder.Select) builder.Select {
	return o.handler.callback
}

// Parent returns the resource parent path.
func (o preloadOperation) Parent() string {
	size := len(o.levels) - 1
	return strings.Join(o.levels[0:size], ".")
}

// Name returns the resource name.
func (o preloadOperation) Name() string {
	size := len(o.levels) - 1
	if size < 0 {
		return ""
	}
	return o.levels[size]
}

func getPreloadGroupOperations(handlers []PreloadHandler) []preloadGroupOperation {
	groups := []preloadGroupOperation{}

	for i := range handlers {

		// Clean up preload path and create a slice of it's level.
		// "User.Profile.Avatar" -> ["User", "Profile", "Avatar"]
		handlers[i].field = strings.Trim(handlers[i].field, ".")
		levels := strings.Split(handlers[i].field, ".")

		// Increase levels slice length if required.
		for len(levels) > len(groups) {
			groups = append(groups, preloadGroupOperation{})
		}

		// Create preload operation and attach it to the correct level.
		op := preloadOperation{
			levels:  levels,
			handler: handlers[i],
		}

		idx := op.Level() - 1
		groups[idx][op.Path()] = op

	}

	for i := range handlers {

		// Fill preload operations to ensure we have a sequential walkthrough.
		// [0] -> "User"
		// [1] -> "User.Profile"
		// [2] -> "User.Profile.Avatar"

		list := strings.Split(handlers[i].field, ".")
		for i := 0; i < len(list); i++ {
			n := i + 1
			levels := list[0:n]
			path := strings.Join(levels, ".")
			_, ok := groups[i][path]
			if !ok {
				op := preloadOperation{
					levels:  levels,
					handler: WithPreloadField(path),
				}
				groups[i][op.Path()] = op
			}
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
	dest interface{}, group preloadGroupOperation) error {

	handler, err := getPreloadHandler(ctx, driver, dest)
	if err != nil {
		return err
	}

	associations := handler.schema.associations

	for _, operation := range group {
		if operation.Level() != 1 {
			return errors.Wrapf(ErrPreloadInvalidPath, "cannot execute preload of '%s'", operation.Path())
		}

		reference, ok := associations[operation.Name()]
		if !ok {
			return errors.Wrapf(ErrPreloadInvalidPath, "'%s' is not a valid association", operation.Path())
		}

		err := handler.preload(reference, operation.Unscoped(), operation.Callback())
		if err != nil {
			return err
		}
	}

	return nil
}

// executePreloadWalker will executes a preload from the second level and beyond...
func executePreloadWalker(ctx context.Context, driver Driver,
	dest interface{}, group preloadGroupOperation) error {

	for _, operation := range group {
		walker := reflectx.NewWalker(dest)
		defer walker.Close()

		err := walker.Find(operation.Parent(), func(values interface{}) error {
			op := WithPreloadCallback(operation.Name(), operation.Callback())
			op.unscoped = operation.Unscoped()
			return preload(ctx, driver, values, op)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type preloadHandler struct {
	ctx    context.Context
	driver Driver
	model  Model
	schema *Schema
	dest   interface{}
}

func (handler *preloadHandler) preload(reference Reference, unscoped bool,
	callback func(query builder.Select) builder.Select) error {

	if reference.IsAssociationType(AssociationTypeOne) {
		return handler.preloadOne(reference, unscoped, callback)
	}
	return handler.preloadMany(reference, unscoped, callback)
}

func (handler *preloadHandler) preloadOne(reference Reference, unscoped bool,
	callback func(query builder.Select) builder.Select) error {

	if reference.IsLocal() {
		return handler.preloadOneLocal(reference, unscoped, callback)
	}
	return handler.preloadOneRemote(reference, unscoped, callback)
}

func (handler *preloadHandler) preloadOneLocal(reference Reference, unscoped bool,
	callback func(query builder.Select) builder.Select) error {

	remote := reference.Remote()
	local := reference.Local()

	err := preloadCheckLocalForeignKey(reference, local, remote)
	if err != nil {
		return err
	}

	builder := callback(loukoum.Select(remote.Columns()).From(remote.TableName()))
	if remote.HasDeletedKey() && !unscoped {
		builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
	}

	switch local.ForeignKeyType() {
	case FKStringType:

		preloader := reflectx.NewStringPreloader(reference.FieldName(), reference.Type(), handler.dest)
		defer preloader.Close()

		return handler.preloadString(preloader, reference, builder,
			getPreloadForEachCallbackLocalString(preloader, reference))

	case FKIntegerType:

		preloader := reflectx.NewIntegerPreloader(reference.FieldName(), reference.Type(), handler.dest)
		defer preloader.Close()

		return handler.preloadInteger(preloader, reference, builder,
			getPreloadForEachCallbackLocalInteger(preloader, reference))

	case FKOptionalStringType:

		preloader := reflectx.NewStringPreloader(reference.FieldName(), reference.Type(), handler.dest)
		defer preloader.Close()

		return handler.preloadString(preloader, reference, builder,
			getPreloadForEachCallbackLocalOptionalString(preloader, reference))

	case FKOptionalIntegerType:

		preloader := reflectx.NewIntegerPreloader(reference.FieldName(), reference.Type(), handler.dest)
		defer preloader.Close()

		return handler.preloadInteger(preloader, reference, builder,
			getPreloadForEachCallbackLocalOptionalInteger(preloader, reference))

	default:
		return errors.Errorf("'%s' is a unsupported foreign key type for preload", reference.Type())
	}
}

func (handler *preloadHandler) preloadOneRemote(reference Reference, unscoped bool,
	callback func(query builder.Select) builder.Select) error {

	remote := reference.Remote()
	local := reference.Local()

	err := preloadCheckRemoteForeignKey(reference, local, remote)
	if err != nil {
		return err
	}

	builder := callback(loukoum.Select(remote.Columns()).From(remote.TableName()))
	if remote.HasDeletedKey() && !unscoped {
		builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
	}

	switch local.PrimaryKeyType() {
	case PKStringType:

		preloader := reflectx.NewStringPreloader(reference.FieldName(), reference.Type(), handler.dest)
		defer preloader.Close()

		return handler.preloadString(preloader, reference, builder,
			getPreloadForEachCallbackRemoteString(preloader, reference))

	case PKIntegerType:

		preloader := reflectx.NewIntegerPreloader(reference.FieldName(), reference.Type(), handler.dest)
		defer preloader.Close()

		return handler.preloadInteger(preloader, reference, builder,
			getPreloadForEachCallbackRemoteInteger(preloader, reference))

	default:
		return errors.Errorf("'%s' is a unsupported primary key type for preload", reference.Type())
	}
}

func (handler *preloadHandler) preloadMany(reference Reference, unscoped bool,
	callback func(query builder.Select) builder.Select) error {

	remote := reference.Remote()
	local := reference.Local()

	err := preloadCheckRemoteForeignKey(reference, local, remote)
	if err != nil {
		return err
	}

	builder := callback(loukoum.Select(remote.Columns()).From(remote.TableName()))
	if remote.HasDeletedKey() && !unscoped {
		builder = builder.Where(loukoum.Condition(remote.DeletedKeyPath()).IsNull(true))
	}

	switch local.PrimaryKeyType() {
	case PKStringType:

		preloader := reflectx.NewStringPreloader(reference.FieldName(), reference.Type(), handler.dest)
		defer preloader.Close()

		return handler.preloadString(preloader, reference, builder,
			getPreloadForEachCallbackRemoteString(preloader, reference))

	case PKIntegerType:

		preloader := reflectx.NewIntegerPreloader(reference.FieldName(), reference.Type(), handler.dest)
		defer preloader.Close()

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
