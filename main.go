package main

import (
	"log"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/ikanadev/ikanaapi/apps/boliviaencrisis"
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
	config := GetConfig()
	app := echo.New()
	db := sqlx.MustConnect("postgres", config.DBConn)
	defer db.Close()
	migrateDB(config)

	app.Use(middleware.CORS())
	boliviaencrisis.SetupServer(app, db)
	app.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Running with cron :)")
	})
	setupPriceCron(db)
	panicIfErr(app.Start("0.0.0.0:" + config.Port))
}

func migrateDB(config Config) {
	migrator, err := migrate.New(config.MigrationsSource, config.DBConn)
	panicIfErr(err)
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		panicIfErr(err)
	}
}
