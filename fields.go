package sqlxx

import (
	"fmt"
	"strings"

	"reflect"

	"github.com/serenize/snaker"
	"github.com/ulule/sqlxx/reflekt"
)

// Field is a field.
type Field struct {
	// The reflect.StructField instance
	StructField reflect.StructField
	// The reflect Type of the field
	Type reflect.Type
	// The field struct tags
	Tags reflekt.FieldTags

	// Model is the zero-valued field's model used to generate schema from.
	Model Model
	// Model name that contains this field
	ModelName string
	// Table name of the model that contains this field
	TableName string

	// The field name
	Name string
	// The database column name
	ColumnName string

	// Does this field is excluded? (anonymous, private, non-sql...)
	IsExcluded bool
	// Does this field is a primary key?
	IsPrimaryKey bool
	// Does this field is a foreign key?
	IsForeignKey bool
	// Does this field is an association?
	IsAssociation bool

	// The association struct instance
	Association *Association
}

// String returns struct instance string representation.
func (f Field) String() string {
	p := fmt.Sprintf("field(model:%s table:%s name:%s column:%s)", f.ModelName, f.TableName, f.Name, f.ColumnName)
	if f.IsAssociation {
		return fmt.Sprintf("%s -- assoc(%s)", p, f.Association)
	}
	return fmt.Sprintf(p)
}

// ColumnPath returns database full column path.
func (f Field) ColumnPath() string {
	return fmt.Sprintf("%s.%s", f.TableName, f.ColumnName)
}

// NewField returns full column name from model, field and tag.
func NewField(model Model, name string) (Field, error) {
	structField, fieldFound := reflekt.GetIndirectValue(model).Type().FieldByName(name)
	if !fieldFound {
		return Field{}, fmt.Errorf("field '%s' not found in model", name)
	}

	var (
		tags                = reflekt.GetFieldTags(structField, SupportedTags, TagsMapping)
		columnName          = snaker.CamelToSnake(name)
		columnNameOverrided = false
	)

	//
	// Custom column name
	//

	if v := tags.GetByKey(SQLXStructTagName, StructTagSQLXField); v != "" {
		columnName = v
		columnNameOverrided = true
	}
	//
	// Excluded?
	//

	isExcluded := false
	if tag := tags.GetByKey(SQLXStructTagName, StructTagSQLXField); structField.PkgPath != "" || tag == "-" {
		isExcluded = true
	}

	//
	// Foreign key?
	//

	isFK := false
	if tags.HasKey(StructTagName, StructTagForeignKey) || (len(name) > len(PrimaryKeyFieldName) && strings.HasSuffix(name, PrimaryKeyFieldName)) {
		isFK = true
	}

	//
	// Type
	//

	fieldType := structField.Type
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	field := Field{
		StructField:   structField,
		Type:          fieldType,
		Tags:          tags,
		Name:          name,
		Model:         model,
		ModelName:     GetModelName(model),
		TableName:     model.TableName(),
		ColumnName:    columnName,
		IsExcluded:    isExcluded,
		IsPrimaryKey:  name == PrimaryKeyFieldName || len(tags.GetByKey(StructTagName, StructTagPrimaryKey)) != 0,
		IsForeignKey:  isFK,
		IsAssociation: false,
	}

	// If it's not an association, early return
	modelType := GetModelFromType(fieldType)
	if modelType == nil {
		return field, nil
	}

	//
	// Association
	//

	var associationType AssociationType

	switch field.Type.Kind() {
	case reflect.Struct:
		associationType = AssociationTypeOne
	case reflect.Slice:
		associationType = AssociationTypeMany
	default:
		associationType = AssociationTypeUndefined
	}

	if associationType == AssociationTypeUndefined {
		return field, fmt.Errorf("unable to guess the association type for field %s", name)
	}

	field.IsAssociation = true

	var (
		assocModel     = GetModelFromType(field.Type)
		assocModelName = GetModelName(assocModel)
	)

	assocSchema, err := GetSchema(assocModel)
	if err != nil {
		return field, err
	}

	field.Association = &Association{
		Type:            associationType,
		Schema:          assocSchema,
		Model:           assocModel,
		ModelName:       assocModelName,
		TableName:       assocModel.TableName(),
		PrimaryKeyField: assocSchema.PrimaryKeyField,
		FieldName:       assocSchema.PrimaryKeyField.Name,
		ColumnName:      assocSchema.PrimaryKeyField.ColumnName,
		FKFieldName:     fmt.Sprintf("%sID", field.Name),
		FKColumnName:    fmt.Sprintf("%s_id", snaker.CamelToSnake(field.Name)),
	}

	if !columnNameOverrided {
		field.ColumnName = fmt.Sprintf("%s_id", columnName)
	}

	if field.Association.Type == AssociationTypeMany {
		var (
			tableName      = field.TableName
			assocTableName = field.Association.TableName
		)

		field.Association.TableName = tableName
		field.TableName = assocTableName
		field.ColumnName = fmt.Sprintf("%s_%s", snaker.CamelToSnake(field.ModelName), field.Association.PrimaryKeyField.ColumnName)
		field.Association.FKFieldName = fmt.Sprintf("%s%s", field.ModelName, PrimaryKeyFieldName)
	}

	return field, nil
}
