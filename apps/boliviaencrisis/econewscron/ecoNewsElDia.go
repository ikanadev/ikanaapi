package econewscron

import (
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
	"github.com/ikanadev/ikanaapi/config"
	"github.com/jmoiron/sqlx"
)

func handleElDiaNews(db *sqlx.DB, config config.Config) {
	ecoNews := getElDiaNews()
	ecoNews = filterUnparsedNews(ecoNews, db)

	var wg sync.WaitGroup
	wg.Add(len(ecoNews))
	for i := range ecoNews {
		go func(index int) {
			getElDiaNewDetails(ecoNews[index])
			generateAIEcoNewData(ecoNews[index], config.OpenAIKey)
			wg.Done()
		}(i)
	}
	wg.Wait()

	saveEcoNews(db, ecoNews)
}

func getElDiaNewDetails(ecoNew *EconomicNew) {
	c := colly.NewCollector()
	c.OnHTML("div.info div.content", func(e *colly.HTMLElement) {
		ecoNew.Content = &e.Text
	})
	c.Visit(ecoNew.URL)
}

func getElDiaNews() []*EconomicNew {
	c := colly.NewCollector()

	news := make([]*EconomicNew, 0)
	baseUrl := "https://www.eldia.com.bo"

	c.OnHTML("div.pub-content div.postBox", func(e *colly.HTMLElement) {
		imgUrl := e.ChildAttr("img", "src")
		url := e.ChildAttr("a.title", "href")
		title := e.ChildText("h2")
		date := time.Now().UTC()
		if len(url) > 35 {
			dateStr := url[25:35]
			parsedDate, err := time.Parse("2006-01-02", dateStr)
			if err == nil {
				date = parsedDate.UTC()
			}
		}

		ecoNew := EconomicNew{
			ID:      uuid.New(),
			Company: "El Día",
			Image:   &imgUrl,
			Title:   title,
			URL:     url,
			Date:    &date,
			Tags:    make([]string, 0),
		}
		news = append(news, &ecoNew)
	})

	c.Visit(baseUrl + "/economia")
	return news
}
