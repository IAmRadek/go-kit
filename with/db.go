package with

import (
	"context"
	"database/sql"
	"errors"
)

//ErrRollback is used stop with.Tx
var ErrRollback = errors.New("ErrRollback")

//TxFunc is func used to perform inside transaction
type TxFunc func(tx *sql.Tx) error

//Tx (with.Tx()) encapsulates common pattern of transaction
func Tx(dbConn *sql.DB, fn TxFunc, opts ...*sql.TxOptions) error {
	return TxContext(context.Background(), dbConn, fn, opts...)
}

//TxContext encapsulates common pattern of transaction
func TxContext(ctx context.Context, dbConn *sql.DB, fn TxFunc, opts ...*sql.TxOptions) (err error) {
	var opt *sql.TxOptions = nil
	if len(opts) > 0 {
		opt = opts[0]
	}

	tx, err := dbConn.BeginTx(ctx, opt)
	if err != nil {
		return
	}

	//defer is used in case fn panics
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
		return
	}()

	err = fn(tx)

	return
}
