package activities

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
)

type Repository interface {
	CreateActivity(ctx context.Context, data *Activity) (*Activity, error)
	GetActivityById(ctx context.Context, uuid string) (*Activity, error)
	GetAllActivitiesByBoardId(ctx context.Context, uuid string) ([]*Activity, error)
	UpdateActivityById(ctx context.Context, uuid string, data *Activity) error
	DeleteActivityById(ctx context.Context, uuid string) error
}

type repositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &repositoryImpl{
		queries: queries,
	}
}

func (repo *repositoryImpl) CreateActivity(ctx context.Context, data *Activity) (*Activity, error) {
	boardUuid, err := utility.UuidFromString(data.BoardUuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", data.BoardUuid, "error", err)
		return nil, err
	}

	userUuid, err := utility.UuidFromString(data.UserUuid)
	if err != nil {
		log.Error("Failed parsing UserUUID", "uuid", data.UserUuid, "error", err)
		return nil, err
	}

	params := db.CreateActivityParams{
		Url:       data.Url,
		Title:     data.Title,
		ImageUrl:  data.ImageUrl,
		BoardUuid: boardUuid,
		UserUuid:  userUuid,
		StartAt:   utility.TimestamptzFromTime(data.StartAt),
		EndAt:     utility.TimestamptzFromTime(data.EndAt),
		Location:  data.Location,
		Price:       utility.NumericFromString(data.Price),
		Currency:    utility.NullCurrencyFromPtr(data.Currency),
		Description: data.Description,
	}

	entity, err := repo.queries.CreateActivity(ctx, params)
	if err != nil {
		log.Error("Failed creating entity", "url", data.Url, "error", err)
		return nil, err
	}

	return FromEntity(entity), nil
}

func (repo *repositoryImpl) GetActivityById(ctx context.Context, uuid string) (*Activity, error) {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return nil, err
	}

	ent, err := repo.queries.GetActivityByUuid(ctx, id)
	if err != nil {
		log.Error("Failed fetching an entity", "uuid", uuid, "error", err)
		return nil, err
	}

	return FromGetActivityRow(ent), nil
}

func (repo *repositoryImpl) GetAllActivitiesByBoardId(ctx context.Context, uuid string) ([]*Activity, error) {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return nil, err
	}

	ent, err := repo.queries.FindAllActivitiesByBoardUuid(ctx, id)
	if err != nil {
		log.Error("Failed fetching an entity", "uuid", uuid, "error", err)
		return nil, err
	}

	data := []*Activity{}
	for _, val := range ent {
		data = append(data, FromFindActivitiesRow(val))
	}

	return data, nil
}

func (repo *repositoryImpl) UpdateActivityById(ctx context.Context, uuid string, data *Activity) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return err
	}

	params := db.UpdateActivityByUuidParams{
		Url:              data.Url,
		Title:            data.Title,
		ImageUrl:         data.ImageUrl,
		Notes:            data.Notes,
		Status:           data.Status,
		BookingReference: data.BookingReference,
		Selected:         data.Selected,
		Uuid:             id,
		StartAt:          utility.TimestamptzFromTime(data.StartAt),
		EndAt:            utility.TimestamptzFromTime(data.EndAt),
		Location:         data.Location,
		Price:            utility.NumericFromString(data.Price),
		Currency:         utility.NullCurrencyFromPtr(data.Currency),
		Description:      data.Description,
	}

	return repo.queries.UpdateActivityByUuid(ctx, params)
}

func (repo *repositoryImpl) DeleteActivityById(ctx context.Context, uuid string) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return err
	}

	return repo.queries.DeleteActivityByUuid(ctx, id)
}
