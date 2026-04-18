package accomodations

import "context"

type Service interface {
	CreateAccomodation(ctx context.Context, url string) (*Accomodation, error)
	GetAccomodationById(ctx context.Context, uuid string) (*AggregatedAccomodation, error)
	GetAccomodationsByBoardId(ctx context.Context, uuid string) ([]*Accomodation, error)
	UpdateAccomodationById(ctx context.Context, uuid string, data *Accomodation) error
	DeleteAccomodationById(ctx context.Context, uuid string) error
}

type accomodationServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &accomodationServiceImpl{
		repo: repo,
	}
}

func (svc *accomodationServiceImpl) CreateAccomodation(ctx context.Context, url string) (*Accomodation, error) {
	// TODO: Scrape the data
	data := Accomodation{}

	return svc.repo.CreateAccomodation(ctx, &data)
}

func (svc *accomodationServiceImpl) GetAccomodationById(ctx context.Context, uuid string) (*AggregatedAccomodation, error) {
	accomodation, err := svc.repo.GetAccomodationById(ctx, uuid)
	return nil, nil
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
