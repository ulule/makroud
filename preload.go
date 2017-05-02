package sqlxx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Preloader is a custom preloader.
type Preloader func(d Driver) (Driver, error)

// Preload preloads related fields.
func Preload(driver Driver, out interface{}, paths ...string) error {
	var (
		err error
		// isSlice           = IsSlice(out)
		assocs         []Field
		assocsOfAssocs = map[string]Field{}
	)

	if !GetIndirectValue(out).CanAddr() {
		return errors.New("model instance must be addressable (pointer required)")
	}

	schema, err := GetSchema(out)
	if err != nil {
		return err
	}

	for _, path := range paths {
		assoc, ok := schema.Associations[path]
		if !ok {
			return fmt.Errorf("%s is not a valid association", path)
		}

		splits := strings.Split(path, ".")

		if len(splits) == 1 {
			assocs = append(assocs, assoc)
		}

		if len(splits) == 2 {
			assocsOfAssocs[splits[0]] = assoc
		}
	}

	err = PreloadAssociations(driver, out, assocs)
	if err != nil {
		return err
	}

	// if isSlice {
	// 	for _, child := range childAssociations {
	// 		spew.Dump(child)

	// 		// if err = preloadSlice(driver, out, schema, child); err != nil {
	// 		// 	return err
	// 		// }
	// 	}
	// 	return nil
	// }

	for k, v := range assocsOfAssocs {
		// At this step, value can be either a value or a pointer.
		value, err := GetFieldValue(reflect.ValueOf(out), k)
		if err != nil {
			return err
		}

		// We must reflect to check the addressability.
		reflected := reflect.ValueOf(value)
		isValue := false

		// Value is a value? Create a pointer to.
		if !reflected.CanAddr() {
			value = Copy(value)
			isValue = true
		}

		err = Preload(driver, value, v.Name)
		if err != nil {
			return err
		}

		// Relation was initially a value.
		if isValue {
			value = reflect.Indirect(reflect.ValueOf(value)).Interface()
		}

		err = SetFieldValue(out, k, value)
		if err != nil {
			return err
		}
	}

	return nil
}

// func preloadSlice(driver Driver, out interface{}, schema Schema, field Field) error {
// 	var (
// 		err                 error
// 		slice               = reflect.ValueOf(out).Elem()
// 		childrenRelationPKs []interface{}
// 	)

// 	for i := 0; i < slice.Len(); i++ {
// 		var (
// 			// Article
// 			itemValue            = slice.Index(i)
// 			item                 = itemValue.Interface()
// 			itemPKFieldName      = schema.PrimaryField.Name
// 			itemChildFieldName   = child.field
// 			itemChildPKFieldName = child.relation.RelatedFKField()
// 		)

// 		// Retrieve Article.ID
// 		itemPK, err := GetFieldValue(item, itemPKFieldName)
// 		if err != nil {
// 			return err
// 		}

// 		// Retrieve Article.User previously fetched
// 		itemChild, err := GetFieldValue(item, itemChildFieldName)
// 		if err != nil {
// 			return err
// 		}

// 		// Retrieve Article.UserID
// 		itemChildPK, err := GetFieldValue(itemChild, child.relation.ParentSchema.PrimaryField.Name)
// 		if err != nil {
// 			return err
// 		}

// 		// Retrieve Article.User.APIKeyID (for SELECT IN)
// 		itemChildRelationPK, err := GetFieldValue(itemChild, child.relation.RelatedFKField())
// 		if err != nil {
// 			return err
// 		}

// 		relationships = append(relationships, Relationship{
// 			item:                       item,
// 			itemValue:                  itemValue,
// 			itemPK:                     itemPK,
// 			itemChild:                  itemChild,
// 			itemChildFieldName:         itemChildFieldName,
// 			itemChildPK:                itemChildPK,
// 			itemChildPKFieldName:       itemChildPKFieldName,
// 			itemChildRelationPK:        itemChildRelationPK,
// 			itemChildRelationFieldName: child.relationField,
// 		})

// 		var exists bool
// 		for _, v := range childrenRelationPKs {
// 			if v == itemChildRelationPK {
// 				exists = true
// 				break
// 			}
// 		}

// 		if !exists {
// 			childrenRelationPKs = append(childrenRelationPKs, itemChildRelationPK)
// 		}

// 	}

// 	// Build a []APIKey slice
// 	t := reflect.SliceOf(GetIndirectType(reflect.TypeOf(child.relation.Model)))
// 	s := reflect.New(t)
// 	s.Elem().Set(reflect.MakeSlice(t, 0, 0))

// 	// SELECT * FROM from api_keys WHERE id IN childrenRelationPKs
// 	if err = FindByParams(driver, s.Interface(), map[string]interface{}{child.relation.Schema.PrimaryField.ColumnName: childrenRelationPKs}); err != nil {
// 		return err
// 	}

// 	instances := s.Elem()

// 	// Set the relations
// 	for _, relationship := range relationships {
// 		for i := 0; i < instances.Len(); i++ {
// 			instance := instances.Index(i)

// 			// APIKey.ID
// 			instancePK, err := GetFieldValue(instance, child.relation.Schema.PrimaryField.Name)
// 			if err != nil {
// 				return err
// 			}

// 			if relationship.itemChildRelationPK == instancePK {
// 				itemChildCopy := Copy(relationship.itemChild)
// 				if err = SetFieldValue(itemChildCopy, relationship.itemChildRelationFieldName, instance.Interface()); err != nil {
// 					return err
// 				}
// 				if err = SetFieldValue(relationship.itemValue, relationship.itemChildFieldName, itemChildCopy); err != nil {
// 					return err
// 				}
// 			}
// 		}
// 	}

// 	return nil
// }
