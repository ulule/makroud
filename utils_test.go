package sqlxx_test

import (
	"reflect"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestIntToInt64(t *testing.T) {
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
	}

	for _, valid := range valids {
		v, err := sqlxx.IntToInt64(valid)
		assert.NoError(t, err)
		assert.Equal(t, v, int64(1))
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
		assert.Error(t, err)
		assert.Equal(t, int64(0), v)
	}
}