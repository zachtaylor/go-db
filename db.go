package db // import "ztaylor.me/db"

import (
	"database/sql"
	"errors"
	"strings"

	// package db only uses mysql
	_ "github.com/go-sql-driver/mysql"
)

// ErrPatchTable is returned by Patch when the patch table doesn't exist
var ErrPatchTable = errors.New("patch table does not exist")

// ErrSQLPanic is returned by ExecTx when it encounters a panic
var ErrSQLPanic = errors.New("sql panic")

// ErrTxEmpty is returned by ExecTx when tx has no statements
var ErrTxEmpty = errors.New("tx is empty")

// DB == sql.DB
type DB = sql.DB

// Result == sql.Result
type Result = sql.Result

// Scanner provides a header for generic SQL data set
type Scanner interface {
	Scan(...interface{}) error
}

// New creates a db connection using mysql
func New(dataSourceName string) (*DB, error) {
	return sql.Open("mysql", dataSourceName)
}

// Use prepares a db connection by issuing SQL "USE" command
func Use(db *DB, table string) (Result, error) {
	return db.Exec("USE " + table)
}

// Open creates a DB connection for the dsn and table name
func Open(dsn string, table string) (*DB, error) {
	if db, err := New(dsn); err != nil {
		return nil, err
	} else if _, err = Use(db, table); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

// Patch returns the current patch number for the database
func Patch(db *DB) (int, error) {
	if patch, err := scanPatch(db); err == nil {
		return patch, nil
	} else if e := err.Error(); len(e) > 10 && e[:10] == "Error 1146" {
		return -1, ErrPatchTable
	} else {
		return -1, err
	}
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
	for _, stmt := range strings.Split(sql, `;`) {
		if stmt = strings.Trim(stmt, "\n\r \t"); stmt != "" {
			isEmpty = false
			if _, err = tx.Exec(stmt); err != nil {
				break
			}
		}
	}
	return err
}

func scanPatch(db *DB) (int, error) {
	var patch int
	row := db.QueryRow("SELECT * FROM patch")
	err := row.Scan(&patch)
	if err != nil {
		return 0, err
	}
	return patch, nil
}

// BuildDSN returns a formatted DSN string
func BuildDSN(user, password, host, port string) string {
	return user + `:` + password + `@tcp(` + host + `:` + port + `)/`
}

// CreatePatchTable creates the patch table
func CreatePatchTable(db *DB) error {
	return ExecTx(db, `CREATE TABLE patch(patch INTEGER) ENGINE=InnoDB; INSERT INTO patch (patch) VALUES (0);`)
}
