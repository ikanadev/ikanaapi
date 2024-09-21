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

type MainPageData struct {
	USTDPrice       int64       `json:"ustdPrice"`
	LastUSTDRecords []USTDPrice `json:"lastUSTDRecords"`
}
