package handler

import (
	"database/sql"
	"encoding/json"

	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
	"github.com/labstack/echo"
	sqlite3 "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func openDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./db/main.db")
	util.CheckInternalFatalError(err)
	return db
}

type handerWithAccount func(user model.Account, c echo.Context) error

func updateAccount(user model.Account, c echo.Context) error {

	var reqBody model.AccountCreateRequest
	d := json.NewDecoder(c.Request().Body)
	err := d.Decode(&reqBody)
	util.CheckInternalFatalError(err)

	return nil
}

var UpdateAccount = checkAuth(updateAccount)

func checkAuth(handler handerWithAccount) echo.HandlerFunc {
	return func(c echo.Context) error {
		db := openDb()
		defer db.Close()

		token := c.Request().Header.Get("x-karimono-token")

		smit, err := db.Prepare(`
select account.name, account.id from access_token join account on access_token.account_id = account.id 
		where session_token = ?
	`)

		rows, err := smit.Query(token)
		defer rows.Close()
		util.CheckInternalFatalError(err)

		var account model.Account
		for rows.Next() {
			rows.Scan(&account.Name, &account.Id)
		}

		if account.Id == "" {
			err := model.ErrorResponse{"invalid auth"}
			c.JSON(403, err)
			return nil
		}
		return handler(account, c)
	}
}

func CreateAccountHandler(c echo.Context) error {
	d := json.NewDecoder(c.Request().Body)
	var res model.AccountCreateRequest
	d.Decode(&res)

	db := openDb()
	defer db.Close()

	pass, err := bcrypt.GenerateFromPassword([]byte(res.Password), bcrypt.DefaultCost)
	util.CheckInternalFatalError(err)

	tx, err := db.Begin()
	util.CheckInternalFatalError(err)

	smit, err := db.Prepare("insert into account(name, password_hash) values(?,?)")
	_, err = smit.Exec(res.Name, string(pass))
	if serr, ok := err.(sqlite3.Error); ok && serr.ExtendedCode == sqlite3.ErrConstraintUnique {
		resbody := model.ErrorResponse{"name has aliready taken"}
		c.JSON(400, resbody)
		return nil
	}
	util.CheckInternalFatalError(err)

	token, err := renewToken(res.Name, db)
	util.CheckInternalFatalError(err)

	tx.Commit()

	result := model.AccountCreateResponse{
		AccessToken: token,
	}

	c.JSON(200, result)
	return nil
}

func renewToken(name string, db *sql.DB) (string, error) {
	token := util.RandString(100)

	query := `
	insert or replace into access_token (account_id, session_token)
	select id, ? from account where name = ? 
	`
	smit, err := db.Prepare(query)
	util.CheckInternalFatalError(err)
	_, err = smit.Exec(token, name)
	util.CheckInternalFatalError(err)

	return token, nil
}

func RenewAccessTokenHandler(c echo.Context) error {
	d := json.NewDecoder(c.Request().Body)
	var requestBody model.AccountCreateRequest
	err := d.Decode(&requestBody)
	util.CheckInternalFatalError(err)

	db := openDb()
	defer db.Close()

	smit, err := db.Prepare(`
		select password_hash, name, id from account where name = ?
	`)

	rows, err := smit.Query(requestBody.Name)
	defer rows.Close()

	util.CheckInternalFatalError(err)

	var account model.Account
	var passHash []byte
	for rows.Next() {
		rows.Scan(&passHash, &account.Name, &account.Id)
	}

	err = bcrypt.CompareHashAndPassword(passHash, []byte(requestBody.Password))
	if err != nil {
		msg := model.ErrorResponse{"id or pass is wrong"}
		c.JSON(400, msg)
		return nil
	}
	util.CheckInternalFatalError(err)

	token, err := renewToken(account.Name, db)
	result := model.AccountCreateResponse{
		AccessToken: token,
	}

	c.JSON(200, result)

	return nil
}
