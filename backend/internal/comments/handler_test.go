package comments_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	mock_comments "github.com/PegasusMKD/travel-dream-board/internal/comments/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupCommentsRouter(svc comments.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := comments.NewHandler(svc)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func TestCommentsHandler_CreateComment(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_comments.Service)
		router := setupCommentsRouter(mockSvc)

		body := &comments.Comment{Content: "Test comment", CommentedOn: db.CommentedOnActivities, CommentedOnUuid: "act-1"}
		bodyBytes, _ := json.Marshal(body)

		expectedComment := &comments.Comment{Uuid: "com-1", Content: "Test comment"}

		mockSvc.On("CreateComment", mock.Anything, mock.MatchedBy(func(c *comments.Comment) bool {
			return c.Content == "Test comment"
		})).Return(expectedComment, nil)

		req, _ := http.NewRequest(http.MethodPost, "/comments/", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Bad JSON", func(t *testing.T) {
		mockSvc := new(mock_comments.Service)
		router := setupCommentsRouter(mockSvc)

		req, _ := http.NewRequest(http.MethodPost, "/comments/", bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_comments.Service)
		router := setupCommentsRouter(mockSvc)

		body := &comments.Comment{Content: "Test comment", CommentedOn: db.CommentedOnActivities, CommentedOnUuid: "act-1"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("CreateComment", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPost, "/comments/", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCommentsHandler_UpdateCommentByUuid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_comments.Service)
		router := setupCommentsRouter(mockSvc)

		uuid := "com-1"
		body := map[string]string{"content": "Updated content"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateCommentByUuid", mock.Anything, uuid, "Updated content").Return(nil)

		req, _ := http.NewRequest(http.MethodPatch, "/comments/"+uuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Bad JSON", func(t *testing.T) {
		mockSvc := new(mock_comments.Service)
		router := setupCommentsRouter(mockSvc)

		uuid := "com-1"
		req, _ := http.NewRequest(http.MethodPatch, "/comments/"+uuid, bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_comments.Service)
		router := setupCommentsRouter(mockSvc)

		uuid := "com-1"
		body := map[string]string{"content": "Updated content"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateCommentByUuid", mock.Anything, uuid, mock.Anything).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPatch, "/comments/"+uuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestCommentsHandler_DeleteCommentByUuid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mock_comments.Service)
		router := setupCommentsRouter(mockSvc)

		uuid := "com-1"

		mockSvc.On("DeleteCommentByUuid", mock.Anything, uuid).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/comments/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mock_comments.Service)
		router := setupCommentsRouter(mockSvc)

		uuid := "com-1"

		mockSvc.On("DeleteCommentByUuid", mock.Anything, uuid).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodDelete, "/comments/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
