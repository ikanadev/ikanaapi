package boliviaencrisis

import (
	"github.com/google/uuid"
	"github.com/ikanadev/ikanaapi/apps/boliviaencrisis/econewscron"
	"github.com/ikanadev/ikanaapi/common"
)

type USDTPrice struct {
	ID    uuid.UUID `json:"id"`
	Price int64     `json:"price"`
	common.TimeData
}

type MainPageData struct {
	USDTPrice          int64                      `json:"usdtPrice"`
	USDTPriceLastWeek  int64                      `json:"usdtPriceLastWeek"`
	USDTPriceLastMonth int64                      `json:"usdtPriceLastMonth"`
	LastUSDTRecords    []USDTPrice                `json:"lastUsdtRecords"`
	EcoNews            []*econewscron.EconomicNew `json:"ecoNews"`
}
