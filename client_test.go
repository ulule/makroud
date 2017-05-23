package sqlxx_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	env := setup(t)
	defer env.teardown()

	is := require.New(t)

	err := env.driver.Ping()
	is.NoError(err)
}
