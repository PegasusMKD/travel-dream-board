package scrapeaudit

import "github.com/PegasusMKD/travel-dream-board/internal/db"

type ScrapeAudit struct {
	Uuid string

	Url    string
	Host   string
	Status db.ScrapeStatus

	Title       *string
	Description *string
	ImageUrl    *string
	SiteName    *string
}

func FromEntity(entity *db.ScrapeAudit) *ScrapeAudit {
	return &ScrapeAudit{
		Uuid: entity.Uuid.String(),

		Url:    entity.Url,
		Host:   entity.Host,
		Status: entity.Status,

		Title:       entity.Title,
		Description: entity.Description,
		ImageUrl:    entity.ImageUrl,
		SiteName:    entity.SiteName,
	}
}

type ScrapeResult struct {
	InitialUrl string
	ActualUrl  string

	Title       *string
	Description *string
	ImageUrl    *string
	SiteName    *string
}
