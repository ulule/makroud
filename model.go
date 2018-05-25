package sqlxx

import (
	"fmt"
)

// AssociationType define an association type.
type AssociationType uint8

// Association types.
const (
	AssociationTypeUndefined = AssociationType(iota)
	AssociationTypeOne
	AssociationTypeMany
	AssociationTypeManyToMany
)

func (e AssociationType) String() string {
	switch e {
	case AssociationTypeUndefined:
		return "undefined"
	case AssociationTypeOne:
		return "one"
	case AssociationTypeMany:
		return "many"
	case AssociationTypeManyToMany:
		return "many-to-many"
	default:
		panic(fmt.Sprintf("sqlxx: unknown association type: %d", e))
	}
}

// Model represents a database table.
type Model interface {
	TableName() string
}

// ModelOpts define a configuration for model.
type ModelOpts struct {
	PrimaryKey string
	CreatedKey string
	UpdatedKey string
	DeletedKey string
}
