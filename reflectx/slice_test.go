package reflectx_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/sqlxx/reflectx"
)

import "reflect"

func TestReflectx_IsSlice(t *testing.T) {
	is := require.New(t)

	si := []int{}
	si8 := []int8{}
	si16 := []int16{}
	si32 := []int32{}
	si64 := []int64{}
	sui := []uint{}
	sui8 := []uint8{}
	sui16 := []uint16{}
	sui32 := []uint32{}
	sui64 := []uint64{}
	sf32 := []float32{}
	sf64 := []float64{}
	sb := []bool{}
	ss := []string{}
	st := []time.Time{}
	sni64 := []sql.NullInt64{}
	snf64 := []sql.NullFloat64{}
	snb := []sql.NullBool{}
	sns := []sql.NullString{}
	snt := []pq.NullTime{}
	se := []Elements{}

	valids := []interface{}{
		si,
		&si,
		si8,
		&si8,
		si16,
		&si16,
		si32,
		&si32,
		si64,
		&si64,
		sui,
		&sui,
		sui8,
		&sui8,
		sui16,
		&sui16,
		sui32,
		&sui32,
		sui64,
		&sui64,
		sf32,
		&sf32,
		sf64,
		&sf64,
		sb,
		&sb,
		ss,
		&ss,
		st,
		&st,
		sni64,
		&sni64,
		snf64,
		&snf64,
		snb,
		&snb,
		sns,
		&sns,
		snt,
		&snt,
		se,
		&se,
	}

	for i, valid := range valids {
		v := reflectx.IsSlice(valid)
		is.True(v, fmt.Sprintf("loop #%d", i))
	}

	vi := int(6)
	vi8 := int8(6)
	vi16 := int16(6)
	vi32 := int32(6)
	vi64 := int64(6)
	vui := uint(6)
	vui8 := uint8(6)
	vui16 := uint16(6)
	vui32 := uint32(6)
	vui64 := uint64(6)
	vf32 := float32(6)
	vf64 := float64(6)
	vb := true
	vs := "foo"
	vt := time.Now()
	vni64 := sql.NullInt64{Valid: true, Int64: 6}
	vnf64 := sql.NullFloat64{Valid: true, Float64: 6}
	vnb := sql.NullBool{Valid: true, Bool: true}
	vns := sql.NullString{Valid: true, String: "foo"}
	vnt := pq.NullTime{Valid: true, Time: time.Now()}
	ve := Elements{xD: 6}

	invalids := []interface{}{
		vi,
		&vi,
		vi8,
		&vi8,
		vi16,
		&vi16,
		vi32,
		&vi32,
		vi64,
		&vi64,
		vui,
		&vui,
		vui8,
		&vui8,
		vui16,
		&vui16,
		vui32,
		&vui32,
		vui64,
		&vui64,
		vf32,
		&vf32,
		vf64,
		&vf64,
		vb,
		&vb,
		vs,
		&vs,
		vt,
		&vt,
		vni64,
		&vni64,
		vnf64,
		&vnf64,
		vnb,
		&vnb,
		vns,
		&vns,
		vnt,
		&vnt,
		ve,
		&ve,
	}

	for i, invalid := range invalids {
		v := reflectx.IsSlice(invalid)
		is.False(v, fmt.Sprintf("loop #%d", i))
	}

}

func TestReflectx_GetSliceType(t *testing.T) {
	is := require.New(t)

	si := []int{}
	si8 := []int8{}
	si16 := []int16{}
	si32 := []int32{}
	si64 := []int64{}
	sui := []uint{}
	sui8 := []uint8{}
	sui16 := []uint16{}
	sui32 := []uint32{}
	sui64 := []uint64{}
	sf32 := []float32{}
	sf64 := []float64{}
	sb := []bool{}
	ss := []string{}
	st := []time.Time{}
	sni64 := []sql.NullInt64{}
	snf64 := []sql.NullFloat64{}
	snb := []sql.NullBool{}
	sns := []sql.NullString{}
	snt := []pq.NullTime{}
	se := []Elements{}

	list := []struct {
		slice interface{}
		kind  interface{}
	}{
		{
			slice: si,
			kind:  int(0),
		},
		{
			slice: &si,
			kind:  int(0),
		},
		{
			slice: si8,
			kind:  int8(0),
		},
		{
			slice: &si8,
			kind:  int8(0),
		},
		{
			slice: si16,
			kind:  int16(0),
		},
		{
			slice: &si16,
			kind:  int16(0),
		},
		{
			slice: si32,
			kind:  int32(0),
		},
		{
			slice: &si32,
			kind:  int32(0),
		},
		{
			slice: si64,
			kind:  int64(0),
		},
		{
			slice: &si64,
			kind:  int64(0),
		},
		{
			slice: sui,
			kind:  uint(0),
		},
		{
			slice: &sui,
			kind:  uint(0),
		},
		{
			slice: sui8,
			kind:  uint8(0),
		},
		{
			slice: &sui8,
			kind:  uint8(0),
		},
		{
			slice: sui16,
			kind:  uint16(0),
		},
		{
			slice: &sui16,
			kind:  uint16(0),
		},
		{
			slice: sui32,
			kind:  uint32(0),
		},
		{
			slice: &sui32,
			kind:  uint32(0),
		},
		{
			slice: sui64,
			kind:  uint64(0),
		},
		{
			slice: &sui64,
			kind:  uint64(0),
		},
		{
			slice: sf32,
			kind:  float32(0),
		},
		{
			slice: &sf32,
			kind:  float32(0),
		},
		{
			slice: sf64,
			kind:  float64(0),
		},
		{
			slice: &sf64,
			kind:  float64(0),
		},
		{
			slice: sb,
			kind:  false,
		},
		{
			slice: &sb,
			kind:  false,
		},
		{
			slice: ss,
			kind:  "",
		},
		{
			slice: &ss,
			kind:  "",
		},
		{
			slice: st,
			kind:  time.Time{},
		},
		{
			slice: &st,
			kind:  time.Time{},
		},
		{
			slice: sni64,
			kind:  sql.NullInt64{},
		},
		{
			slice: &sni64,
			kind:  sql.NullInt64{},
		},
		{
			slice: snf64,
			kind:  sql.NullFloat64{},
		},
		{
			slice: &snf64,
			kind:  sql.NullFloat64{},
		},
		{
			slice: snb,
			kind:  sql.NullBool{},
		},
		{
			slice: &snb,
			kind:  sql.NullBool{},
		},
		{
			slice: sns,
			kind:  sql.NullString{},
		},
		{
			slice: &sns,
			kind:  sql.NullString{},
		},
		{
			slice: snt,
			kind:  pq.NullTime{},
		},
		{
			slice: &snt,
			kind:  pq.NullTime{},
		},
		{
			slice: se,
			kind:  Elements{},
		},
		{
			slice: &se,
			kind:  Elements{},
		},
	}

	for i := range list {
		v := reflectx.GetSliceType(list[i].slice)
		match := v == reflect.TypeOf(list[i].kind)
		is.True(match, fmt.Sprintf("loop #%d", i))
	}

}
