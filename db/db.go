package db

import (
	"database/sql"
	"errors"
)

var ErrBadAffectedCount = errors.New("bad affected count")

// CheckAffected checks if result.RowsAffected() are equal to the expected.
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
