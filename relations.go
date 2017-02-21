package sqlxx

import (
	"fmt"
	"reflect"

	"github.com/serenize/snaker"
	"github.com/ulule/sqlxx/reflekt"
)

// ----------------------------------------------------------------------------
// Relation
// ----------------------------------------------------------------------------

// Relation represents an related field between two models.
type Relation struct {
	// The relation field name (if field is Author and model is User, field name is Author)
	Name string
	// The relation type
	Type RelationType
	// The related model
	Model Model
	// The related schema
	Schema Schema
	// The parent model
	ParentModel Model
	// The parent schema
	ParentSchema Schema
	// The foreign key field
	FK Field
	// The foreign key reference field
	Reference Field
}

// RelatedFKField returns related FK field
func (r Relation) RelatedFKField() string {
	if !r.IsOne() {
		return fmt.Sprintf("%sID", reflekt.GetIndirectType(reflect.TypeOf(r.ParentModel)).Name())
	}
	return fmt.Sprintf("%sID", r.Name)
}

// IsOne returns true if the relation is a "one" relation.
// Used to handle SELECT IN at preloading.
func (r Relation) IsOne() bool {
	_, ok := RelationsOne[r.Type]
	return ok
}

func (r Relation) String() string {
	return fmt.Sprintf("model:%s parent:%s field:%s fk:%s ref:%s",
		reflect.TypeOf(r.Model).Name(),
		reflect.TypeOf(r.ParentModel).Name(),
		r.Name,
		r.FK.ColumnPath(),
		r.Reference.ColumnPath())
}

// NewRelation creates a new relation.
func NewRelation(schema Schema, model Model, meta FieldMeta, typ RelationType) (Relation, error) {
	var (
		err       error
		modelType = reflekt.GetIndirectType(model)
		refModel  = GetModelFromType(meta.Type)
		refType   = reflekt.GetIndirectType(refModel)
	)

	refStructField, ok := refType.FieldByName("ID")
	if !ok {
		return Relation{}, fmt.Errorf("Field %s does not exist", meta.Name)
	}

	refMeta := GetFieldMeta(refStructField, SupportedTags, TagsMapping)

	refSchema, err := GetSchema(refModel)
	if err != nil {
		return Relation{}, err
	}

	relation := Relation{
		Name:         meta.Name,
		Type:         typ,
		ParentModel:  model,
		ParentSchema: schema,
		Model:        refModel,
		Schema:       refSchema,
	}

	reversed := !relation.IsOne()

	if reversed {
		relation.FK, err = NewField(refModel, refMeta)
		if err != nil {
			return relation, err
		}

		// Defaults to "<model>_id"
		relation.FK.ColumnName = fmt.Sprintf("%s_%s", snaker.CamelToSnake(GetModelName(model)), relation.Schema.PrimaryField.ColumnName)

		relation.Reference, err = NewField(reflect.New(modelType).Interface().(Model), refMeta)
		if err != nil {
			return relation, err
		}

	} else {
		relation.FK, err = NewField(model, meta)
		if err != nil {
			return relation, err
		}

		// Defaults to "fieldname_id"
		relation.FK.ColumnName = fmt.Sprintf("%s_id", relation.FK.ColumnName)
		if reversed {
			relation.FK.ColumnName = relation.Schema.PrimaryField.ColumnName
		}

		// Get the SQLX one if any.
		if customName := relation.FK.Tags.GetByKey(SQLXStructTagName, "field"); len(customName) != 0 {
			relation.FK.ColumnName = customName
		}

		relation.Reference, err = NewField(reflect.New(refType).Interface().(Model), refMeta)
		if err != nil {
			return relation, err
		}
	}

	return relation, nil
}

// ----------------------------------------------------------------------------
// Relation queries
// ----------------------------------------------------------------------------

// RelationQuery is a relation query
type RelationQuery struct {
	relation Relation
	query    string
	args     []interface{}
	params   map[string]interface{}
	fetchOne bool
}

// RelationQueries are a slice of relation query ready to be ordered by level
type RelationQueries []RelationQuery

// getRelationQueries returns relation queries ASC sorted by their level
func getRelationQueries(out interface{}, relations []Relation) (RelationQueries, error) {
	queries := RelationQueries{}

	for _, relation := range relations {
		var (
			err         error
			params      = map[string]interface{}{}
			fkFieldName = relation.RelatedFKField()
			columnName  = relation.Reference.ColumnName
			pks         = []interface{}{}
		)

		// If we have a many relation, let's reverse
		if !relation.IsOne() {
			columnName = relation.FK.ColumnName
			fkFieldName = relation.Schema.PrimaryField.Name
		}

		// Out is a slice, we must iterate over items and retrieve pk for each one.
		// Out is a struct, just retrieve pk

		isSlice := reflekt.IsSlice(out)

		if !isSlice {
			pks, err = GetPrimaryKeys(out, fkFieldName)
			if err != nil {
				return nil, err
			}
		} else {
			value := reflect.ValueOf(out).Elem()

			for i := 0; i < value.Len(); i++ {
				values, err := GetPrimaryKeys(value.Index(i).Interface(), fkFieldName)
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
			params[columnName] = pks
		} else {
			params[columnName] = pks[0]
		}

		query, args, err := whereQuery(relation.Model, params, relation.IsOne() && !isSlice)
		if err != nil {
			return nil, err
		}

		queries = append(queries, RelationQuery{
			relation: relation,
			query:    query,
			args:     args,
			params:   params,
			fetchOne: relation.IsOne(),
		})
	}

	return queries, nil
}

// getRelationType returns RelationType for the given reflect.Type.
func getRelationType(typ reflect.Type) RelationType {
	if typ.Kind() == reflect.Slice {
		typ = typ.Elem()

		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}

		if _, isModel := reflect.New(typ).Interface().(Model); isModel {
			return RelationTypeManyToOne
		}

		return RelationTypeUnknown
	}

	if _, isModel := reflect.New(typ).Interface().(Model); isModel {
		return RelationTypeOneToMany
	}

	return RelationTypeUnknown
}
