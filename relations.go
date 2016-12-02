package sqlxx

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/oleiade/reflections"
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
	// The related model
	Model Model
	// The parent model
	ParentModel Model
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
	return fmt.Sprintf("field:%s fk:%s ref:%s", r.Name, r.FK.ColumnPath(), r.Reference.ColumnPath())
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
		Name:        meta.Name,
		Type:        typ,
		Model:       refModel,
		ParentModel: model,
		Schema:      refSchema,
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
	path     string
	query    string
	args     []interface{}
	params   map[string]interface{}
	fetchOne bool
	level    int
}

// RelationQueries are a slice of relation query ready to be ordered by level
type RelationQueries []RelationQuery

// Sort interface
func (rq RelationQueries) Len() int           { return len(rq) }
func (rq RelationQueries) Less(i, j int) bool { return rq[i].level < rq[j].level }
func (rq RelationQueries) Swap(i, j int)      { rq[i], rq[j] = rq[j], rq[i] }

// getRelationQueries returns relation queries ASC sorted by their level
func getRelationQueries(schema Schema, primaryKeys []interface{}, fields ...string) (RelationQueries, error) {
	var (
		pkCount = len(primaryKeys)
		paths   = schema.RelationPaths()
	)

	queries := RelationQueries{}

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
			columnName = relation.FK.ColumnName
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

		queries = append(queries, RelationQuery{
			relation: relation,
			path:     field,
			query:    query,
			args:     args,
			params:   params,
			fetchOne: relation.IsOne(),
			level:    len(strings.Split(field, ".")),
		})
	}

	sort.Sort(queries)

	return queries, nil
}

// preloadRelations preloads relations of out from queries.
func preloadRelations(driver Driver, out interface{}, queries RelationQueries) error {
	var (
		err          error
		currentLevel = 1
	)

	for _, rq := range queries {
		if rq.level == currentLevel {
			if err = setRelation(driver, out, rq); err != nil {
				return err
			}
		} else {
			// Here get the parent model
			// Build the struct / slice
			// Performs setRelation on.
			// Set created struct / slice to out
		}
		currentLevel = rq.level
	}

	return nil
}

// setRelation performs query and populates the given out with values.
func setRelation(driver Driver, out interface{}, rq RelationQuery) error {
	var (
		err      error
		instance interface{}
	)

	isMany := !rq.relation.IsOne()

	if isMany {
		instance = reflekt.CloneType(rq.relation.Model, reflect.Slice)
	} else {
		instance = reflekt.CloneType(rq.relation.Model)
	}

	// Populate instance with data
	if err = fetchRelation(driver, instance, rq); err != nil {
		return err
	}

	if isMany {
		if err = reflections.SetField(out, rq.relation.Name, reflect.ValueOf(instance).Elem().Interface()); err != nil {
			return err
		}
	} else {
		if err = reflections.SetField(out, rq.relation.Name, InterfaceToModel(instance)); err != nil {
			return err
		}
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
