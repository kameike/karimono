package domain

import (
	. "github.com/kameike/karimono/error"
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

type AccountCreateRequester interface {
	AccountPasswordProvider
	AccountIdProvider
}

type AccountSignInRequester interface {
	AccountPasswordProvider
	AccountIdProvider
}

type AccountDescriber interface {
	AccountIdProvider
}

type AccountAccessTokenProvider interface {
	AccountAccessToken() string
}

type AuthDomain interface {
	CreateAccount(AccountCreateRequester) (*model.Account, error)
	SignInAccount(AccountSignInRequester) (*model.Account, error)
}

type applicationAuthDomain struct {
	repo repository.DataRepository
}

func (self *applicationAuthDomain) CreateAccount(req AccountCreateRequester) (*model.Account, error) {
	createReq := repository.InsertAccountRequest{
		Id:                req.AccountId(),
		EncryptedPassword: hashPassword(req.AccountPassword()),
	}

	err := self.repo.InsertAccount(createReq)
	if err != nil {
		return nil, err
	}

	token := newToken()

	account := model.Account{
		Token: token,
		Name:  req.AccountId(),
	}

	self.repo.UpdateOrReplaceAccessToken(repository.UpdateOrReqlaceAccessTokenRequest{
		AccountName: account.Name,
		NewToken:    token,
	})

	return &account, nil
}

func (self *applicationAuthDomain) SignInAccount(req AccountSignInRequester) (*model.Account, error) {
	return nil, nil
}

func hashPassword(plainPass string) string {
	pass, err := bcrypt.GenerateFromPassword([]byte(plainPass), bcrypt.DefaultCost)
	util.CheckInternalFatalError(err)

	return string(pass)
}

func checkPasswordHash(plainPass string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPass))

	if err != nil {
		return ApplicationError{ErrorInvalidTeamPassword}
	}
	return nil
}

func newToken() string {
	return util.RandString(100)
}
