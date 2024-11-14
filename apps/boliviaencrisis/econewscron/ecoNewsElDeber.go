package econewscron

import (
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/google/uuid"
)

type ElDeberSource struct {
	Name string
}

func newElDeberSource() *ElDeberSource {
	return &ElDeberSource{
		Name: "El Deber",
	}
}

func (laRazonSource *ElDeberSource) GetEcoNews() []*EconomicNew {
	c := colly.NewCollector()

	news := make([]*EconomicNew, 0)
	baseURL := "https://eldeber.com.bo"

	c.OnHTML("div.component--medium", func(e *colly.HTMLElement) {
		url := baseURL + e.ChildAttr("a.nota-link", "href")
		title := e.ChildText("a.nota-link")
		imgUrl := baseURL + e.ChildAttr("img", "src")
		ecoNew := EconomicNew{
			ID:      uuid.New(),
			Company: laRazonSource.Name,
			Image:   &imgUrl,
			Title:   title,
			URL:     url,
			Tags:    make([]string, 0),
		}
		news = append(news, &ecoNew)
	})

	c.Visit(baseURL + "/economia/")
	return news
}

func (laRazonSource *ElDeberSource) GetEcoNewDetails(ecoNew *EconomicNew) {
	c := colly.NewCollector()

	c.OnHTML("div.dateNote.mobile", func(e *colly.HTMLElement) {
		date := time.Now().UTC()
		dateParts := strings.Split(strings.TrimSpace(e.Text), ",")
		if len(dateParts) > 1 {
			dateParts := strings.Split(strings.TrimSpace(dateParts[0]), " ")
			month, ok := monthMap[strings.ToLower(dateParts[2])]
			if ok && len(dateParts) == 5 && month != "" {
				formattedDate := fmt.Sprintf("%s %s %s", dateParts[0], month, dateParts[4])
				parsedDate, err := time.Parse("2 January 2006", formattedDate)
				if err == nil {
					date = parsedDate.UTC()
				}
			}
		}
		ecoNew.Date = &date
	})

	content := ""
	c.OnHTML("div.text-editor p", func(e *colly.HTMLElement) {
		content += e.Text + "\n"
	})

	c.Visit(ecoNew.URL)

	ecoNew.Content = &content
}
