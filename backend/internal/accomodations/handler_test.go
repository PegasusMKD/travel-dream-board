package accomodations_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	mock_accomodations "github.com/PegasusMKD/travel-dream-board/internal/accomodations/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAccomodationsRouter(svc accomodations.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("user_uuid", "user-123")
		c.Next()
	})

	handler := accomodations.NewHandler(svc)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func setupAccomodationsRouterNoUser(svc accomodations.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := accomodations.NewHandler(svc)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func TestAccomodationsHandler_CreateAccomodation(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		url := "http://example.com"
		boardUuid := "board-1"
		userUuid := "user-123"

		expectedAccomodation := &accomodations.Accomodation{Title: "Test Acc"}

		mockSvc.On("CreateAccomodation", mock.Anything, url, []byte(nil), "", boardUuid, userUuid).Return(expectedAccomodation, nil)

		req, _ := http.NewRequest(http.MethodPost, "/accomodations/?url="+url+"&boardUuid="+boardUuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing Url", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		req, _ := http.NewRequest(http.MethodPost, "/accomodations/?boardUuid=board-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing Board Uuid", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		req, _ := http.NewRequest(http.MethodPost, "/accomodations/?url=http://example.com", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing User UUID", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouterNoUser(mockSvc)

		req, _ := http.NewRequest(http.MethodPost, "/accomodations/?url=test&boardUuid=board-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		url := "http://example.com"
		boardUuid := "board-1"
		mockSvc.On("CreateAccomodation", mock.Anything, url, []byte(nil), "", boardUuid, mock.Anything).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPost, "/accomodations/?url="+url+"&boardUuid="+boardUuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAccomodationsHandler_GetAccomodationById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		uuid := "acc-1"
		expectedAccomodation := &accomodations.AggregatedAccomodation{Accomodation: accomodations.Accomodation{Uuid: uuid}}

		mockSvc.On("GetAccomodationById", mock.Anything, uuid).Return(expectedAccomodation, nil)

		req, _ := http.NewRequest(http.MethodGet, "/accomodations/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		uuid := "acc-1"
		mockSvc.On("GetAccomodationById", mock.Anything, uuid).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodGet, "/accomodations/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAccomodationsHandler_UpdateAccomodationById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		uuid := "acc-1"
		body := &accomodations.Accomodation{Title: "Updated", Url: "http://example.com", BoardUuid: "board-1"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateAccomodationById", mock.Anything, uuid, mock.MatchedBy(func(b *accomodations.Accomodation) bool {
			return b.Title == "Updated"
		})).Return(nil)

		req, _ := http.NewRequest(http.MethodPatch, "/accomodations/"+uuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Bad JSON", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		uuid := "acc-1"
		req, _ := http.NewRequest(http.MethodPatch, "/accomodations/"+uuid, bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code) // The handler currently returns 500 for bad JSON
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		uuid := "acc-1"
		body := &accomodations.Accomodation{Title: "Updated", Url: "http://example.com", BoardUuid: "board-1"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateAccomodationById", mock.Anything, uuid, mock.Anything).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPatch, "/accomodations/"+uuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestAccomodationsHandler_DeleteAccomodationById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		uuid := "acc-1"

		mockSvc.On("DeleteAccomodationById", mock.Anything, uuid).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/accomodations/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_accomodations.Service)
		router := setupAccomodationsRouter(mockSvc)

		uuid := "acc-1"

		mockSvc.On("DeleteAccomodationById", mock.Anything, uuid).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodDelete, "/accomodations/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
