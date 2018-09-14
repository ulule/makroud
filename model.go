package makroud

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
