package main

import (
	"os"

	"github.com/kameike/karimono/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

	e.POST("/account", handler.SignUp)       //アカウント作成
	e.POST("/token", handler.SignIn)         //アカウント作成
	e.PUT("/account", handler.UpdateAccount) //アカウント情報のアップデート

	e.POST("/teams/id/menbers", stub)   // チームにサインインする
	e.GET("/teams/id/menbers", stub)    // チームメンバーを一覧する
	e.DELETE("/teams/id/menbers", stub) // チームから抜ける

	e.GET("/teams", handler.GetTeams)    // 参加しているチームの情報を見る
	e.POST("/teams", stub)               // チームを作成する
	e.PUT("/teams/id", stub)             // チーム情報のアップデート
	e.GET("/teams/id", stub)             // チーム情報を取得する
	e.GET("/teams/id/histories", stub)   // チームで起きたことの履歴を見る
	e.GET("/teams/id/borrowings", stub)  // チームでアイテムを借りる
	e.POST("/teams/id/borrowings", stub) // チームでアイテムを借りる

	e.GET("/borrowings", stub)      // 自分が借りているものを一覧する
	e.POST("/returning/uuid", stub) // アイテムを返す

	e.Start(":1323") //ポート番号指定してね
}

func stub(c echo.Context) error {
	return nil
}
