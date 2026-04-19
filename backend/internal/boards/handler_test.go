package boards_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/boards"
	mock_boards "github.com/PegasusMKD/travel-dream-board/internal/boards/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupBoardsRouter(svc boards.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Add a dummy middleware to inject user_uuid for testing
	r.Use(func(c *gin.Context) {
		c.Set("user_uuid", "user-123")
		c.Next()
	})

	handler := boards.NewHandler(svc)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func setupBoardsRouterNoUser(svc boards.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := boards.NewHandler(svc)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func TestBoardsHandler_CreateBoard(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		body := &boards.Board{Name: "Test Board", LocationName: "Test Location"}
		bodyBytes, _ := json.Marshal(body)

		userUuid := "user-123"
		expectedBoard := &boards.Board{Name: "Test Board", LocationName: "Test Location", UserUuid: &userUuid, Uuid: "uuid-1"}

		mockSvc.On("CreateBoard", mock.Anything, mock.MatchedBy(func(b *boards.Board) bool {
			return b.Name == "Test Board" && b.LocationName == "Test Location" && *b.UserUuid == userUuid
		})).Return(expectedBoard, nil)

		req, _ := http.NewRequest(http.MethodPost, "/boards/", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Bad JSON", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		req, _ := http.NewRequest(http.MethodPost, "/boards/", bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Missing User UUID", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouterNoUser(mockSvc)

		body := &boards.Board{Name: "Test Board", LocationName: "Test Location"}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/boards/", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		body := &boards.Board{Name: "Test Board", LocationName: "Test Location"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("CreateBoard", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPost, "/boards/", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestBoardsHandler_GetBoardById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		boardUuid := "uuid-1"
		expectedBoard := &boards.AggregatedBoard{
			Board: boards.Board{Uuid: boardUuid, Name: "Test Board"},
		}

		mockSvc.On("GetBoardById", mock.Anything, boardUuid).Return(expectedBoard, nil)

		req, _ := http.NewRequest(http.MethodGet, "/boards/"+boardUuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing UUID", func(t *testing.T) {
		// Handled by Gin routing usually, but test empty param logic if possible
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		boardUuid := "uuid-1"
		mockSvc.On("GetBoardById", mock.Anything, boardUuid).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodGet, "/boards/"+boardUuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestBoardsHandler_GetAllBoards(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		userUuid := "user-123"
		expectedBoards := []*boards.Board{{Uuid: "uuid-1", Name: "Test Board"}}

		mockSvc.On("GetAllBoards", mock.Anything, userUuid).Return(expectedBoards, nil)

		req, _ := http.NewRequest(http.MethodGet, "/boards/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing User UUID", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouterNoUser(mockSvc)

		req, _ := http.NewRequest(http.MethodGet, "/boards/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		userUuid := "user-123"
		mockSvc.On("GetAllBoards", mock.Anything, userUuid).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodGet, "/boards/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestBoardsHandler_UpdateBoardById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		boardUuid := "uuid-1"
		body := &boards.Board{Name: "Updated Board", LocationName: "Updated Location"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateBoardById", mock.Anything, boardUuid, mock.MatchedBy(func(b *boards.Board) bool {
			return b.Name == "Updated Board" && b.LocationName == "Updated Location"
		})).Return(nil)

		req, _ := http.NewRequest(http.MethodPatch, "/boards/"+boardUuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Bad JSON", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		boardUuid := "uuid-1"
		req, _ := http.NewRequest(http.MethodPatch, "/boards/"+boardUuid, bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		boardUuid := "uuid-1"
		body := &boards.Board{Name: "Updated Board", LocationName: "Updated Location"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateBoardById", mock.Anything, boardUuid, mock.Anything).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPatch, "/boards/"+boardUuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestBoardsHandler_DeleteBoardById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		boardUuid := "uuid-1"

		mockSvc.On("DeleteBoardById", mock.Anything, boardUuid).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/boards/"+boardUuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_boards.Service)
		router := setupBoardsRouter(mockSvc)

		boardUuid := "uuid-1"

		mockSvc.On("DeleteBoardById", mock.Anything, boardUuid).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodDelete, "/boards/"+boardUuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
