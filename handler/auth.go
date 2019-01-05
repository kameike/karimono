package handler

import (
	"github.com/kameike/karimono/model"
)

func (self *Handler) updateAccount() error {
	return nil
}

type AccountAuthorizeRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type accountAuthorizedResponse struct {
	AccessToken string        `json:"accessToken"`
	Account     model.Account `json:"account"`
}

func (self AccountAuthorizeRequest) AccountId() string {
	return self.Name
}

func (self AccountAuthorizeRequest) AccountPassword() string {
	return self.Password
}

type accountNameCheckRequest struct {
	Name string `json:"name"`
}

type accountNameCheckResponse struct {
	Availrable bool `json:"available"`
}

func (r accountNameCheckRequest) AccountId() string {
	return r.Name
}

func (self *Handler) validateAccount() error {
	var req accountNameCheckRequest
	self.bodyAsJson(&req)

	a := self.provider.GetAuthDomain()
	res := a.CheckNameAvailable(req)

	self.renderJson(accountNameCheckResponse{
		Availrable: res,
	})

	return nil
}

func (self *Handler) createAccount() error {
	var res AccountAuthorizeRequest
	self.bodyAsJson(&res)

	authDomain := self.provider.GetAuthDomain()
	me, err := authDomain.CreateAccount(res)

	if err != nil {
		self.renderError(err)
		return nil
	}

	account := me.ToAccount()

	self.renderJson(accountAuthorizedResponse{
		AccessToken: me.Token,
		Account:     account,
	})
	return nil
}

func (self *Handler) renewAccessToken() error {
	var req AccountAuthorizeRequest
	self.bodyAsJson(&req)

	authDomain := self.provider.GetAuthDomain()
	me, err := authDomain.SignInAccount(req)

	if err != nil {
		self.renderError(err)
		return nil
	}

	account := me.ToAccount()

	self.renderJson(accountAuthorizedResponse{
		AccessToken: me.Token,
		Account:     account,
	})
	return nil
}
