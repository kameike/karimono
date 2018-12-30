package handler

import (
	"database/sql"
	"github.com/kameike/karimono/util"
	"testing"
)

func Test(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	util.CheckInternalFatalError(err)
	print(db)
}
