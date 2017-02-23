package sqlxx

import (
	"fmt"
	"reflect"
)

// AssociationType is an association type.
type AssociationType uint8

// Association types
const (
	AssociationTypeUndefined = AssociationType(iota)
	AssociationTypeOne
	AssociationTypeMany
	AssociationTypeManyToMany
)

// Association is a field association.
type Association struct {
	Type            AssociationType
	Model           Model
	ModelName       string
	TableName       string
	PrimaryKeyField Field
	FieldName       string
	ColumnName      string
}

// ColumnPath returns database full column path.
func (a Association) ColumnPath() string {
	return fmt.Sprintf("%s.%s", a.TableName, a.ColumnName)
}

// NewAssociation returns a new Association instance for the given struct field.
// And a boolean either the given field is an association or not.
func NewAssociation(f reflect.StructField) (*Association, bool, error) {
	var (
		t               = f.Type
		associationType = AssociationTypeUndefined
	)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct || t.Kind() == reflect.Slice {
		associationType = AssociationTypeOne

		if t.Kind() == reflect.Slice {
			associationType = AssociationTypeMany
			t = t.Elem()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
		}

		if _, ok := reflect.New(t).Interface().(Model); !ok {
			associationType = AssociationTypeUndefined
		}
	}

	if associationType == AssociationTypeUndefined {
		return nil, false, nil
	}

	var (
		model     = GetModelFromType(t)
		modelName = GetModelName(model)
	)

	schema, err := GetSchema(model)
	if err != nil {
		return nil, true, err
	}

	return &Association{
		Type:            associationType,
		Model:           model,
		ModelName:       modelName,
		TableName:       model.TableName(),
		PrimaryKeyField: schema.PrimaryKeyField,
		FieldName:       schema.PrimaryKeyField.Name,
		ColumnName:      schema.PrimaryKeyField.ColumnName,
	}, true, nil
}
