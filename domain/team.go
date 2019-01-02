package domain

import (
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/repository"
)

// Interfaces
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

type TeamDomain interface {
	UpdateTeamPassword(TeamPasswordUpdateRequester) (*model.Team, error)
	UpdateTeamName(TeamNameUpdateRequester) (*model.Team, error)
	KickAccount(AccountDescriable) error
	GetTeamInfo() (*model.Team, error)
	GetHistories() ([]model.Hisotry, error)
}

type applicationTeamDomain struct {
	account    model.Account
	team       model.Team
	repository repository.DataRepository
}

func createApplicationTeamDomain(account model.Account, team model.Team, repository repository.DataRepository) TeamDomain {
	domain := applicationTeamDomain{
		account:    account,
		team:       team,
		repository: repository,
	}

	return &domain
}
func (self *applicationTeamDomain) UpdateTeamPassword(TeamPasswordUpdateRequester) (*model.Team, error) {
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
