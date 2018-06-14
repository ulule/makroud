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

func TestReflectx_IsZero(t *testing.T) {
	is := require.New(t)

	zi := int(0)
	zi8 := int8(0)
	zi16 := int16(0)
	zi32 := int32(0)
	zi64 := int64(0)
	zui := uint(0)
	zui8 := uint8(0)
	zui16 := uint16(0)
	zui32 := uint32(0)
	zui64 := uint64(0)
	zf32 := float32(0)
	zf64 := float64(0)
	zb := false
	zs := ""
	zt := time.Time{}
	zni64 := sql.NullInt64{}
	znf64 := sql.NullFloat64{}
	znb := sql.NullBool{}
	zns := sql.NullString{}
	znt := pq.NullTime{}

	valids := []interface{}{
		zi,
		&zi,
		zi8,
		&zi8,
		zi16,
		&zi16,
		zi32,
		&zi32,
		zi64,
		&zi64,
		zui,
		&zui,
		zui8,
		&zui8,
		zui16,
		&zui16,
		zui32,
		&zui32,
		zui64,
		&zui64,
		zf32,
		&zf32,
		zf64,
		&zf64,
		zb,
		&zb,
		zs,
		&zs,
		zt,
		&zt,
		zni64,
		&zni64,
		znf64,
		&znf64,
		znb,
		&znb,
		zns,
		&zns,
		znt,
		&znt,
	}

	for i, valid := range valids {
		v := reflectx.IsZero(valid)
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
	}

	for i, invalid := range invalids {
		v := reflectx.IsZero(invalid)
		is.False(v, fmt.Sprintf("loop #%d", i))
	}
}
