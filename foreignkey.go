package sqlxx

import "strings"

// ForeignKey contains foreign key field information.
type ForeignKey struct {
	// The foreign key field name
	Name string
	// The foreign key column name
	ColumnName string
	// The field name of the foreign key reference
	RelatedFieldName string
	// The column name of the foreign key reference
	RelatedColumnName string
}

// NewForeignKey returns a ForeignKey instance.
func NewForeignKey(f Field) *ForeignKey {
	return &ForeignKey{
		Name:       f.Name,
		ColumnName: f.ColumnName,
	}
}

// IsForeignKey returns true if the given fields looks like a foreign key or
// had been explicitly set as foreign key field.
func IsForeignKey(f Field) bool {
	if f.Tags.HasKey(StructTagName, "fk") {
		return true
	}

	// Typically MyFieldID/MyFieldPK
	if len(f.Name) > len(PrimaryKeyFieldName) && strings.HasSuffix(f.Name, PrimaryKeyFieldName) {
		return true
	}

	return false
}
