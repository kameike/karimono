package repository

import (
	"database/sql"
	"github.com/kameike/karimono/util"
)

type DataBaseProxy interface {
	BeginTransaction()
	EndTransaction()
	CancelTransaction()
}

type DataRepository interface {
	DataBaseProxy
	AuthDataRepository
	AccountDataRepository
	TeamDataRepository
	BorrowingDataRepository
	HistoryDataRepository
}

func CreateApplicationDataRepository() DataRepository {
	db := openDb()

	repo := applicationDataRepository{
		_db: db,
	}
	return &repo
}

func openDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/main.db")
	util.CheckInternalFatalError(err)
	return db
}

type applicationDataRepository struct {
	_db *sql.DB
	_tx *sql.Tx
}

type queryExecter interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Prepare(query string) (*sql.Stmt, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func (self *applicationDataRepository) db() queryExecter {
	if self._tx != nil {
		return self._tx
	}
	return self._db
}

func (self *applicationDataRepository) BeginTransaction() {
	tx, err := self._db.Begin()
	util.CheckInternalFatalError(err)
	self._tx = tx
}

func (self *applicationDataRepository) EndTransaction() {
	self._tx.Commit()
	self._tx = nil
}

func (self *applicationDataRepository) CancelTransaction() {
	if self._tx == nil {
		return
	}

	self._tx.Rollback()
	self._tx = nil
}
