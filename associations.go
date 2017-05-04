package sqlxx

import (
	"fmt"
	"reflect"
)

// Queries is a list of Query instances
type Queries []Query

// Query is a relation query
type Query struct {
	Field    Field
	Query    string
	Args     []interface{}
	Params   map[string]interface{}
	FetchOne bool
}

// String returns struct instance string representation.
func (q Query) String() string {
	return q.Query
}

// getAssociationPrimaryKeys returns primary keys for a given association.
func getAssociationPrimaryKeys(instance interface{}, field Field) ([]int64, error) {
	var (
		err    error
		values []interface{}
		pks    []int64
	)

	if !IsSlice(instance) {
		values, err = GetPrimaryKeys(instance, field.RelationFieldName())
		if err != nil {
			return nil, err
		}
	} else {
		slc := reflect.ValueOf(instance).Elem()

		for i := 0; i < slc.Len(); i++ {
			v, err := GetPrimaryKeys(slc.Index(i).Interface(), field.RelationFieldName())
			if err != nil {
				return nil, err
			}
			values = append(values, v...)
		}
	}

	for _, value := range values {
		pk, err := IntToInt64(value)
		if err != nil {
			return nil, err
		}

		if pk != int64(0) {
			pks = append(pks, pk)
		}
	}

	return pks, nil
}

// GetAssociationQueries returns relation queries ASC sorted by their level
func GetAssociationQueries(out interface{}, fields []Field) (Queries, error) {
	var (
		err     error
		queries Queries
	)

	for _, field := range fields {
		err = checkAssociation(field)
		if err != nil {
			return nil, err
		}

		pks, err := getAssociationPrimaryKeys(out, field)
		if err != nil {
			return nil, err
		}

		if len(pks) == 0 {
			continue
		}

		params := map[string]interface{}{}
		columnName := field.RelationColumnName()

		if len(pks) > 1 {
			params[columnName] = pks
		} else {
			params[columnName] = pks[0]
		}

		query, args, err := whereQuery(field.RelationModel(), params, field.IsAssociationTypeOne() && !IsSlice(out))
		if err != nil {
			return nil, err
		}

		queries = append(queries, Query{
			Field:    field,
			Query:    query,
			Args:     args,
			Params:   params,
			FetchOne: field.IsAssociationTypeOne(),
		})
	}

	return queries, nil
}

// preloadAssociations preloads relations of out from queries.
func preloadAssociations(driver Driver, out interface{}, fields []Field) error {
	queries, err := GetAssociationQueries(out, fields)
	if err != nil {
		return err
	}

	for _, query := range queries {
		err := setAssociation(driver, out, query)
		if err != nil {
			return err
		}
	}

	return nil
}

// setAssociation performs query and populates the given out with values.
func setAssociation(driver Driver, out interface{}, q Query) error {
	err := checkAssociation(q.Field)
	if err != nil {
		return err
	}

	isSlice := IsSlice(out)
	assoc := q.Field.CreateAssociation(isSlice)

	err = fetchAssociation(driver, assoc, q)
	if err != nil {
		return err
	}

	if !isSlice {
		return SetFieldValue(out, q.Field.OneToAssociationFieldName(), reflect.ValueOf(assoc).Elem().Interface())
	}

	instances := reflect.ValueOf(out).Elem()

	if !q.Field.IsAssociationTypeMany() {
		assocs := reflect.ValueOf(assoc).Elem()

		for i := 0; i < instances.Len(); i++ {
			instance := instances.Index(i).Addr()

			fk, err := GetInt64PrimaryKey(instance.Interface(), q.Field.RelationFieldName())
			if err != nil {
				return err
			}

			for ii := 0; ii < assocs.Len(); ii++ {
				pk, err := GetInt64PrimaryKey(assocs.Index(ii).Interface(), q.Field.RelationPrimaryKeyFieldName())
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

	assocs := reflect.ValueOf(assoc).Elem()

	for i := 0; i < instances.Len(); i++ {
		instance := instances.Index(i).Addr()

		pk, err := GetInt64PrimaryKey(instance.Interface(), q.Field.RelationFieldName())
		if err != nil {
			return err
		}

		slc := reflect.ValueOf(MakeSlice(q.Field.ParentModel()))

		for ii := 0; ii < assocs.Len(); ii++ {
			assocv := assocs.Index(ii).Addr()

			fk, err := GetInt64PrimaryKey(assocv.Interface(), q.Field.ForeignKey.FieldName)
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

// fetchAssociation fetches the given relation.
func fetchAssociation(driver Driver, out interface{}, query Query) error {
	if query.FetchOne && !IsSlice(out) {
		return driver.Get(out, driver.Rebind(query.Query), query.Args...)
	}

	return driver.Select(out, driver.Rebind(query.Query), query.Args...)
}

func checkAssociation(field Field) error {
	if !field.IsAssociation {
		return fmt.Errorf("field '%s' is not an association", field.Name)
	}

	if field.ForeignKey == nil {
		return fmt.Errorf("no ForeignKey instance found for field %s", field.Name)
	}

	return nil
}
