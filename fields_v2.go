package sqlxx

type XField struct {
	// Model name that contains this field
	ModelName string
	// Table name of the model that contains this field
	TableName string
	// The field name
	FieldName string
	// The database column name
	ColumnName string
	// Does this field is a primary key?
	IsPrimaryKey bool
	// TODO
	IsArchiveKey bool
	// IsForeignKey define if field should behave like a foreign key
	IsForeignKey bool
	// Default define the default statement to use if value in model is undefined.
	Default string
}

// FieldOption is used to define field configuration.
type FieldOption interface {
	apply(*fieldOp)
}

// fieldOption is a private implementation of FieldOption using a callback to configure a Field.
type fieldOption func(*fieldOp)

func (o fieldOption) apply(op *fieldOp) {
	o(op)
}

func IsForeignKey(reference string) FieldOption {
	return fieldOption(func(op *fieldOp) {
		op.IsForeignKey = true
		op.ReferenceName = reference
	})
}

func IsArchiveKey() FieldOption {
	return fieldOption(func(op *fieldOp) {
		op.IsArchiveKey = true
	})
}

// HasDefault define a default value for a field.
func HasDefault(value string) FieldOption {
	return fieldOption(func(op *fieldOp) {
		op.Default = value
	})
}

// NOTE
// ArchiveKey -> define a generator so we can use boolean or time.Time, for example...

type XPrimaryKey struct {
	ModelName  string
	TableName  string
	FieldName  string
	ColumnName string
	Type       PrimaryKeyType
}

// ForeignKey is a composite object that define a foreign key for a model.
//
// For example: If we have an User with an Avatar, we could have this foreign key defined in User's schema.
//
//     ForeignKey {
//         ModelName: User,
//         TableName: users,
//         FieldName: AvatarID,
//         ColumnName: avatar_id,
//         ReferenceName: Avatar,
//     }
//
type XForeignKey struct {
	// ModelName define the model name of this foreign key.
	ModelName string
	// TableName define the database table.
	TableName string
	// FieldName define the struct field name used as foreign key.
	FieldName string
	// ColumnName define the database column name of this foreign key.
	ColumnName string
	// ReferenceName define the reference model name.
	ReferenceName string
}

type XAssociationReference struct {
	// ModelName define the model name of this association.
	ModelName string
	// FieldName define the struct field name used to contains the value(s) of this association.
	FieldName string
	// Source define the field used to ... TODO
	Source XField
	// Reference define the field used to ... TODO
	Reference XField
	// Type define the association type. (ie: has once, has many, etc...)
	Type AssociationType
}
