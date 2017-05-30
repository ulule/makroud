package sqlxx

import (
	"reflect"

	//"github.com/davecgh/go-spew/spew"

	"github.com/pkg/errors"
)

// required:
//   -> get pk
//     -> get pk_value
//     -> get pk_key

type XSchema struct {
	Name         string
	Table        string
	PrimaryKey   XPrimaryKey
	Fields       map[string]XField
	Associations map[string]XAssociationReference
	ForeignKeys  map[string]XForeignKey
}

func (e XSchema) FindForeignKeyWithReference(reference string) (XForeignKey, bool) {
	for _, fk := range e.ForeignKeys {
		if fk.ReferenceName == reference {
			return fk, true
		}
	}
	return XForeignKey{}, false
}

type SchemaBuilder interface {
	SetTableName(name, table string) SchemaBuilder
	SetPrimaryKey(name, column string, nature PrimaryKeyType) SchemaBuilder
	AddField(name, column string, options ...FieldOption) SchemaBuilder
	AddAssociation(name, reference string, nature AssociationType) SchemaBuilder
	Create(driver Driver, model XModel) (*XSchema, error)
}

func NewSchemaBuilder() SchemaBuilder {
	return &schemaBuilder{
		fields:       make(map[string]fieldOp),
		associations: make(map[string]associationOp),
	}
}

type schemaBuilder struct {
	hasTableName bool
	modelName    string
	tableName    string
	hasPk        bool
	pkName       string
	pkColumn     string
	pkType       PrimaryKeyType
	fields       map[string]fieldOp
	associations map[string]associationOp
}

func (e *schemaBuilder) SetTableName(name, table string) SchemaBuilder {
	e.hasTableName = true
	e.modelName = name
	e.tableName = table
	return e
}

func (e *schemaBuilder) SetPrimaryKey(name, column string, nature PrimaryKeyType) SchemaBuilder {
	e.hasPk = true
	e.pkName = name
	e.pkColumn = column
	e.pkType = nature
	return e
}

func (e *schemaBuilder) AddField(name, column string, options ...FieldOption) SchemaBuilder {
	_, ok := e.fields[name]
	if ok {
		return e
	}

	field := fieldOp{
		FieldName:  name,
		ColumnName: column,
	}

	for i := range options {
		options[i].apply(&field)
	}

	e.fields[name] = field
	return e
}

func (e *schemaBuilder) AddAssociation(name, reference string, nature AssociationType) SchemaBuilder {
	_, ok := e.associations[name]
	if ok {
		return e
	}

	association := associationOp{
		FieldName:     name,
		ReferenceName: reference,
		Type:          nature,
	}

	e.associations[name] = association
	return e
}

func (e *schemaBuilder) Create(driver Driver, model XModel) (*XSchema, error) {
	if !e.hasTableName {
		return nil, errors.New("sqlxx: please define table name on model")
	}
	if !e.hasPk {
		return nil, errors.New("sqlxx: please define primary key on model")
	}

	schema := &XSchema{
		Name:         e.modelName,
		Table:        e.tableName,
		Fields:       make(map[string]XField),
		Associations: make(map[string]XAssociationReference),
		ForeignKeys:  make(map[string]XForeignKey),
	}

	schema.Fields[e.pkName] = XField{
		ModelName:    e.modelName,
		TableName:    e.tableName,
		FieldName:    e.pkName,
		ColumnName:   e.pkColumn,
		IsPrimaryKey: true,
	}

	schema.PrimaryKey = XPrimaryKey{
		ModelName:  e.modelName,
		TableName:  e.tableName,
		FieldName:  e.pkName,
		ColumnName: e.pkColumn,
		Type:       e.pkType,
	}

	for name, field := range e.fields {
		if field.IsForeignKey {
			schema.ForeignKeys[name] = XForeignKey{
				ModelName:     e.modelName,
				TableName:     e.tableName,
				FieldName:     field.FieldName,
				ColumnName:    field.ColumnName,
				ReferenceName: field.ReferenceName,
			}
		}
		schema.Fields[name] = XField{
			ModelName:    e.modelName,
			TableName:    e.tableName,
			FieldName:    field.FieldName,
			ColumnName:   field.ColumnName,
			IsPrimaryKey: false,
			IsArchiveKey: field.IsArchiveKey,
			IsForeignKey: field.IsForeignKey,
			Default:      field.Default,
		}
	}

	for _, association := range e.associations {
		switch association.Type {
		case AssociationTypeOne:
			err := e.createHasOneAssociations(driver, model, association, schema)
			if err != nil {
				return nil, err
			}
		case AssociationTypeMany:
			err := e.createHasManyAssociations(driver, model, association, schema)
			if err != nil {
				return nil, err
			}
		default:
			panic("unsupported")
		}
	}

	return schema, nil
}

// getReferenceSchema return Schema from given reference's type.
func getReferenceSchema(driver Driver, name string, element reflect.Type) (XSchema, error) {
	zero := XSchema{}

	model, ok := GetZero(element).Interface().(XModel)
	if !ok {
		return zero, errors.Errorf("sqlxx: field '%s' require a valid sqlxx model as reference", name)
	}

	schema, err := XGetSchema(driver, model)
	if err != nil {
		return zero, errors.Wrapf(err, "sqlxx: field '%s' require a valid sqlxx model as reference", name)
	}

	return schema, nil
}

func (e *schemaBuilder) createHasOneAssociations(driver Driver, model XModel,
	association associationOp, schema *XSchema) error {

	name := association.FieldName
	element, ok := GetFieldByName(model, name)
	if !ok {
		return errors.Errorf("sqlxx: field '%s' not found in given model", name)
	}

	if element.Type.Kind() != reflect.Struct {
		return errors.Errorf("sqlxx: field '%s' should be a struct", name)
	}

	target, err := getReferenceSchema(driver, name, element.Type)
	if err != nil {
		return err
	}

	source := XField{}
	reference := XField{}
	hasSource := false
	hasReference := false

	// First, we try to obtain the foreign key from target's schema.
	fk, ok := target.FindForeignKeyWithReference(e.modelName)
	if ok {

		source, ok = target.Fields[fk.FieldName]
		if !ok {
			return errors.Errorf("sqlxx: cannot find foreign key in schema: %s", name)
		}

		reference, ok = schema.Fields[schema.PrimaryKey.FieldName]
		if !ok {
			return errors.Errorf("sqlxx: cannot find foreign key in schema: %s", name)
		}

		hasSource = true
		hasReference = true

	}

	// Unless the foreign key is defined in current schema.
	fk, ok = schema.FindForeignKeyWithReference(association.ReferenceName)
	if ok {

		source, ok = schema.Fields[fk.FieldName]
		if !ok {
			return errors.Errorf("sqlxx: cannot find foreign key in schema: %s", name)
		}

		reference, ok = target.Fields[target.PrimaryKey.FieldName]
		if !ok {
			return errors.Errorf("sqlxx: cannot find foreign key in schema: %s", name)
		}

		hasSource = true
		hasReference = true

	}

	if !hasSource || !hasReference {
		return errors.Errorf("sqlxx: cannot find foreign key in schema: %s", name)
	}

	schema.Associations[name] = XAssociationReference{
		ModelName: e.modelName,
		FieldName: association.FieldName,
		Type:      association.Type,
		Source:    source,
		Reference: reference,
	}

	return nil
}

func (e *schemaBuilder) createHasManyAssociations(driver Driver, model XModel,
	association associationOp, schema *XSchema) error {

	name := association.FieldName
	element, ok := GetFieldByName(model, name)
	if !ok {
		return errors.Errorf("sqlxx: field '%s' not found in given model", name)
	}

	if element.Type.Kind() != reflect.Slice {
		return errors.Errorf("sqlxx: field '%s' should be a slice", name)
	}

	target, err := getReferenceSchema(driver, name, element.Type.Elem())
	if err != nil {
		return err
	}

	fk, ok := target.FindForeignKeyWithReference(e.modelName)
	if !ok {
		return errors.Wrapf(err, "sqlxx: cannot find foreign key in reference: %s", name)
	}

	source, ok := target.Fields[fk.FieldName]
	if !ok {
		return errors.Wrapf(err, "sqlxx: cannot find foreign key in schema: %s", name)
	}

	reference, ok := schema.Fields[schema.PrimaryKey.FieldName]
	if !ok {
		return errors.Errorf("sqlxx: cannot find foreign key in schema: %s", name)
	}

	schema.Associations[name] = XAssociationReference{
		ModelName: e.modelName,
		FieldName: association.FieldName,
		Type:      association.Type,
		Source:    source,
		Reference: reference,
	}

	return nil
}

type associationOp struct {
	FieldName     string
	ReferenceName string
	Type          AssociationType
}

type fieldOp struct {
	// The field name
	FieldName string
	// The database column name
	ColumnName string
	// TODO
	IsArchiveKey bool
	// Default define the default statement to use if value in model is undefined.
	Default string
	// IsForeignKey define if field should behave like a foreign key.
	IsForeignKey bool
	// ReferenceName define the reference models name if field is a foreign key.
	ReferenceName string
}

// GetSchema returns the given schema from global cache
// If the given schema does not exists, returns false as bool.
func XGetSchema(driver Driver, model XModel) (XSchema, error) {
	if !driver.hasCache() {
		return XnewSchema(driver, model)
	}

	schema, found := driver.cache().XGetSchema(model)
	if found {
		return schema, nil
	}

	schema, err := XnewSchema(driver, model)
	if err != nil {
		return schema, err
	}

	driver.cache().XSetSchema(schema)

	return schema, nil
}

func XnewSchema(driver Driver, model XModel) (XSchema, error) {
	builder := NewSchemaBuilder()
	model.CreateSchema(builder)

	schema, err := builder.Create(driver, model)
	if schema == nil {
		schema = &XSchema{}
	}

	return *schema, err
}
