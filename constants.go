package sqlxx

import (
	"database/sql"
	"reflect"

	"github.com/lib/pq"
)

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

// NullFieldTypes are field considered as NULL.
var NullFieldTypes = map[reflect.Type]bool{
	reflect.TypeOf(sql.NullBool{}):    true,
	reflect.TypeOf(sql.NullFloat64{}): true,
	reflect.TypeOf(sql.NullInt64{}):   true,
	reflect.TypeOf(sql.NullString{}):  true,
	reflect.TypeOf(pq.NullTime{}):     true,
}

// PrimaryKeyFieldName is the implicit primary key field name.
const PrimaryKeyFieldName = "ID"
