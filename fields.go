package sqlxx

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"github.com/serenize/snaker"

	"github.com/ulule/sqlxx/reflectx"
)

// Field is a field.
type Field struct {
	modelName     string
	tableName     string
	fieldName     string
	columnPath    string
	columnName    string
	isPrimaryKey  bool
	isForeignKey  bool
	isAssociation bool
	isExcluded    bool
	hasDefault    bool
	hasULID       bool
	isCreatedKey  bool
	isUpdatedKey  bool
	isDeletedKey  bool
	rtype         reflect.Type

	// // The reflect.StructField instance
	// StructField reflect.StructField
	// // The reflect Type of the field
	// Type reflect.Type
	// // The field struct tags
	// Tags FieldTags
	// // Schema is the related model schema
	// Schema *Schema
	// // Model is the zero-valued field's model used to generate schema from.
	// Model Model

	// // The association type
	// AssociationType AssociationType
	// // ForeignKey contains foreign key relations information
	// //ForeignKey *ForeignKey
	// // DestinationField is the value destination field if the field is an association.
	// DestinationField string
}

// ModelName define the model name of this field.
func (field Field) ModelName() string {
	return field.modelName
}

// FieldName define the struct field name used for this field.
func (field Field) FieldName() string {
	return field.fieldName
}

// TableName returns the model name's table name of this field.
func (field Field) TableName() string {
	return field.tableName
}

// ColumnPath returns the field's full column path.
func (field Field) ColumnPath() string {
	return field.columnPath
}

// ColumnName returns the field's column name.
func (field Field) ColumnName() string {
	return field.columnName
}

// IsPrimaryKey returns if the field is a primary key.
func (field Field) IsPrimaryKey() bool {
	return field.isPrimaryKey
}

// IsExcluded returns if the field excluded. (anonymous, private...)
func (field Field) IsExcluded() bool {
	return field.isExcluded
}

// IsForeignKey returns if the field is a foreign key.
func (field Field) IsForeignKey() bool {
	return field.isForeignKey
}

// IsAssociation returns if the field is an association.
func (field Field) IsAssociation() bool {
	return field.isAssociation
}

// HasDefault returns if the field has a default value and should be in returning statement.
func (field Field) HasDefault() bool {
	return field.hasDefault
}

// HasULID returns if the field has a ulid type for it's primary key.
func (field Field) HasULID() bool {
	return field.hasULID
}

// IsCreatedKey returns if the field is a created key.
func (field Field) IsCreatedKey() bool {
	return field.isCreatedKey
}

// IsUpdatedKey returns if the field is a updated key.
func (field Field) IsUpdatedKey() bool {
	return field.isUpdatedKey
}

// IsDeletedKey returns if the field is a deleted key.
func (field Field) IsDeletedKey() bool {
	return field.isDeletedKey
}

// Type returns the reflect's type of the field.
func (field Field) Type() reflect.Type {
	return field.rtype
}

// String returns a human readable version of current instance.
func (field Field) String() string {
	buffer := &bytes.Buffer{}
	debugField(field).write(buffer)
	return buffer.String()
}

// NewField returns full column name from model, field and tag.
func NewField(driver Driver, schema *Schema, model Model, name string, args ...ModelOpts) (*Field, error) {
	if schema == nil {
		return nil, errors.New("schema is required to generate a field instance")
	}

	opts := defaultModelOpts()
	if len(args) > 0 {
		opts = args[0]
	}

	field, ok := reflectx.GetFieldByName(model, name)
	if !ok {
		return nil, errors.Errorf("field '%s' not found in model", name)
	}

	tags := GetTags(field)

	rtype := field.Type
	if rtype.Kind() == reflect.Ptr {
		rtype = rtype.Elem()
	}

	modelName := reflectx.GetIndirectTypeName(model)
	tableName := model.TableName()

	columnName := tags.GetByKey(TagName, TagKeyColumn)
	if columnName == "" {
		columnName = snaker.CamelToSnake(name)
	}

	columnPath := fmt.Sprintf("%s.%s", tableName, columnName)

	isPrimaryKey := tags.HasKey(TagName, TagKeyPrimaryKey)
	isForeignKey := tags.HasKey(TagName, TagKeyForeignKey)
	isExcluded := tags.HasKey(TagName, TagKeyIgnored) || field.PkgPath != ""
	hasDefault := tags.HasKey(TagName, TagKeyDefault)
	hasULID := tags.GetByKey(TagName, TagKeyPrimaryKey) == TagKeyULID

	isCreatedKey := columnName == opts.CreatedKey
	isUpdatedKey := columnName == opts.UpdatedKey
	isDeletedKey := columnName == opts.DeletedKey

	hasDefault = hasDefault || isCreatedKey || isUpdatedKey

	instance := &Field{
		modelName:    modelName,
		tableName:    tableName,
		fieldName:    name,
		columnName:   columnName,
		columnPath:   columnPath,
		isPrimaryKey: isPrimaryKey,
		isForeignKey: isForeignKey,
		isExcluded:   isExcluded,
		isCreatedKey: isCreatedKey,
		isUpdatedKey: isUpdatedKey,
		isDeletedKey: isDeletedKey,
		hasDefault:   hasDefault,
		hasULID:      hasULID,
		rtype:        rtype,
	}

	// Early return if the field type is not an association.
	reference := ToModel(rtype)
	if reference == nil {
		return instance, nil
	}

	// associationType := AssociationTypeUndefined
	//
	// if field.Type.Kind() == reflect.Struct {
	// 	associationType = AssociationTypeOne
	// }
	//
	// if field.Type.Kind() == reflect.Slice {
	// 	associationType = AssociationTypeMany
	// }
	//
	// if associationType == AssociationTypeUndefined {
	// 	return field, errors.Errorf("unable to guess the association type for field %s", name)
	// }
	//
	// field.IsAssociation = true
	// field.AssociationType = associationType
	//
	// // author -> author_id
	// field.ColumnName = fmt.Sprintf("%s_%s", field.ColumnName, strings.ToLower(PrimaryKeyFieldName))
	//
	// field.ForeignKey, err = NewForeignKey(driver, field)
	// if err != nil {
	// 	return field, err
	// }
	//
	// return field, nil
	return nil, nil
}

// ----------------------------------------------------------------------------
// Foreign key field
// ----------------------------------------------------------------------------

// // ForeignKey is a foreign key
// type ForeignKey struct {
// 	Schema               *Schema
// 	Model                Model
// 	ModelName            string
// 	TableName            string
// 	FieldName            string
// 	ColumnName           string
// 	AssociationFieldName string
// 	Reference            *ForeignKey
// }

// func (fk ForeignKey) String() string {
// 	return fmt.Sprintf("{{ model:%s tb:%s field:%s col:%s assoc:%s -- reference: %s }}",
// 		fk.ModelName, fk.TableName,
// 		fk.FieldName, fk.ColumnName,
// 		fk.AssociationFieldName, fk.Reference,
// 	)
// }

// ColumnPath is the foreign key column path.
// func (fk ForeignKey) ColumnPath() string {
// 	return fmt.Sprintf("%s.%s", fk.TableName, fk.ColumnName)
// }
//
// // NewForeignKey returns a new ForeignKey instance from the given field instance.
// func NewForeignKey(driver Driver, field Field) (*ForeignKey, error) {
//
// 	fmt.Printf("::4 %+v\n", field)
//
// 	var (
// 		referenceModel     = ToModel(field.Type)
// 		referenceModelName = reflect.Indirect(reflect.ValueOf(referenceModel)).Type().Name()
// 		referenceTableName = referenceModel.TableName()
// 	)
//
// 	referenceSchema, err := GetSchema(driver, referenceModel)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	fmt.Printf("::5 %+v\n", referenceSchema)
//
// 	fieldName := field.Tags.GetByKey(StructTagName, StructTagForeignKey)
// 	if fieldName == "" {
// 		fieldName = fmt.Sprintf("%s%s", field.FieldName, PrimaryKeyFieldName)
// 	}
//
// 	fmt.Printf("::6 %+v\n", fieldName)
//
// 	// Article.Author(User)
// 	if field.AssociationType == AssociationTypeOne {
// 		fmt.Printf("::7 has_once\n")
// 		fk := &ForeignKey{
// 			Schema:               &referenceSchema,
// 			Model:                field.Model,      // Article model
// 			ModelName:            field.ModelName,  // Article
// 			TableName:            field.TableName,  // articles
// 			FieldName:            fieldName,        // AuthorID
// 			ColumnName:           field.ColumnName, // author_id
// 			AssociationFieldName: field.FieldName,  // Author
//
// 			Reference: &ForeignKey{
// 				Model:      referenceModel,                       // User model
// 				ModelName:  referenceModelName,                   // User
// 				TableName:  referenceTableName,                   // users
// 				FieldName:  PrimaryKeyFieldName,                  // ID
// 				ColumnName: strings.ToLower(PrimaryKeyFieldName), // id
// 			},
// 		}
// 		spew.Dump(fk)
// 		return fk, nil
// 	}
//
// 	// User.Avatars(Avatar) -- Avatar.UserID
// 	if field.AssociationType == AssociationTypeMany {
// 		fmt.Printf("::7 has_many\n")
// 		fk := &ForeignKey{
// 			Model:                referenceModel,                                                                                   // Avatar model
// 			ModelName:            referenceModelName,                                                                               // Avatar
// 			TableName:            referenceTableName,                                                                               // avatars
// 			FieldName:            fmt.Sprintf("%s%s", field.ModelName, PrimaryKeyFieldName),                                        // UserID
// 			ColumnName:           fmt.Sprintf("%s_%s", snaker.CamelToSnake(field.ModelName), strings.ToLower(PrimaryKeyFieldName)), // user_id
// 			AssociationFieldName: field.ModelName,                                                                                  // User
//
// 			Reference: &ForeignKey{
// 				Model:                field.Model,                          // User model
// 				ModelName:            field.ModelName,                      // User
// 				TableName:            field.TableName,                      // users
// 				FieldName:            PrimaryKeyFieldName,                  // ID
// 				ColumnName:           strings.ToLower(PrimaryKeyFieldName), // id
// 				AssociationFieldName: field.FieldName,                      // Avatars
// 			},
// 		}
// 		spew.Dump(fk)
// 		return fk, nil
// 	}
// 	fmt.Printf("::7 has_none ?\n")
// 	return nil, nil
// }
