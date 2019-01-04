package repository

import (
	"database/sql"

	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type AccountDataRepository interface {
	CreateTeam(CreateTeamRequest) error
	GetTeam(GetTeamRequest) (*model.Team, error)
	GetTeams(GetTeamsRequest) ([]model.Team, error)
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

type GetTeamsRequest struct {
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

func (self *applicationDataRepository) GetTeams(req GetTeamsRequest) ([]model.Team, error) {
	query := `
	select team.id, team.name from team 
	join account_team on team.id = account_team.team_id
	join account on account.id = account_team.account_id
	where account.name = ?
	`

	rows, err := self.db().Query(query, req.TeamName)
	util.CheckInternalFatalError(err)

	var teams []model.Team

	for rows.Next() {
		var team model.Team
		rows.Scan(&team.Id, &team.Name)
		teams = append(teams, team)
	}

	return teams, nil
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

func (self *applicationDataRepository) GetTeamPasswordHash(req GetTeamPasswordHashRequest) (string, error) {
	query := `
	select password_hash from team
	where name = ?
	`
	row := self.db().QueryRow(query, req.TeamName)

	var pass string
	err := row.Scan(&pass)
	util.CheckInternalFatalError(err)

	return pass, nil
}
