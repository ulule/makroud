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
		return fmt.Sprintf("%sID", reflekt.ReflectType(r.ParentModel).Name())
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

// makeRelation creates a new relation.
func makeRelation(schema Schema, model Model, meta reflekt.FieldMeta, typ RelationType) (Relation, error) {
	var (
		err       error
		modelType = reflekt.ReflectType(model)
		refModel  = TypeToModel(meta.Type)
		refType   = reflekt.ReflectType(refModel)
	)

	refStructField, ok := refType.FieldByName("ID")
	if !ok {
		return Relation{}, fmt.Errorf("Field %s does not exist", meta.Name)
	}

	refMeta := reflekt.GetFieldMeta(refStructField, SupportedTags, TagsMapping)

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
		relation.FK, err = makeField(refModel, refMeta)
		if err != nil {
			return relation, err
		}

		// Defaults to "<model>_id"
		relation.FK.ColumnName = fmt.Sprintf("%s_%s", snaker.CamelToSnake(reflect.TypeOf(model).Name()), relation.Schema.PrimaryField.ColumnName)

		relation.Reference, err = makeField(reflect.New(modelType).Interface().(Model), refMeta)
		if err != nil {
			return relation, err
		}

	} else {
		relation.FK, err = makeField(model, meta)
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

		relation.Reference, err = makeField(reflect.New(refType).Interface().(Model), refMeta)
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
		if !reflekt.IsSlice(out) {
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
			return nil, err
		}

		if len(pks) > 1 {
			params[columnName] = pks
		} else {
			params[columnName] = pks[0]
		}

		query, args, err := whereQuery(relation.Model, params, relation.IsOne())
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

// preloadRelations preloads relations of out from queries.
func preloadRelations(driver Driver, out interface{}, relations []Relation) error {
	var err error

	queries, err := getRelationQueries(out, relations)
	if err != nil {
		return err
	}

	for _, rq := range queries {
		if err = setRelation(driver, out, rq); err != nil {
			return err
		}
	}

	return nil
}

// setRelation performs query and populates the given out with values.
func setRelation(driver Driver, out interface{}, rq RelationQuery) error {
	var (
		err      error
		instance interface{}
	)

	var (
		// Example: []User{}
		isSlice = reflekt.IsSlice(out)
		// Example: User.Avatars
		isMany = !rq.relation.IsOne()
	)

	//
	// If it's a many, let's clone the type as out for sqlx Get() / Select()
	//

	if isMany {
		instance = reflekt.CloneType(rq.relation.Model, reflect.Slice)
	} else {
		instance = reflekt.CloneType(rq.relation.Model)
	}

	//
	// Fetch all related relations
	//

	if err = fetchRelation(driver, instance, rq); err != nil {
		return err
	}

	//
	// Instance
	//

	if !isSlice {
		// User.Avatar
		if !isMany {
			return reflekt.SetFieldValue(out, rq.relation.Name, InterfaceToModel(instance))
		}

		// User.Avatars
		return reflekt.SetFieldValue(out, rq.relation.Name, reflect.ValueOf(instance).Elem().Interface())
	}

	//
	// Slice
	//

	// Users.Avatar
	if !isMany {
		value := reflect.ValueOf(out).Elem()

		for i := 0; i < value.Len(); i++ {
			item := value.Index(i)
			if !item.CanSet() {
				continue
			}

			val := reflect.Indirect(reflect.ValueOf(instance))
			if val.IsValid() {
				field := item.FieldByName(rq.relation.Name)
				field.Set(val)
			}
		}

		return nil
	}

	//
	// Users.Avatars
	//

	// Users
	items := reflect.ValueOf(out).Elem()

	// Avatars
	relatedItems := reflect.ValueOf(instance).Elem()

	// Iterate over slice items (Users)
	for i := 0; i < items.Len(); i++ {
		item := items.Index(i)

		if !item.CanSet() {
			continue
		}

		itemPK, err := reflekt.GetFieldValue(item.Interface().(Model), rq.relation.Reference.Name)
		if err != nil {
			return err
		}

		// Build the related items's item
		itemRelatedItems := []reflect.Value{}

		// Iterate over related items (Avatars)
		for ii := 0; ii < relatedItems.Len(); ii++ {
			var (
				relatedItem         = relatedItems.Index(ii)
				relatedItemInstance = relatedItem.Interface().(Model)
			)

			relatedFK, err := reflekt.GetFieldValue(relatedItemInstance, rq.relation.RelatedFKField())
			if err != nil {
				return err
			}

			// Compare User's avatar
			if itemPK == relatedFK {
				itemRelatedItems = append(itemRelatedItems, relatedItem)
			}
		}

		//
		// Build the related model instance slice and set it to related item.
		//

		newSlice := reflekt.MakeSlice(rq.relation.Model)
		newSliceValue := reflect.ValueOf(newSlice)

		for _, related := range itemRelatedItems {
			newSliceValue = reflect.Append(newSliceValue, related)
		}

		field := item.FieldByName(rq.relation.Name)
		field.Set(newSliceValue)
	}

	return nil
}

// fetchRelation fetches the given relation.
func fetchRelation(driver Driver, out interface{}, rq RelationQuery) error {
	var err error

	if rq.fetchOne {
		if err = driver.Get(out, driver.Rebind(rq.query), rq.args...); err != nil {
			return err
		}
	} else {
		if err = driver.Select(out, driver.Rebind(rq.query), rq.args...); err != nil {
			return err
		}
	}

	return nil
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
