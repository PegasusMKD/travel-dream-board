package activities

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	scrapeprocess "github.com/PegasusMKD/travel-dream-board/internal/scrape_process"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Service interface {
	CreateActivity(ctx context.Context, url string, imageBytes []byte, imageExt string, boardUuid, userUuid string) (*Activity, error)
	GetActivityById(ctx context.Context, uuid string) (*AggregatedActivity, error)
	GetActivitiesByBoardId(ctx context.Context, uuid string) ([]*Activity, error)
	UpdateActivityById(ctx context.Context, uuid string, data *Activity) error
	DeleteActivityById(ctx context.Context, uuid string) error
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

func (svc *accomodationServiceImpl) CreateActivity(ctx context.Context, url string, imageBytes []byte, imageExt string, boardUuid, userUuid string) (*Activity, error) {
	data := Activity{
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

			data.StartAt = utility.ParseWallClockTime(extractedData.StartAt)
			data.EndAt = utility.ParseWallClockTime(extractedData.EndAt)
			data.Location = extractedData.Location

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

		data.StartAt = extractedData.StartAt
		data.EndAt = extractedData.EndAt
		data.Location = extractedData.Location

		data.Price = extractedData.Price
		data.Currency = extractedData.Currency
		data.Description = extractedData.Description
	}

	return svc.repo.CreateActivity(ctx, &data)
}

func (svc *accomodationServiceImpl) GetActivityById(ctx context.Context, uuid string) (*AggregatedActivity, error) {
	accomodation, err := svc.repo.GetActivityById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	comms, err := svc.commentsService.FindAllCommentsByRelatedEntity(ctx, db.CommentedOnActivities, uuid)
	if err != nil {
		return nil, err
	}

	votes, err := svc.votesService.FindAllVotesByRelatedEntity(ctx, db.VotedOnActivities, uuid)
	if err != nil {
		return nil, err
	}

	return &AggregatedActivity{
		Activity: *accomodation,
		Comments: comms,
		Votes:    votes,
	}, nil
}

func (svc *accomodationServiceImpl) GetActivitiesByBoardId(ctx context.Context, uuid string) ([]*Activity, error) {
	return svc.repo.GetAllActivitiesByBoardId(ctx, uuid)
}

func (svc *accomodationServiceImpl) UpdateActivityById(ctx context.Context, uuid string, data *Activity) error {
	return svc.repo.UpdateActivityById(ctx, uuid, data)
}

func (svc *accomodationServiceImpl) DeleteActivityById(ctx context.Context, uuid string) error {
	return svc.repo.DeleteActivityById(ctx, uuid)
}
