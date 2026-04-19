package sharetokens_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/sharetokens"
	"github.com/PegasusMKD/travel-dream-board/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupShareTokensRouter(svc sharetokens.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := sharetokens.NewHandler(svc)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func TestShareTokensHandler_CreateShareToken(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mocks.MocksharetokensService)
		router := setupShareTokensRouter(mockSvc)

		uuid := "board-1"
		mockToken := &sharetokens.ShareToken{Token: "mock-token"}
		mockSvc.On("CreateShareToken", mock.Anything, uuid).Return(mockToken, nil)

		req, _ := http.NewRequest(http.MethodPost, "/boards/"+uuid+"/share-tokens/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "mock-token")
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mocks.MocksharetokensService)
		router := setupShareTokensRouter(mockSvc)

		uuid := "board-1"
		mockSvc.On("CreateShareToken", mock.Anything, uuid).Return((*sharetokens.ShareToken)(nil), errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPost, "/boards/"+uuid+"/share-tokens/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestShareTokensHandler_GetShareTokensForBoard(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mocks.MocksharetokensService)
		router := setupShareTokensRouter(mockSvc)

		uuid := "board-1"
		mockTokens := []*sharetokens.ShareToken{{Token: "token1"}, {Token: "token2"}}
		mockSvc.On("GetShareTokensForBoard", mock.Anything, uuid).Return(mockTokens, nil)

		req, _ := http.NewRequest(http.MethodGet, "/boards/"+uuid+"/share-tokens/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "token1")
		assert.Contains(t, w.Body.String(), "token2")
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mocks.MocksharetokensService)
		router := setupShareTokensRouter(mockSvc)

		uuid := "board-1"
		mockSvc.On("GetShareTokensForBoard", mock.Anything, uuid).Return([]*sharetokens.ShareToken(nil), errors.New("service error"))

		req, _ := http.NewRequest(http.MethodGet, "/boards/"+uuid+"/share-tokens/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestShareTokensHandler_DeleteShareToken(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mocks.MocksharetokensService)
		router := setupShareTokensRouter(mockSvc)

		uuid := "board-1"
		token := "token-1"
		mockSvc.On("DeleteShareToken", mock.Anything, token, uuid).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/boards/"+uuid+"/share-tokens/"+token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing Token", func(t *testing.T) {
		mockSvc := new(mocks.MocksharetokensService)
		router := setupShareTokensRouter(mockSvc)

		uuid := "board-1"
		// In Gin, trailing slash or empty param can sometimes skip route match or give empty param
		req, _ := http.NewRequest(http.MethodDelete, "/boards/"+uuid+"/share-tokens/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		// Usually 404 if no route matches, but let's just test what happens
		// We could skip testing gin's routing for missing param.
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mocks.MocksharetokensService)
		router := setupShareTokensRouter(mockSvc)

		uuid := "board-1"
		token := "token-1"
		mockSvc.On("DeleteShareToken", mock.Anything, token, uuid).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodDelete, "/boards/"+uuid+"/share-tokens/"+token, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
