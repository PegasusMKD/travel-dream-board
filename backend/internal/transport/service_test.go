package transport_test

import (
	"context"
	"errors"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	mock_comments "github.com/PegasusMKD/travel-dream-board/internal/comments/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/scrape_audit"
	"github.com/PegasusMKD/travel-dream-board/mocks"

	"github.com/PegasusMKD/travel-dream-board/internal/transport"
	mock_transport "github.com/PegasusMKD/travel-dream-board/internal/transport/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
	mock_votes "github.com/PegasusMKD/travel-dream-board/internal/votes/mocks"
	"github.com/stretchr/testify/assert"
)

func TestTransportService_CreateTransport(t *testing.T) {
	ctx := context.Background()
	url := "http://example.com"
	boardUuid := "board-1"
	userUuid := "user-1"
	title := "Scraped Title"

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_transport.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := transport.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		scrapedData := &scrapeaudit.ScrapeResult{Title: &title}
		mockScrape.On("Scrape", ctx, url).Return(scrapedData, nil)

		expectedTransport := &transport.Transport{
			Url:       url,
			Title:     title,
			BoardUuid: boardUuid,
			UserUuid:  userUuid,
		}

		createdTransport := &transport.Transport{Uuid: "trans-1", Title: title}
		mockRepo.On("CreateTransport", ctx, expectedTransport).Return(createdTransport, nil)

		result, err := svc.CreateTransport(ctx, url, []byte(nil), "", boardUuid, userUuid)
		assert.NoError(t, err)
		assert.Equal(t, createdTransport, result)
	})

	t.Run("Scrape Error", func(t *testing.T) {
		mockRepo := new(mock_transport.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := transport.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		mockScrape.On("Scrape", ctx, url).Return(nil, errors.New("scrape error"))

		result, err := svc.CreateTransport(ctx, url, []byte(nil), "", boardUuid, userUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTransportService_GetTransportById(t *testing.T) {
	ctx := context.Background()
	transUuid := "trans-1"

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_transport.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := transport.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		trans := &transport.Transport{Uuid: transUuid, Title: "Test"}
		mockComms := []*comments.Comment{{Uuid: "comm-1"}}
		mockVts := []*votes.Vote{{Uuid: "vote-1"}}

		mockRepo.On("GetTransportById", ctx, transUuid).Return(trans, nil)
		mockComments.On("FindAllCommentsByRelatedEntity", ctx, db.CommentedOnTransport, transUuid).Return(mockComms, nil)
		mockVotes.On("FindAllVotesByRelatedEntity", ctx, db.VotedOnTransport, transUuid).Return(mockVts, nil)

		result, err := svc.GetTransportById(ctx, transUuid)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, transUuid, result.Uuid)
		assert.Len(t, result.Comments, 1)
		assert.Len(t, result.Votes, 1)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo := new(mock_transport.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := transport.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		mockRepo.On("GetTransportById", ctx, transUuid).Return(nil, errors.New("repo error"))

		result, err := svc.GetTransportById(ctx, transUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTransportService_GetTransportsByBoardId(t *testing.T) {
	mockRepo := new(mock_transport.Repository)
	svc := transport.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	boardUuid := "board-1"
	expectedTrans := []*transport.Transport{{Uuid: "trans-1"}}

	mockRepo.On("GetAllTransportsByBoardId", ctx, boardUuid).Return(expectedTrans, nil)

	result, err := svc.GetTransportsByBoardId(ctx, boardUuid)
	assert.NoError(t, err)
	assert.Equal(t, expectedTrans, result)
}

func TestTransportService_UpdateTransportById(t *testing.T) {
	mockRepo := new(mock_transport.Repository)
	svc := transport.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	transUuid := "trans-1"
	transData := &transport.Transport{Title: "Updated"}

	mockRepo.On("UpdateTransportById", ctx, transUuid, transData).Return(nil)

	err := svc.UpdateTransportById(ctx, transUuid, transData)
	assert.NoError(t, err)
}

func TestTransportService_DeleteTransportById(t *testing.T) {
	mockRepo := new(mock_transport.Repository)
	svc := transport.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	transUuid := "trans-1"

	mockRepo.On("DeleteTransportById", ctx, transUuid).Return(nil)

	err := svc.DeleteTransportById(ctx, transUuid)
	assert.NoError(t, err)
}
