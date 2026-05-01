package memories

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
)

type Repository interface {
	CreateMemory(ctx context.Context, boardUuid, userUuid, imageUrl string) (*Memory, error)
	GetMemoryByUuid(ctx context.Context, uuid string) (*Memory, error)
	FindAllByBoardUuid(ctx context.Context, boardUuid string) ([]*Memory, error)
	DeleteMemoryByUuid(ctx context.Context, uuid string) error
}

type repositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &repositoryImpl{queries: queries}
}

func (r *repositoryImpl) CreateMemory(ctx context.Context, boardUuid, userUuid, imageUrl string) (*Memory, error) {
	bId, err := utility.UuidFromString(boardUuid)
	if err != nil {
		return nil, err
	}
	uId, err := utility.UuidFromString(userUuid)
	if err != nil {
		return nil, err
	}

	entity, err := r.queries.CreateMemory(ctx, db.CreateMemoryParams{
		BoardUuid:  bId,
		UploadedBy: uId,
		ImageUrl:   imageUrl,
	})
	if err != nil {
		return nil, err
	}
	return FromEntity(entity), nil
}

func (r *repositoryImpl) GetMemoryByUuid(ctx context.Context, uuid string) (*Memory, error) {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		return nil, err
	}
	entity, err := r.queries.GetMemoryByUuid(ctx, id)
	if err != nil {
		return nil, err
	}
	return FromEntity(entity), nil
}

func (r *repositoryImpl) FindAllByBoardUuid(ctx context.Context, boardUuid string) ([]*Memory, error) {
	id, err := utility.UuidFromString(boardUuid)
	if err != nil {
		return nil, err
	}
	rows, err := r.queries.FindAllMemoriesByBoardUuid(ctx, id)
	if err != nil {
		return nil, err
	}
	out := make([]*Memory, 0, len(rows))
	for _, row := range rows {
		out = append(out, FromEntity(row))
	}
	return out, nil
}

func (r *repositoryImpl) DeleteMemoryByUuid(ctx context.Context, uuid string) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		return err
	}
	return r.queries.DeleteMemoryByUuid(ctx, id)
}
