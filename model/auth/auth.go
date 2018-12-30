package auth

import (
	"github.com/kameike/karimono/model"
)

type ContextProvidor interface {
}

type AuthDomain interface {
	CreateAccount() (model.Account, error)
	UpdateAccount() (model.Account, error)
	GetAccount() (model.Account, error)
}

type TeamDomain interface {
	CreateTeam() (model.Team, error)
	UpdateTeam() (model.Team, error)
	GetBorrowings() ([]model.Hisotry, error)
	GetHistory() ([]model.Hisotry, error)
}

type AccountDomain interface {
	JoinTeam() (model.Team, error)
	LeaveTeam() error
	GetHistory() ([]model.Hisotry, error)
	GetBorrowings() ([]model.Hisotry, error)
}

type BorrowingDomain interface {
	BorrowItem(itemName string) (model.Borrowing, error)
	ReturnItem() (model.Borrowing, error)
}

type DomainsProvider interface {
	GetAuthDomain() AuthDomain
	GetAccountDomain() AccountDomain
	GetTeamDomain(model.Team) TeamDomain
	GetBorrowingDomain(model.Team) BorrowingDomain
	CloseSession()
}

// type ApplicationDomainProvidor struct{}
// type MockDomainProvidor struct{}
//
// func (p ApplicationDomainProvidor) CreateNewSession() Domains {
//
// 	return Domains{
// 		AuthDomain: nil,
// 	}
// }
//
// func (p ApplicationDomainProvidor) CloseSession() {
// }
