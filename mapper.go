package sqlxx

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

type Mapper map[string]interface{}
type MapHandler func(mapper Mapper) (bool, error)

func MapBool(key string, callback func(value bool)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}
		if raw == nil {
			return false, nil
		}

		value := false
		switch e := raw.(type) {
		case bool:
			value = e
		case *bool:
			value = *e
		default:
			return false, errors.New("cannot map value as 'bool'")
		}

		callback(value)
		return true, nil
	}
}

// nolint: dupl
func MapNullBool(key string, callback func(value sql.NullBool)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}

		value := sql.NullBool{}
		if raw == nil {
			callback(value)
			return true, nil
		}

		switch e := raw.(type) {
		case bool:
			value.Valid = true
			value.Bool = e
		case *bool:
			value.Valid = true
			value.Bool = *e
		case sql.NullBool:
			value = e
		case *sql.NullBool:
			value = *e
		default:
			return false, errors.New("cannot map value as optional 'bool'")
		}

		callback(value)
		return true, nil
	}
}

// nolint: dupl
func MapInt64(key string, callback func(value int64)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}
		if raw == nil {
			return false, nil
		}

		value := int64(0)
		switch e := raw.(type) {
		case int64:
			value = e
		case *int64:
			value = *e
		default:
			return false, errors.New("cannot map value as 'int64'")
		}

		callback(value)
		return true, nil
	}
}

// nolint: dupl
func MapNullInt64(key string, callback func(value sql.NullInt64)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}

		value := sql.NullInt64{}

		if raw == nil {
			callback(value)
			return true, nil
		}

		switch e := raw.(type) {
		case int64:
			value.Valid = true
			value.Int64 = e
		case *int64:
			value.Valid = true
			value.Int64 = *e
		case sql.NullInt64:
			value = e
		case *sql.NullInt64:
			value = *e
		default:
			return false, errors.New("cannot map value as optional 'int64'")
		}

		callback(value)
		return true, nil
	}
}

// nolint: dupl
func MapFloat64(key string, callback func(value float64)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}
		if raw == nil {
			return false, nil
		}

		value := float64(0)
		switch e := raw.(type) {
		case float64:
			value = e
		case *float64:
			value = *e
		default:
			return false, errors.New("cannot map value as 'float64'")
		}

		callback(value)
		return true, nil
	}
}

// nolint: dupl
func MapNullFloat64(key string, callback func(value sql.NullFloat64)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}

		value := sql.NullFloat64{}

		if raw == nil {
			callback(value)
			return true, nil
		}

		switch e := raw.(type) {
		case float64:
			value.Valid = true
			value.Float64 = e
		case *float64:
			value.Valid = true
			value.Float64 = *e
		case sql.NullFloat64:
			value = e
		case *sql.NullFloat64:
			value = *e
		default:
			return false, errors.New("cannot map value as optional 'float64'")
		}

		callback(value)
		return true, nil
	}
}

// nolint: dupl
func MapString(key string, callback func(value string)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}
		if raw == nil {
			return false, nil
		}

		value := ""
		switch e := raw.(type) {
		case string:
			value = e
		case *string:
			value = *e
		default:
			return false, errors.New("cannot map value as 'string'")
		}

		callback(value)
		return true, nil
	}
}

// nolint: dupl
func MapNullString(key string, callback func(value sql.NullString)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}

		value := sql.NullString{}

		if raw == nil {
			callback(value)
			return true, nil
		}

		switch e := raw.(type) {
		case string:
			value.Valid = true
			value.String = e
		case *string:
			value.Valid = true
			value.String = *e
		case sql.NullString:
			value = e
		case *sql.NullString:
			value = *e
		default:
			return false, errors.New("cannot map value as optional 'string'")
		}

		callback(value)
		return true, nil
	}
}

// nolint: dupl
func MapTime(key string, callback func(value time.Time)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}
		if raw == nil {
			return false, nil
		}

		value := time.Time{}
		switch e := raw.(type) {
		case time.Time:
			value = e
		case *time.Time:
			value = *e
		default:
			return false, errors.New("cannot map value as 'time'")
		}

		callback(value)
		return true, nil
	}
}

// nolint: dupl
func MapNullTime(key string, callback func(value pq.NullTime)) MapHandler {
	return func(mapper Mapper) (bool, error) {
		raw, ok := mapper[key]
		if !ok {
			return false, nil
		}

		value := pq.NullTime{}

		if raw == nil {
			callback(value)
			return true, nil
		}

		switch e := raw.(type) {
		case time.Time:
			value.Valid = true
			value.Time = e
		case *time.Time:
			value.Valid = true
			value.Time = *e
		case pq.NullTime:
			value = e
		case *pq.NullTime:
			value = *e
		default:
			return false, errors.New("cannot map value as optional 'time'")
		}

		callback(value)
		return true, nil
	}
}

func Map(mapper Mapper, handlers ...MapHandler) error {
	for _, handler := range handlers {
		_, err := handler(mapper)
		if err != nil {
			return err
		}
	}
	return nil
}
