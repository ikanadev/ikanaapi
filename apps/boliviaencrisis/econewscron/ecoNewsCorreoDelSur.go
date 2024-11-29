package econewscron

import (
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

type CorreoDelSurSource struct {
	Name string
}

func newCorreoDelSurSource() *CorreoDelSurSource {
	return &CorreoDelSurSource{
		Name: "Correo Del Sur",
	}
}

func (correoDelSurSource *CorreoDelSurSource) GetEcoNews() []*EconomicNew {
	c := colly.NewCollector()
	baseURL := "https://correodelsur.com"
	ecoNews := make([]*EconomicNew, 0)

	c.OnHTML("div.order-2 ku-card", func(e *colly.HTMLElement) {
		title := e.Attr("title")
		url := e.Attr("url")
		imgUrl := e.Attr("image-url")
		ecoNew := EconomicNew{
			ID:      uuid.New(),
			Company: correoDelSurSource.Name,
			Title:   title,
			URL:     url,
			Image:   &imgUrl,
			Tags:    make([]string, 0),
		}
		ecoNews = append(ecoNews, &ecoNew)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error correo del sur: ", err)
	})

	c.Visit(baseURL + "/economia")
	return ecoNews
}

func (correoDelSurSource *CorreoDelSurSource) GetEcoNewDetails(ecoNew *EconomicNew) {
	c := colly.NewCollector()
	c.OnHTML("div.text-sm.text-neutral-500.mb-4 span:nth-child(2)", func(e *colly.HTMLElement) {
		date := time.Now().UTC()
		dateStr := strings.Split(strings.TrimSpace(e.Text), " ")
		parsedDate, err := time.Parse("02/01/2006", dateStr[0])
		if err == nil {
			date = parsedDate.UTC()
		}
		ecoNew.Date = &date
	})
	c.OnHTML("section.uk-container div.font-sans", func(e *colly.HTMLElement) {
		ecoNew.Content = &e.Text
	})
	c.Visit(ecoNew.URL)
}
