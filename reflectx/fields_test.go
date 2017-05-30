package reflectx_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx/reflectx"
)

type Element interface {
	Name() string
}

type Elements struct {
	Spirit     bool
	Air        string
	Fire       string
	Water      string
	Earth      string
	Combustion int64
	Lava       int64
	Plant      int64
	Snow       int64
	Blast      interface{}
	Metal      interface{}
	Ice        interface{}
	Sound      interface{}
	lightning  []byte
	sand       []byte
	blood      []byte
	vaccum     []byte
}

func TestReflectx_GetFields(t *testing.T) {
	is := require.New(t)

	expected := []string{
		"Spirit",
		"Air",
		"Fire",
		"Water",
		"Earth",
		"Combustion",
		"Lava",
		"Plant",
		"Snow",
		"Blast",
		"Metal",
		"Ice",
		"Sound",
	}

	fields, err := reflectx.GetFields(&Elements{})
	is.NoError(err)
	is.NotEmpty(fields)
	is.Equal(expected, fields)

	fields, err = reflectx.GetFields(Elements{})
	is.NoError(err)
	is.NotEmpty(fields)
	is.Equal(expected, fields)

	fields, err = reflectx.GetFields(func() bool {
		return false
	})
	is.Error(err)
	is.Empty(fields)

	fields, err = reflectx.GetFields("hello world!")
	is.Error(err)
	is.Empty(fields)

	fields, err = reflectx.GetFields(Elements{}.sand)
	is.Error(err)
	is.Empty(fields)
}

func TestReflectx_GetFieldByName(t *testing.T) {
	is := require.New(t)

	field, ok := reflectx.GetFieldByName(&Elements{}, "Snow")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Snow", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Int64, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "Snow")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Snow", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Int64, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(&Elements{}, "Blast")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Blast", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Interface, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "Blast")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Blast", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Interface, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(&Elements{}, "lightning")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("lightning", field.Name)
	is.False(field.Anonymous)
	is.NotEmpty(field.PkgPath)
	is.Equal(reflect.Slice, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "lightning")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("lightning", field.Name)
	is.False(field.Anonymous)
	is.NotEmpty(field.PkgPath)
	is.Equal(reflect.Slice, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "gravity")
	is.False(ok)
	is.Empty(field)
}

func TestReflectx_GetFieldValue(t *testing.T) {
	is := require.New(t)

	elements := &Elements{
		Spirit: false,
		Air:    "air",
		Snow:   546,
		Sound:  "89db",
		blood:  []byte("hello"),
	}

	value, err := reflectx.GetFieldValue(elements, "Spirit")
	is.NoError(err)
	is.Equal(false, value)

	value, err = reflectx.GetFieldValue(elements, "Air")
	is.NoError(err)
	is.Equal("air", value)

	value, err = reflectx.GetFieldValue(elements, "Snow")
	is.NoError(err)
	is.Equal(int64(546), value)

	value, err = reflectx.GetFieldValue(elements, "Sound")
	is.NoError(err)
	is.Equal("89db", value)

	value, err = reflectx.GetFieldValue(elements, "blood")
	is.Error(err)
	is.Nil(value)
}
