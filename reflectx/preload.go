package reflectx

import (
	"reflect"

	"github.com/pkg/errors"
)

type Preloader interface {
	ForEach(callback func(element reflect.Value) error) error
	OnExecute(kind reflect.Type, callback func(element interface{}) error) error
	OnUpdate(callback func(element interface{}) error) error

	StringIndexes() []string
	AddStringIndex(id string, element reflect.Value) error
	UpdateValueForStringIndex(name string, id string, element interface{}) error

	Int64Indexes() []int64
	AddInt64Index(id int64, element reflect.Value) error
	UpdateValueForInt64Index(name string, id int64, element interface{}) error
}

type StructPreloader struct {
	value     interface{}
	relations reflect.Value
	mapString map[string][]reflect.Value
	mapInt64  map[int64][]reflect.Value
}

func NewStructPreloader(value interface{}) *StructPreloader {
	return &StructPreloader{
		value:     value,
		mapString: map[string][]reflect.Value{},
		mapInt64:  map[int64][]reflect.Value{},
	}
}

func (w *StructPreloader) ForEach(callback func(element reflect.Value) error) error {
	return callback(CreateReflectPointer(w.value))
}

func (w *StructPreloader) OnExecute(kind reflect.Type, callback func(element interface{}) error) error {
	elem := kind

	if elem.Kind() == reflect.Slice {

		// If output type is a slice, create a new slice with it's indirect type.
		// For example, a slice with "[]*Foobar" as type will create a new slice with "[]Foobar" as type.

		elem = elem.Elem()
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct {
			return errors.Errorf("cannot execute a preload this type: %s", elem)
		}

	} else if elem.Kind() == reflect.Struct || elem.Kind() == reflect.Ptr {

		// If output type is not a slice, so either a struct or a pointer to a struct,
		// create a new slice with it's indirect type.
		// For example, a pointer with "*Foobar" as type will create a new slice with "[]Foobar" as type.

		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct {
			return errors.Errorf("cannot execute a preload this type: %s", elem)
		}

	} else {
		return errors.Errorf("cannot execute a preload this type: %s", kind)
	}

	w.relations = NewReflectSlice(elem)
	return callback(w.relations.Interface())
}

func (w *StructPreloader) OnUpdate(callback func(element interface{}) error) error {
	if w.relations.Kind() == reflect.Ptr {
		w.relations = w.relations.Elem()
	}

	if w.relations.Len() == 0 {
		return nil
	}

	val := w.relations.Index(0).Addr()
	return callback(val.Interface())
}

func (w *StructPreloader) StringIndexes() []string {
	return getStringIndexes(w.mapString)
}

func (w *StructPreloader) AddStringIndex(id string, element reflect.Value) error {
	return addStringIndex(w.mapString, id, element)
}

func (w *StructPreloader) UpdateValueForStringIndex(name string, id string, element interface{}) error {
	values, ok := w.mapString[id]
	if !ok {
		return errors.Errorf("cannot find element with primary key: '%s'", id)
	}
	if len(values) != 1 {
		return errors.Errorf("only one element was expected for primary key: '%s'", id)
	}

	return PushFieldValue(values[0].Interface(), name, element, false)
}

func (w *StructPreloader) Int64Indexes() []int64 {
	return getInt64Indexes(w.mapInt64)
}

func (w *StructPreloader) AddInt64Index(id int64, element reflect.Value) error {
	return addInt64Index(w.mapInt64, id, element)
}

func (w *StructPreloader) UpdateValueForInt64Index(name string, id int64, element interface{}) error {
	values, ok := w.mapInt64[id]
	if !ok {
		return errors.Errorf("cannot find element with primary key: '%d'", id)
	}
	if len(values) != 1 {
		return errors.Errorf("only one element was expected for primary key: '%d'", id)
	}

	return PushFieldValue(values[0].Interface(), name, element, false)
}

type SlicePreloader struct {
	value     interface{}
	relations reflect.Value
	mapString map[string][]reflect.Value
	mapInt64  map[int64][]reflect.Value
}

func NewSlicePreloader(value interface{}) *SlicePreloader {
	return &SlicePreloader{
		value:     value,
		mapString: map[string][]reflect.Value{},
		mapInt64:  map[int64][]reflect.Value{},
	}
}

func (w *SlicePreloader) ForEach(callback func(element reflect.Value) error) error {
	slice := GetIndirectValue(w.value)

	for i := 0; i < slice.Len(); i++ {
		value := slice.Index(i)

		if value.Kind() == reflect.Interface {
			value = reflect.ValueOf(value.Interface())
			if value.IsNil() {
				continue
			}
		}

		if value.Kind() != reflect.Ptr && value.CanAddr() {
			value = value.Addr()
		}

		err := callback(value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *SlicePreloader) OnExecute(kind reflect.Type, callback func(element interface{}) error) error {
	elem := kind

	if elem.Kind() == reflect.Slice {

		// If output type is a slice, create a new slice with it's indirect type.
		// For example, a slice with "[]*Foobar" as type will create a new slice with "[]Foobar" as type.

		elem = elem.Elem()
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct {
			return errors.Errorf("cannot execute a preload this type: %s", elem)
		}

	} else if elem.Kind() == reflect.Struct || elem.Kind() == reflect.Ptr {

		// If output type is not a slice, so either a struct or a pointer to a struct,
		// create a new slice with it's indirect type.
		// For example, a pointer with "*Foobar" as type will create a new slice with "[]Foobar" as type.

		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct {
			return errors.Errorf("cannot execute a preload this type: %s", elem)
		}

	} else {
		return errors.Errorf("cannot execute a preload this type: %s", kind)
	}

	w.relations = NewReflectSlice(elem)

	return callback(w.relations.Interface())
}

func (w *SlicePreloader) OnUpdate(callback func(element interface{}) error) error {
	if w.relations.Kind() == reflect.Ptr {
		w.relations = w.relations.Elem()
	}

	for i := 0; i < w.relations.Len(); i++ {
		val := w.relations.Index(i).Addr()

		err := callback(val.Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *SlicePreloader) StringIndexes() []string {
	return getStringIndexes(w.mapString)
}

func (w *SlicePreloader) AddStringIndex(id string, element reflect.Value) error {
	return addStringIndex(w.mapString, id, element)
}

func (w *SlicePreloader) UpdateValueForStringIndex(name string, id string, element interface{}) error {
	values, ok := w.mapString[id]
	if !ok {
		return errors.Errorf("cannot find element with primary key: '%s'", id)
	}

	for i := range values {
		err := PushFieldValue(values[i].Interface(), name, element, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *SlicePreloader) Int64Indexes() []int64 {
	return getInt64Indexes(w.mapInt64)
}

func (w *SlicePreloader) AddInt64Index(id int64, element reflect.Value) error {
	return addInt64Index(w.mapInt64, id, element)
}

func (w *SlicePreloader) UpdateValueForInt64Index(name string, id int64, element interface{}) error {
	values, ok := w.mapInt64[id]
	if !ok {
		return errors.Errorf("cannot find element with primary key: '%d'", id)
	}

	for i := range values {
		err := PushFieldValue(values[i].Interface(), name, element, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func getInt64Indexes(values map[int64][]reflect.Value) []int64 {
	list := make([]int64, 0, len(values))
	for id := range values {
		list = append(list, id)
	}
	return list
}

func getStringIndexes(values map[string][]reflect.Value) []string {
	list := make([]string, 0, len(values))
	for id := range values {
		list = append(list, id)
	}
	return list
}

func addInt64Index(values map[int64][]reflect.Value, id int64, element reflect.Value) error {
	list, ok := values[id]
	if !ok {
		list = []reflect.Value{}
	}

	list = append(list, element)
	values[id] = list

	return nil
}

func addStringIndex(values map[string][]reflect.Value, id string, element reflect.Value) error {
	list, ok := values[id]
	if !ok {
		list = []reflect.Value{}
	}

	list = append(list, element)
	values[id] = list

	return nil
}
