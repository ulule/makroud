package reflectx_test

import (
	"reflect"
	"testing"
	"time"

	"database/sql"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx/reflectx"
)

type Elements struct {
	Spirit     bool
	Air        uint8
	Water      uint16
	Earth      uint32
	Fire       uint64
	Combustion *uint8
	Snow       *uint16
	Plant      *uint32
	Lava       *uint64
	Blast      int8
	Sound      int16
	Ice        int32
	Metal      int64
	Cristal    *int8
	Gas        *int16
	Acid       *int32
	Petrolium  *int64
	Sand       string
	Lightning  *string
	Vaccum     float32
	Blood      *float32
	Plastic    float64
	Glass      *float64
	Oxygen     pq.NullTime
	Poison     *pq.NullTime
	Vapor      sql.NullInt64
	Dust       *sql.NullInt64
	Laser      sql.NullString
	Deflect    *sql.NullString
	Corrosion  sql.NullBool
	Rubber     *sql.NullBool
	Fiber      sql.NullFloat64
	Dioxide    *sql.NullFloat64
	Gravity    ElementInterface
	Pressure   ElementStruct
	Push       []byte
	Absortion  []rune
	Fragture   []int64
	Friction   []string
	String     time.Time
	Tension    *time.Time
	xA         []byte
	xB         bool
	xC         string
	xD         int64
}

type ElementInterface interface {
}

type ElementStruct struct {
}

func TestReflectx_GetFields(t *testing.T) {
	is := require.New(t)

	expected := []string{
		"Spirit",
		"Air",
		"Water",
		"Earth",
		"Fire",
		"Combustion",
		"Snow",
		"Plant",
		"Lava",
		"Blast",
		"Sound",
		"Ice",
		"Metal",
		"Cristal",
		"Gas",
		"Acid",
		"Petrolium",
		"Sand",
		"Lightning",
		"Vaccum",
		"Blood",
		"Plastic",
		"Glass",
		"Oxygen",
		"Poison",
		"Vapor",
		"Dust",
		"Laser",
		"Deflect",
		"Corrosion",
		"Rubber",
		"Fiber",
		"Dioxide",
		"Gravity",
		"Pressure",
		"Push",
		"Absortion",
		"Fragture",
		"Friction",
		"String",
		"Tension",
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

	fields, err = reflectx.GetFields(Elements{}.Sand)
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
	is.Equal(reflect.Ptr, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "Snow")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Snow", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Ptr, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(&Elements{}, "Gravity")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Gravity", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Interface, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "Gravity")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Gravity", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Interface, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(&Elements{}, "Metal")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Metal", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Int64, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "Metal")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("Metal", field.Name)
	is.False(field.Anonymous)
	is.Empty(field.PkgPath)
	is.Equal(reflect.Int64, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(&Elements{}, "xA")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("xA", field.Name)
	is.False(field.Anonymous)
	is.NotEmpty(field.PkgPath)
	is.Equal(reflect.Slice, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "xB")
	is.True(ok)
	is.NotEmpty(field)
	is.Equal("xB", field.Name)
	is.False(field.Anonymous)
	is.NotEmpty(field.PkgPath)
	is.Equal(reflect.Bool, field.Type.Kind())

	field, ok = reflectx.GetFieldByName(Elements{}, "xE")
	is.False(ok)
	is.Empty(field)
}

func TestReflectx_GetFieldValue(t *testing.T) {
	is := require.New(t)

	a := false
	b := "89db"
	c := time.Now()
	d := uint8(3)
	e := sql.NullFloat64{
		Valid:   true,
		Float64: 3.333,
	}
	f := ElementStruct{}

	elements := &Elements{
		Spirit:   a,
		Air:      d,
		Tension:  &c,
		Sand:     b,
		Fiber:    e,
		Pressure: f,
		xA:       []byte("hello"),
	}

	value, err := reflectx.GetFieldValue(elements, "Spirit")
	is.NoError(err)
	is.Equal(false, value)

	value, err = reflectx.GetFieldValue(elements, "Air")
	is.NoError(err)
	is.Equal(uint8(3), value)

	value, err = reflectx.GetFieldValue(elements, "Tension")
	is.NoError(err)
	is.Equal(&c, value)

	value, err = reflectx.GetFieldValue(elements, "Sand")
	is.NoError(err)
	is.Equal("89db", value)

	value, err = reflectx.GetFieldValue(elements, "Pressure")
	is.NoError(err)
	is.Equal(f, value)

	value, err = reflectx.GetFieldValue(elements, "xA")
	is.Error(err)
	is.Nil(value)
}

func TestReflectx_UpdateFieldValue(t *testing.T) {
	is := require.New(t)

	v1 := false
	v2 := true
	v3 := uint8(10)
	v4 := uint8(20)
	v5 := int32(200)
	v6 := int64(2000)
	v7 := sql.NullInt64{
		Valid: true,
		Int64: int64(v3),
	}
	v8 := sql.NullInt64{
		Valid: true,
		Int64: v6,
	}
	v9 := "hello"
	v10 := "world"
	v11 := sql.NullString{
		Valid:  true,
		String: v9,
	}
	v12 := sql.NullString{
		Valid:  true,
		String: v10,
	}

	elements := Elements{}

	err := reflectx.UpdateFieldValue(elements, "Spirit", v1)
	is.Error(err)

	err = reflectx.UpdateFieldValue(elements, "Spirit", v2)
	is.Error(err)

	err = reflectx.UpdateFieldValue(&elements, "Spirit", v1)
	is.NoError(err)
	is.Equal(v1, elements.Spirit)

	err = reflectx.UpdateFieldValue(&elements, "Spirit", v2)
	is.NoError(err)
	is.Equal(v2, elements.Spirit)

	err = reflectx.UpdateFieldValue(&elements, "Air", v3)
	is.NoError(err)
	is.Equal(v3, elements.Air)

	err = reflectx.UpdateFieldValue(&elements, "Air", &v4)
	is.NoError(err)
	is.Equal(v4, elements.Air)

	err = reflectx.UpdateFieldValue(&elements, "Combustion", v3)
	is.NoError(err)
	is.Equal(&v3, elements.Combustion)

	err = reflectx.UpdateFieldValue(&elements, "Combustion", &v4)
	is.NoError(err)
	is.Equal(&v4, elements.Combustion)

	err = reflectx.UpdateFieldValue(&elements, "Ice", &v5)
	is.NoError(err)
	is.Equal(v5, elements.Ice)

	err = reflectx.UpdateFieldValue(&elements, "Ice", &v3)
	is.Error(err)
	is.Equal(v5, elements.Ice)

	err = reflectx.UpdateFieldValue(&elements, "Vapor", &v3)
	is.NoError(err)
	is.Equal(v7, elements.Vapor)

	err = reflectx.UpdateFieldValue(&elements, "Dust", v6)
	is.NoError(err)
	is.Equal(&v8, elements.Dust)

	err = reflectx.UpdateFieldValue(&elements, "Deflect", v9)
	is.NoError(err)
	is.Equal(&v11, elements.Deflect)

	err = reflectx.UpdateFieldValue(&elements, "Laser", v9)
	is.NoError(err)
	is.Equal(v11, elements.Laser)

	err = reflectx.UpdateFieldValue(&elements, "Laser", v12)
	is.NoError(err)
	is.Equal(v12, elements.Laser)

	err = reflectx.UpdateFieldValue(&elements, "xC", v10)
	is.Error(err)
	is.Empty(elements.xC)

}
