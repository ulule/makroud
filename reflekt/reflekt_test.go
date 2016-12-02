package reflekt

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeSlice(t *testing.T) {
	is := assert.New(t)

	type Article struct{}

	results := []struct {
		value    interface{}
		expected interface{}
	}{
		{Article{}, []Article{}},
		{&Article{}, []Article{}},
	}

	for _, r := range results {
		is.IsType(r.expected, MakeSlice(r.value))
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

func TestGetFieldTags(t *testing.T) {
	is := assert.New(t)

	tagged := struct {
		ID int `db:"foo_field" sqlxx:"primary_key; ignored; foo:bar; hello : you; bool:false;    id:1"`
	}{1}

	field, _ := reflect.TypeOf(tagged).FieldByName("ID")
	tags := GetFieldTags(field, []string{"db", "sqlxx"}, map[string]string{"db": "field"})

	is.Equal("foo_field", tags.GetByKey("db", "field"))
	is.Equal("true", tags.GetByKey("sqlxx", "primary_key"))
	is.Equal("true", tags.GetByKey("sqlxx", "ignored"))
	is.Equal("bar", tags.GetByKey("sqlxx", "foo"))
	is.Equal("you", tags.GetByKey("sqlxx", "hello"))
	is.Equal("false", tags.GetByKey("sqlxx", "bool"))
	is.Equal("1", tags.GetByKey("sqlxx", "id"))
}
