package reflekt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeSlice(t *testing.T) {
	is := assert.New(t)

	type Article struct{}

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
		actual := MakeSlice(r.value)
		is.IsType(r.expected, actual)
	}
}
