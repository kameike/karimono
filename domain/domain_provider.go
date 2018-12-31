package domain

import ()

type DomainsProvider interface {
	GetAuthDomain() AuthDomain
	GetAccountDomain() (AccountDomain, error)
	GetTeamDomain(TeamIdProvider) (TeamDomain, error)
	GetBorrowingDomain(TeamIdProvider) (BorrowingDomain, error)
	CloseSession()
}
