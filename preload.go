package sqlxx

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
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

	mapping := map[int][]Field{}
	var queries Queries

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
		fields := mapping[level]

		if level == 1 {
			for _, field := range fields {
				q, err := preloadOneToAssociation(driver, out, field)
				queries = append(queries, q...)
				if err != nil {
					return queries, err
				}
			}
		}

		if level == 2 {
			if IsSlice(out) {
				for _, field := range fields {
					q, err := preloadManyToAssociation(driver, out, field)
					queries = append(queries, q...)
					if err != nil {
						return queries, err
					}
				}
			} else {
				for _, field := range fields {
					value, err := GetFieldValue(reflect.ValueOf(out), field.DestinationField)
					if err != nil {
						return queries, err
					}

					reflected := reflect.ValueOf(value)
					isValue := false

					if !reflected.CanAddr() {
						value = Copy(value)
						isValue = true
					}

					q, err := Preload(driver, value, field.Name)
					queries = append(queries, q...)
					if err != nil {
						return queries, err
					}

					if isValue {
						value = reflect.Indirect(reflect.ValueOf(value)).Interface()
					}

					err = SetFieldValue(out, field.DestinationField, value)
					if err != nil {
						return queries, err
					}
				}
			}
		}
	}

	return queries, nil
}
