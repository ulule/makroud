package sqlxx

import (
	"fmt"
	"strings"

	"reflect"

	"github.com/serenize/snaker"
)

// Field is a field.
type Field struct {
	// The reflect.StructField instance
	StructField reflect.StructField
	// The reflect Type of the field
	Type reflect.Type
	// The field struct tags
	Tags FieldTags

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
	// Does this field is a foreign key (the foreign key ID)?
	IsForeignKey bool

	// Does this field is an association?
	IsAssociation bool
	// The association type
	AssociationType AssociationType
	// ForeignKey contains foreign key relations information
	ForeignKey *ForeignKey
}

// String returns struct instance string representation.
func (f Field) String() string {
	return fmt.Sprintf("field(model:%s table:%s name:%s column:%s)",
		f.ModelName,
		f.TableName,
		f.Name,
		f.ColumnName)
}

// IsAssociationTypeOne returns true if the field is an AssociationTypeOne.
func (f Field) IsAssociationTypeOne() bool {
	return f.AssociationType == AssociationTypeOne
}

// IsAssociationTypeMany returns true if the field is an AssociationTypeMany.
func (f Field) IsAssociationTypeMany() bool {
	return f.AssociationType == AssociationTypeMany
}

// ColumnPath returns database full column path.
func (f Field) ColumnPath() string {
	return fmt.Sprintf("%s.%s", f.TableName, f.ColumnName)
}

// NewField returns full column name from model, field and tag.
func NewField(model Model, name string) (Field, error) {
	structField, fieldFound := GetIndirectValue(model).Type().FieldByName(name)
	if !fieldFound {
		return Field{}, fmt.Errorf("field '%s' not found in model", name)
	}

	var (
		err          error
		tags         = GetFieldTags(structField, SupportedTags, TagsMapping)
		columnName   = snaker.CamelToSnake(name)
		isExcluded   = false
		isForeignKey = false
		fieldType    = structField.Type
	)

	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	if v := tags.GetByKey(SQLXStructTagName, StructTagSQLXField); v != "" {
		columnName = v
	}

	if tag := tags.GetByKey(SQLXStructTagName, StructTagSQLXField); structField.PkgPath != "" || tag == "-" {
		isExcluded = true
	}

	if tags.HasKey(StructTagName, StructTagForeignKey) || (len(name) > len(PrimaryKeyFieldName) && strings.HasSuffix(name, PrimaryKeyFieldName)) {
		isForeignKey = true
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
		IsForeignKey:  isForeignKey,
		IsAssociation: false,
	}

	// Early return if the field type is not an association
	if modelType := GetModelFromType(field.Type); modelType == nil {
		return field, nil
	}

	associationType := AssociationTypeUndefined

	if field.Type.Kind() == reflect.Struct {
		associationType = AssociationTypeOne
	}

	if field.Type.Kind() == reflect.Slice {
		associationType = AssociationTypeMany
	}

	if associationType == AssociationTypeUndefined {
		return field, fmt.Errorf("unable to guess the association type for field %s", name)
	}

	field.IsAssociation = true
	field.AssociationType = associationType

	// author -> author_id
	field.ColumnName = fmt.Sprintf("%s_%s", field.ColumnName, strings.ToLower(PrimaryKeyFieldName))

	field.ForeignKey, err = NewForeignKey(field)
	if err != nil {
		return field, err
	}

	return field, nil
}

// ----------------------------------------------------------------------------
// Foreign key field
// ----------------------------------------------------------------------------

// ForeignKey is a foreign key
type ForeignKey struct {
	Model                Model
	ModelName            string
	TableName            string
	FieldName            string
	ColumnName           string
	AssociationFieldName string
	Reference            *ForeignKey
}

func (fk ForeignKey) String() string {
	return fmt.Sprintf("model:%s field:%s assoc:%s -- reference: %s", fk.ModelName, fk.FieldName, fk.AssociationFieldName, fk.Reference)
}

// ColumnPath is the foreign key column path.
func (fk ForeignKey) ColumnPath() string {
	return fmt.Sprintf("%s.%s", fk.TableName, fk.ColumnName)
}

// NewForeignKey returns a new ForeignKey instance from the given field instance.
func NewForeignKey(field Field) (*ForeignKey, error) {
	var (
		referenceModel     = GetModelFromType(field.Type)
		referenceModelName = GetModelName(referenceModel)
		referenceTableName = referenceModel.TableName()
	)

	// Article.Author(User)
	if field.AssociationType == AssociationTypeOne {
		return &ForeignKey{
			Model:                field.Model,                                          // Article model
			ModelName:            field.ModelName,                                      // Article
			TableName:            field.TableName,                                      // articles
			FieldName:            fmt.Sprintf("%s%s", field.Name, PrimaryKeyFieldName), // AuthorID
			ColumnName:           field.ColumnName,                                     // author_id
			AssociationFieldName: field.Name,                                           // Author

			Reference: &ForeignKey{
				Model:      referenceModel,                       // User model
				ModelName:  referenceModelName,                   // User
				TableName:  referenceTableName,                   // users
				FieldName:  PrimaryKeyFieldName,                  // ID
				ColumnName: strings.ToLower(PrimaryKeyFieldName), // id
			},
		}, nil
	}

	// User.Avatars(Avatar) -- Avatar.UserID
	if field.AssociationType == AssociationTypeMany {
		return &ForeignKey{
			Model:                referenceModel,                                                                                   // Avatar model
			ModelName:            referenceModelName,                                                                               // Avatar
			TableName:            referenceTableName,                                                                               // avatars
			FieldName:            fmt.Sprintf("%s%s", field.ModelName, PrimaryKeyFieldName),                                        // UserID
			ColumnName:           fmt.Sprintf("%s_%s", snaker.CamelToSnake(field.ModelName), strings.ToLower(PrimaryKeyFieldName)), // user_id
			AssociationFieldName: field.ModelName,                                                                                  // User

			Reference: &ForeignKey{
				Model:                field.Model,                          // User model
				ModelName:            field.ModelName,                      // User
				TableName:            field.TableName,                      // users
				FieldName:            PrimaryKeyFieldName,                  // ID
				ColumnName:           strings.ToLower(PrimaryKeyFieldName), // id
				AssociationFieldName: field.Name,                           // Avatars
			},
		}, nil
	}

	return nil, nil
}
