package main

import (
	"github.com/kameike/karimono/handler"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

func main() {
	os.Getenv("KARIMONO_PATH")
	os.Mkdir("tmp", 0666)
	// Echoのインスタンス作る
	e := echo.New()

	// 全てのリクエストで差し込みたいミドルウェア（ログとか）はここ
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/item/borrow", handler.Borrow())
	e.POST("/item/return", handler.Return())

	e.GET("/items", handler.Items())

	// サーバー起動
	e.Start(":1323") //ポート番号指定してね
}
