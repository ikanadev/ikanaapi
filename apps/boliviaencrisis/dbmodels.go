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
