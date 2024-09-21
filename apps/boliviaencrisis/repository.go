package boliviaencrisis

import (
	"github.com/jmoiron/sqlx"
)

type BoliviaCrisisRepository interface {
	GetAllUSTDPrices() ([]USTDPrice, error)
	GetLastUSTDPrices() ([]USTDPrice, error)
}

type BoliviaCrisisRepositoryImpl struct {
	db *sqlx.DB
}

// GetLastUSTDPrices implements BoliviaCrisisRepository.
func (r BoliviaCrisisRepositoryImpl) GetLastUSTDPrices() ([]USTDPrice, error) {
	var dbPrices []DbUSTDPrice
	err := r.db.Select(&dbPrices, "SELECT * FROM ustd_price order by created_at desc limit 15;")
	if err != nil {
		return nil, err
	}

	prices := make([]USTDPrice, len(dbPrices))
	for i, dbPrice := range dbPrices {
		prices[i] = dbPrice.ToUSTDPrice()
	}

	return prices, err
}

// GetAllUSTDPrices implements BoliviaCrisisRepository.
func (r BoliviaCrisisRepositoryImpl) GetAllUSTDPrices() ([]USTDPrice, error) {
	var dbPrices []DbUSTDPrice
	err := r.db.Select(&dbPrices, "SELECT * FROM ustd_price order by created_at desc;")
	if err != nil {
		return nil, err
	}

	prices := make([]USTDPrice, len(dbPrices))
	for i, dbPrice := range dbPrices {
		prices[i] = dbPrice.ToUSTDPrice()
	}

	return prices, nil
}

func newBoliviaCrisisRepositoryImpl(db *sqlx.DB) BoliviaCrisisRepositoryImpl {
	return BoliviaCrisisRepositoryImpl{db}
}
