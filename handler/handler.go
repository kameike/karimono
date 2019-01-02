package handler

import (
	"encoding/json"

	"github.com/kameike/karimono/domain"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

func SignUp(c echo.Context) error {
	return createHandler(c).createAccount()
}

func SignIn(c echo.Context) error {
	return createHandler(c).renewAccessToken()
}

func GeneralHandler(c echo.Context) error {
	return nil
}

type TokenProvider struct {
	context echo.Context
}

func (p *TokenProvider) AccountAccessToken() string {
	token := p.context.Request().Header.Get("x-karimono-token")
	return token
}

func (p *TokenProvider) HasToken() bool {
	return p.AccountAccessToken() != ""
}

type Handler struct {
	provider domain.DomainsProvider
	context  echo.Context
}

func createHandler(c echo.Context) *Handler {
	tokenProvider := TokenProvider{
		context: c,
	}

	handler := Handler{
		provider: domain.CreateApplicatoinDomains(&tokenProvider),
	}

	return &handler
}

func (self *Handler) bodyAsJson(target interface{}) {
	c := self.context
	d := json.NewDecoder(c.Request().Body)
	d.Decode(target)
}

func (self *Handler) renderError(err error) {
	c := self.context
	c.String(400, "damedame")
}

func (self *Handler) renderJson(target interface{}) {
	c := self.context
	c.JSON(200, target)
}
