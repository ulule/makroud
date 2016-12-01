package sqlxx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeModel(t *testing.T) {
	is := assert.New(t)

	results := []struct {
		value    interface{}
		expected interface{}
	}{
		{&Article{}, Article{}},
		{Article{}, Article{}},
	}

	for _, r := range results {
		actual := makeModel(reflect.TypeOf(r.value))
		is.IsType(r.expected, actual)
	}
}

func TestMakeSlice(t *testing.T) {
	is := assert.New(t)

	articles := []Article{}
	articlesPtrs := []*Article{}

	results := []struct {
		value    interface{}
		expected interface{}
	}{
		{articles, []Article{}},
		{&articles, []Article{}},
		{articlesPtrs, []Article{}},
	}

	for _, r := range results {
		actual := makeSlice(r.value)
		is.IsType(r.expected, actual)
	}
}
