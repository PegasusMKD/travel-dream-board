package activities_test

import (
	"context"
	"errors"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/activities"
	mock_activities "github.com/PegasusMKD/travel-dream-board/internal/activities/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	mock_comments "github.com/PegasusMKD/travel-dream-board/internal/comments/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/scrape_audit"
	"github.com/PegasusMKD/travel-dream-board/mocks"

	"github.com/PegasusMKD/travel-dream-board/internal/votes"
	mock_votes "github.com/PegasusMKD/travel-dream-board/internal/votes/mocks"
	"github.com/stretchr/testify/assert"
)

func TestActivitiesService_CreateActivity(t *testing.T) {
	ctx := context.Background()
	url := "http://example.com"
	boardUuid := "board-1"
	userUuid := "user-1"
	title := "Scraped Title"

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_activities.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := activities.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		scrapedData := &scrapeaudit.ScrapeResult{Title: &title}
		mockScrape.On("Scrape", ctx, url).Return(scrapedData, nil)

		expectedActivity := &activities.Activity{
			Url:       url,
			Title:     title,
			BoardUuid: boardUuid,
			UserUuid:  userUuid,
		}

		createdActivity := &activities.Activity{Uuid: "act-1", Title: title}
		mockRepo.On("CreateActivity", ctx, expectedActivity).Return(createdActivity, nil)

		result, err := svc.CreateActivity(ctx, url, []byte(nil), "", boardUuid, userUuid)
		assert.NoError(t, err)
		assert.Equal(t, createdActivity, result)
	})

	t.Run("Scrape Error", func(t *testing.T) {
		mockRepo := new(mock_activities.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := activities.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		mockScrape.On("Scrape", ctx, url).Return(nil, errors.New("scrape error"))

		result, err := svc.CreateActivity(ctx, url, []byte(nil), "", boardUuid, userUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestActivitiesService_GetActivityById(t *testing.T) {
	ctx := context.Background()
	actUuid := "act-1"

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_activities.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := activities.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		act := &activities.Activity{Uuid: actUuid, Title: "Test"}
		mockComms := []*comments.Comment{{Uuid: "comm-1"}}
		mockVts := []*votes.Vote{{Uuid: "vote-1"}}

		mockRepo.On("GetActivityById", ctx, actUuid).Return(act, nil)
		mockComments.On("FindAllCommentsByRelatedEntity", ctx, db.CommentedOnActivities, actUuid).Return(mockComms, nil)
		mockVotes.On("FindAllVotesByRelatedEntity", ctx, db.VotedOnActivities, actUuid).Return(mockVts, nil)

		result, err := svc.GetActivityById(ctx, actUuid)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, actUuid, result.Uuid)
		assert.Len(t, result.Comments, 1)
		assert.Len(t, result.Votes, 1)
	})

	t.Run("Repo Error", func(t *testing.T) {
		mockRepo := new(mock_activities.Repository)
		mockComments := new(mock_comments.Service)
		mockVotes := new(mock_votes.Service)
		mockScrape := new(mocks.MockscrapeprocessService)

		svc := activities.NewService(mockRepo, mockComments, mockVotes, mockScrape)

		mockRepo.On("GetActivityById", ctx, actUuid).Return(nil, errors.New("repo error"))

		result, err := svc.GetActivityById(ctx, actUuid)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestActivitiesService_GetActivitiesByBoardId(t *testing.T) {
	mockRepo := new(mock_activities.Repository)
	svc := activities.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	boardUuid := "board-1"
	expectedActs := []*activities.Activity{{Uuid: "act-1"}}

	mockRepo.On("GetAllActivitiesByBoardId", ctx, boardUuid).Return(expectedActs, nil)

	result, err := svc.GetActivitiesByBoardId(ctx, boardUuid)
	assert.NoError(t, err)
	assert.Equal(t, expectedActs, result)
}

func TestActivitiesService_UpdateActivityById(t *testing.T) {
	mockRepo := new(mock_activities.Repository)
	svc := activities.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	actUuid := "act-1"
	actData := &activities.Activity{Title: "Updated"}

	mockRepo.On("UpdateActivityById", ctx, actUuid, actData).Return(nil)

	err := svc.UpdateActivityById(ctx, actUuid, actData)
	assert.NoError(t, err)
}

func TestActivitiesService_DeleteActivityById(t *testing.T) {
	mockRepo := new(mock_activities.Repository)
	svc := activities.NewService(mockRepo, nil, nil, nil)

	ctx := context.Background()
	actUuid := "act-1"

	mockRepo.On("DeleteActivityById", ctx, actUuid).Return(nil)

	err := svc.DeleteActivityById(ctx, actUuid)
	assert.NoError(t, err)
}
