package mimic

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strconv"
	"sync"

	"github.com/go-xorm/core"
)

// Forked from https://github.com/volatiletech/boilbench

func init() {
	sql.Register("mimic", &mimic{})
}

var (
	mutex   = sync.Mutex{}
	dsns    = map[string]QueryResult{}
	counter = 0
)

type QueryResult struct {
	*Result
	*Query
	NumInput int
}

type Result struct {
	NumRows int
}

type Query struct {
	Cols []string
	Vals [][]driver.Value
}

type mimic struct {
}

func (m *mimic) Open(dsn string) (driver.Conn, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if len(dsn) == 0 {
		dsn = strconv.Itoa(counter)
		counter++
	}

	conn := &mimicConn{
		dsns[dsn],
	}

	return conn, nil
}

type mimicConn struct {
	Q QueryResult
}

func (m *mimicConn) Prepare(query string) (driver.Stmt, error) {
	return &mimicStmt{m.Q}, nil
}

func (m *mimicConn) Close() error              { return nil }
func (m *mimicConn) Begin() (driver.Tx, error) { return nil, errors.New("tx not supported") }

type mimicStmt struct {
	Q QueryResult
}

func (m *mimicStmt) Close() error  { return nil }
func (m *mimicStmt) NumInput() int { return m.Q.NumInput }
func (m *mimicStmt) Exec(args []driver.Value) (driver.Result, error) {
	if m.Q.Result == nil {
		return nil, errors.New("statement was not a result type")
	}

	return &mimicResult{m.Q.Result.NumRows}, nil
}

func (m *mimicStmt) Query(args []driver.Value) (driver.Rows, error) {
	if m.Q.Query == nil {
		return nil, errors.New("statement was not a query type")
	}

	return &mimicRows{columns: m.Q.Query.Cols, values: m.Q.Query.Vals}, nil
}

type mimicResult struct {
	rowsAffected int
}

func (m *mimicResult) LastInsertId() (int64, error) {
	return int64(m.rowsAffected), nil
}

func (m *mimicResult) RowsAffected() (int64, error) {
	return int64(m.rowsAffected), nil
}

type mimicRows struct {
	cursor  int
	columns []string
	values  [][]driver.Value
}

func (m *mimicRows) Columns() []string { return m.columns }
func (m *mimicRows) Close() error      { return nil }
func (m *mimicRows) Next(dest []driver.Value) error {
	if m.cursor == len(m.values) {
		return io.EOF
	}

	for i, val := range m.values[m.cursor] {
		dest[i] = val
	}
	m.cursor++

	return nil
}

func NewResult(q QueryResult) string {
	mutex.Lock()
	defer mutex.Unlock()
	dsn := strconv.Itoa(counter)
	counter++
	dsns[dsn] = q
	return dsn
}

func NewQuery(q QueryResult) string {
	mutex.Lock()
	defer mutex.Unlock()
	dsn := strconv.Itoa(counter)
	counter++
	dsns[dsn] = q
	return dsn
}

func NewResultDSN(dsn string, q QueryResult) {
	mutex.Lock()
	defer mutex.Unlock()
	dsns[dsn] = q
}

func NewQueryDSN(dsn string, q QueryResult) {
	mutex.Lock()
	defer mutex.Unlock()
	dsns[dsn] = q
}

type XormDriver struct {
}

func (x *XormDriver) Parse(a string, b string) (*core.Uri, error) {
	return &core.Uri{DbType: core.POSTGRES}, nil
}
