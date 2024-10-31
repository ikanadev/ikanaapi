package econewscron

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

func filterUnparsedNews(news []*EconomicNew, db *sqlx.DB) []*EconomicNew {
	urls := make([]string, len(news))
	for i := range news {
		urls[i] = news[i].URL
	}
	sql := `
	SELECT ARRAY(
		SELECT EXISTS (
			SELECT 1 FROM economic_new WHERE url = u.url
		) FROM unnest($1::text[]) as u(url)
	);`

	exists := make([]bool, len(news))
	err := db.Get(pq.Array(&exists), sql, pq.Array(urls))
	if err != nil {
		log.Printf("Error quering existing news: %v", err)
		return make([]*EconomicNew, 0)
	}

	filteredNews := make([]*EconomicNew, 0)
	for i := range news {
		if !exists[i] {
			filteredNews = append(filteredNews, news[i])
		}
	}

	return filteredNews
}

func saveEcoNews(db *sqlx.DB, news []*EconomicNew) {
	if len(news) == 0 {
		return
	}
	dbNews := make([]*DbEconomicNew, len(news))
	for i := range news {
		dbNews[i] = news[i].ToDbEconomicNew()
		dbNews[i].CreatedAt = time.Now().UTC()
	}
	sql := `
		INSERT INTO economic_new (id, title, date, url, image, summary, company, tags, sentiment, created_at, deleted_at)
		VALUES (:id, :title, :date, :url, :image, :summary, :company, :tags, :sentiment, :created_at, :deleted_at);
	`
	_, err := db.NamedExec(sql, dbNews)
	if err != nil {
		log.Printf("Error saving economic news: %v", err)
		return
	}
}

func generateAIEcoNewData(ecoNew *EconomicNew) {
	ecoNew.Summary = nil
	ecoNew.Sentiment = 0
	ecoNew.Tags = make([]string, 0)
}
