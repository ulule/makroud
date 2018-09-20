package reflectx_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/makroud/reflectx"
)

func TestReflectx_Preload(t *testing.T) {
	is := require.New(t)

	type A struct {
		Value int64
		Link  int64
	}

	type B struct {
		Value string
		Link  string
	}

	type WrapA struct {
		Hash  int64
		LinkA *A
	}

	type WrapB struct {
		Hash  string
		LinkB *B
	}

	{

		wb := &WrapB{
			Hash: "01CQ6KTTVHG17C97QMQXV3MW4G",
		}
		b := &B{
			Value: "foobar",
			Link:  "01CQ6KTTVHG17C97QMQXV3MW4G",
		}

		name := "LinkB"
		kind, err := reflectx.GetFieldReflectTypeByName(WrapB{}, name)
		is.NoError(err)
		is.NotEmpty(kind)

		preloader := reflectx.NewStringPreloader(name, kind, wb)
		defer preloader.Close()

		count := 0
		err = preloader.ForEach(func(element reflectx.PreloadValue) error {
			pk, err := reflectx.GetFieldValueString(element.Unwrap(), "Hash")
			if err != nil {
				return err
			}
			count++
			return preloader.AddIndex(pk, element)
		})
		is.NoError(err)
		is.Equal(1, count)

		indexes := preloader.Indexes()
		is.NotEmpty(indexes)
		is.Len(indexes, 1)
		is.Equal(wb.Hash, indexes[0])

		count = 0
		err = preloader.OnExecute(func(element interface{}) error {
			slice := reflect.ValueOf(element)
			reflectx.AppendReflectSlice(slice, b)
			count++
			return nil
		})
		is.NoError(err)
		is.Equal(1, count)

		count = 0
		err = preloader.OnUpdate(func(element interface{}) error {
			fk, err := reflectx.GetFieldValueString(element, "Link")
			if err != nil {
				return err
			}
			count++
			return preloader.UpdateValueOnIndex(fk, element)
		})
		is.NoError(err)
		is.Equal(1, count)
		is.NotNil(wb.LinkB)
		is.Equal(b, wb.LinkB)

	}
	{

		wa := &WrapA{
			Hash: 500,
		}
		a := &A{
			Value: 42,
			Link:  500,
		}

		name := "LinkA"
		kind, err := reflectx.GetFieldReflectTypeByName(WrapA{}, name)
		is.NoError(err)
		is.NotEmpty(kind)

		preloader := reflectx.NewIntegerPreloader(name, kind, wa)
		defer preloader.Close()

		count := 0
		err = preloader.ForEach(func(element reflectx.PreloadValue) error {
			pk, err := reflectx.GetFieldValueInt64(element.Unwrap(), "Hash")
			if err != nil {
				return err
			}
			count++
			return preloader.AddIndex(pk, element)
		})
		is.NoError(err)
		is.Equal(1, count)

		indexes := preloader.Indexes()
		is.NotEmpty(indexes)
		is.Len(indexes, 1)
		is.Equal(wa.Hash, indexes[0])

		count = 0
		err = preloader.OnExecute(func(element interface{}) error {
			slice := reflect.ValueOf(element)
			reflectx.AppendReflectSlice(slice, a)
			count++
			return nil
		})
		is.NoError(err)
		is.Equal(1, count)

		count = 0
		err = preloader.OnUpdate(func(element interface{}) error {
			fk, err := reflectx.GetFieldValueInt64(element, "Link")
			if err != nil {
				return err
			}
			count++
			return preloader.UpdateValueOnIndex(fk, element)
		})
		is.NoError(err)
		is.Equal(1, count)
		is.NotNil(wa.LinkA)
		is.Equal(a, wa.LinkA)

	}
	{
		lwa := &[]WrapA{
			{
				Hash: 500,
			},
			{
				Hash: 501,
			},
			{
				Hash: 502,
			},
			{
				Hash: 503,
			},
			{
				Hash: 505,
			},
		}

		a1 := &A{
			Value: 42,
			Link:  500,
		}
		a2 := &A{
			Value: 44,
			Link:  501,
		}
		a3 := &A{
			Value: 46,
			Link:  502,
		}
		a4 := &A{
			Value: 48,
			Link:  503,
		}

		name := "LinkA"
		kind, err := reflectx.GetFieldReflectTypeByName(WrapA{}, name)
		is.NoError(err)
		is.NotEmpty(kind)

		preloader := reflectx.NewIntegerPreloader(name, kind, lwa)
		defer preloader.Close()

		count := 0
		err = preloader.ForEach(func(element reflectx.PreloadValue) error {
			pk, err := reflectx.GetFieldValueInt64(element.Unwrap(), "Hash")
			if err != nil {
				return err
			}
			count++
			return preloader.AddIndex(pk, element)
		})
		is.NoError(err)
		is.Equal(5, count)

		indexes := preloader.Indexes()
		is.NotEmpty(indexes)
		is.Len(indexes, 5)
		sort.Slice(indexes, func(i int, j int) bool {
			return indexes[i] < indexes[j]
		})
		is.Equal((*lwa)[0].Hash, indexes[0])
		is.Equal((*lwa)[1].Hash, indexes[1])
		is.Equal((*lwa)[2].Hash, indexes[2])
		is.Equal((*lwa)[3].Hash, indexes[3])
		is.Equal((*lwa)[4].Hash, indexes[4])

		count = 0
		err = preloader.OnExecute(func(element interface{}) error {
			slice := reflect.ValueOf(element)
			reflectx.AppendReflectSlice(slice, a1)
			reflectx.AppendReflectSlice(slice, a2)
			reflectx.AppendReflectSlice(slice, a3)
			reflectx.AppendReflectSlice(slice, a4)
			count++
			return nil
		})
		is.NoError(err)
		is.Equal(1, count)

		count = 0
		err = preloader.OnUpdate(func(element interface{}) error {
			fk, err := reflectx.GetFieldValueInt64(element, "Link")
			if err != nil {
				return err
			}
			count++
			return preloader.UpdateValueOnIndex(fk, element)
		})
		is.NoError(err)
		is.Equal(4, count)
		is.NotNil((*lwa)[0].LinkA)
		is.Equal(a1, (*lwa)[0].LinkA)
		is.NotNil((*lwa)[1].LinkA)
		is.Equal(a2, (*lwa)[1].LinkA)
		is.NotNil((*lwa)[2].LinkA)
		is.Equal(a3, (*lwa)[2].LinkA)
		is.NotNil((*lwa)[3].LinkA)
		is.Equal(a4, (*lwa)[3].LinkA)
		is.Nil((*lwa)[4].LinkA)

	}
}
