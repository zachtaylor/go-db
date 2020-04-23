package env // import "ztaylor.me/db/env"

import (
	"ztaylor.me/db"
	"ztaylor.me/env"
)

// USER is name of env var
const USER = "USER"

// PASSWORD is name of env var
const PASSWORD = "PASSWORD"

// HOST is name of env var
const HOST = "HOST"

// PORT is name of env var
const PORT = "PORT"

// NAME is name of env var
const NAME = "NAME"

// Service exports `env.Service`
type Service = env.Service

// NewService returns an `env.Service` with empty settings
func NewService() env.Service {
	return env.Service{
		USER:     "",
		PASSWORD: "",
		HOST:     "",
		PORT:     "",
		NAME:     "",
	}
}

// BuildDSN uses env.Service to build database DSN
func BuildDSN(env env.Service) string {
	return db.BuildDSN(env[USER], env[PASSWORD], env[HOST], env[PORT], env[NAME])
}
