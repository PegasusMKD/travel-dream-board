package transport

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
)

type Repository interface {
	CreateTransport(ctx context.Context, data *Transport) (*Transport, error)
	GetTransportById(ctx context.Context, uuid string) (*Transport, error)
	GetAllTransportsByBoardId(ctx context.Context, uuid string) ([]*Transport, error)
	UpdateTransportById(ctx context.Context, uuid string, data *Transport) error
	DeleteTransportById(ctx context.Context, uuid string) error
}

type repositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &repositoryImpl{
		queries: queries,
	}
}

func (repo *repositoryImpl) CreateTransport(ctx context.Context, data *Transport) (*Transport, error) {
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

	params := db.CreateTransportParams{
		Url:       data.Url,
		Title:     data.Title,
		ImageUrl:  data.ImageUrl,
		BoardUuid: boardUuid,
		UserUuid:  userUuid,

		OutboundDepartingLocation: data.OutboundDepartingLocation,
		OutboundArrivingLocation:  data.OutboundArrivingLocation,
		OutboundDepartingAt:       utility.TimestamptzFromTime(data.OutboundDepartingAt),
		OutboundArrivingAt:        utility.TimestamptzFromTime(data.OutboundArrivingAt),
		InboundDepartingLocation:  data.InboundDepartingLocation,
		InboundArrivingLocation:   data.InboundArrivingLocation,
		InboundDepartingAt:        utility.TimestamptzFromTime(data.InboundDepartingAt),
		InboundArrivingAt:         utility.TimestamptzFromTime(data.InboundArrivingAt),
	}

	entity, err := repo.queries.CreateTransport(ctx, params)
	if err != nil {
		log.Error("Failed creating entity", "url", data.Url, "error", err)
		return nil, err
	}

	return FromEntity(entity), nil
}

func (repo *repositoryImpl) GetTransportById(ctx context.Context, uuid string) (*Transport, error) {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return nil, err
	}

	ent, err := repo.queries.GetTransportByUuid(ctx, id)
	if err != nil {
		log.Error("Failed fetching an entity", "uuid", uuid, "error", err)
		return nil, err
	}

	return FromGetTransportRow(ent), nil
}

func (repo *repositoryImpl) GetAllTransportsByBoardId(ctx context.Context, uuid string) ([]*Transport, error) {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return nil, err
	}

	ent, err := repo.queries.FindAllTransportByBoardUuid(ctx, id)
	if err != nil {
		log.Error("Failed fetching an entity", "uuid", uuid, "error", err)
		return nil, err
	}

	data := []*Transport{}
	for _, val := range ent {
		data = append(data, FromFindTransportRow(val))
	}

	return data, nil
}

func (repo *repositoryImpl) UpdateTransportById(ctx context.Context, uuid string, data *Transport) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing passed UUID", "uuid", uuid, "error", err)
		return err
	}

	params := db.UpdateTransportByUuidParams{
		Url:              data.Url,
		Title:            data.Title,
		ImageUrl:         data.ImageUrl,
		Notes:            data.Notes,
		Status:           data.Status,
		BookingReference: data.BookingReference,
		Selected:         data.Selected,
		Uuid:             id,

		OutboundDepartingLocation: data.OutboundDepartingLocation,
		OutboundArrivingLocation:  data.OutboundArrivingLocation,
		OutboundDepartingAt:       utility.TimestamptzFromTime(data.OutboundDepartingAt),
		OutboundArrivingAt:        utility.TimestamptzFromTime(data.OutboundArrivingAt),
		InboundDepartingLocation:  data.InboundDepartingLocation,
		InboundArrivingLocation:   data.InboundArrivingLocation,
		InboundDepartingAt:        utility.TimestamptzFromTime(data.InboundDepartingAt),
		InboundArrivingAt:         utility.TimestamptzFromTime(data.InboundArrivingAt),
	}

	return repo.queries.UpdateTransportByUuid(ctx, params)
}

func (repo *repositoryImpl) DeleteTransportById(ctx context.Context, uuid string) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return err
	}

	return repo.queries.DeleteTransportByUuid(ctx, id)
}
