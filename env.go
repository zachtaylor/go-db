package db

import "ztaylor.me/env"

// DB_USER is name of env var
const DB_USER = "DB_USER"

// DB_PASSWORD is name of env var
const DB_PASSWORD = "DB_PASSWORD"

// DB_HOST is name of env var
const DB_HOST = "DB_HOST"

// DB_PORT is name of env var
const DB_PORT = "DB_PORT"

// DB_TABLE is name of env var
const DB_TABLE = "DB_TABLE"

func envDataSourceName() string {
	user := env.Get(DB_USER)
	password := env.Get(DB_PASSWORD)
	host := env.Get(DB_HOST)
	port := env.Get(DB_PORT)
	return user + `:` + password + `@tcp(` + host + `:` + port + `)/`
}

// OpenEnv uses Open and env to get the database connection settings
func OpenEnv(table string) (*DB, error) {
	return Open(envDataSourceName(), table)
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
