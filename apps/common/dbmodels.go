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

type DbPublicFeedback struct {
	ID        uuid.UUID      `db:"id"`
	App       string         `db:"app"`
	UserID    string         `db:"user_id"`
	Ips       pq.StringArray `db:"ips"`
	Section   string         `db:"section"`
	Content   string         `db:"content"`
	CreatedAt time.Time      `db:"created_at"`
}

func (publicFeedback DbPublicFeedback) ToPublicFeedback() PublicFeedback {
	return PublicFeedback{
		ID:        publicFeedback.ID,
		CreatedAt: publicFeedback.CreatedAt,
		PublicFeedbackData: PublicFeedbackData{
			App:     publicFeedback.App,
			UserID:  publicFeedback.UserID,
			Ips:     publicFeedback.Ips,
			Section: publicFeedback.Section,
			Content: publicFeedback.Content,
		},
	}
}
