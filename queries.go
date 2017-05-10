package sqlxx

import (
	"fmt"
	"sort"
	"strings"
)

// Queries is a list of Query instances
type Queries []Query

// ByTable returns queries for a given table.
func (q Queries) ByTable(name string) (*Query, bool) {
	for i := range q {
		if strings.Contains(q[i].Query, fmt.Sprintf("FROM %v", name)) {
			return &q[i], true
		}
	}

	return nil, false
}

// Query is a relation query
type Query struct {
	// sqlx things
	Query  string
	Args   []interface{}
	Params map[string]interface{}

	// Associations
	Field    Field
	FetchOne bool
}

// String returns struct instance string representation.
func (q Query) String() string {
	msg := []string{
		fmt.Sprintf("Query:\t%v", q.Query),
	}

	if len(q.Args) > 0 {
		msg = append(msg, fmt.Sprintf("Args:\t%v", q.Args))
	}

	if len(q.Params) > 0 {
		msg = append(msg, fmt.Sprintf("Params:\t%v", q.Params))
	}

	if q.Field.Name != "" {
		msg = append(msg, fmt.Sprintf("Field:\t%v", q.Field))
	}

	return fmt.Sprintf("\n%v\n", strings.Join(msg, "\n"))
}

// ----------------------------------------------------------------------------
// Columns
// ----------------------------------------------------------------------------

// Columns is a list of table columns.
type Columns []string

// Returns string representation of slice.
func (c Columns) String() string {
	sort.Strings(c)
	return strings.Join(c, ", ")
}

// ----------------------------------------------------------------------------
// Where clauses
// ----------------------------------------------------------------------------

// Conditions is a list of query conditions
type Conditions []string

// String returns conditions as AND query.
func (c Conditions) String() string {
	return strings.Join(c, " AND ")
}

// OR returns conditions as OR query.
func (c Conditions) OR() string {
	return strings.Join(c, " OR ")
}
