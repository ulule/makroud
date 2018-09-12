package sqlxx

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/stmt"
)

// QueryBuilderConfiguration defines a template used by a QueryBuilder to generate queries with loukoum.
type QueryBuilderConfiguration struct {
	model Model
}

// NewQueryBuilder creates a new QueryBuilderConfiguration using given models.
func NewQueryBuilder(model Model) QueryBuilderConfiguration {
	return QueryBuilderConfiguration{
		model: model,
	}
}

// New creates a new QueryBuilder from current configuration.
func (c QueryBuilderConfiguration) New(driver Driver) *QueryBuilder {
	helper := &QueryBuilder{}

	schema, err := GetSchema(driver, c.model)
	if err != nil {
		panic(errors.Wrap(err, "sqlxx: cannot fetch schema informations"))
	}

	table := stmt.NewTable(c.model.TableName())

	columns := []stmt.Column{}
	list := schema.ColumnPaths()
	for i := range list {
		columns = append(columns, stmt.NewColumn(list[i]))
	}

	helper.table = table
	helper.columns = columns
	helper.query = builder.NewSelect().From(helper.table)

	return helper
}

// QueryBuilder is an engine that generate query using a Model and a loukoum builder.
type QueryBuilder struct {
	schema   Schema
	columns  []stmt.Column
	table    stmt.Table
	query    builder.Select
	unscoped bool
	limit    int64
	offset   int64
	cursor   string
	sort     string
	orders   []stmt.Order
}

// Columns returns models columns.
func (b *QueryBuilder) Columns() []stmt.Column {
	return b.columns
}

// Table returns models table.
func (b *QueryBuilder) Table() stmt.Table {
	return b.table
}

// Column adds table namespace to given column name.
func (b *QueryBuilder) Column(name string) string {
	return fmt.Sprint(b.table.Name, ".", name)
}

// Unscope remove filter on archived rows.
func (b *QueryBuilder) Unscope() *QueryBuilder {
	b.unscoped = true
	return b
}

// With adds an equality condition for given key.
func (b *QueryBuilder) With(key string, value interface{}) *QueryBuilder {
	b.Where(loukoum.Condition(b.Column(key)).Equal(value))
	return b
}

// Without adds a not equality condition for given key.
func (b *QueryBuilder) Without(key string, value interface{}) *QueryBuilder {
	b.Where(loukoum.Condition(b.Column(key)).NotEqual(value))
	return b
}

// In adds an equality condition for given key.
func (b *QueryBuilder) In(key string, value ...interface{}) *QueryBuilder {
	b.Where(loukoum.Condition(b.Column(key)).In(value...))
	return b
}

// NotIn adds a not equality condition for given key.
func (b *QueryBuilder) NotIn(key string, value ...interface{}) *QueryBuilder {
	b.Where(loukoum.Condition(b.Column(key)).NotIn(value...))
	return b
}

// Where adds a where condition on query.
func (b *QueryBuilder) Where(condition stmt.Expression) *QueryBuilder {
	b.query = b.query.Where(condition)
	return b
}

// Offset defines query offset.
func (b *QueryBuilder) Offset(offset int64) *QueryBuilder {
	b.offset = offset
	return b
}

// GetOffset returns query offset.
func (b *QueryBuilder) GetOffset() int64 {
	return b.offset
}

// Limit defines query limit.
func (b *QueryBuilder) Limit(limit int64) *QueryBuilder {
	b.limit = limit
	return b
}

// GetLimit returns query limit.
func (b *QueryBuilder) GetLimit() int64 {
	return b.limit
}

// Order defines query order.
func (b *QueryBuilder) Order(orders ...stmt.Order) {
	b.orders = append(b.orders, orders...)
}

func (b *QueryBuilder) Execute(value interface{}) error {
	if !b.unscoped && b.schema.HasDeletedKey() {
		b.query = b.query.Where(loukoum.Condition(b.schema.DeletedKeyPath()).IsNull(true))
	}
	if b.offset != 0 {
		b.query = b.query.Offset(b.offset)
	}
	if b.limit != 0 {
		b.query = b.query.Limit(b.limit)
	}
	if len(b.orders) > 0 {
		b.query = b.query.OrderBy(b.orders...)
	}

	b.query = b.query.Columns(b.Columns())

	// TODO query exec

	return nil
}

// TODO Sort, Pagination (?)

// GetSort returns query sort.
func (b *QueryBuilder) GetSort() string {
	return b.sort
}

// GetCursor returns query cursor.
func (b *QueryBuilder) GetCursor() string {
	return b.cursor
}
