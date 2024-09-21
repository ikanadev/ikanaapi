package boliviaencrisis

import (
	"github.com/google/uuid"
	"github.com/ikanadev/ikanaapi/common"
)

type DbUSTDPrice struct {
	ID    uuid.UUID `db:"id"`
	Price int64     `db:"price"`
	common.DBTimeData
}

func (dbPrice DbUSTDPrice) ToUSTDPrice() USTDPrice {
	return USTDPrice{
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
