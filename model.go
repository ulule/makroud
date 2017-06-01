package sqlxx

// PrimaryKeyType define a primary key type.
type PrimaryKeyType uint8

// PrimaryKey types.
const (
	// PrimaryKeyInteger define an integer primary key.
	PrimaryKeyInteger = PrimaryKeyType(iota)
	// PrimaryKeyUUID is unsupported at the moment.
	PrimaryKeyUUID
)

func (e PrimaryKeyType) String() string {
	switch e {
	case PrimaryKeyInteger:
		return "integer"
	}
	panic("sqlxx: unknown type")
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
