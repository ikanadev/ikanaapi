package boliviaencrisis

import (
	"github.com/google/uuid"
	"github.com/ikanadev/ikanaapi/common"
)

type USTDPrice struct {
	ID    uuid.UUID `json:"id"`
	Price int64     `json:"price"`
	common.TimeData
}
