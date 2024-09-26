package common

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func parseForwarededFor(c echo.Context) []string {
	forwardedForStr := c.Request().Header.Get("X-Forwarded-For")
	if forwardedForStr == "" {
		return make([]string, 0)
	}
	forwardedParts := strings.Split(forwardedForStr, ",")
	for i := range forwardedParts {
		forwardedParts[i] = strings.TrimSpace(forwardedParts[i])
	}
	return forwardedParts
}

func postPageViewRecord(repo CommonRepository) echo.HandlerFunc {
	type reqData struct {
		App    string `json:"app"`
		UserID string `json:"userId"`
		URL    string `json:"url"`
	}
	return func(c echo.Context) error {
		var reqData reqData
		if err := c.Bind(&reqData); err != nil {
			return err
		}
		ips := parseForwarededFor(c)
		data := PageViewRecordData{
			App:    reqData.App,
			UserID: reqData.UserID,
			URL:    reqData.URL,
			Ips:    ips,
		}

		if err := repo.SavePageViewRecord(data); err != nil {
			return err
		}
		return c.NoContent(http.StatusCreated)
	}
}
