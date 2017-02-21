package sqlxx

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/ulule/sqlxx/reflekt"
)

// Preloader is a custom preloader.
type Preloader func(d Driver) (Driver, error)

// Preload preloads related fields.
func Preload(driver Driver, out interface{}, fields ...string) error {
	var err error

	schema, err := GetSchema(out)
	if err != nil {
		return err
	}

	type ChildRelation struct {
		field         string
		relationField string
		relation      Relation
	}

	var (
		relations      []Relation
		childRelations []ChildRelation
		relationPaths  = schema.RelationPaths()
	)

	for _, field := range fields {
		relation, ok := relationPaths[field]
		if !ok {
			return fmt.Errorf("%s is not a valid relation", field)
		}

		splits := strings.Split(field, ".")
		if len(splits) == 1 {
			relations = append(relations, relation)
		}

		if len(splits) == 2 {
			childRelations = append(childRelations, ChildRelation{
				field:         splits[0],
				relationField: strings.Join(splits[1:], "."),
				relation:      relation,
			})
		}
	}

	if err = preloadRelations(driver, out, relations); err != nil {
		return err
	}

	type Relationship struct {
		item                       interface{}
		itemValue                  reflect.Value
		itemPK                     interface{}
		itemChild                  interface{}
		itemChildFieldName         string
		itemChildPK                interface{}
		itemChildPKFieldName       string
		itemChildRelationPK        interface{}
		itemChildRelationFieldName string
	}

	for _, child := range childRelations {
		// Articles
		if reflekt.IsSlice(out) {
			var (
				slice               = reflect.ValueOf(out).Elem()
				relationships       []Relationship
				childrenRelationPKs []interface{}
			)

			for i := 0; i < slice.Len(); i++ {
				var (
					// Article
					itemValue            = slice.Index(i)
					item                 = itemValue.Interface()
					itemPKFieldName      = schema.PrimaryField.Name
					itemChildFieldName   = child.field
					itemChildPKFieldName = child.relation.RelatedFKField()
				)

				// Retrieve Article.ID
				itemPK, err := reflekt.GetFieldValue(item, itemPKFieldName)
				if err != nil {
					return err
				}

				// Retrieve Article.User previously fetched
				itemChild, err := reflekt.GetFieldValue(item, itemChildFieldName)
				if err != nil {
					return err
				}

				// Retrieve Article.UserID
				itemChildPK, err := reflekt.GetFieldValue(itemChild, child.relation.ParentSchema.PrimaryField.Name)
				if err != nil {
					return err
				}

				// Retrieve Article.User.APIKeyID (for SELECT IN)
				itemChildRelationPK, err := reflekt.GetFieldValue(itemChild, child.relation.RelatedFKField())
				if err != nil {
					return err
				}

				relationships = append(relationships, Relationship{
					item:                       item,
					itemValue:                  itemValue,
					itemPK:                     itemPK,
					itemChild:                  itemChild,
					itemChildFieldName:         itemChildFieldName,
					itemChildPK:                itemChildPK,
					itemChildPKFieldName:       itemChildPKFieldName,
					itemChildRelationPK:        itemChildRelationPK,
					itemChildRelationFieldName: child.relationField,
				})

				var exists bool

				for _, v := range childrenRelationPKs {
					if v == itemChildRelationPK {
						exists = true
						break
					}
				}

				if !exists {
					childrenRelationPKs = append(childrenRelationPKs, itemChildRelationPK)
				}

			}

			// Build a []APIKey slice
			t := reflect.SliceOf(reflekt.GetIndirectType(reflect.TypeOf(child.relation.Model)))
			s := reflect.New(t)
			s.Elem().Set(reflect.MakeSlice(t, 0, 0))

			// SELECT * FROM from api_keys WHERE id IN childrenRelationPKs
			if err = FindByParams(driver, s.Interface(), map[string]interface{}{child.relation.Schema.PrimaryField.ColumnName: childrenRelationPKs}); err != nil {
				return err
			}

			instances := s.Elem()

			// Set the relations
			for _, relationship := range relationships {
				for i := 0; i < instances.Len(); i++ {
					instance := instances.Index(i)

					// APIKey.ID
					instancePK, err := reflekt.GetFieldValue(instance, child.relation.Schema.PrimaryField.Name)
					if err != nil {
						return err
					}

					if relationship.itemChildRelationPK == instancePK {
						itemChildCopy := reflekt.Copy(relationship.itemChild)
						if err = reflekt.SetFieldValue(itemChildCopy, relationship.itemChildRelationFieldName, instance.Interface()); err != nil {
							return err
						}
						if err = reflekt.SetFieldValue(relationship.itemValue, relationship.itemChildFieldName, itemChildCopy); err != nil {
							return err
						}
					}
				}
			}

			return nil
		}

		instance, err := reflekt.GetFieldValue(out, child.field)
		if err != nil {
			return err
		}

		cp := reflekt.Copy(instance)

		if err = Preload(driver, cp, child.relationField); err != nil {
			return err
		}

		if err = reflekt.SetFieldValue(out, child.field, cp); err != nil {
			return err
		}
	}

	return nil
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

	if isMany || isSlice {
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

	// user.Avatars || user.Avatar
	if !isSlice {
		return reflekt.SetFieldValue(out, rq.relation.Name, reflect.ValueOf(instance).Elem().Interface())
	}

	//
	// Slice
	//

	// users.Avatar
	if !isMany {
		value := reflect.ValueOf(out).Elem()

		// user.Avatar
		if !isSlice {
			for i := 0; i < value.Len(); i++ {
				if err := reflekt.SetFieldValue(value.Index(i), rq.relation.Name, instance); err != nil {
					return err
				}
			}
		} else {
			instancesMap := map[interface{}]reflect.Value{}

			items := reflect.ValueOf(instance).Elem()

			for i := 0; i < items.Len(); i++ {
				value, err := reflekt.GetFieldValue(items.Index(i), rq.relation.Reference.Name)

				if err != nil {
					return nil
				}

				instancesMap[value] = items.Index(i)
			}

			for i := 0; i < value.Len(); i++ {
				val, err := reflekt.GetFieldValue(value.Index(i), rq.relation.RelatedFKField())

				if err != nil {
					return nil
				}

				switch val.(type) {
				case sql.NullInt64:
					val = int(val.(sql.NullInt64).Int64)
				}

				instance, ok := instancesMap[val]

				if ok {
					if err := reflekt.SetFieldValue(value.Index(i), rq.relation.Name, instance.Interface()); err != nil {
						return err
					}
				}
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

	if rq.fetchOne && !reflekt.IsSlice(out) {
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
