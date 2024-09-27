package common

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCommonRepositoryImpl(t *testing.T) {
	t.Run("SavePageViewRecord", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error when opening a stub database connection: %s\n", err)
		}
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repo := newCommonRepositoryImpl(sqlxDB)

		data := PageViewRecordData{
			App:    "app",
			UserID: "12345",
			URL:    "/url",
			Ips:    []string{"127.0.0.1"},
		}

		mock.ExpectExec("INSERT INTO page_view_record").
			WithArgs(sqlmock.AnyArg(), data.App, data.UserID, data.URL, pq.StringArray(data.Ips), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err = repo.SavePageViewRecord(data)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("SavePageViewRecord error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("error when opening a stub database connection: %s\n", err)
		}
		defer db.Close()

		sqlxDB := sqlx.NewDb(db, "sqlmock")
		repo := newCommonRepositoryImpl(sqlxDB)

		data := PageViewRecordData{
			App:    "app",
			UserID: "12345",
			URL:    "/url",
			Ips:    []string{"127.0.0.1"},
		}

		mock.ExpectExec("INSERT INTO page_view_record").
			WithArgs(sqlmock.AnyArg(), data.App, data.UserID, data.URL, pq.StringArray(data.Ips), sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf("error"))

		err = repo.SavePageViewRecord(data)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
