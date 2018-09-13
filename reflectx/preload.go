package reflectx

import (
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

// A PreloadValue is an alias type to reflect.Value used to manipulate a reference while executing a preload.
type PreloadValue reflect.Value

// Unwrap returns underlying value.
func (v PreloadValue) Unwrap() interface{} {
	return reflect.Value(v).Interface()
}

// A StringPreloader is used by preload mechanism to attach values received from query to the correct reference.
// This preloader handles everything related to reflection, that means the query engine is handled on
// another components. It uses a string mapper internally.
type StringPreloader struct {
	name      string
	kind      reflect.Type
	value     interface{}
	relations reflect.Value
	mapper    map[string][]reflect.Value
}

// NewStringPreloader creates a new StringPreloader using given value as root reference and given name as the element
// to preload on root reference.
func NewStringPreloader(name string, kind reflect.Type, value interface{}) *StringPreloader {
	preloader := stringPreloaderPool.Get().(*StringPreloader)
	preloader.name = name
	preloader.kind = kind
	preloader.value = value
	return preloader
}

// ForEach will scan every root references (that means a simple struct or a slice of struct) and creates a zero
// value of the element to preload (identified by name).
// The given callback should extract the primary, or the foreign key, and call AddIndex() so later we can find what
// element should be linked to the parent references using UpdateValueOnIndex().
func (p *StringPreloader) ForEach(callback func(element PreloadValue) error) error {
	return preloadForEach(p.value, p.name, callback)
}

// OnExecute creates a slice of the element to preload and execute the given callback.
// The callback must populate the slice with required elements so we can use the results with OnUpdate() later.
func (p *StringPreloader) OnExecute(callback func(element interface{}) error) error {
	elem, err := preloadOnExecuteGetType(p.kind)
	if err != nil {
		return err
	}

	p.relations = NewReflectSlice(elem)
	return callback(p.relations.Interface())
}

// OnUpdate executes given callback on each elements found by OnExecute().
// The callback must use UpdateValueOnIndex() with the correct primary key or foreign key so we can link
// the element to the correct parent.
func (p *StringPreloader) OnUpdate(callback func(element interface{}) error) error {
	return preloadOnUpdate(p.relations, callback)
}

// Indexes returns a list of keys.
func (p *StringPreloader) Indexes() []string {
	list := make([]string, 0, len(p.mapper))
	for id := range p.mapper {
		list = append(list, id)
	}

	return list
}

// AddIndex appends given element to be updated if a value is found with this id.
func (p *StringPreloader) AddIndex(id string, element PreloadValue) error {
	list, ok := p.mapper[id]
	if !ok {
		list = []reflect.Value{}
	}

	list = append(list, reflect.Value(element))
	p.mapper[id] = list

	return nil
}

// UpdateValueOnIndex adds given element on its parent.
func (p *StringPreloader) UpdateValueOnIndex(id string, element interface{}) error {
	values, ok := p.mapper[id]
	if !ok {
		return errors.Errorf("cannot find element with primary key: '%s'", id)
	}

	for i := range values {
		err := PushFieldValue(values[i].Interface(), p.name, element, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close will cleanup current preloader.
func (p *StringPreloader) Close() {
	if p != nil {
		p.name = ""
		p.kind = nil
		p.value = nil
		p.relations = reflect.Value{}
		p.mapper = map[string][]reflect.Value{}
		stringPreloaderPool.Put(p)
	}
}

// A IntegerPreloader is used by preload mechanism to attach values received from query to the correct reference.
// This preloader handles everything related to reflection, that means the query engine is handled on
// another components. It uses a integer mapper internally.
type IntegerPreloader struct {
	name      string
	kind      reflect.Type
	value     interface{}
	relations reflect.Value
	mapper    map[int64][]reflect.Value
}

// NewIntegerPreloader creates a new IntegerPreloader using given value as root reference and given name as the element
// to preload on root reference.
func NewIntegerPreloader(name string, kind reflect.Type, value interface{}) *IntegerPreloader {
	preloader := integerPreloaderPool.Get().(*IntegerPreloader)
	preloader.name = name
	preloader.kind = kind
	preloader.value = value
	return preloader
}

// ForEach will scan every root references (that means a simple struct or a slice of struct) and creates a zero
// value of the element to preload (identified by name).
// The given callback should extract the primary, or the foreign key, and call AddIndex() so later we can find what
// element should be linked to the parent references using UpdateValueOnIndex().
func (p *IntegerPreloader) ForEach(callback func(element PreloadValue) error) error {
	return preloadForEach(p.value, p.name, callback)
}

// OnExecute creates a slice of the element to preload and execute the given callback.
// The callback must populate the slice with required elements so we can use the results with OnUpdate() later.
func (p *IntegerPreloader) OnExecute(callback func(element interface{}) error) error {
	elem, err := preloadOnExecuteGetType(p.kind)
	if err != nil {
		return err
	}

	p.relations = NewReflectSlice(elem)
	return callback(p.relations.Interface())
}

// OnUpdate executes given callback on each elements found by OnExecute().
// The callback must use UpdateValueOnIndex() with the correct primary key or foreign key so we can link
// the element to the correct parent.
func (p *IntegerPreloader) OnUpdate(callback func(element interface{}) error) error {
	return preloadOnUpdate(p.relations, callback)
}

// Indexes returns a list of keys.
func (p *IntegerPreloader) Indexes() []int64 {
	list := make([]int64, 0, len(p.mapper))
	for id := range p.mapper {
		list = append(list, id)
	}

	return list
}

// AddIndex appends given element to be updated if a value is found with this id.
func (p *IntegerPreloader) AddIndex(id int64, element PreloadValue) error {
	list, ok := p.mapper[id]
	if !ok {
		list = []reflect.Value{}
	}

	list = append(list, reflect.Value(element))
	p.mapper[id] = list

	return nil
}

// UpdateValueOnIndex adds given element on its parent.
func (p *IntegerPreloader) UpdateValueOnIndex(id int64, element interface{}) error {
	values, ok := p.mapper[id]
	if !ok {
		return errors.Errorf("cannot find element with primary key: '%d'", id)
	}

	for i := range values {
		err := PushFieldValue(values[i].Interface(), p.name, element, false)
		if err != nil {
			return err
		}
	}

	return nil
}

// Close will cleanup current preloader.
func (p *IntegerPreloader) Close() {
	if p != nil {
		p.name = ""
		p.kind = nil
		p.value = nil
		p.relations = reflect.Value{}
		p.mapper = map[int64][]reflect.Value{}
		integerPreloaderPool.Put(p)
	}
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

// A string preloader pool to reduce memory allocation pressure.
var stringPreloaderPool = &sync.Pool{
	New: func() interface{} {
		return &StringPreloader{
			mapper: map[string][]reflect.Value{},
		}
	},
}

// A integer preloader pool to reduce memory allocation pressure.
var integerPreloaderPool = &sync.Pool{
	New: func() interface{} {
		return &IntegerPreloader{
			mapper: map[int64][]reflect.Value{},
		}
	},
}
