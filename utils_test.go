package sqlxx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnakeCase(t *testing.T) {
	var results = []struct {
		in  string
		out string
	}{
		{"FooBar", "foo_bar"},
		{"ID", "id"},
		{"UserID", "user_id"},
		{"User1234", "user1234"},
		{"blahBlah", "blah_blah"},
	}

	for _, tt := range results {
		s := toSnakeCase(tt.in)
		assert.Equal(t, s, tt.out)
	}
}
