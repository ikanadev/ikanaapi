package econewscron

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type DbEconomicNew struct {
	ID        uuid.UUID      `db:"id"`
	Title     string         `db:"title"`
	URL       string         `db:"url"`
	Date      *time.Time     `db:"date"`
	Image     *string        `db:"image"`
	Summary   *string        `db:"summary"`
	Company   string         `db:"company"`
	Tags      pq.StringArray `db:"tags"`
	Sentiment int8           `db:"sentiment"`
	CreatedAt time.Time      `db:"created_at"`
	DeletedAt *time.Time     `db:"deleted_at"`
}

func (dbEcoNew *DbEconomicNew) ToEconomicNew() *EconomicNew {
	return &EconomicNew{
		ID:        dbEcoNew.ID,
		Title:     dbEcoNew.Title,
		URL:       dbEcoNew.URL,
		Date:      dbEcoNew.Date,
		Image:     dbEcoNew.Image,
		Summary:   dbEcoNew.Summary,
		Company:   dbEcoNew.Company,
		Tags:      dbEcoNew.Tags,
		Sentiment: dbEcoNew.Sentiment,
		CreatedAt: dbEcoNew.CreatedAt,
		DeletedAt: dbEcoNew.DeletedAt,
	}
}

type EconomicNew struct {
	ID        uuid.UUID  `json:"id"`
	Title     string     `json:"title"`
	URL       string     `json:"url"`
	Date      *time.Time `json:"date"`
	Image     *string    `json:"image"`
	Content   *string    // used only to save content in memory to ask AI for summary
	Summary   *string    `json:"summary"`
	Company   string     `json:"company"`
	Tags      []string   `json:"tags"`
	Sentiment int8       `json:"sentiment"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

func (ecoNew *EconomicNew) ToDbEconomicNew() *DbEconomicNew {
	return &DbEconomicNew{
		ID:        ecoNew.ID,
		Title:     ecoNew.Title,
		URL:       ecoNew.URL,
		Date:      ecoNew.Date,
		Image:     ecoNew.Image,
		Summary:   ecoNew.Summary,
		Company:   ecoNew.Company,
		Tags:      pq.StringArray(ecoNew.Tags),
		Sentiment: ecoNew.Sentiment,
		CreatedAt: ecoNew.CreatedAt,
		DeletedAt: ecoNew.DeletedAt,
	}
}
