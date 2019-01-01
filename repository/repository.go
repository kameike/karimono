package repository

import (
	"database/sql"

	"github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type InsertAccountRequest struct {
	Id                string
	EncryptedPassword string
}

type UpdateOrReqlaceAccessTokenRequest struct {
	Account  model.Account
	NewToken string
}

type DataRepository interface {
	BeginTransaction()
	EndTransaction()

	InsertAccount(InsertAccountRequest) error
	InsertOrReplaceAccessToken(UpdateOrReqlaceAccessTokenRequest)
}

func CreateApplicationDataRepository() DataRepository {
	db := openDb()

	repo := applicationDataRepository{
		db: db,
	}

	return &repo
}

func openDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/main.db")
	util.CheckInternalFatalError(err)
	return db
}

type applicationDataRepository struct {
	db *sql.DB
}

func (self *applicationDataRepository) BeginTransaction() {
}

func (self *applicationDataRepository) EndTransaction() {
}

func (self *applicationDataRepository) InsertAccount(req InsertAccountRequest) error {
	smit, err := self.db.Prepare("insert into account(name, password_hash) values(?,?)")
	_, err = smit.Exec(req.Id, req.EncryptedPassword)
	if serr, ok := err.(sqlite3.Error); ok && serr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return apperror.ApplicationError{Code: apperror.AccountNameAlreadyTaken}
	}
	util.CheckInternalFatalError(err)
	return nil
}

func (self *applicationDataRepository) InsertOrReplaceAccessToken(req UpdateOrReqlaceAccessTokenRequest) {

	query := `
	insert or replace into access_token (account_id, session_token)
	select id, ? from account where name = ? 
	`
	smit, err := self.db.Prepare(query)
	util.CheckInternalFatalError(err)
	_, err = smit.Exec(req.Account.Id, req.NewToken)
	util.CheckInternalFatalError(err)

}
