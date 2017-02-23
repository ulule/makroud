package sqlxx

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/ulule/sqlxx/reflekt"
)

// Association is a field association.
type Association struct {
	Type AssociationType

	// Model
	Schema          Schema
	Model           Model
	ModelName       string
	TableName       string
	PrimaryKeyField Field

	// Field
	FieldName  string
	ColumnName string

	// Foreign key
	FKFieldName  string
	FKColumnName string
}

// String representation
func (a Association) String() string {
	return fmt.Sprintf("model:%s table:%s field:%s columun:%s fk:%s fk_column:%s",
		a.ModelName, a.TableName, a.FieldName, a.ColumnName, a.FKFieldName, a.FKColumnName)
}

// ColumnPath returns database full column path.
func (a Association) ColumnPath() string {
	return fmt.Sprintf("%s.%s", a.TableName, a.ColumnName)
}

// IsOne returns true if the association is an AssociationTypeOne.
func (a Association) IsOne() bool {
	return a.Type == AssociationTypeOne
}

// IsMany returns true if the association is an AssociationTypeMany.
func (a Association) IsMany() bool {
	return a.Type == AssociationTypeMany
}

// ----------------------------------------------------------------------------
// SQL queries
// ----------------------------------------------------------------------------

type (
	// AssociationQuery is a relation query
	AssociationQuery struct {
		Field    Field
		Query    string
		Args     []interface{}
		Params   map[string]interface{}
		FetchOne bool
	}
	// AssociationQueries are a slice of relation query ready to be ordered by level
	AssociationQueries []AssociationQuery
)

func (aq AssociationQuery) String() string {
	return aq.Query
}

// GetAssociationQueries returns relation queries ASC sorted by their level
func GetAssociationQueries(out interface{}, fields []Field) (AssociationQueries, error) {
	var (
		queries = AssociationQueries{}
		isSlice = reflekt.IsSlice(out)
	)

	for _, field := range fields {
		var (
			err        error
			params     = map[string]interface{}{}
			pks        = []interface{}{}
			fieldName  = field.Association.FKFieldName
			columnName = field.Association.ColumnName
			model      = field.Association.Model
			isOne      = field.Association.IsOne()
		)

		if !isSlice {
			pks, err = GetPrimaryKeys(out, fieldName)
			if err != nil {
				return nil, err
			}
		} else {
			value := reflect.ValueOf(out).Elem()

			for i := 0; i < value.Len(); i++ {
				values, err := GetPrimaryKeys(value.Index(i).Interface(), fieldName)
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

		query, args, err := whereQuery(model, params, isOne && !isSlice)
		if err != nil {
			return nil, err
		}

		queries = append(queries, AssociationQuery{
			Field:    field,
			Query:    query,
			Args:     args,
			Params:   params,
			FetchOne: isOne,
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
		if err := SetAssociation(driver, out, query); err != nil {
			return err
		}
	}

	return nil
}

// ----------------------------------------------------------------------------
// Getter / setter
// ----------------------------------------------------------------------------

// SetAssociation performs query and populates the given out with values.
func SetAssociation(driver Driver, out interface{}, query AssociationQuery) error {
	if !query.Field.IsAssociation {
		return fmt.Errorf("cannot set association for field: %v", query.Field)
	}

	var (
		err       error
		instance  interface{}
		isSlice   = reflekt.IsSlice(out)
		assoc     = query.Field.Association
		destField = query.Field.Name
		isMany    = query.Field.Association.IsMany()
	)

	if assoc.IsMany() || isSlice {
		instance = reflekt.CloneType(assoc.Model, reflect.Slice)
	} else {
		instance = reflekt.CloneType(assoc.Model)
	}

	if err = FetchAssociation(driver, instance, query); err != nil {
		return err
	}

	// user.Avatars || user.Avatar
	if !isSlice {
		return reflekt.SetFieldValue(out, destField, reflect.ValueOf(instance).Elem().Interface())
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
				if err := reflekt.SetFieldValue(value.Index(i), query.Field.Name, instance); err != nil {
					return err
				}
			}
		} else {
			var (
				instancesMap = map[interface{}]reflect.Value{}
				items        = reflect.ValueOf(instance).Elem()
			)

			for i := 0; i < items.Len(); i++ {
				value, err := reflekt.GetFieldValue(items.Index(i), query.Field.Association.FieldName)
				if err != nil {
					return nil
				}
				instancesMap[value] = items.Index(i)
			}

			for i := 0; i < value.Len(); i++ {
				val, err := reflekt.GetFieldValue(value.Index(i), query.Field.Association.FKFieldName)
				if err != nil {
					return nil
				}

				switch val.(type) {
				case sql.NullInt64:
					val = int(val.(sql.NullInt64).Int64)
				}

				instance, ok := instancesMap[val]
				if ok {
					if err := reflekt.SetFieldValue(value.Index(i), query.Field.Name, instance.Interface()); err != nil {
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

		itemPK, err := reflekt.GetFieldValue(item.Interface().(Model), query.Field.Association.FieldName)
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

			relatedFK, err := reflekt.GetFieldValue(relatedItemInstance, query.Field.Association.FieldName)
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
			newSlice      = reflekt.MakeSlice(query.Field.Model)
			newSliceValue = reflect.ValueOf(newSlice)
			field         = item.FieldByName(query.Field.Name)
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
	if query.FetchOne && !reflekt.IsSlice(out) {
		if err := driver.Get(out, driver.Rebind(query.Query), query.Args...); err != nil {
			return err
		}
		return nil
	}

	return driver.Select(out, driver.Rebind(query.Query), query.Args...)
}
