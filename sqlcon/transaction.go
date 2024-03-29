package sqlcon

import (
	"github.com/jmoiron/sqlx"
)

// A Txfn is a function that will be called with an initialized `Transaction` object
// that can be used for executing statements and queries against a database.
type TxFn func(tx *sqlx.Tx) error

// WithTransaction creates a new transaction and handles rollback/commit based on the
// error object returned by the `TxFn`
func WithTransaction(db *sqlx.DB, fn TxFn) (err error) {
	tx, err := db.Beginx()
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
