package boards

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
	"github.com/jackc/pgx/v5/pgtype"
)

type Repository interface {
	CreateBoard(ctx context.Context, data *Board) (*Board, error)
	GetBoardById(ctx context.Context, uuid string) (*Board, error)
	GetAllBoards(ctx context.Context, userUuid string) ([]*Board, error)
	UpdateBoardById(ctx context.Context, uuid string, data *Board) error
	DeleteBoardById(ctx context.Context, uuid string) error
	GetShareToken(ctx context.Context, token string) (*db.ShareToken, error)
}

type repositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &repositoryImpl{
		queries: queries,
	}
}

func (repo *repositoryImpl) CreateBoard(ctx context.Context, data *Board) (*Board, error) {
	var userUuid pgtype.UUID
	if data.UserUuid != nil {
		err := userUuid.Scan(*data.UserUuid)
		if err != nil {
			log.Error("Failed parsing UserUuid", "uuid", *data.UserUuid, "error", err)
			return nil, err
		}
	}

	params := db.CreateBoardParams{
		Name:         data.Name,
		LocationName: data.LocationName,
		UserUuid:     userUuid,
	}

	ent, err := repo.queries.CreateBoard(ctx, params)
	if err != nil {
		log.Error("Failed creating board", "name", data.Name, "location_name", data.LocationName, "error", err)
		return nil, err
	}

	board := FromEntity(&ent)
	return board, nil
}

func (repo *repositoryImpl) GetBoardById(ctx context.Context, uuid string) (*Board, error) {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return nil, err
	}

	ent, err := repo.queries.GetBoardById(ctx, id)
	if err != nil {
		log.Error("Failed fetching an entity", "uuid", uuid, "error", err)
		return nil, err
	}

	return FromEntity(&ent), nil
}

func (repo *repositoryImpl) GetAllBoards(ctx context.Context, userUuid string) ([]*Board, error) {
	id, err := utility.UuidFromString(userUuid)
	if err != nil {
		log.Error("Failed parsing provided id", "uuid", userUuid, "error", err)
		return nil, err
	}

	ents, err := repo.queries.GetAllBoards(ctx, id)
	if err != nil {
		log.Error("Failed fetching boards", "error", err)
		return nil, err
	}

	data := []*Board{}
	for _, val := range ents {
		data = append(data, FromEntity(&val))
	}

	return data, nil
}

func (repo *repositoryImpl) UpdateBoardById(ctx context.Context, uuid string, data *Board) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing provided id", "uuid", uuid, "error", err)
		return err
	}

	params := db.UpdateBoardByIdParams{
		Name:         data.Name,
		Details:      data.Details,
		LocationName: data.LocationName,
		StartsAt:     utility.DateFromTime(data.StartsAt),
		LastsUntil:   utility.DateFromTime(data.LastsUntil),
		ThumbnailUrl: data.ThumbnailUrl,
		Uuid:         id,
	}

	return repo.queries.UpdateBoardById(ctx, params)
}

func (repo *repositoryImpl) DeleteBoardById(ctx context.Context, uuid string) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing provided id", "uuid", uuid, "error", err)
		return err
	}

	return repo.queries.DeleteBoardById(ctx, id)
}

func (repo *repositoryImpl) GetShareToken(ctx context.Context, token string) (*db.ShareToken, error) {
	ent, err := repo.queries.GetShareToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return &ent, nil
}
