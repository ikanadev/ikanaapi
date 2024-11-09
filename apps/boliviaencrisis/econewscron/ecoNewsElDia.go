package econewscron

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

type ElDiaSource struct {
	Name string
}

func newElDiaSource() *ElDiaSource {
	return &ElDiaSource{
		Name: "El Día",
	}
}

func (elDiaSource *ElDiaSource) GetEcoNews() []*EconomicNew {
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

func (elDiaSource *ElDiaSource) GetEcoNewDetails(ecoNew *EconomicNew) {
	c := colly.NewCollector()
	c.OnHTML("div.info div.content", func(e *colly.HTMLElement) {
		ecoNew.Content = &e.Text
	})
	c.Visit(ecoNew.URL)
}
