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

// OpenEnv uses env to get the database connection settings
func OpenEnv() (*DB, error) {
	if db, err := New(envDataSourceName()); err != nil {
		return nil, err
	} else if _, err = Use(db, env.Get(DB_TABLE)); err != nil {
		return nil, err
	} else {
		return db, nil
	}
}
