package sqlxx

import (
	"reflect"
	"time"
)

func now() time.Time {
	return time.Now()
}

func isZeroValue(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}
