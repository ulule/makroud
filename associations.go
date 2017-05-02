package sqlxx

import (
	"database/sql"
	"fmt"
	"reflect"
)

// AssociationQueries are a slice of relation query ready to be ordered by level
type AssociationQueries []AssociationQuery

// AssociationQuery is a relation query
type AssociationQuery struct {
	Field    Field
	Query    string
	Args     []interface{}
	Params   map[string]interface{}
	FetchOne bool
}

// String returns struct instance string representation.
func (aq AssociationQuery) String() string {
	return aq.Query
}

// GetAssociationQueries returns relation queries ASC sorted by their level
func GetAssociationQueries(out interface{}, fields []Field) (AssociationQueries, error) {
	var (
		queries = AssociationQueries{}
		isSlice = IsSlice(out)
	)

	for _, field := range fields {
		if !field.IsAssociation {
			return nil, fmt.Errorf("field '%s' is not an association", field.Name)
		}

		if field.ForeignKey == nil {
			return nil, fmt.Errorf("no ForeignKey instance found for field %s", field.Name)
		}

		var (
			err    error
			params = map[string]interface{}{}
			pks    = []interface{}{}
		)

		if !isSlice {
			// For Category.User, should be: Category.UserID
			pks, err = GetPrimaryKeys(out, field.ForeignKey.FieldName)
			if err != nil {
				return nil, err
			}
		} else {
			value := reflect.ValueOf(out).Elem()
			for i := 0; i < value.Len(); i++ {
				values, err := GetPrimaryKeys(value.Index(i).Interface(), field.ForeignKey.FieldName)
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
			// For Category.User, should be: users.id
			params[field.ForeignKey.Reference.ColumnName] = pks
		} else {
			params[field.ForeignKey.Reference.ColumnName] = pks[0]
		}

		query, args, err := whereQuery(field.ForeignKey.Reference.Model, params, field.IsAssociationTypeOne() && !isSlice)
		if err != nil {
			return nil, err
		}

		queries = append(queries, AssociationQuery{
			Field:    field,
			Query:    query,
			Args:     args,
			Params:   params,
			FetchOne: field.IsAssociationTypeOne(),
		})
	}

	return queries, nil
}

// ----------------------------------------------------------------------------
// Preloading
// ----------------------------------------------------------------------------

// PreloadAssociations preloads relations of out from queries.
func PreloadAssociations(driver Driver, out interface{}, fields []Field) error {
	queries, err := GetAssociationQueries(out, fields)
	if err != nil {
		return err
	}

	for _, query := range queries {
		err := SetAssociation(driver, out, query)
		if err != nil {
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Getter / setter
// ----------------------------------------------------------------------------

// SetAssociation performs query and populates the given out with values.
func SetAssociation(driver Driver, out interface{}, q AssociationQuery) error {
	if !q.Field.IsAssociation {
		return fmt.Errorf("cannot set association for field: %v", q.Field)
	}

	var (
		err      error
		instance interface{}
		isSlice  = IsSlice(out)
	)

	if q.Field.IsAssociationTypeMany() || isSlice {
		instance = CloneType(q.Field.ForeignKey.Reference.Model, reflect.Slice)
	} else {
		instance = CloneType(q.Field.ForeignKey.Reference.Model)
	}

	if err = FetchAssociation(driver, instance, q); err != nil {
		return err
	}

	// user.Avatars || user.Avatar
	if !isSlice {
		return SetFieldValue(out, q.Field.ForeignKey.AssociationFieldName, reflect.ValueOf(instance).Elem().Interface())
	}

	//
	// Slice
	//

	// users.Avatar
	if !q.Field.IsAssociationTypeMany() {
		value := reflect.ValueOf(out).Elem()

		// user.Avatar
		if !isSlice {
			for i := 0; i < value.Len(); i++ {
				err := SetFieldValue(value.Index(i), q.Field.ForeignKey.AssociationFieldName, instance)
				if err != nil {
					return err
				}
			}
		} else {
			var (
				instancesMap = map[interface{}]reflect.Value{}
				items        = reflect.ValueOf(instance).Elem()
			)

			for i := 0; i < items.Len(); i++ {
				value, err := GetFieldValue(items.Index(i), q.Field.ForeignKey.Reference.FieldName)
				if err != nil {
					return err
				}
				instancesMap[value] = items.Index(i)
			}

			for i := 0; i < value.Len(); i++ {
				val, err := GetFieldValue(value.Index(i), q.Field.ForeignKey.FieldName)
				if err != nil {
					return err
				}

				switch val.(type) {
				case sql.NullInt64:
					val = int(val.(sql.NullInt64).Int64)
				}

				instance, ok := instancesMap[val]
				if ok {
					err := SetFieldValue(value.Index(i), q.Field.ForeignKey.AssociationFieldName, instance.Interface())
					if err != nil {
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

		itemPK, err := GetFieldValue(item.Interface().(Model), q.Field.ForeignKey.FieldName)
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

			relatedFK, err := GetFieldValue(relatedItemInstance, q.Field.ForeignKey.FieldName)
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

		var (
			newSlice      = MakeSlice(q.Field.Model)
			newSliceValue = reflect.ValueOf(newSlice)
			field         = item.FieldByName(q.Field.Name)
		)

		for _, related := range itemRelatedItems {
			newSliceValue = reflect.Append(newSliceValue, related)
		}

		field.Set(newSliceValue)
	}

	return nil
}

// FetchAssociation fetches the given relation.
func FetchAssociation(driver Driver, out interface{}, query AssociationQuery) error {
	if query.FetchOne && !IsSlice(out) {
		return driver.Get(out, driver.Rebind(query.Query), query.Args...)
	}

	return driver.Select(out, driver.Rebind(query.Query), query.Args...)
}
