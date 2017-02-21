package sqlxx

import (
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
