package transport

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	scrapeprocess "github.com/PegasusMKD/travel-dream-board/internal/scrape_process"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Service interface {
	CreateTransport(ctx context.Context, url string, imageBytes []byte, imageExt string, boardUuid, userUuid string) (*Transport, error)
	GetTransportById(ctx context.Context, uuid string) (*AggregatedTransport, error)
	GetTransportsByBoardId(ctx context.Context, uuid string) ([]*Transport, error)
	UpdateTransportById(ctx context.Context, uuid string, data *Transport) error
	DeleteTransportById(ctx context.Context, uuid string) error
}

type accomodationServiceImpl struct {
	repo Repository

	commentsService comments.Service
	votesService    votes.Service

	scrapeService scrapeprocess.Service
}

func NewService(repo Repository, commentsService comments.Service, votesService votes.Service, scrapeService scrapeprocess.Service) Service {
	return &accomodationServiceImpl{
		repo: repo,

		commentsService: commentsService,
		votesService:    votesService,

		scrapeService: scrapeService,
	}
}

func (svc *accomodationServiceImpl) CreateTransport(ctx context.Context, url string, imageBytes []byte, imageExt string, boardUuid, userUuid string) (*Transport, error) {
	data := Transport{
		UserUuid:  userUuid,
		BoardUuid: boardUuid,
	}

	if len(imageBytes) > 0 {
		extractedData, err := svc.scrapeService.ExtractFromImage(ctx, url, imageBytes)
		if err != nil {
			return nil, err
		}

		if extractedData != nil {
			data.Title = extractedData.Title
			data.Url = url
			data.ImageUrl = &url

			data.OutboundDepartingLocation = extractedData.OutboundDepartingLocation
			data.OutboundArrivingLocation = extractedData.OutboundArrivingLocation
			data.OutboundDepartingAt = utility.ParseWallClockTime(extractedData.OutboundDepartingAt)
			data.OutboundArrivingAt = utility.ParseWallClockTime(extractedData.OutboundArrivingAt)
			data.InboundDepartingLocation = extractedData.InboundDepartingLocation
			data.InboundArrivingLocation = extractedData.InboundArrivingLocation
			data.InboundDepartingAt = utility.ParseWallClockTime(extractedData.InboundDepartingAt)
			data.InboundArrivingAt = utility.ParseWallClockTime(extractedData.InboundArrivingAt)

			data.OutboundDurationMinutes = utility.ParseDurationMinutes(extractedData.OutboundDurationMinutes)
			data.InboundDurationMinutes = utility.ParseDurationMinutes(extractedData.InboundDurationMinutes)

			data.Price = extractedData.Price
			data.Currency = utility.ParseCurrencyCode(extractedData.Currency)
			if extractedData.Description != "" {
				desc := extractedData.Description
				data.Description = &desc
			}
		} else {
			data.Title = "Uploaded Image"
			data.Url = url
			data.ImageUrl = &url
		}
	} else {
		extractedData, err := svc.scrapeService.Scrape(ctx, url)
		if err != nil {
			return nil, err
		}
		data.Url = url
		data.Title = *extractedData.Title
		data.ImageUrl = extractedData.ImageUrl

		data.OutboundDepartingLocation = extractedData.OutboundDepartingLocation
		data.OutboundArrivingLocation = extractedData.OutboundArrivingLocation
		data.OutboundDepartingAt = extractedData.OutboundDepartingAt
		data.OutboundArrivingAt = extractedData.OutboundArrivingAt
		data.InboundDepartingLocation = extractedData.InboundDepartingLocation
		data.InboundArrivingLocation = extractedData.InboundArrivingLocation
		data.InboundDepartingAt = extractedData.InboundDepartingAt
		data.InboundArrivingAt = extractedData.InboundArrivingAt

		data.OutboundDurationMinutes = extractedData.OutboundDurationMinutes
		data.InboundDurationMinutes = extractedData.InboundDurationMinutes

		data.Price = extractedData.Price
		data.Currency = extractedData.Currency
		data.Description = extractedData.Description
	}

	return svc.repo.CreateTransport(ctx, &data)
}

func (svc *accomodationServiceImpl) GetTransportById(ctx context.Context, uuid string) (*AggregatedTransport, error) {
	accomodation, err := svc.repo.GetTransportById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	comms, err := svc.commentsService.FindAllCommentsByRelatedEntity(ctx, db.CommentedOnTransport, uuid)
	if err != nil {
		return nil, err
	}

	votes, err := svc.votesService.FindAllVotesByRelatedEntity(ctx, db.VotedOnTransport, uuid)
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
