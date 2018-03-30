package sqlxx_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestMapper_MapBool(t *testing.T) {
	is := require.New(t)

	isActiveKey := "is_active"
	isActiveVal := true
	isActiveValPtr := &isActiveVal

	scenarios := []struct {
		mapper   sqlxx.Mapper
		valid    bool
		expected bool
	}{
		{
			mapper: map[string]interface{}{
				isActiveKey: nil,
			},
		},
		{
			mapper: map[string]interface{}{
				isActiveKey: isActiveVal,
			},
			valid:    true,
			expected: true,
		},
		{
			mapper: map[string]interface{}{
				isActiveKey: isActiveValPtr,
			},
			valid:    true,
			expected: true,
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := false
		handler := sqlxx.MapBool(isActiveKey, func(value bool) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.Equal(scenario.valid, ok, message)

		is.Equal(scenario.expected, actual, message)

	}

}

func TestMapper_MapNullBool(t *testing.T) {
	is := require.New(t)

	isActiveKey := "is_active"
	isActiveVal := true
	isActiveValPtr := &isActiveVal
	isActiveOpt := sql.NullBool{Valid: true, Bool: true}
	isActiveOptPtr := &isActiveOpt
	isActiveNull := sql.NullBool{Valid: false, Bool: false}
	isActiveNullPtr := &isActiveNull

	scenarios := []struct {
		mapper   sqlxx.Mapper
		expected sql.NullBool
	}{
		{
			mapper: map[string]interface{}{
				isActiveKey: nil,
			},
			expected: sql.NullBool{
				Valid: false,
				Bool:  false,
			},
		},
		{
			mapper: map[string]interface{}{
				isActiveKey: isActiveNull,
			},
			expected: sql.NullBool{
				Valid: false,
				Bool:  false,
			},
		},
		{
			mapper: map[string]interface{}{
				isActiveKey: isActiveNullPtr,
			},
			expected: sql.NullBool{
				Valid: false,
				Bool:  false,
			},
		},
		{
			mapper: map[string]interface{}{
				isActiveKey: isActiveVal,
			},
			expected: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		},
		{
			mapper: map[string]interface{}{
				isActiveKey: isActiveValPtr,
			},
			expected: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		},
		{
			mapper: map[string]interface{}{
				isActiveKey: isActiveOpt,
			},
			expected: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		},
		{
			mapper: map[string]interface{}{
				isActiveKey: isActiveOptPtr,
			},
			expected: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := sql.NullBool{}
		handler := sqlxx.MapNullBool(isActiveKey, func(value sql.NullBool) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.True(ok, message)

		is.Equal(scenario.expected, actual, message)

	}

}

func TestMapper_MapInt64(t *testing.T) {
	is := require.New(t)

	userIDKey := "user_id"
	userIDVal := int64(24)
	userIDValPtr := &userIDVal

	scenarios := []struct {
		mapper   sqlxx.Mapper
		valid    bool
		expected int64
	}{
		{
			mapper: map[string]interface{}{
				userIDKey: nil,
			},
		},
		{
			mapper: map[string]interface{}{
				userIDKey: userIDVal,
			},
			valid:    true,
			expected: 24,
		},
		{
			mapper: map[string]interface{}{
				userIDKey: userIDValPtr,
			},
			valid:    true,
			expected: 24,
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := int64(0)
		handler := sqlxx.MapInt64(userIDKey, func(value int64) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.Equal(scenario.valid, ok, message)

		is.Equal(scenario.expected, actual, message)

	}

}

func TestMapper_MapNullInt64(t *testing.T) {
	is := require.New(t)

	userIDKey := "user_id"
	userIDVal := int64(24)
	userIDValPtr := &userIDVal
	userIDOpt := sql.NullInt64{Valid: true, Int64: 24}
	userIDOptPtr := &userIDOpt
	userIDNull := sql.NullInt64{Valid: false, Int64: 0}
	userIDNullPtr := &userIDNull

	scenarios := []struct {
		mapper   sqlxx.Mapper
		expected sql.NullInt64
	}{
		{
			mapper: map[string]interface{}{
				userIDKey: nil,
			},
			expected: sql.NullInt64{
				Valid: false,
				Int64: 0,
			},
		},
		{
			mapper: map[string]interface{}{
				userIDKey: userIDNull,
			},
			expected: sql.NullInt64{
				Valid: false,
				Int64: 0,
			},
		},
		{
			mapper: map[string]interface{}{
				userIDKey: userIDNullPtr,
			},
			expected: sql.NullInt64{
				Valid: false,
				Int64: 0,
			},
		},
		{
			mapper: map[string]interface{}{
				userIDKey: userIDVal,
			},
			expected: sql.NullInt64{
				Valid: true,
				Int64: 24,
			},
		},
		{
			mapper: map[string]interface{}{
				userIDKey: userIDValPtr,
			},
			expected: sql.NullInt64{
				Valid: true,
				Int64: 24,
			},
		},
		{
			mapper: map[string]interface{}{
				userIDKey: userIDOpt,
			},
			expected: sql.NullInt64{
				Valid: true,
				Int64: 24,
			},
		},
		{
			mapper: map[string]interface{}{
				userIDKey: userIDOptPtr,
			},
			expected: sql.NullInt64{
				Valid: true,
				Int64: 24,
			},
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := sql.NullInt64{}
		handler := sqlxx.MapNullInt64(userIDKey, func(value sql.NullInt64) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.True(ok, message)

		is.Equal(scenario.expected, actual, message)

	}

}

func TestMapper_MapFloat64(t *testing.T) {
	is := require.New(t)

	matchKey := "match"
	matchVal := float64(18.8495559215)
	matchValPtr := &matchVal

	scenarios := []struct {
		mapper   sqlxx.Mapper
		valid    bool
		expected float64
	}{
		{
			mapper: map[string]interface{}{
				matchKey: nil,
			},
		},
		{
			mapper: map[string]interface{}{
				matchKey: matchVal,
			},
			valid:    true,
			expected: 18.8495559215,
		},
		{
			mapper: map[string]interface{}{
				matchKey: matchValPtr,
			},
			valid:    true,
			expected: 18.8495559215,
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := float64(0)
		handler := sqlxx.MapFloat64(matchKey, func(value float64) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.Equal(scenario.valid, ok, message)

		if scenario.valid {
			is.InEpsilon(scenario.expected, actual, 0.01, message)
		}

	}

}

func TestMapper_MapNullFloat64(t *testing.T) {
	is := require.New(t)

	matchKey := "match"
	matchVal := float64(18.8495559215)
	matchValPtr := &matchVal
	matchOpt := sql.NullFloat64{Valid: true, Float64: 18.8495559215}
	matchOptPtr := &matchOpt
	matchNull := sql.NullFloat64{Valid: false, Float64: 0}
	matchNullPtr := &matchNull

	scenarios := []struct {
		mapper   sqlxx.Mapper
		expected sql.NullFloat64
	}{
		{
			mapper: map[string]interface{}{
				matchKey: nil,
			},
			expected: sql.NullFloat64{
				Valid:   false,
				Float64: 0,
			},
		},
		{
			mapper: map[string]interface{}{
				matchKey: matchNull,
			},
			expected: sql.NullFloat64{
				Valid:   false,
				Float64: 0,
			},
		},
		{
			mapper: map[string]interface{}{
				matchKey: matchNullPtr,
			},
			expected: sql.NullFloat64{
				Valid:   false,
				Float64: 0,
			},
		},
		{
			mapper: map[string]interface{}{
				matchKey: matchVal,
			},
			expected: sql.NullFloat64{
				Valid:   true,
				Float64: 18.8495559215,
			},
		},
		{
			mapper: map[string]interface{}{
				matchKey: matchValPtr,
			},
			expected: sql.NullFloat64{
				Valid:   true,
				Float64: 18.8495559215,
			},
		},
		{
			mapper: map[string]interface{}{
				matchKey: matchOpt,
			},
			expected: sql.NullFloat64{
				Valid:   true,
				Float64: 18.8495559215,
			},
		},
		{
			mapper: map[string]interface{}{
				matchKey: matchOptPtr,
			},
			expected: sql.NullFloat64{
				Valid:   true,
				Float64: 18.8495559215,
			},
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := sql.NullFloat64{}
		handler := sqlxx.MapNullFloat64(matchKey, func(value sql.NullFloat64) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.True(ok, message)

		is.Equal(scenario.expected.Valid, actual.Valid, message)
		if actual.Valid {
			is.InEpsilon(scenario.expected.Float64, actual.Float64, 0.01, message)
		}

	}

}

func TestMapper_MapString(t *testing.T) {
	is := require.New(t)

	passwordKey := "password"
	passwordVal := "c2NyeXB0AA4AAAAIAAAAAa3LiE0YwIQWkQYes96CSWPYtqC/eFpfFU/EoTJwj"
	passwordValPtr := &passwordVal

	scenarios := []struct {
		mapper   sqlxx.Mapper
		valid    bool
		expected string
	}{
		{
			mapper: map[string]interface{}{
				passwordKey: nil,
			},
		},
		{
			mapper: map[string]interface{}{
				passwordKey: passwordVal,
			},
			valid:    true,
			expected: "c2NyeXB0AA4AAAAIAAAAAa3LiE0YwIQWkQYes96CSWPYtqC/eFpfFU/EoTJwj",
		},
		{
			mapper: map[string]interface{}{
				passwordKey: passwordValPtr,
			},
			valid:    true,
			expected: "c2NyeXB0AA4AAAAIAAAAAa3LiE0YwIQWkQYes96CSWPYtqC/eFpfFU/EoTJwj",
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := ""
		handler := sqlxx.MapString(passwordKey, func(value string) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.Equal(scenario.valid, ok, message)
		is.Equal(scenario.expected, actual, message)

	}

}

func TestMapper_MapNullString(t *testing.T) {
	is := require.New(t)

	passwordKey := "password"
	passwordVal := "c2NyeXB0AA4AAAAIAAAAAa3LiE0YwIQWkQYes96CSWPYtqC/eFpfFU/EoTJwj"
	passwordValPtr := &passwordVal
	passwordOpt := sql.NullString{Valid: true, String: "c2NyeXB0AA4AAAAIAAAAAa3LiE0YwIQWkQYes96CSWPYtqC/eFpfFU/EoTJwj"}
	passwordOptPtr := &passwordOpt
	passwordNull := sql.NullString{Valid: false, String: ""}
	passwordNullPtr := &passwordNull

	scenarios := []struct {
		mapper   sqlxx.Mapper
		expected sql.NullString
	}{
		{
			mapper: map[string]interface{}{
				passwordKey: nil,
			},
			expected: sql.NullString{
				Valid:  false,
				String: "",
			},
		},
		{
			mapper: map[string]interface{}{
				passwordKey: passwordNull,
			},
			expected: sql.NullString{
				Valid:  false,
				String: "",
			},
		},
		{
			mapper: map[string]interface{}{
				passwordKey: passwordNullPtr,
			},
			expected: sql.NullString{
				Valid:  false,
				String: "",
			},
		},
		{
			mapper: map[string]interface{}{
				passwordKey: passwordVal,
			},
			expected: sql.NullString{
				Valid:  true,
				String: "c2NyeXB0AA4AAAAIAAAAAa3LiE0YwIQWkQYes96CSWPYtqC/eFpfFU/EoTJwj",
			},
		},
		{
			mapper: map[string]interface{}{
				passwordKey: passwordValPtr,
			},
			expected: sql.NullString{
				Valid:  true,
				String: "c2NyeXB0AA4AAAAIAAAAAa3LiE0YwIQWkQYes96CSWPYtqC/eFpfFU/EoTJwj",
			},
		},
		{
			mapper: map[string]interface{}{
				passwordKey: passwordOpt,
			},
			expected: sql.NullString{
				Valid:  true,
				String: "c2NyeXB0AA4AAAAIAAAAAa3LiE0YwIQWkQYes96CSWPYtqC/eFpfFU/EoTJwj",
			},
		},
		{
			mapper: map[string]interface{}{
				passwordKey: passwordOptPtr,
			},
			expected: sql.NullString{
				Valid:  true,
				String: "c2NyeXB0AA4AAAAIAAAAAa3LiE0YwIQWkQYes96CSWPYtqC/eFpfFU/EoTJwj",
			},
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := sql.NullString{}
		handler := sqlxx.MapNullString(passwordKey, func(value sql.NullString) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.True(ok, message)

		is.Equal(scenario.expected, actual, message)

	}

}

func TestMapper_MapTime(t *testing.T) {
	is := require.New(t)

	snapshotKey := "snapshot"
	snapshotVal := time.Date(2017, 6, 6, 15, 0, 0, 0, time.UTC)
	snapshotValPtr := &snapshotVal

	scenarios := []struct {
		mapper   sqlxx.Mapper
		valid    bool
		expected time.Time
	}{
		{
			mapper: map[string]interface{}{
				snapshotKey: nil,
			},
		},
		{
			mapper: map[string]interface{}{
				snapshotKey: snapshotVal,
			},
			valid:    true,
			expected: time.Date(2017, 6, 6, 15, 0, 0, 0, time.UTC),
		},
		{
			mapper: map[string]interface{}{
				snapshotKey: snapshotValPtr,
			},
			valid:    true,
			expected: time.Date(2017, 6, 6, 15, 0, 0, 0, time.UTC),
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := time.Time{}
		handler := sqlxx.MapTime(snapshotKey, func(value time.Time) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.Equal(scenario.valid, ok, message)
		is.Equal(scenario.expected, actual, message)

	}

}

func TestMapper_MapNullTime(t *testing.T) {
	is := require.New(t)

	snapshotKey := "password"
	snapshotVal := time.Date(2017, 6, 6, 15, 0, 0, 0, time.UTC)
	snapshotValPtr := &snapshotVal
	snapshotOpt := pq.NullTime{Valid: true, Time: time.Date(2017, 6, 6, 15, 0, 0, 0, time.UTC)}
	snapshotOptPtr := &snapshotOpt
	snapshotNull := pq.NullTime{Valid: false, Time: time.Time{}}
	snapshotNullPtr := &snapshotNull

	scenarios := []struct {
		mapper   sqlxx.Mapper
		expected pq.NullTime
	}{
		{
			mapper: map[string]interface{}{
				snapshotKey: nil,
			},
			expected: pq.NullTime{
				Valid: false,
				Time:  time.Time{},
			},
		},
		{
			mapper: map[string]interface{}{
				snapshotKey: snapshotNull,
			},
			expected: pq.NullTime{
				Valid: false,
				Time:  time.Time{},
			},
		},
		{
			mapper: map[string]interface{}{
				snapshotKey: snapshotNullPtr,
			},
			expected: pq.NullTime{
				Valid: false,
				Time:  time.Time{},
			},
		},
		{
			mapper: map[string]interface{}{
				snapshotKey: snapshotVal,
			},
			expected: pq.NullTime{
				Valid: true,
				Time:  time.Date(2017, 6, 6, 15, 0, 0, 0, time.UTC),
			},
		},
		{
			mapper: map[string]interface{}{
				snapshotKey: snapshotValPtr,
			},
			expected: pq.NullTime{
				Valid: true,
				Time:  time.Date(2017, 6, 6, 15, 0, 0, 0, time.UTC),
			},
		},
		{
			mapper: map[string]interface{}{
				snapshotKey: snapshotOpt,
			},
			expected: pq.NullTime{
				Valid: true,
				Time:  time.Date(2017, 6, 6, 15, 0, 0, 0, time.UTC),
			},
		},
		{
			mapper: map[string]interface{}{
				snapshotKey: snapshotOptPtr,
			},
			expected: pq.NullTime{
				Valid: true,
				Time:  time.Date(2017, 6, 6, 15, 0, 0, 0, time.UTC),
			},
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #(%d)", (i + 1))

		actual := pq.NullTime{}
		handler := sqlxx.MapNullTime(snapshotKey, func(value pq.NullTime) {
			actual = value
		})

		ok, err := handler(scenario.mapper)
		is.NoError(err, message)
		is.True(ok, message)

		is.Equal(scenario.expected, actual, message)

	}

}
