package handler

import (
	"database/sql"
	"github.com/kameike/karimono/model"
	"github.com/kameike/karimono/util"
	"github.com/labstack/echo"
)

type handerWithAccount func(user model.Account, db *sql.DB, c echo.Context) error

type withDbConnInjected func(conn *sql.DB, c echo.Context) error

func injectDbConn(w withDbConnInjected) echo.HandlerFunc {
	db := openDb()
	return func(c echo.Context) error {
		defer db.Close()
		return w(db, c)
	}
}

func checkAuth(handler handerWithAccount) withDbConnInjected {
	return func(db *sql.DB, c echo.Context) error {
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
			rows.Scan(&account.Id, &account.Id)
		}

		if account.Id == "" {
			err := model.ErrorResponse{"invalid auth"}
			c.JSON(403, err)
			return nil
		}
		return handler(account, db, c)
	}
}
