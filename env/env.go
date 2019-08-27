package dbe

import (
	"ztaylor.me/db"
	"ztaylor.me/env"
)

// DB_USER is name of env var
const DB_USER = "DB_USER"

// DB_PASSWORD is name of env var
const DB_PASSWORD = "DB_PASSWORD"

// DB_HOST is name of env var
const DB_HOST = "DB_HOST"

// DB_PORT is name of env var
const DB_PORT = "DB_PORT"

// DB_NAME is name of env var
const DB_NAME = "DB_NAME"

// BuildDSN uses env.Service to build database DSN
func BuildDSN(env env.Service) string {
	return db.BuildDSN(env.Get(DB_USER), env.Get(DB_PASSWORD), env.Get(DB_HOST), env.Get(DB_PORT), env.Get(DB_NAME))
}
