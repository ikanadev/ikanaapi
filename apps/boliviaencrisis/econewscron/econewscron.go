package econewscron

import (
	"sync"
	"time"

	"github.com/ikanadev/ikanaapi/config"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
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
	c := cron.New(cron.WithLocation(time.UTC))
	cronExp := "0 14,16,18,20,24 * * *" // each day at 10,12,14,16 & 18 hours
	_, err := c.AddFunc(cronExp, func() {
		ecoNewsCron.fetchNews()
	})
	if err != nil {
		panic(err)
	}
	c.Start()
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
