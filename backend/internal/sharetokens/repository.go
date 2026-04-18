package sharetokens

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
)

type Repository interface {
	CreateShareToken(ctx context.Context, token string, boardUuid string) (*ShareToken, error)
	GetShareToken(ctx context.Context, token string) (*ShareToken, error)
	DeleteShareToken(ctx context.Context, token string, boardUuid string) error
	GetShareTokensForBoard(ctx context.Context, boardUuid string) ([]*ShareToken, error)
}

type repositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &repositoryImpl{
		queries: queries,
	}
}

func (repo *repositoryImpl) CreateShareToken(ctx context.Context, token string, boardUuid string) (*ShareToken, error) {
	bId, err := utility.UuidFromString(boardUuid)
	if err != nil {
		return nil, err
	}

	params := db.CreateShareTokenParams{
		Token:     token,
		BoardUuid: bId,
	}

	ent, err := repo.queries.CreateShareToken(ctx, params)
	if err != nil {
		return nil, err
	}

	return FromEntity(&ent), nil
}

func (repo *repositoryImpl) GetShareToken(ctx context.Context, token string) (*ShareToken, error) {
	ent, err := repo.queries.GetShareToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return FromEntity(&ent), nil
}

func (repo *repositoryImpl) DeleteShareToken(ctx context.Context, token string, boardUuid string) error {
	bId, err := utility.UuidFromString(boardUuid)
	if err != nil {
		return err
	}

	params := db.DeleteShareTokenParams{
		Token:     token,
		BoardUuid: bId,
	}

	return repo.queries.DeleteShareToken(ctx, params)
}

func (repo *repositoryImpl) GetShareTokensForBoard(ctx context.Context, boardUuid string) ([]*ShareToken, error) {
	bId, err := utility.UuidFromString(boardUuid)
	if err != nil {
		return nil, err
	}

	ents, err := repo.queries.GetShareTokensForBoard(ctx, bId)
	if err != nil {
		return nil, err
	}

	data := make([]*ShareToken, 0, len(ents))
	for _, val := range ents {
		data = append(data, FromEntity(&val))
	}

	return data, nil
}
