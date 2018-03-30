package sqlxx

import (
	"fmt"
)

// PrimaryKeyType define a primary key type.
type PrimaryKeyType uint8

// PrimaryKey types.
const (
	// PrimaryKeyInteger define an integer primary key.
	PrimaryKeyInteger = PrimaryKeyType(iota)
	// PrimaryKeyULID define an ulid primary key.
	PrimaryKeyULID
)

func (e PrimaryKeyType) String() string {
	switch e {
	case PrimaryKeyInteger:
		return "int64"
	case PrimaryKeyULID:
		return "ulid"
	default:
		panic(fmt.Sprintf("sqlxx: unknown type: %d", e))
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
	}
	panic("sqlxx: unknown type")
}

// Model represents a database table.
type Model interface {
	TableName() string
}

// Model represents a database table.
type XModel interface {
	// CreateSchema will define table schema using a builder.
	CreateSchema(builder SchemaBuilder)
	// WriteModel will update model with given values.
	WriteModel(mapper Mapper) error
}

// Models represents a list of Model.
type XModels interface {
	// Append will create a new Model in list using given mapper.
	Append(mapper Mapper) error
	// Model returns a empty Model instance.
	Model() XModel
}
