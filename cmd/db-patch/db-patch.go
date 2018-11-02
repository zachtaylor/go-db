package main

import (
	"errors"
	"io/ioutil"
	"strconv"
	"time"

	"ztaylor.me/db"
	"ztaylor.me/env"
	"ztaylor.me/log"
)

// PATCH_DIR is name of env var
const PATCH_DIR = "PATCH_DIR"

var errNotSQL = errors.New("file is not .sql type")

func main() {
	log.SetLevel("debug")
	log.Add("DB_HOST", env.Get(db.DB_HOST)).Add("DB_TABLE", env.Get(db.DB_TABLE)).Debug("db-patch: starting...")

	conn, err := db.OpenEnv()
	if err != nil {
		log.Add("Error", err).Error("db-patch: error opening db")
		return
	}

	patch, err := db.Patch(conn)
	if err == db.ErrPatchTable {
		if err := db.CreatePatchTable(conn); err != nil {
			log.Add("Error", err).Error("db-patch: failed to create patch table")
			return
		}
		log.Info("db-patch: created patch table")
	} else if err != nil {
		log.Add("Error", err).Error("db-patch: unrecognized error")
		return
	}
	log.Info("db-patch: current patch#" + strconv.Itoa(patch))

	patch++
	for patches := getPatches(); patches[patch] != ""; patch++ {
		patchFile := patches[patch]
		log := log.Add("PatchID", patch).Add("PatchFile", patchFile)
		tStart := time.Now()
		sql, err := ioutil.ReadFile(patchFile)

		if err = db.ExecTx(conn, string(sql)); err != nil {
			log.Add("Error", err).Error("db-patch: failed")
			return
		} else if _, err = conn.Exec("UPDATE patch SET patch=?", patch); err != nil {
			log.Add("Error", err).Error("db-patch: failed to update patch version")
			return
		}

		log.Add("Time", time.Now().Sub(tStart)).Info("db-patch: patch applied")
	}

	log.Info("db-patch: finished")
}

func getPatches() map[int]string {
	patches := make(map[int]string)
	if dir := env.Get(PATCH_DIR); dir != "" {
		if files, err := ioutil.ReadDir(dir); err != nil {
			log.Add("PatchDir", dir).Add("Error", err).Error("db-patch: cannot read patch path")
		} else {
			for _, f := range files {
				if id, err := parsePatchID(f.Name()); err != nil {
					log.Add("PatchDir", dir).Add("File", f.Name()).Add("Error", err).Warn("db-patch: cannot parse patch id")
				} else {
					patches[id] = dir + f.Name()
				}
			}
		}
	} else {
		log.Error("db-patch: no patches to apply")
	}
	return patches
}

func parsePatchID(fileName string) (int, error) {
	if fileName[len(fileName)-4:] != ".sql" {
		log.Warn("sqlite-patch: file is not .sql type")
		return 0, errNotSQL
	}
	patchid, err := strconv.ParseInt(fileName[:4], 10, 64)
	return int(patchid), err
}
