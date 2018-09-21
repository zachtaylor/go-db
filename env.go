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
	return BuildDSN(env.Get(DB_USER), env.Get(DB_PASSWORD), env.Get(DB_HOST), env.Get(DB_PORT))
}

// OpenEnv uses Open and env to get the database connection settings
func OpenEnv() (*DB, error) {
	return Open(envDataSourceName(), env.Get(DB_TABLE))
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
