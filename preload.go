package sqlxx

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// Preload preloads related fields.
func Preload(driver Driver, out interface{}, paths ...string) error {
	_, err := preload(driver, out, paths...)
	return err
}

// PreloadWithQueries preloads related fields and returns performed queries.
func PreloadWithQueries(driver Driver, out interface{}, paths ...string) (Queries, error) {
	return preload(driver, out, paths...)
}

// Preload preloads related fields.
func preload(driver Driver, out interface{}, paths ...string) (Queries, error) {
	if !reflect.Indirect(reflect.ValueOf(out)).CanAddr() {
		return nil, errors.New("model instance must be addressable (pointer required)")
	}

	schema, err := GetSchema(out)
	if err != nil {
		return nil, err
	}

	var (
		queries Queries
		isSlice = IsSlice(out)
		mapping = map[int][]Field{}
	)

	for _, path := range paths {
		field, ok := schema.Associations[path]
		if !ok {
			return nil, fmt.Errorf("%s is not a valid association", path)
		}

		splits := strings.Split(path, ".")
		level := len(splits)

		_, ok = mapping[level]
		if !ok {
			mapping[level] = []Field{}
		}

		field.DestinationField = splits[0]

		mapping[level] = append(mapping[level], field)
	}

	var levels []int
	for level := range mapping {
		levels = append(levels, level)
	}
	sort.Ints(levels)

	for _, level := range levels {
		var q Queries

		if !isSlice {
			q, err = preloadSingle(driver, out, level, mapping[level])
		} else {
			q, err = preloadSlice(driver, out, level, mapping[level])
		}

		queries = append(queries, q...)

		if err != nil {
			return queries, err
		}
	}

	return queries, nil
}

// ----------------------------------------------------------------------------
// Single instance preload
// ----------------------------------------------------------------------------

func preloadSingle(driver Driver, out interface{}, level int, fields []Field) (Queries, error) {
	var queries Queries

	for _, field := range fields {
		if level > 1 {
			relation, err := GetFieldValue(out, field.DestinationField)
			if err != nil {
				return queries, err
			}

			var (
				relationOut = Copy(relation)
				isSlice     = IsSlice(relation)
			)

			if field.IsAssociationTypeOne() {
				if isSlice {
					q, err := preloadSliceOne(driver, relationOut, field)
					queries = append(queries, q...)
					if err != nil {
						return queries, err
					}
				} else {
					q, err := preloadSingleOne(driver, relationOut, field)
					queries = append(queries, q...)
					if err != nil {
						return queries, err
					}
				}
			}

			err = SetFieldValue(out, field.DestinationField, relationOut)
			if err != nil {
				return queries, err
			}
		} else {
			if field.IsAssociationTypeOne() {
				q, err := preloadSingleOne(driver, out, field)
				queries = append(queries, q...)
				if err != nil {
					return queries, err
				}
			} else {
				q, err := preloadSingleMany(driver, out, field)
				queries = append(queries, q...)
				if err != nil {
					return queries, err
				}
			}
		}
	}

	return queries, nil
}

func preloadSingleOne(driver Driver, out interface{}, field Field) (Queries, error) {
	var queries Queries

	if !field.IsValidAssociation() {
		return nil, fmt.Errorf("field %s is not a valid association", field.Name)
	}

	fk, err := GetFieldValueInt64(out, field.ForeignKey.FieldName)
	if err != nil {
		return nil, err
	}

	if fk == int64(0) {
		return nil, err
	}

	params := map[string]interface{}{field.ForeignKey.Reference.ColumnName: fk}

	query, args, err := whereQuery(field.ForeignKey.Reference.Model, params, field.IsAssociationTypeOne())
	if err != nil {
		return nil, err
	}

	q := Query{
		Field:    field,
		Query:    query,
		Args:     args,
		Params:   params,
		FetchOne: field.IsAssociationTypeOne(),
	}

	queries = append(queries, q)

	relation := Copy(field.ForeignKey.Reference.Model)

	err = driver.Get(relation, driver.Rebind(q.Query), q.Args...)
	if err != nil {
		return queries, err
	}

	err = SetFieldValue(out, field.ForeignKey.AssociationFieldName, relation)
	if err != nil {
		return queries, err
	}

	return queries, nil
}

func preloadSingleMany(driver Driver, out interface{}, field Field) (Queries, error) {
	var queries Queries

	fk, err := GetFieldValueInt64(out, field.Schema.PrimaryKeyField.Name)
	if err != nil {
		return nil, err
	}

	if fk == int64(0) {
		return queries, nil
	}

	t := reflect.SliceOf(GetIndirectType(reflect.TypeOf(field.ForeignKey.Model)))
	relations := reflect.New(t)
	relations.Elem().Set(reflect.MakeSlice(t, 0, 0))

	q, err := FindByParamsWithQueries(driver, relations.Interface(), map[string]interface{}{field.ForeignKey.ColumnName: fk})
	queries = append(queries, q...)
	if err != nil {
		return queries, err
	}

	err = SetFieldValue(out, field.ForeignKey.Reference.AssociationFieldName, relations.Interface())
	if err != nil {
		return queries, err
	}

	return queries, nil
}

// ----------------------------------------------------------------------------
// Slice of instances preload
// ----------------------------------------------------------------------------

func preloadSlice(driver Driver, out interface{}, level int, fields []Field) (Queries, error) {
	var (
		queries Queries
		slc     reflect.Value
		value   = reflect.ValueOf(out)
	)

	if value.Kind() == reflect.Slice {
		slc = value
	} else {
		slc = value.Elem()
	}

	for _, field := range fields {
		if level > 1 {
			var (
				relations []interface{}
				mapping   = map[int64][]interface{}{}
			)

			//
			// Build relations preload slice
			//

			for i := 0; i < slc.Len(); i++ {
				instance := slc.Index(i).Interface()

				pk, err := GetFieldValueInt64(instance, field.Schema.PrimaryKeyField.Name)
				if err != nil {
					return queries, err
				}

				relation, err := GetFieldValue(instance, field.DestinationField)
				if err != nil {
					return queries, err
				}

				relationOut := Copy(relation)
				mapping[pk] = append(mapping[pk], relationOut)
				relations = append(relations, relationOut)
			}

			//
			// Preload
			//

			if field.IsAssociationTypeOne() {
				q, err := preloadSliceOne(driver, relations, field)
				queries = append(queries, q...)
				if err != nil {
					return queries, err
				}
			} else {
				q, err := preloadSliceMany(driver, relations, field)
				queries = append(queries, q...)
				if err != nil {
					return queries, err
				}
			}

			//
			// Set it back
			//

			for i := 0; i < slc.Len(); i++ {
				instance := slc.Index(i).Addr().Interface()

				pk, err := GetFieldValueInt64(instance, field.Schema.PrimaryKeyField.Name)
				if err != nil {
					return queries, err
				}

				instanceRelations := mapping[pk]

				if field.IsAssociationTypeOne() && len(instanceRelations) > 0 {
					err = SetFieldValue(instance, field.DestinationField, instanceRelations[0])
					if err != nil {
						return queries, err
					}
				}
			}
		} else {
			if field.IsAssociationTypeOne() {
				q, err := preloadSliceOne(driver, out, field)
				queries = append(queries, q...)
				if err != nil {
					return queries, err
				}
			} else {
				q, err := preloadSliceMany(driver, out, field)
				queries = append(queries, q...)
				if err != nil {
					return queries, err
				}
			}
		}
	}

	return queries, nil
}

func preloadSliceOne(driver Driver, out interface{}, field Field) (Queries, error) {
	var slc reflect.Value
	if reflect.ValueOf(out).Kind() == reflect.Slice {
		slc = reflect.ValueOf(out)
	} else {
		slc = reflect.ValueOf(out).Elem()
	}

	var (
		queries     Queries
		foreignKeys []int64
		mapping     = map[int64]map[int64]reflect.Value{} // pk -> fk -> pk instance value
	)

	//
	// Build mapping
	//

	for i := 0; i < slc.Len(); i++ {
		v := slc.Index(i)

		if v.Kind() == reflect.Interface {
			v = reflect.ValueOf(v.Interface())
		}

		if v.Kind() != reflect.Ptr && v.CanAddr() {
			v = v.Addr()
		}

		instance := v.Interface()

		pk, err := GetFieldValueInt64(instance, field.Schema.PrimaryKeyField.Name)
		if err != nil {
			return nil, err
		}

		fk, err := GetFieldValueInt64(instance, field.ForeignKey.FieldName)
		if err != nil {
			return nil, err
		}

		if fk != 0 && !InInt64Slice(foreignKeys, fk) {
			foreignKeys = append(foreignKeys, fk)
		}

		_, ok := mapping[pk]
		if !ok {
			mapping[pk] = map[int64]reflect.Value{}
		}

		mapping[pk][fk] = v
	}

	//
	// Perform queries (SELECT IN)
	//

	relationType := reflect.SliceOf(GetIndirectType(reflect.TypeOf(field.ForeignKey.Reference.Model)))
	relations := reflect.New(relationType)
	relations.Elem().Set(reflect.MakeSlice(relationType, 0, 0))

	q, err := FindByParamsWithQueries(driver, relations.Interface(), map[string]interface{}{field.ForeignKey.Reference.ColumnName: foreignKeys})
	queries = append(queries, q...)
	if err != nil {
		return queries, err
	}

	//
	// Iterate over instances and set related relation
	//

	relations = relations.Elem()

	for _, fkMap := range mapping {
		for i := 0; i < relations.Len(); i++ {
			var (
				relationValue = relations.Index(i).Addr()
				relation      = relationValue.Interface()
			)

			relationPK, err := GetFieldValueInt64(relation, field.ForeignKey.Schema.PrimaryKeyField.Name)
			if err != nil {
				return queries, err
			}

			instanceValue, ok := fkMap[relationPK]
			if !ok {
				continue
			}

			err = SetFieldValue(instanceValue.Interface(), field.ForeignKey.AssociationFieldName, relation)
			if err != nil {
				return queries, err
			}
		}
	}

	return queries, nil
}

func preloadSliceMany(driver Driver, out interface{}, field Field) (Queries, error) {
	var (
		slc         = reflect.ValueOf(out).Elem()
		queries     Queries
		foreignKeys []int64                     // As it's reversed, here foreign keys are instances primary keys
		mapping     = map[int64]reflect.Value{} // fk -> fk instance value
	)

	//
	// Build mapping
	//

	for i := 0; i < slc.Len(); i++ {
		instanceValue := slc.Index(i)

		if instanceValue.Kind() != reflect.Ptr && instanceValue.CanAddr() {
			instanceValue = instanceValue.Addr()
		}

		instance := instanceValue.Interface()

		fk, err := GetFieldValueInt64(instance, field.Schema.PrimaryKeyField.Name)
		if err != nil {
			return nil, err
		}

		if fk != 0 && !InInt64Slice(foreignKeys, fk) {
			foreignKeys = append(foreignKeys, fk)
			mapping[fk] = instanceValue
		}
	}

	//
	// Perform queries (SELECT IN)
	//

	relationType := reflect.SliceOf(GetIndirectType(reflect.TypeOf(field.ForeignKey.Model)))
	relations := reflect.New(relationType)
	relations.Elem().Set(reflect.MakeSlice(relationType, 0, 0))

	q, err := FindByParamsWithQueries(driver, relations.Interface(), map[string]interface{}{field.ForeignKey.ColumnName: foreignKeys})
	queries = append(queries, q...)
	if err != nil {
		return queries, err
	}

	//
	// Iterate over instances and set related relation
	//

	relations = relations.Elem()

	instancesRelations := map[int64][]reflect.Value{}

	for instancePK := range mapping {
		for i := 0; i < relations.Len(); i++ {
			var (
				relationValue = relations.Index(i).Addr()
				relation      = relationValue.Interface()
			)

			fk, err := GetFieldValueInt64(relation, field.ForeignKey.FieldName)
			if err != nil {
				return queries, err
			}

			if fk == instancePK {
				instancesRelations[instancePK] = append(instancesRelations[instancePK], relationValue)
			}
		}
	}

	for instancePK, instanceRelations := range instancesRelations {
		instanceValue := mapping[instancePK]

		t := reflect.SliceOf(GetIndirectType(reflect.TypeOf(field.ForeignKey.Model)))
		slc := reflect.New(t).Elem()
		slc.Set(reflect.MakeSlice(t, 0, 0))

		for _, relationValue := range instanceRelations {
			reflect.Append(slc, relationValue.Elem())
		}

		err := SetFieldValue(instanceValue.Interface(), field.ForeignKey.Reference.AssociationFieldName, slc.Interface())
		if err != nil {
			return queries, err
		}
	}

	return queries, nil
}
