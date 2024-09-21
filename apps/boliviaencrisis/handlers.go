package boliviaencrisis

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getUSTDPrices(repo BoliviaCrisisRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		prices, err := repo.GetAllUSTDPrices()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, prices)
	}
}

func getMainPageData(repo BoliviaCrisisRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		prices, err := repo.GetLastUSTDPrices()
		if err != nil {
			return err
		}

		var lastPrice int64 = 0
		if len(prices) > 0 {
			lastPrice = prices[0].Price
		}

		return c.JSON(http.StatusOK, MainPageData{
			USTDPrice:       lastPrice,
			LastUSTDRecords: prices,
		})
	}
}
