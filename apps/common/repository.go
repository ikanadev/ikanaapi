package common

import (
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CommonRepository interface {
	SavePageViewRecord(data PageViewRecordData) error
}

type CommonRepositoryImpl struct {
	db *sqlx.DB
}

// SavePageViewRecord implements CommonRepository.
func (r CommonRepositoryImpl) SavePageViewRecord(data PageViewRecordData) error {
	dbData := DbPageViewRecord{
		ID:        uuid.New(),
		App:       data.App,
		UserID:    data.UserID,
		URL:       data.URL,
		CreatedAt: time.Now().UTC(),
	}
	sql := `INSERT INTO page_view_record
	(id, app, user_id, url, created_at)
	VALUES
	(:id, :app, :user_id, :url, :created_at);`
	_, err := r.db.NamedExec(sql, dbData)
	return err
}

func newCommonRepositoryImpl(db *sqlx.DB) CommonRepositoryImpl {
	return CommonRepositoryImpl{db}
}
