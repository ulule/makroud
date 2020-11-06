// Package hooks provides hooks into the driver.Connector. It is useful for
// instrumenting.
package hooks

import (
	"context"
	"database/sql/driver"
)

// Wrap returns a new Connector wrapping c.
func Wrap(c driver.Connector) *Connector { return &Connector{wrapped: c} }

// A Connector wraps an existing connector.
type Connector struct {
	BeforeConnect func(ctx context.Context) context.Context
	AfterConnect  func(ctx context.Context, conn driver.Conn, err error)

	BeforeExec func(ctx context.Context, query string, args []driver.NamedValue) context.Context
	AfterExec  func(ctx context.Context, result driver.Result, err error)

	BeforeQuery func(ctx context.Context, query string, args []driver.NamedValue) context.Context
	AfterQuery  func(ctx context.Context, rows driver.Rows, err error)

	BeforeBegin func(ctx context.Context, opts driver.TxOptions) context.Context
	AfterBegin  func(ctx context.Context, tx driver.Tx, err error)

	BeforeCommit func(ctx context.Context) context.Context
	AfterCommit  func(ctx context.Context, err error)

	BeforeRollback func(ctx context.Context) context.Context
	AfterRollback  func(ctx context.Context, err error)

	wrapped driver.Connector
}

// Connect implements database/sql/driver.Connector.
func (connector *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	if connector.BeforeConnect != nil {
		ctx = connector.BeforeConnect(ctx)
	}
	c, err := connector.wrapped.Connect(ctx)
	if connector.AfterConnect != nil {
		connector.AfterConnect(ctx, c, err)
	}
	return &conn{
		wrapped:        c,
		BeforeExec:     connector.BeforeExec,
		AfterExec:      connector.AfterExec,
		BeforeQuery:    connector.BeforeQuery,
		AfterQuery:     connector.AfterQuery,
		BeforeBegin:    connector.BeforeBegin,
		AfterBegin:     connector.AfterBegin,
		BeforeCommit:   connector.BeforeCommit,
		AfterCommit:    connector.AfterCommit,
		BeforeRollback: connector.BeforeRollback,
		AfterRollback:  connector.AfterRollback,
	}, err
}

// Driver implements database/sql/driver.Connector.
func (connector *Connector) Driver() driver.Driver { return connector.wrapped.Driver() }

type conn struct {
	wrapped driver.Conn

	BeforeExec     func(ctx context.Context, query string, args []driver.NamedValue) context.Context
	AfterExec      func(ctx context.Context, result driver.Result, err error)
	BeforeQuery    func(ctx context.Context, query string, args []driver.NamedValue) context.Context
	AfterQuery     func(ctx context.Context, rows driver.Rows, err error)
	BeforeBegin    func(ctx context.Context, opts driver.TxOptions) context.Context
	AfterBegin     func(ctx context.Context, tx driver.Tx, err error)
	BeforeCommit   func(ctx context.Context) context.Context
	AfterCommit    func(ctx context.Context, err error)
	BeforeRollback func(ctx context.Context) context.Context
	AfterRollback  func(ctx context.Context, err error)
}

func (c *conn) Begin() (driver.Tx, error) {
	return c.wrapped.Begin()
}

func (c *conn) Close() error {
	return c.wrapped.Close()
}

func (c *conn) Prepare(query string) (driver.Stmt, error) {
	return c.wrapped.Prepare(query)
}

var (
	_ driver.ExecerContext  = &conn{}
	_ driver.QueryerContext = &conn{}
	_ driver.ConnBeginTx    = &conn{}
)

func (c *conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	if c.BeforeExec != nil {
		ctx = c.BeforeExec(ctx, query, args)
	}
	result, err := c.wrapped.(driver.ExecerContext).ExecContext(ctx, query, args)
	if c.AfterExec != nil {
		c.AfterExec(ctx, result, err)
	}
	return result, err
}

func (c *conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	if c.BeforeQuery != nil {
		ctx = c.BeforeQuery(ctx, query, args)
	}
	rows, err := c.wrapped.(driver.QueryerContext).QueryContext(ctx, query, args)
	if c.AfterQuery != nil {
		c.AfterQuery(ctx, rows, err)
	}
	return rows, err
}

func (c *conn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if c.BeforeBegin != nil {
		ctx = c.BeforeBegin(ctx, opts)
	}
	t, err := c.wrapped.(driver.ConnBeginTx).BeginTx(ctx, opts)
	if c.AfterBegin != nil {
		c.AfterBegin(ctx, t, err)
	}
	return &tx{
		wrapped:        t,
		ctx:            ctx,
		BeforeCommit:   c.BeforeCommit,
		AfterCommit:    c.AfterCommit,
		BeforeRollback: c.BeforeRollback,
		AfterRollback:  c.AfterRollback,
	}, err
}

type tx struct {
	wrapped driver.Tx
	ctx     context.Context

	BeforeCommit   func(ctx context.Context) context.Context
	AfterCommit    func(ctx context.Context, err error)
	BeforeRollback func(ctx context.Context) context.Context
	AfterRollback  func(ctx context.Context, err error)
}

func (tx *tx) Commit() error {
	ctx := tx.ctx
	if tx.BeforeCommit != nil {
		ctx = tx.BeforeCommit(ctx)
	}
	err := tx.wrapped.Commit()
	if tx.AfterCommit != nil {
		tx.AfterCommit(ctx, err)
	}
	return err
}

func (tx *tx) Rollback() error {
	ctx := tx.ctx
	if tx.BeforeRollback != nil {
		ctx = tx.BeforeRollback(ctx)
	}
	err := tx.wrapped.Rollback()
	if tx.AfterRollback != nil {
		tx.AfterRollback(ctx, err)
	}
	return err
}
