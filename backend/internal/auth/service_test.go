package auth_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/auth"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type mockRoundTripper struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestAuthService_ProcessGoogleCallback(t *testing.T) {
	mockRepo := new(mocks.MockauthRepository)
	jwtSecret := "my-secret-key"

	oauthConfig := &oauth2.Config{
		ClientID:     "mock-client-id",
		ClientSecret: "mock-client-secret",
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://mock.token.url",
			AuthURL:  "https://mock.auth.url",
		},
	}
	svc := auth.NewService(mockRepo, oauthConfig, jwtSecret)

	t.Run("Success", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					if req.URL.String() == "https://mock.token.url" {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(`{"access_token": "mock-token", "token_type": "Bearer", "expires_in": 3600}`)),
							Header:     make(http.Header),
						}, nil
					}
					if req.URL.String() == "https://www.googleapis.com/oauth2/v2/userinfo" {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(`{"id": "123", "email": "test@test.com", "name": "Test User", "picture": "http://example.com/pic"}`)),
							Header:     make(http.Header),
						}, nil
					}
					return nil, errors.New("unexpected url: " + req.URL.String())
				},
			},
		}

		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockClient)

		mockID := uuid.New()
		pgMockID := pgtype.UUID{Bytes: mockID, Valid: true}
		email := "test@test.com"
		mockRepo.On("UpsertUser", ctx, mock.Anything).Return(&db.User{
			Uuid:  pgMockID,
			Email: &email,
		}, nil).Once()

		user, token, err := svc.ProcessGoogleCallback(ctx, "mock-code")

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.NotEmpty(t, token)
		assert.Equal(t, pgMockID, user.Uuid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Exchange Error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					if req.URL.String() == "https://mock.token.url" {
						return &http.Response{
							StatusCode: http.StatusBadRequest,
							Body:       io.NopCloser(bytes.NewBufferString(`{"error": "invalid_grant"}`)),
							Header:     make(http.Header),
						}, nil
					}
					return nil, errors.New("unexpected url: " + req.URL.String())
				},
			},
		}

		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockClient)

		user, token, err := svc.ProcessGoogleCallback(ctx, "mock-code")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
	})

	t.Run("UserInfo API Error", func(t *testing.T) {
		mockClient := &http.Client{
			Transport: &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					if req.URL.String() == "https://mock.token.url" {
						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(`{"access_token": "mock-token", "token_type": "Bearer", "expires_in": 3600}`)),
							Header:     make(http.Header),
						}, nil
					}
					if req.URL.String() == "https://www.googleapis.com/oauth2/v2/userinfo" {
						return nil, errors.New("api error")
					}
					return nil, errors.New("unexpected url: " + req.URL.String())
				},
			},
		}

		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, mockClient)

		user, token, err := svc.ProcessGoogleCallback(ctx, "mock-code")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Empty(t, token)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	mockRepo := new(mocks.MockauthRepository)
	oauthConfig := &oauth2.Config{}
	jwtSecret := "my-secret-key"

	svc := auth.NewService(mockRepo, oauthConfig, jwtSecret)

	t.Run("Valid Token", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "user-1",
		})
		tokenString, _ := token.SignedString([]byte(jwtSecret))

		sub, err := svc.ValidateToken(tokenString)
		assert.NoError(t, err)
		assert.Equal(t, "user-1", sub)
	})

	t.Run("Invalid Secret", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "user-1",
		})
		tokenString, _ := token.SignedString([]byte("wrong-secret"))

		sub, err := svc.ValidateToken(tokenString)
		assert.Error(t, err)
		assert.Empty(t, sub)
	})

	t.Run("Missing Sub Claim", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": "test@test.com",
		})
		tokenString, _ := token.SignedString([]byte(jwtSecret))

		sub, err := svc.ValidateToken(tokenString)
		assert.Error(t, err)
		assert.Empty(t, sub)
		assert.Equal(t, "sub claim missing or not string", err.Error())
	})

	t.Run("Invalid Token String", func(t *testing.T) {
		sub, err := svc.ValidateToken("invalid.token.string")
		assert.Error(t, err)
		assert.Empty(t, sub)
	})
}
