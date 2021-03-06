package main

import (
	"os"

	"github.com/kameike/karimono/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	path := os.Getenv("KARIMONO_DB_PATH")
	if path == "" {
		panic("failed to init due to evn KARIMONO_DB_PATH is not set")
	}

	// Echoのインスタンス作る
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/account", handler.SignUp)                     //アカウント作成
	e.POST("/account/validation", handler.ValidateAccount) //アカウント作成
	e.POST("/token", handler.SignIn)                       //アカウント作成
	e.PUT("/account", handler.UpdateAccount)               //アカウント情報のアップデート

	e.GET("/health", health)

	e.POST("/teams", handler.CreateTeam) // チームにサインインする
	//	e.POST("/teams/validation", handler.SignUp) //アカウント作成

	e.POST("/teams/menbers", handler.JoinTeam)             // チームに参加する
	e.GET("/teams/:id/menbers", handler.GetTeamBorrowings) // チームメンバーを一覧する
	e.DELETE("/teams/:id/menbers", stub)                   // チームから抜ける

	e.GET("/teams", handler.GetTeams) // 参加しているチームの情報を見る

	e.PUT("/teams/:id", stub)                                 // チーム情報のアップデート
	e.GET("/teams/:id", stub)                                 // チーム情報を取得する
	e.GET("/teams/:id/histories", stub)                       // チームで起きたことの履歴を見る
	e.GET("/teams/:id/borrowings", handler.GetTeamBorrowings) // チームでアイテムを借りる
	e.POST("/borrowings", handler.CreateBorrowing)            // チームでアイテムを借りる

	e.GET("/borrowings", handler.GetAccountBorrowing)        // 自分が借りているものを一覧する
	e.DELETE("/borrowings/:idHash", handler.ReturnBorrowing) // アイテムを返す
	e.DELETE("/borrowings", stub)                            // アイテムを返す

	e.GET("/histories", handler.GetAccountHistories) // 歴史

	e.Start(":1323") //ポート番号指定してね
}

func stub(c echo.Context) error {
	return nil
}

type Health struct {
	Status string
}

func health(c echo.Context) error {
	c.JSON(200, Health{"ok"})
	return nil
}
