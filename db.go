package db // import "ztaylor.me/db"

import (
	"database/sql"

	"ztaylor.me/cast"
)

// Version is the module version
const Version = "v0.0.9"

// ErrPatchTable is returned by Patch when the patch table doesn't exist
var ErrPatchTable = cast.NewError(nil, `table "patch" does not exist`)

// ErrSQLPanic is returned by ExecTx when it encounters a panic
var ErrSQLPanic = cast.NewError(nil, `sql panic`)

// ErrTxEmpty is returned by ExecTx when tx has no statements
var ErrTxEmpty = cast.NewError(nil, `patch file contains no statements`)

// DB == sql.DB
type DB = sql.DB

// Result == sql.Result
type Result = sql.Result

// Scanner provides a header for generic SQL data set
type Scanner interface {
	Scan(...interface{}) error
}

// Service is a database driver
type Service interface {
	Open(string) (*DB, error)
}

// Patch returns the current patch number for the database
//
// returns -1, ErrPatchTable if the table doesn't exist
func Patch(db *DB) (int, error) {
	if patch, err := scanPatch(db); err == nil {
		return patch, nil
	} else if e := err.Error(); len(e) > 10 && e[:10] == "Error 1146" {
		return -1, ErrPatchTable
	} else {
		return -1, err
	}
}
func scanPatch(db *DB) (patch int, err error) {
	err = db.QueryRow("SELECT * FROM patch").Scan(&patch)
	return
}

// ExecTx uses database transaction to apply SQL statements
func ExecTx(db *DB, sql string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	isEmpty := true
	defer func() {
		if p := recover(); p != nil {
			err = ErrSQLPanic
		}
		if isEmpty {
			err = ErrTxEmpty
		}
		if err == nil {
			err = tx.Commit()
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	for _, stmt := range cast.Split(sql, `;`) {
		if stmt = cast.Trim(stmt, "\n\r \t"); stmt != "" {
			isEmpty = false
			if _, err = tx.Exec(stmt); err != nil {
				break
			}
		}
	}
	return err
}

// BuildDSN returns a formatted DSN string
func BuildDSN(user, password, host, port, table string) string {
	return user + `:` + password + `@tcp(` + host + `:` + port + `)/` + table
}
