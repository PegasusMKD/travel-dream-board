package scrapeaudit

import (
	"context"

	"github.com/PegasusMKD/travel-dream-board/internal/db"
	"github.com/PegasusMKD/travel-dream-board/internal/utility"
)

type Repository interface {
	CreateScrapeAudit(ctx context.Context, url, host string) (*ScrapeAudit, error)
	UpdateScrapeAuditByUuid(ctx context.Context, uuid string, status db.ScrapeStatus, data *ScrapeResult) error
}

type scrapeAuditRepositoryImpl struct {
	queries *db.Queries
}

func NewRepository(queries *db.Queries) Repository {
	return &scrapeAuditRepositoryImpl{
		queries: queries,
	}
}

func (repo *scrapeAuditRepositoryImpl) CreateScrapeAudit(ctx context.Context, url, host string) (*ScrapeAudit, error) {
	params := db.CreateScrapeAuditParams{
		Url:  url,
		Host: host,
	}

	ent, err := repo.queries.CreateScrapeAudit(ctx, params)
	if err != nil {
		log.Error("Failed to create scrape audit", "url", url, "host", host, "error", err)
	}

	return FromEntity(&ent), nil
}

func (repo *scrapeAuditRepositoryImpl) UpdateScrapeAuditByUuid(ctx context.Context, uuid string, status db.ScrapeStatus, data *ScrapeResult) error {
	id, err := utility.UuidFromString(uuid)
	if err != nil {
		log.Error("Failed parsing UUID", "uuid", uuid, "error", err)
		return err
	}

	params := db.UpdateScrapeAuditByUuidParams{
		Uuid:        id,
		Status:      status,
		Title:       data.Title,
		ImageUrl:    data.ImageUrl,
		Description: data.Description,
		SiteName:    data.SiteName,
	}

	return repo.queries.UpdateScrapeAuditByUuid(ctx, params)
}
