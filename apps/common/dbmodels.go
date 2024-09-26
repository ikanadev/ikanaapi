package common

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type DbPageViewRecord struct {
	ID        uuid.UUID      `db:"id"`
	App       string         `db:"app"`
	UserID    string         `db:"user_id"`
	URL       string         `db:"url"`
	Ips       pq.StringArray `db:"ips"`
	CreatedAt time.Time      `db:"created_at"`
}

func (pageViewRecord DbPageViewRecord) ToPageViewRecord() PageViewRecord {
	return PageViewRecord{
		ID:        pageViewRecord.ID,
		CreatedAt: pageViewRecord.CreatedAt,
		PageViewRecordData: PageViewRecordData{
			App:    pageViewRecord.App,
			UserID: pageViewRecord.UserID,
			URL:    pageViewRecord.URL,
			Ips:    pageViewRecord.Ips,
		},
	}
}
