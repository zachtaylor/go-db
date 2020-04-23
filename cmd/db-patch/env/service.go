package env

import (
	"io/ioutil"

	"ztaylor.me/cast"
	"ztaylor.me/db/env"
	"ztaylor.me/log"
)

// PATCH_DIR is name of env var
const PATCH_DIR = "PATCH_DIR"

// Service exports `env.Service`
type Service = env.Service

// NewService returns an `env.Service` with empty settings, and embedded db settings
func NewService() env.Service {
	return env.Service{
		PATCH_DIR: "",
	}.Merge("DB_", env.NewService())
}

// GetPatches scans PATCH_DIR, returns map(patchid->filename)
func GetPatches(env env.Service, logger log.Service) map[int]string {
	patches := make(map[int]string)
	if dir := env[PATCH_DIR]; dir != "" {
		log := logger.New().Add(PATCH_DIR, dir)
		if files, err := ioutil.ReadDir(dir); err != nil {
			log.Add("Error", err).Error("failed to read PATCH_DIR")
		} else {
			for _, f := range files {
				if name := f.Name(); len(name) < 8 {
					// file name too short
				} else if ext := name[len(name)-4:]; ext != ".sql" {
					// file name does not end with ".sql"
				} else if id := cast.IntS(name[:4]); id < 1 {
					log.Add("File", name).Warn("cannot parse patch id")
				} else {
					patches[id] = dir + name
				}
			}
		}
	}
	return patches
}
