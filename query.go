package sqlxx

import "strings"

// ----------------------------------------------------------------------------
// Columns
// ----------------------------------------------------------------------------

// Columns is a list of table columns.
type Columns []string

// Returns string representation of slice.
func (c Columns) String() string {
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
