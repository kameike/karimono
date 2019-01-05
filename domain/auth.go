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

type NameCheckRequester interface {
	AccountIdProvider
}

type AuthDomain interface {
	CreateAccount(AccountCreateRequester) (*model.Me, error)
	SignInAccount(AccountSignInRequester) (*model.Me, error)
	CheckNameAvailable(NameCheckRequester) bool
}

type applicationAuthDomain struct {
	repo repository.DataRepository
}

func (self *applicationAuthDomain) CreateAccount(req AccountCreateRequester) (*model.Me, error) {
	createReq := repository.InsertAccountRequest{
		Id:                req.AccountId(),
		EncryptedPassword: hashPassword(req.AccountPassword()),
	}

	err := self.repo.InsertAccount(createReq)
	if err != nil {
		return nil, err
	}

	token := newToken()

	self.repo.UpdateOrReplaceAccessToken(repository.UpdateOrReqlaceAccessTokenRequest{
		AccountName: req.AccountId(),
		NewToken:    token,
	})

	account, err := self.repo.GetAccountWithSecretInfo(repository.GetAccountRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}

	return account, nil
}
func (self *applicationAuthDomain) CheckNameAvailable(req NameCheckRequester) bool {
	r := self.repo

	res, err := r.GetAccountWithName(repository.GetAccountRequestWithName{
		Name: req.AccountId(),
	})

	if serr, ok := err.(ApplicationError); ok && serr.Code != ErrorDataNotFount {
		return false
	}

	if res == nil {
		return true
	}

	if err != nil {
		return false
	}

	return false
}

func (self *applicationAuthDomain) SignInAccount(req AccountSignInRequester) (*model.Me, error) {
	r := self.repo

	me, err := r.GetAccountWithName(repository.GetAccountRequestWithName{
		Name: req.AccountId(),
	})
	if err != nil {
		return nil, err
	}

	err = checkPasswordHash(req.AccountPassword(), me.PasswordHash)

	if err != nil {
		return nil, err
	}

	token := newToken()

	self.repo.UpdateOrReplaceAccessToken(repository.UpdateOrReqlaceAccessTokenRequest{
		AccountName: req.AccountId(),
		NewToken:    token,
	})

	me, err = self.repo.GetAccountWithSecretInfo(repository.GetAccountRequest{
		Token: token,
	})
	if err != nil {
		return nil, err
	}
	return me, nil
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
	return util.RandString(128)
}
