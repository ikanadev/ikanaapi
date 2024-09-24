package boliviaencrisis

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func getUSDTPrices(repo BoliviaCrisisRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		prices, err := repo.GetAllUSDTPrices()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, prices)
	}
}

func getMainPageData(repo BoliviaCrisisRepository) echo.HandlerFunc {
	return func(c echo.Context) error {
		prices, err := repo.GetLastUSDTPrices()
		if err != nil {
			return err
		}

		var lastPrice int64 = 0
		if len(prices) > 0 {
			lastPrice = prices[len(prices)-1].Price
		}

		return c.JSON(http.StatusOK, MainPageData{
			USDTPrice:       lastPrice,
			LastUSDTRecords: prices,
		})
	}
}
