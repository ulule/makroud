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
	if !GetIndirectValue(out).CanAddr() {
		return errors.New("model instance must be addressable (pointer required)")
	}

	schema, err := GetSchema(out)
	if err != nil {
		return err
	}

	assocs := []Field{}
	assocsOfAssocs := map[string]Field{}

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

	//
	// 1st level
	//

	err = PreloadAssociations(driver, out, assocs)
	if err != nil {
		return err
	}

	//
	// 2nd level: associations of association
	//

	if IsSlice(out) {
		for k, v := range assocsOfAssocs {
			err = preloadAssociationForSlice(driver, out, schema, k, v)
			if err != nil {
				return err
			}
		}
		return nil
	}

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

func preloadAssociationForSlice(driver Driver, out interface{}, schema Schema, fieldName string, field Field) error {
	slice := reflect.ValueOf(out).Elem()

	type rel struct {
		item  reflect.Value // ex: pointer to Article
		assoc reflect.Value // ex: pointer to Article.User
		field string        // ex: APIKey (for Article.User.APIKey)
	}

	// fk -> relation
	relations := map[int64]rel{}

	// Articles
	for i := 0; i < slice.Len(); i++ {
		value := slice.Index(i)

		if value.Kind() != reflect.Ptr && value.CanAddr() {
			value = value.Addr()
		}

		// Retrieve Article.User previously fetched
		assocValue, assocPtr, err := getFieldValues(value.Interface(), fieldName)
		if err != nil {
			return err
		}

		// Retrieve Article.User.APIKeyID (for SELECT IN)
		fkv, err := GetFieldValue(assocValue.Interface(), field.ForeignKey.FieldName)
		if err != nil {
			return err
		}

		fk, err := IntToInt64(fkv)
		if err != nil {
			return err
		}

		if fk != int64(0) {
			relations[fk] = rel{
				item:  value,
				assoc: assocPtr,
				field: field.ForeignKey.AssociationFieldName,
			}
		}
	}

	var fks []int64
	for k := range relations {
		fks = append(fks, k)
	}

	// Build a []APIKey slice
	fkAssocType := reflect.SliceOf(GetIndirectType(reflect.TypeOf(field.ForeignKey.Reference.Model)))
	fkAssocs := reflect.New(fkAssocType)
	fkAssocs.Elem().Set(reflect.MakeSlice(fkAssocType, 0, 0))

	// SELECT * FROM from api_keys WHERE id IN childrenRelationPKs
	// TODO: HERE GET PRIMAREY KEY OF REFERENCE
	err := FindByParams(driver, fkAssocs.Interface(), map[string]interface{}{"id": fks})
	if err != nil {
		return err
	}

	// Slice ptr -> value
	fkAssocs = fkAssocs.Elem()

	for i := 0; i < fkAssocs.Len(); i++ {
		assoc := fkAssocs.Index(i)

		// APIKey.ID
		pkv, err := GetFieldValue(assoc, "ID")
		if err != nil {
			return err
		}

		pk, err := IntToInt64(pkv)
		if err != nil {
			return err
		}

		if pk != int64(0) {
			relation, ok := relations[pk]
			if ok {
				err := SetFieldValue(relation.assoc.Interface(), relation.field, assoc.Interface())
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
