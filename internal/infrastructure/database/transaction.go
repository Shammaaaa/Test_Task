package database

import (
	"context"
	"database/sql"
)

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type TxFn func(Transaction) error

func WithTransaction(ctx context.Context, db *sql.DB, fn TxFn) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {

			_ = tx.Rollback()
			panic(p)
		} else if err != nil {

			_ = tx.Rollback()
		} else {

			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
