package reflectx_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx/reflectx"
)

func TestReflectx_GetIndirectType(t *testing.T) {
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
	run := 'c'
	ns := sql.NullString{Valid: true, String: "c"}
	e := foo{value: 3}

	scenarios := []struct {
		input  interface{}
		output reflect.Type
	}{
		{
			input:  i,
			output: reflect.TypeOf(i),
		},
		{
			input:  &i,
			output: reflect.TypeOf(i),
		},
		{
			input:  i8,
			output: reflect.TypeOf(i8),
		},
		{
			input:  &i8,
			output: reflect.TypeOf(i8),
		},
		{
			input:  i16,
			output: reflect.TypeOf(i16),
		},
		{
			input:  &i16,
			output: reflect.TypeOf(i16),
		},
		{
			input:  i32,
			output: reflect.TypeOf(i32),
		},
		{
			input:  &i32,
			output: reflect.TypeOf(i32),
		},
		{
			input:  i64,
			output: reflect.TypeOf(i64),
		},
		{
			input:  &i64,
			output: reflect.TypeOf(i64),
		},
		{
			input:  ui,
			output: reflect.TypeOf(ui),
		},
		{
			input:  &ui,
			output: reflect.TypeOf(ui),
		},
		{
			input:  ui8,
			output: reflect.TypeOf(ui8),
		},
		{
			input:  &ui8,
			output: reflect.TypeOf(ui8),
		},
		{
			input:  ui16,
			output: reflect.TypeOf(ui16),
		},
		{
			input:  &ui16,
			output: reflect.TypeOf(ui16),
		},
		{
			input:  ui32,
			output: reflect.TypeOf(ui32),
		},
		{
			input:  &ui32,
			output: reflect.TypeOf(ui32),
		},
		{
			input:  ui64,
			output: reflect.TypeOf(ui64),
		},
		{
			input:  &ui64,
			output: reflect.TypeOf(ui64),
		},
		{
			input:  f32,
			output: reflect.TypeOf(f32),
		},
		{
			input:  &f32,
			output: reflect.TypeOf(f32),
		},
		{
			input:  f64,
			output: reflect.TypeOf(f64),
		},
		{
			input:  &f64,
			output: reflect.TypeOf(f64),
		},
		{
			input:  ni64,
			output: reflect.TypeOf(ni64),
		},
		{
			input:  &ni64,
			output: reflect.TypeOf(ni64),
		},
		{
			input:  nf64,
			output: reflect.TypeOf(nf64),
		},
		{
			input:  &nf64,
			output: reflect.TypeOf(nf64),
		},
		{
			input:  str,
			output: reflect.TypeOf(str),
		},
		{
			input:  &str,
			output: reflect.TypeOf(str),
		},
		{
			input:  run,
			output: reflect.TypeOf(run),
		},
		{
			input:  &run,
			output: reflect.TypeOf(run),
		},
		{
			input:  ns,
			output: reflect.TypeOf(ns),
		},
		{
			input:  &ns,
			output: reflect.TypeOf(ns),
		},
		{
			input:  e,
			output: reflect.TypeOf(e),
		},
		{
			input:  &e,
			output: reflect.TypeOf(e),
		},
	}

	for i, scenario := range scenarios {
		v := reflectx.GetIndirectType(scenario.input)
		is.Equal(scenario.output, v, fmt.Sprintf("loop #%d", i))
	}

}

func TestReflectx_GetFlattenValue(t *testing.T) {
	is := require.New(t)

	type foo struct {
		value int
	}

	e1 := foo{value: 3}
	e2 := &e1
	e3 := &e2
	e4 := &e3
	e5 := &e4
	e6 := &e5

	l1 := []foo{{value: 3}}
	l2 := &l1
	l3 := &l2
	l4 := &l3
	l5 := &l4
	l6 := &l5

	scenarios := []struct {
		input  interface{}
		output interface{}
	}{
		{
			input:  e1,
			output: foo{value: 3},
		},
		{
			input:  e2,
			output: &foo{value: 3},
		},
		{
			input:  e3,
			output: &foo{value: 3},
		},
		{
			input:  e4,
			output: &foo{value: 3},
		},
		{
			input:  e5,
			output: &foo{value: 3},
		},
		{
			input:  e6,
			output: &foo{value: 3},
		},
		{
			input:  l1,
			output: []foo{{value: 3}},
		},
		{
			input:  l2,
			output: &[]foo{{value: 3}},
		},
		{
			input:  l3,
			output: &[]foo{{value: 3}},
		},
		{
			input:  l4,
			output: &[]foo{{value: 3}},
		},
		{
			input:  l5,
			output: &[]foo{{value: 3}},
		},
		{
			input:  l6,
			output: &[]foo{{value: 3}},
		},
	}

	for i, scenario := range scenarios {
		v := reflectx.GetFlattenValue(scenario.input)
		is.Equal(scenario.output, v, fmt.Sprintf("loop #%d", i))
	}

}
