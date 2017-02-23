package sqlxx

import (
	"fmt"
	"reflect"

	"github.com/ulule/sqlxx/reflekt"
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

// IsOne returns true if the association is an AssociationTypeOne.
func (a Association) IsOne() bool {
	return a.Type == AssociationTypeOne
}

// IsMany returns true if the association is an AssociationTypeMany.
func (a Association) IsMany() bool {
	return a.Type == AssociationTypeMany
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

// ----------------------------------------------------------------------------
// Queries
// ----------------------------------------------------------------------------

type (
	// AssociationQuery is a relation query
	AssociationQuery struct {
		Field    Field
		Query    string
		Args     []interface{}
		Params   map[string]interface{}
		FetchOne bool
	}
	// AssociationQueries are a slice of relation query ready to be ordered by level
	AssociationQueries []AssociationQuery
)

// GetAssociationQueries returns relation queries ASC sorted by their level
func GetAssociationQueries(out interface{}, fields []Field) (AssociationQueries, error) {
	var (
		queries = AssociationQueries{}
		isSlice = reflekt.IsSlice(out)
	)

	for _, field := range fields {
		var (
			err    error
			params = map[string]interface{}{}
			pks    = []interface{}{}
		)

		// Out is a slice, we must iterate over items and retrieve pk for each one.
		// Out is a struct, just retrieve pk

		if !isSlice {
			pks, err = GetPrimaryKeys(out, field.Association.FieldName)
			if err != nil {
				return nil, err
			}
		} else {
			value := reflect.ValueOf(out).Elem()
			for i := 0; i < value.Len(); i++ {
				values, err := GetPrimaryKeys(value.Index(i).Interface(), field.Name)
				if err != nil {
					return nil, err
				}
				pks = append(pks, values...)
			}
		}

		// Zero
		if len(pks) == 0 {
			continue
		}

		if len(pks) > 1 {
			params[field.ColumnName] = pks
		} else {
			params[field.ColumnName] = pks[0]
		}

		query, args, err := whereQuery(field.Model, params, field.Association.IsOne() && !isSlice)
		if err != nil {
			return nil, err
		}

		queries = append(queries, AssociationQuery{
			Field:    field,
			Query:    query,
			Args:     args,
			Params:   params,
			FetchOne: field.Association.IsOne(),
		})
	}

	return queries, nil
}
