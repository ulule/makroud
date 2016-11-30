package sqlxx

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTags(t *testing.T) {
	var (
		field reflect.StructField
		tags  Tags
	)

	is := assert.New(t)

	s := struct {
		ID int `db:"foo_field" sqlxx:"primary_key; ignored; foo:bar; hello : you; bool:false;    id:1"`
	}{1}

	field, _ = reflect.TypeOf(s).FieldByName("ID")
	tags = makeTags(field)
	is.Equal("foo_field", tags.GetByKey(SQLXStructTagName, "field"))
	is.Equal("true", tags.GetByKey(StructTagName, "primary_key"))
	is.Equal("true", tags.GetByKey(StructTagName, "ignored"))
	is.Equal("bar", tags.GetByKey(StructTagName, "foo"))
	is.Equal("you", tags.GetByKey(StructTagName, "hello"))
	is.Equal("false", tags.GetByKey(StructTagName, "bool"))
	is.Equal("1", tags.GetByKey(StructTagName, "id"))
}
