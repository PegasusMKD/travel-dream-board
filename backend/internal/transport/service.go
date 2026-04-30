package transport

import (
	"context"
	"time"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	scrapeprocess "github.com/PegasusMKD/travel-dream-board/internal/scrape_process"
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
			data.OutboundDepartingAt = parseLegTimeStr(extractedData.OutboundDepartingAt)
			data.OutboundArrivingAt = parseLegTimeStr(extractedData.OutboundArrivingAt)
			data.InboundDepartingLocation = extractedData.InboundDepartingLocation
			data.InboundArrivingLocation = extractedData.InboundArrivingLocation
			data.InboundDepartingAt = parseLegTimeStr(extractedData.InboundDepartingAt)
			data.InboundArrivingAt = parseLegTimeStr(extractedData.InboundArrivingAt)
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
	}

	return svc.repo.CreateTransport(ctx, &data)
}

func parseLegTimeStr(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	layouts := []string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04:05", "2006-01-02 15:04:05Z07:00", "2006-01-02 15:04:05"}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, *s); err == nil {
			// Wall-clock: keep Y/M/D H:M:S components, re-anchor as UTC so the
			// value survives timezone conversions unchanged (14:30 stays 14:30).
			wall := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.UTC)
			return &wall
		}
	}
	return nil
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
