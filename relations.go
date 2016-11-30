package sqlxx

import "reflect"

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

// makeRelation creates a new relation.
func makeRelation(model Model, meta Meta, typ RelationType) (Relation, error) {
	var err error

	relation := Relation{
		Name: meta.Name,
		Type: typ,
	}

	relation.FK, err = makeForeignKeyField(model, meta)
	if err != nil {
		return relation, err
	}

	relation.Model = getModelType(meta.Type)

	schema, err := GetSchema(relation.Model)
	if err != nil {
		return relation, err
	}

	relation.Schema = schema

	relation.Reference, err = makeReferenceField(relation.Model, "ID")
	if err != nil {
		return relation, err
	}

	return relation, nil
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
