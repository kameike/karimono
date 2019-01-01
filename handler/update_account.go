package handler

import (
	"database/sql"

	"github.com/kameike/karimono/model"
	"github.com/labstack/echo"
)

var UpdateAccount = injectDbConn(checkAuth(updateAccount))

func updateAccount(user model.Account, db *sql.DB, c echo.Context) error {
	return nil
}

// func updateAccount(user model.Account, db *sql.DB, c echo.Context) error {
// 	var reqBody model.AccountCreateRequest
// 	d := json.NewDecoder(c.Request().Body)
// 	err := d.Decode(&reqBody)
// 	util.CheckInternalFatalError(err)
//
// 	tx, err := db.Begin()
// 	util.CheckInternalFatalError(err)
//
// 	pass, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
// 	util.CheckInternalFatalError(err)
//
// 	smit, err := db.Prepare(`
// 		update set password_hash = ?, name = ? where name = ?
// 	`)
// 	util.CheckInternalFatalError(err)
//
// 	_, err = smit.Exec(pass, reqBody.Name, user.Id)
// 	util.CheckInternalFatalError(err)
//
// 	token, err := renewToken(reqBody.Name, db)
// 	util.CheckInternalFatalError(err)
//
// 	err = tx.Commit()
// 	util.CheckInternalFatalError(err)
//
// 	model := model.AccountCreateResponse{
// 		AccessToken: token,
// 	}
//
// 	c.JSON(200, model)
// 	return nil
// }
