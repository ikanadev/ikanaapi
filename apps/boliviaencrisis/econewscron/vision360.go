package econewscron

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

func getVision360NewDetails(ecoNew *EconomicNew) {
	c := colly.NewCollector()
	monthMap := map[string]string{
		"enero":      "January",
		"febrero":    "February",
		"marzo":      "March",
		"abril":      "April",
		"mayo":       "May",
		"junio":      "June",
		"julio":      "July",
		"agosto":     "August",
		"septiembre": "September",
		"octubre":    "octubre",
		"noviembre":  "November",
		"diciembre":  "December",
	}
	c.OnHTML("div.noticia-fecha", func(e *colly.HTMLElement) {
		dateParts := strings.Split(strings.TrimSpace(e.Text), " ")
		if len(dateParts) != 6 {
			return
		}
		month, ok := monthMap[strings.ToLower(dateParts[3])]
		if !ok {
			return
		}
		formattedDate := fmt.Sprintf("%s %s %s", dateParts[1], month, dateParts[5])
		parsedDate, err := time.Parse("2 January 2006", formattedDate)
		if err != nil {
			return
		}
		ecoNew.Date = &parsedDate
	})
	c.OnHTML("div.noticia-contenido", func(e *colly.HTMLElement) {
		ecoNew.Content = &e.Text
	})

	c.Visit(ecoNew.URL)
}

func getVision360News() []*EconomicNew {
	c := colly.NewCollector()
	news := make([]*EconomicNew, 0)
	baseURL := "https://www.vision360.bo"
	c.OnHTML("article.listado-noticias-relacionadas", func(e *colly.HTMLElement) {
		title := ""
		title += e.ChildText("h3.text-noticia-simple-volanta")
		if len(title) > 0 {
			title += ", "
		}
		title += e.ChildText("h2.text-noticia-simple-titulo")
		url := baseURL + e.ChildAttr("a", "href")
		img := e.ChildAttr("img", "src")

		news = append(news, &EconomicNew{
			ID:      uuid.New(),
			Title:   title,
			URL:     url,
			Image:   &img,
			Company: "Visión 360",
		})
	})

	c.Visit("https://www.vision360.bo/economia")
	return news
}