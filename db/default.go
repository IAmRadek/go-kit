package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"sync"
	"time"
)

var dbConn *sql.DB
var dbConnOnce sync.Once

func Init(driverName, dataSourceName string) (err error) {
	dbConnOnce.Do(func() {
		dbConn, err = sql.Open(driverName, dataSourceName)
	})
	return
}

func Begin() (*sql.Tx, error) { return dbConn.Begin() }
func BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return dbConn.BeginTx(ctx, opts)
}
func Close() error                                               { return dbConn.Close() }
func Conn(ctx context.Context) (*sql.Conn, error)                { return dbConn.Conn(ctx) }
func Driver() driver.Driver                                      { return dbConn.Driver() }
func Exec(query string, args ...interface{}) (sql.Result, error) { return dbConn.Exec(query, args...) }
func ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return dbConn.ExecContext(ctx, query, args...)
}
func Ping() error                             { return dbConn.Ping() }
func PingContext(ctx context.Context) error   { return dbConn.PingContext(ctx) }
func Prepare(query string) (*sql.Stmt, error) { return dbConn.Prepare(query) }
func PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return dbConn.PrepareContext(ctx, query)
}
func Query(query string, args ...interface{}) (*sql.Rows, error) { return dbConn.Query(query, args...) }
func QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return dbConn.QueryContext(ctx, query, args...)
}
func QueryRow(query string, args ...interface{}) *sql.Row { return dbConn.QueryRow(query, args...) }
func QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return dbConn.QueryRowContext(ctx, query, args...)
}
func SetConnMaxLifetime(d time.Duration) { dbConn.SetConnMaxLifetime(d) }
func SetMaxIdleConns(n int)              { dbConn.SetMaxIdleConns(n) }
func SetMaxOpenConns(n int)              { dbConn.SetMaxOpenConns(n) }
func Stats() sql.DBStats                 { return dbConn.Stats() }
