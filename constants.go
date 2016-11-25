package sqlxx

// Struct tags
const (
	StructTagName     = "sqlxx"
	SQLXStructTagName = "db"
)

// SupportedTags are supported tags.
var SupportedTags = []string{
	StructTagName,
	SQLXStructTagName,
}

// RelationType is a field relation type.
type RelationType int

// Field types.
const (
	RelationTypeUnknown RelationType = iota
	RelationTypeOneToOne
	RelationTypeOneToMany
	RelationTypeManyToOne
	RelationTypeManyToMany
)

// RelationTypes are supported relations types.
var RelationTypes = map[RelationType]bool{
	RelationTypeOneToOne:   true,
	RelationTypeOneToMany:  true,
	RelationTypeManyToOne:  true,
	RelationTypeManyToMany: true,
}