package handler

import (
	"database/sql"

	"github.com/kameike/karimono/util"
	_ "github.com/mattn/go-sqlite3"
)

func openDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/main.db")
	util.CheckInternalFatalError(err)
	return db
}
