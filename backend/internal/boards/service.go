package boards

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
)

type Service interface {
	CreateBoard(ctx context.Context, data *Board) (*Board, error)
	GetBoardById(ctx context.Context, uuid string) (*AggregatedBoard, error)
	GetAllBoards(ctx context.Context) ([]*Board, error)
	UpdateBoardById(ctx context.Context, uuid string, data *Board) error
	DeleteBoardById(ctx context.Context, uuid string) error
}

type serviceImpl struct {
	repo            Repository
	accomodationSvc accomodations.Service
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (svc *serviceImpl) CreateBoard(ctx context.Context, data *Board) (*Board, error) {
	return svc.repo.CreateBoard(ctx, data)
}

func (svc *serviceImpl) GetBoardById(ctx context.Context, uuid string) (*AggregatedBoard, error) {
	board, err := svc.repo.GetBoardById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	accoms, err := svc.accomodationSvc.GetAccomodationsByBoardId(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &AggregatedBoard{
		Board:         *board,
		Accomodations: accoms,
	}, nil
}

func (svc *serviceImpl) GetAllBoards(ctx context.Context) ([]*Board, error) {
	return svc.repo.GetAllBoards(ctx)
}

func (svc *serviceImpl) UpdateBoardById(ctx context.Context, uuid string, data *Board) error {
	return svc.repo.UpdateBoardById(ctx, uuid, data)
}

func (svc *serviceImpl) DeleteBoardById(ctx context.Context, uuid string) error {
	return svc.repo.DeleteBoardById(ctx, uuid)
}
