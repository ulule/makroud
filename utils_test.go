package sqlxx_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestUtils_IntToInt64(t *testing.T) {
	is := require.New(t)

	valids := []interface{}{
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		float32(1),
		float64(1),
		sql.NullInt64{Valid: true, Int64: 1},
	}

	for _, valid := range valids {
		v, err := sqlxx.IntToInt64(valid)
		is.NoError(err)
		is.Equal(v, int64(1))
	}

	str := "hello"
	type A struct{}

	invalids := []interface{}{
		nil,
		str,
		&str,
		reflect.ValueOf(1),
		A{},
		&A{},
	}

	for _, invalid := range invalids {
		v, err := sqlxx.IntToInt64(invalid)
		is.Error(err)
		is.Equal(int64(0), v)
	}
}

func TestUtils_MakePointer(t *testing.T) {
	is := require.New(t)

	type embedType struct {
		value int
	}

	type anyType struct {
		value    int
		embed    embedType
		embedPtr *embedType
	}

	any := anyType{value: 1}
	anyPtr := &anyType{value: 1}

	results := []interface{}{
		sqlxx.MakePointer(any),
		sqlxx.MakePointer(anyPtr),
	}

	for _, r := range results {
		is.Equal(1, r.(*anyType).value)
		is.Equal(reflect.ValueOf(r).Kind(), reflect.Ptr)
		is.Equal(reflect.ValueOf(r).Type().Elem(), reflect.TypeOf(anyType{}))
	}

	anyWithEmbed := anyType{value: 1, embed: embedType{value: 2}}
	anyWithEmbedPtr := anyType{value: 1, embedPtr: &embedType{value: 2}}

	results = []interface{}{
		sqlxx.MakePointer(anyWithEmbed.embed),
		sqlxx.MakePointer(anyWithEmbedPtr.embedPtr),
	}

	for _, r := range results {
		is.Equal(2, r.(*embedType).value)
		is.Equal(reflect.ValueOf(r).Kind(), reflect.Ptr)
		is.Equal(reflect.ValueOf(r).Type().Elem(), reflect.TypeOf(embedType{}))
	}
}
