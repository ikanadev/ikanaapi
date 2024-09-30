package common

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCommonRepository struct {
	mock.Mock
}

// SavePublicFeedback implements CommonRepository.
func (m *MockCommonRepository) SavePublicFeedback(data PublicFeedbackData) error {
	args := m.Called(data)
	return args.Error(0)
}

// SavePageViewRecord implements CommonRepository.
func (m *MockCommonRepository) SavePageViewRecord(data PageViewRecordData) error {
	args := m.Called(data)
	return args.Error(0)
}

func TestCommonServices(t *testing.T) {
	t.Run("PostPageViewRecord", func(t *testing.T) {
		mockRepo := new(MockCommonRepository)
		// request data
		jsonData := `{
			"userId": "12345",
			"app": "crisis",
			"url": "/"
		}`
		pageViewData := PageViewRecordData{
			App:    "crisis",
			UserID: "12345",
			URL:    "/",
			Ips:    []string{"127.0.0.1"},
		}
		// request
		req := httptest.NewRequest(http.MethodPost, "/common/page_view", strings.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderXForwardedFor, "127.0.0.1")
		rec := httptest.NewRecorder()
		context := echo.New().NewContext(req, rec)
		// expectations
		mockRepo.On("SavePageViewRecord", pageViewData).Return(nil)
		// test
		handler := postPageViewRecord(mockRepo)
		if assert.NoError(t, handler(context)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
		mockRepo.AssertExpectations(t)
	})

	t.Run("PostPublicFeedback", func(t *testing.T) {
		mockRepo := new(MockCommonRepository)
		// request data
		jsonData := `{
			"app": "crisis",
			"userId": "12345",
			"content": "test content"
		}`
		publicFeedbackData := PublicFeedbackData{
			App:     "crisis",
			UserID:  "12345",
			Content: "test content",
			Ips:     []string{"127.0.0.1"},
		}
		// request
		req := httptest.NewRequest(http.MethodPost, "/common/public_feedback", strings.NewReader(jsonData))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderXForwardedFor, "127.0.0.1")
		rec := httptest.NewRecorder()
		context := echo.New().NewContext(req, rec)
		// expectations
		mockRepo.On("SavePublicFeedback", publicFeedbackData).Return(nil)
		// test
		handler := postPublicFeedback(mockRepo)
		if assert.NoError(t, handler(context)) {
			assert.Equal(t, http.StatusCreated, rec.Code)
		}
		mockRepo.AssertExpectations(t)
	})
}
