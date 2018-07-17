package sqlrepo

import (
	"context"
	"database/sql"
)

//Repo provides utility methods to help working with the SQL package.
type Repo struct {
	db *sql.DB

	d Dialect
}

//New returns a new Repo that connects to db and uses the Dialect d.
func New(db *sql.DB, d Dialect) *Repo {
	return &Repo{
		db: db,
		d:  d,
	}
}

//QueryContext is the QueryerContext implementation.
//It calls QueryContext on r's underlying sql.DB.
func (r *Repo) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return r.db.QueryContext(ctx, Normalize(r.d, query), args...)
}

//QueryRowContext is the QueryerContext implementation.
//It calls QueryRowContext on r's underlying sql.DB.
func (r *Repo) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return r.db.QueryRowContext(ctx, Normalize(r.d, query), args...)
}

//ExecContext is the QueryExecerContext implementation.
//It calls ExecContext on r's underlying sql.DB.
func (r *Repo) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return r.db.ExecContext(ctx, Normalize(r.d, query), args...)
}

//BeginContext starts a transaction in r.
//The interface returned wraps the QueryExecerContext and the necessary
//methods on the sql.Tx type.
//it calls BeginTx on r's underlying sql.DB.
func (r *Repo) BeginContext(ctx context.Context) (Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

//TxWorkContext executes work inside of a transaction with ctx with committing and
//rollback handled for you.
func (r *Repo) TxWorkContext(ctx context.Context, work func(QueryExecerContext) error) error {
	sqlTx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	err = work(&tx{sqlTx, r.d})
	if err != nil {
		sqlTx.Rollback()
		return err
	}
	return sqlTx.Commit()
}

//Scanner is a wrapper interface around sql.Rows and sql.Row.
type Scanner interface {
	//Scan is the method to get a row into data.
	Scan(dst ...interface{}) error
}

//QueryerContext provides the query methods with Contexts for a sql database.
type QueryerContext interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)

	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

//ExecerContext provides the execute methods with Contexts for a sql database.
type ExecerContext interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

//QueryExecerContext provides the query and exec with Contexts methods for a sql
//database.
type QueryExecerContext interface {
	QueryerContext
	ExecerContext
}

//Tx is a QueryExecer with the Commit and Rollback methods required for transactions.
//It also provides a Stmt method so prepared statements can be executed within
//this transaction.
type Tx interface {
	//QueryExecerContext provides the query and exec context methods.
	QueryExecerContext

	//Commit commits the transaction.
	Commit() error

	//Rollback rolls back the transaction.
	Rollback() error
}

//tx provides transaction support with a dialect.
type tx struct {
	sqlTx *sql.Tx

	d Dialect
}

func (t *tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return t.sqlTx.QueryContext(ctx, Normalize(t.d, query), args...)
}

func (t *tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return t.sqlTx.QueryRowContext(ctx, Normalize(t.d, query), args...)
}

func (t *tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return t.sqlTx.ExecContext(ctx, Normalize(t.d, query), args...)
}

func (t *tx) Commit() error {
	return t.sqlTx.Commit()
}

func (t *tx) Rollback() error {
	return t.sqlTx.Rollback()
}
