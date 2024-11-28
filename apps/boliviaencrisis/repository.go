package boliviaencrisis

import (
	"time"

	"github.com/ikanadev/ikanaapi/apps/boliviaencrisis/econewscron"
	"github.com/jmoiron/sqlx"
)

type BoliviaCrisisRepository interface {
	GetAllUSDTPrices() ([]USDTPrice, error)
	GetLastUSDTPrices() ([]USDTPrice, error)
	GetUSDTPriceByDate(date time.Time) (int64, error)
	GetLatestEcoNews() ([]*econewscron.EconomicNew, error)
}

type BoliviaCrisisRepositoryImpl struct {
	db *sqlx.DB
}

// GetUSDTPriceByDate implements BoliviaCrisisRepository.
func (r BoliviaCrisisRepositoryImpl) GetUSDTPriceByDate(date time.Time) (int64, error) {
	query := `
	SELECT COALESCE(ROUND(AVG(price)), 0) as average
	from ustd_price
	where DATE(created_at at time zone 'UTC-4') = $1;`

	var average int64
	err := r.db.Get(&average, query, date.Format("2006-01-02"))

	return average, err
}

// GetLastUSDTPrices implements BoliviaCrisisRepository.
func (r BoliviaCrisisRepositoryImpl) GetLastUSDTPrices() ([]USDTPrice, error) {
	var dbPrices []DbUSDTPrice
	err := r.db.Select(&dbPrices, "SELECT * FROM ustd_price order by created_at desc limit 7;")
	if err != nil {
		return nil, err
	}

	prices := make([]USDTPrice, len(dbPrices))
	for i, dbPrice := range dbPrices {
		prices[len(prices)-i-1] = dbPrice.ToUSDTPrice()
	}

	return prices, err
}

// GetAllUSDTPrices implements BoliviaCrisisRepository.
func (r BoliviaCrisisRepositoryImpl) GetAllUSDTPrices() ([]USDTPrice, error) {
	var dbPrices []DbUSDTPrice
	err := r.db.Select(&dbPrices, "SELECT * FROM ustd_price order by created_at desc;")
	if err != nil {
		return nil, err
	}

	prices := make([]USDTPrice, len(dbPrices))
	for i, dbPrice := range dbPrices {
		prices[i] = dbPrice.ToUSDTPrice()
	}

	return prices, nil
}

// GetLatestEcoNews implements BoliviaCrisisRepository.
func (r BoliviaCrisisRepositoryImpl) GetLatestEcoNews() ([]*econewscron.EconomicNew, error) {
	var dbNews []econewscron.DbEconomicNew
	query := `
	SELECT * FROM
	economic_new
	WHERE summary IS NOT NULL
	AND trim(summary) != ''
	AND tags IS NOT NULL
	AND array_length(tags, 1) > 0
	order by created_at desc;`
	err := r.db.Select(&dbNews, query)
	if err != nil {
		return nil, err
	}

	news := make([]*econewscron.EconomicNew, len(dbNews))
	for i := range dbNews {
		news[i] = dbNews[i].ToEconomicNew()
	}

	return news, nil
}

func newBoliviaCrisisRepositoryImpl(db *sqlx.DB) BoliviaCrisisRepositoryImpl {
	return BoliviaCrisisRepositoryImpl{db}
}
