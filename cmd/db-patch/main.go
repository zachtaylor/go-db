package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"ztaylor.me/cast"
	"ztaylor.me/db"
	main_env "ztaylor.me/db/cmd/db-patch/env"
	db_env "ztaylor.me/db/env"
	"ztaylor.me/db/mysql"
	"ztaylor.me/log"
)

// HelpMessage is printed when you use arg "-help" or -"h"
var HelpMessage = `
	db-patch runs database migrations
	internally uses (or can create) table "patch" to manage revision number

	--- options
	[name]			[default]			[comment]

	-help, -h		false				print this help page and then quit

	-PATCH_DIR		"./"				directory to load patch files from

	-DB_USER		(required)			username to use when connecting to database

	-DB_PASSWORD		(required)			password to use when conencting to database

	-DB_HOST		(required)			database host ip address

	-DB_PORT		(required)			port to open database host ip with mysql

	-DB_NAME		(required)			database name to connect to within database server
`

func main() {
	env := main_env.NewService().ParseDefault()
	logger := log.StdOutService(log.LevelDebug)
	logger.Formatter().CutSourcePath(0)
	logger.New().With(cast.JSON{
		"DB_NAME":   env["DB_NAME"],
		"PATCH_DIR": env[main_env.PATCH_DIR],
	}).Debug("db-patch", db.Version)

	if cast.Bool(env["help"]) || cast.Bool(env["h"]) {
		fmt.Print(HelpMessage)
		return
	}

	conn, err := mysql.Open(db_env.BuildDSN(env.Match("DB_")))
	if conn == nil {
		logger.New().Add("Error", err).Error("failed to open db")
		return
	}
	logger.New().With(cast.JSON{
		db_env.HOST: env[db_env.HOST],
		db_env.NAME: env[db_env.NAME],
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

	patches := main_env.GetPatches(env, logger)
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
