package sqlxx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetModelFromType(t *testing.T) {
	is := assert.New(t)

	results := []struct {
		value    interface{}
		expected interface{}
	}{
		{&Article{}, Article{}},
		{Article{}, Article{}},
	}

	for _, r := range results {
		actual := GetModelFromType(reflect.TypeOf(r.value))
		is.IsType(r.expected, actual)
	}
}
