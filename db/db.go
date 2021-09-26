package db

import (
	"database/sql"
	"errors"

	goMysql "github.com/go-sql-driver/mysql"
	myMysql "github.com/ziutek/mymysql/mysql"
)

var ErrBadAffectedCount = errors.New("bad affected count")

//CheckAffected checks if result.RowsAffected() are equal to the expected
func CheckAffected(rslt sql.Result, sqlErr error, expected ...int) error {
	if sqlErr != nil {
		return sqlErr
	}

	n, err := rslt.RowsAffected()
	if err != nil {
		return err
	}

	for _, i := range expected {
		if n == int64(i) {
			return nil
		}
	}

	return ErrBadAffectedCount
}

//IsMySQLDuplicate checks if mysql error is ER_DUP_ENTRY mysql error
func IsMySQLDuplicate(err error) bool {
	if val, ok := err.(*myMysql.Error); ok && val.Code == myMysql.ER_DUP_ENTRY {
		return true
	}

	if val, ok := err.(*goMysql.MySQLError); ok && val.Number == 1062 {
		return true
	}

	return false
}
