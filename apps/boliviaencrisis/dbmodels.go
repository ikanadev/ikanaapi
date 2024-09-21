package boliviaencrisis

import (
	"github.com/google/uuid"
	"github.com/ikanadev/ikanaapi/common"
)

type DbUSDTPrice struct {
	ID    uuid.UUID `db:"id"`
	Price int64     `db:"price"`
	common.DBTimeData
}

func (dbPrice DbUSDTPrice) ToUSDTPrice() USDTPrice {
	return USDTPrice{
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
