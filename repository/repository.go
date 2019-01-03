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
	Token string
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

type CreateTeamAccountHistoryRequest struct {
	TeamName    string
	AccountName string
	History     string
}

type CreateTeamHistoryRequest struct {
	TeamName string
	History  string
}

type CreateAccountHistoryRequest struct {
	AccountName string
	History     string
}

type CreateBorrowingRequest struct {
	ItemName         string
	BorrowingId      string
	BorrwoingAccount string
	BorrwoedTeam     string
	Memo             string
}

type ReturnBorrowingRequest struct {
	BorrowingId string
}

type GetTeamBorrowingRequest struct {
	TeamName string
}

type DataBaseProxy interface {
	BeginTransaction()
	EndTransaction()
	CancelTransaction()
}

type AuthDataReposity interface {
	InsertAccount(InsertAccountRequest) error
	UpdateOrReplaceAccessToken(UpdateOrReqlaceAccessTokenRequest)
	UpdateAccountPassword(UpdateAccountPasswordRequest) error
	UpdateAccountId(UpdateAccountIdRequest) error
	GetAccountWithSecretInfo(GetAccountRequest) (*model.Me, error)

	CheckAuth(AuthCheckRequest) error
	CheckAccountTeamRelation(CheckAccountTeamRelationRequest) error
}

type DataRepository interface {
	AuthDataReposity
	DataBaseProxy

	GetTeamPasswordHash(GetTeamPasswordHashRequest) (string, error)
	CreateTeam(CreateTeamRequest) error
	CreateTeamAccountReleation(CreateTeamAccountReleationRequest) error
	DeleteTeamAccountReleation(DeleteTeamAccountReleationRequest) error
	GetTeam(GetTeamRequest) (*model.Team, error)

	GetAccountHistory(GetAccountHistoryRequest) ([]model.Hisotry, error)
	GetTeamHistory(GetTeamHistoryRequest) ([]model.Hisotry, error)

	CreateTeamAccountHistory(CreateTeamAccountHistoryRequest) error
	CreateTeamHistory(CreateTeamHistoryRequest) error
	CreateAccountHistory(CreateAccountHistoryRequest) error

	CreateBorrowing(CreateBorrowingRequest) error
	ReturnBorrowing(ReturnBorrowingRequest) error

	GetAccountBorrowing(GetAccountBorrowingRequset) ([]model.Borrowing, error)
	GetTeamBorrowing(GetTeamBorrowingRequest) ([]model.Borrowing, error)
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

func (self *applicationDataRepository) CheckAuth(req AuthCheckRequest) error {
	query := `
	select account.name, account.id from access_token join account on access_token.account_id = account.id
	where token = ?
	`

	rows, err := self.db().Query(query, req.AccessToken)
	defer rows.Close()
	util.CheckInternalFatalError(err)

	if rows.Next() {
		return nil
	} else {
		return ApplicationError{ErrorInvalidAccessToken}
	}
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

func (self *applicationDataRepository) GetAccountWithSecretInfo(req GetAccountRequest) (*model.Me, error) {
	query := `
	select account.id, account.name, account.password_hash, access_token.token from account
	join access_token on access_token.account_id = account.id
	where access_token.token = ?
	`
	row := self.db().QueryRow(query, req.Token)

	var account model.Me
	err := row.Scan(&account.Id, &account.Name, &account.PasswordHash, &account.Token)

	if err == sql.ErrNoRows {
		return nil, ApplicationError{ErrorDataNotFount}
	}

	util.CheckInternalFatalError(err)

	return &account, nil
}

func (self *applicationDataRepository) UpdateOrReplaceAccessToken(req UpdateOrReqlaceAccessTokenRequest) {
	query := `
	insert or replace into access_token (account_id, token)
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

func (self *applicationDataRepository) UpdateAccountPassword(req UpdateAccountPasswordRequest) error {
	query := `
	update account set password_hash = ? 
	where name = ?
	`
	result, err := self.db().Exec(query, req.HashedPassword, req.AccountName)
	util.CheckInternalFatalError(err)

	count, err := result.RowsAffected()
	if count == 0 {
		return ApplicationError{ErrorDataNotFount}
	}
	return nil
}

func (self *applicationDataRepository) UpdateAccountId(req UpdateAccountIdRequest) error {
	query := `
	update account set name = ?
	where name = ?
	`
	_, err := self.db().Exec(query, req.NewAccountName, req.OldAccountName)

	if serr, ok := err.(sqlite3.Error); ok && serr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return ApplicationError{ErrorAccountNameAlreadyTaken}
	}
	util.CheckInternalFatalError(err)

	return nil
}

func (self *applicationDataRepository) CreateTeamAccountHistory(req CreateTeamAccountHistoryRequest) error {
	query := `
	insert into history (team_id, account_id, text)
	select team.id, account.id, ? from account
	join team
	where team.name = ? and account.name = ?
	`

	result, err := self.db().Exec(query, req.History, req.TeamName, req.AccountName)
	util.CheckInternalFatalError(err)

	count, err := result.RowsAffected()

	if count == 0 {
		return ApplicationError{ErrorDataNotFount}
	}

	return nil
}

func (self *applicationDataRepository) CreateTeamHistory(req CreateTeamHistoryRequest) error {
	query := `
	insert into history (team_id, text)
	select team.id, ? from team
	where team.name = ?
	`
	result, err := self.db().Exec(query, req.History, req.TeamName)
	util.CheckInternalFatalError(err)

	count, err := result.RowsAffected()

	if count == 0 {
		return ApplicationError{ErrorDataNotFount}
	}

	return nil
}

func (self *applicationDataRepository) CreateAccountHistory(req CreateAccountHistoryRequest) error {
	query := `
	insert into history (account_id, text)
	select account.id, ? from account
	where account.name = ?
	`
	result, err := self.db().Exec(query, req.History, req.AccountName)
	util.CheckInternalFatalError(err)

	count, err := result.RowsAffected()

	if count == 0 {
		return ApplicationError{ErrorDataNotFount}
	}

	return nil
}

func (self *applicationDataRepository) GetTeamHistory(req GetTeamHistoryRequest) ([]model.Hisotry, error) {
	query := `
	select text, created_at from history
	where team_id in
	(
		select id from team
		where name = ?
	)
	`

	rows, err := self.db().Query(query, req.TeamName)
	defer rows.Close()
	util.CheckInternalFatalError(err)

	result := []model.Hisotry{}

	for rows.Next() {
		history := model.Hisotry{}
		rows.Scan(&history.Text, &history.Timestamp)
		result = append(result, history)
	}

	return result, nil
}

func (self *applicationDataRepository) GetAccountHistory(req GetAccountHistoryRequest) ([]model.Hisotry, error) {
	query := `
	select text, created_at from history
	where account_id in
	(
		select id from account
		where name = ?
	)
	`

	rows, err := self.db().Query(query, req.AccountName)
	defer rows.Close()
	util.CheckInternalFatalError(err)

	result := []model.Hisotry{}

	for rows.Next() {
		history := model.Hisotry{}
		rows.Scan(&history.Text, &history.Timestamp)
		result = append(result, history)
	}

	return result, nil
}

func (self *applicationDataRepository) CheckAccountTeamRelation(req CheckAccountTeamRelationRequest) error {
	query := `
	 select 1 from account_team
	 where (team_id, account_id) in
	 (
	 	select team.id, account.id from team
		left join account
	 	where team.name = ? and account.name = ?
	 )
	 `

	rows, err := self.db().Query(query, req.TeamName, req.AccountName)
	util.CheckInternalFatalError(err)

	if rows.Next() {
		return nil
	} else {
		return ApplicationError{ErrorNoRelationBetweenUserAndTeam}
	}
}

func (self *applicationDataRepository) CreateBorrowing(req CreateBorrowingRequest) error {
	query := `
	insert into borrowing(account_id, team_id, hashed_id, name, has_return, memo)
	select account.id, team.id, ?, ?, false, ? from team
	left join account
	where team.name = ? and account.name = ?
	`
	_, err := self.db().Exec(query, req.BorrowingId, req.ItemName, req.Memo, req.BorrwoedTeam, req.BorrwoingAccount)
	util.CheckInternalFatalError(err)

	return nil
}

func (self *applicationDataRepository) GetAccountBorrowing(req GetAccountBorrowingRequset) ([]model.Borrowing, error) {
	query := `
	select borrowing.name, borrowing.hashed_id, account.id, account.name, team.id, team.name from borrowing
	join account on account.id = borrowing.account_id
	join team on team.id = borrowing.team_id
	where account.name = ? and borrowing.has_return = false
	`
	rows, err := self.db().Query(query, req.AccountName)
	defer rows.Close()
	util.CheckInternalFatalError(err)

	var borrowings []model.Borrowing

	for rows.Next() {
		account := model.Account{}
		team := model.Team{}
		borrowing := model.Borrowing{}
		rows.Scan(
			&borrowing.ItemName,
			&borrowing.Uuid,
			&account.Id,
			&account.Name,
			&team.Id,
			&team.Name)
		borrowing.Account = account
		borrowing.Team = team

		borrowings = append(borrowings, borrowing)
	}

	return borrowings, nil
}

func (self *applicationDataRepository) ReturnBorrowing(req ReturnBorrowingRequest) error {
	query := `
	update borrowing set has_return = true where hashed_id = ?
	`
	_, err := self.db().Exec(query, req.BorrowingId)
	util.CheckInternalFatalError(err)

	return nil
}

func (self *applicationDataRepository) GetTeamBorrowing(req GetTeamBorrowingRequest) ([]model.Borrowing, error) {
	query := `
	select borrowing.name, borrowing.hashed_id, account.id, account.name, team.id, team.name from borrowing
	join account on account.id = borrowing.account_id
	join team on team.id = borrowing.team_id
	where team.name = ? and borrowing.has_return = false
	`
	rows, err := self.db().Query(query, req.TeamName)
	defer rows.Close()
	util.CheckInternalFatalError(err)

	var borrowings []model.Borrowing

	for rows.Next() {
		account := model.Account{}
		team := model.Team{}
		borrowing := model.Borrowing{}
		rows.Scan(
			&borrowing.ItemName,
			&borrowing.Uuid,
			&account.Id,
			&account.Name,
			&team.Id,
			&team.Name)
		borrowing.Account = account
		borrowing.Team = team

		borrowings = append(borrowings, borrowing)
	}

	return borrowings, nil
}
