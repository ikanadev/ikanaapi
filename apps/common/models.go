package common

import (
	"time"

	"github.com/google/uuid"
)

type PageViewRecordData struct {
	App    string   `json:"app"`
	UserID string   `json:"userId"`
	URL    string   `json:"url"`
	Ips    []string `json:"ips"`
}

type PageViewRecord struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	PageViewRecordData
}
