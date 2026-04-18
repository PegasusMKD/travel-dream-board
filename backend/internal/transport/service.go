package transport

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Service interface {
	CreateTransport(ctx context.Context, url string, userUuid string) (*Transport, error)
	GetTransportById(ctx context.Context, uuid string) (*AggregatedTransport, error)
	GetTransportsByBoardId(ctx context.Context, uuid string) ([]*Transport, error)
	UpdateTransportById(ctx context.Context, uuid string, data *Transport) error
	DeleteTransportById(ctx context.Context, uuid string) error
}

type accomodationServiceImpl struct {
	repo Repository

	commentsSvc comments.Service
	votesSvc    votes.Service
}

func NewService(repo Repository, commentsSvc comments.Service, votesSvc votes.Service) Service {
	return &accomodationServiceImpl{
		repo:        repo,
		commentsSvc: commentsSvc,
		votesSvc:    votesSvc,
	}
}

func (svc *accomodationServiceImpl) CreateTransport(ctx context.Context, url string, userUuid string) (*Transport, error) {
	// TODO: Scrape the data
	data := Transport{
		Url:      url,
		UserUuid: userUuid,
	}

	return svc.repo.CreateTransport(ctx, &data)
}

func (svc *accomodationServiceImpl) GetTransportById(ctx context.Context, uuid string) (*AggregatedTransport, error) {
	accomodation, err := svc.repo.GetTransportById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	comms, err := svc.commentsSvc.FindAllCommentsByRelatedEntity(ctx, db.CommentedOnTransport, uuid)
	if err != nil {
		return nil, err
	}

	votes, err := svc.votesSvc.FindAllVotesByRelatedEntity(ctx, db.VotedOnTransport, uuid)
	if err != nil {
		return nil, err
	}

	return &AggregatedTransport{
		Transport: *accomodation,
		Comments:  comms,
		Votes:     votes,
	}, nil
}

func (svc *accomodationServiceImpl) GetTransportsByBoardId(ctx context.Context, uuid string) ([]*Transport, error) {
	return svc.repo.GetAllTransportsByBoardId(ctx, uuid)
}

func (svc *accomodationServiceImpl) UpdateTransportById(ctx context.Context, uuid string, data *Transport) error {
	return svc.repo.UpdateTransportById(ctx, uuid, data)
}

func (svc *accomodationServiceImpl) DeleteTransportById(ctx context.Context, uuid string) error {
	return svc.repo.DeleteTransportById(ctx, uuid)
}
