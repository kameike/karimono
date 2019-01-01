package domain

import (
	"database/sql"

	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/repository"
	"github.com/kameike/karimono/util"
	"golang.org/x/crypto/bcrypt"
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

type AccountCreateRequester interface {
	AccountPasswordProvider
	AccountIdProvider
}

type AccountDescriber interface {
	AccountIdProvider
}

type AccountAccessTokenProvider interface {
	Token() string
}

type AuthDomain interface {
	CreateAccount(AccountCreateRequester) (*model.Account, error)
	UpdateAccountPassword(AccountPasswordProvider) (*model.Account, error)
	UpdateAccountId(AccountIdProvider) (*model.Account, error)
	GetAccount() (*model.Account, error)
	RenewAccessToken() (*model.Account, error)
}

type applicationAuthDomain struct {
	repo repository.DataRepository
}

func (self *applicationAuthDomain) CreateAccount(req AccountCreateRequester) (*model.Account, error) {
	self.repo.BeginTransaction()
	defer self.repo.EndTransaction()

	createReq := repository.InsertAccountRequest{
		Id:                req.AccountId(),
		EncryptedPassword: encryptPassword(req.AccountPassword()),
	}

	err := self.repo.InsertAccount(createReq)
	if err != nil {
		return nil, err
	}

	token := newToken()

	account := model.Account{
		Token: token,
		Id:    req.AccountId(),
	}

	self.repo.InsertOrReplaceAccessToken(repository.UpdateOrReqlaceAccessTokenRequest{
		Account:  account,
		NewToken: token,
	})

	return &account, nil
}

func (self *applicationAuthDomain) UpdateAccountPassword(rep AccountPasswordProvider) (*model.Account, error) {
	return nil, nil
}

func (self *applicationAuthDomain) UpdateAccountId(req AccountIdProvider) (*model.Account, error) {
	return nil, nil
}

func encryptPassword(plainPass string) string {
	pass, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
	util.CheckInternalFatalError(err)

	return string(pass)
}

func newToken() string {
	return util.RandString(100)
}

func (self *applicationAuthDomain) GetAccount() (*model.Account, error) {
	return nil, nil
}

func (self *applicationAuthDomain) RenewAccessToken() (*model.Account, error) {
	return nil, nil
}

func renewToken(name string, db *sql.DB) (string, error) {
	token := util.RandString(100)

	return token, nil
}
