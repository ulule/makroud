package reflectx_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx/reflectx"
)

func TestReflectx_MakePointer(t *testing.T) {
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
		reflectx.MakePointer(any),
		reflectx.MakePointer(anyPtr),
	}

	for _, result := range results {
		elem, ok := result.(*anyType)
		is.True(ok)
		is.Equal(1, elem.value)
		is.Equal(reflect.ValueOf(result).Kind(), reflect.Ptr)
		is.Equal(reflect.ValueOf(result).Type().Elem(), reflect.TypeOf(anyType{}))
	}

	anyWithEmbed := anyType{value: 1, embed: embedType{value: 2}}
	anyWithEmbedPtr := anyType{value: 1, embedPtr: &embedType{value: 2}}

	results = []interface{}{
		reflectx.MakePointer(anyWithEmbed.embed),
		reflectx.MakePointer(anyWithEmbedPtr.embedPtr),
	}

	for _, result := range results {
		elem, ok := result.(*embedType)
		is.True(ok)
		is.Equal(2, elem.value)
		is.Equal(reflect.ValueOf(result).Kind(), reflect.Ptr)
		is.Equal(reflect.ValueOf(result).Type().Elem(), reflect.TypeOf(embedType{}))
	}

}
