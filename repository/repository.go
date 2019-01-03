package repository

import (
	"database/sql"

	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type InsertAccountRequest struct {
	Id                string
	EncryptedPassword string
}

type UpdateOrReqlaceAccessTokenRequest struct {
	AccountName string
	NewToken    string
}

type AuthCheckRequest struct {
	AccessToken string
}

type CreateTeamRequest struct {
	EncryptedPassword string
	Name              string
}

type CreateTeamAccountReleationRequest struct {
	AccountName string
	TeamName    string
}

type DeleteTeamAccountReleationRequest struct {
	TeamName    string
	AccountName string
}

type GetTeamPasswordHashRequest struct {
	TeamName string
}

type GetTeamRequest struct {
	TeamName string
}

type UpdateAccountPasswordRequest struct {
	AccountName    string
	HashedPassword string
}

type UpdateAccountIdRequest struct {
	OldAccountName string
	NewAccountName string
}

type GetAccountRequest struct {
	AccountName string
}

type GetAccountHistoryRequest struct {
	AccountName string
}

type GetTeamHistoryRequest struct {
	TeamName string
}

type GetAccountBorrowingRequset struct {
	AccountName string
}

type CheckAccountTeamRelationRequest struct {
	AccountName string
	TeamName    string
}

type DataRepository interface {
	BeginTransaction()
	EndTransaction()
	CancelTransaction()

	InsertAccount(InsertAccountRequest) error
	UpdateOrReplaceAccessToken(UpdateOrReqlaceAccessTokenRequest)
	CheckAuth(AuthCheckRequest) (*model.Account, error)
	UpdateAccountPassword(UpdateAccountPasswordRequest) error
	UpdateAccountId(UpdateAccountIdRequest) error
	GetAccount(GetAccountRequest) (*model.Account, error)

	CheckAccountTeamRelation(CheckAccountTeamRelationRequest) error

	GetTeamPasswordHash(GetTeamPasswordHashRequest) (string, error)
	CreateTeam(CreateTeamRequest) error
	CreateTeamAccountReleation(CreateTeamAccountReleationRequest) error
	DeleteTeamAccountReleation(DeleteTeamAccountReleationRequest) error
	GetTeam(GetTeamRequest) (*model.Team, error)

	GetAccountHistory(GetAccountHistoryRequest) ([]model.Hisotry, error)
	GetTeamHistory(GetTeamHistoryRequest) ([]model.Hisotry, error)
	GetAccountBorrowing(GetAccountBorrowingRequset) ([]model.Borrowing, error)
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
	self._tx.Rollback()
	self._tx = nil
}

func (self *applicationDataRepository) CheckAuth(req AuthCheckRequest) (*model.Account, error) {
	smit, err := self.db().Prepare(`
select account.name, account.id from access_token join account on access_token.account_id = account.id
		where session_token = ?
	`)
	defer smit.Close()

	rows, err := smit.Query(req.AccessToken)
	defer rows.Close()
	util.CheckInternalFatalError(err)

	var account model.Account
	for rows.Next() {
		rows.Scan(&account.Name, &account.Id)
	}

	if account.Name == "" {
		return nil, ApplicationError{ErrorInvalidAccessToken}
	}

	return &account, nil
}
func (self *applicationDataRepository) InsertAccount(req InsertAccountRequest) error {
	if req.Id == "" {
		return ApplicationError{ErrorInvalidAccountName}
	}

	smit, err := self.db().Prepare("insert into account(name, password_hash) values(?,?)")
	_, err = smit.Exec(req.Id, req.EncryptedPassword)
	if serr, ok := err.(sqlite3.Error); ok && serr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return ApplicationError{ErrorAccountNameAlreadyTaken}
	}
	util.CheckInternalFatalError(err)
	return nil
}

func (self *applicationDataRepository) GetAccount(req GetAccountRequest) (*model.Account, error) {
	query := `
		select id, name, password_hash from account where name = ?
	`
	row := self.db().QueryRow(query, req.AccountName)

	var account model.Account
	err := row.Scan(&account.Id, &account.Name, &account.PasswordHash)

	if err == sql.ErrNoRows {
		return nil, ApplicationError{ErrorDataNotFount}
	}

	return &account, nil
}

func (self *applicationDataRepository) UpdateOrReplaceAccessToken(req UpdateOrReqlaceAccessTokenRequest) {
	query := `
	insert or replace into access_token (account_id, session_token)
	select id, ? from account where name = ? 
	`
	smit, err := self.db().Prepare(query)
	util.CheckInternalFatalError(err)
	_, err = smit.Exec(req.NewToken, req.AccountName)
	util.CheckInternalFatalError(err)
}

func (self *applicationDataRepository) GetTeamPasswordHash(GetTeamPasswordHashRequest) (string, error) {
	return "", nil
}

func (self *applicationDataRepository) CreateTeam(req CreateTeamRequest) error {
	query := `
	insert into team (name, password_hash) values (?, ?)
	`
	_, err := self.db().Exec(query, req.Name, req.EncryptedPassword)
	if serr, ok := err.(sqlite3.Error); ok && serr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return ApplicationError{ErrorTeamNameAlreadyTaken}
	}
	util.CheckInternalFatalError(err)

	return nil
}

func (self *applicationDataRepository) CreateTeamAccountReleation(req CreateTeamAccountReleationRequest) error {
	query := `
	insert into account_team (account_id, team_id)
	select * from 
		(select id as account_id from account where name = ?)
		join
		(select id as team_id from team where name = ?)
	`
	result, err := self.db().Exec(query, req.AccountName, req.TeamName)
	if serr, ok := err.(sqlite3.Error); ok && serr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return ApplicationError{ErrorAlreadyJoin}
	}
	util.CheckInternalFatalError(err)

	effectedCount, err := result.RowsAffected()
	util.CheckInternalFatalError(err)
	if effectedCount == 0 {
		return ApplicationError{ErrorDataNotFount}
	}

	return nil
}

func (self *applicationDataRepository) DeleteTeamAccountReleation(req DeleteTeamAccountReleationRequest) error {
	query := `
	delete from account_team where id in
	(select account_team.id from account_team
		join account on account_team.account_id = account.id
		join team on account_team.team_id = team.id
	where team.name = ? and account.name = ?)
	`
	_, err := self.db().Exec(query, req.TeamName, req.AccountName)
	util.CheckInternalFatalError(err)

	return nil
}

func (self *applicationDataRepository) GetTeam(req GetTeamRequest) (*model.Team, error) {
	query := `
	select name, id from team where name = ?
	`
	row := self.db().QueryRow(query, req.TeamName)

	team := &model.Team{}
	err := row.Scan(&team.Name, &team.Id)

	if err == sql.ErrNoRows {
		return nil, ApplicationError{ErrorDataNotFount}
	}
	util.CheckInternalFatalError(err)

	return team, nil
}

func (self *applicationDataRepository) UpdateAccountPassword(UpdateAccountPasswordRequest) error {
	return nil
}

func (self *applicationDataRepository) UpdateAccountId(UpdateAccountIdRequest) error {
	return nil
}

func (self *applicationDataRepository) GetAccountHistory(GetAccountHistoryRequest) ([]model.Hisotry, error) {
	return nil, nil
}

func (self *applicationDataRepository) GetTeamHistory(GetTeamHistoryRequest) ([]model.Hisotry, error) {
	return nil, nil
}

func (self *applicationDataRepository) GetAccountBorrowing(GetAccountBorrowingRequset) ([]model.Borrowing, error) {
	return nil, nil
}

func (self *applicationDataRepository) CheckAccountTeamRelation(CheckAccountTeamRelationRequest) error {
	return nil
}
