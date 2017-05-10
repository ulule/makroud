package sqlxx_test

import (
	"reflect"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx"
)

func TestModel_TypeToModel(t *testing.T) {
	results := []struct {
		value    interface{}
		expected interface{}
	}{
		{&Article{}, Article{}},
		{Article{}, Article{}},
	}

	for _, r := range results {
		actual := sqlxx.ToModel(reflect.TypeOf(r.value))
		assert.IsType(t, r.expected, actual)
	}
}
