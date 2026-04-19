package votes_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
	"github.com/PegasusMKD/travel-dream-board/mocks"
	"github.com/stretchr/testify/assert"
)

func TestVotesService_CreateVote(t *testing.T) {
	mockRepo := new(mocks.MockvotesRepository)
	svc := votes.NewService(mockRepo)

	ctx := context.Background()
	voteData := &votes.Vote{Rank: 1}
	expectedVote := &votes.Vote{Uuid: "vote-1", Rank: 1}

	mockRepo.On("CreateVote", ctx, voteData).Return(expectedVote, nil)

	result, err := svc.CreateVote(ctx, voteData)
	assert.NoError(t, err)
	assert.Equal(t, expectedVote, result)
	mockRepo.AssertExpectations(t)
}

func TestVotesService_FindAllVotesByRelatedEntity(t *testing.T) {
	mockRepo := new(mocks.MockvotesRepository)
	svc := votes.NewService(mockRepo)

	ctx := context.Background()
	relatedType := db.VotedOnActivities
	uuid := "act-1"
	expectedVotes := []*votes.Vote{{Uuid: "vote-1"}}

	mockRepo.On("FindAllVotesByRelatedEntity", ctx, relatedType, uuid).Return(expectedVotes, nil)

	result, err := svc.FindAllVotesByRelatedEntity(ctx, relatedType, uuid)
	assert.NoError(t, err)
	assert.Equal(t, expectedVotes, result)
	mockRepo.AssertExpectations(t)
}

func TestVotesService_UpdateVoteByUuid(t *testing.T) {
	mockRepo := new(mocks.MockvotesRepository)
	svc := votes.NewService(mockRepo)

	ctx := context.Background()
	uuid := "vote-1"
	var rank int32 = 2

	mockRepo.On("UpdateVoteByUuid", ctx, uuid, rank).Return(nil)

	err := svc.UpdateVoteByUuid(ctx, uuid, rank)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestVotesService_DeleteVoteByUuid(t *testing.T) {
	mockRepo := new(mocks.MockvotesRepository)
	svc := votes.NewService(mockRepo)

	ctx := context.Background()
	uuid := "vote-1"

	mockRepo.On("DeleteVoteByUuid", ctx, uuid).Return(nil)

	err := svc.DeleteVoteByUuid(ctx, uuid)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
