package common

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func setupTestData(t *testing.T) (CommonRepositoryImpl, sqlmock.Sqlmock) {
	db, dbMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error when opening a stub database connection: %s\n", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := newCommonRepositoryImpl(sqlxDB)
	return repo, dbMock
}

func TestCommonRepositoryImpl(t *testing.T) {
	t.Run("SavePageViewRecord", func(t *testing.T) {
		repo, dbMock := setupTestData(t)
		data := PageViewRecordData{
			App:    "app",
			UserID: "12345",
			URL:    "/url",
			Ips:    []string{"127.0.0.1"},
		}

		dbMock.ExpectExec("INSERT INTO page_view_record").
			WithArgs(sqlmock.AnyArg(), data.App, data.UserID, data.URL, pq.StringArray(data.Ips), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.SavePageViewRecord(data)

		assert.NoError(t, err)
		assert.NoError(t, dbMock.ExpectationsWereMet())
	})

	t.Run("SavePageViewRecord error", func(t *testing.T) {
		repo, dbMock := setupTestData(t)
		data := PageViewRecordData{
			App:    "app",
			UserID: "12345",
			URL:    "/url",
			Ips:    []string{"127.0.0.1"},
		}

		dbMock.ExpectExec("INSERT INTO page_view_record").
			WithArgs(sqlmock.AnyArg(), data.App, data.UserID, data.URL, pq.StringArray(data.Ips), sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf("error"))
		err := repo.SavePageViewRecord(data)

		assert.Error(t, err)
		assert.NoError(t, dbMock.ExpectationsWereMet())
	})

	t.Run("SavePublicFeedback", func(t *testing.T) {
		repo, dbMock := setupTestData(t)
		data := PublicFeedbackData{
			App:     "app",
			UserID:  "12345",
			Ips:     []string{"127.0.0.1"},
			Content: "content",
		}

		dbMock.ExpectExec("INSERT INTO public_feedback").
			WithArgs(
				sqlmock.AnyArg(),
				data.App,
				data.UserID,
				pq.StringArray(data.Ips),
				data.Section,
				data.Content,
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		err := repo.SavePublicFeedback(data)

		assert.NoError(t, err)
		assert.NoError(t, dbMock.ExpectationsWereMet())
	})

	t.Run("SavePublicFeedback", func(t *testing.T) {
		repo, dbMock := setupTestData(t)
		data := PublicFeedbackData{
			App:     "app",
			UserID:  "12345",
			Ips:     []string{"127.0.0.1"},
			Content: "content",
		}

		dbMock.ExpectExec("INSERT INTO public_feedback").
			WithArgs(
				sqlmock.AnyArg(),
				data.App,
				data.UserID,
				pq.StringArray(data.Ips),
				data.Section,
				data.Content,
				sqlmock.AnyArg(),
			).
			WillReturnError(fmt.Errorf("error"))
		err := repo.SavePublicFeedback(data)

		assert.Error(t, err)
		assert.NoError(t, dbMock.ExpectationsWereMet())
	})
}
