package domain

import (
	. "github.com/kameike/karimono/error"
	"github.com/kameike/karimono/repository"
)

type DomainsProvider interface {
	GetAuthDomain() AuthDomain
	GetAccountDomain() (AccountDomain, error)
	GetTeamDomain(TeamIdProvider) (TeamDomain, error)
	GetTeamProviderViaTeamName(req string) (TeamDomain, error)
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

	err := self.repository.CheckAuth(repository.AuthCheckRequest{
		AccessToken: self.tokenState.AccountAccessToken(),
	})

	if err != nil {
		return nil, err
	}

	result, err := self.repository.GetAccountWithSecretInfo(repository.GetAccountRequest{
		Token: self.tokenState.AccountAccessToken(),
	})
	if err != nil {
		return nil, err
	}

	domain := createAccountApplicatoinDomain(*result, self.repository)
	self.accountDomain = domain

	return domain, nil
}

type internalTeamIdprovider struct {
	id int
}

func (p internalTeamIdprovider) TeamId() int {
	return p.id
}

func (self *applicationDomainProvider) GetTeamProviderViaTeamName(req string) (TeamDomain, error) {
	r := self.repository

	t, err := r.GetTeam(repository.GetTeamRequest{
		TeamName: req,
	})

	if err != nil {
		return nil, err
	}

	return self.GetTeamDomain(internalTeamIdprovider{t.Id})
}

func (self *applicationDomainProvider) GetTeamDomain(req TeamIdProvider) (TeamDomain, error) {
	accountDomain, err := self.GetAccountDomain()
	if err != nil {
		return nil, err
	}

	account := accountDomain.GetAccount()

	r := self.repository

	team, err := r.GetTeamWithId(repository.GetTeamWithIdRequest{
		Id: req.TeamId(),
	})

	if err != nil {
		return nil, err
	}

	err = r.CheckAccountTeamRelation(repository.CheckAccountTeamRelationRequest{
		AccountName: account.Name,
		TeamName:    team.Name,
	})

	if err != nil {
		return nil, err
	}

	domain := createApplicationTeamDomain(*account, *team, r)

	return domain, nil
}

func (self *applicationDomainProvider) CloseSession() {
}
