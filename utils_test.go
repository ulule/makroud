package sqlxx_test

import (
	"database/sql"
	"reflect"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestUtils_GetPrimaryKeys(t *testing.T) {
	_, fixtures, shutdown := dbConnection(t)
	defer shutdown()

	pks, err := sqlxx.GetPrimaryKeys(&fixtures.Articles, "ID")
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{1, 2, 3, 4, 5}, pks)

	pks, err = sqlxx.GetPrimaryKeys(&fixtures.Articles[0], "ID")
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{1}, pks)
}

func TestUtils_IntToInt64(t *testing.T) {
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
