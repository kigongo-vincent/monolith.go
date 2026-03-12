package db

import (
	"context"
	"database/sql"
)

// DB is the ORM wrapper interface. Handlers receive this; repos use it for queries.
type DB interface {
	Exec(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) *sql.Row
}

// DBImpl wraps *sql.DB to implement DB.
type DBImpl struct {
	*sql.DB
}

// Exec runs a query without returning rows.
func (d *DBImpl) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return d.DB.ExecContext(ctx, query, args...)
}

// Query runs a query that returns rows.
func (d *DBImpl) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return d.DB.QueryContext(ctx, query, args...)
}

// QueryRow runs a query that returns at most one row.
func (d *DBImpl) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	return d.DB.QueryRowContext(ctx, query, args...)
}

// New opens a database from driver and dsn (e.g. "sqlite", "file:monolith.db").
func New(driver, dsn string) (*DBImpl, error) {
	conn, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DBImpl{DB: conn}, nil
}
