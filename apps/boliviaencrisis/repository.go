package boliviaencrisis

import (
	"github.com/jmoiron/sqlx"
)

type BoliviaCrisisRepository interface {
	GetAllUSDTPrices() ([]USDTPrice, error)
	GetLastUSDTPrices() ([]USDTPrice, error)
}

type BoliviaCrisisRepositoryImpl struct {
	db *sqlx.DB
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

func newBoliviaCrisisRepositoryImpl(db *sqlx.DB) BoliviaCrisisRepositoryImpl {
	return BoliviaCrisisRepositoryImpl{db}
}
