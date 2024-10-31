package econewscron

import "github.com/jmoiron/sqlx"

func SetupEconomicNewsCron(db *sqlx.DB) {
	vision360News := getVision360News()
	vision360News = filterUnparsedNews(vision360News, db)
	for i := range vision360News {
		getVision360NewDetails(vision360News[i])
		generateAIEcoNewData(vision360News[i])
	}
	saveEcoNews(db, vision360News)
}
