package boards

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/accomodations"
	"github.com/PegasusMKD/travel-dream-board/internal/activities"
	"github.com/PegasusMKD/travel-dream-board/internal/transport"
)

type Service interface {
	CreateBoard(ctx context.Context, data *Board) (*Board, error)
	GetBoardById(ctx context.Context, uuid string) (*AggregatedBoard, error)
	GetAllBoards(ctx context.Context, userUuid string) ([]*Board, error)
	UpdateBoardById(ctx context.Context, uuid string, data *Board) error
	DeleteBoardById(ctx context.Context, uuid string) error
}

type serviceImpl struct {
	repo Repository

	accomodationSvc accomodations.Service
	activitiesSvc   activities.Service
	transportSvc    transport.Service
}

func NewService(repo Repository, accomodationSvc accomodations.Service, activitiesSvc activities.Service, transportSvc transport.Service) Service {
	return &serviceImpl{
		repo:            repo,
		accomodationSvc: accomodationSvc,
		activitiesSvc:   activitiesSvc,
		transportSvc:    transportSvc,
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

	activs, err := svc.activitiesSvc.GetActivitiesByBoardId(ctx, uuid)
	if err != nil {
		return nil, err
	}

	transports, err := svc.transportSvc.GetTransportsByBoardId(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &AggregatedBoard{
		Board:         *board,
		Accomodations: accoms,
		Activities:    activs,
		Transport:     transports,
	}, nil
}

func (svc *serviceImpl) GetAllBoards(ctx context.Context, userUuid string) ([]*Board, error) {
	return svc.repo.GetAllBoards(ctx, userUuid)
}

func (svc *serviceImpl) UpdateBoardById(ctx context.Context, uuid string, data *Board) error {
	return svc.repo.UpdateBoardById(ctx, uuid, data)
}

func (svc *serviceImpl) DeleteBoardById(ctx context.Context, uuid string) error {
	return svc.repo.DeleteBoardById(ctx, uuid)
}
