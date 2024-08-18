package common

import "time"

type TimeData struct {
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
	ArchivedAt *time.Time `json:"archivedAt"`
	DeletedAt  *time.Time `json:"deletedAt"`
}
type DBTimeData struct {
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
	ArchivedAt *time.Time `db:"archived_at"`
	DeletedAt  *time.Time `db:"deleted_at"`
}
