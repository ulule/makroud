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

// TagsMapping is the reflekt.Tags mapping to handle struct tag without key:value format
var TagsMapping = map[string]string{
	"db": "field",
}

// Struct tag key names.
const (
	StructTagPrimaryKey = "primary_key"
	StructTagIgnored    = "ignored"
	StructTagDefault    = "default"
)

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

// RelationsOne are One relations.
var RelationsOne = map[RelationType]bool{
	RelationTypeOneToOne:  true,
	RelationTypeOneToMany: true,
}

// RelationsMany are Many relations.
var RelationsMany = map[RelationType]bool{
	RelationTypeManyToOne:  true,
	RelationTypeManyToMany: true,
}
