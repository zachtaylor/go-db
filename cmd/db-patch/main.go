package main

import (
	"io/ioutil"
	"time"

	"ztaylor.me/cast"
	"ztaylor.me/db"
	"ztaylor.me/env"
	"ztaylor.me/log"
)

// PATCH_DIR is name of env var
const PATCH_DIR = "PATCH_DIR"

func main() {
	env := env.Global()
	conn, err := db.OpenEnv(env)
	logger := log.StdOutService(log.LevelDebug)
	if conn == nil {
		logger.New().Add("Error", err).Error("failed to open db")
		return
	}
	logger.New().With(log.Fields{
		"DB_HOST": env.Get(db.DB_HOST),
		"DB_NAME": env.Get(db.DB_NAME),
	}).Info("starting...")

	// get current patch info
	patch, err := db.Patch(conn)
	if err == db.ErrPatchTable {
		logger.New().Warn(err.Error())
		if err := db.CreatePatchTable(conn); err != nil {
			logger.New().Add("Error", err).Error("failed to create patch table")
			return
		}
		logger.New().Info("created patch table")
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
		log := logger.New().With(log.Fields{
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
			logger.New().With(log.Fields{
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
					logger.New().With(log.Fields{
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
