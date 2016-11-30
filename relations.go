package sqlxx

import (
	"fmt"
	"reflect"
)

// ----------------------------------------------------------------------------
// Relation
// ----------------------------------------------------------------------------

// Relation represents an related field between two models.
type Relation struct {
	// The relation field name (if field is Author and model is User, field name is Author)
	Name string
	// The related model
	Model Model
	// The related schema
	Schema Schema
	// The relation type
	Type RelationType
	// The foreign key field
	FK Field
	// The foreign key reference field
	Reference Field
}

// IsOne returns true if the relation is a "one" relation.
// Used to handle SELECT IN at preloading.
func (r Relation) IsOne() bool {
	_, ok := RelationsOne[r.Type]
	return ok
}

func (r Relation) String() string {
	return fmt.Sprintf("name:%s fk:%s ref:%s", r.Name, r.FK.ColumnPath(), r.Reference.ColumnPath())
}

// makeRelation creates a new relation.
func makeRelation(model Model, meta Meta, typ RelationType) (Relation, error) {
	var err error

	relation := Relation{
		Name:  meta.Name,
		Type:  typ,
		Model: getModelType(meta.Type),
	}

	relation.Schema, err = GetSchema(relation.Model)
	if err != nil {
		return relation, err
	}

	reversed := !relation.IsOne()

	relation.FK, err = makeForeignKeyField(model, meta, relation.Schema, reversed)
	if err != nil {
		return relation, err
	}

	relation.Reference, err = makeReferenceField(relation.Model, "ID")
	if err != nil {
		return relation, err
	}

	return relation, nil
}

// ----------------------------------------------------------------------------
// Relation queries
// ----------------------------------------------------------------------------

// RelationQuery is a relation query
type RelationQuery struct {
	query string
	args  []interface{}
}

// GetRelationQueries returns conditions for the given relations.
func GetRelationQueries(schema Schema, primaryKeys []interface{}, fields ...string) ([]RelationQuery, error) {
	var (
		pkCount = len(primaryKeys)
		paths   = schema.RelationPaths()
	)

	queries := []RelationQuery{}

	for _, field := range fields {
		relation, ok := paths[field]
		if !ok {
			return nil, fmt.Errorf("%s is not a valid relation", field)
		}

		var (
			params     = map[string]interface{}{}
			columnName = relation.Reference.ColumnName
		)

		// If we have a many relation, let's reverse
		if !relation.IsOne() {
			columnName = relation.Reference.ColumnName
		}

		if pkCount == 1 {
			params[columnName] = primaryKeys[0]
		} else {
			params[columnName] = primaryKeys
		}

		query, args, err := whereQuery(relation.Model, params, relation.IsOne())
		if err != nil {
			return nil, err
		}

		queries = append(queries, RelationQuery{query: query, args: args})
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
