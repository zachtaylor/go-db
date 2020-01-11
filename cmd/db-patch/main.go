package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"ztaylor.me/cast"
	"ztaylor.me/db"
	dbe "ztaylor.me/db/env"
	"ztaylor.me/db/mysql"
	"ztaylor.me/env"
	enviro "ztaylor.me/env"
	"ztaylor.me/log"
)

// PATCH_DIR is name of env var
const PATCH_DIR = "PATCH_DIR"

// Version is the version of the binary
const Version = "0.0.2"

func main() {
	env := enviro.NewDefaultService()
	enviro.ParseFlags(env)

	if env.Default(`version`, `false`) == `true` || env.Default(`v`, `false`) == `true` {
		fmt.Printf(`db-patch version ` + Version + ` db version ` + db.Version)
		return
	}

	enviro.ParseFile(env, ".env")
	conn, err := mysql.Open(dbe.BuildDSN(env))
	logger := log.StdOutService(log.LevelDebug)
	if conn == nil {
		logger.New().Add("Error", err).Error("failed to open db")
		return
	}
	logger.New().With(cast.JSON{
		dbe.DB_HOST: env.Get(dbe.DB_HOST),
		dbe.DB_NAME: env.Get(dbe.DB_NAME),
		PATCH_DIR:   env.Get(PATCH_DIR),
		"Version":   db.Version,
	}).Info("starting")

	// get current patch info
	patch, err := db.Patch(conn)
	if err == db.ErrPatchTable {
		logger.New().Warn(err.Error())
		if err := mysql.CreatePatchTable(conn); err != nil {
			logger.New().Add("Error", err).Error("failed to create patch table")
			return
		}
		logger.New().Info("created patch table")
		patch = 0 // reset patch=-1 during error
	} else if err != nil {
		logger.New().Add("Error", err).Error("failed to identify patch number")
		return
	}
	logger.New().Info("found patch#" + cast.StringI(patch))

	patches := getPatches(env, logger)
	if len(patches) < 1 {
		logger.New().Error("no patches found")
		return
	}

	// apply patches
	for patch++; patches[patch] != ""; patch++ {
		patchFile := patches[patch]
		log := logger.New().With(cast.JSON{
			"PatchID":   patch,
			"PatchFile": patchFile,
		})
		tStart := time.Now()
		sql, err := ioutil.ReadFile(patchFile)
		if err = db.ExecTx(conn, string(sql)); err != nil {
			log.Add("Error", err).Error("failed to patch")
			return
		} else if _, err = conn.Exec("UPDATE patch SET patch=?", patch); err != nil {
			log.Add("Error", err).Error("failed to update patch number")
			return
		}
		log.Add("Time", time.Now().Sub(tStart)).Info("applied patch")
	}

	logger.New().Add("Patch", patch-1).Info("done")
}

// getPatches scans PATCH_DIR, returns map(patchid->filename)
func getPatches(env env.Service, logger log.Service) map[int]string {
	patches := make(map[int]string)
	if dir := env.Get(PATCH_DIR); dir != "" {
		if files, err := ioutil.ReadDir(dir); err != nil {
			logger.New().With(cast.JSON{
				"PATCH_DIR": dir,
				"Error":     err,
			}).Error("failed to read PATCH_DIR")
		} else {
			for _, f := range files {
				if name := f.Name(); len(name) < 8 {
					// file name too short
				} else if ext := name[len(name)-4:]; ext != ".sql" {
					// file name does not end with ".sql"
				} else if id := cast.IntS(name[:4]); id < 1 {
					logger.New().With(cast.JSON{
						"PATCH_DIR": dir,
						"File":      name,
					}).Warn("cannot parse patch id")
				} else {
					patches[id] = dir + name
				}
			}
		}
	}
	return patches
}
