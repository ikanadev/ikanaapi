package boliviaencrisis

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Server struct {
	app  *echo.Group
	repo BoliviaCrisisRepository
}

func (s Server) setupRoutes() {
	s.app.GET("/prices/usdt", getUSDTPrices(s.repo))
	s.app.GET("/main", getMainPageData(s.repo))
}

func newServer(app *echo.Group, repo BoliviaCrisisRepository) Server {
	return Server{app, repo}
}

func SetupServer(app *echo.Echo, db *sqlx.DB) {
	appGroup := app.Group("crisis")
	repo := newBoliviaCrisisRepositoryImpl(db)
	server := newServer(appGroup, repo)
	server.setupRoutes()
}
