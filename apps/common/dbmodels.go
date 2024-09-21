package common

import (
	"time"

	"github.com/google/uuid"
)

type DbPageViewRecord struct {
	ID        uuid.UUID `db:"id"`
	App       string    `db:"app"`
	UserID    string    `db:"user_id"`
	URL       string    `db:"url"`
	CreatedAt time.Time `db:"created_at"`
}

func (pageViewRecord DbPageViewRecord) ToPageViewRecord() PageViewRecord {
	return PageViewRecord{
		ID:        pageViewRecord.ID,
		CreatedAt: pageViewRecord.CreatedAt,
		PageViewRecordData: PageViewRecordData{
			App:    pageViewRecord.App,
			UserID: pageViewRecord.UserID,
			URL:    pageViewRecord.URL,
		},
	}
}
