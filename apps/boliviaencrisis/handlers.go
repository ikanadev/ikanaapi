package boliviaencrisis

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getUSTDPrices(repo BoliviaCrisisRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		prices, err := repo.GetUSTDPrices()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, prices)
	}
}
