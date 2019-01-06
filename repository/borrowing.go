package repository

import (
	"database/sql"

	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
)

type GetAccountBorrowingRequset struct {
	AccountName string
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

type GetBorrowingRequest struct {
	Id string
}

type GetBorrowingRequestWithName struct {
	TeamName    string
	AccountName string
	ItemName    string
}

type BorrowingDataRepository interface {
	CreateBorrowing(CreateBorrowingRequest) error
	ReturnBorrowing(ReturnBorrowingRequest) error
	GetAccountBorrowing(GetAccountBorrowingRequset) ([]model.Borrowing, error)
	GetTeamBorrowing(GetTeamBorrowingRequest) ([]model.Borrowing, error)
	GetBorrowing(GetBorrowingRequest) (*model.Borrowing, error)
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
	select borrowing.memo, borrowing.name, borrowing.hashed_id, account.id, account.name, team.id, team.name from borrowing
	join account on account.id = borrowing.account_id
	join team on team.id = borrowing.team_id
	where account.name = ? and borrowing.has_return = false
	`
	rows, err := self.db().Query(query, req.AccountName)
	defer rows.Close()
	util.CheckInternalFatalError(err)

	borrowings := make([]model.Borrowing, 0, 100)

	for rows.Next() {
		account := model.Account{}
		team := model.Team{}
		borrowing := model.Borrowing{}
		rows.Scan(
			&borrowing.Memo,
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

func (self *applicationDataRepository) GetBorrowingWithName(req GetBorrowingRequestWithName) (*model.Borrowing, error) {
	// query := `
	// 	select borrowing.memo,  borrowing.name, borrowing.hashed_id, account.id, account.name, team.id, team.name from borrowing
	// 	join account on account.id = borrowing.account_id
	// 	join team on team.id = borrowing.team_id
	// 	where borrowing.hashed_id = ?
	// 	`

	return nil, nil
}

func (self *applicationDataRepository) GetBorrowing(req GetBorrowingRequest) (*model.Borrowing, error) {
	query := `
	select borrowing.memo,  borrowing.name, borrowing.hashed_id, account.id, account.name, team.id, team.name from borrowing
	join account on account.id = borrowing.account_id
	join team on team.id = borrowing.team_id
	where borrowing.hashed_id = ?
	`
	row := self.db().QueryRow(query, req.Id)

	account := model.Account{}
	team := model.Team{}
	borrowing := model.Borrowing{}
	err := row.Scan(
		&borrowing.Memo,
		&borrowing.ItemName,
		&borrowing.Uuid,
		&account.Id,
		&account.Name,
		&team.Id,
		&team.Name)
	borrowing.Account = account
	borrowing.Team = team

	if err == sql.ErrNoRows {
		return nil, ApplicationError{ErrorDataNotFount}
	}
	util.CheckInternalFatalError(err)

	return &borrowing, nil
}

func (self *applicationDataRepository) GetTeamBorrowing(req GetTeamBorrowingRequest) ([]model.Borrowing, error) {
	query := `
	select borrowing.memo, borrowing.name, borrowing.hashed_id, account.id, account.name, team.id, team.name from borrowing
	join account on account.id = borrowing.account_id
	join team on team.id = borrowing.team_id
	where team.name = ? and borrowing.has_return = false
	`
	rows, err := self.db().Query(query, req.TeamName)
	defer rows.Close()
	util.CheckInternalFatalError(err)

	borrowings := make([]model.Borrowing, 0, 100)

	for rows.Next() {
		account := model.Account{}
		team := model.Team{}
		borrowing := model.Borrowing{}
		rows.Scan(
			&borrowing.Memo,
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
