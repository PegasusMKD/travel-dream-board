package comments

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
)

type Service interface {
	CreateComment(ctx context.Context, data *Comment) (*Comment, error)
	FindAllCommentsByRelatedEntity(ctx context.Context, _type db.CommentedOn, uuid string) ([]*Comment, error)
	UpdateCommentByUuid(ctx context.Context, uuid string, content string) error
	DeleteCommentByUuid(ctx context.Context, uuid string) error
}

type commentsServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &commentsServiceImpl{
		repo: repo,
	}
}

func (svc *commentsServiceImpl) CreateComment(ctx context.Context, data *Comment) (*Comment, error) {
	return svc.repo.CreateComment(ctx, data)
}

func (svc *commentsServiceImpl) FindAllCommentsByRelatedEntity(ctx context.Context, _type db.CommentedOn, uuid string) ([]*Comment, error) {
	return svc.repo.FindAllCommentsByRelatedEntity(ctx, _type, uuid)
}

func (svc *commentsServiceImpl) UpdateCommentByUuid(ctx context.Context, uuid string, content string) error {
	return svc.repo.UpdateCommentByUuid(ctx, uuid, content)
}

func (svc *commentsServiceImpl) DeleteCommentByUuid(ctx context.Context, uuid string) error {
	return svc.repo.DeleteCommentByUuid(ctx, uuid)
}
