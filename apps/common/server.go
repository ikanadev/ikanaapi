package common

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Server struct {
	app  *echo.Group
	repo CommonRepository
}

func (s Server) setupRoutes() {
	s.app.POST("/page_view", postPageViewRecord(s.repo))
	s.app.POST("/public_feedback", postPublicFeedback(s.repo))
}

func newServer(app *echo.Group, repo CommonRepository) Server {
	return Server{app, repo}
}

func SetupServer(app *echo.Echo, db *sqlx.DB) {
	appGroup := app.Group("common")
	repo := newCommonRepositoryImpl(db)
	server := newServer(appGroup, repo)
	server.setupRoutes()
}
