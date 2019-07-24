package reflectx_test

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/makroud/reflectx"
)

func TestReflectx_IsStruct(t *testing.T) {
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
	run := 'x'
	str := "x"
	bl := true
	sl := []int{}
	ar := [5]int{}
	mp := map[string]bool{}
	ch := make(chan struct{})
	f := func() bool {
		return false
	}
	p := &time.Time{}

	invalids := []interface{}{
		i,
		i8,
		i16,
		i32,
		i64,
		ui,
		ui8,
		ui16,
		ui32,
		ui64,
		f32,
		f64,
		run,
		str,
		bl,
		sl,
		ar,
		mp,
		ch,
		f,
		p,
		reflect.ValueOf(i),
		reflect.ValueOf(i8),
		reflect.ValueOf(i16),
		reflect.ValueOf(i32),
		reflect.ValueOf(i64),
		reflect.ValueOf(ui),
		reflect.ValueOf(ui8),
		reflect.ValueOf(ui16),
		reflect.ValueOf(ui32),
		reflect.ValueOf(ui64),
		reflect.ValueOf(f32),
		reflect.ValueOf(f64),
		reflect.ValueOf(run),
		reflect.ValueOf(str),
		reflect.ValueOf(bl),
		reflect.ValueOf(sl),
		reflect.ValueOf(ar),
		reflect.ValueOf(mp),
		reflect.ValueOf(ch),
		reflect.ValueOf(f),
		reflect.ValueOf(p),
	}

	for i, invalid := range invalids {
		v := reflectx.IsStruct(invalid)
		is.False(v, fmt.Sprintf("loop #%d", i))
	}

	v1 := time.Time{}
	v2 := sql.NullInt64{}
	v3 := Elements{}
	v4 := pq.NullTime{}
	v5 := struct{}{}

	valids := []interface{}{
		v1,
		v2,
		v3,
		v4,
		v5,
		reflect.ValueOf(v1),
		reflect.ValueOf(v2),
		reflect.ValueOf(v3),
		reflect.ValueOf(v4),
		reflect.ValueOf(v5),
	}

	for i, valid := range valids {
		v := reflectx.IsStruct(valid)
		is.True(v, fmt.Sprintf("loop #%d", i))
	}
}
