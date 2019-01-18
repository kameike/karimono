package domain

import (
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/repository"
	"github.com/kameike/karimono/util"
)

type TeamIdProvider interface {
	TeamId() int
}

type TeamNameProvider interface {
	TeamName() string
}

type TeamDescriable interface {
	TeamNameProvider
	TeamIdProvider
}

type TeamPasswordUpdateRequester interface {
	TeamPasswordProvider
}

type TeamNameUpdateRequester interface {
	TeamNameProvider
}

type AccountDescriable interface {
	AccountIdProvider
}

type BorrowMemoProvider interface {
	Memo() string
}

type BorrowItemProvider interface {
	ItemName() string
}

type BorrowItemRequester interface {
	BorrowMemoProvider
	BorrowItemProvider
}

type TeamDomainRequester interface {
	TeamIdProvider
}

type TeamDomain interface {
	UpdateTeamPassword(TeamPasswordUpdateRequester) (*model.Team, error)
	UpdateTeamName(TeamNameUpdateRequester) (*model.Team, error)
	KickAccount(AccountDescriable) error
	GetTeamInfo() (*model.Team, error)
	GetTeamBorrowings() ([]model.Borrowing, error)
	GetHistories() ([]model.Hisotry, error)
	BorrowItem(BorrowItemRequester) (*model.Borrowing, error)
	GetTeamMenbers() ([]model.Account, error)
}

type applicationTeamDomain struct {
	account    model.Me
	team       model.Team
	repository repository.DataRepository
}

func createApplicationTeamDomain(account model.Me, team model.Team, repository repository.DataRepository) TeamDomain {
	domain := applicationTeamDomain{
		account:    account,
		team:       team,
		repository: repository,
	}

	return &domain
}

func (self *applicationTeamDomain) UpdateTeamPassword(TeamPasswordUpdateRequester) (*model.Team, error) {
	// r := self.repository
	return nil, nil

}
func (self *applicationTeamDomain) UpdateTeamName(TeamNameUpdateRequester) (*model.Team, error) {
	return nil, nil
}
func (self *applicationTeamDomain) KickAccount(AccountDescriable) error {
	return nil
}
func (self *applicationTeamDomain) GetTeamInfo() (*model.Team, error) {
	return nil, nil
}
func (self *applicationTeamDomain) GetHistories() ([]model.Hisotry, error) {
	return nil, nil
}

func (self *applicationTeamDomain) BorrowItem(req BorrowItemRequester) (*model.Borrowing, error) {

	id := randomBorrowingId()

	r := self.repository
	err := r.CreateBorrowing(repository.CreateBorrowingRequest{
		BorrowingId:      id,
		BorrwoedTeam:     self.team.Name,
		BorrwoingAccount: self.account.Name,
		ItemName:         req.ItemName(),
		Memo:             req.Memo(),
	})

	if err != nil {
		return nil, err
	}

	borrowing, err := r.GetBorrowing(repository.GetBorrowingRequest{
		Id: id,
	})

	if err != nil {
		return nil, err
	}

	return borrowing, nil
}

func (self *applicationTeamDomain) GetTeamMenbers() ([]model.Account, error) {
	r := self.repository
	return r.GetTeamMembers(repository.GetTeamMenbersRequest{self.team.Id})
}

func (self *applicationTeamDomain) GetTeamBorrowings() ([]model.Borrowing, error) {
	r := self.repository

	return r.GetTeamBorrowing(repository.GetTeamBorrowingRequest{
		TeamName: self.team.Name,
	})
}

func randomBorrowingId() string {
	return util.RandString(48)
}
