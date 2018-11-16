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

// OpenEnv uses an env.Provider to Open() a database connection
func OpenEnv(env env.Provider) (*DB, error) {
	return Open(BuildDSN(env.Get(DB_USER), env.Get(DB_PASSWORD), env.Get(DB_HOST), env.Get(DB_PORT)), env.Get(DB_TABLE))
}
