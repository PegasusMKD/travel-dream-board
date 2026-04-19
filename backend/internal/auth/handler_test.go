package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/auth"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/mocks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

func setupAuthRouter(svc auth.Service, oauth *oauth2.Config) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := auth.NewHandler(svc, oauth)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func TestAuthHandler_GoogleLogin(t *testing.T) {
	mockSvc := new(mocks.MockauthService)
	oauthConf := &oauth2.Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
		RedirectURL: "http://localhost/callback",
		Scopes:      []string{"email"},
	}

	router := setupAuthRouter(mockSvc, oauthConf)

	req, _ := http.NewRequest(http.MethodGet, "/auth/google/login", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code)
	// Should have a Set-Cookie header for oauthstate
	assert.Contains(t, w.Header().Get("Set-Cookie"), "oauthstate=")
	// Location should be google auth url
	assert.Contains(t, w.Header().Get("Location"), "https://accounts.google.com/o/oauth2/auth")
}

func TestAuthHandler_GoogleCallback(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mocks.MockauthService)
		oauthConf := &oauth2.Config{}
		router := setupAuthRouter(mockSvc, oauthConf)

		email := "test@example.com"
		userUuid, _ := uuid.NewRandom()
		pgUuid := pgtype.UUID{Bytes: userUuid, Valid: true}
		user := db.User{Uuid: pgUuid, Name: "Test User", Email: &email}
		token := "mock-jwt-token"

		mockSvc.On("ProcessGoogleCallback", mock.Anything, "mock-code").Return(&user, token, nil)

		req, _ := http.NewRequest(http.MethodGet, "/auth/google/callback?state=mock-state&code=mock-code", nil)
		req.AddCookie(&http.Cookie{Name: "oauthstate", Value: "mock-state"})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "Successfully logged in")
		assert.Contains(t, w.Header().Get("Set-Cookie"), "jwt_token=mock-jwt-token")
		mockSvc.AssertExpectations(t)
	})

	t.Run("Missing State Cookie", func(t *testing.T) {
		mockSvc := new(mocks.MockauthService)
		oauthConf := &oauth2.Config{}
		router := setupAuthRouter(mockSvc, oauthConf)

		req, _ := http.NewRequest(http.MethodGet, "/auth/google/callback?state=mock-state&code=mock-code", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("State Mismatch", func(t *testing.T) {
		mockSvc := new(mocks.MockauthService)
		oauthConf := &oauth2.Config{}
		router := setupAuthRouter(mockSvc, oauthConf)

		req, _ := http.NewRequest(http.MethodGet, "/auth/google/callback?state=different-state&code=mock-code", nil)
		req.AddCookie(&http.Cookie{Name: "oauthstate", Value: "mock-state"})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Missing Code", func(t *testing.T) {
		mockSvc := new(mocks.MockauthService)
		oauthConf := &oauth2.Config{}
		router := setupAuthRouter(mockSvc, oauthConf)

		req, _ := http.NewRequest(http.MethodGet, "/auth/google/callback?state=mock-state", nil)
		req.AddCookie(&http.Cookie{Name: "oauthstate", Value: "mock-state"})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mocks.MockauthService)
		oauthConf := &oauth2.Config{}
		router := setupAuthRouter(mockSvc, oauthConf)

		mockSvc.On("ProcessGoogleCallback", mock.Anything, "mock-code").Return((*db.User)(nil), "", assert.AnError)

		req, _ := http.NewRequest(http.MethodGet, "/auth/google/callback?state=mock-state&code=mock-code", nil)
		req.AddCookie(&http.Cookie{Name: "oauthstate", Value: "mock-state"})
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
