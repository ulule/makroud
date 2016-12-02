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

func TestGetFieldValues(t *testing.T) {
	is := assert.New(t)

	type Article struct {
		ID   int
		Name string
	}

	articles := []Article{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "two"},
	}

	ids, err := GetFieldValues(&articles, "ID")
	is.Nil(err)
	is.Equal([]interface{}{1, 2}, ids)

	ids, err = GetFieldValues(&articles, "Name")
	is.Nil(err)
	is.Equal([]interface{}{"one", "two"}, ids)

	ids, err = GetFieldValues(&articles[0], "Name")
	is.Nil(err)
	is.Equal([]interface{}{"one"}, ids)
}
