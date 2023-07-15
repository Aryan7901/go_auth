package main

import (
	"auth/helpers/constants"
	auth "auth/routes/auth"
	"fmt"

	"embed"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//go:embed templates/*
var templates embed.FS

func main() {
	app := echo.New()
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	var err, env_err error
	if env_err != nil {
		fmt.Println(env_err)
	}
	db, err := sqlx.Open("pgx", constants.DbConn)
	config := constants.AppConfig{Db: db, Templates: templates}
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	auth.Register(app, &config)
	app.Logger.Fatal(app.Start(":" + constants.Port))
}
