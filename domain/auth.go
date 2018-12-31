package domain

import (
	"github.com/kameike/karimono/model"
)

// Interfaces
type AccountPasswordProvider interface {
	AccountPassword() string
}

type AccountIdProvider interface {
	AccountId() string
}

type AccountIdUpdateRequester interface {
	AccountIdProvider
}

type AccountPasswordUpdateRequester interface {
	AccountPasswordProvider
}

type AuthDomain interface {
	CreateAccount() (model.Account, error)
	UpdateAccountPassword(AccountPasswordProvider) (model.Account, error)
	UpdateAccountId(AccountIdProvider) (model.Account, error)
	GetAccount() (model.Account, error)
}
