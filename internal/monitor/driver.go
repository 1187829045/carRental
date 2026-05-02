package monitor

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

const DriverName = "monitor-mysql"

func init() {
	sql.Register(DriverName, &monitorDriver{base: &mysql.MySQLDriver{}})
}

type monitorDriver struct {
	base driver.Driver
}

func (d *monitorDriver) Open(name string) (driver.Conn, error) {
	conn, err := d.base.Open(name)
	if c := Default(); c != nil {
		c.RecordConnect(err)
	}
	if err != nil {
		return nil, err
	}
	return &monitorConn{Conn: conn}, nil
}

type monitorConn struct {
	driver.Conn
}

func (c *monitorConn) Close() error {
	err := c.Conn.Close()
	if err == nil {
		if col := Default(); col != nil {
			col.RecordClose()
		}
	}
	return err
}

func (c *monitorConn) Ping(ctx context.Context) error {
	if p, ok := c.Conn.(driver.Pinger); ok {
		return p.Ping(ctx)
	}
	return nil
}

func (c *monitorConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if bt, ok := c.Conn.(driver.ConnBeginTx); ok {
		return bt.BeginTx(ctx, opts)
	}
	if opts.ReadOnly || opts.Isolation != driver.IsolationLevel(0) {
		return nil, errors.New("unsupported transaction options")
	}
	return c.Conn.Begin()
}

func (c *monitorConn) ResetSession(ctx context.Context) error {
	if rs, ok := c.Conn.(driver.SessionResetter); ok {
		return rs.ResetSession(ctx)
	}
	return nil
}

func (c *monitorConn) CheckNamedValue(nv *driver.NamedValue) error {
	if checker, ok := c.Conn.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(nv)
	}
	return driver.ErrSkip
}

func (c *monitorConn) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if pc, ok := c.Conn.(driver.ConnPrepareContext); ok {
		stmt, err := pc.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		return &monitorStmt{Stmt: stmt, query: query}, nil
	}
	stmt, err := c.Conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &monitorStmt{Stmt: stmt, query: query}, nil
}

func (c *monitorConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	ec, ok := c.Conn.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	done := beginSQL(query)
	res, err := ec.ExecContext(ctx, query, args)
	done(ctx, err, resultRowsAffected(res))
	return res, err
}

func (c *monitorConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	qc, ok := c.Conn.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	done := beginSQL(query)
	rows, err := qc.QueryContext(ctx, query, args)
	done(ctx, err, 0)
	return rows, err
}

type monitorStmt struct {
	driver.Stmt
	query string
}

func (s *monitorStmt) Exec(args []driver.Value) (driver.Result, error) {
	done := beginSQL(s.query)
	res, err := s.Stmt.Exec(args)
	done(context.Background(), err, resultRowsAffected(res))
	return res, err
}

func (s *monitorStmt) Query(args []driver.Value) (driver.Rows, error) {
	done := beginSQL(s.query)
	rows, err := s.Stmt.Query(args)
	done(context.Background(), err, 0)
	return rows, err
}

func (s *monitorStmt) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	if se, ok := s.Stmt.(driver.StmtExecContext); ok {
		done := beginSQL(s.query)
		res, err := se.ExecContext(ctx, args)
		done(ctx, err, resultRowsAffected(res))
		return res, err
	}
	return nil, driver.ErrSkip
}

func (s *monitorStmt) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	if sq, ok := s.Stmt.(driver.StmtQueryContext); ok {
		done := beginSQL(s.query)
		rows, err := sq.QueryContext(ctx, args)
		done(ctx, err, 0)
		return rows, err
	}
	return nil, driver.ErrSkip
}

func beginSQL(query string) func(context.Context, error, int64) {
	if col := Default(); col != nil {
		return col.BeginSQL(query)
	}
	start := time.Now()
	return func(context.Context, error, int64) {
		_ = start
	}
}

func resultRowsAffected(res driver.Result) int64 {
	if res == nil {
		return 0
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0
	}
	return n
}
