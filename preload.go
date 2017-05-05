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
func Preload(driver Driver, out interface{}, paths ...string) (Queries, error) {
	if !reflect.Indirect(reflect.ValueOf(out)).CanAddr() {
		return nil, errors.New("model instance must be addressable (pointer required)")
	}

	schema, err := GetSchema(out)
	if err != nil {
		return nil, err
	}

	var (
		assocs         []Field
		assocsOfAssocs = map[string]Field{}
		queries        Queries
	)

	for _, path := range paths {
		assoc, ok := schema.Associations[path]
		if !ok {
			return nil, fmt.Errorf("%s is not a valid association", path)
		}

		splits := strings.Split(path, ".")

		if len(splits) == 1 {
			assocs = append(assocs, assoc)
		}

		if len(splits) == 2 {
			assocsOfAssocs[splits[0]] = assoc
		}
	}

	q, err := preloadOneToAssociations(driver, out, assocs)
	queries = append(queries, q...)
	if err != nil {
		return queries, err
	}

	if IsSlice(out) {
		for k, v := range assocsOfAssocs {
			q, err = preloadManyToAssociations(driver, out, schema, k, v)
			queries = append(queries, q...)
			if err != nil {
				return queries, err
			}

		}
		return queries, nil
	}

	for k, v := range assocsOfAssocs {
		value, err := GetFieldValue(reflect.ValueOf(out), k)
		if err != nil {
			return queries, err
		}

		reflected := reflect.ValueOf(value)
		isValue := false

		if !reflected.CanAddr() {
			value = Copy(value)
			isValue = true
		}

		q, err = Preload(driver, value, v.Name)
		queries = append(queries, q...)
		if err != nil {
			return queries, err
		}

		if isValue {
			value = reflect.Indirect(reflect.ValueOf(value)).Interface()
		}

		err = SetFieldValue(out, k, value)
		if err != nil {
			return queries, err
		}
	}

	return queries, nil
}
