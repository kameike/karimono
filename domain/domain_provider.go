package domain

import (
	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/repository"
)

type DomainsProvider interface {
	GetAuthDomain() AuthDomain
	GetAccountDomain() (AccountDomain, error)
	GetTeamDomain(TeamIdProvider) (TeamDomain, error)
	CloseSession()
}

type TokenStatusProvider interface {
	AccountAccessTokenProvider
	HasToken() bool
}

func CreateApplicatoinDomains(token TokenStatusProvider) DomainsProvider {
	repo := repository.CreateApplicationDataRepository()

	authDomain := applicationAuthDomain{
		repo: repo,
	}

	provider := applicationDomainProvider{
		authDomain: &authDomain,
		tokenState: token,
		repository: repo,
	}

	return &provider
}

type applicationDomainProvider struct {
	repository    repository.DataRepository
	tokenState    TokenStatusProvider
	authDomain    AuthDomain
	accountDomain AccountDomain
	teamDomain    TeamDomain
}

func (self *applicationDomainProvider) GetAuthDomain() AuthDomain {
	return self.authDomain
}

func (self *applicationDomainProvider) GetAccountDomain() (AccountDomain, error) {
	if self.accountDomain != nil {
		return self.accountDomain, nil
	}

	if !self.tokenState.HasToken() {
		return nil, ApplicationError{ErrorInvalidAccessToken}
	}

	result, err := self.repository.CheckAuth(repository.AuthCheckRequest{
		AccessToken: self.tokenState.AccountAccessToken(),
	})

	if err != nil {
		return nil, err
	}

	domain := createAccountApplicatoinDomain(*result, self.repository)
	self.accountDomain = domain

	return domain, nil
}

func (self *applicationDomainProvider) GetTeamDomain(req TeamIdProvider) (TeamDomain, error) {
	accountDomain, err := self.GetAccountDomain()
	if err != nil {
		return nil, err
	}

	account := accountDomain.GetAccount()

	r := self.repository

	err = r.CheckAccountTeamRelation(repository.CheckAccountTeamRelationRequest{
		AccountName: account.Name,
		TeamName:    req.TeamId(),
	})

	if err != nil {
		return nil, err
	}

	team, err := r.GetTeam(repository.GetTeamRequest{
		TeamName: req.TeamId(),
	})

	if err != nil {
		return nil, err
	}

	domain := createApplicationTeamDomain(*account, *team, r)

	return domain, nil
}

func (self *applicationDomainProvider) CloseSession() {
}
