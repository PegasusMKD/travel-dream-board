package transport_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/transport"
	mock_transport "github.com/PegasusMKD/travel-dream-board/internal/transport/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTransportRouter(svc transport.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Add a dummy middleware to inject user_uuid for testing
	r.Use(func(c *gin.Context) {
		c.Set("user_uuid", "user-123")
		c.Next()
	})

	handler := transport.NewHandler(svc)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func setupTransportRouterNoUser(svc transport.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := transport.NewHandler(svc)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func TestTransportHandler_CreateTransport(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		url := "http://example.com"
		boardUuid := "board-1"
		userUuid := "user-123"

		expectedTransport := &transport.Transport{Title: "Test Transport"}

		mockSvc.On("CreateTransport", mock.Anything, url, boardUuid, userUuid).Return(expectedTransport, nil)

		req, _ := http.NewRequest(http.MethodPost, "/transport/?url="+url+"&boardUuid="+boardUuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing Url", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		req, _ := http.NewRequest(http.MethodPost, "/transport/?boardUuid=board-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing Board Uuid", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		req, _ := http.NewRequest(http.MethodPost, "/transport/?url=http://example.com", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Missing User UUID", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouterNoUser(mockSvc)

		req, _ := http.NewRequest(http.MethodPost, "/transport/?url=test&boardUuid=board-1", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		url := "http://example.com"
		boardUuid := "board-1"
		mockSvc.On("CreateTransport", mock.Anything, url, boardUuid, mock.Anything).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPost, "/transport/?url="+url+"&boardUuid="+boardUuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTransportHandler_GetTransportById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		uuid := "act-1"
		expectedTransport := &transport.AggregatedTransport{Transport: transport.Transport{Uuid: uuid}}

		mockSvc.On("GetTransportById", mock.Anything, uuid).Return(expectedTransport, nil)

		req, _ := http.NewRequest(http.MethodGet, "/transport/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		uuid := "act-1"
		mockSvc.On("GetTransportById", mock.Anything, uuid).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodGet, "/transport/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTransportHandler_UpdateTransportById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		uuid := "act-1"
		body := &transport.Transport{Title: "Updated", Url: "http://example.com", BoardUuid: "board-1"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateTransportById", mock.Anything, uuid, mock.MatchedBy(func(b *transport.Transport) bool {
			return b.Title == "Updated"
		})).Return(nil)

		req, _ := http.NewRequest(http.MethodPatch, "/transport/"+uuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Bad JSON", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		uuid := "act-1"
		req, _ := http.NewRequest(http.MethodPatch, "/transport/"+uuid, bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		uuid := "act-1"
		body := &transport.Transport{Title: "Updated", Url: "http://example.com", BoardUuid: "board-1"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateTransportById", mock.Anything, uuid, mock.Anything).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPatch, "/transport/"+uuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestTransportHandler_DeleteTransportById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		uuid := "act-1"

		mockSvc.On("DeleteTransportById", mock.Anything, uuid).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/transport/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_transport.Service)
		router := setupTransportRouter(mockSvc)

		uuid := "act-1"

		mockSvc.On("DeleteTransportById", mock.Anything, uuid).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodDelete, "/transport/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
