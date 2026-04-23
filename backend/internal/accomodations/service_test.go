package accomodations_test

import (
	"context"
	"errors"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	mock_accomodations "github.com/PegasusMKD/travel-dream-board/internal/accomodations/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	mock_comments "github.com/PegasusMKD/travel-dream-board/internal/comments/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/scrape_audit"
	"github.com/PegasusMKD/travel-dream-board/mocks"

	"github.com/PegasusMKD/travel-dream-board/internal/votes"
	mock_votes "github.com/PegasusMKD/travel-dream-board/internal/votes/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAccomodationsService_CreateAccomodation(t *testing.T) {
	ctx := context.Background()
	url := "http://example.com"
	boardUuid := "board-1"
	userUuid := "user-1"
	title := "Scraped Title"

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_accomodations.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := accomodations.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		scrapedData := &scrapeaudit.ScrapeResult{Title: &title}
		mockScrape.On("Scrape", ctx, url).Return(scrapedData, nil)

		expectedAccomodation := &accomodations.Accomodation{
			Url:       url,
			Title:     title,
			BoardUuid: boardUuid,
			UserUuid:  userUuid,
		}

		createdAccomodation := &accomodations.Accomodation{Uuid: "acc-1", Title: title}
		mockRepo.On("CreateAccomodation", ctx, expectedAccomodation).Return(createdAccomodation, nil)

		result, err := svc.CreateAccomodation(ctx, url, []byte(nil), "", boardUuid, userUuid)
		assert.NoError(t, err)
		assert.Equal(t, createdAccomodation, result)
	})

	t.Run("Scrape Error", func(t *testing.T) {
		mockRepo := new(mock_accomodations.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := accomodations.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		mockScrape.On("Scrape", ctx, url).Return(nil, errors.New("scrape error"))

		result, err := svc.CreateAccomodation(ctx, url, []byte(nil), "", boardUuid, userUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestAccomodationsService_GetAccomodationById(t *testing.T) {
	ctx := context.Background()
	accUuid := "acc-1"

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_accomodations.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := accomodations.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		acc := &accomodations.Accomodation{Uuid: accUuid, Title: "Test"}
		mockComms := []*comments.Comment{{Uuid: "comm-1"}}
		mockVts := []*votes.Vote{{Uuid: "vote-1"}}

		mockRepo.On("GetAccomodationById", ctx, accUuid).Return(acc, nil)
		mockComments.On("FindAllCommentsByRelatedEntity", ctx, db.CommentedOnAccomodation, accUuid).Return(mockComms, nil)
		mockVotes.On("FindAllVotesByRelatedEntity", ctx, db.VotedOnAccomodation, accUuid).Return(mockVts, nil)

		result, err := svc.GetAccomodationById(ctx, accUuid)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, accUuid, result.Uuid)
		assert.Len(t, result.Comments, 1)
		assert.Len(t, result.Votes, 1)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo := new(mock_accomodations.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := accomodations.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		mockRepo.On("GetAccomodationById", ctx, accUuid).Return(nil, errors.New("repo error"))

		result, err := svc.GetAccomodationById(ctx, accUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestAccomodationsService_GetAccomodationsByBoardId(t *testing.T) {
	mockRepo := new(mock_accomodations.Repository)
	svc := accomodations.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	boardUuid := "board-1"
	expectedAccs := []*accomodations.Accomodation{{Uuid: "acc-1"}}

	mockRepo.On("GetAllAccomodationsByBoardId", ctx, boardUuid).Return(expectedAccs, nil)

	result, err := svc.GetAccomodationsByBoardId(ctx, boardUuid)
	assert.NoError(t, err)
	assert.Equal(t, expectedAccs, result)
}

func TestAccomodationsService_UpdateAccomodationById(t *testing.T) {
	mockRepo := new(mock_accomodations.Repository)
	svc := accomodations.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	accUuid := "acc-1"
	accData := &accomodations.Accomodation{Title: "Updated"}

	mockRepo.On("UpdateAccomodationById", ctx, accUuid, accData).Return(nil)

	err := svc.UpdateAccomodationById(ctx, accUuid, accData)
	assert.NoError(t, err)
}

func TestAccomodationsService_DeleteAccomodationById(t *testing.T) {
	mockRepo := new(mock_accomodations.Repository)
	svc := accomodations.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	accUuid := "acc-1"

	mockRepo.On("DeleteAccomodationById", ctx, accUuid).Return(nil)

	err := svc.DeleteAccomodationById(ctx, accUuid)
	assert.NoError(t, err)
}
