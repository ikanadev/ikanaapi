package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func postPageViewRecord(repo CommonRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		var reqData PageViewRecordData
		if err := c.Bind(&reqData); err != nil {
			return err
		}
		if err := repo.SavePageViewRecord(reqData); err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	}
}
