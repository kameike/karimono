package repository

import (
	"database/sql"

	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type CheckAccountTeamRelationRequest struct {
	AccountName string
	TeamName    string
}

type AuthCheckRequest struct {
	AccessToken string
}

type GetAccountRequest struct {
	Token string
}

type InsertAccountRequest struct {
	Id                string
	EncryptedPassword string
}

type AuthDataRepository interface {
	CheckAccountTeamRelation(CheckAccountTeamRelationRequest) error
	CheckAuth(AuthCheckRequest) error
	GetAccountWithSecretInfo(GetAccountRequest) (*model.Me, error)
	InsertAccount(InsertAccountRequest) error
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
