package repository

import (
	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type AccountDataRepository interface {
	CreateTeam(CreateTeamRequest) error
	GetTeam(GetTeamRequest) (*model.Team, error)
	UpdateOrReplaceAccessToken(UpdateOrReqlaceAccessTokenRequest)
	UpdateAccountId(UpdateAccountIdRequest) error
	UpdateAccountPassword(UpdateAccountPasswordRequest) error
}

type UpdateAccountPasswordRequest struct {
	AccountName    string
	HashedPassword string
}

type UpdateAccountIdRequest struct {
	OldAccountName string
	NewAccountName string
}

type CreateTeamRequest struct {
	EncryptedPassword string
	Name              string
}

type GetTeamRequest struct {
	TeamName string
}
type UpdateOrReqlaceAccessTokenRequest struct {
	AccountName string
	NewToken    string
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
