package sqlxx

// AssociationType is an association type.
type AssociationType uint8

func (a AssociationType) String() string {
	return map[AssociationType]string{
		AssociationTypeUndefined:  "undefined",
		AssociationTypeOne:        "one",
		AssociationTypeMany:       "many",
		AssociationTypeManyToMany: "many-to-many",
	}[a]
}

// Association types
const (
	AssociationTypeUndefined = AssociationType(iota)
	AssociationTypeOne
	AssociationTypeMany
	AssociationTypeManyToMany
)

// Constants
const (
	StructTagName       = "sqlxx"
	SQLXStructTagName   = "db"
	StructTagPrimaryKey = "primary_key"
	StructTagIgnored    = "ignored"
	StructTagDefault    = "default"
	StructTagForeignKey = "fk"
	StructTagSQLXField  = "field"
)

// PrimaryKeyFieldName is the default field name for primary keys
const PrimaryKeyFieldName = "ID"

// SupportedTags are supported tags.
var SupportedTags = []string{
	StructTagName,
	SQLXStructTagName,
}

// TagsMapping is the reflekt.Tags mapping to handle struct tag without key:value format
var TagsMapping = map[string]string{
	"db": "field",
}
