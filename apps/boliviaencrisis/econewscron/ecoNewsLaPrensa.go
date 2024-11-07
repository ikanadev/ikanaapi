package econewscron

import (
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"github.com/ikanadev/ikanaapi/config"
	"github.com/jmoiron/sqlx"
)

func handleLaPrensaNews(db *sqlx.DB, config config.Config) {
	ecoNews := getLaPrensaNews()
	ecoNews = filterUnparsedNews(ecoNews, db)

	var wg sync.WaitGroup
	wg.Add(len(ecoNews))
	for i := range ecoNews {
		go func(index int) {
			getLaPrensaNewDetails(ecoNews[index])
			generateAIEcoNewData(ecoNews[index], config.OpenAIKey)
			wg.Done()
		}(i)
	}
	wg.Wait()
	saveEcoNews(db, ecoNews)
}

func getLaPrensaNewDetails(ecoNew *EconomicNew) {
	c := colly.NewCollector()
	c.OnHTML("article div.field--name-body", func(e *colly.HTMLElement) {
		ecoNew.Content = &e.Text
	})
	c.Visit(ecoNew.URL)
}

func getLaPrensaNews() []*EconomicNew {
	c := colly.NewCollector()
	news := make([]*EconomicNew, 0)
	baseURL := "https://laprensa.bo"
	c.OnHTML(".views-row", func(e *colly.HTMLElement) {
		imageUrl := baseURL + e.ChildAttr("img.image-field", "src")
		title := e.ChildText("h2 a")
		url := baseURL + e.ChildAttr("a", "href")
		dateStr := e.ChildAttr("time", "datetime")
		date := time.Now().UTC()
		if len(dateStr) > 10 {
			dateStr = dateStr[:10]
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				date = parsedDate.UTC()
			}
		}

		ecoNew := EconomicNew{
			ID:      uuid.New(),
			Company: "La Prensa",
			Image:   &imageUrl,
			Title:   title,
			URL:     url,
			Date:    &date,
			Tags:    make([]string, 0),
		}
		news = append(news, &ecoNew)
	})

	c.Visit(baseURL + "/economia")

	return news
}
