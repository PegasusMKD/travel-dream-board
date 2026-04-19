package votes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
	"github.com/PegasusMKD/travel-dream-board/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupVotesRouter(svc votes.Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	handler := votes.NewHandler(svc)
	handler.RegisterRoutes(r.Group(""))
	return r
}

func TestVotesHandler_CreateVote(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mocks.MockvotesService)
		router := setupVotesRouter(mockSvc)

		body := &votes.Vote{Rank: 1, VotedOn: db.VotedOnActivities, VotedOnUuid: "act-1"}
		bodyBytes, _ := json.Marshal(body)

		expectedVote := &votes.Vote{Uuid: "vote-1", Rank: 1}

		mockSvc.On("CreateVote", mock.Anything, mock.MatchedBy(func(v *votes.Vote) bool {
			return v.Rank == 1
		})).Return(expectedVote, nil)

		req, _ := http.NewRequest(http.MethodPost, "/votes/", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Bad JSON", func(t *testing.T) {
		mockSvc := new(mocks.MockvotesService)
		router := setupVotesRouter(mockSvc)

		req, _ := http.NewRequest(http.MethodPost, "/votes/", bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mocks.MockvotesService)
		router := setupVotesRouter(mockSvc)

		body := &votes.Vote{Rank: 1, VotedOn: db.VotedOnActivities, VotedOnUuid: "act-1"}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("CreateVote", mock.Anything, mock.Anything).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPost, "/votes/", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestVotesHandler_UpdateVoteByUuid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mocks.MockvotesService)
		router := setupVotesRouter(mockSvc)

		uuid := "vote-1"
		body := map[string]int32{"rank": 2}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateVoteByUuid", mock.Anything, uuid, int32(2)).Return(nil)

		req, _ := http.NewRequest(http.MethodPatch, "/votes/"+uuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Bad JSON", func(t *testing.T) {
		mockSvc := new(mocks.MockvotesService)
		router := setupVotesRouter(mockSvc)

		uuid := "vote-1"
		req, _ := http.NewRequest(http.MethodPatch, "/votes/"+uuid, bytes.NewBuffer([]byte("{bad json}")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mocks.MockvotesService)
		router := setupVotesRouter(mockSvc)

		uuid := "vote-1"
		body := map[string]int32{"rank": 2}
		bodyBytes, _ := json.Marshal(body)

		mockSvc.On("UpdateVoteByUuid", mock.Anything, uuid, mock.Anything).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodPatch, "/votes/"+uuid, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestVotesHandler_DeleteVoteByUuid(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mocks.MockvotesService)
		router := setupVotesRouter(mockSvc)

		uuid := "vote-1"

		mockSvc.On("DeleteVoteByUuid", mock.Anything, uuid).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, "/votes/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(mocks.MockvotesService)
		router := setupVotesRouter(mockSvc)

		uuid := "vote-1"

		mockSvc.On("DeleteVoteByUuid", mock.Anything, uuid).Return(errors.New("service error"))

		req, _ := http.NewRequest(http.MethodDelete, "/votes/"+uuid, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
