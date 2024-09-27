package boliviaencrisis

import (
	"net/http"
	"time"

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

		// Bolivia UTC time
		location, err := time.LoadLocation("America/La_Paz")
		if err != nil {
			return err
		}

		now := time.Now().In(location)
		lastWeekPrice, err := repo.GetUSDTPriceByDate(now.AddDate(0, 0, -7))
		if err != nil {
			return err
		}
		lastMonthPrice, err := repo.GetUSDTPriceByDate(now.AddDate(0, -1, 0))
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, MainPageData{
			USDTPrice:          lastPrice,
			USDTPriceLastWeek:  lastWeekPrice,
			USDTPriceLastMonth: lastMonthPrice,
			LastUSDTRecords:    prices,
		})
	}
}
