package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/ikanadev/ikanaapi/apps/boliviaencrisis"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	config := GetConfig()
	app := echo.New()
	db := sqlx.MustConnect("postgres", config.DBConn)
	migrateDB(config)

	app.Use(middleware.CORS())
	boliviaencrisis.SetupServer(app, db)
	panicIfErr(app.Start(config.Port))
}

func migrateDB(config Config) {
	migrator, err := migrate.New(config.MigrationsURL, config.DBConn)
	panicIfErr(err)
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		panicIfErr(err)
	}
}
