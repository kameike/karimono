package main

import (
	"database/sql"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"

	"github.com/kameike/karimono/model"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	sqlite3 "github.com/mattn/go-sqlite3"
)

func main() {
	os.Getenv("KARIMONO_PATH")
	err := os.MkdirAll("tmp", 0770)
	if err != nil {
		println(err.Error())
		panic("failt to make tmp")
	}
	// Echoのインスタンス作る
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/account", createAccountHandler)  //アカウント作成
	e.POST("/token", renewAccessTokenHandler) //アカウント作成
	e.PUT("/account", checkAuth)              //アカウント情報のアップデート

	e.POST("/teams/id/menbers", stub())   // チームにサインインする
	e.GET("/teams/id/menbers", stub())    // チームメンバーを一覧する
	e.DELETE("/teams/id/menbers", stub()) // チームから抜ける

	e.GET("/teams", stub())                // 参加しているチームの情報を見る
	e.POST("/teams/id", stub())            // チームを作成する
	e.PUT("/teams/id", stub())             // チーム情報のアップデート
	e.GET("/teams/id", checkAuth)          // チーム情報を取得する
	e.GET("/teams/histories", stub())      // チームで起きたことの履歴を見る
	e.GET("/teams/id/borrowings", stub())  // チームでアイテムを借りる
	e.POST("/teams/id/borrowings", stub()) // チームでアイテムを借りる

	e.GET("/borrowings", stub())      // 自分が借りているものを一覧する
	e.POST("/returning/uuid", stub()) // アイテムを返す

	e.Start(":1323") //ポート番号指定してね
}

type HanderWithAcoount func(user model.Account, c echo.Context) error

func authenticateUser() {
}

func checkAuth(c echo.Context) error {
	db, err := sql.Open("sqlite3", "./db/main.db")
	defer db.Close()
	checkInternalFatalError(err)

	token := c.Request().Header.Get("x-karimono-token")

	smit, err := db.Prepare(`
		select account.name, account.id from access_token join account on access_token.account_id = account.id 
		where session_token = ?
	`)

	rows, err := smit.Query(token)
	defer rows.Close()
	checkInternalFatalError(err)

	var account model.Account
	for rows.Next() {
		rows.Scan(&account.Name, &account.Id)
	}

	if account.Id == "" {
		c.String(400, "dame")
		return nil
	}

	c.String(200, "hello "+account.Name)

	return nil
}

func stub() echo.HandlerFunc {
	return checkAuth
}

func renewAccessTokenHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "./db/main.db")
	defer db.Close()
	checkInternalFatalError(err)

	d := json.NewDecoder(c.Request().Body)
	var requestBody model.AccountCreateRequest
	err = d.Decode(&requestBody)
	checkInternalFatalError(err)

	smit, err := db.Prepare(`
		select password_hash, name, id from account where name = ?
	`)

	rows, err := smit.Query(requestBody.Name)
	defer rows.Close()

	checkInternalFatalError(err)

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
	checkInternalFatalError(err)

	token, err := renewToken(account.Name, db)
	result := model.AccountCreateResponse{
		AccessToken: token,
	}

	c.JSON(200, result)

	return nil
}

func createAccountHandler(c echo.Context) error {
	db, err := sql.Open("sqlite3", "./db/main.db")
	defer db.Close()
	checkInternalFatalError(err)

	d := json.NewDecoder(c.Request().Body)
	var res model.AccountCreateRequest
	d.Decode(&res)

	pass, err := bcrypt.GenerateFromPassword([]byte(res.Password), bcrypt.DefaultCost)
	checkInternalFatalError(err)

	tx, err := db.Begin()
	checkInternalFatalError(err)

	smit, err := db.Prepare("insert into account(name, password_hash) values(?,?)")
	_, err = smit.Exec(res.Name, string(pass))
	if serr, ok := err.(sqlite3.Error); ok && serr.ExtendedCode == sqlite3.ErrConstraintUnique {
		resbody := model.ErrorResponse{"name has aliready taken"}
		c.JSON(400, resbody)
		return nil
	}

	checkInternalFatalError(err)

	token, err := renewToken(res.Name, db)
	checkInternalFatalError(err)

	tx.Commit()

	result := model.AccountCreateResponse{
		AccessToken: token,
	}

	c.JSON(200, result)
	return nil
}

func renewToken(name string, db *sql.DB) (string, error) {
	token := RandStringRunes(100)
	query := `
	insert or replace into access_token (account_id, session_token)
	select id, ? from account where name = ? 
	`
	smit, err := db.Prepare(query)
	checkInternalFatalError(err)
	_, err = smit.Exec(token, name)
	checkInternalFatalError(err)

	return token, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func checkInternalFatalError(err error) {
	if err != nil {
		println("======")
		println(err.Error())
		println("======")
		panic(err.Error())
	}
}
