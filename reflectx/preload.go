package reflectx

import (
	"reflect"

	"github.com/pkg/errors"
)

type Preloader interface {
	ForEach(name string, callback func(element reflect.Value) error) error
	OnExecute(kind reflect.Type, callback func(element interface{}) error) error
	OnUpdate(callback func(element interface{}) error) error

	StringIndexes() []string
	AddStringIndex(id string, element reflect.Value) error
	UpdateValueForStringIndex(name string, id string, element interface{}) error

	Int64Indexes() []int64
	AddInt64Index(id int64, element reflect.Value) error
	UpdateValueForInt64Index(name string, id int64, element interface{}) error
}

// type StringPreloader interface {
// 	Indexes() []string
// 	AddIndex(id string, element reflect.Value) error
// 	UpdateValueOnIndex(name string, id string, element interface{}) error
// }
//
// type Int64Preloader interface {
// 	Indexes() []int64
// 	AddIndex(id int64, element reflect.Value) error
// 	UpdateValueOnIndex(name string, id int64, element interface{}) error
// }

type preloader struct {
	value     interface{}
	relations reflect.Value
	mapString map[string][]reflect.Value
	mapInt64  map[int64][]reflect.Value
}

func NewPreloader(value interface{}) Preloader {
	return &preloader{
		value:     value,
		mapString: map[string][]reflect.Value{},
		mapInt64:  map[int64][]reflect.Value{},
	}
}

func (p *preloader) ForEach(name string, callback func(element reflect.Value) error) error {
	if IsSlice(p.value) {
		return p.forEach(name, callback)
	}
	return p.forOne(name, callback)
}

func (p *preloader) makeZeroValue(element interface{}, name string) error {
	value, err := getDestinationReflectValue(element, name)
	if err != nil {
		return err
	}

	if value.Kind() == reflect.Ptr && value.IsNil() && IsSlice(value) {
		value.Set(NewReflectSlice(GetSliceType(value)))
	}

	return nil
}

func (p *preloader) forOne(name string, callback func(element reflect.Value) error) error {
	err := p.makeZeroValue(p.value, name)
	if err != nil {
		return err
	}

	return callback(CreateReflectPointer(p.value))
}

func (p *preloader) forEach(name string, callback func(element reflect.Value) error) error {
	slice := GetIndirectValue(p.value)

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

		err := p.makeZeroValue(value.Interface(), name)
		if err != nil {
			return err
		}

		err = callback(value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *preloader) OnExecute(kind reflect.Type, callback func(element interface{}) error) error {
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

	p.relations = NewReflectSlice(elem)
	return callback(p.relations.Interface())
}

func (p *preloader) OnUpdate(callback func(element interface{}) error) error {
	if p.relations.Kind() == reflect.Ptr {
		p.relations = p.relations.Elem()
	}

	for i := 0; i < p.relations.Len(); i++ {
		val := p.relations.Index(i).Addr()

		err := callback(val.Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *preloader) StringIndexes() []string {
	list := make([]string, 0, len(p.mapString))
	for id := range p.mapString {
		list = append(list, id)
	}

	return list
}

func (p *preloader) AddStringIndex(id string, element reflect.Value) error {
	list, ok := p.mapString[id]
	if !ok {
		list = []reflect.Value{}
	}

	list = append(list, element)
	p.mapString[id] = list

	return nil
}

func (p *preloader) UpdateValueForStringIndex(name string, id string, element interface{}) error {
	values, ok := p.mapString[id]
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

func (p *preloader) Int64Indexes() []int64 {
	list := make([]int64, 0, len(p.mapInt64))
	for id := range p.mapInt64 {
		list = append(list, id)
	}

	return list
}

func (p *preloader) AddInt64Index(id int64, element reflect.Value) error {
	list, ok := p.mapInt64[id]
	if !ok {
		list = []reflect.Value{}
	}

	list = append(list, element)
	p.mapInt64[id] = list

	return nil
}

func (p *preloader) UpdateValueForInt64Index(name string, id int64, element interface{}) error {
	values, ok := p.mapInt64[id]
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
