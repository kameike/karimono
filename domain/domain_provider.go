package domain

import "github.com/kameike/karimono/repository"

type DomainsProvider interface {
	GetAuthDomain() AuthDomain
	GetAccountDomain() (AccountDomain, error)
	GetTeamDomain(TeamIdProvider) (TeamDomain, error)
	GetBorrowingDomain(TeamIdProvider) (BorrowingDomain, error)
	CloseSession()
}

func CreateApplicatoinDomains() DomainsProvider {
	repo := repository.CreateApplicationDataRepository()

	authDomain := applicationAuthDomain{
		repo: repo,
	}

	provider := applicationDomainProvider{
		authDomain: &authDomain,
	}

	return &provider
}

type applicationDomainProvider struct {
	authDomain      AuthDomain
	accountDomain   AccountDomain
	teamDomain      TeamDomain
	borrowingDomain BorrowingDomain
}

func (self *applicationDomainProvider) GetAuthDomain() AuthDomain {
	return self.authDomain
}

func (self *applicationDomainProvider) GetAccountDomain() (AccountDomain, error) {
	return nil, nil
}

func (self *applicationDomainProvider) GetTeamDomain(TeamIdProvider) (TeamDomain, error) {
	return nil, nil
}

func (self *applicationDomainProvider) GetBorrowingDomain(TeamIdProvider) (BorrowingDomain, error) {
	return nil, nil
}

func (self *applicationDomainProvider) CloseSession() {
}
