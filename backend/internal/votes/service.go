package votes

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
)

type Service interface {
	CreateVote(ctx context.Context, data *Vote) (*Vote, error)
	FindAllVotesByRelatedEntity(ctx context.Context, _type db.VotedOn, uuid string) ([]*Vote, error)
	UpdateVoteByUuid(ctx context.Context, uuid string, rank int32) error
	DeleteVoteByUuid(ctx context.Context, uuid string) error
}

type votesServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &votesServiceImpl{
		repo: repo,
	}
}

func (svc *votesServiceImpl) CreateVote(ctx context.Context, data *Vote) (*Vote, error) {
	return svc.repo.CreateVote(ctx, data)
}

func (svc *votesServiceImpl) FindAllVotesByRelatedEntity(ctx context.Context, _type db.VotedOn, uuid string) ([]*Vote, error) {
	return svc.repo.FindAllVotesByRelatedEntity(ctx, _type, uuid)
}

func (svc *votesServiceImpl) UpdateVoteByUuid(ctx context.Context, uuid string, rank int32) error {
	return svc.repo.UpdateVoteByUuid(ctx, uuid, rank)
}

func (svc *votesServiceImpl) DeleteVoteByUuid(ctx context.Context, uuid string) error {
	return svc.repo.DeleteVoteByUuid(ctx, uuid)
}
