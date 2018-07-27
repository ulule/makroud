package sqlxx

import (
	"github.com/pkg/errors"
)

// TODO:
//
//     |--------------|--------------|-----------------|
//     |    Source    |    Action    |    Reference    |
//     |--------------|--------------|-----------------|
//     |      1       |      ->      |        1        |
//     |      1       |      <-      |        1        |
//     |    Source    |    Action    |    Reference    |
//     |    Source    |    Action    |    Reference    |
//     |    Source    |    Action    |    Reference    |

// Preload preloads related fields.
func Preload(driver Driver, out interface{}, paths ...string) error {
	_, err := PreloadWithQueries(driver, out, paths...)
	return err
}

// PreloadWithQueries preloads related fields and returns performed queries.
func PreloadWithQueries(driver Driver, out interface{}, paths ...string) (Queries, error) {
	return nil, errors.New("sqlxx: cannot execute preload")
}
