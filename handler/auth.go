package handler

import (
	"encoding/json"

	"github.com/kameike/karimono/domain"
	"github.com/kameike/karimono/model"
	"github.com/labstack/echo"
)

var CreateAccountHandler = createHandler().createAccount

var RenewAccessTokenHandler = createHandler().renewAccessToken

type Handler struct {
	provider domain.DomainsProvider
}

func createHandler() *Handler {
	handler := Handler{
		provider: domain.CreateApplicatoinDomains(),
	}

	return &handler
}

type Contents struct {
	code    int
	content []byte
}

func bodyAsJson(target interface{}, c echo.Context) {
	d := json.NewDecoder(c.Request().Body)
	d.Decode(target)
}

func (self *Handler) createAccount(c echo.Context) error {
	var res model.AccountCreateRequest
	bodyAsJson(&res, c)

	authDomain := self.provider.GetAuthDomain()
	account, err := authDomain.CreateAccount(res)

	if err == nil {
		renderError(err, c)
		return nil
	}

	c.JSON(200, account)

	return nil
}

func (self *Handler) renewAccessToken(c echo.Context) error {
	return nil
}

func renderError(err error, c echo.Context) {
	c.String(400, "damedame")
}

// func CreateAccountHandler(c echo.Context) error {
//
// 	db := openDb()
// 	defer db.Close()
//
// 	pass, err := bcrypt.GenerateFromPassword([]byte(res.Password), bcrypt.DefaultCost)
// 	util.CheckInternalFatalError(err)
//
// 	tx, err := db.Begin()
// 	util.CheckInternalFatalError(err)
//
// 	smit, err := db.Prepare("insert into account(name, password_hash) values(?,?)")
// 	_, err = smit.Exec(res.Name, string(pass))
// 	if serr, ok := err.(sqlite3.Error); ok && serr.ExtendedCode == sqlite3.ErrConstraintUnique {
// 		resbody := model.ErrorResponse{"name has aliready taken"}
// 		c.JSON(400, resbody)
// 		return nil
// 	}
// 	util.CheckInternalFatalError(err)
//
// 	token, err := renewToken(res.Name, db)
// 	util.CheckInternalFatalError(err)
//
// 	tx.Commit()
//
// 	result := model.AccountCreateResponse{
// 		AccessToken: token,
// 	}
//
// 	c.JSON(200, result)
// 	return nil
// }

// func RenewAccessTokenHandler(c echo.Context) error {
// 	d := json.NewDecoder(c.Request().Body)
// 	var requestBody model.AccountCreateRequest
// 	err := d.Decode(&requestBody)
// 	util.CheckInternalFatalError(err)
//
// 	db := openDb()
// 	defer db.Close()
//
// 	smit, err := db.Prepare(`
// 		select password_hash, name, id from account where name = ?
// 	`)
//
// 	rows, err := smit.Query(requestBody.Name)
// 	defer rows.Close()
//
// 	util.CheckInternalFatalError(err)
//
// 	var account model.Account
// 	var passHash []byte
// 	for rows.Next() {
// 		name := ""
// 		rows.Scan(&passHash, &name, &account.Id)
// 	}
//
// 	err = bcrypt.CompareHashAndPassword(passHash, []byte(requestBody.Password))
// 	if err != nil {
// 		msg := model.ErrorResponse{"id or pass is wrong"}
// 		c.JSON(400, msg)
// 		return nil
// 	}
// 	util.CheckInternalFatalError(err)
//
// 	token, err := renewToken(account.Id, db)
// 	result := model.AccountCreateResponse{
// 		AccessToken: token,
// 	}
//
// 	c.JSON(200, result)
//
// 	return nil
// }
