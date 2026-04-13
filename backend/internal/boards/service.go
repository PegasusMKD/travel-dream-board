package boards

import "context"

type Service interface {
	CreateBoard(ctx context.Context, data *Board) (*Board, error)
	GetBoardById(ctx context.Context, uuid string) (*Board, error)
	GetAllBoards(ctx context.Context) ([]*Board, error)
	UpdateBoardById(ctx context.Context, uuid string, data *Board) error
	DeleteBoardById(ctx context.Context, uuid string) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{
		repo: repo,
	}
}

func (svc *serviceImpl) CreateBoard(ctx context.Context, data *Board) (*Board, error) {
	return svc.repo.CreateBoard(ctx, data)
}

func (svc *serviceImpl) GetBoardById(ctx context.Context, uuid string) (*Board, error) {
	return svc.repo.GetBoardById(ctx, uuid)
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
