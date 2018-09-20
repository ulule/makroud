package reflectx_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ulule/makroud/reflectx"
)

func TestReflectx_ToInt64(t *testing.T) {
	is := require.New(t)

	i := int(1)
	i8 := int8(1)
	i16 := int16(1)
	i32 := int32(1)
	i64 := int64(1)
	ui := uint(1)
	ui8 := uint8(1)
	ui16 := uint16(1)
	ui32 := uint32(1)
	ui64 := uint64(1)
	f32 := float32(1)
	f64 := float64(1)
	ni64 := sql.NullInt64{Valid: true, Int64: 1}
	nf64 := sql.NullFloat64{Valid: true, Float64: float64(1)}

	valids := []interface{}{
		i,
		&i,
		i8,
		&i8,
		i16,
		&i16,
		i32,
		&i32,
		i64,
		&i64,
		ui,
		&ui,
		ui8,
		&ui8,
		ui16,
		&ui16,
		ui32,
		&ui32,
		ui64,
		&ui64,
		f32,
		&f32,
		f64,
		&f64,
		ni64,
		&ni64,
		nf64,
		&nf64,
	}

	for i, valid := range valids {
		m := fmt.Sprintf("loop #%d", i)

		v, err := reflectx.ToInt64(valid)
		is.NoError(err, m)
		is.Equal(int64(1), v, m)

		v, ok := reflectx.ToOptionalInt64(valid)
		is.True(ok, m)
		is.Equal(int64(1), v, m)
	}

	str := "hello"
	type A struct{}
	foo := true
	zni64 := sql.NullInt64{}
	znf64 := sql.NullFloat64{}
	nf64 = sql.NullFloat64{Valid: true, Float64: float64(1.5)}

	invalids := []interface{}{
		nil,
		str,
		&str,
		reflect.ValueOf(1),
		A{},
		&A{},
		foo,
		&foo,
		[]A{},
		&[]A{},
		zni64,
		&zni64,
		znf64,
		&znf64,
	}

	for i, invalid := range invalids {
		m := fmt.Sprintf("loop #%d", i)

		v, err := reflectx.ToInt64(invalid)
		is.Error(err, m)
		is.Equal(int64(0), v, m)

		v, ok := reflectx.ToOptionalInt64(invalid)
		is.False(ok, m)
		is.Equal(int64(0), v, m)
	}
}

func TestReflectx_ToString(t *testing.T) {
	is := require.New(t)

	str := "c"
	run := 'c'
	ns := sql.NullString{Valid: true, String: "c"}

	valids := []interface{}{
		str,
		&str,
		run,
		&run,
		ns,
		&ns,
	}

	for i, valid := range valids {
		m := fmt.Sprintf("loop #%d", i)

		v, err := reflectx.ToString(valid)
		is.NoError(err, m)
		is.Equal("c", v, m)

		v, ok := reflectx.ToOptionalString(valid)
		is.True(ok, m)
		is.Equal("c", v, m)
	}

	type A struct{}
	foo := true
	zns := sql.NullString{}

	invalids := []interface{}{
		nil,
		reflect.ValueOf("c"),
		A{},
		&A{},
		foo,
		&foo,
		[]A{},
		&[]A{},
		zns,
		&zns,
	}

	for i, invalid := range invalids {
		m := fmt.Sprintf("loop #%d", i)

		v, err := reflectx.ToString(invalid)
		is.Error(err, m)
		is.Equal("", v, m)

		v, ok := reflectx.ToOptionalString(invalid)
		is.False(ok, m)
		is.Equal("", v, m)
	}
}

func TestReflectx_GetType(t *testing.T) {
	is := require.New(t)

	type foo struct {
		value int
	}

	i := int(1)
	i8 := int8(1)
	i16 := int16(1)
	i32 := int32(1)
	i64 := int64(1)
	ui := uint(1)
	ui8 := uint8(1)
	ui16 := uint16(1)
	ui32 := uint32(1)
	ui64 := uint64(1)
	f32 := float32(1)
	f64 := float64(1)
	ni64 := sql.NullInt64{Valid: true, Int64: 1}
	nf64 := sql.NullFloat64{Valid: true, Float64: float64(1)}
	str := "c"
	ns := sql.NullString{Valid: true, String: "c"}
	x := true
	e := foo{value: 3}
	t0 := time.Now()

	scenarios := []struct {
		input    reflect.Type
		expected reflectx.Type
	}{
		{
			input:    reflect.TypeOf(i),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&i),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(i8),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&i8),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(i16),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&i16),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(i32),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&i32),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(i64),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&i64),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(ui),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&ui),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(ui8),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&ui8),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(ui16),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&ui16),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(ui32),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&ui32),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(ui64),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(&ui64),
			expected: reflectx.Int64Type,
		},
		{
			input:    reflect.TypeOf(f32),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(&f32),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(f64),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(&f64),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(ni64),
			expected: reflectx.OptionalInt64Type,
		},
		{
			input:    reflect.TypeOf(&ni64),
			expected: reflectx.OptionalInt64Type,
		},
		{
			input:    reflect.TypeOf(nf64),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(&nf64),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(str),
			expected: reflectx.StringType,
		},
		{
			input:    reflect.TypeOf(&str),
			expected: reflectx.StringType,
		},
		{
			input:    reflect.TypeOf(ns),
			expected: reflectx.OptionalStringType,
		},
		{
			input:    reflect.TypeOf(&ns),
			expected: reflectx.OptionalStringType,
		},
		{
			input:    reflect.TypeOf(x),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(&x),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(e),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(&e),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(t0),
			expected: reflectx.UnsupportedType,
		},
		{
			input:    reflect.TypeOf(&t0),
			expected: reflectx.UnsupportedType,
		},
	}

	for i, scenario := range scenarios {
		m := fmt.Sprintf("loop #%d", i)
		v := reflectx.GetType(scenario.input)
		is.Equal(scenario.expected, v, m)
	}
}
