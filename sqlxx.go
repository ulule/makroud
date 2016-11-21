package sqlxx

import (
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oleiade/reflections"
)

// Struct tags
const (
	StructTagName     = "sqlxx"
	SQLXStructTagName = "db"
)

// Regexes
var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// Driver can either be a *sqlx.DB or a *sqlx.Tx.
type Driver interface {
	sqlx.Execer
	sqlx.Queryer
	sqlx.Preparer
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	DriverName() string
	Get(dest interface{}, query string, args ...interface{}) error
	MustExec(query string, args ...interface{}) sql.Result
	NamedExec(query string, arg interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	Preparex(query string) (*sqlx.Stmt, error)
	Rebind(query string) string
	Select(dest interface{}, query string, args ...interface{}) error
}

// Model represents a database table.
type Model interface {
	TableName() string
}

// ModelSchema is a model schema.
type ModelSchema struct {
	Columns      map[string]string
	Associations map[string]RelatedField
}

// RelatedField represents an related field between two models.
type RelatedField struct {
	FK          string
	FKReference string
}

// Preloader is a custom preloader.
type Preloader func(d Driver) (Driver, error)

// GetByParams executes a WHERE with params and populates the given model
// instance with related data.
func GetByParams(driver Driver, out Model, params map[string]string) error {
	return nil
}

// Preload preloads related fields.
func Preload(driver Driver, out Model, related ...string) error {
	return nil
}

// PreloadFuncs preloads with the given preloader functions.
func PreloadFuncs(driver Driver, out Model, preloaders ...Preloader) error {
	return nil
}

// GetModelSchema returns model's table columns, extracted by reflection.
// The returned map is modelFieldName -> table_name.column_name
func GetModelSchema(model Model) (*ModelSchema, error) {
	fields, err := reflections.Fields(model)
	if err != nil {
		return nil, err
	}

	schema := &ModelSchema{
		Columns:      map[string]string{},
		Associations: map[string]RelatedField{},
	}

	for _, field := range fields {
		kind, err := reflections.GetFieldKind(model, field)
		if err != nil {
			return nil, err
		}

		// Associations

		if kind == reflect.Struct || kind == reflect.Ptr {
			relatedField, err := newRelatedField(model, field)
			if err != nil {
				return nil, err
			}

			schema.Associations[field] = *relatedField

			continue
		}

		// Columns

		tag, err := reflections.GetFieldTag(model, field, SQLXStructTagName)
		if err != nil {
			return nil, err
		}

		schema.Columns[field] = columnName(model, field, tag, false, false)
	}

	return schema, nil
}

// newRelatedField creates a new related field.
func newRelatedField(model Model, field string) (*RelatedField, error) {
	relatedValue, err := reflections.GetField(model, field)
	if err != nil {
		return nil, err
	}

	dbTag, err := reflections.GetFieldTag(model, field, SQLXStructTagName)
	if err != nil {
		return nil, err
	}

	tag, err := reflections.GetFieldTag(model, field, StructTagName)
	if err != nil {
		return nil, err
	}

	related := relatedValue.(Model)

	return &RelatedField{
		FK:          columnName(model, field, dbTag, true, false),
		FKReference: columnName(related, field, tag, true, true),
	}, nil
}

// columnName returns full column name from model, field and tag.
func columnName(model Model, field string, tag string, isRelated bool, isReference bool) string {
	// Retrieve the model type
	reflectType := reflect.ValueOf(model).Type()

	// If it's a pointer, we must get the elem to avoid double pointer errors
	if reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	// Then we can safely cast
	reflected := reflect.New(reflectType).Interface().(Model)

	hasTag := len(tag) > 0

	// Build column name from tag or field
	column := tag
	if !hasTag {
		column = toSnakeCase(field)
	}

	// It's not a related field, early return
	if !isRelated {
		return fmt.Sprintf("%s.%s", reflected.TableName(), column)
	}

	// Reference primary key fields are "id" and "field_id"
	if isReference {
		column = "id"

		if hasTag {
			column = tag
		}

		return fmt.Sprintf("%s.%s", reflected.TableName(), column)
	}

	// It's a foreign key
	column = fmt.Sprintf("%s_id", column)
	if hasTag {
		column = tag
	}

	return fmt.Sprintf("%s.%s", reflected.TableName(), column)
}

// toSnakeCase converts camelcased string to snakecase.
func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
