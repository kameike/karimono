package handler

import (
	"encoding/json"

	"github.com/kameike/karimono/domain"
	"github.com/kameike/karimono/util"
	"github.com/labstack/echo"
	_ "github.com/mattn/go-sqlite3"
)

func SignUp(c echo.Context) error {
	return createHandler(c).createAccount()
}

func SignIn(c echo.Context) error {
	return createHandler(c).renewAccessToken()
}

func ValidateAccount(c echo.Context) error {
	return createHandler(c).validateAccount()
}

func UpdateAccount(c echo.Context) error {
	return createHandler(c).updateAccount()
}

func CreateTeam(c echo.Context) error {
	return handleWithAccountDomain(c, createTeam)
}

func JoinTeam(c echo.Context) error {
	return handleWithAccountDomain(c, joinTeam)
}

func GetAccountBorrowing(c echo.Context) error {
	return handleWithAccountDomain(c, getBorrowings)
}

func GetAccountHistories(c echo.Context) error {
	return handleWithAccountDomain(c, getHistory)
}

func GetTeams(c echo.Context) error {
	return handleWithAccountDomain(c, getTeams)
}

type accountDomainHandler func(domain.AccountDomain, *Handler)

func handleWithAccountDomain(c echo.Context, handler accountDomainHandler) error {
	h := createHandler(c)

	a, err := h.provider.GetAccountDomain()
	if err != nil {
		h.renderError(err)
		return nil
	}

	handler(a, h)

	return nil
}

type errorResponse struct {
	Message string `json:"message"`
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
		context:  c,
	}

	return &handler
}

func (self *Handler) bodyAsJson(target interface{}) {
	c := self.context
	d := json.NewDecoder(c.Request().Body)
	err := d.Decode(target)
	util.CheckInternalFatalError(err)
}

func (self *Handler) renderError(err error) {
	c := self.context
	c.String(400, "damedame")
}

func (self *Handler) renderJson(target interface{}) {
	c := self.context
	c.JSON(200, target)
}
