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

// Version is the version of the binary
const Version = "0.0.3"

// PATCH_DIR is name of env var
const PATCH_DIR = "PATCH_DIR"

// ENV is name of env var
const ENV = "env"

// HelpMessage is printed when you use arg "-help" or -"h"
var HelpMessage = `
	db-patch runs database migrations
	internally uses (or can create) table "patch" to manage revision number

	--- options
	[name]			[default]			[comment]

	-help, -h		false				print this help page and then quit

	-env			".env"				env file to be loaded before opening a database connection

	-PATCH_DIR		"./"				directory to load patch files from

	-DB_USER		(required)			username to use when connecting to database

	-DB_PASSWORD		(required)			password to use when conencting to database

	-DB_HOST		(required)			database host ip address

	-DB_PORT		(required)			port to open database host ip with mysql

	-DB_NAME		(required)			database name to connect to within database server
`

func main() {
	env := enviro.NewDefaultService()
	enviro.ParseFlags(env)
	logger := log.StdOutService(log.LevelDebug)
	logger.Formatter().CutSourcePath(0)
	logger.New().With(cast.JSON{
		"module": db.Version,
	}).Debug("db-patch version", Version)

	if env.Default("help", "false") == "true" || env.Default("h", "false") == "true" {
		fmt.Print(HelpMessage)
		return
	}

	if envFile := env.Default(ENV, ".env"); envFile == "" {
		logger.New().Debug("no env")
	} else if err := enviro.ParseFile(env, envFile); err != nil {
		logger.New().With(cast.JSON{
			"envFile": envFile,
			"error":   err,
		}).Error("failed to load environment")
		return
	} else {
		logger.New().Debug("loaded env")
	}

	conn, err := mysql.Open(dbe.BuildDSN(env))
	if conn == nil {
		logger.New().Add("Error", err).Error("failed to open db")
		return
	}
	logger.New().With(cast.JSON{
		dbe.DB_HOST: env.Get(dbe.DB_HOST),
		dbe.DB_NAME: env.Get(dbe.DB_NAME),
	}).Info("opened connection")

	// get current patch info
	patch, err := db.Patch(conn)
	if err == db.ErrPatchTable {
		logger.New().Warn(err.Error())
		ansbuf := "?"
		for ansbuf != "y" && ansbuf != "" && ansbuf != "n" {
			fmt.Print(`patch table does not exist, create patch table now? (y/n): `)
			fmt.Scanln(&ansbuf)
			ansbuf = cast.Trim(ansbuf, " \t")
		}
		if ansbuf == "n" {
			logger.New().Info("exit without creating patch table")
			return
		}
		if err := mysql.CreatePatchTable(conn); err != nil {
			logger.New().Add("Error", err).Error("failed to create patch table")
			return
		}
		logger.New().Info("created patch table")
		patch = 0 // reset patch=-1 from the error state
	} else if err != nil {
		logger.New().Add("Error", err).Error("failed to identify patch number")
		return
	} else {
		logger.New().Info("found patch#", patch)
	}

	patches := getPatches(env, logger)
	if len(patches) < 1 {
		logger.New().Error("no patches found")
		return
	}

	for i := patch + 1; patches[i] != ""; i++ {
		fmt.Println("queue patch#", i, " 	file:", patches[i])
	}

	// ask about patches
	ansbuf := "?"
	for ansbuf != "y" && ansbuf != "" && ansbuf != "n" {
		fmt.Print("Apply patches? [y/n] (default 'y'): ")
		fmt.Scanln(&ansbuf)
		ansbuf = cast.Trim(ansbuf, " \t\n")
	}
	if ansbuf == "n" {
		logger.New().Info("not applying patches")
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
