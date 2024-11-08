package econewscron

import (
	"github.com/ikanadev/ikanaapi/config"
	"github.com/jmoiron/sqlx"
)

func SetupEconomicNewsCron(db *sqlx.DB, config config.Config) {
	// handleVision360News(db, config)
	// handleLaPrensaNews(db, config)
	handleElDiaNews(db, config)
}
