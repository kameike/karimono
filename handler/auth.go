package handler

func (self *Handler) updateAccount() error {
	return nil
}

type accountCreateRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (self accountCreateRequest) AccountId() string {
	return self.Name
}

func (self accountCreateRequest) AccountPassword() string {
	return self.Password
}

type accountCreateResponse struct {
	AccessToken string `json:"accessToken"`
}

func (self *Handler) createAccount() error {
	var res accountCreateRequest
	self.bodyAsJson(&res)

	authDomain := self.provider.GetAuthDomain()
	account, err := authDomain.CreateAccount(res)

	if err == nil {
		self.renderError(err)
		return nil
	}

	self.renderJson(account)
	return nil
}

func (self *Handler) renewAccessToken() error {
	return nil
}
