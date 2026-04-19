package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/middleware"
	"github.com/PegasusMKD/travel-dream-board/internal/sharetokens"
	"github.com/PegasusMKD/travel-dream-board/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequireAuth(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockAuthSvc := new(mocks.MockauthService)
		mockAuthSvc.On("ValidateToken", "valid-token").Return("user-123", nil)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(middleware.RequireAuth(mockAuthSvc))
		r.GET("/test", func(c *gin.Context) {
			val, _ := c.Get("user_uuid")
			c.String(http.StatusOK, val.(string))
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "valid-token"})
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "user-123", w.Body.String())
	})

	t.Run("Missing Cookie", func(t *testing.T) {
		mockAuthSvc := new(mocks.MockauthService)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(middleware.RequireAuth(mockAuthSvc))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid Token", func(t *testing.T) {
		mockAuthSvc := new(mocks.MockauthService)
		mockAuthSvc.On("ValidateToken", "invalid-token").Return("", errors.New("invalid"))

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.Use(middleware.RequireAuth(mockAuthSvc))
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "invalid-token"})
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRequireBoardAccess(t *testing.T) {
	t.Run("Success with JWT", func(t *testing.T) {
		mockAuthSvc := new(mocks.MockauthService)
		mockShareSvc := new(mocks.MocksharetokensService)
		mockAuthSvc.On("ValidateToken", "valid-token").Return("user-123", nil)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.GET("/boards/:uuid/test", middleware.RequireBoardAccess(mockAuthSvc, mockShareSvc), func(c *gin.Context) {
			val, _ := c.Get("user_uuid")
			c.String(http.StatusOK, val.(string))
		})

		req, _ := http.NewRequest(http.MethodGet, "/boards/00000000-0000-0000-0000-000000000000/test", nil)
		req.AddCookie(&http.Cookie{Name: "jwt_token", Value: "valid-token"})
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "user-123", w.Body.String())
	})

	t.Run("Success with Share Token Header", func(t *testing.T) {
		mockAuthSvc := new(mocks.MockauthService)
		mockShareSvc := new(mocks.MocksharetokensService)

		validUUIDStr := "11111111-1111-1111-1111-111111111111"
		mockShareToken := &sharetokens.ShareToken{BoardUuid: validUUIDStr}
		mockShareSvc.On("GetShareToken", mock.Anything, "valid-share-token").Return(mockShareToken, nil)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.GET("/boards/:uuid/test", middleware.RequireBoardAccess(mockAuthSvc, mockShareSvc), func(c *gin.Context) {
			val, _ := c.Get("share_token")
			c.String(http.StatusOK, val.(string))
		})

		req, _ := http.NewRequest(http.MethodGet, "/boards/"+validUUIDStr+"/test", nil)
		req.Header.Set("X-Share-Token", "valid-share-token")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "valid-share-token", w.Body.String())
	})

	t.Run("Forbidden No Tokens", func(t *testing.T) {
		mockAuthSvc := new(mocks.MockauthService)
		mockShareSvc := new(mocks.MocksharetokensService)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.GET("/boards/:uuid/test", middleware.RequireBoardAccess(mockAuthSvc, mockShareSvc), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/boards/00000000-0000-0000-0000-000000000000/test", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Forbidden Invalid Share Token", func(t *testing.T) {
		mockAuthSvc := new(mocks.MockauthService)
		mockShareSvc := new(mocks.MocksharetokensService)

		mockShareSvc.On("GetShareToken", mock.Anything, "invalid-share-token").Return(nil, errors.New("not found"))

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.GET("/boards/:uuid/test", middleware.RequireBoardAccess(mockAuthSvc, mockShareSvc), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req, _ := http.NewRequest(http.MethodGet, "/boards/00000000-0000-0000-0000-000000000000/test", nil)
		req.Header.Set("X-Share-Token", "invalid-share-token")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Forbidden Share Token Board Mismatch", func(t *testing.T) {
		mockAuthSvc := new(mocks.MockauthService)
		mockShareSvc := new(mocks.MocksharetokensService)

		mockShareToken := &sharetokens.ShareToken{BoardUuid: "22222222-2222-2222-2222-222222222222"}
		mockShareSvc.On("GetShareToken", mock.Anything, "valid-share-token").Return(mockShareToken, nil)

		gin.SetMode(gin.TestMode)
		r := gin.New()
		r.GET("/boards/:uuid/test", middleware.RequireBoardAccess(mockAuthSvc, mockShareSvc), func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		// Requesting a different board uuid than what the token is for
		req, _ := http.NewRequest(http.MethodGet, "/boards/11111111-1111-1111-1111-111111111111/test", nil)
		req.Header.Set("X-Share-Token", "valid-share-token")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}
