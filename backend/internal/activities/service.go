package activities

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Service interface {
	CreateActivity(ctx context.Context, url string) (*Activity, error)
	GetActivityById(ctx context.Context, uuid string) (*AggregatedActivity, error)
	GetActivitiesByBoardId(ctx context.Context, uuid string) ([]*Activity, error)
	UpdateActivityById(ctx context.Context, uuid string, data *Activity) error
	DeleteActivityById(ctx context.Context, uuid string) error
}

type accomodationServiceImpl struct {
	repo Repository

	commentsSvc comments.Service
	votesSvc    votes.Service
}

func NewService(repo Repository, commentsSvc comments.Service, votesSvc votes.Service) Service {
	return &accomodationServiceImpl{
		repo:        repo,
		commentsSvc: commentsSvc,
		votesSvc:    votesSvc,
	}
}

func (svc *accomodationServiceImpl) CreateActivity(ctx context.Context, url string) (*Activity, error) {
	// TODO: Scrape the data
	data := Activity{}

	return svc.repo.CreateActivity(ctx, &data)
}

func (svc *accomodationServiceImpl) GetActivityById(ctx context.Context, uuid string) (*AggregatedActivity, error) {
	accomodation, err := svc.repo.GetActivityById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	comms, err := svc.commentsSvc.FindAllCommentsByRelatedEntity(ctx, db.CommentedOnActivities, uuid)
	if err != nil {
		return nil, err
	}

	votes, err := svc.votesSvc.FindAllVotesByRelatedEntity(ctx, db.VotedOnActivities, uuid)
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
