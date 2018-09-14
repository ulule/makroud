package makroud_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/makroud"
)

func TestPing(t *testing.T) {
	Setup(t)(func(driver makroud.Driver) {
		is := require.New(t)

		err := driver.Ping()
		is.NoError(err)
	})
}
