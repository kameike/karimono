package handler

import "github.com/kameike/karimono/model"

func (self *Handler) updateAccount() error {
	return nil
}

type accountAuthorizeRequest struct {
	name     string
	password string
}

type accountAuthorizedResponse struct {
	AccessToken string `json:"accessToken"`
	Account     model.Account
}

func (self accountAuthorizeRequest) AccountId() string {
	return self.name
}

func (self accountAuthorizeRequest) AccountPassword() string {
	return self.password
}

func (self *Handler) createAccount() error {
	var res accountAuthorizeRequest
	self.bodyAsJson(&res)

	authDomain := self.provider.GetAuthDomain()
	me, err := authDomain.CreateAccount(res)

	if err == nil {
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
	var req accountAuthorizeRequest
	self.bodyAsJson(&req)

	authDomain := self.provider.GetAuthDomain()
	me, err := authDomain.SignInAccount(req)

	if err == nil {
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
