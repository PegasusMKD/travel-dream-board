package accomodations

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	scrapeprocess "github.com/PegasusMKD/travel-dream-board/internal/scrape_process"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Service interface {
	CreateAccomodation(ctx context.Context, url, boardUuid, userUuid string) (*Accomodation, error)
	GetAccomodationById(ctx context.Context, uuid string) (*AggregatedAccomodation, error)
	GetAccomodationsByBoardId(ctx context.Context, uuid string) ([]*Accomodation, error)
	UpdateAccomodationById(ctx context.Context, uuid string, data *Accomodation) error
	DeleteAccomodationById(ctx context.Context, uuid string) error
}

type accomodationServiceImpl struct {
	repo Repository

	commentsService comments.Service
	votesService    votes.Service

	scrapeService scrapeprocess.Service
}

func NewService(repo Repository, commentsService comments.Service, votesService votes.Service, scrapeService scrapeprocess.Service) Service {
	return &accomodationServiceImpl{
		repo:            repo,
		commentsService: commentsService,
		votesService:    votesService,
		scrapeService:   scrapeService,
	}
}

func (svc *accomodationServiceImpl) CreateAccomodation(ctx context.Context, url, boardUuid, userUuid string) (*Accomodation, error) {
	extractedData, err := svc.scrapeService.Scrape(ctx, url)
	if err != nil {
		return nil, err
	}

	//
	data := Accomodation{
		Url:      url,
		Title:    *extractedData.Title,
		ImageUrl: extractedData.ImageUrl,

		UserUuid:  userUuid,
		BoardUuid: boardUuid,
	}

	return svc.repo.CreateAccomodation(ctx, &data)
}

func (svc *accomodationServiceImpl) GetAccomodationById(ctx context.Context, uuid string) (*AggregatedAccomodation, error) {
	accomodation, err := svc.repo.GetAccomodationById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	comms, err := svc.commentsService.FindAllCommentsByRelatedEntity(ctx, db.CommentedOnAccomodation, uuid)
	if err != nil {
		return nil, err
	}

	votes, err := svc.votesService.FindAllVotesByRelatedEntity(ctx, db.VotedOnAccomodation, uuid)
	if err != nil {
		return nil, err
	}

	return &AggregatedAccomodation{
		Accomodation: *accomodation,
		Comments:     comms,
		Votes:        votes,
	}, nil
}

func (svc *accomodationServiceImpl) GetAccomodationsByBoardId(ctx context.Context, uuid string) ([]*Accomodation, error) {
	return svc.repo.GetAllAccomodationsByBoardId(ctx, uuid)
}

func (svc *accomodationServiceImpl) UpdateAccomodationById(ctx context.Context, uuid string, data *Accomodation) error {
	return svc.repo.UpdateAccomodationById(ctx, uuid, data)
}

func (svc *accomodationServiceImpl) DeleteAccomodationById(ctx context.Context, uuid string) error {
	return svc.repo.DeleteAccomodationById(ctx, uuid)
}
