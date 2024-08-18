package boliviaencrisis

import (
	"github.com/ikanadev/ikanaapi/common"
	"github.com/jmoiron/sqlx"
)

type BoliviaCrisisRepository interface {
	GetUSTDPrices() ([]USTDPrice, error)
}

type BoliviaCrisisRepositoryImpl struct {
	db *sqlx.DB
}

// GetUSTDPrices implements BoliviaCrisisRepository.
func (r BoliviaCrisisRepositoryImpl) GetUSTDPrices() ([]USTDPrice, error) {
	var dbPrices []DbUSTDPrice
	err := r.db.Select(&dbPrices, "SELECT * FROM ustd_price order by created_at desc")
	if err != nil {
		return nil, err
	}

	prices := make([]USTDPrice, len(dbPrices))
	for i, dbPrice := range dbPrices {
		prices[i] = USTDPrice{
			ID:    dbPrice.ID,
			Price: dbPrice.Price,
			TimeData: common.TimeData{
				CreatedAt:  dbPrice.CreatedAt,
				UpdatedAt:  dbPrice.UpdatedAt,
				ArchivedAt: dbPrice.ArchivedAt,
				DeletedAt:  dbPrice.DeletedAt,
			},
		}
	}

	return prices, nil
}

func newBoliviaCrisisRepositoryImpl(db *sqlx.DB) BoliviaCrisisRepositoryImpl {
	return BoliviaCrisisRepositoryImpl{db}
}
