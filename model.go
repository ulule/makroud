package sqlxx

import (
	"fmt"
)

// PrimaryKeyType define a primary key type.
type PrimaryKeyType uint8

// PrimaryKey types.
const (
	// PrimaryKeyIntegerType uses an integer as primary key.
	PrimaryKeyIntegerType = PrimaryKeyType(iota)
	// PrimaryKeyString uses a string as primary key.
	PrimaryKeyStringType
)

func (e PrimaryKeyType) String() string {
	switch e {
	case PrimaryKeyIntegerType:
		return "int64"
	case PrimaryKeyStringType:
		return "string"
	default:
		panic(fmt.Sprintf("sqlxx: unknown primary key type: %d", e))
	}
}

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
