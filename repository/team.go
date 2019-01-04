package repository

import (
	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/util"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type GetTeamPasswordHashRequest struct {
	TeamName string
}

type CreateTeamAccountReleationRequest struct {
	AccountName string
	TeamName    string
}

type DeleteTeamAccountReleationRequest struct {
	TeamName    string
	AccountName string
}

type TeamDataRepository interface {
	CreateTeamAccountReleation(CreateTeamAccountReleationRequest) error
	DeleteTeamAccountReleation(DeleteTeamAccountReleationRequest) error
	GetTeamPasswordHash(GetTeamPasswordHashRequest) (string, error)
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
