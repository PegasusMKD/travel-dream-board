package accomodations

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/comments"
	"github.com/PegasusMKD/travel-dream-board/internal/db"
	scrapeprocess "github.com/PegasusMKD/travel-dream-board/internal/scrape_process"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
	"github.com/PegasusMKD/travel-dream-board/internal/votes"
)

type Service interface {
	CreateAccomodation(ctx context.Context, url string, imageBytes []byte, imageExt string, boardUuid, userUuid string) (*Accomodation, error)
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

func (svc *accomodationServiceImpl) CreateAccomodation(ctx context.Context, url string, imageBytes []byte, imageExt, boardUuid, userUuid string) (*Accomodation, error) {
	data := Accomodation{
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
			data.Url = url // In handler we will pass the local uploaded path as `url`
			// Set the ImageUrl to be the local url since we uploaded it
			data.ImageUrl = &url

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

		data.Price = extractedData.Price
		data.Currency = extractedData.Currency
		data.Description = extractedData.Description
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
