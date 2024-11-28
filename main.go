package main

import (
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ikanadev/ikanaapi/apps/boliviaencrisis"
	"github.com/ikanadev/ikanaapi/apps/boliviaencrisis/econewscron"
	"github.com/ikanadev/ikanaapi/apps/common"
	"github.com/ikanadev/ikanaapi/config"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	// postgres connection and migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func panicIfErr(err error) {
	if err != nil {
		log.Println("ERROR", err)
		panic(err)
	}
}

func main() {
	config := config.GetConfig()
	app := echo.New()
	// app.Debug = true
	db := sqlx.MustConnect("postgres", config.DBConn)
	defer db.Close()
	migrateDB(config)

	app.Use(middleware.CORS())
	app.Use(middleware.Logger())
	boliviaencrisis.SetupServer(app, db)
	common.SetupServer(app, db)

	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Running...")
	})

	setupPriceCron(db)
	newsCron := econewscron.NewEcoNewsCron(db, config)
	newsCron.SetupCron()
	panicIfErr(app.Start("0.0.0.0:" + config.Port))
}

func migrateDB(config config.Config) {
	migrator, err := migrate.New(config.MigrationsSource, config.DBConn)
	panicIfErr(err)
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		panicIfErr(err)
	}
}
