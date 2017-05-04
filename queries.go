package sqlxx

import (
	"fmt"
	"strings"
)

// Queries is a list of Query instances
type Queries []Query

// Query is a relation query
type Query struct {
	// SQLX things
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
