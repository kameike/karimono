package repository

import (
	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
)

type GetAccountHistoryRequest struct {
	AccountName string
}

type GetTeamHistoryRequest struct {
	TeamName string
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

type HistoryDataRepository interface {
	CreateAccountHistory(CreateAccountHistoryRequest) error
	CreateTeamAccountHistory(CreateTeamAccountHistoryRequest) error
	CreateTeamHistory(CreateTeamHistoryRequest) error
	GetAccountHistory(GetAccountHistoryRequest) ([]model.Hisotry, error)
	GetTeamHistory(GetTeamHistoryRequest) ([]model.Hisotry, error)
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
