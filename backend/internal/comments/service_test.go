package comments_test

import (
	"context"
	"testing"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	mock_comments "github.com/PegasusMKD/travel-dream-board/internal/comments/mocks"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/stretchr/testify/assert"
)

func TestCommentsService_CreateComment(t *testing.T) {
	mockRepo := new(mock_comments.Repository)
	svc := comments.NewService(mockRepo)

	ctx := context.Background()
	commentData := &comments.Comment{Content: "Test comment"}
	expectedComment := &comments.Comment{Uuid: "comm-1", Content: "Test comment"}

	mockRepo.On("CreateComment", ctx, commentData).Return(expectedComment, nil)

	result, err := svc.CreateComment(ctx, commentData)
	assert.NoError(t, err)
	assert.Equal(t, expectedComment, result)
	mockRepo.AssertExpectations(t)
}

func TestCommentsService_FindAllCommentsByRelatedEntity(t *testing.T) {
	mockRepo := new(mock_comments.Repository)
	svc := comments.NewService(mockRepo)

	ctx := context.Background()
	relatedType := db.CommentedOnActivities
	uuid := "act-1"
	expectedComments := []*comments.Comment{{Uuid: "comm-1"}}

	mockRepo.On("FindAllCommentsByRelatedEntity", ctx, relatedType, uuid).Return(expectedComments, nil)

	result, err := svc.FindAllCommentsByRelatedEntity(ctx, relatedType, uuid)
	assert.NoError(t, err)
	assert.Equal(t, expectedComments, result)
	mockRepo.AssertExpectations(t)
}

func TestCommentsService_UpdateCommentByUuid(t *testing.T) {
	mockRepo := new(mock_comments.Repository)
	svc := comments.NewService(mockRepo)

	ctx := context.Background()
	uuid := "comm-1"
	content := "Updated content"

	mockRepo.On("UpdateCommentByUuid", ctx, uuid, content).Return(nil)

	err := svc.UpdateCommentByUuid(ctx, uuid, content)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCommentsService_DeleteCommentByUuid(t *testing.T) {
	mockRepo := new(mock_comments.Repository)
	svc := comments.NewService(mockRepo)

	ctx := context.Background()
	uuid := "comm-1"

	mockRepo.On("DeleteCommentByUuid", ctx, uuid).Return(nil)

	err := svc.DeleteCommentByUuid(ctx, uuid)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
