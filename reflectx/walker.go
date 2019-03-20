package reflectx

import (
	"reflect"
	"strings"
	"sync"

	"github.com/pkg/errors"
)

// A walker pool to reduce memory allocation pressure.
var walkerPool = &sync.Pool{
	New: func() interface{} {
		return &Walker{}
	},
}

// A Walker will execute a tree traversal of root value to extract values from given path.
// Once every values are extracted from given path, it will execute a callback with a pointer to a slice
// of theses values in Find.
//
// For example, given a slice of this type:
//
// type A struct {
//     B []struct {
//         C []struct {
//             Foobar int64
//         }
//     }
// }
//
// Using Find("B.C", callback) will execute given callback with every occurrences of C found.
//
type Walker struct {
	// value is the root value from tree.
	value interface{}
	// slice is a collection of values from given path.
	slice *reflect.Value
}

// NewWalker creates a new walker with given root value.
func NewWalker(value interface{}) *Walker {
	walker := walkerPool.Get().(*Walker)
	walker.value = value
	walker.slice = nil
	return walker
}

// Find will extract values from given path and execute given callback with a pointer to a slice.
// If the walker cannot extract values from given path
func (w *Walker) Find(path string, callback func(values interface{}) error) error {
	path = strings.Trim(path, ".")
	levels := strings.Split(path, ".")

	err := w.onEach(levels)
	if err != nil {
		return errors.Wrapf(err, "cannot find values from path: %s", path)
	}

	if w.slice == nil {
		return nil
	}

	return callback(w.slice.Interface())
}

// Close will cleanup current walker.
func (w *Walker) Close() {
	if w != nil {
		w.value = nil
		w.slice = nil
		walkerPool.Put(w)
	}
}

func (w *Walker) onEach(levels []string) error {
	if IsSlice(w.value) {

		slice := GetIndirectValue(w.value)

		for i := 0; i < slice.Len(); i++ {
			value, ok := w.normalizeReflectValue(slice.Index(i))
			if !ok {
				continue
			}

			err := w.find(value.Interface(), levels)
			if err != nil {
				return err
			}
		}

		return nil

	}

	return w.find(w.value, levels)
}

func (w *Walker) find(value interface{}, levels []string) error {

	leaf, err := GetFieldValueWithName(GetIndirectValue(value), levels[0])
	if err != nil {
		return err
	}

	// Ignore zero value or nil value.
	if leaf == nil || IsZero(leaf) {
		return nil
	}

	if IsSlice(leaf) {
		return w.findMany(leaf, levels)
	}

	return w.findOne(leaf, levels)
}

func (w *Walker) findMany(value interface{}, levels []string) error {
	slice := GetIndirectValue(value)
	for i := 0; i < slice.Len(); i++ {

		item, ok := w.normalizeReflectValue(slice.Index(i))
		if !ok {
			continue
		}

		if len(levels) >= 2 {
			err := w.find(item.Interface(), levels[1:])
			if err != nil {
				return err
			}
		} else {
			err := w.push(item)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *Walker) findOne(value interface{}, levels []string) error {
	item, ok := w.normalizeReflectValue(GetFlattenReflectValue(value))
	if !ok {
		return nil
	}

	if len(levels) >= 2 {
		return w.find(item.Interface(), levels[1:])
	}

	return w.push(item)
}

func (w *Walker) push(value reflect.Value) error {
	if w.slice == nil {
		slice := NewReflectSlice(GetReflectPointerType(value))
		w.slice = &slice
	}
	AppendReflectSlice((*w.slice), value.Interface())
	return nil
}

func (Walker) normalizeReflectValue(value reflect.Value) (reflect.Value, bool) {
	item := value
	if item.Kind() == reflect.Interface {
		item = reflect.ValueOf(item.Interface())
		if item.IsNil() {
			return value, false
		}
	}

	if item.Kind() != reflect.Ptr && item.CanAddr() {
		item = item.Addr()
	}

	return item, true
}
