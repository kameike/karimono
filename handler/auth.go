package handler

import (
	"github.com/kameike/karimono/model"
)

func (self *Handler) updateAccount() error {
	return nil
}

func (self *Handler) createAccount() error {
	var res model.AccountCreateRequest
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
