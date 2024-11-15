package econewscron

import (
	"sync"

	"github.com/ikanadev/ikanaapi/config"
	"github.com/jmoiron/sqlx"
)

type EcoNewsCron struct {
	sources []EcoNewsSource
	db      *sqlx.DB
	config  config.Config
}

func NewEcoNewsCron(db *sqlx.DB, config config.Config) *EcoNewsCron {
	sources := []EcoNewsSource{
		newVision360Source(),
		newElDiaSource(),
		newLaPrensaSource(),
		newElDeberSource(),
		newCorreoDelSurSource(),
	}
	return &EcoNewsCron{
		sources: sources,
		db:      db,
		config:  config,
	}
}

func (ecoNewsCron *EcoNewsCron) SetupCron() {
	ecoNewsCron.fetchNews()
}

func (ecoNewsCron *EcoNewsCron) fetchNews() {
	var wgSources sync.WaitGroup
	wgSources.Add(len(ecoNewsCron.sources))

	for _, s := range ecoNewsCron.sources {
		go func(source EcoNewsSource) {
			defer wgSources.Done()
			ecoNews := source.GetEcoNews()
			ecoNews = filterUnparsedNews(ecoNews, ecoNewsCron.db)

			var wgNews sync.WaitGroup
			wgNews.Add(len(ecoNews))

			for i := range ecoNews {
				go func(index int) {
					source.GetEcoNewDetails(ecoNews[index])
					generateAIEcoNewData(ecoNews[index], ecoNewsCron.config.OpenAIKey)
					wgNews.Done()
				}(i)
			}
			wgNews.Wait()
			saveEcoNews(ecoNewsCron.db, ecoNews)
		}(s)
	}

	wgSources.Wait()
}
