package scrapeaudit

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
)

type Service interface {
	CreateScrapeAudit(ctx context.Context, url, host string) (*ScrapeAudit, error)
	UpdateScrapeAuditByUuid(ctx context.Context, uuid string, status db.ScrapeStatus, data *ScrapeResult) error
}

type scrapeAuditServiceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &scrapeAuditServiceImpl{
		repo: repo,
	}
}

func (svc *scrapeAuditServiceImpl) CreateScrapeAudit(ctx context.Context, url, host string) (*ScrapeAudit, error) {
	return svc.repo.CreateScrapeAudit(ctx, url, host)
}

func (svc *scrapeAuditServiceImpl) UpdateScrapeAuditByUuid(ctx context.Context, uuid string, status db.ScrapeStatus, data *ScrapeResult) error {
	return svc.repo.UpdateScrapeAuditByUuid(ctx, uuid, status, data)
}
