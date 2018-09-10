package reflectx

import (
	"reflect"

	"github.com/pkg/errors"
)

type PreloadValue reflect.Value

func (v PreloadValue) Unwrap() interface{} {
	return reflect.Value(v).Interface()
}

type Preloader interface {
	ForEach(name string, callback func(element PreloadValue) error) error
	OnExecute(kind reflect.Type, callback func(element interface{}) error) error
	OnUpdate(callback func(element interface{}) error) error

	StringIndexes() []string
	AddStringIndex(id string, element PreloadValue) error
	UpdateValueForStringIndex(name string, id string, element interface{}) error

	Int64Indexes() []int64
	AddInt64Index(id int64, element PreloadValue) error
	UpdateValueForInt64Index(name string, id int64, element interface{}) error
}

type StringPreloader interface {
	ForEach(name string, callback func(element PreloadValue) error) error
	OnExecute(kind reflect.Type, callback func(element interface{}) error) error
	OnUpdate(callback func(element interface{}) error) error

	Indexes() []string
	AddIndex(id string, element PreloadValue) error
	UpdateValueOnIndex(name string, id string, element interface{}) error
}

type IntegerPreloader interface {
	ForEach(name string, callback func(element PreloadValue) error) error
	OnExecute(kind reflect.Type, callback func(element interface{}) error) error
	OnUpdate(callback func(element interface{}) error) error

	Indexes() []int64
	AddIndex(id int64, element PreloadValue) error
	UpdateValueOnIndex(name string, id int64, element interface{}) error
}

type stringPreloader struct {
	value     interface{}
	relations reflect.Value
	mapper    map[string][]reflect.Value
}

func NewStringPreloader(value interface{}) StringPreloader {
	return &stringPreloader{
		value:  value,
		mapper: map[string][]reflect.Value{},
	}
}

func (p *stringPreloader) ForEach(name string, callback func(element PreloadValue) error) error {
	return preloadForEach(p.value, name, callback)
}

func (p *stringPreloader) OnExecute(kind reflect.Type, callback func(element interface{}) error) error {
	elem, err := preloadOnExecuteGetType(kind)
	if err != nil {
		return err
	}

	p.relations = NewReflectSlice(elem)
	return callback(p.relations.Interface())
}

func (p *stringPreloader) OnUpdate(callback func(element interface{}) error) error {
	return preloadOnUpdate(p.relations, callback)
}

func (p *stringPreloader) Indexes() []string {
	list := make([]string, 0, len(p.mapper))
	for id := range p.mapper {
		list = append(list, id)
	}

	return list
}

func (p *stringPreloader) AddIndex(id string, element PreloadValue) error {
	list, ok := p.mapper[id]
	if !ok {
		list = []reflect.Value{}
	}

	list = append(list, reflect.Value(element))
	p.mapper[id] = list

	return nil
}

func (p *stringPreloader) UpdateValueOnIndex(name string, id string, element interface{}) error {
	values, ok := p.mapper[id]
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

type integerPreloader struct {
	value     interface{}
	relations reflect.Value
	mapper    map[int64][]reflect.Value
}

func NewIntegerPreloader(value interface{}) IntegerPreloader {
	return &integerPreloader{
		value:  value,
		mapper: map[int64][]reflect.Value{},
	}
}

func (p *integerPreloader) ForEach(name string, callback func(element PreloadValue) error) error {
	return preloadForEach(p.value, name, callback)
}

func (p *integerPreloader) OnExecute(kind reflect.Type, callback func(element interface{}) error) error {
	elem, err := preloadOnExecuteGetType(kind)
	if err != nil {
		return err
	}

	p.relations = NewReflectSlice(elem)
	return callback(p.relations.Interface())
}

func (p *integerPreloader) OnUpdate(callback func(element interface{}) error) error {
	return preloadOnUpdate(p.relations, callback)
}

func (p *integerPreloader) Indexes() []int64 {
	list := make([]int64, 0, len(p.mapper))
	for id := range p.mapper {
		list = append(list, id)
	}

	return list
}

func (p *integerPreloader) AddIndex(id int64, element PreloadValue) error {
	list, ok := p.mapper[id]
	if !ok {
		list = []reflect.Value{}
	}

	list = append(list, reflect.Value(element))
	p.mapper[id] = list

	return nil
}

func (p *integerPreloader) UpdateValueOnIndex(name string, id int64, element interface{}) error {
	values, ok := p.mapper[id]
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

func preloadForEach(value interface{}, name string, callback func(element PreloadValue) error) error {
	if IsSlice(value) {
		return preloadExecForEach(value, name, callback)
	}
	return preloadExecForOne(value, name, callback)
}

func preloadExecForOne(value interface{}, name string, callback func(element PreloadValue) error) error {
	err := preloadMakeZeroValue(value, name)
	if err != nil {
		return err
	}

	return callback(PreloadValue(CreateReflectPointer(value)))
}

func preloadExecForEach(value interface{}, name string, callback func(element PreloadValue) error) error {
	slice := GetIndirectValue(value)

	for i := 0; i < slice.Len(); i++ {
		elem := slice.Index(i)

		if elem.Kind() == reflect.Interface {
			elem = reflect.ValueOf(elem.Interface())
			if elem.IsNil() {
				continue
			}
		}

		if elem.Kind() != reflect.Ptr && elem.CanAddr() {
			elem = elem.Addr()
		}

		err := preloadMakeZeroValue(elem.Interface(), name)
		if err != nil {
			return err
		}

		err = callback(PreloadValue(elem))
		if err != nil {
			return err
		}
	}

	return nil
}

func preloadMakeZeroValue(element interface{}, name string) error {
	value, err := getDestinationReflectValue(element, name)
	if err != nil {
		return err
	}

	if value.Kind() == reflect.Ptr && value.IsNil() && IsSlice(value) {
		value.Set(NewReflectSlice(GetSliceType(value)))
	}

	return nil
}

func preloadOnExecuteGetType(kind reflect.Type) (reflect.Type, error) {
	elem := kind

	if elem.Kind() == reflect.Slice {

		// If output type is a slice, create a new slice with it's indirect type.
		// For example, a slice with "[]*Foobar" as type will create a new slice with "[]Foobar" as type.

		elem = elem.Elem()
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct {
			return kind, errors.Errorf("cannot execute a preload this type: %s", elem)
		}

	} else if elem.Kind() == reflect.Struct || elem.Kind() == reflect.Ptr {

		// If output type is not a slice, so either a struct or a pointer to a struct,
		// create a new slice with it's indirect type.
		// For example, a pointer with "*Foobar" as type will create a new slice with "[]Foobar" as type.

		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct {
			return kind, errors.Errorf("cannot execute a preload this type: %s", elem)
		}

	} else {
		return kind, errors.Errorf("cannot execute a preload this type: %s", kind)
	}

	return elem, nil
}

func preloadOnUpdate(relations reflect.Value, callback func(element interface{}) error) error {
	if relations.Kind() == reflect.Ptr {
		relations = relations.Elem()
	}

	for i := 0; i < relations.Len(); i++ {
		val := relations.Index(i).Addr()

		err := callback(val.Interface())
		if err != nil {
			return err
		}
	}

	return nil
}
