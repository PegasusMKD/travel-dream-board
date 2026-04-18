package activities

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	scrapeprocess "github.com/PegasusMKD/travel-dream-board/internal/scrape_process"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Service interface {
	CreateActivity(ctx context.Context, url, boardUuid, userUuid string) (*Activity, error)
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

func (svc *accomodationServiceImpl) CreateActivity(ctx context.Context, url, boardUuid, userUuid string) (*Activity, error) {
	extractedData, err := svc.scrapeService.Scrape(ctx, url)
	if err != nil {
		return nil, err
	}

	//
	data := Activity{
		Url:      url,
		Title:    *extractedData.Title,
		ImageUrl: extractedData.ImageUrl,

		UserUuid:  userUuid,
		BoardUuid: boardUuid,
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
