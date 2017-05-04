package sqlxx

import (
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
		err     error
		queries AssociationQueries
		isSlice = IsSlice(out)
	)

	for _, field := range fields {
		if !field.IsAssociation {
			return nil, fmt.Errorf("field '%s' is not an association", field.Name)
		}

		if field.ForeignKey == nil {
			return nil, fmt.Errorf("no ForeignKey instance found for field %s", field.Name)
		}

		params := map[string]interface{}{}
		var pks []interface{}

		if !isSlice {
			pks, err = GetPrimaryKeys(out, field.RelationFieldName())
			if err != nil {
				return nil, err
			}
		} else {
			slc := reflect.ValueOf(out).Elem()
			for i := 0; i < slc.Len(); i++ {
				v, err := GetPrimaryKeys(slc.Index(i).Interface(), field.RelationFieldName())
				if err != nil {
					return nil, err
				}
				pks = append(pks, v...)
			}
		}

		if len(pks) == 0 {
			continue
		}

		columnName := field.RelationColumnName()
		if len(pks) > 1 {
			params[columnName] = pks
		} else {
			params[columnName] = pks[0]
		}

		query, args, err := whereQuery(field.RelationModel(), params, field.IsAssociationTypeOne() && !isSlice)
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
		err     error
		isSlice = IsSlice(out)
		assoc   = q.Field.CreateAssociation(isSlice)
	)

	err = FetchAssociation(driver, assoc, q)
	if err != nil {
		return err
	}

	// Single instance

	if !isSlice {
		return SetFieldValue(out, q.Field.OneToAssociationFieldName(), reflect.ValueOf(assoc).Elem().Interface())
	}

	// Slice of instances

	instances := reflect.ValueOf(out).Elem()

	// OneTo

	if !q.Field.IsAssociationTypeMany() {
		assocs := reflect.ValueOf(assoc).Elem()

		for i := 0; i < instances.Len(); i++ {
			instance := instances.Index(i).Addr()

			fkv, err := GetFieldValue(instance.Interface(), q.Field.ForeignKey.FieldName)
			if err != nil {
				return err
			}

			fk, err := IntToInt64(fkv)
			if err != nil {
				return err
			}

			for ii := 0; ii < assocs.Len(); ii++ {
				pkv, err := GetFieldValue(assocs.Index(ii).Interface(), "ID")
				if err != nil {
					return err
				}

				pk, err := IntToInt64(pkv)
				if err != nil {
					return err
				}

				if fk == pk {
					err = SetFieldValue(instance.Interface(), q.Field.OneToAssociationFieldName(), assocs.Index(ii).Interface())
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	}

	// ManyTo

	assocs := reflect.ValueOf(assoc).Elem()

	for i := 0; i < instances.Len(); i++ {
		instance := instances.Index(i).Addr()

		pkv, err := GetFieldValue(instance.Interface(), q.Field.RelationFieldName())
		if err != nil {
			return err
		}

		pk, err := IntToInt64(pkv)
		if err != nil {
			return err
		}

		slc := reflect.ValueOf(MakeSlice(q.Field.ParentModel()))

		for ii := 0; ii < assocs.Len(); ii++ {
			assocv := assocs.Index(ii).Addr()

			fkv, err := GetFieldValue(assocv.Interface(), q.Field.ForeignKey.FieldName)
			if err != nil {
				return err
			}

			fk, err := IntToInt64(fkv)
			if err != nil {
				return err
			}

			if pk == fk {
				slc = reflect.Append(slc, assocv.Elem())
			}
		}

		err = SetFieldValue(instance.Interface(), q.Field.ManyToAssociationFieldName(), slc.Interface())
		if err != nil {
			return err
		}
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
