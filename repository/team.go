package repository

import (
	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type GetTeamPasswordHashRequest struct {
	TeamName string
}

type GetTeamMenbersRequest struct {
	TeamId int
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
	GetTeamMembers(req GetTeamMenbersRequest) ([]model.Account, error)
}

func (self *applicationDataRepository) GetTeamMembers(req GetTeamMenbersRequest) ([]model.Account, error) {
	res := make([]model.Account, 0, 100)

	query := `
	select account.id, account.name from account
	join account_team on account.id = account_team.account_id
	where account_team.team_id = ?
	`
	rows, err := self.db().Query(query, req.TeamId)
	util.CheckInternalFatalError(err)

	for rows.Next() {
		a := model.Account{}
		rows.Scan(&a.Id, &a.Name)
		res = append(res, a)
	}

	return res, nil
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
