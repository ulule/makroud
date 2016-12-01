package sqlxx

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRelationQueries(t *testing.T) {
	is := assert.New(t)

	// Test order by level

	rq := RelationQueries{
		RelationQuery{level: 1},
		RelationQuery{level: 1},
		RelationQuery{level: 1},
		RelationQuery{level: 30},
		RelationQuery{level: 20},
		RelationQuery{level: 4},
		RelationQuery{level: 4},
	}

	sort.Sort(rq)

	is.Equal(rq[0].level, 1)
	is.Equal(rq[1].level, 1)
	is.Equal(rq[2].level, 1)
	is.Equal(rq[3].level, 4)
	is.Equal(rq[4].level, 4)
	is.Equal(rq[5].level, 20)
	is.Equal(rq[6].level, 30)
}
