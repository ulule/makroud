package sqlxx

// Queries is a list of Query instances
type Queries []Query

// Query is a relation query
type Query struct {
	Field    Field
	Query    string
	Args     []interface{}
	Params   map[string]interface{}
	FetchOne bool
}

// String returns struct instance string representation.
func (q Query) String() string {
	return q.Query
}
