package sharetokens_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/sharetokens"
	"github.com/PegasusMKD/travel-dream-board/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShareTokensService_CreateShareToken(t *testing.T) {
	mockRepo := new(mocks.MocksharetokensRepository)
	svc := sharetokens.NewService(mockRepo)

	ctx := context.Background()
	boardUuid := "board-1"
	expectedToken := &sharetokens.ShareToken{Token: "test-token", BoardUuid: boardUuid}

	mockRepo.On("CreateShareToken", ctx, mock.AnythingOfType("string"), boardUuid).Return(expectedToken, nil)

	result, err := svc.CreateShareToken(ctx, boardUuid)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, result)
	mockRepo.AssertExpectations(t)
}

func TestShareTokensService_GetShareToken(t *testing.T) {
	mockRepo := new(mocks.MocksharetokensRepository)
	svc := sharetokens.NewService(mockRepo)

	ctx := context.Background()
	token := "test-token"
	expectedToken := &sharetokens.ShareToken{Token: token}

	mockRepo.On("GetShareToken", ctx, token).Return(expectedToken, nil)

	result, err := svc.GetShareToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, expectedToken, result)
	mockRepo.AssertExpectations(t)
}

func TestShareTokensService_DeleteShareToken(t *testing.T) {
	mockRepo := new(mocks.MocksharetokensRepository)
	svc := sharetokens.NewService(mockRepo)

	ctx := context.Background()
	token := "test-token"
	boardUuid := "board-1"

	mockRepo.On("DeleteShareToken", ctx, token, boardUuid).Return(nil)

	err := svc.DeleteShareToken(ctx, token, boardUuid)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestShareTokensService_GetShareTokensForBoard(t *testing.T) {
	mockRepo := new(mocks.MocksharetokensRepository)
	svc := sharetokens.NewService(mockRepo)

	ctx := context.Background()
	boardUuid := "board-1"
	expectedTokens := []*sharetokens.ShareToken{{Token: "test-token"}}

	mockRepo.On("GetShareTokensForBoard", ctx, boardUuid).Return(expectedTokens, nil)

	result, err := svc.GetShareTokensForBoard(ctx, boardUuid)
	assert.NoError(t, err)
	assert.Equal(t, expectedTokens, result)
	mockRepo.AssertExpectations(t)
}
